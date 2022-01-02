package main

import (
    "log"
    "os"
    "strings"
    "time"

    "github.com/denali-capital/grizzly/exchanges/binanceus"
    "github.com/denali-capital/grizzly/exchanges/kraken"
    "github.com/denali-capital/grizzly/exchanges/kucoin"

    "github.com/denali-capital/grizzly/model/bootstrap"
    "github.com/denali-capital/grizzly/model/nn"
    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    _ "github.com/joho/godotenv/autoload"
)

// each exchange will have their own module that implements Exchange interface above
// can have module to compute statistics in background and access them
// Neural network module that handles predict / fit functionality

const configPath string = "config"
const secretKeySuffix string = "_SECRET_KEY"
const sleepDuration time.Duration = 100 * time.Millisecond

func grizzly(exchange1 *types.Exchange, exchange2 *types.Exchange, allowedAssetPairs []types.AssetPair, killerInstinct *nn.KillerInstinct) {
    for {
        if killerInstinct.Predict()

        time.Sleep(sleepDuration)
    }
}

func main() {
    // need to create array of exchange objects
    // pass this array into function that computes statistics in background

    implementedExchanges := util.DiscoverTypes("github.com/denali-capital/grizzly/exchanges")

    exchangeList := util.ReadCsvFile(configPath + "/exchanges.csv")
    if exchangeList[0][0] != "exchange" || exchangeList[0][1] != "api_key" || exchangeList[0][2] != "fees" {
        log.Fatalln("Labels must be \"exchange,api_key,fees,...\"")
    }

    zippedExchangeList := util.Zip(exchangeList...)
    allowedExchanges := util.StringIntersection(implementedExchanges, zippedExchangeList[0][1:])

    if len(allowedExchanges) == 0 {
        log.Fatalln("No exchanges are both implemented and have requisite config info")
    }

    exchangeInfo := make(map[string][]string)
    for i := 1; i < len(exchangeList); i++ {
        info := exchangeList[i]
        exchangeName := info[0]
        if util.Contains(allowedExchanges, exchangeName) {
            if info[1] == "" {
                log.Fatalln("API key not provided for %v", exchangeName)
            }
            if info[2] == "" {
                log.Println("warning: Fees not provided for %v, assuming 0", exchangeName)
            }
            exchangeInfo[exchangeName] := info
        }
    }

    assetPairsList := util.ReadCsvFile(configPath + "/assetpairs.csv")
    if assetPairsList[0][0] != "canonical" || assetPairsList[0][1] != "ISO4217" {
        log.Fatalln("Labels must be \"canonical,ISO4217,...\"")
    }

    exchangeIndices := make(map[string]uint)
    assetPairTranslators := make(map[string]types.AssetPairTranslator)
    for i, exchangeName := range assetPairsList[0][1:] {
        exchangeIndices[exchangeName] = i + 1
        assetPairTranslators[exchangeName] := make(types.AssetPairTranslator)
    }

    assetPairCanonicalTranslator := make(types.AssetPairTranslator)
    zippedAssetPairsList := util.Zip(assetPairsList...)
    for i, assetPairCanonical := range zippedAssetPairsList[0][1:] {
        assetPairCanonicalTranslator[i + 1] = assetPairCanonical
        for exchangeName, translator := range assetPairTranslators {
            assetPairSpecific := zippedAssetPairsList[exchangeIndices[exchangeName]][i + 1]
            if assetPairSpecific != "" {
                translator[i + 1] = assetPairSpecific
            }
        }
    }

    exchanges := make([]*Exchange, len(allowedExchanges))
    for i, exchangeName := range allowedExchanges {
        apiKey := exchangeInfo[exchangeName][1]
        secretKeyEnvVar := strings.ToUpper(exchangeName) + secretKeySuffix
        secretKey := os.Getenv(secretKeyEnvVar)
        if secretKey == "" {
            log.Fatalln("Secret key not provided for %v (searching for %v)", exchangeName, secretKeyEnvVar)
        }
        switch exchangeName {
        case "BinanceUS":
            exchanges[i] = binanceus.NewBinanceUS(apiKey, secretKey, assetPairTranslators["BinanceUS"])
        case "Kraken":
            exchanges[i] = kraken.NewKraken(apiKey, secretKey, assetPairTranslators["Kraken"], assetPairTranslators["ISO4217"])
        case "KuCoin":
            apiPassphrase := os.Getenv("KUCOIN_API_PASSPHRASE")
            if apiPassphrase = "" {
                log.Fatalln("KuCoin API Passphrase not provided")
            }
            exchanges[i] = kucoin.NewKuCoin(apiKey, secretKey, apiPassphrase, assetPairTranslators["KuCoin"])
        default:
            log.Fatalln("Exchange implementation not found for %v", exchangeName)
        }
    }

    // run algo

    killerInstinct := nn.NewKillerInstinct()

    for exchangePair := range util.ExchangeCombinations(exchanges, 2) {
        commonAssetPairs := util.AssetPairIntersection(
            assetPairTranslators[exchangePair[0].String()].GetAssetPairs(),
            assetPairTranslators[exchangePair[1].String()].GetAssetPairs()
        )

        // start go routines and predictions here
        go grizzly(exchangePair[0], exchangePair[1], commonAssetPairs)
    }
}
