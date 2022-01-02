package kucoin

import (
	"fmt"
	"os"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
	"github.com/joho/godotenv"
)

const apiKey string = "61cff0fbd5e581000154c4a5"

// TODO: test use data stuff
func TestKuCoin(t *testing.T) {
	err := godotenv.Load("../../.env")
    if err != nil {
        t.Fatalf("Error loading .env file\n%v\n", err)
    }
	kuCoin := NewKuCoin(apiKey, os.Getenv("KUCOIN" + grizzlytesting.SecretKeySuffix), os.Getenv("KUCOIN_API_PASSPHRASE"), grizzlytesting.KuCoinAssetPairTranslator)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testKuCoinGetHistoricalSpreads(t, kuCoin)
	})
	t.Run("GetCurrentSpread", func(t *testing.T) {
		testGetCurrentSpread(t, kuCoin)
	})
	// t.Run("GetOrderBooks", func(t *testing.T) {
	// 	testGetOrderBooks(t, kuCoin)
	// })
	t.Run("GetLatency", func(t *testing.T) {
		testGetLatency(t, kuCoin)
	})
	t.Run("GetBalances", func(t *testing.T) {
		testGetBalances(t, kuCoin)
	})
}

func testKuCoinGetHistoricalSpreads(t *testing.T, kuCoin *KuCoin) {
	historicalSpreads := kuCoin.GetHistoricalSpreads(grizzlytesting.KuCoinAssetPairs, grizzlytesting.SampleDuration, grizzlytesting.Samples)
	if len(historicalSpreads) == 0 {
		t.Fatalf("HistoricalSpreads should not be empty")
	}
	for assetPair, historicalSpread := range historicalSpreads {
		if uint(len(historicalSpread)) != grizzlytesting.Samples {
			t.Fatalf("There should be %v samples", grizzlytesting.Samples)
		}
		fmt.Printf("%v : %v\n", grizzlytesting.KuCoinAssetPairTranslator[assetPair], historicalSpread)
	}
	fmt.Println(historicalSpreads)
}

func testGetCurrentSpread(t *testing.T, kuCoin *KuCoin) {
	spread := kuCoin.GetCurrentSpread(grizzlytesting.ETHUSDT)
	fmt.Println(spread)
}

// func testGetOrderBooks(t *testing.T, kuCoin *KuCoin) {
// 	orderBooks := kuCoin.GetOrderBooks(grizzlytesting.AssetPairs)
// 	if len(orderBooks) == 0 {
// 		t.Fatalf("OrderBooks should not be empty")
// 	}
// 	for assetPair, orderBook := range orderBooks {
// 		fmt.Printf("%v: %v\n", grizzlytesting.KuCoinAssetPairTranslator[assetPair], *orderBook)
// 	}
// }

func testGetLatency(t *testing.T, kuCoin *KuCoin) {
	latency := kuCoin.GetLatency()
	fmt.Println(latency)
	time.Sleep(grizzlytesting.LatencyDuration)
	latency = kuCoin.GetLatency()
	fmt.Println(latency)
}

// func testExecuteOrders(t *testing.T, kuCoin *KuCoin) {}

// func testGetOrderStatuses(t *testing.T, kuCoin *KuCoin) {}

// func testCancelOrders(t *testing.T, kuCoin *KuCoin) {}

func testGetBalances(t *testing.T, kuCoin *KuCoin) {
	balances := kuCoin.GetBalances()
	fmt.Println(balances)
}
