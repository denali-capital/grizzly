package util

import (
    "log"
    "math"

    "github.com/montanaflynn/stats"
    "github.com/denali-capital/grizzly/types"
)

// TODO: can do better than midpoint prices based on the paper

func ComputePriceVolatility(spreads []types.Spread) float64 {
    midpoints := make([]float64, len(spreads))

    for i, spread := range spreads {
        midpoints[i] = (spread.Bid + spread.Ask) / 2
    }

    sdev, err := stats.StdDevS(midpoints)

    if err != nil {
        // this only fails if the len of slice is zero
        log.Fatalln(err)
    }

    return sdev
}

func ComputeSlippage(orderbook *types.OrderBook, quantity float64) float64 {
    midpoint := (orderbook.Bids[0].Price + orderbook.Asks[0].Price) / 2
    idealCost := quantity * midpoint

    // compute slippage on buy side
    buyCost := computeOrderCost(orderbook.Asks, quantity)

    // compute slippage on sell side
    sellCost := computeOrderCost(orderbook.Bids, quantity)

    return (buyCost - sellCost) / (2 * idealCost)
}

func computeOrderCost(side []types.OrderBookEntry, amountToFill float64) float64 {
    cumulativeCost := float64(0)

    for _, entry := range side {
        if amountToFill < entry.Quantity {
            cumulativeCost += amountToFill * entry.Price
            // the order is now entirely filled
            break
        }

        cumulativeCost += entry.Quantity * entry.Price
        amountToFill -= entry.Quantity
    }

    return cumulativeCost
}

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
