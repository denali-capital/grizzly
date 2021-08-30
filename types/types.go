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
    filled   bool
    price    *float64
    quantity *float64
}

type Exchange interface {
    // functionality we want each exchange to implement

    // * getting data
    GetHistoricalSpreads(symbol []Symbol, seconds time.Duration) [][]Spread
    GetCurrentSpread(symbol []Symbol) []Spread
    GetLatency() time.Duration
    // get data for liquidity measurements
    // look into liquidity calculations

    // * executing and checking status of orders
    ExecuteOrders(orders []Order) []OrderId
    GetOrderStatuses(orderIds []OrderId) map[OrderId]OrderStatus
    CancelOrders(orderIds []OrderId)

    GetBalances() map[Asset]float64
}
