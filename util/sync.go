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
