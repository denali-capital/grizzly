package kraken

import (
	"log"
	"sort"

	"github.com/shopspring/decimal"
)

func remove(slice []DecimalOrderBookEntry, s int) []DecimalOrderBookEntry {
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

func RemovePriceFromBids(bids []DecimalOrderBookEntry, price decimal.Decimal) []DecimalOrderBookEntry {
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price.LessThanOrEqual(price)
	})
	if i < len(bids) && bids[i].Price.Equals(price) {
		return remove(bids, i)
	} else {
		return bids
	}
}

func InsertPriceInBids(bids []DecimalOrderBookEntry, price decimal.Decimal, quantity decimal.Decimal) []DecimalOrderBookEntry {
	bids = RemovePriceFromBids(bids, price)
	decimalOrderBookEntry := DecimalOrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(bids), func(i int) bool {
		return bids[i].Price.LessThan(price)
	})
	bids = append(bids, DecimalOrderBookEntry{})
	copy(bids[i + 1:], bids[i:])
	bids[i] = decimalOrderBookEntry
	return bids
}

func RemovePriceFromAsks(asks []DecimalOrderBookEntry, price decimal.Decimal) []DecimalOrderBookEntry {
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price.GreaterThanOrEqual(price)
	})
	if i < len(asks) && asks[i].Price.Equals(price) {
		return remove(asks, i)
	} else {
		return asks
	}
}

func InsertPriceInAsks(asks []DecimalOrderBookEntry, price decimal.Decimal, quantity decimal.Decimal) []DecimalOrderBookEntry {
	asks = RemovePriceFromAsks(asks, price)
	decimalOrderBookEntry := DecimalOrderBookEntry{
		Price: price,
		Quantity: quantity,
	}
	i := sort.Search(len(asks), func(i int) bool {
		return asks[i].Price.GreaterThan(price)
	})
	asks = append(asks, DecimalOrderBookEntry{})
	copy(asks[i + 1:], asks[i:])
	asks[i] = decimalOrderBookEntry
	return asks
}
