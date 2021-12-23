package testing

import (
	"time"

	"github.com/denali-capital/grizzly/types"
)

const (
	BTCUSD types.AssetPair = iota
	ADAUSDT
	ETHUSDC
	DOGEUSD
)

var AssetPairs []types.AssetPair = []types.AssetPair{BTCUSD, ADAUSDT, ETHUSDC}
var Iso4217Translator types.AssetPairTranslator = types.AssetPairTranslator{
	BTCUSD: "XBT/USD",
	ADAUSDT: "ADA/USDT",
	ETHUSDC: "ETH/USDC",
	DOGEUSD: "XDG/USD",
}

const SleepDuration time.Duration = 3 * time.Second
const SampleDuration time.Duration = 2 * time.Second
const LatencyDuration time.Duration = time.Second

const Samples uint = 10
