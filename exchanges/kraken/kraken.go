package kraken

import (
    "crypto/hmac"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/base64"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/shopspring/decimal"
)

// docs: https://docs.kraken.com/rest/
const RESTEndpoint string = "https://api.kraken.com"

type Kraken struct {
    AssetPairTranslator      types.AssetPairTranslator
    ISO4217Translator        types.AssetPairTranslator

    apiKey                   string
    secretKey                string
    spreadRecorder           types.SpreadRecorder
    orderBookRecorder        types.OrderBookRecorder
    latencyEstimator         *util.EwmaEstimator
    orderIdToOrderTranslator *util.ConcurrentOrderIdToOrderPtrMap
    // add timeouts
    httpClient               *http.Client
}

func NewKraken(apiKey, secretKey string, assetPairTranslator types.AssetPairTranslator, iso4217Translator types.AssetPairTranslator) *Kraken {
    assetPairs := assetPairTranslator.GetAssetPairs()
    return &Kraken{
        AssetPairTranslator: assetPairTranslator,
        apiKey: apiKey,
        secretKey: secretKey,
        spreadRecorder: NewKrakenSpreadRecorder(assetPairs, iso4217Translator, 200),
        orderBookRecorder: NewKrakenOrderBookRecorder(assetPairs, iso4217Translator, 1000),
        latencyEstimator: util.NewEwmaEstimator(0.125, 0.25, 4),
        orderIdToOrderTranslator: util.NewConcurrentOrderIdToOrderPtrMap(),
        httpClient: &http.Client{},
    }
}

func (k *Kraken) String() string {
    return "Kraken"
}

func checkError(bodyJson map[string]interface{}) {
    errors := bodyJson["error"].([]interface{})
    if len(errors) > 0 {
        log.Fatalln(errors)
    }
}

func (k *Kraken) getHistoricalSpread(assetPair types.AssetPair, duration time.Duration, samples uint, channel chan types.SpreadResponse) {
    if samples == 0 || duration <= 0 {
        channel <- types.SpreadResponse{assetPair, []types.Spread{}}
    }

    rawHistoricalSpreads, ok := k.spreadRecorder.GetHistoricalSpreads(assetPair)
    if len(rawHistoricalSpreads) == 0 {
        if !ok {
            k.spreadRecorder.RegisterAssetPair(assetPair)
        }
        channel <- types.SpreadResponse{assetPair, rawHistoricalSpreads}
        return
    }

    channel <- types.SpreadResponse{assetPair, util.GetSpreadSamples(rawHistoricalSpreads, duration, samples)}
}

func (k *Kraken) GetHistoricalSpreads(assetPairs []types.AssetPair, duration time.Duration, samples uint) map[types.AssetPair][]types.Spread {
    channel := make(chan types.SpreadResponse)
    for _, assetPair := range assetPairs {
        go k.getHistoricalSpread(assetPair, duration, samples, channel)
    }

    historicalSpreads := make(map[types.AssetPair][]types.Spread)
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        historicalSpreads[response.AssetPair] = response.HistoricalSpreads
    }
    return historicalSpreads
}

func (k *Kraken) GetCurrentSpread(assetPair types.AssetPair) types.Spread {
    spread, ok := k.spreadRecorder.GetCurrentSpread(assetPair)
    if !ok {
        k.spreadRecorder.RegisterAssetPair(assetPair)
        bodyJson := util.HttpGetAndGetBody(k.httpClient, util.ParseUrlWithQuery(RESTEndpoint + "/0/public/Ticker", url.Values{
            "pair": []string{k.AssetPairTranslator[assetPair]},
        }))
        checkError(bodyJson)
    
        data := bodyJson["result"].(map[string]interface{})[k.AssetPairTranslator[assetPair]].(map[string]interface{})
    
        bid, err := decimal.NewFromString(data["b"].([]interface{})[0].(string))
        if err != nil {
            log.Fatalln(err)
        }
        ask, err := decimal.NewFromString(data["a"].([]interface{})[0].(string))
        if err != nil {
            log.Fatalln(err)
        }
    
        return types.Spread{
            Bid: bid,
            Ask: ask,
            Timestamp: time.Now(),
        }
    }

    return spread
}

func (k *Kraken) getOrderBook(assetPair types.AssetPair, channel chan types.OrderBookResponse) {
    orderBook, ok := k.orderBookRecorder.GetOrderBook(assetPair)
    if !ok {
        k.orderBookRecorder.RegisterAssetPair(assetPair)
    }
    channel <- types.OrderBookResponse{assetPair, &orderBook}
}

func (k *Kraken) GetOrderBooks(assetPairs []types.AssetPair) map[types.AssetPair]*types.OrderBook {
    channel := make(chan types.OrderBookResponse)
    for _, assetPair := range assetPairs {
        go k.getOrderBook(assetPair, channel)
    }

    orderBooks := make(map[types.AssetPair]*types.OrderBook)
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        orderBooks[response.AssetPair] = response.OrderBook
    }
    return orderBooks
}

func (k *Kraken) GetLatency() time.Duration {
    start := time.Now()

    bodyJson := util.HttpGetAndGetBody(k.httpClient, RESTEndpoint + "/0/public/Time")
    checkError(bodyJson)

    duration := time.Since(start)

    k.latencyEstimator.Sample(float64(duration.Milliseconds()))

    return time.Duration(k.latencyEstimator.GetEstimate()) * time.Millisecond
}

func parseOrderType(ot types.OrderType) string {
    if (ot == types.Buy) {
        return "buy"
    }
    return "sell"
}

func (k *Kraken) getKrakenSignature(urlPath string, values url.Values) string {
    b64DecodedSecret, err := base64.StdEncoding.DecodeString(k.secretKey)
    if err != nil {
        log.Fatalln(err)
    }

    sha := sha256.New()
    sha.Write([]byte(values.Get("nonce") + values.Encode()))
    shasum := sha.Sum(nil)

    mac := hmac.New(sha512.New, b64DecodedSecret)
    mac.Write(append([]byte(urlPath), shasum...))
    macsum := mac.Sum(nil)
    return base64.StdEncoding.EncodeToString(macsum)
}

func (k *Kraken) executeOrder(order types.Order, channel chan types.OrderIdResponse) {
    queryParams := url.Values{
        "pair": []string{k.AssetPairTranslator[order.AssetPair]},
        "type": []string{parseOrderType(order.OrderType)},
        "ordertype": []string{"limit"},
        "price": []string{order.Price.String()},
        "volume": []string{order.Quantity.String()},
        "nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }
    request, err := http.NewRequest("POST", RESTEndpoint + "/0/private/AddOrder", strings.NewReader(queryParams.Encode()))
    if err != nil {
        log.Fatalln(err)
    }

    signature := k.getKrakenSignature("/0/private/AddOrder", queryParams)
    request.Header.Set("API-Sign", signature)
    request.Header.Set("API-Key", k.apiKey)

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    data := bodyJson["result"].(map[string]interface{})["txid"].([]interface{})
    id := data[0].(string)

    k.orderIdToOrderTranslator.Store(types.OrderId(id), &order)

    channel <- types.OrderIdResponse{order, types.OrderId(id)}
}

func (k *Kraken) ExecuteOrders(orders []types.Order) map[types.Order]types.OrderId {
    channel := make(chan types.OrderIdResponse)
    for _, order := range orders {
        go k.executeOrder(order, channel)
    }

    orderIds := make(map[types.Order]types.OrderId)
    for i := 0; i < len(orders); i++ {
        response := <- channel
        orderIds[response.Order] = response.OrderId
    }
    return orderIds
}

func (k *Kraken) GetOrderStatuses(orderIds []types.OrderId) map[types.OrderId]types.OrderStatus {
    if len(orderIds) == 0 {
        return make(map[types.OrderId]types.OrderStatus)
    }

    orderIdStrings := make([]string, len(orderIds))
    for i, orderId := range orderIds {
        orderIdStrings[i] = string(orderId)
    }
    queryParams := url.Values{
        "txid": []string{strings.Join(orderIdStrings, ",")},
        "nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }
    request, err := http.NewRequest("POST", RESTEndpoint + "/0/private/QueryOrders", strings.NewReader(queryParams.Encode()))
    if err != nil {
        log.Fatalln(err)
    }

    signature := k.getKrakenSignature("/0/private/QueryOrders", queryParams)
    request.Header.Set("API-Sign", signature)
    request.Header.Set("API-Key", k.apiKey)

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    data := bodyJson["result"].(map[types.OrderId]map[string]interface{})
    orderStatuses := make(map[types.OrderId]types.OrderStatus)
    for id, rawOrderData := range data {
        original, ok := k.orderIdToOrderTranslator.Load(id)
        if !ok {
            log.Fatalf("order with id %v not found\n", id)
        }
        orderStatus := types.OrderStatus{
            Original: original,
        }

        switch status := rawOrderData["status"].(string); status {
        case "pending":
            orderStatus.Status = types.Pending
        case "open":
            orderStatus.Status = types.Unfilled
        case "closed":
            price, err := decimal.NewFromString(rawOrderData["price"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            quantity, err := decimal.NewFromString(rawOrderData["vol_exec"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            orderStatus.Status = types.Filled
            orderStatus.FilledPrice = &price
            orderStatus.FilledQuantity = &quantity
            k.orderIdToOrderTranslator.Delete(id)
        case "canceled":
            orderStatus.Status = types.Canceled
            k.orderIdToOrderTranslator.Delete(id)
        case "expired":
            orderStatus.Status = types.Expired
            k.orderIdToOrderTranslator.Delete(id)
        }
        orderStatuses[id] = orderStatus
    }
    return orderStatuses
}

func (k *Kraken) CancelOrders(orderIds []types.OrderId) {
    if len(orderIds) == 0 {
        return
    }

    orderIdStrings := make([]string, len(orderIds))
    for i, orderId := range orderIds {
        orderIdStrings[i] = string(orderId)
    }
    queryParams := url.Values{
        "txid": []string{strings.Join(orderIdStrings, ",")},
        "nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }
    request, err := http.NewRequest("POST", RESTEndpoint + "/0/private/CancelOrder", strings.NewReader(queryParams.Encode()))
    if err != nil {
        log.Fatalln(err)
    }

    signature := k.getKrakenSignature("/0/private/CancelOrder", queryParams)
    request.Header.Set("API-Sign", signature)
    request.Header.Set("API-Key", k.apiKey)

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    for _, orderId := range orderIds {
        k.orderIdToOrderTranslator.Delete(orderId)
    }
}

func (k *Kraken) GetBalances() map[types.Asset]decimal.Decimal {
    queryParams := url.Values{
        "nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }
    request, err := http.NewRequest("POST", RESTEndpoint + "/0/private/Balance", strings.NewReader(queryParams.Encode()))
    if err != nil {
        log.Fatalln(err)
    }

    signature := k.getKrakenSignature("/0/private/Balance", queryParams)
    request.Header.Set("API-Sign", signature)
    request.Header.Set("API-Key", k.apiKey)

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    if _, ok := bodyJson["result"]; ok {
        data := bodyJson["result"].(map[types.Asset]string)

        balances := make(map[types.Asset]decimal.Decimal)
        for asset, balanceString := range data {
            balances[asset], err = decimal.NewFromString(balanceString)
            if err != nil {
                log.Fatalln(err)
            }
        }
        return balances
    }
    return make(map[types.Asset]decimal.Decimal)
}
