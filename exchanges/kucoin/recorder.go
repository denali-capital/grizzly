package kucoin

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/gorilla/websocket"
    "github.com/shopspring/decimal"
)

// docs: https://docs.kucoin.com/#websocket-feed

type kuCoinMessage struct {
    Id       string `json:"id"`
    Type     string `json:"type"`
    Topic    string `json:"topic,omitempty"`
    Response bool   `json:"response,omitempty"`
}

func initializeWebSocketConnection(httpClient *http.Client) (time.Duration, *websocket.Conn) {
    request, err := http.NewRequest("POST", RESTEndpoint + "/api/v1/bullet-public", nil)
    if err != nil {
        log.Fatalln(err)
    }

    bodyJson := util.DoHttpAndGetBody(httpClient, request)
    checkError(bodyJson)

    data := bodyJson["data"].(map[string]interface{})
    instanceServer := data["instanceServers"].([]interface{})[0].(map[string]interface{})

    endpoint := util.ParseUrlWithQuery(instanceServer["endpoint"].(string), url.Values{
        "token": []string{data["token"].(string)},
    })

    webSocketConnection, _, err := websocket.DefaultDialer.Dial(endpoint, http.Header{})
    if err != nil {
        log.Fatalln(err)
    }

    var initialResponse map[string]interface{}
    err = webSocketConnection.ReadJSON(&initialResponse)
    if err != nil {
        log.Fatalln(err)
    }
    if initialResponse["type"].(string) != "welcome" {
        log.Fatalln(initialResponse)
    }

    return time.Duration(int64(instanceServer["pingInterval"].(float64))) * time.Millisecond, webSocketConnection
}

// should do separate connection for each asset pair?
type kuCoinWebSocketRecorder struct {
    sync.Mutex
    webSocketConnection *websocket.Conn
    assetPairTranslator types.AssetPairTranslator
    // map[string]chan map[string]interface{}
    channels            *sync.Map
}

func (k *kuCoinWebSocketRecorder) ping(interval time.Duration) {
    for {
        id := strconv.FormatInt(time.Now().UnixMilli(), 10)
        payloadJson, err := json.Marshal(kuCoinMessage{
            Id: id,
            Type: "ping",
        })
        if err != nil {
            log.Fatalln(err)
        }
        k.Lock()
        k.webSocketConnection.WriteMessage(1, payloadJson)
        for {
            var resp map[string]interface{}
            err := k.webSocketConnection.ReadJSON(&resp)
            if err != nil {
                log.Fatalln(err)
            }
            if resp["type"].(string) == "error" {
                log.Fatalln(resp)
            }
            if msgId, ok := resp["id"]; ok {
                if msgId.(string) != id {
                    log.Fatalf("id mismatch between sent %v and received %v\n", id, msgId)
                }
                if tpe := resp["type"].(string); tpe != "pong" {
                    log.Fatalf("expected pong, got %v\n", tpe)
                }
                break
            }
            topic := resp["topic"].(string)
            channel, ok := k.channels.Load(topic)
            if !ok {
                log.Fatalf("channel not found for topic %v\n", topic)
            }
            channel.(chan map[string]interface{}) <- util.MapCopy(resp)
        }
        k.Unlock()
        time.Sleep(interval)
    }
}

func (k *kuCoinWebSocketRecorder) record() {
    var resp map[string]interface{}
    for {
        k.Lock()
        err := k.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if resp["type"].(string) == "error" {
            log.Fatalln(resp)
        }
        topic := resp["topic"].(string)
        channel, ok := k.channels.Load(topic)
        if !ok {
            log.Fatalf("channel not found for topic %v\n", topic)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
        k.Unlock()
    }
}

type KuCoinSpreadRecorder struct {
    kuCoinWebSocketRecorder
    capacity                 uint
    // map[types.AssetPair]*util.ConcurrentFixedSizeSpreadQueue
    historicalSpreads        *sync.Map
}

func NewKuCoinSpreadRecorder(httpClient *http.Client, assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, capacity uint) *KuCoinSpreadRecorder {
    pingInterval, webSocketConnection := initializeWebSocketConnection(httpClient)

    assetPairNames := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        assetPairNames[i] = assetPairTranslator[assetPair]
    }

    id := strconv.FormatInt(time.Now().UnixMilli(), 10)
    payloadJson, err := json.Marshal(kuCoinMessage{
        Id: id,
        Type: "subscribe",
        Topic: "/market/ticker:" + strings.Join(assetPairNames, ","),
        Response: true,
    })
    if err != nil {
        log.Fatalln(err)
    }
    webSocketConnection.WriteMessage(1, payloadJson)

    var initialResponse map[string]interface{}
    err = webSocketConnection.ReadJSON(&initialResponse)
    if err != nil {
        log.Fatalln(err)
    }
    if !(initialResponse["id"].(string) == id && initialResponse["type"].(string) == "ack") {
        log.Fatalln(initialResponse)
    }

    channels := &sync.Map{}
    historicalSpreads := &sync.Map{}
    for _, assetPair := range assetPairs {
        channel := make(chan map[string]interface{})
        historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(capacity)

        channels.Store("/market/ticker:" + assetPairTranslator[assetPair], channel)
        historicalSpreads.Store(assetPair, historicalSpread)

        go processSpreadUpdates(historicalSpread, channel)
    }

    kuCoinSpreadRecorder := &KuCoinSpreadRecorder{
        kuCoinWebSocketRecorder: kuCoinWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            assetPairTranslator: assetPairTranslator,
            channels: channels,
        },
        capacity: capacity,
        historicalSpreads: historicalSpreads,
    }

    go kuCoinSpreadRecorder.ping(pingInterval)
    go kuCoinSpreadRecorder.record()

    return kuCoinSpreadRecorder
}

func processSpreadUpdates(historicalSpread *util.ConcurrentFixedSizeSpreadQueue, channel chan map[string]interface{}) {
    for {
        select {
        case resp := <- channel:
            rawSpread := resp["data"].(map[string]interface{})
            bid, err := decimal.NewFromString(rawSpread["bestBid"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            ask, err := decimal.NewFromString(rawSpread["bestAsk"].(string))
            if err != nil {
                log.Fatalln(err)
            }

            historicalSpread.Push(types.Spread{
                Bid: bid,
                Ask: ask,
                Timestamp: time.UnixMilli(int64(rawSpread["time"].(float64))),
            })
        }
    }
}

func (k *KuCoinSpreadRecorder) GetHistoricalSpreads(assetPair types.AssetPair) ([]types.Spread, bool) {
    result, ok := k.historicalSpreads.Load(assetPair)
    if !ok {
        return make([]types.Spread, 0), false
    }
    return result.(*util.ConcurrentFixedSizeSpreadQueue).Data(), true
}

func (k *KuCoinSpreadRecorder) GetCurrentSpread(assetPair types.AssetPair) (types.Spread, bool) {
    result, ok := k.historicalSpreads.Load(assetPair)
    if !ok {
        return types.Spread{}, false
    }
    return result.(*util.ConcurrentFixedSizeSpreadQueue).Back(), true
}

func (k *KuCoinSpreadRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := k.historicalSpreads.Load(assetPair); ok {
        return
    }

    topic := "/market/ticker:" + k.assetPairTranslator[assetPair]

    id := strconv.FormatInt(time.Now().UnixMilli(), 10)
    payloadJson, err := json.Marshal(kuCoinMessage{
        Id: id,
        Type: "subscribe",
        Topic: topic,
        Response: true,
    })
    if err != nil {
        log.Fatalln(err)
    }
    k.Lock()
    defer k.Unlock()
    k.webSocketConnection.WriteMessage(1, payloadJson)
    for {
        var resp map[string]interface{}
        err := k.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if resp["type"].(string) == "error" {
            log.Fatalln(resp)
        }
        if msgId, ok := resp["id"]; ok {
            if msgId.(string) != id {
                log.Fatalf("id mismatch between sent %v and received %v\n", id, msgId)
            }
            if tpe := resp["type"].(string); tpe != "ack" {
                log.Fatalf("expected ack, got %v\n", tpe)
            }
            break
        }
        topic := resp["topic"].(string)
        channel, ok := k.channels.Load(topic)
        if !ok {
            log.Fatalf("channel not found for topic %v\n", topic)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
    }
    channel := make(chan map[string]interface{})
    historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(k.capacity)

    k.channels.Store(topic, channel)
    k.historicalSpreads.Store(assetPair, historicalSpread)

    go processSpreadUpdates(historicalSpread, channel)
}

type KuCoinOrderBookRecorder struct {
    kuCoinWebSocketRecorder
    apiKey                  string
    secretKey               string
    apiPassphrase           string
    httpClient              *http.Client
    depth                   uint
    // map[types.AssetPair]*util.ConcurrentOrderBook
    orderBooks              *sync.Map
}

func getOrderBookSnapshot(httpClient *http.Client, apiKey, secretKey, apiPassphrase string, assetPair types.AssetPair, assetPairTranslator types.AssetPairTranslator, depth uint, channel chan util.ConcurrentOrderBookResponse) {
    time := strconv.FormatInt(time.Now().UnixMilli(), 10)
    path := util.ParseUrlWithQuery("/api/v3/market/orderbook/level2", url.Values{
        "symbol": []string{assetPairTranslator[assetPair]},
    })

    signature, passphrase := getKuCoinSignatureAndPassphrase(secretKey, apiPassphrase, time, "GET", path, "")

    request, err := http.NewRequest("GET", RESTEndpoint + path, nil)
    if err != nil {
        log.Fatalln(err)
    }
    request.Header.Set("KC-API-SIGN", signature)
    request.Header.Set("KC-API-TIMESTAMP", time)
    request.Header.Set("KC-API-KEY", apiKey)
    request.Header.Set("KC-API-PASSPHRASE", passphrase)
    request.Header.Set("KC-API-KEY-VERSION", "2")

    bodyJson := util.DoHttpAndGetBody(httpClient, request)
    checkError(bodyJson)

    data := bodyJson["data"].(map[string]interface{})

    sequence, err := strconv.ParseUint(data["sequence"].(string), 10, 64)
    if err != nil {
        log.Fatalln(err)
    }
    asks := make([]types.OrderBookEntry, 0)
    bids := make([]types.OrderBookEntry, 0)
    for _, rawOrderBookEntry := range data["asks"].([]interface{}) {
        price, err := decimal.NewFromString(rawOrderBookEntry.([]interface{})[0].(string))
        if err != nil {
            log.Fatalln(err)
        }
        quantity, err := decimal.NewFromString(rawOrderBookEntry.([]interface{})[1].(string))
        if err != nil {
            log.Fatalln(err)
        }
        asks = append(asks, types.OrderBookEntry{
            Price: price,
            Quantity: quantity,
            UpdateId: uint(sequence),
        })
        if uint(len(asks)) == depth {
            break
        }
    }
    for _, rawOrderBookEntry := range data["bids"].([]interface{}) {
        price, err := decimal.NewFromString(rawOrderBookEntry.([]interface{})[0].(string))
        if err != nil {
            log.Fatalln(err)
        }
        quantity, err := decimal.NewFromString(rawOrderBookEntry.([]interface{})[1].(string))
        if err != nil {
            log.Fatalln(err)
        }
        bids = append(bids, types.OrderBookEntry{
            Price: price,
            Quantity: quantity,
            UpdateId: uint(sequence),
        })
        if uint(len(bids)) == depth {
            break
        }
    }
    concurrentOrderBook := util.NewConcurrentOrderBook(bids, asks)
    concurrentOrderBook.LastUpdateId = uint(sequence)

    channel <- util.ConcurrentOrderBookResponse{assetPair, concurrentOrderBook}
}

func getOrderBookSnapshots(httpClient *http.Client, apiKey, secretKey, apiPassphrase string, assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, depth uint) (map[types.AssetPair]*util.ConcurrentOrderBook, *sync.Map) {
    channel := make(chan util.ConcurrentOrderBookResponse)
    for _, assetPair := range assetPairs {
        go getOrderBookSnapshot(httpClient, apiKey, secretKey, apiPassphrase, assetPair, assetPairTranslator, depth, channel)
    }

    orderBooks := make(map[types.AssetPair]*util.ConcurrentOrderBook)
    // map[types.AssetPair]*util.ConcurrentOrderBook
    syncOrderBooks := &sync.Map{}
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        orderBooks[response.AssetPair] = response.ConcurrentOrderBook
        syncOrderBooks.Store(response.AssetPair, response.ConcurrentOrderBook)
    }
    return orderBooks, syncOrderBooks
}

func NewKuCoinOrderBookRecorder(httpClient *http.Client, apiKey, secretKey, apiPassphrase string, assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, depth uint) *KuCoinOrderBookRecorder {
    pingInterval, webSocketConnection := initializeWebSocketConnection(httpClient)

    assetPairNames := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        assetPairNames[i] = assetPairTranslator[assetPair]
    }

    id := strconv.FormatInt(time.Now().UnixMilli(), 10)
    payloadJson, err := json.Marshal(kuCoinMessage{
        Id: id,
        Type: "subscribe",
        Topic: "/market/level2:" + strings.Join(assetPairNames, ","),
        Response: true,
    })
    if err != nil {
        log.Fatalln(err)
    }
    webSocketConnection.WriteMessage(1, payloadJson)

    var initialResponse map[string]interface{}
    err = webSocketConnection.ReadJSON(&initialResponse)
    if err != nil {
        log.Fatalln(err)
    }
    if !(initialResponse["id"].(string) == id && initialResponse["type"].(string) == "ack") {
        log.Fatalln(initialResponse)
    }

    channels := &sync.Map{}
    orderBooks, syncOrderBooks := getOrderBookSnapshots(httpClient, apiKey, secretKey, apiPassphrase, assetPairs, assetPairTranslator, depth)
    for _, assetPair := range assetPairs {
        channel := make(chan map[string]interface{})

        channels.Store("/market/level2:" + assetPairTranslator[assetPair], channel)

        go processOrderBookUpdates(httpClient, assetPair, assetPairTranslator, orderBooks[assetPair], channel, depth)
    }

    kuCoinOrderBookRecorder := &KuCoinOrderBookRecorder{
        kuCoinWebSocketRecorder: kuCoinWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            assetPairTranslator: assetPairTranslator,
            channels: channels,
        },
        apiKey: apiKey,
        secretKey: secretKey,
        apiPassphrase: apiPassphrase,
        httpClient: httpClient,
        depth: depth,
        orderBooks: syncOrderBooks,
    }

    go kuCoinOrderBookRecorder.ping(pingInterval)
    go kuCoinOrderBookRecorder.record()

    return kuCoinOrderBookRecorder
}

func processOrderBookUpdates(httpClient *http.Client, assetPair types.AssetPair, assetPairTranslator types.AssetPairTranslator, concurrentOrderBook *util.ConcurrentOrderBook, channel chan map[string]interface{}, depth uint) {
    for {
        select {
        case resp := <- channel:
            changes := resp["data"].(map[string]interface{})["changes"].(map[string]interface{})
            bids := concurrentOrderBook.GetBids()
            asks := concurrentOrderBook.GetAsks()
            maxSequence := uint64(0)
            for _, rawOrderBookEntry := range changes["bids"].([]interface{}) {
                sequence, err := strconv.ParseUint(rawOrderBookEntry.([]interface{})[2].(string), 10, 64)
                if err != nil {
                    log.Fatalln(err)
                }
                if uint(sequence) < concurrentOrderBook.LastUpdateId {
                    continue
                }
                if sequence > maxSequence {
                    maxSequence = sequence
                }
                price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                if quantity.Equal(decimal.Zero) {
                    bids = util.RemovePriceFromBids(bids, price)
                } else {
                    bids = util.InsertPriceInBids(bids, types.OrderBookEntry{
                        Price: price,
                        Quantity: quantity,
                        UpdateId: uint(sequence),
                    })
                    bids = bids[:util.MinUint(depth, uint(len(bids)))]
                }
            }
            for _, rawOrderBookEntry := range changes["asks"].([]interface{}) {
                sequence, err := strconv.ParseUint(rawOrderBookEntry.([]interface{})[2].(string), 10, 64)
                if err != nil {
                    log.Fatalln(err)
                }
                if uint(sequence) < concurrentOrderBook.LastUpdateId {
                    continue
                }
                if sequence > maxSequence {
                    maxSequence = sequence
                }
                price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                if quantity.Equal(decimal.Zero) {
                    asks = util.RemovePriceFromAsks(asks, price)
                } else {
                    asks = util.InsertPriceInAsks(asks, types.OrderBookEntry{
                        Price: price,
                        Quantity: quantity,
                        UpdateId: uint(sequence),
                    })
                    asks = asks[:util.MinUint(depth, uint(len(asks)))]
                }
            }
            if maxSequence > 0 {
                concurrentOrderBook.LastUpdateId = uint(maxSequence)
            }
            concurrentOrderBook.SetBidsAndAsks(bids[:util.MinUint(depth, uint(len(bids)))], asks[:util.MinUint(depth, uint(len(asks)))])
        }
    }
}

func (k *KuCoinOrderBookRecorder) GetOrderBook(assetPair types.AssetPair) (types.OrderBook, bool) {
    result, ok := k.orderBooks.Load(assetPair)
    if !ok {
        return types.OrderBook{}, false
    }
    return result.(*util.ConcurrentOrderBook).Data(), true
}

func (k *KuCoinOrderBookRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := k.orderBooks.Load(assetPair); ok {
        return
    }

    topic := "/market/level2:" + k.assetPairTranslator[assetPair]

    id := strconv.FormatInt(time.Now().UnixMilli(), 10)
    payloadJson, err := json.Marshal(kuCoinMessage{
        Id: id,
        Type: "subscribe",
        Topic: topic,
        Response: true,
    })
    if err != nil {
        log.Fatalln(err)
    }
    k.Lock()
    defer k.Unlock()
    k.webSocketConnection.WriteMessage(1, payloadJson)
    for {
        var resp map[string]interface{}
        err := k.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if resp["type"].(string) == "error" {
            log.Fatalln(resp)
        }
        if msgId, ok := resp["id"]; ok {
            if msgId.(string) != id {
                log.Fatalf("id mismatch between sent %v and received %v\n", id, msgId)
            }
            if tpe := resp["type"].(string); tpe != "ack" {
                log.Fatalf("expected ack, got %v\n", tpe)
            }
            break
        }
        topic := resp["topic"].(string)
        channel, ok := k.channels.Load(topic)
        if !ok {
            log.Fatalf("channel not found for topic %v\n", topic)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
    }
    channel := make(chan map[string]interface{})

    k.channels.Store(topic, channel)

    snapshotChannel := make(chan util.ConcurrentOrderBookResponse)
    go getOrderBookSnapshot(k.httpClient, k.apiKey, k.secretKey, k.apiPassphrase, assetPair, k.assetPairTranslator, k.depth, snapshotChannel)
    select {
    case resp := <- snapshotChannel:
        k.orderBooks.Store(assetPair, resp.ConcurrentOrderBook)
        go processOrderBookUpdates(k.httpClient, assetPair, k.assetPairTranslator, resp.ConcurrentOrderBook, channel, k.depth)
    }
}
