package binanceus

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/shopspring/decimal"
)

// docs: https://github.com/binance-us/binance-official-api-docs/blob/master/rest-api.md
const RESTEndpoint string = "https://api.binance.us"

type BinanceUS struct {
    AssetPairTranslator      types.AssetPairTranslator

    apiKey                   string
    secretKey                string
    spreadRecorder           types.SpreadRecorder
    orderBookRecorder        types.OrderBookRecorder
    latencyEstimator         *util.EwmaEstimator
    orderIdToOrderTranslator *util.ConcurrentOrderIdToOrderPtrMap
    // add timeouts
    httpClient               *http.Client
}

func NewBinanceUS(apiKey, secretKey string, assetPairTranslator types.AssetPairTranslator) *BinanceUS {
    assetPairs := assetPairTranslator.GetAssetPairs()
    httpClient := &http.Client{}
    return &BinanceUS{
        AssetPairTranslator: assetPairTranslator,
        apiKey: apiKey,
        secretKey: secretKey,
        spreadRecorder: NewBinanceUSSpreadRecorder(assetPairs, assetPairTranslator, 200),
        orderBookRecorder: NewBinanceUSOrderBookRecorder(httpClient, assetPairs, assetPairTranslator, 1000),
        latencyEstimator: util.NewEwmaEstimator(0.125, 0.25, 4),
        orderIdToOrderTranslator: util.NewConcurrentOrderIdToOrderPtrMap(),
        httpClient: httpClient,
    }
}

func (b *BinanceUS) String() string {
    return "BinanceUS"
}

func checkError(bodyJson map[string]interface{}) {
    if _, ok := bodyJson["code"]; ok {
        log.Fatalln(bodyJson)
    }
}

func (b *BinanceUS) getHistoricalSpread(assetPair types.AssetPair, duration time.Duration, samples uint, channel chan types.SpreadResponse) {
    if samples == 0 || duration <= 0 {
        channel <- types.SpreadResponse{assetPair, []types.Spread{}}
    }

    rawHistoricalSpreads, ok := b.spreadRecorder.GetHistoricalSpreads(assetPair)
    if len(rawHistoricalSpreads) == 0 {
        if !ok {
            b.spreadRecorder.RegisterAssetPair(assetPair)
        }
        channel <- types.SpreadResponse{assetPair, rawHistoricalSpreads}
        return
    }

    channel <- types.SpreadResponse{assetPair, util.GetSpreadSamples(rawHistoricalSpreads, duration, samples)}
}

func (b *BinanceUS) GetHistoricalSpreads(assetPairs []types.AssetPair, duration time.Duration, samples uint) map[types.AssetPair][]types.Spread {
    channel := make(chan types.SpreadResponse)
    for _, assetPair := range assetPairs {
        go b.getHistoricalSpread(assetPair, duration, samples, channel)
    }

    historicalSpreads := make(map[types.AssetPair][]types.Spread)
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        historicalSpreads[response.AssetPair] = response.HistoricalSpreads
    }
    return historicalSpreads
}

func (b *BinanceUS) GetCurrentSpread(assetPair types.AssetPair) types.Spread {
    bodyJson := util.HttpGetAndGetBody(b.httpClient, util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/ticker/bookTicker", url.Values{
        "symbol": []string{b.AssetPairTranslator[assetPair]},
    }))
    checkError(bodyJson)

    bid, err := decimal.NewFromString(bodyJson["bidPrice"].(string))
    if err != nil {
        log.Fatalln(err)
    }
    ask, err := decimal.NewFromString(bodyJson["askPrice"].(string))
    if err != nil {
        log.Fatalln(err)
    }

    return types.Spread{
        Bid: bid,
        Ask: ask,
    }
}

func (b *BinanceUS) getOrderBook(assetPair types.AssetPair, channel chan types.OrderBookResponse) {
    orderBook, ok := b.orderBookRecorder.GetOrderBook(assetPair)
    if !ok {
        b.orderBookRecorder.RegisterAssetPair(assetPair)
    }
    channel <- types.OrderBookResponse{assetPair, &orderBook}
}

func (b *BinanceUS) GetOrderBooks(assetPairs []types.AssetPair) map[types.AssetPair]*types.OrderBook {
    channel := make(chan types.OrderBookResponse)
    for _, assetPair := range assetPairs {
        go b.getOrderBook(assetPair, channel)
    }

    orderBooks := make(map[types.AssetPair]*types.OrderBook)
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        orderBooks[response.AssetPair] = response.OrderBook
    }
    return orderBooks
}

func (b *BinanceUS) GetLatency() time.Duration {
    start := time.Now()

    bodyJson := util.HttpGetAndGetBody(b.httpClient, RESTEndpoint + "/api/v3/ping")
    checkError(bodyJson)

    duration := time.Since(start)

    b.latencyEstimator.Sample(float64(duration.Milliseconds()))

    return time.Duration(b.latencyEstimator.GetEstimate()) * time.Millisecond
}

func (b *BinanceUS) parseOrderType(ot types.OrderType) string {
    if (ot == types.Buy) {
        return "BUY"
    }
    return "SELL"
}

func (b *BinanceUS) getBinanceUSSignature(values url.Values) string {
    mac := hmac.New(sha256.New, []byte(b.secretKey))
    mac.Write([]byte(values.Encode()))
    return fmt.Sprintf("%x", mac.Sum(nil))
}

func (b *BinanceUS) executeOrder(order types.Order, channel chan types.OrderIdResponse) {
    queryParams := url.Values{
        "symbol": []string{b.AssetPairTranslator[order.AssetPair]},
        "side": []string{b.parseOrderType(order.OrderType)},
        "type": []string{"LIMIT"},
        "timeInForce": []string{"GTC"},
        "price": []string{order.Price.String()},
        "quantity": []string{order.Quantity.String()},
        "timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }

    signature := b.getBinanceUSSignature(queryParams)

    request, err := http.NewRequest("POST", util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("X-MBX-APIKEY", b.apiKey)

    bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
    checkError(bodyJson)

    orderId := types.OrderId(strconv.FormatUint(uint64(bodyJson["orderId"].(float64)), 10))

    b.orderIdToOrderTranslator.Store(orderId, &order)

    channel <- types.OrderIdResponse{order, orderId}
}

func (b *BinanceUS) ExecuteOrders(orders []types.Order) map[types.Order]types.OrderId {
    channel := make(chan types.OrderIdResponse)
    for _, order := range orders {
        go b.executeOrder(order, channel)
    }

    orderIds := make(map[types.Order]types.OrderId)
    for i := 0; i < len(orders); i++ {
        response := <- channel
        orderIds[response.Order] = response.OrderId
    }
    return orderIds
}

func (b *BinanceUS) getOrderStatus(orderId types.OrderId, channel chan types.OrderStatusResponse) {
    order, ok := b.orderIdToOrderTranslator.Load(orderId)
    if !ok {
        log.Fatalf("order with id %v not found\n", orderId)
    }
    queryParams := url.Values{
        "symbol": []string{b.AssetPairTranslator[order.AssetPair]},
        "orderId": []string{string(orderId)},
        "timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }

    signature := b.getBinanceUSSignature(queryParams)

    request, err := http.NewRequest("GET", util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("X-MBX-APIKEY", b.apiKey)

    bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
    checkError(bodyJson)

    orderStatus := types.OrderStatus{
        Original: order,
    }

    switch status := bodyJson["status"].(string); status {
    case "NEW":
        orderStatus.Status = types.Unfilled
    case "PARTIALLY_FILLED":
        price, err := decimal.NewFromString(bodyJson["price"].(string))
        if err != nil {
            log.Fatalln(err)
        }
        quantity, err := decimal.NewFromString(bodyJson["executedQty"].(string))
        if err != nil {
            log.Fatalln(err)
        }
        orderStatus.Status = types.PartiallyFilled
        orderStatus.FilledPrice = &price
        orderStatus.FilledQuantity = &quantity
    case "FILLED":
        price, err := decimal.NewFromString(bodyJson["price"].(string))
        if err != nil {
            log.Fatalln(err)
        }
        quantity, err := decimal.NewFromString(bodyJson["executedQty"].(string))
        if err != nil {
            log.Fatalln(err)
        }
        orderStatus.Status = types.Filled
        orderStatus.FilledPrice = &price
        orderStatus.FilledQuantity = &quantity
        b.orderIdToOrderTranslator.Delete(orderId)
    case "CANCELED":
        orderStatus.Status = types.Canceled
        b.orderIdToOrderTranslator.Delete(orderId)
    case "EXPIRED":
        orderStatus.Status = types.Expired
        b.orderIdToOrderTranslator.Delete(orderId)
    case "REJECTED":
        log.Fatalf("order %v was rejected by BinanceUS\n", *order)
    }

    channel <- types.OrderStatusResponse{orderId, orderStatus}
}

func (b *BinanceUS) GetOrderStatuses(orderIds []types.OrderId) map[types.OrderId]types.OrderStatus {
    channel := make(chan types.OrderStatusResponse)
    for _, orderId := range orderIds {
        go b.getOrderStatus(orderId, channel)
    }

    orderStatuses := make(map[types.OrderId]types.OrderStatus)
    for i := 0; i < len(orderIds); i++ {
        response := <- channel
        orderStatuses[response.OrderId] = response.OrderStatus
    }
    return orderStatuses
}

func (b *BinanceUS) cancelOrder(orderId types.OrderId) {
    order, ok := b.orderIdToOrderTranslator.Load(orderId)
    if !ok {
        log.Fatalf("order with id %v not found\n", orderId)
    }
    queryParams := url.Values{
        "symbol": []string{b.AssetPairTranslator[order.AssetPair]},
        "orderId": []string{string(orderId)},
        "timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }

    signature := b.getBinanceUSSignature(queryParams)

    request, err := http.NewRequest("DELETE", util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("X-MBX-APIKEY", b.apiKey)

    bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
    checkError(bodyJson)

    b.orderIdToOrderTranslator.Delete(orderId)
}

func (b *BinanceUS) CancelOrders(orderIds []types.OrderId) {
    for _, orderId := range orderIds {
        go b.cancelOrder(orderId)
    }
}

func (b *BinanceUS) GetBalances() map[types.Asset]decimal.Decimal {
    queryParams := url.Values{
        "timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
    }

    signature := b.getBinanceUSSignature(queryParams)

    request, err := http.NewRequest("GET", util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/account", queryParams) + "&signature=" + signature, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("X-MBX-APIKEY", b.apiKey)

    bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
    checkError(bodyJson)

    balances := make(map[types.Asset]decimal.Decimal)
    for _, rawData := range bodyJson["balances"].([]interface{}) {
        data := rawData.(map[string]string)
        free, err := decimal.NewFromString(data["free"])
        if err != nil {
            log.Fatalln(err)
        }
        locked, err := decimal.NewFromString(data["locked"])
        if err != nil {
            log.Fatalln(err)
        }
        balances[types.Asset(data["asset"])] = free.Add(locked)
    }
    return balances
}
