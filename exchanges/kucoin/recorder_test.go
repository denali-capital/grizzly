package kucoin

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	grizzlytesting "github.com/denali-capital/grizzly/testing"
	"github.com/joho/godotenv"
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

func TestKuCoinOrderBookRecorder(t *testing.T) {
	err := godotenv.Load("../../.env")
    if err != nil {
        t.Fatalf("Error loading .env file\n%v\n", err)
    }
	kuCoinOrderBookRecorder := NewKuCoinOrderBookRecorder(&http.Client{}, apiKey, os.Getenv("KUCOIN" + grizzlytesting.SecretKeySuffix), os.Getenv("KUCOIN_API_PASSPHRASE"), grizzlytesting.KuCoinAssetPairs, grizzlytesting.KuCoinAssetPairTranslator, 100)
	time.Sleep(grizzlytesting.SleepDuration)
	t.Run("GetOrderBook", func(t *testing.T) {
		testGetOrderBook(t, kuCoinOrderBookRecorder)
	})
	t.Run("RegisterAssetPair", func(t *testing.T) {
		testOrderBookRegisterAssetPair(t, kuCoinOrderBookRecorder)
	})
}

func testGetOrderBook(t *testing.T, kuCoinOrderBookRecorder *KuCoinOrderBookRecorder) {
	for _, assetPair := range grizzlytesting.KuCoinAssetPairs {
		translatedPair := grizzlytesting.KuCoinAssetPairTranslator[assetPair]
		orderBook, ok := kuCoinOrderBookRecorder.GetOrderBook(assetPair)
		if !ok {
			t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
		}
		fmt.Printf("%v: %v\n", translatedPair, orderBook)
	}
}

func testOrderBookRegisterAssetPair(t *testing.T, kuCoinOrderBookRecorder *KuCoinOrderBookRecorder) {
	translatedPair := grizzlytesting.KuCoinAssetPairTranslator[grizzlytesting.BTCUSDC]
	kuCoinOrderBookRecorder.RegisterAssetPair(grizzlytesting.BTCUSDC)
	time.Sleep(grizzlytesting.SleepDuration)
	orderBook, ok := kuCoinOrderBookRecorder.GetOrderBook(grizzlytesting.BTCUSDC)
	if !ok {
		t.Fatalf("AssetPair %v should be recorded\n", translatedPair)
	}
	fmt.Printf("%v: %v\n", translatedPair, orderBook)
}
