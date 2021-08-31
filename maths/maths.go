package maths

import (
    "github.com/montanaflynn/stats"
    "github.com/denali-capital/grizzly/types"
)

// TODO: can do better than midpoint prices based on the paper

func ComputePriceVolitility(spreads []types.Spread) float64 {
    midpoints := make([]float64, len(spreads))

    for i, spread := range spreads {
        midpoints[i] = (spread.Bid + spread.Ask) / 2
    }

    sdev, err := stats.StdDevS(midpoints)

    if err != nil {
        // this only fails if the len of slice is zero
        panic(err)
    }

    return sdev
}

func ComputeSlippage(orderbook types.OrderBook, quantity float64) float64 {
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