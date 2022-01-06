package util

import (
    "sync"

    "github.com/denali-capital/grizzly/types"
)

type ConcurrentOrderIdToOrderPtrMap struct {
    sync.RWMutex
    internal map[types.OrderId]*types.Order
}

func NewConcurrentOrderIdToOrderPtrMap() *ConcurrentOrderIdToOrderPtrMap {
    return &ConcurrentOrderIdToOrderPtrMap{
        internal: make(map[types.OrderId]*types.Order),
    }
}

func (c *ConcurrentOrderIdToOrderPtrMap) Load(key types.OrderId) (value *types.Order, ok bool) {
    c.RLock()
    defer c.RUnlock()
    result, ok := c.internal[key]
    return result, ok
}

func (c *ConcurrentOrderIdToOrderPtrMap) Delete(key types.OrderId) {
    c.Lock()
    defer c.Unlock()
    delete(c.internal, key)
}

func (c *ConcurrentOrderIdToOrderPtrMap) Store(key types.OrderId, value *types.Order) {
    c.Lock()
    defer c.Unlock()
    c.internal[key] = value
}

type ConcurrentFixedSizeSpreadQueue struct {
    sync.RWMutex
    internal []types.Spread
    capacity uint
}

func NewConcurrentFixedSizeSpreadQueue(capacity uint) *ConcurrentFixedSizeSpreadQueue {
    return &ConcurrentFixedSizeSpreadQueue{
        internal: make([]types.Spread, 0, capacity),
        capacity: capacity,
    }
}

func (c *ConcurrentFixedSizeSpreadQueue) Push(spread types.Spread) {
    c.Lock()
    defer c.Unlock()
    if uint(len(c.internal)) == c.capacity {
        c.internal = c.internal[1:]
    }
    c.internal = append(c.internal, spread)
}

func (c *ConcurrentFixedSizeSpreadQueue) Pop() {
    c.Lock()
    defer c.Unlock()
    c.internal = c.internal[1:]
}

func (c *ConcurrentFixedSizeSpreadQueue) Data() []types.Spread {
    c.RLock()
    defer c.RUnlock()
    tmp := make([]types.Spread, len(c.internal))
    copy(tmp, c.internal)
    return tmp
}

func (c *ConcurrentFixedSizeSpreadQueue) Back() types.Spread {
    c.RLock()
    defer c.RUnlock()
    return c.internal[len(c.internal) - 1]
}

type ConcurrentOrderBook struct {
    sync.RWMutex
    internal     types.OrderBook
    LastUpdateId uint
}

func NewConcurrentOrderBook(bids []types.OrderBookEntry, asks []types.OrderBookEntry) *ConcurrentOrderBook {
    return &ConcurrentOrderBook{
        internal: types.OrderBook{
            Bids: bids,
            Asks: asks,
        },
    }
}

func (c *ConcurrentOrderBook) GetAsks() []types.OrderBookEntry {
    c.RLock()
    defer c.RUnlock()
    return c.internal.Asks
}

func (c *ConcurrentOrderBook) GetBids() []types.OrderBookEntry {
    c.RLock()
    defer c.RUnlock()
    return c.internal.Bids
}

func (c *ConcurrentOrderBook) SetBidsAndAsks(bids []types.OrderBookEntry, asks []types.OrderBookEntry) {
    c.Lock()
    defer c.Unlock()
    c.internal.Bids = bids
    c.internal.Asks = asks
}

func (c *ConcurrentOrderBook) FilterAndMerge(other *ConcurrentOrderBook, prefer bool) {
    c.Lock()
    defer c.Unlock()
    lastUpdateId := other.LastUpdateId

    var bids []types.OrderBookEntry
    var asks []types.OrderBookEntry
    if prefer {
        bids = other.GetBids()
        bids = append(bids, c.internal.Bids...)
        asks = other.GetAsks()
        asks = append(asks, c.internal.Asks...)
    } else {
        bids = c.internal.Bids
        bids = append(bids, other.GetBids()...)
        asks = c.internal.Asks
        asks = append(asks, other.GetAsks()...)
    }
    processed := make(map[string]struct{})
    w := 0
    for _, orderBookEntry := range bids {
        if orderBookEntry.UpdateId > lastUpdateId {
            priceIdentifier := orderBookEntry.Price.String()
            if _, exists := processed[priceIdentifier]; !exists {
                processed[priceIdentifier] = struct{}{}
                bids[w] = orderBookEntry
                w++
            }
        }
    }
    bids = bids[:w]
    processed = make(map[string]struct{})
    w = 0
    for _, orderBookEntry := range asks {
        if orderBookEntry.UpdateId > lastUpdateId {
            priceIdentifier := orderBookEntry.Price.String()
            if _, exists := processed[priceIdentifier]; !exists {
                processed[priceIdentifier] = struct{}{}
                asks[w] = orderBookEntry
                w++
            }
        }
    }
    asks = asks[:w]
    c.internal.Bids = bids
    c.internal.Asks = asks
}

func (c *ConcurrentOrderBook) Data() types.OrderBook {
    c.RLock()
    defer c.RUnlock()
    return c.internal
}

type ConcurrentOrderBookResponse struct {
    AssetPair           types.AssetPair
    ConcurrentOrderBook *ConcurrentOrderBook
}
