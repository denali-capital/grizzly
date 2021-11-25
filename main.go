package main

import (
	"log"

	"github.com/denali-capital/grizzly/exchanges"
	"github.com/denali-capital/grizzly/util"
)

// each exchange will have their own module that implements Exchange interface above
// can have module to compute statistics in background and access them
// Neural network module that handles predict / fit functionality

const configPath string = "config"

func main() {
	// need to create array of exchange objects
	// pass this array into function that computes statistics in background

	implementedExchanges := util.DiscoverTypes("github.com/denali-capital/grizzly/exchanges")

	exchangeList := util.ReadCsvFile(configPath + "/exchanges.csv")
	if exchangeList[0] != "exchange" || exchangeList[1] != "api_key" {
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
	if assetPairsList[0] != "canonical" {
		log.Fatalln("Labels must be \"canonical,...\"")
	}

	// make map of exchange name to AssetPairTranslators to use below

	exchanges := make([]*Exchange, len(allowedExchanges))
	for i, exchangeName := range allowedExchanges {
		switch exchangeName {
		case "BinanceUS":
			exchanges[i] = exchanges.NewBinanceUS()
		case "Kraken":
			exchanges[i] = exchanges.NewKraken()
		default:
			log.Fatalln("Exchange implementation not found for %v", exchangeName)
		}
	}

	// run algo
}
