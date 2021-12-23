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

func (f *FixedSizeSpreadQueue) Push(spread types.Spread) {
    if uint(len(f.internal)) == f.capacity {
        f.Pop()
    }
    f.internal = append(f.internal, spread)
}

func (f *FixedSizeSpreadQueue) Pop() {
    f.internal = f.internal[1:]
}

func (f *FixedSizeSpreadQueue) Data() []types.Spread {
    tmp := make([]types.Spread, len(f.internal))
    copy(tmp, f.internal)
    return tmp
}
