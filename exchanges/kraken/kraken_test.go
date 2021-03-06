package kraken

import (
	"fmt"
	"os"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
	"github.com/joho/godotenv"
)

const apiKey string = "c8Mrlv+qder9EzFm+1trJRthtsYzgSBYHNP8opkB0O5FR+gS3UY52ex0"

// TODO: test use data stuff
func TestKraken(t *testing.T) {
	err := godotenv.Load("../../.env")
    if err != nil {
        t.Fatalf("Error loading .env file\n%v\n", err)
    }
	kraken := NewKraken(apiKey, os.Getenv("KRAKEN" + grizzlytesting.SecretKeySuffix), grizzlytesting.KrakenAssetPairTranslator, grizzlytesting.Iso4217Translator)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testKrakenGetHistoricalSpreads(t, kraken)
	})
	t.Run("GetCurrentSpread", func(t *testing.T) {
		testGetCurrentSpread(t, kraken)
	})
	t.Run("GetOrderBooks", func(t *testing.T) {
		testGetOrderBooks(t, kraken)
	})
	t.Run("GetLatency", func(t *testing.T) {
		testGetLatency(t, kraken)
	})
	t.Run("GetBalances", func(t *testing.T) {
		testGetBalances(t, kraken)
	})
}

func testKrakenGetHistoricalSpreads(t *testing.T, kraken *Kraken) {
	historicalSpreads := kraken.GetHistoricalSpreads(grizzlytesting.AssetPairs, grizzlytesting.SampleDuration, grizzlytesting.Samples)
	if len(historicalSpreads) == 0 {
		t.Fatalf("HistoricalSpreads should not be empty")
	}
	for assetPair, historicalSpread := range historicalSpreads {
		if uint(len(historicalSpread)) != grizzlytesting.Samples {
			t.Fatalf("There should be %v samples", grizzlytesting.Samples)
		}
		fmt.Printf("%v : %v\n", grizzlytesting.KrakenAssetPairTranslator[assetPair], historicalSpread)
	}
	fmt.Println(historicalSpreads)
}

func testGetCurrentSpread(t *testing.T, kraken *Kraken) {
	spread := kraken.GetCurrentSpread(grizzlytesting.BTCUSD)
	fmt.Println(spread)
}

func testGetOrderBooks(t *testing.T, kraken *Kraken) {
	orderBooks := kraken.GetOrderBooks(grizzlytesting.AssetPairs)
	if len(orderBooks) == 0 {
		t.Fatalf("OrderBooks should not be empty")
	}
	for assetPair, orderBook := range orderBooks {
		fmt.Printf("%v: %v\n", grizzlytesting.KrakenAssetPairTranslator[assetPair], *orderBook)
	}
}

func testGetLatency(t *testing.T, kraken *Kraken) {
	latency := kraken.GetLatency()
	fmt.Println(latency)
	time.Sleep(grizzlytesting.LatencyDuration)
	latency = kraken.GetLatency()
	fmt.Println(latency)
}

// func testExecuteOrders(t *testing.T, kraken *Kraken) {}

// func testGetOrderStatuses(t *testing.T, kraken *Kraken) {}

// func testCancelOrders(t *testing.T, kraken *Kraken) {}

func testGetBalances(t *testing.T, kraken *Kraken) {
	balances := kraken.GetBalances()
	fmt.Println(balances)
}
