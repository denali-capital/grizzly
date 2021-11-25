package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/denali-capital/grizzly/types"
	"github.com/denali-capital/grizzly/util"
)

// docs: https://docs.kraken.com/rest/
const KrakenEndpoint string = "https://api.kraken.com"

type Kraken struct {
	AssetPairTranslator      types.AssetPairTranslator

    apiKey                   string
	secretKey                string
	latencyEstimator         *util.EwmaEstimator
	orderIdToOrderTranslator *util.ConcurrentOrderIdToOrderPtrMap
	// add timeouts
	httpClient               *http.Client
}

func NewKraken(apiKey, secretKey string, assetPairTranslator types.AssetPairTranslator) *Kraken {
	return &Kraken{
		AssetPairTranslator: assetPairTranslator,
		apiKey: apiKey,
		secretKey: secretKey,
		latencyEstimator: util.NewEwmaEstimator(0.125, 0.25, 4),
		orderIdToOrderTranslator: util.NewConcurrentOrderIdToOrderPtrMap(),
		httpClient: &http.Client{},
	}
}

func (k *Kraken) String() string {
	return "Kraken"
}

func (k *Kraken) checkError(bodyJson map[string]interface{}) {
	errors := bodyJson["error"].([]interface{})
	if len(errors) > 0 {
		log.Fatalln(errors)
	}
}

// pointerize?
func partitionMatch(currentSpreads [][]interface{}, timestamps []float64) []types.Spread {
	spreads := make([]types.Spread, len(timestamps))
	for i, timestamp := range timestamps {
		timestampInteger, timestampFraction := math.Modf(timestamp)
		// fraction is in ns
		timestampTime := time.Unix(int64(timestampInteger), int64(timestampFraction * 1000000))
		spread := currentSpreads[uint(float64(len(currentSpreads)) * timestampFraction)]
		bid, err := strconv.ParseFloat(spread[1].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		ask, err := strconv.ParseFloat(spread[2].(string), 64)
		if err != nil {
			log.Fatalln(err)
		}
		spreads[i] = types.Spread{
			Bid: bid,
			Ask: ask,
			Timestamp: &timestampTime,
		}
	}
	return spreads
}

func (k *Kraken) getHistoricalSpread(assetPair types.AssetPair, duration time.Duration, samples uint, channel chan types.SpreadResponse) {
	if samples == 0 || duration <= 0 {
		channel <- types.SpreadResponse{assetPair, []types.Spread{}}
	}

	bodyJson := util.HttpGetAndGetBody(k.httpClient, util.ParseUrlWithQuery(KrakenEndpoint + "/0/public/Spread", url.Values{
		"pair": []string{k.AssetPairTranslator[assetPair]},
	}))
	k.checkError(bodyJson)

	mostRecentTimestamp := bodyJson["result"].(map[string]interface{})["last"].(float64)
	data := bodyJson["result"].(map[string]interface{})[k.AssetPairTranslator[assetPair]].([]interface{})

	// make timestamp sample list
	// seconds per sample
	period := duration.Seconds() / float64(samples)

	// check enough samples exist
	if leastRecentTimestamp := data[0].([]interface{})[0].(float64); mostRecentTimestamp - float64(samples - 1) * period < leastRecentTimestamp {
		log.Println("warning: duration is too long, using longest possible duration instead")
		period = (mostRecentTimestamp - leastRecentTimestamp) / float64(samples)
	}

	timestamps := make([][]float64, 0)
	sameIntegerTimestamps := make([]float64, 0)
	lastTimestamp := -1.0
	for i := int(samples - 1); i >= 0; i-- {
		currentTimestamp := mostRecentTimestamp - float64(i) * period
		if uint(currentTimestamp) == uint(lastTimestamp) {
			sameIntegerTimestamps = append(sameIntegerTimestamps, currentTimestamp)
		} else {
			if len(sameIntegerTimestamps) > 0 {
				timestamps = append(timestamps, sameIntegerTimestamps)
				sameIntegerTimestamps = make([]float64, 0, cap(sameIntegerTimestamps))
			}
			sameIntegerTimestamps = append(sameIntegerTimestamps, currentTimestamp)
		}
		lastTimestamp = currentTimestamp
	}
	timestamps = append(timestamps, sameIntegerTimestamps)

	// get spreads according to sample list
	currentTimestampIndex := 0
	currentSpreads := make([][]interface{}, 0)
	historicalSpreads := make([]types.Spread, 0, samples)
	for _, rawSpread := range data {
		spread := rawSpread.([]interface{})
		timestamp := uint(spread[0].(float64))
		if timestamp == uint(timestamps[currentTimestampIndex][0]) {
			currentSpreads = append(currentSpreads, spread)
		} else if len(currentSpreads) > 0 {
			historicalSpreads = append(historicalSpreads, partitionMatch(currentSpreads, timestamps[currentTimestampIndex])...)
			currentTimestampIndex++
			currentSpreads = make([][]interface{}, 0, cap(currentSpreads))
			if timestamp == uint(timestamps[currentTimestampIndex][0]) {
				currentSpreads = append(currentSpreads, spread)
			}
		}
	}

	// add last spreads
	if len(currentSpreads) == 0 {
		last := uint(data[len(data) - 1].([]interface{})[0].(float64))
		currentSpreads = append(currentSpreads, data[len(data) - 1].([]interface{}))
		for i := int(len(data) - 2); i >= 0; i-- {
			if last == uint(data[i].([]interface{})[0].(float64)) {
				currentSpreads = append(currentSpreads, data[i].([]interface{}))
			} else {
				break
			}
		}
		timestampDiff := last - uint(timestamps[currentTimestampIndex][0])
		for i := range timestamps[currentTimestampIndex] {
			timestamps[currentTimestampIndex][i] += float64(timestampDiff)
		}
	}
	historicalSpreads = append(historicalSpreads, partitionMatch(currentSpreads, timestamps[currentTimestampIndex])...)
	channel <- types.SpreadResponse{assetPair, historicalSpreads}
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
	bodyJson := util.HttpGetAndGetBody(k.httpClient, util.ParseUrlWithQuery(KrakenEndpoint + "/0/public/Ticker", url.Values{
		"pair": []string{k.AssetPairTranslator[assetPair]},
	}))
	k.checkError(bodyJson)

	data := bodyJson["result"].(map[string]interface{})[k.AssetPairTranslator[assetPair]].(map[string]interface{})

	bid, err := strconv.ParseFloat(data["b"].([]interface{})[0].(string), 64)
	if err != nil {
		log.Fatalln(err)
	}
	ask, err := strconv.ParseFloat(data["a"].([]interface{})[0].(string), 64)
	if err != nil {
		log.Fatalln(err)
	}

	return types.Spread{
		Bid: bid,
		Ask: ask,
	}
}

func (k *Kraken) getOrderBook(assetPair types.AssetPair, channel chan types.OrderBookResponse) {
	bodyJson := util.HttpGetAndGetBody(k.httpClient, util.ParseUrlWithQuery(KrakenEndpoint + "/0/public/Depth", url.Values{
		"pair": []string{k.AssetPairTranslator[assetPair]},
	}))
	k.checkError(bodyJson)

	data := bodyJson["result"].(map[string]interface{})[k.AssetPairTranslator[assetPair]].(map[string]interface{})

	orderBook := types.OrderBook{}
	for _, rawOrderBookEntry := range data["asks"].([]interface{}) {
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
	for _, rawOrderBookEntry := range data["bids"].([]interface{}) {
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

func (k *Kraken) GetOrderBooks(assetPairs []types.AssetPair) map[types.AssetPair]types.OrderBook {
	channel := make(chan types.OrderBookResponse)
	for _, assetPair := range assetPairs {
		go k.getOrderBook(assetPair, channel)
	}

	orderBooks := make(map[types.AssetPair]types.OrderBook)
	for i := 0; i < len(assetPairs); i++ {
		response := <- channel
		orderBooks[response.AssetPair] = response.OrderBook
	}
	return orderBooks
}

func (k *Kraken) GetLatency() time.Duration {
    start := time.Now()

    bodyJson := util.HttpGetAndGetBody(k.httpClient, KrakenEndpoint + "/0/public/Time")
    k.checkError(bodyJson)

    duration := time.Since(start)

    k.latencyEstimator.Sample(float64(duration.Milliseconds()))

    return time.Duration(k.latencyEstimator.GetEstimate()) * time.Millisecond
}

func (k *Kraken) parseOrderType(ot types.OrderType) string {
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
		"type": []string{k.parseOrderType(order.OrderType)},
		"ordertype": []string{"limit"},
		"price": []string{strconv.FormatFloat(order.Price, 'f', -1, 64)},
		"volume": []string{strconv.FormatFloat(order.Quantity, 'f', -1, 64)},
		"nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	request, err := http.NewRequest("POST", KrakenEndpoint + "/0/private/AddOrder", strings.NewReader(queryParams.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	signature := k.getKrakenSignature("/0/private/AddOrder", queryParams)
	request.Header.Set("API-Sign", signature)
	request.Header.Set("API-Key", k.apiKey)

	bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
	k.checkError(bodyJson)

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
	orderIdStrings := make([]string, len(orderIds))
	for i, orderId := range orderIds {
		orderIdStrings[i] = string(orderId)
	}
	queryParams := url.Values{
		"txid": []string{strings.Join(orderIdStrings, ",")},
		"nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	request, err := http.NewRequest("POST", KrakenEndpoint + "/0/private/QueryOrders", strings.NewReader(queryParams.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	signature := k.getKrakenSignature("/0/private/QueryOrders", queryParams)
	request.Header.Set("API-Sign", signature)
	request.Header.Set("API-Key", k.apiKey)

	bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
	k.checkError(bodyJson)

	data := bodyJson["result"].(map[types.OrderId]map[string]interface{})
	orderStatuses := make(map[types.OrderId]types.OrderStatus)
	for id, rawOrderData := range data {
		original, ok := k.orderIdToOrderTranslator.Load(id)
		if !ok {
			log.Fatalln("order with id %v not found", id)
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
			price, err := strconv.ParseFloat(rawOrderData["price"].(string), 64)
			if err != nil {
				log.Fatalln(err)
			}
			quantity, err := strconv.ParseFloat(rawOrderData["vol_exec"].(string), 64)
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

func (k *Kraken) cancelOrder(orderId types.OrderId) {
	queryParams := url.Values{
		"txid": []string{string(orderId)},
		"nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	request, err := http.NewRequest("POST", KrakenEndpoint + "/0/private/CancelOrder", strings.NewReader(queryParams.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	signature := k.getKrakenSignature("/0/private/CancelOrder", queryParams)
	request.Header.Set("API-Sign", signature)
	request.Header.Set("API-Key", k.apiKey)

	bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
	k.checkError(bodyJson)

	k.orderIdToOrderTranslator.Delete(orderId)
}

func (k *Kraken) CancelOrders(orderIds []types.OrderId) {
	for _, orderId := range orderIds {
		go k.cancelOrder(orderId)
	}
}

func (k *Kraken) GetBalances() map[types.Asset]float64 {
	queryParams := url.Values{
		"nonce": []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	request, err := http.NewRequest("POST", KrakenEndpoint + "/0/private/Balance", strings.NewReader(queryParams.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	signature := k.getKrakenSignature("/0/private/Balance", queryParams)
	request.Header.Set("API-Sign", signature)
	request.Header.Set("API-Key", k.apiKey)

	bodyJson := util.DoHttpAndGetBody(k.httpClient, request)
	k.checkError(bodyJson)

	if _, ok := bodyJson["result"]; ok {
		data := bodyJson["result"].(map[types.Asset]string)

		balances := make(map[types.Asset]float64)
		for asset, balanceString := range data {
			balances[asset], err = strconv.ParseFloat(balanceString, 64)
			if err != nil {
				log.Fatalln(err)
			}
		}
		return balances
	}
	return map[types.Asset]float64{}
}
