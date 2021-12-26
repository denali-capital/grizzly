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

type ConcurrentOrderBook struct {
    sync.RWMutex
    internal types.OrderBook
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

func (c *ConcurrentOrderBook) Data() types.OrderBook {
    c.RLock()
    defer c.RUnlock()
    return c.internal
}
