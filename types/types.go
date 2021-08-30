package types

import "time"

type Spread struct {
    bid       uint
    ask       uint
    timestamp time.Time
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
    orderType OrderType
    symbol    Symbol
    price     float64
    quantity  float64
}

type OrderId uint

type OrderStatus struct {
    filled         bool
    filledPrice    *float64
    filledQuantity *float64
    original       *Order
}

type OrderBook struct {
    bids []OrderBookEntry
    asks []OrderBookEntry
}

type OrderBookEntry struct {
    price    float64
    quantity float64
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
