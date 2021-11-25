package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/denali-capital/grizzly/types"
	"github.com/denali-capital/grizzly/util"
)

// docs: https://github.com/binance-us/binance-official-api-docs/blob/master/rest-api.md
const BinanceUSEndpoint string = "https://api.binance.us"

type BinanceUS struct {
	AssetPairTranslator      types.AssetPairTranslator

    apiKey                   string
	secretKey                string
	spreadRecorder			 *util.SpreadRecorder
	latencyEstimator         *util.EwmaEstimator
	orderIdToOrderTranslator *util.ConcurrentOrderIdToOrderPtrMap
	// add timeouts
	httpClient               *http.Client
}

func NewBinanceUS(apiKey, secretKey string, assetPairTranslator types.AssetPairTranslator) *BinanceUS {
	assetPairs := make([]types.AssetPair, len(assetPairTranslator))
	i := 0
	for assetPair := range assetPairTranslator {
		assetPairs[i] = assetPair
		i++
	}
	binanceUS := &BinanceUS{
		AssetPairTranslator: assetPairTranslator,
		apiKey: apiKey,
		secretKey: secretKey,
		latencyEstimator: util.NewEwmaEstimator(0.125, 0.25, 4),
		orderIdToOrderTranslator: util.NewConcurrentOrderIdToOrderPtrMap(),
		httpClient: &http.Client{},
	}
	// parameterize capacity and period?
	binanceUS.spreadRecorder = util.NewSpreadRecorder(
		assetPairs,
		func(assetPair types.AssetPair)types.Spread{
			spread := binanceUS.GetCurrentSpread(assetPair)
			currentTime := time.Now()
			spread.Timestamp = &currentTime
			return spread
		},
		200,
		time.Duration(200) * time.Millisecond,
	)
	return binanceUS
}

func (b *BinanceUS) String() string {
	return "BinanceUS"
}

func (b *BinanceUS) checkError(bodyJson map[string]interface{}) {
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

	mostRecentTimestamp := rawHistoricalSpreads[len(rawHistoricalSpreads) - 1].Timestamp

	// make timestamp sample list
	// nanoseconds per sample
	period := time.Duration(duration.Nanoseconds() / int64(samples))

	// check enough samples exist
	if leastRecentTimestamp := rawHistoricalSpreads[0].Timestamp; mostRecentTimestamp.Add(time.Duration(-(samples - 1)) * period).Before(*leastRecentTimestamp) {
		log.Println("warning: duration is too long, using longest possible duration instead")
		period = time.Duration(mostRecentTimestamp.Sub(*leastRecentTimestamp).Nanoseconds() / int64(samples))
	}

	timestamps := make([]time.Time, samples)
	for i := uint(0); i < samples; i++ {
		timestamps[i] = mostRecentTimestamp.Add(time.Duration(-(samples - i - 1)) * period)
	}

	// get spreads according to sample list
	historicalSpreads := make([]types.Spread, samples)
	currentTimestampIndex := 0
	for i, spread := range rawHistoricalSpreads {
		for (timestamps[currentTimestampIndex].After(*spread.Timestamp) || timestamps[currentTimestampIndex].Equal(*spread.Timestamp)) && (i + 1 >= len(rawHistoricalSpreads) || timestamps[currentTimestampIndex].Before(*rawHistoricalSpreads[i + 1].Timestamp)) {
			effectiveSpread := spread
			effectiveSpread.Timestamp = &timestamps[currentTimestampIndex]
			historicalSpreads[currentTimestampIndex] = effectiveSpread
			currentTimestampIndex++
			if currentTimestampIndex == len(timestamps) {
				break
			}
		}
		if currentTimestampIndex == len(timestamps) {
			break
		}
	}

	channel <- types.SpreadResponse{assetPair, historicalSpreads}
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
	bodyJson := util.HttpGetAndGetBody(b.httpClient, util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/ticker/bookTicker", url.Values{
		"symbol": []string{b.AssetPairTranslator[assetPair]},
	}))
	b.checkError(bodyJson)

	bid, err := strconv.ParseFloat(bodyJson["bidPrice"].(string), 64)
	if err != nil {
		log.Fatalln(err)
	}
	ask, err := strconv.ParseFloat(bodyJson["askPrice"].(string), 64)
	if err != nil {
		log.Fatalln(err)
	}

	return types.Spread{
		Bid: bid,
		Ask: ask,
	}
}

func (b *BinanceUS) getOrderBook(assetPair types.AssetPair, channel chan types.OrderBookResponse) {
	bodyJson := util.HttpGetAndGetBody(b.httpClient, util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/depth", url.Values{
		"symbol": []string{b.AssetPairTranslator[assetPair]},
	}))
	b.checkError(bodyJson)

	orderBook := types.OrderBook{}
	for _, rawOrderBookEntry := range bodyJson["asks"].([]interface{}) {
		price, err := strconv.ParseFloat(rawOrderBookEntry.([]interface{})[0].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		quantity, err := strconv.ParseFloat(rawOrderBookEntry.([]interface{})[1].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		orderBook.Asks = append(orderBook.Asks, types.OrderBookEntry{
			Price: price,
			Quantity: quantity,
		})
	}
	for _, rawOrderBookEntry := range bodyJson["bids"].([]interface{}) {
		price, err := strconv.ParseFloat(rawOrderBookEntry.([]interface{})[0].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		quantity, err := strconv.ParseFloat(rawOrderBookEntry.([]interface{})[1].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		orderBook.Bids = append(orderBook.Bids, types.OrderBookEntry{
			Price: price,
			Quantity: quantity,
		})
	}
	channel <- types.OrderBookResponse{assetPair, orderBook}
}

func (b *BinanceUS) GetOrderBooks(assetPairs []types.AssetPair) map[types.AssetPair]types.OrderBook {
	channel := make(chan types.OrderBookResponse)
	for _, assetPair := range assetPairs {
		go b.getOrderBook(assetPair, channel)
	}

	orderBooks := make(map[types.AssetPair]types.OrderBook)
	for i := 0; i < len(assetPairs); i++ {
		response := <- channel
		orderBooks[response.AssetPair] = response.OrderBook
	}
	return orderBooks
}

func (b *BinanceUS) GetLatency() time.Duration {
    start := time.Now()

    bodyJson := util.HttpGetAndGetBody(b.httpClient, BinanceUSEndpoint + "/api/v3/ping")
    b.checkError(bodyJson)

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
		"price": []string{strconv.FormatFloat(order.Price, 'f', -1, 64)},
		"quantity": []string{strconv.FormatFloat(order.Quantity, 'f', -1, 64)},
		"timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}

	signature := b.getBinanceUSSignature(queryParams)

	request, err := http.NewRequest("POST", util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-MBX-APIKEY", b.apiKey)

	bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
	b.checkError(bodyJson)

    id := bodyJson["orderId"].(uint)

    b.orderIdToOrderTranslator.Store(types.OrderId(id), &order)

    channel <- types.OrderIdResponse{order, types.OrderId(id)}
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
		log.Fatalln("order with id %v not found", orderId)
	}
	queryParams := url.Values{
		"symbol": []string{b.AssetPairTranslator[order.AssetPair]},
		"orderId": []string{string(orderId)},
		"timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}

	signature := b.getBinanceUSSignature(queryParams)

	request, err := http.NewRequest("GET", util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-MBX-APIKEY", b.apiKey)

	bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
	b.checkError(bodyJson)

    orderStatus := types.OrderStatus{
    	Original: order,
    }

    switch status := bodyJson["status"].(string); status {
    case "NEW":
    	orderStatus.Status = types.Unfilled
    case "PARTIALLY_FILLED":
    	price, err := strconv.ParseFloat(bodyJson["price"].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		quantity, err := strconv.ParseFloat(bodyJson["executedQty"].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
    	orderStatus.Status = types.PartiallyFilled
    	orderStatus.FilledPrice = &price
    	orderStatus.FilledQuantity = &quantity
    case "FILLED":
    	price, err := strconv.ParseFloat(bodyJson["price"].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		quantity, err := strconv.ParseFloat(bodyJson["executedQty"].(string), 64)
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
    	log.Fatalln("order %v was rejected by BinanceUS", *order)
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
		log.Fatalln("order with id %v not found", orderId)
	}
	queryParams := url.Values{
		"symbol": []string{b.AssetPairTranslator[order.AssetPair]},
		"orderId": []string{string(orderId)},
		"timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}

	signature := b.getBinanceUSSignature(queryParams)

	request, err := http.NewRequest("DELETE", util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/order", queryParams) + "&signature=" + signature, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-MBX-APIKEY", b.apiKey)

	bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
	b.checkError(bodyJson)

	b.orderIdToOrderTranslator.Delete(orderId)
}

func (b *BinanceUS) CancelOrders(orderIds []types.OrderId) {
	for _, orderId := range orderIds {
		go b.cancelOrder(orderId)
	}
}

func (b *BinanceUS) GetBalances() map[types.Asset]float64 {
	queryParams := url.Values{
		"timestamp": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}

	signature := b.getBinanceUSSignature(queryParams)

	request, err := http.NewRequest("GET", util.ParseUrlWithQuery(BinanceUSEndpoint + "/api/v3/account", queryParams) + "&signature=" + signature, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-MBX-APIKEY", b.apiKey)

	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))

	bodyJson := util.DoHttpAndGetBody(b.httpClient, request)
	b.checkError(bodyJson)

	balances := make(map[types.Asset]float64)
	for _, rawData := range bodyJson["balances"].([]interface{}) {
		data := rawData.(map[string]string)
		free, err := strconv.ParseFloat(data["free"], 64)
		if err != nil {
			log.Fatalln(err)
		}
		locked, err := strconv.ParseFloat(data["locked"], 64)
		if err != nil {
			log.Fatalln(err)
		}
		balances[types.Asset(data["asset"])] = free + locked
	}
	return balances
}
