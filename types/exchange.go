package types

import "time"

type Exchange interface {
    // * exchange specific information
    String() string

    // * getting data
    GetHistoricalSpreads(assetPairs []AssetPair, duration time.Duration, samples uint) map[AssetPair][]Spread // WebSocket
    GetCurrentSpread(assetPair AssetPair) Spread
    GetOrderBooks(assetPairs []AssetPair) map[AssetPair]*OrderBook // WebSocket
    GetLatency() time.Duration

    // * deal with orders
    ExecuteOrders(orders []Order) map[Order]OrderId
    GetOrderStatuses(orderIds []OrderId) map[OrderId]OrderStatus
    CancelOrders(orderIds []OrderId)

    // * getting account info
    GetBalances() map[Asset]float64
}

// add closing?
type AssetPairRecorder interface {
    RegisterAssetPair(assetPair AssetPair)
}

type SpreadRecorder interface {
    AssetPairRecorder
    GetHistoricalSpreads(assetPair AssetPair) ([]Spread, bool)
}

type OrderBookRecorder interface {
    AssetPairRecorder
    GetOrderBook(assetPair AssetPair) (OrderBook, bool)
}
