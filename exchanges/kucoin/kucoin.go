package kucoin

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// docs: https://docs.kucoin.com/
const RESTEndpoint string = "https://api.kucoin.com"

type KuCoin struct {
    AssetPairTranslator      types.AssetPairTranslator

    apiKey                   string
    secretKey                string
    apiPassphrase            string
    spreadRecorder           types.SpreadRecorder
    orderBookRecorder        types.OrderBookRecorder
    latencyEstimator         *util.EwmaEstimator
    orderIdToOrderTranslator *util.ConcurrentOrderIdToOrderPtrMap
    // add timeouts
    httpClient               *http.Client
}

func NewKuCoin(apiKey, secretKey, apiPassphrase string, assetPairTranslator types.AssetPairTranslator) *KuCoin {
    assetPairs := assetPairTranslator.GetAssetPairs()
    httpClient := &http.Client{}
    return &KuCoin{
        AssetPairTranslator: assetPairTranslator,
        apiKey: apiKey,
        secretKey: secretKey,
        apiPassphrase: apiPassphrase,
        spreadRecorder: NewKuCoinSpreadRecorder(httpClient, assetPairs, assetPairTranslator, 200),
        orderBookRecorder: NewKuCoinOrderBookRecorder(httpClient, apiKey, secretKey, apiPassphrase, assetPairs, assetPairTranslator, 1000),
        latencyEstimator: util.NewEwmaEstimator(0.125, 0.25, 4),
        orderIdToOrderTranslator: util.NewConcurrentOrderIdToOrderPtrMap(),
        httpClient: httpClient,
    }
}

func (k *KuCoin) String() string {
    return "KuCoin"
}

func checkError(bodyJson map[string]interface{}) {
    if bodyJson["code"].(string) != "200000" {
        log.Fatalln(bodyJson)
    }
}

func (k *KuCoin) getHistoricalSpread(assetPair types.AssetPair, duration time.Duration, samples uint, channel chan types.SpreadResponse) {
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

func (k *KuCoin) GetHistoricalSpreads(assetPairs []types.AssetPair, duration time.Duration, samples uint) map[types.AssetPair][]types.Spread {
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

func (k *KuCoin) GetCurrentSpread(assetPair types.AssetPair) types.Spread {
    bodyJson := util.HttpGetAndGetBody(k.httpClient, util.ParseUrlWithQuery(RESTEndpoint + "/api/v1/market/orderbook/level1", url.Values{
        "symbol": []string{k.AssetPairTranslator[assetPair]},
    }))
    checkError(bodyJson)

    data := bodyJson["data"].(map[string]interface{})
    bid, err := decimal.NewFromString(data["bestBid"].(string))
    if err != nil {
        log.Fatalln(err)
    }
    ask, err := decimal.NewFromString(data["bestAsk"].(string))
    if err != nil {
        log.Fatalln(err)
    }

    return types.Spread{
        Bid: bid,
        Ask: ask,
    }
}

func (k *KuCoin) getOrderBook(assetPair types.AssetPair, channel chan types.OrderBookResponse) {
    orderBook, ok := k.orderBookRecorder.GetOrderBook(assetPair)
    if !ok {
        k.orderBookRecorder.RegisterAssetPair(assetPair)
    }
    channel <- types.OrderBookResponse{assetPair, &orderBook}
}
 
func (k *KuCoin) GetOrderBooks(assetPairs []types.AssetPair) map[types.AssetPair]*types.OrderBook {
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

func (k *KuCoin) GetLatency() time.Duration {
    start := time.Now()

    bodyJson := util.HttpGetAndGetBody(k.httpClient, RESTEndpoint + "/api/v1/timestamp")
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

func getKuCoinSignatureAndPassphrase(secretKey, apiPassphrase, time, method, path, data string) (string, string) {
    toSign := time + method + path + data
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(toSign))
    signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    mac = hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(apiPassphrase))
    passphrase := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    return signature, passphrase
}

func (k *KuCoin) executeOrder(order types.Order, channel chan types.OrderIdResponse) {
    data, err := json.Marshal(map[string]interface{}{
        "clientOid": uuid.NewString(),
        "side": parseOrderType(order.OrderType),
        "symbol": k.AssetPairTranslator[order.AssetPair],
        "price": order.Price.String(),
        "size": order.Quantity.String(),
    })
    if err != nil {
        log.Fatalln(err)
    }

    time := strconv.FormatInt(time.Now().UnixMilli(), 10)

    signature, passphrase := getKuCoinSignatureAndPassphrase(k.secretKey, k.apiPassphrase, time, "POST", "/api/v1/orders", string(data))

    request, err := http.NewRequest("POST", RESTEndpoint + "/api/v1/orders", bytes.NewReader(data))
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("KC-API-SIGN", signature)
    request.Header.Set("KC-API-TIMESTAMP", time)
    request.Header.Set("KC-API-KEY", k.apiKey)
    request.Header.Set("KC-API-PASSPHRASE", passphrase)
    request.Header.Set("KC-API-KEY-VERSION", "2")

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    jsonData := bodyJson["data"].(map[string]interface{})

    orderId := types.OrderId(jsonData["orderId"].(string))

    k.orderIdToOrderTranslator.Store(orderId, &order)

    channel <- types.OrderIdResponse{order, orderId}
}

func (k *KuCoin) ExecuteOrders(orders []types.Order) map[types.Order]types.OrderId {
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

func (k *KuCoin) getOrderStatus(orderId types.OrderId, channel chan types.OrderStatusResponse) {
    order, ok := k.orderIdToOrderTranslator.Load(orderId)
    if !ok {
        log.Fatalf("order with id %v not found\n", orderId)
    }

    time := strconv.FormatInt(time.Now().UnixMilli(), 10)
    path := "/api/v1/orders/" + string(orderId)

    signature, passphrase := getKuCoinSignatureAndPassphrase(k.secretKey, k.apiPassphrase, time, "GET", path, "")

    request, err := http.NewRequest("GET", RESTEndpoint + path, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("KC-API-SIGN", signature)
    request.Header.Set("KC-API-TIMESTAMP", time)
    request.Header.Set("KC-API-KEY", k.apiKey)
    request.Header.Set("KC-API-PASSPHRASE", passphrase)
    request.Header.Set("KC-API-KEY-VERSION", "2")

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    data := bodyJson["data"].(map[string]interface{})

    orderStatus := types.OrderStatus{
        Original: order,
    }

    if data["isActive"].(bool) {
        orderStatus.Status = types.Unfilled
    } else {
        if data["cancelExist"].(bool) {
            orderStatus.Status = types.Canceled
            k.orderIdToOrderTranslator.Delete(orderId)
        } else {
            price, err := decimal.NewFromString(data["price"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            quantity, err := decimal.NewFromString(data["size"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            orderStatus.Status = types.Filled
            orderStatus.FilledPrice = &price
            orderStatus.FilledQuantity = &quantity
            k.orderIdToOrderTranslator.Delete(orderId)
        }
    }

    channel <- types.OrderStatusResponse{orderId, orderStatus}
}

func (k *KuCoin) GetOrderStatuses(orderIds []types.OrderId) map[types.OrderId]types.OrderStatus {
    channel := make(chan types.OrderStatusResponse)
    for _, orderId := range orderIds {
        go k.getOrderStatus(orderId, channel)
    }

    orderStatuses := make(map[types.OrderId]types.OrderStatus)
    for i := 0; i < len(orderIds); i++ {
        response := <- channel
        orderStatuses[response.OrderId] = response.OrderStatus
    }
    return orderStatuses
}

func (k *KuCoin) cancelOrder(orderId types.OrderId) {
    time := strconv.FormatInt(time.Now().UnixMilli(), 10)
    path := "/api/v1/orders/" + string(orderId)

    signature, passphrase := getKuCoinSignatureAndPassphrase(k.secretKey, k.apiPassphrase, time, "DELETE", path, "")

    request, err := http.NewRequest("DELETE", RESTEndpoint + path, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("KC-API-SIGN", signature)
    request.Header.Set("KC-API-TIMESTAMP", time)
    request.Header.Set("KC-API-KEY", k.apiKey)
    request.Header.Set("KC-API-PASSPHRASE", passphrase)
    request.Header.Set("KC-API-KEY-VERSION", "2")

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    k.orderIdToOrderTranslator.Delete(orderId)
}

func (k *KuCoin) CancelOrders(orderIds []types.OrderId) {
    for _, orderId := range orderIds {
        go k.cancelOrder(orderId)
    }
}

func (k *KuCoin) GetBalances() map[types.Asset]decimal.Decimal {
    time := strconv.FormatInt(time.Now().UnixMilli(), 10)

    signature, passphrase := getKuCoinSignatureAndPassphrase(k.secretKey, k.apiPassphrase, time, "GET", "/api/v1/accounts", "")

    request, err := http.NewRequest("GET", RESTEndpoint + "/api/v1/accounts", nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("KC-API-SIGN", signature)
    request.Header.Set("KC-API-TIMESTAMP", time)
    request.Header.Set("KC-API-KEY", k.apiKey)
    request.Header.Set("KC-API-PASSPHRASE", passphrase)
    request.Header.Set("KC-API-KEY-VERSION", "2")

    bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
    checkError(bodyJson)

    data := bodyJson["data"].([]interface{})

    balances := make(map[types.Asset]decimal.Decimal)
    for _, rawData := range data {
        data := rawData.(map[string]string)
        balance, err := decimal.NewFromString(data["balance"])
        if err != nil {
            log.Fatalln(err)
        }
        balances[types.Asset(data["currency"])] = balances[types.Asset(data["currency"])].Add(balance)
    }
    return balances
}
