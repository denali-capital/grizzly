package kucoin

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
)

func TestKuCoinSpreadRecorder(t *testing.T) {
	kuCoinSpreadRecorder := NewKuCoinSpreadRecorder(&http.Client{}, grizzlytesting.KuCoinAssetPairs, grizzlytesting.KuCoinAssetPairTranslator, 10)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetHistoricalSpreads", func(t *testing.T) {
		testRecorderGetHistoricalSpreads(t, kuCoinSpreadRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testSpreadRegisterAssetPair(t, kuCoinSpreadRecorder)
	})
}

func testRecorderGetHistoricalSpreads(t *testing.T, kuCoinSpreadRecorder *KuCoinSpreadRecorder) {
	for _, assetPair := range grizzlytesting.KuCoinAssetPairs {
		translatedPair := grizzlytesting.KuCoinAssetPairTranslator[assetPair]
		historicalSpreads, ok := kuCoinSpreadRecorder.GetHistoricalSpreads(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
	}
}

func testSpreadRegisterAssetPair(t *testing.T, kuCoinSpreadRecorder *KuCoinSpreadRecorder) {
	translatedPair := grizzlytesting.KuCoinAssetPairTranslator[grizzlytesting.BTCUSDC]
	kuCoinSpreadRecorder.RegisterAssetPair(grizzlytesting.BTCUSDC)
	time.Sleep(grizzlytesting.SleepDuration)
	historicalSpreads, ok := kuCoinSpreadRecorder.GetHistoricalSpreads(grizzlytesting.BTCUSDC)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, historicalSpreads)
}
