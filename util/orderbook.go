package util

import (
	"log"
	"sort"
	"strconv"

	"github.com/denali-capital/grizzly/types"
)

// general orderbook stuff
func remove(slice []types.OrderBookEntry, s int) []types.OrderBookEntry {
	// preserves the order
	return append(slice[:s], slice[s + 1:]...)
}

func GetPriceAndQuantity(rawOrderBookEntry []interface{}) (float64, float64) {
    price, err := strconv.ParseFloat(rawOrderBookEntry[0].(string), 64)
    if err != nil {
        log.Fatalln(err)
    }
    quantity, err := strconv.ParseFloat(rawOrderBookEntry[1].(string), 64)
    if err != nil {
        log.Fatalln(err)
    }
    return price, quantity
}

func RemovePriceFromBids(bids []types.OrderBookEntry, price float64) []types.OrderBookEntry {
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price <= price
	})
	if i < len(bids) && bids[i].Price == price {
		return remove(bids, i)
	} else {
		return bids
	}
}

func InsertPriceInBids(bids []types.OrderBookEntry, price float64, quantity float64) []types.OrderBookEntry {
	bids = RemovePriceFromBids(bids, price)
	orderBookEntry := types.OrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price < price
	})
	bids = append(bids, types.OrderBookEntry{})
	copy(bids[i + 1:], bids[i:])
	bids[i] = orderBookEntry
	return bids
}

func RemovePriceFromAsks(asks []types.OrderBookEntry, price float64) []types.OrderBookEntry {
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price >= price
	})
	if i < len(asks) && asks[i].Price == price {
		return remove(asks, i)
	} else {
		return asks
	}
}

func InsertPriceInAsks(asks []types.OrderBookEntry, price float64, quantity float64) []types.OrderBookEntry {
	asks = RemovePriceFromAsks(asks, price)
	orderBookEntry := types.OrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price > price
	})
	asks = append(asks, types.OrderBookEntry{})
	copy(asks[i + 1:], asks[i:])
	asks[i] = orderBookEntry
	return asks
}
