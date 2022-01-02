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
    CheckError(bodyJson)

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
    var resp map[string]interface{}
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

    reverseAssetPairTranslator := util.ReverseAssetPairTranslator(assetPairTranslator)
    assetPairNames := make([]string, len(assetPairs))
    for i, assetPair := range assetPairs {
        streamName := assetPairTranslator[assetPair]
        assetPairNames[i] = streamName
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
    for _, assetPairName := range assetPairNames {
        channel := make(chan map[string]interface{})
        historicalSpread := util.NewConcurrentFixedSizeSpreadQueue(capacity)

        channels.Store("/market/ticker:" + assetPairName, channel)
        historicalSpreads.Store(reverseAssetPairTranslator[assetPairName], historicalSpread)

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
            timestampTime := time.UnixMilli(int64(rawSpread["time"].(float64)))

            historicalSpread.Push(types.Spread{
                Bid: bid,
                Ask: ask,
                Timestamp: &timestampTime,
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
    var resp map[string]interface{}
    for {
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
