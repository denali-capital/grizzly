package util

import (
	"log"
	"sort"

	"github.com/denali-capital/grizzly/types"
	"github.com/shopspring/decimal"
)

func remove(slice []types.OrderBookEntry, s int) []types.OrderBookEntry {
	return append(slice[:s], slice[s+1:]...)
	// preserves the order
}

func GetPriceAndQuantity(rawOrderBookEntry []interface{}) (decimal.Decimal, decimal.Decimal) {
	price, err := decimal.NewFromString(rawOrderBookEntry[0].(string))
	if err != nil {
		log.Fatal(err)
	}
	quantity, err := decimal.NewFromString(rawOrderBookEntry[1].(string))
	if err != nil {
		log.Fatal(err)
	}
	return price, quantity
}

func RemovePriceFromBids(bids []types.OrderBookEntry, price decimal.Decimal) []types.OrderBookEntry {
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price.LessThanOrEqual(price)
	})
	if i < len(bids) && bids[i].Price.Equals(price) {
		return remove(bids, i)
	} else {
		return bids
	}
}

func InsertPriceInBids(bids []types.OrderBookEntry, price decimal.Decimal, quantity decimal.Decimal) []types.OrderBookEntry {
	bids = RemovePriceFromBids(bids, price)
	orderBookEntry := types.OrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price.LessThan(price)
	})
	bids = append(bids, types.OrderBookEntry{})
	copy(bids[i + 1:], bids[i:])
	bids[i] = orderBookEntry
	return bids
}

func RemovePriceFromAsks(asks []types.OrderBookEntry, price decimal.Decimal) []types.OrderBookEntry {
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price.GreaterThanOrEqual(price)
	})
	if i < len(asks) && asks[i].Price.Equals(price) {
		return remove(asks, i)
	} else {
		return asks
	}
}

func InsertPriceInAsks(asks []types.OrderBookEntry, price decimal.Decimal, quantity decimal.Decimal) []types.OrderBookEntry {
	asks = RemovePriceFromAsks(asks, price)
	orderBookEntry := types.OrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price.GreaterThan(price)
	})
	asks = append(asks, types.OrderBookEntry{})
	copy(asks[i + 1:], asks[i:])
	asks[i] = orderBookEntry
	return asks
}
