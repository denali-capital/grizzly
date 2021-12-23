package types

type SpreadResponse struct {
    AssetPair         AssetPair
    HistoricalSpreads []Spread
}

type OrderBookResponse struct {
    AssetPair AssetPair
    OrderBook *OrderBook
}

type OrderIdResponse struct {
    Order   Order
    OrderId OrderId
}

type OrderStatusResponse struct {
    OrderId     OrderId
    OrderStatus OrderStatus
}

type PredictionResponse struct {
    Index       uint
    Prediction  float32
}
