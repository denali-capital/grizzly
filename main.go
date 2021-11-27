package main

import (
    "log"
    "os"
    "strings"

    "github.com/denali-capital/grizzly/exchanges"
    "github.com/denali-capital/grizzly/types"
    "github.com/denali-capital/grizzly/util"
    "github.com/joho/godotenv"
)

// each exchange will have their own module that implements Exchange interface above
// can have module to compute statistics in background and access them
// Neural network module that handles predict / fit functionality

const configPath string = "config"
const secretKeySuffix string = "_SECRET_KEY"

func main() {
    // need to create array of exchange objects
    // pass this array into function that computes statistics in background

    implementedExchanges := util.DiscoverTypes("github.com/denali-capital/grizzly/exchanges")

    exchangeList := util.ReadCsvFile(configPath + "/exchanges.csv")
    if exchangeList[0][0] != "exchange" || exchangeList[0][1] != "api_key" {
        log.Fatalln("Labels must be \"exchange,api_key,...\"")
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
            exchangeInfo[exchangeName] := info
        }
    }

    assetPairsList := util.ReadCsvFile(configPath + "/assetpairs.csv")
    if assetPairsList[0][0] != "canonical" {
        log.Fatalln("Labels must be \"canonical,...\"")
    }

    exchangeIndices := make(map[string]uint)
    assetPairTranslators := make(map[string]types.AssetPairTranslator)
    for i, exchangeName := range assetPairsList[0][1:] {
        exchangeIndices[exchangeName] = i + 1
        assetPairTranslators[exchangeName] := make(types.AssetPairTranslator)
    }

    assetPairCanonicalTranslator := make(map[types.AssetPair]string)
    zippedAssetPairsList := util.Zip(assetPairsList...)
    for i, assetPairCanonical := range zippedAssetPairsList[0][1:] {
        assetPairIndices[i + 1] = assetPairCanonical
        for exchangeName, translator := range assetPairTranslators {
            assetPairSpecific := zippedAssetPairsList[exchangeIndices[exchangeName]][i + 1]
            if assetPairSpecific != "" {
                translator[i + 1] = assetPairSpecific
            }
        }
    }

    err := godotenv.Load()
    if err != nil {
        log.Fatalln("Error loading .env file", err)
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
            exchanges[i] = exchanges.NewBinanceUS(apiKey, secretKey, assetPairTranslators["BinanceUS"])
        case "Kraken":
            exchanges[i] = exchanges.NewKraken(apiKey, secretKey, assetPairTranslators["Kraken"])
        default:
            log.Fatalln("Exchange implementation not found for %v", exchangeName)
        }
    }

    // run algo

    for exchangePair := range util.ExchangeCombinations(exchanges, 2) {
        commonAssetPairs := util.AssetPairIntersection(
            assetPairTranslators[exchangePair[0].String()].GetAssetPairs(),
            assetPairTranslators[exchangePair[1].String()].GetAssetPairs()
        )

        // start go routine and predictions here
        // todo: add reading from fees.csv
    }
}
