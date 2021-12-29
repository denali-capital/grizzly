package binanceus

import (
	"fmt"
	"os"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
	"github.com/joho/godotenv"
)

const apiKey string = "JJqJgQeJhZCgXXZp6SBea3Dje8mmsWW0hWEjuVkYLOuZjxVPf71tNIE5AJU0GsYt"

// TODO: test use data stuff
func TestBinanceUS(t *testing.T) {
	err := godotenv.Load("../../.env")
    if err != nil {
        t.Fatalf("Error loading .env file\n%v\n", err)
    }
	binanceUS := NewBinanceUS(apiKey, os.Getenv("BINANCEUS" + grizzlytesting.SecretKeySuffix), grizzlytesting.BinanceUSAssetPairTranslator)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testBinanceUSGetHistoricalSpreads(t, binanceUS)
	})
	t.Run("GetCurrentSpread", func(t *testing.T) {
		testGetCurrentSpread(t, binanceUS)
	})
	t.Run("GetOrderBooks", func(t *testing.T) {
		testGetOrderBooks(t, binanceUS)
	})
	t.Run("GetLatency", func(t *testing.T) {
		testGetLatency(t, binanceUS)
	})
	t.Run("GetBalances", func(t *testing.T) {
		testGetBalances(t, binanceUS)
	})
}

func testBinanceUSGetHistoricalSpreads(t *testing.T, binanceUS *BinanceUS) {
	historicalSpreads := binanceUS.GetHistoricalSpreads(grizzlytesting.AssetPairs, grizzlytesting.SampleDuration, grizzlytesting.Samples)
	if len(historicalSpreads) == 0 {
		t.Fatalf("HistoricalSpreads should not be empty")
	}
	for assetPair, historicalSpread := range historicalSpreads {
		if uint(len(historicalSpread)) != grizzlytesting.Samples {
			t.Fatalf("There should be %v samples", grizzlytesting.Samples)
		}
		fmt.Printf("%v : %v\n", grizzlytesting.BinanceUSAssetPairTranslator[assetPair], historicalSpread)
	}
	fmt.Println(historicalSpreads)
}

func testGetCurrentSpread(t *testing.T, binanceUS *BinanceUS) {
	spread := binanceUS.GetCurrentSpread(grizzlytesting.BTCUSD)
	fmt.Println(spread)
}

func testGetOrderBooks(t *testing.T, binanceUS *BinanceUS) {
	orderBooks := binanceUS.GetOrderBooks(grizzlytesting.AssetPairs)
	if len(orderBooks) == 0 {
		t.Fatalf("OrderBooks should not be empty")
	}
	for assetPair, orderBook := range orderBooks {
		fmt.Printf("%v: %v\n", grizzlytesting.BinanceUSAssetPairTranslator[assetPair], *orderBook)
	}
}

func testGetLatency(t *testing.T, binanceUS *BinanceUS) {
	latency := binanceUS.GetLatency()
	fmt.Println(latency)
	time.Sleep(grizzlytesting.LatencyDuration)
	latency = binanceUS.GetLatency()
	fmt.Println(latency)
}

// func testExecuteOrders(t *testing.T, binanceUS *BinanceUS) {}

// func testGetOrderStatuses(t *testing.T, binanceUS *BinanceUS) {}

// func testCancelOrders(t *testing.T, binanceUS *BinanceUS) {}

func testGetBalances(t *testing.T, binanceUS *BinanceUS) {
	balances := binanceUS.GetBalances()
	fmt.Println(balances)
}
