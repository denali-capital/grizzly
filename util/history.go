package util

import (
    "sync"
    "time"

    "github.com/denali-capital/grizzly/types"
)

type SpreadRecorder struct {
    sync.RWMutex
    spreadFunction         func(types.AssetPair) types.Spread
    capacity               uint
    period                 time.Duration
    spreadQueues           map[types.AssetPair]*FixedSizeSpreadQueue
}

func NewSpreadRecorder(assetPairs []types.AssetPair, spreadFunction func(types.AssetPair) types.Spread, capacity uint, period time.Duration) *SpreadRecorder {
    spreadQueues := make(map[types.AssetPair]*FixedSizeSpreadQueue)
    for _, assetPair := range assetPairs {
        spreadQueues[assetPair] = NewFixedSizeSpreadQueue(capacity)
    }
    spreadRecorder := &SpreadRecorder{
        spreadFunction: spreadFunction,
        capacity: capacity,
        period: period,
        spreadQueues: spreadQueues,
    }
    go spreadRecorder.record()
    return spreadRecorder
}

func (s *SpreadRecorder) record() {
    for {
        s.Lock()
        for assetPair, queue := range s.spreadQueues {
            queue.Enqueue(s.spreadFunction(assetPair))
        }
        s.Unlock()
        time.Sleep(s.period)
    }
}

func (s *SpreadRecorder) GetHistoricalSpreads(assetPair types.AssetPair) ([]types.Spread, bool) {
    s.RLock()
    defer s.RUnlock()
    result, ok := s.spreadQueues[assetPair]
    if !ok {
        return make([]types.Spread, 0), false
    }
    return result.Data(), true
}

func (s *SpreadRecorder) RegisterAssetPair(assetPair types.AssetPair) {
    s.Lock()
    defer s.Unlock()
    if _, ok := s.spreadQueues[assetPair]; !ok {
        s.spreadQueues[assetPair] = NewFixedSizeSpreadQueue(s.capacity)
    }
}

func (s *SpreadRecorder) SetCapacity(capacity uint) {
    s.Lock()
    defer s.Unlock()
    s.capacity = capacity
    for _, queue := range s.spreadQueues {
        queue.SetCapacity(capacity)
    }
}

func (s *SpreadRecorder) SetPeriod(period time.Duration) {
    s.Lock()
    defer s.Unlock()
    s.period = period
}
