package types

import "time"

type Spread struct {
    Bid       float64
    Ask       float64
    Timestamp time.Time
}

type OrderType uint

const (
    Buy OrderType = iota
    Sell
)

type Asset struct {

}

type Symbol struct {
    
}

type Order struct {
    OrderType OrderType
    Symbol    Symbol
    Price     float64
    Quantity  float64
}

type OrderId uint

type OrderStatus struct {
    Filled         bool
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
    // get data
    GetHistoricalSpreads(symbol []Symbol, seconds time.Duration) map[Symbol][]Spread
    GetCurrentSpread(symbol []Symbol) map[Symbol]Spread
    GetOrderBook(symbol []Symbol) map[Symbol]OrderBook

    GetLatency() time.Duration

    // deal with orders
    ExecuteOrders(orders []Order) []OrderId
    GetOrderStatuses(orderIds []OrderId) map[OrderId]OrderStatus
    CancelOrders(orderIds []OrderId)

    GetBalances() map[Asset]float64
}
