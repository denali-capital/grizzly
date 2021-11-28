package util

import (
    "fmt"
    "strings"

    "github.com/denali-capital/grizzly/types"
)

func Zip(slices ...[]int) ([][]int, error) {
    if len(slices) == 0 {
        return [][]int{}, nil
    }

    length := len(slices[0])
    for i := 1; i < len(slices); i++ {
        if len(slices[i]) != length {
            return nil, fmt.Errorf("zip: arguments must be of same length")
        }
    }

    r := make([][]int, length)

    for i, e := range slices[0] {
        a := make([]int, len(slices))
        a[0] = e
        for j := 1; j < len(slices); j++ {
            a[j] = slices[j][i]
        }
        r[i] = a
    }

    return r, nil
}

func StringIntersection(a []string, b []string) []string {
    set := make([]string, 0)
    hash := make(map[string]bool)

    for i := 0; i < a.Len(); i++ {
        hash[a[i]] = true
    }

    for i := 0; i < b.Len(); i++ {
        if _, found := hash[b[i]]; found {
            set = append(set, b[i])
        }
    }

    return set
}

func AssetPairIntersection(a []types.AssetPair, b []types.AssetPair) []types.AssetPair {
    set := make([]types.AssetPair, 0)
    hash := make(map[types.AssetPair]bool)

    for i := 0; i < a.Len(); i++ {
        hash[a[i]] = true
    }

    for i := 0; i < b.Len(); i++ {
        if _, found := hash[b[i]]; found {
            set = append(set, b[i])
        }
    }

    return set
}

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func addCombinations(c chan []*types.Exchange, exchanges []*types.Exchange, k uint, init []*types.Exchange) {
    if k == 0 {
        c <- init
        return
    }

    otherExchanges := make([]*types.Exchange, len(exchanges))
    copy(otherExchanges, exchanges)
    for _, exchange := range exchanges {
        otherExchanges = otherExchanges[1:]
        addCombinations(c, otherExchanges, k - 1, append(init, exchange))
    }
}

func ExchangeCombinations(exchanges []*types.Exchange, k uint) <-chan []*types.Exchange {
    c := make(chan []*types.Exchange)

    go func(c chan []*types.Exchange){
        defer close(c)

        addCombinations(c, exchanges, k, []*types.Exchanges{})
    }(c)

    return c
}
