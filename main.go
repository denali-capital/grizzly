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
	allowedExchanges := util.CaseInsensitiveIntersection(implementedExchanges, zippedExchangeList[0][1:], true)

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

	zippedAssetPairsList := util.Zip(assetPairsList...)
	for i, _ := range zippedAssetPairsList[0][1:] {
		for exchangeName, translator := range assetPairTranslators {
			translator[i + 1] := zippedAssetPairsList[exchangeIndices[exchangeName]][i + 1]
		}
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}

	exchanges := make([]*Exchange, len(allowedExchanges))
	for i, exchangeName := range allowedExchanges {
		apiKey := exchangeInfo[exchangeName][1]
		secretKey := os.Getenv(strings.ToUpper(exchangeName) + secretKeySuffix)
		if secretKey == "" {
			log.Fatalln("Secret key not provided for %v", exchangeName)
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
}
