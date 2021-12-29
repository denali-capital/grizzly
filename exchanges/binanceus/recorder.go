package binanceus

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/gorilla/websocket"
    "github.com/shopspring/decimal"
)

// docs: https://docs.binanceUS.com/websockets
const WebSocketEndpoint string = "wss://stream.binance.us:9443"
const CombinedStreamIndicator string = "/stream?streams="

type binanceUSSubscriptionMessage struct {
    Method string   `json:"method"`
    Params []string `json:"params"`
    Id     uint     `json:"id"`
}

// should do separate connection for each asset pair?
type binanceUSWebSocketRecorder struct {
    sync.Mutex
    webSocketConnection *websocket.Conn
    assetPairTranslator types.AssetPairTranslator
    // map[string]chan map[string]interface{}
    channels            *sync.Map
    id                  uint
}

func (b *binanceUSWebSocketRecorder) record() {
    var resp map[string]interface{}
    for {
        b.Lock()
        err := b.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if _, ok := resp["code"]; ok {
            log.Fatalln(resp)
        }
        streamName := resp["stream"].(string)
        channel, ok := b.channels.Load(streamName)
        if !ok {
            log.Fatalf("channel not found for streamName %v\n", streamName)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
        b.Unlock()
    }
}

type BinanceUSSpreadRecorder struct {
    binanceUSWebSocketRecorder
    capacity                 uint
    // map[types.AssetPair]*util.ConcurrentFixedSizeSpreadQueue
    historicalSpreads        *sync.Map
}

func NewBinanceUSSpreadRecorder(assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, capacity uint) *BinanceUSSpreadRecorder {
    streamTranslator := make(map[string]types.AssetPair)
    streams := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        streamName := strings.ToLower(assetPairTranslator[assetPair]) + "@bookTicker"
        streamTranslator[streamName] = assetPair
        streams[i] = streamName
    }
    endpoint := WebSocketEndpoint + CombinedStreamIndicator + strings.Join(streams, "/")

    webSocketConnection, _, err := websocket.DefaultDialer.Dial(endpoint, http.Header{})
    if err != nil {
        log.Fatalln(err)
    }

    channels := &sync.Map{}
    historicalSpreads := &sync.Map{}
    for _, streamName := range streams {
        channel := make(chan map[string]interface{})
        historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(capacity)

        channels.Store(streamName, channel)
        historicalSpreads.Store(streamTranslator[streamName], historicalSpread)

        go processSpreadUpdates(historicalSpread, channel)
    }

    binanceUSSpreadRecorder := &BinanceUSSpreadRecorder{
        binanceUSWebSocketRecorder: binanceUSWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            assetPairTranslator: assetPairTranslator,
            channels: channels,
        },
        capacity: capacity,
        historicalSpreads: historicalSpreads,
    }

    go binanceUSSpreadRecorder.record()

    return binanceUSSpreadRecorder
}

func processSpreadUpdates(historicalSpread *util.ConcurrentFixedSizeSpreadQueue, channel chan map[string]interface{}) {
    for {
        select {
        case resp := <- channel:
            rawSpread := resp["data"].(map[string]interface{})
            bid, err := decimal.NewFromString(rawSpread["b"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            ask, err := decimal.NewFromString(rawSpread["a"].(string))
            if err != nil {
                log.Fatalln(err)
            }
            timestamp := time.Now()

            historicalSpread.Push(types.Spread{
                Bid: bid,
                Ask: ask,
                Timestamp: &timestamp,
            })
        }
    }
}

func (b *BinanceUSSpreadRecorder) GetHistoricalSpreads(assetPair types.AssetPair) ([]types.Spread, bool) {
    result, ok := b.historicalSpreads.Load(assetPair)
    if !ok {
        return make([]types.Spread, 0), false
    }
    return result.(*util.ConcurrentFixedSizeSpreadQueue).Data(), true
}

func (b *BinanceUSSpreadRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := b.historicalSpreads.Load(assetPair); ok {
        return
    }

    streamName := strings.ToLower(b.assetPairTranslator[assetPair]) + "@bookTicker"

    payloadJson, err := json.Marshal(binanceUSSubscriptionMessage{
        Method: "SUBSCRIBE",
        Params: []string{streamName},
        Id: b.id,
    })
    if err != nil {
        log.Fatalln(err)
    }
    b.Lock()
    defer b.Unlock()
    b.webSocketConnection.WriteMessage(1, payloadJson)
    var resp map[string]interface{}
    for {
        err := b.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if _, ok := resp["code"]; ok {
            log.Fatalln(resp)
        }
        if id, ok := resp["id"]; ok {
            if uint(id.(float64)) != b.id {
                log.Fatalf("id mismatch between sent %v and received %v\n", b.id, id)
            }
            b.id++
            break
        }
        streamName := resp["stream"].(string)
        channel, ok := b.channels.Load(streamName)
        if !ok {
            log.Fatalf("channel not found for streamName %v\n", streamName)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
    }
    channel := make(chan map[string]interface{})
    historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(b.capacity)

    b.channels.Store(streamName, channel)
    b.historicalSpreads.Store(assetPair, historicalSpread)

    go processSpreadUpdates(historicalSpread, channel)
}

type BinanceUSOrderBookRecorder struct {
    binanceUSWebSocketRecorder
    httpClient              *http.Client
    depth                   uint
    // map[types.AssetPair]*util.ConcurrentOrderBook
    orderBooks              *sync.Map
}

type concurrentOrderBookResponse struct {
    assetPair           types.AssetPair
    concurrentOrderBook *util.ConcurrentOrderBook
}

var SnapshotLimits [8]uint = [8]uint{
    5, 10, 20, 50, 100, 500, 1000, 5000,
}

func selectLimit(depth uint) uint {
    n := len(SnapshotLimits)
    i := sort.Search(n, func(i int) bool {
        return depth <= SnapshotLimits[i]
    })
    if i == n {
        return SnapshotLimits[n - 1]
    }
    return SnapshotLimits[i]
}

func getOrderBookSnapshot(httpClient *http.Client, assetPair types.AssetPair, assetPairTranslator types.AssetPairTranslator, limit uint, channel chan concurrentOrderBookResponse) {
    bodyJson := util.HttpGetAndGetBody(httpClient, util.ParseUrlWithQuery(RESTEndpoint + "/api/v3/depth", url.Values{
        "symbol": []string{assetPairTranslator[assetPair]},
        "limit": []string{strconv.FormatUint(uint64(limit), 10)},
    }))
    if _, ok := bodyJson["code"]; ok {
        log.Fatalln(bodyJson)
    }

    lastUpdateId := uint(bodyJson["lastUpdateId"].(float64))
    asks := make([]types.OrderBookEntry, 0)
    bids := make([]types.OrderBookEntry, 0)
    for _, rawOrderBookEntry := range bodyJson["asks"].([]interface{}) {
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
            UpdateId: lastUpdateId,
        })
    }
    for _, rawOrderBookEntry := range bodyJson["bids"].([]interface{}) {
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
            UpdateId: lastUpdateId,
        })
    }
    concurrentOrderBook := util.NewConcurrentOrderBook(bids, asks)

    channel <- concurrentOrderBookResponse{assetPair, concurrentOrderBook}
}

func getOrderBookSnapshots(httpClient *http.Client, assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, depth uint) (map[types.AssetPair]*util.ConcurrentOrderBook, *sync.Map) {
    channel := make(chan concurrentOrderBookResponse)
    for _, assetPair := range assetPairs {
        go getOrderBookSnapshot(httpClient, assetPair, assetPairTranslator, selectLimit(depth), channel)
    }

    orderBooks := make(map[types.AssetPair]*util.ConcurrentOrderBook)
    // map[types.AssetPair]*util.ConcurrentOrderBook
    syncOrderBooks := &sync.Map{}
    for i := 0; i < len(assetPairs); i++ {
        response := <- channel
        orderBooks[response.assetPair] = response.concurrentOrderBook
        syncOrderBooks.Store(response.assetPair, response.concurrentOrderBook)
    }
    return orderBooks, syncOrderBooks
}

func NewBinanceUSOrderBookRecorder(httpClient *http.Client, assetPairs []types.AssetPair, assetPairTranslator types.AssetPairTranslator, depth uint) *BinanceUSOrderBookRecorder {
    streamTranslator := make(map[string]types.AssetPair)
    streams := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        streamName := strings.ToLower(assetPairTranslator[assetPair]) + "@depth"
        streamTranslator[streamName] = assetPair
        streams[i] = streamName
    }
    endpoint := WebSocketEndpoint + CombinedStreamIndicator + strings.Join(streams, "/")

    webSocketConnection, _, err := websocket.DefaultDialer.Dial(endpoint, http.Header{})
    if err != nil {
        log.Fatalln(err)
    }

    channels := &sync.Map{}
    orderBooks, syncOrderBooks := getOrderBookSnapshots(httpClient, assetPairs, assetPairTranslator, depth)
    for _, streamName := range streams {
        assetPair := streamTranslator[streamName]
        channel := make(chan map[string]interface{})

        channels.Store(streamName, channel)

        go processOrderBookUpdates(httpClient, assetPair, assetPairTranslator, orderBooks[assetPair], channel, depth)
    }

    binanceUSOrderBookRecorder := &BinanceUSOrderBookRecorder{
        binanceUSWebSocketRecorder: binanceUSWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            assetPairTranslator: assetPairTranslator,
            channels: channels,
        },
        httpClient: httpClient,
        depth: depth,
        orderBooks: syncOrderBooks,
    }

    go binanceUSOrderBookRecorder.record()

    return binanceUSOrderBookRecorder
}

func validateUpdateId(eventId uint, lastUpdateId uint) bool {
    if eventId - lastUpdateId != 1 && lastUpdateId != 0 {
        return false
    }
    return true
}

func processOrderBookUpdates(httpClient *http.Client, assetPair types.AssetPair, assetPairTranslator types.AssetPairTranslator, concurrentOrderBook *util.ConcurrentOrderBook, channel chan map[string]interface{}, depth uint) {
    for {
        select {
        case resp := <- channel:
            data := resp["data"].(map[string]interface{})
            bids := concurrentOrderBook.GetBids()
            asks := concurrentOrderBook.GetAsks()
            if validateUpdateId(uint(data["U"].(float64)), concurrentOrderBook.LastUpdateId) {
                lastUpdateId := uint(data["u"].(float64))
                for _, rawOrderBookEntry := range data["b"].([]interface{}) {
                    price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                    if quantity.Equal(decimal.Zero) {
                        bids = util.RemovePriceFromBids(bids, price)
                    } else {
                        bids = util.InsertPriceInBids(bids, types.OrderBookEntry{
                            Price: price,
                            Quantity: quantity,
                            UpdateId: lastUpdateId,
                        })
                        bids = bids[:util.MinUint(depth, uint(len(bids)))]
                    }
                }
                for _, rawOrderBookEntry := range data["a"].([]interface{}) {
                    price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                    if quantity.Equal(decimal.Zero) {
                        asks = util.RemovePriceFromAsks(asks, price)
                    } else {
                        asks = util.InsertPriceInAsks(asks, types.OrderBookEntry{
                            Price: price,
                            Quantity: quantity,
                            UpdateId: lastUpdateId,
                        })
                        asks = asks[:util.MinUint(depth, uint(len(asks)))]
                    }
                }
                concurrentOrderBook.LastUpdateId = lastUpdateId
                concurrentOrderBook.SetBidsAndAsks(bids[:util.MinUint(depth, uint(len(bids)))], asks[:util.MinUint(depth, uint(len(asks)))])
            } else {
                channel := make(chan concurrentOrderBookResponse)
                go getOrderBookSnapshot(httpClient, assetPair, assetPairTranslator, selectLimit(depth), channel)
                select {
                case resp := <- channel:
                    concurrentOrderBook.FilterAndMerge(resp.concurrentOrderBook, true)
                }
            }
        }
    }
}

func filterUpdateId(updateId uint) func(types.OrderBookEntry) bool {
    return func(orderBookEntry types.OrderBookEntry) bool {
        return orderBookEntry.UpdateId > updateId
    }
}

func (b *BinanceUSOrderBookRecorder) GetOrderBook(assetPair types.AssetPair) (types.OrderBook, bool) {
    result, ok := b.orderBooks.Load(assetPair)
    if !ok {
        return types.OrderBook{}, false
    }
    return result.(*util.ConcurrentOrderBook).Data(), true
}

func (b *BinanceUSOrderBookRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := b.orderBooks.Load(assetPair); ok {
        return
    }

    streamName := strings.ToLower(b.assetPairTranslator[assetPair]) + "@depth"

    payloadJson, err := json.Marshal(binanceUSSubscriptionMessage{
        Method: "SUBSCRIBE",
        Params: []string{streamName},
        Id: b.id,
    })
    if err != nil {
        log.Fatalln(err)
    }
    b.Lock()
    defer b.Unlock()
    b.webSocketConnection.WriteMessage(1, payloadJson)
    var resp map[string]interface{}
    for {
        err := b.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if _, ok := resp["code"]; ok {
            log.Fatalln(resp)
        }
        if id, ok := resp["id"]; ok {
            if uint(id.(float64)) != b.id {
                log.Fatalf("id mismatch between sent %v and received %v\n", b.id, id)
            }
            b.id++
            break
        }
        streamName := resp["stream"].(string)
        channel, ok := b.channels.Load(streamName)
        if !ok {
            log.Fatalf("channel not found for streamName %v\n", streamName)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
    }
    channel := make(chan map[string]interface{})

    b.channels.Store(streamName, channel)

    snapshotChannel := make(chan concurrentOrderBookResponse)
    go getOrderBookSnapshot(b.httpClient, assetPair, b.assetPairTranslator, selectLimit(b.depth), snapshotChannel)
    select {
    case resp := <- snapshotChannel:
        b.orderBooks.Store(assetPair, resp.concurrentOrderBook)
        go processOrderBookUpdates(b.httpClient, assetPair, b.assetPairTranslator, resp.concurrentOrderBook, channel, b.depth)
    }
}
