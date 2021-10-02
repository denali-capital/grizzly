package types

// pointerize?
type SpreadResponse struct {
	AssetPair         AssetPair
	HistoricalSpreads []Spread
}

// pointerize?
type OrderBookResponse struct {
	AssetPair AssetPair
	OrderBook OrderBook
}

// pointerize?
type OrderIdResponse struct {
	Order   Order
	OrderId OrderId
}

// pointerize?
type OrderStatusResponse struct {
	OrderId     OrderId
	OrderStatus OrderStatus
}
