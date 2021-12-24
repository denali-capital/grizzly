package kraken

import (
    "bytes"
    "encoding/json"
    "fmt"
    "hash/crc32"
    "log"
    "math"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/gorilla/websocket"
    "github.com/shopspring/decimal"
)

// docs: https://docs.kraken.com/websockets
const WebSocketEndpoint string = "wss://ws.kraken.com"

var Heartbeat []byte = []byte{123, 34, 101, 118, 101, 110, 116, 34, 58, 34, 104, 101, 97, 114, 116, 98, 101, 97, 116, 34, 125}

type krakenSubscription struct {
    Name  string `json:"name"`
    Depth uint   `json:"depth,omitempty"`
}

type krakenSubscriptionMessage struct {
    Event        string             `json:"event"`
    Pair         []string           `json:"pair"`
    Subscription krakenSubscription `json:"subscription"`
}

func initializeWebSocketConnection() *websocket.Conn {
    webSocketConnection, _, err := websocket.DefaultDialer.Dial(WebSocketEndpoint, http.Header{})
    if err != nil {
        log.Fatalln(err)
    }

    var initialResponse map[string]interface{}
    err = webSocketConnection.ReadJSON(&initialResponse)
    if err != nil {
        log.Fatalln(err)
    }
    if !(initialResponse["event"].(string) == "systemStatus" && initialResponse["status"].(string) == "online") {
        log.Fatalln(initialResponse)
    }

    return webSocketConnection
}

// should do separate connection for each asset pair?
type krakenWebSocketRecorder struct {
    sync.Mutex
    webSocketConnection *websocket.Conn
    iso4217Translator   types.AssetPairTranslator
    // map[uint]chan []interface{}
    channels        *sync.Map
}

func (k *krakenWebSocketRecorder) record() {
    var resp []interface{}
    for {
        k.Lock()
        _, msg, err := k.webSocketConnection.ReadMessage()
        if err != nil {
            log.Fatalln(err)
        } else if bytes.Compare(Heartbeat, msg) != 0 {
            err = json.Unmarshal(msg, &resp)
            if err != nil {
                log.Fatalln(err)
            }
            channelId := uint(resp[0].(float64))
            channel, ok := k.channels.Load(channelId)
            if !ok {
                log.Fatalf("channel not found for channelId %v\n", channelId)
            }
            channel.(chan []interface{}) <- util.Copy(resp)
        }
        k.Unlock()
    }
}

type KrakenSpreadRecorder struct {
    krakenWebSocketRecorder
    capacity                 uint
    // map[types.AssetPair]*util.FixedSizeSpreadQueue
    historicalSpreads        *sync.Map
}

func NewKrakenSpreadRecorder(assetPairs []types.AssetPair, iso4217Translator types.AssetPairTranslator, capacity uint) *KrakenSpreadRecorder {
    webSocketConnection := initializeWebSocketConnection()

    reverseIso4217Translator := util.ReverseAssetPairTranslator(iso4217Translator)
    iso4217TranslatedPairs := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        iso4217TranslatedPairs[i] = iso4217Translator[assetPair]
    }

    payloadJson, err := json.Marshal(krakenSubscriptionMessage{
        Event: "subscribe",
        Pair: iso4217TranslatedPairs,
        Subscription: krakenSubscription{
            Name: "spread",
        },
    })
    if err != nil {
        log.Fatalln(err)
    }
    webSocketConnection.WriteMessage(1, payloadJson)
    channelIdTranslator := make(map[uint]types.AssetPair)
    channels := &sync.Map{}
    historicalSpreads := &sync.Map{}
    var initialResponse map[string]interface{}
    for i := 0; i < len(iso4217TranslatedPairs); i++ {
        err = webSocketConnection.ReadJSON(&initialResponse)
        if err != nil {
            log.Fatalln(err)
        }
        if !(initialResponse["event"].(string) == "subscriptionStatus" && util.Contains(iso4217TranslatedPairs, initialResponse["pair"].(string)) && initialResponse["status"].(string) == "subscribed") {
            // we assume here that all subscription messages come one right after another
            log.Fatalln(initialResponse)
        }
        channelId := uint(initialResponse["channelID"].(float64))
        assetPair := reverseIso4217Translator[initialResponse["pair"].(string)]
        channelIdTranslator[channelId] = assetPair
        channels.Store(channelId, make(chan []interface{}))
        historicalSpreads.Store(assetPair, util.NewFixedSizeSpreadQueue(capacity))
    }

    krakenSpreadRecorder := &KrakenSpreadRecorder{
        krakenWebSocketRecorder: krakenWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            iso4217Translator: iso4217Translator,
            channels: channels,
        },
        capacity: capacity,
        historicalSpreads: historicalSpreads,
    }

    channels.Range(func(rawChannelId, rawChannel interface{}) bool {
        go krakenSpreadRecorder.processSpreadUpdates(rawChannelId.(uint), channelIdTranslator[rawChannelId.(uint)], rawChannel.(chan []interface{}))
        return true
    })
    go krakenSpreadRecorder.record()

    return krakenSpreadRecorder
}

func (k *KrakenSpreadRecorder) processSpreadUpdates(channelId uint, assetPair types.AssetPair, channel chan []interface{}) {
    for {
        select {
        case resp := <- channel:
            rawSpread := resp[1].([]interface{})
            bid, err := decimal.NewFromString(rawSpread[0].(string))
            if err != nil {
                log.Fatalln(err)
            }
            ask, err := decimal.NewFromString(rawSpread[1].(string))
            if err != nil {
                log.Fatalln(err)
            }
            timestamp, err := strconv.ParseFloat(rawSpread[2].(string), 64)
            if err != nil {
                log.Fatalln(err)
            }
            timestampInteger, timestampFraction := math.Modf(timestamp)
            // fraction is in ns
            timestampTime := time.Unix(int64(timestampInteger), int64(timestampFraction * 1000000))

            historicalSpread, ok := k.historicalSpreads.Load(assetPair)
            if !ok {
                log.Fatalf("historicalSpreads not found for assetPair %v\n", assetPair)
            }
            historicalSpread.(*util.FixedSizeSpreadQueue).Push(types.Spread{
                Bid: bid,
                Ask: ask,
                Timestamp: &timestampTime,
            })
        }
    }
}

func (k *KrakenSpreadRecorder) GetHistoricalSpreads(assetPair types.AssetPair) ([]types.Spread, bool) {
    result, ok := k.historicalSpreads.Load(assetPair)
    if !ok {
        return make([]types.Spread, 0), false
    }
    return result.(*util.FixedSizeSpreadQueue).Data(), true
}

func (k *KrakenSpreadRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := k.historicalSpreads.Load(assetPair); ok {
        return
    }

    iso4217TranslatedPair := k.iso4217Translator[assetPair]

    payloadJson, err := json.Marshal(krakenSubscriptionMessage{
        Event: "subscribe",
        Pair: []string{iso4217TranslatedPair},
        Subscription: krakenSubscription{
            Name: "spread",
        },
    })
    if err != nil {
        log.Fatalln(err)
    }
    k.Lock()
    defer k.Unlock()
    k.webSocketConnection.WriteMessage(1, payloadJson)
    var initialResponse map[string]interface{}
    var resp []interface{}
    for {
        _, msg, err := k.webSocketConnection.ReadMessage()
        if err != nil {
            log.Fatalln(err)
        }
        if bytes.Compare(Heartbeat, msg) != 0 {
            // not a heartbeat
            err := json.Unmarshal(msg, &initialResponse)
            if err != nil {
                _, ok := err.(*json.UnmarshalTypeError)
                if !ok {
                    log.Fatalln(err)
                }
                err = json.Unmarshal(msg, &resp)
                if err != nil {
                    log.Fatalln(err)
                }
                channelId := uint(resp[0].(float64))
                channel, ok := k.channels.Load(channelId)
                if !ok {
                    log.Fatalf("channel not found for channelId %v", channelId)
                }
                channel.(chan []interface{}) <- util.Copy(resp)
                continue
            }
            if !(initialResponse["event"].(string) == "subscriptionStatus" && initialResponse["pair"].(string) == iso4217TranslatedPair && initialResponse["status"].(string) == "subscribed") {
                log.Fatalln(initialResponse)
            }
            break
        }
    }

    channelId := uint(initialResponse["channelID"].(float64))
    channel := make(chan []interface{})
    k.channels.Store(channelId, channel)
    k.historicalSpreads.Store(assetPair, util.NewFixedSizeSpreadQueue(k.capacity))
    go k.processSpreadUpdates(channelId, assetPair, channel)
}

// very much inspired by https://github.com/jurijbajzelj/kraken_ws_orderbook
type KrakenOrderBookRecorder struct {
    krakenWebSocketRecorder
    depth                    uint
    // map[types.AssetPair]*types.OrderBook
    orderBooks        *sync.Map
}

func NewKrakenOrderBookRecorder(assetPairs []types.AssetPair, iso4217Translator types.AssetPairTranslator, depth uint) *KrakenOrderBookRecorder {
    webSocketConnection := initializeWebSocketConnection()

    reverseIso4217Translator := util.ReverseAssetPairTranslator(iso4217Translator)
    iso4217TranslatedPairs := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        iso4217TranslatedPairs[i] = iso4217Translator[assetPair]
    }

    payloadJson, err := json.Marshal(krakenSubscriptionMessage{
        Event: "subscribe",
        Pair: iso4217TranslatedPairs,
        Subscription: krakenSubscription{
            Name: "book",
            Depth: depth,
        },
    })
    if err != nil {
        log.Fatalln(err)
    }
    webSocketConnection.WriteMessage(1, payloadJson)

    channelIdTranslator := make(map[uint]types.AssetPair)
    channels := &sync.Map{}
    var initialResponse map[string]interface{}
    for i := 0; i < len(iso4217TranslatedPairs); i++ {
        err = webSocketConnection.ReadJSON(&initialResponse)
        if err != nil {
            log.Fatalln(err)
        }
        if !(initialResponse["event"].(string) == "subscriptionStatus" && util.Contains(iso4217TranslatedPairs, initialResponse["pair"].(string)) && initialResponse["status"].(string) == "subscribed") {
            // we assume here that all subscription messages come one right after another
            log.Fatalln(initialResponse)
        }
        channelId := uint(initialResponse["channelID"].(float64))
        assetPair := reverseIso4217Translator[initialResponse["pair"].(string)]
        channelIdTranslator[channelId] = assetPair
        channels.Store(channelId, make(chan []interface{}))
    }

    // get initial books
    orderBooks := &sync.Map{}
    var resp []interface{}
    for i := 0; i < len(iso4217TranslatedPairs); i++ {
        for {
            _, msg, err := webSocketConnection.ReadMessage()
            if err != nil {
                log.Fatalln(err)
            }
            if bytes.Compare(Heartbeat, msg) != 0 {
                // not a hearbeat
                err := json.Unmarshal(msg, &resp)
                if err != nil {
                    log.Fatalln(err)
                }
                break
            }
        }

        // we assume here that all initial book messages come right after another
        orderBook := &types.OrderBook{}
        channelId := uint(resp[0].(float64))
        rawOrderBook := resp[1].(map[string]interface{})
        for _, rawOrderBookEntry := range rawOrderBook["as"].([]interface{}) {
            price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
            orderBook.Asks = append(orderBook.Asks, types.OrderBookEntry{
                Price: price,
                Quantity: quantity,
            })
            if uint(len(orderBook.Asks)) == depth {
                break
            }
        }
        for _, rawOrderBookEntry := range rawOrderBook["bs"].([]interface{}) {
            price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
            orderBook.Bids = append(orderBook.Bids, types.OrderBookEntry{
                Price: price,
                Quantity: quantity,
            })
            if uint(len(orderBook.Bids)) == depth {
                break
            }
        }

        orderBooks.Store(channelIdTranslator[channelId], orderBook)
    }

    krakenOrderBookRecorder := &KrakenOrderBookRecorder{
        krakenWebSocketRecorder: krakenWebSocketRecorder{
            webSocketConnection: webSocketConnection,
            iso4217Translator: iso4217Translator,
            channels: channels,
        },
        depth: depth,
        orderBooks: orderBooks,
    }

    channels.Range(func(rawChannelId, rawChannel interface{}) bool {
        go krakenOrderBookRecorder.processOrderBookUpdates(rawChannelId.(uint), channelIdTranslator[rawChannelId.(uint)], rawChannel.(chan []interface{}))
        return true
    })
    go krakenOrderBookRecorder.record()

    return krakenOrderBookRecorder
}

func preFormatDecimal(val decimal.Decimal) string {
    one := decimal.NewFromInt(1)
    ten := decimal.NewFromInt(10)

    offset := 0
    tmp := val
    for tmp.LessThan(one) {
        tmp = tmp.Mul(ten)
        offset++
    }
    return val.StringFixed(int32(val.NumDigits() - len(val.Truncate(0).String()) + offset))
}

func getChecksumInput(bids []types.OrderBookEntry, asks []types.OrderBookEntry) string {
    var str strings.Builder
    for _, orderBookEntry := range asks[:10] {
        price := preFormatDecimal(orderBookEntry.Price)
        price = strings.Replace(price, ".", "", 1)
        price = strings.TrimLeft(price, "0")
        str.WriteString(price)

        quantity := preFormatDecimal(orderBookEntry.Quantity)
        quantity = strings.Replace(quantity, ".", "", 1)
        quantity = strings.TrimLeft(quantity, "0")
        str.WriteString(quantity)
    }
    for _, orderBookEntry := range bids[:10] {
        price := preFormatDecimal(orderBookEntry.Price)
        price = strings.Replace(price, ".", "", 1)
        price = strings.TrimLeft(price, "0")
        str.WriteString(price)

        quantity := preFormatDecimal(orderBookEntry.Quantity)
        quantity = strings.Replace(quantity, ".", "", 1)
        quantity = strings.TrimLeft(quantity, "0")
        str.WriteString(quantity)
    }
    return str.String()
}

func verifyOrderBookChecksum(bids []types.OrderBookEntry, asks []types.OrderBookEntry, checksum string) {
    checksumInput := getChecksumInput(bids, asks)
    crc := crc32.ChecksumIEEE([]byte(checksumInput))
    if fmt.Sprint(crc) != checksum {
        log.Fatalln("order book checksum not the same ", " ", crc, " ", checksum)
    }
}

func (k *KrakenOrderBookRecorder) processOrderBookUpdates(channelId uint, assetPair types.AssetPair, channel chan []interface{}) {
    for {
        select {
        case resp := <- channel:
            orderBook, ok := k.orderBooks.Load(assetPair)
            if !ok {
                log.Fatalf("orderBook not found for assetPair %v", assetPair)
            }
            bids := orderBook.(*types.OrderBook).Bids
            asks := orderBook.(*types.OrderBook).Asks
            if len(resp) == 4 {
                // one of bids or asks is updated
                orderBookDiff := resp[1].(map[string]interface{})
                checksum := orderBookDiff["c"].(string)

                if val, ok := orderBookDiff["b"]; ok {
                    for _, rawOrderBookEntry := range val.([]interface{}) {
                        price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                        if quantity.Equal(decimal.Zero) {
                            bids = util.RemovePriceFromBids(bids, price)
                        } else {
                            if len(rawOrderBookEntry.([]interface{})) == 4 {
                                // it has the 4th element "r" so we just re-append
                                bids = append(bids, types.OrderBookEntry{
                                    Price: price,
                                    Quantity: quantity,
                                })
                            } else {
                                bids = util.InsertPriceInBids(bids, price, quantity)
                                bids = bids[:k.depth]
                            }
                        }
                    }
                } else {
                    for _, rawOrderBookEntry := range orderBookDiff["a"].([]interface{}) {
                        price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                        if quantity.Equal(decimal.Zero) {
                            asks = util.RemovePriceFromAsks(asks, price)
                        } else {
                            if len(rawOrderBookEntry.([]interface{})) == 4 {
                                // it has the 4th element "r" so we just re-append
                                asks = append(asks, types.OrderBookEntry{
                                    Price: price,
                                    Quantity: quantity,
                                })
                            } else {
                                asks = util.InsertPriceInAsks(asks, price, quantity)
                                asks = asks[:k.depth]
                            }
                        }
                    }
                }
                verifyOrderBookChecksum(bids, asks, checksum)
            } else {
                // both bids and asks are updated
                orderBookDiffAsks := resp[1].(map[string]interface{})
                orderBookDiffBids := resp[2].(map[string]interface{})
                checksum := orderBookDiffBids["c"].(string)

                for _, rawOrderBookEntry := range orderBookDiffBids["b"].([]interface{}) {
                    price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                    if quantity.Equal(decimal.Zero) {
                        bids = util.RemovePriceFromBids(bids, price)
                    } else {
                        if len(rawOrderBookEntry.([]interface{})) == 4 {
                            // it has the 4th element "r" so we just re-append
                            bids = append(bids, types.OrderBookEntry{
                                Price: price,
                                Quantity: quantity,
                            })
                        } else {
                            bids = util.InsertPriceInBids(bids, price, quantity)
                            bids = bids[:k.depth]
                        }
                    }
                }
                for _, rawOrderBookEntry := range orderBookDiffAsks["a"].([]interface{}) {
                    price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
                    if quantity.Equal(decimal.Zero) {
                        asks = util.RemovePriceFromAsks(asks, price)
                    } else {
                        if len(rawOrderBookEntry.([]interface{})) == 4 {
                            // it has the 4th element "r" so we just re-append
                            asks = append(asks, types.OrderBookEntry{
                                Price: price,
                                Quantity: quantity,
                            })
                        } else {
                            asks = util.InsertPriceInAsks(asks, price, quantity)
                            asks = asks[:k.depth]
                        }
                    }
                }
                verifyOrderBookChecksum(bids, asks, checksum)
            }
            k.orderBooks.Store(assetPair, &types.OrderBook{
                Bids: bids[:k.depth],
                Asks: asks[:k.depth],
            })
        }
    }
}

func (k *KrakenOrderBookRecorder) GetOrderBook(assetPair types.AssetPair) (types.OrderBook, bool) {
    orderBook, ok := k.orderBooks.Load(assetPair)
    if !ok {
        return types.OrderBook{}, false
    }
    return *orderBook.(*types.OrderBook), true
}

func (k *KrakenOrderBookRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := k.orderBooks.Load(assetPair); ok {
        return
    }

    iso4217TranslatedPair := k.iso4217Translator[assetPair]

    payloadJson, err := json.Marshal(krakenSubscriptionMessage{
        Event: "subscribe",
        Pair: []string{iso4217TranslatedPair},
        Subscription: krakenSubscription{
            Name: "book",
            Depth: k.depth,
        },
    })
    if err != nil {
        log.Fatalln(err)
    }
    k.Lock()
    defer k.Unlock()
    k.webSocketConnection.WriteMessage(1, payloadJson)
    var initialResponse map[string]interface{}
    var resp []interface{}
    for {
        _, msg, err := k.webSocketConnection.ReadMessage()
        if err != nil {
            log.Fatalln(err)
        }
        if bytes.Compare(Heartbeat, msg) != 0 {
            // not a heartbeat
            err := json.Unmarshal(msg, &initialResponse)
            if err != nil {
                _, ok := err.(*json.UnmarshalTypeError)
                if !ok {
                    log.Fatalln(err)
                }
                err = json.Unmarshal(msg, &resp)
                if err != nil {
                    log.Fatalln(err)
                }
                channelId := uint(resp[0].(float64))
                channel, ok := k.channels.Load(channelId)
                if !ok {
                    log.Fatalf("channel not found for channelId %v", channelId)
                }
                channel.(chan []interface{}) <- util.Copy(resp)
                continue
            }
            if !(initialResponse["event"].(string) == "subscriptionStatus" && initialResponse["pair"].(string) == iso4217TranslatedPair && initialResponse["status"].(string) == "subscribed") {
                log.Fatalln(initialResponse)
            }
            break
        }
    }

    channelId := uint(initialResponse["channelID"].(float64))
    channel := make(chan []interface{})
    k.channels.Store(channelId, channel)

    // get initial book
    for {
        _, msg, err := k.webSocketConnection.ReadMessage()
        if err != nil {
            log.Fatalln(err)
        }
        if bytes.Compare(Heartbeat, msg) != 0 {
            // not a hearbeat
            err := json.Unmarshal(msg, &resp)
            if err != nil {
                log.Fatalln(err)
            }
            if messageChannelId := uint(resp[0].(float64)); messageChannelId != channelId {
                channel, ok := k.channels.Load(messageChannelId)
                if !ok {
                    log.Fatalf("channel not found for channelId %v", messageChannelId)
                }
                channel.(chan []interface{}) <- util.Copy(resp)
                continue
            }
            break
        }
    }
    orderBook := &types.OrderBook{}
    rawOrderBook := resp[1].(map[string]interface{})
    for _, rawOrderBookEntry := range rawOrderBook["as"].([]interface{}) {
        price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
        orderBook.Asks = append(orderBook.Asks, types.OrderBookEntry{
            Price: price,
            Quantity: quantity,
        })
        if uint(len(orderBook.Asks)) == k.depth {
            break
        }
    }
    for _, rawOrderBookEntry := range rawOrderBook["bs"].([]interface{}) {
        price, quantity := util.GetPriceAndQuantity(rawOrderBookEntry.([]interface{}))
        orderBook.Bids = append(orderBook.Bids, types.OrderBookEntry{
            Price: price,
            Quantity: quantity,
        })
        if uint(len(orderBook.Bids)) == k.depth {
            break
        }
    }
    k.orderBooks.Store(assetPair, orderBook)
    go k.processOrderBookUpdates(channelId, assetPair, channel)
}
