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

type AssetPairTranslator map[AssetPair]string

func (a AssetPairTranslator) GetAssetPairs() []AssetPair {
    keys := make([]AssetPair, len(a))

    i := 0
    for k := range a {
        keys[i] = k
        i++
    }

    return keys
}

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

type Observation struct {
    PriceDelta  float64
    Liquidity1  float64
    Liquidity2  float64
    Latency1    float64
    Latency2    float64
    Volatility1 float64
    Volatility2 float64

    // optional
    Label       int32
}
