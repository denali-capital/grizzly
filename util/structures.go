package util

import "github.com/denali-capital/grizzly/types"

type FixedSizeSpreadQueue struct {
    internal []types.Spread
    capacity uint
}

func NewFixedSizeSpreadQueue(capacity uint) *FixedSizeSpreadQueue {
    return &FixedSizeSpreadQueue{
        internal: make([]types.Spread, 0, capacity),
        capacity: capacity,
    }
}

func (f *FixedSizeSpreadQueue) Enqueue(item types.Spread) {
    if uint(len(f.internal)) == f.capacity {
        f.internal = f.internal[1:]
    }
    f.internal = append(f.internal, item)
}

func (f *FixedSizeSpreadQueue) Dequeue() {
    f.internal = f.internal[1:]
}

func (f *FixedSizeSpreadQueue) SetCapacity(capacity uint) {
    f.capacity = capacity
    if uint(len(f.internal)) > capacity {
        f.internal = f.internal[uint(len(f.internal)) - capacity:]
    }
}

func (f *FixedSizeSpreadQueue) Data() []types.Spread {
    tmp := make([]types.Spread, len(f.internal))
    copy(tmp, f.internal)
    return tmp
}
