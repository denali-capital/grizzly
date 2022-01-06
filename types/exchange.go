package types

import (
    "time"

    "github.com/shopspring/decimal"
)

type Exchange interface {
    // * exchange specific information
    String() string

    // * getting data
    GetHistoricalSpreads(assetPairs []AssetPair, duration time.Duration, samples uint) map[AssetPair][]Spread // WebSocket
    GetCurrentSpread(assetPair AssetPair) Spread // Websocket
    GetOrderBooks(assetPairs []AssetPair) map[AssetPair]*OrderBook // WebSocket
    GetLatency() time.Duration

    // * deal with orders
    ExecuteOrders(orders []Order) map[Order]OrderId
    GetOrderStatuses(orderIds []OrderId) map[OrderId]OrderStatus
    CancelOrders(orderIds []OrderId)

    // * getting account info
    GetBalances() map[Asset]decimal.Decimal
}

// add closing?
type AssetPairRecorder interface {
    RegisterAssetPair(assetPair AssetPair)
}

type SpreadRecorder interface {
    AssetPairRecorder
    GetCurrentSpread(assetPair AssetPair) (Spread, bool)
    GetHistoricalSpreads(assetPair AssetPair) ([]Spread, bool)
}

type OrderBookRecorder interface {
    AssetPairRecorder
    GetOrderBook(assetPair AssetPair) (OrderBook, bool)
}
