package util

import (
    "log"
    "math"

    "github.com/montanaflynn/stats"
    "github.com/denali-capital/grizzly/types"
    "github.com/shopspring/decimal"
)

// TODO: can do better than midpoint prices based on the paper
var two decimal.Decimal = decimal.NewFromInt(2)

func ComputePriceVolatility(spreads []types.Spread) float64 {
    midpoints := make([]float64, len(spreads))

    for i, spread := range spreads {
        midpoints[i] = spread.Bid.Add(spread.Ask).Div(two).InexactFloat64()
    }

    sdev, err := stats.StdDevS(midpoints)

    if err != nil {
        // this only fails if the len of slice is zero
        log.Fatalln(err)
    }

    return sdev
}

func ComputeSlippage(orderbook *types.OrderBook, quantity decimal.Decimal) decimal.Decimal {
    midpoint := orderbook.Bids[0].Price.Add(orderbook.Asks[0].Price).Div(two)
    idealCost := quantity.Mul(midpoint)

    // compute slippage on buy side
    buyCost := computeOrderCost(orderbook.Asks, quantity)

    // compute slippage on sell side
    sellCost := computeOrderCost(orderbook.Bids, quantity)

    return buyCost.Sub(sellCost).Div(idealCost.Mul(two))
}

func computeOrderCost(side []types.OrderBookEntry, amountToFill decimal.Decimal) decimal.Decimal {
    cumulativeCost := decimal.Zero

    for _, entry := range side {
        if amountToFill.LessThan(entry.Quantity) {
            cumulativeCost = cumulativeCost.Add(amountToFill.Mul(entry.Price))
            // the order is now entirely filled
            break
        }

        cumulativeCost = cumulativeCost.Add(entry.Quantity.Mul(entry.Price))
        amountToFill = amountToFill.Sub(entry.Quantity)
    }

    return cumulativeCost
}

// use library?
// based off RFC 2988 for estimating RTT
type EwmaEstimator struct {
    estimate  float64
    variation float64
    first     bool
    alpha     float64
    beta      float64
    k         float64
}

func NewEwmaEstimator(alpha, beta, k float64) *EwmaEstimator {
    return &EwmaEstimator{
        first: true,
        alpha: alpha,
        beta: beta,
        k: k,
    }
}

func (e *EwmaEstimator) Sample(sample float64) {
    if e.first {
        e.estimate = sample
        e.first = false
    } else {
        e.variation = (1 - e.beta) * e.variation + e.beta * math.Abs(e.estimate - sample)
        e.estimate = (1 - e.alpha) * e.estimate + e.alpha * sample
    }
}

func (e *EwmaEstimator) GetEstimate() float64 {
    return e.estimate + e.k * e.variation
}
