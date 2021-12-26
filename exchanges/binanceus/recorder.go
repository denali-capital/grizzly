package binanceus

import (
    "encoding/json"
    "log"
    "net/http"
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

func (k *binanceUSWebSocketRecorder) record() {
    var resp map[string]interface{}
    for {
        k.Lock()
        err := k.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if _, ok := resp["code"]; ok {
            log.Fatalln(resp)
        }
        streamName := resp["stream"].(string)
        channel, ok := k.channels.Load(streamName)
        if !ok {
            log.Fatalf("channel not found for streamName %v\n", streamName)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
        k.Unlock()
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
    streamIndicator := make(map[string]interface{})
    for i, assetPair := range assetPairs {
        streamName := strings.ToLower(assetPairTranslator[assetPair]) + "@bookTicker"
        streamTranslator[streamName] = assetPair
        streams[i] = streamName
        streamIndicator[streamName] = nil
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

func (k *BinanceUSSpreadRecorder) GetHistoricalSpreads(assetPair types.AssetPair) ([]types.Spread, bool) {
    result, ok := k.historicalSpreads.Load(assetPair)
    if !ok {
        return make([]types.Spread, 0), false
    }
    return result.(*util.ConcurrentFixedSizeSpreadQueue).Data(), true
}

func (k *BinanceUSSpreadRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    if _, ok := k.historicalSpreads.Load(assetPair); ok {
        return
    }

    streamName := strings.ToLower(k.assetPairTranslator[assetPair]) + "@bookTicker"

    payloadJson, err := json.Marshal(binanceUSSubscriptionMessage{
        Method: "SUBSCRIBE",
        Params: []string{streamName},
        Id: k.id,
    })
    if err != nil {
        log.Fatalln(err)
    }
    k.Lock()
    defer k.Unlock()
    k.webSocketConnection.WriteMessage(1, payloadJson)
    var resp map[string]interface{}
    for {
        err := k.webSocketConnection.ReadJSON(&resp)
        if err != nil {
            log.Fatalln(err)
        }
        if _, ok := resp["code"]; ok {
            log.Fatalln(resp)
        }
        if id, ok := resp["id"]; ok {
            if uint(id.(float64)) != k.id {
                log.Fatalf("id mismatch between sent %v and received %v\n", k.id, id)
            }
            k.id++
            break
        }
        streamName := resp["stream"].(string)
        channel, ok := k.channels.Load(streamName)
        if !ok {
            log.Fatalf("channel not found for streamName %v\n", streamName)
        }
        channel.(chan map[string]interface{}) <- util.MapCopy(resp)
    }
    channel := make(chan map[string]interface{})
    historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(k.capacity)

    k.channels.Store(streamName, channel)
    k.historicalSpreads.Store(assetPair, historicalSpread)

    go processSpreadUpdates(historicalSpread, channel)
}
