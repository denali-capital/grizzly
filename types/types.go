package types

import "time"

type Spread struct {
    Bid       float64
    Ask       float64
    Timestamp *time.Time
}

type OrderType uint

const (
    Buy OrderType = iota
    Sell
)

type Asset string

type AssetPair uint

type Order struct {
    OrderType OrderType
    AssetPair AssetPair
    Price     float64
    Quantity  float64
}

type OrderId string

type StatusType uint

const (
    Pending StatusType = iota
    Unfilled
    PartiallyFilled
    Filled
    Canceled
    Expired
)

type OrderStatus struct {
    Status         StatusType
    FilledPrice    *float64
    FilledQuantity *float64
    Original       *Order
}

type OrderBook struct {
    Bids []OrderBookEntry
    Asks []OrderBookEntry
}

type OrderBookEntry struct {
    Price    float64
    Quantity float64
}

type Exchange interface {
    // * exchange specific information
    String() string

    // * getting data
    GetHistoricalSpreads(assetPairs []AssetPair, duration time.Duration, samples uint) map[AssetPair][]Spread
    GetCurrentSpread(assetPair AssetPair) Spread
    GetOrderBooks(assetPairs []AssetPair) map[AssetPair]OrderBook
    GetLatency() time.Duration

    // * deal with orders
    ExecuteOrders(orders []Order) map[Order]OrderId
    GetOrderStatuses(orderIds []OrderId) map[OrderId]OrderStatus
    CancelOrders(orderIds []OrderId)

    // * getting account info
    GetBalances() map[Asset]float64
}
