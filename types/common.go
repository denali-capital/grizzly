package types

import (
    "time"

    "github.com/shopspring/decimal"
)

type Spread struct {
    Bid       decimal.Decimal
    Ask       decimal.Decimal
    Timestamp time.Time
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
    Price     decimal.Decimal
    Quantity  decimal.Decimal
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
    FilledPrice    *decimal.Decimal
    FilledQuantity *decimal.Decimal
    Original       *Order
}

type OrderBook struct {
    Bids []OrderBookEntry
    Asks []OrderBookEntry
}

type OrderBookEntry struct {
    Price    decimal.Decimal
    Quantity decimal.Decimal
    UpdateId uint
}

type Observation struct {
    PriceDelta  float32
    Liquidity1  float32
    Liquidity2  float32
    Latency1    float32
    Latency2    float32
    Volatility1 float32
    Volatility2 float32

    // optional
    Label       int32
}
