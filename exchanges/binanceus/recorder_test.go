package binanceus

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
)

func TestBinanceUSSpreadRecorder(t *testing.T) {
	binanceUSSpreadRecorder := NewBinanceUSSpreadRecorder(grizzlytesting.AssetPairs, grizzlytesting.BinanceUSAssetPairTranslator, 10)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testRecorderGetHistoricalSpreads(t, binanceUSSpreadRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testSpreadRegisterAssetPair(t, binanceUSSpreadRecorder)
	})
}

func testRecorderGetHistoricalSpreads(t *testing.T, binanceUSSpreadRecorder *BinanceUSSpreadRecorder) {
	for _, assetPair := range grizzlytesting.AssetPairs {
		translatedPair := grizzlytesting.BinanceUSAssetPairTranslator[assetPair]
		historicalSpreads, ok := binanceUSSpreadRecorder.GetHistoricalSpreads(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
	}
}

func testSpreadRegisterAssetPair(t *testing.T, binanceUSSpreadRecorder *BinanceUSSpreadRecorder) {
	translatedPair := grizzlytesting.BinanceUSAssetPairTranslator[grizzlytesting.DOGEUSD]
	binanceUSSpreadRecorder.RegisterAssetPair(grizzlytesting.DOGEUSD)
	time.Sleep(grizzlytesting.SleepDuration)
	historicalSpreads, ok := binanceUSSpreadRecorder.GetHistoricalSpreads(grizzlytesting.DOGEUSD)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
}

func TestBinanceUSOrderBookRecorder(t *testing.T) {
	binanceUSOrderBookRecorder := NewBinanceUSOrderBookRecorder(&http.Client{}, grizzlytesting.AssetPairs, grizzlytesting.BinanceUSAssetPairTranslator, 100)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetOrderBook", func(t *testing.T) {
		testGetOrderBook(t, binanceUSOrderBookRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testOrderBookRegisterAssetPair(t, binanceUSOrderBookRecorder)
	})
}

func testGetOrderBook(t *testing.T, binanceUSOrderBookRecorder *BinanceUSOrderBookRecorder) {
	for _, assetPair := range grizzlytesting.AssetPairs {
		translatedPair := grizzlytesting.BinanceUSAssetPairTranslator[assetPair]
		orderBook, ok := binanceUSOrderBookRecorder.GetOrderBook(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, orderBook)
	}
}

func testOrderBookRegisterAssetPair(t *testing.T, binanceUSOrderBookRecorder *BinanceUSOrderBookRecorder) {
	translatedPair := grizzlytesting.BinanceUSAssetPairTranslator[grizzlytesting.DOGEUSD]
	binanceUSOrderBookRecorder.RegisterAssetPair(grizzlytesting.DOGEUSD)
	time.Sleep(grizzlytesting.SleepDuration)
	orderBook, ok := binanceUSOrderBookRecorder.GetOrderBook(grizzlytesting.DOGEUSD)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, orderBook)
}
