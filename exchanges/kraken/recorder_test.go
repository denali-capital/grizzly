package kraken

import (
	"fmt"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
)

func TestKrakenSpreadRecorder(t *testing.T) {
	krakenSpreadRecorder := NewKrakenSpreadRecorder(grizzlytesting.AssetPairs, grizzlytesting.Iso4217Translator, 10)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testRecorderGetHistoricalSpreads(t, krakenSpreadRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testSpreadRegisterAssetPair(t, krakenSpreadRecorder)
	})
}

func testRecorderGetHistoricalSpreads(t *testing.T, krakenSpreadRecorder *KrakenSpreadRecorder) {
	for _, assetPair := range grizzlytesting.AssetPairs {
		translatedPair := grizzlytesting.Iso4217Translator[assetPair]
		historicalSpreads, ok := krakenSpreadRecorder.GetHistoricalSpreads(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
	}
}

func testSpreadRegisterAssetPair(t *testing.T, krakenSpreadRecorder *KrakenSpreadRecorder) {
	translatedPair := grizzlytesting.Iso4217Translator[grizzlytesting.DOGEUSD]
	krakenSpreadRecorder.RegisterAssetPair(grizzlytesting.DOGEUSD)
	time.Sleep(grizzlytesting.SleepDuration)
	historicalSpreads, ok := krakenSpreadRecorder.GetHistoricalSpreads(grizzlytesting.DOGEUSD)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
}

func TestKrakenOrderBookRecorder(t *testing.T) {
	krakenOrderBookRecorder := NewKrakenOrderBookRecorder(grizzlytesting.AssetPairs, grizzlytesting.Iso4217Translator, 100)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetOrderBook", func(t *testing.T) {
		testGetOrderBook(t, krakenOrderBookRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testOrderBookRegisterAssetPair(t, krakenOrderBookRecorder)
	})
}

func testGetOrderBook(t *testing.T, krakenOrderBookRecorder *KrakenOrderBookRecorder) {
	for _, assetPair := range grizzlytesting.AssetPairs {
		translatedPair := grizzlytesting.Iso4217Translator[assetPair]
		orderBook, ok := krakenOrderBookRecorder.GetOrderBook(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, orderBook)
	}
}

func testOrderBookRegisterAssetPair(t *testing.T, krakenOrderBookRecorder *KrakenOrderBookRecorder) {
	translatedPair := grizzlytesting.Iso4217Translator[grizzlytesting.DOGEUSD]
	krakenOrderBookRecorder.RegisterAssetPair(grizzlytesting.DOGEUSD)
	time.Sleep(grizzlytesting.SleepDuration)
	orderBook, ok := krakenOrderBookRecorder.GetOrderBook(grizzlytesting.DOGEUSD)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, orderBook)
}
