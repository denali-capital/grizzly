package testing

import (
	"time"

	"github.com/denali-capital/grizzly/types"
)

const (
	BTCUSD types.AssetPair = iota
	ETHUSDT
	ADAUSDT
	BTCUSDC
	LTCUSDC
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

var KuCoinAssetPairs []types.AssetPair = []types.AssetPair{ETHUSDT, ADAUSDT, LTCUSDC}
var KuCoinAssetPairTranslator types.AssetPairTranslator = types.AssetPairTranslator{
	ETHUSDT: "ETH-USDT",
	ADAUSDT: "ADA-USDT",
	BTCUSDC: "BTC-USDC",
	LTCUSDC: "LTC-USDC",
}

const SleepDuration time.Duration = 3 * time.Second
const SampleDuration time.Duration = 2 * time.Second
const LatencyDuration time.Duration = time.Second

const Samples uint = 10

const SecretKeySuffix string = "_SECRET_KEY"
