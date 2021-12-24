package testing

import (
	"time"

	"github.com/denali-capital/grizzly/types"
)

const (
	BTCUSD types.AssetPair = iota
	ADAUSDT
	BTCUSDC
	DOGEUSD
)

var AssetPairs []types.AssetPair = []types.AssetPair{BTCUSD, ADAUSDT, BTCUSDC}
var Iso4217Translator types.AssetPairTranslator = types.AssetPairTranslator{
	BTCUSD: "XBT/USD",
	ADAUSDT: "ADA/USDT",
	BTCUSDC: "XBT/USDC",
	DOGEUSD: "XDG/USD",
}
var BinanceUSAssetPairTranslator types.AssetPairTranslator = types.AssetPairTranslator{
	BTCUSD: "BTCUSD",
	ADAUSDT: "ADAUSDT",
	BTCUSDC: "BTCUSDC",
	DOGEUSD: "DOGEUSD",
}
var KrakenAssetPairTranslator types.AssetPairTranslator = types.AssetPairTranslator{
	BTCUSD: "XXBTZUSD",
	ADAUSDT: "ADAUSDT",
	BTCUSDC: "XBTUSDC",
	DOGEUSD: "XDGUSD",
}

const SleepDuration time.Duration = 3 * time.Second
const SampleDuration time.Duration = 2 * time.Second
const LatencyDuration time.Duration = time.Second

const Samples uint = 10
