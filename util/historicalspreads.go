package util

import (
	"log"
	"time"

	"github.com/denali-capital/grizzly/types"
)

func GetSpreadSamples(historicalSpreads []types.Spread, duration time.Duration, samples uint) []types.Spread {
	mostRecentTimestamp := historicalSpreads[len(historicalSpreads) - 1].Timestamp

    // make timestamp sample list
    // nanoseconds per sample
    period := time.Duration(duration.Nanoseconds() / int64(samples))

    // check enough samples exist
    if leastRecentTimestamp := historicalSpreads[0].Timestamp; mostRecentTimestamp.Add(time.Duration(-(samples - 1)) * period).Before(leastRecentTimestamp) {
        log.Println("warning: duration is too long, using longest possible duration instead")
        period = time.Duration(mostRecentTimestamp.Sub(leastRecentTimestamp).Nanoseconds() / int64(samples))
    }

    timestamps := make([]time.Time, samples)
    for i := uint(0); i < samples; i++ {
        timestamps[i] = mostRecentTimestamp.Add(time.Duration(-(samples - i - 1)) * period)
    }

    // get spreads according to sample list
    res := make([]types.Spread, samples)
    currentTimestampIndex := 0
    for i, spread := range historicalSpreads {
        for (timestamps[currentTimestampIndex].After(spread.Timestamp) || timestamps[currentTimestampIndex].Equal(spread.Timestamp)) && (i + 1 >= len(historicalSpreads) || timestamps[currentTimestampIndex].Before(historicalSpreads[i + 1].Timestamp)) {
            effectiveSpread := spread
            effectiveSpread.Timestamp = timestamps[currentTimestampIndex]
            res[currentTimestampIndex] = effectiveSpread
            currentTimestampIndex++
            if currentTimestampIndex == len(timestamps) {
                break
            }
        }
        if currentTimestampIndex == len(timestamps) {
            break
        }
    }

    return res
}
