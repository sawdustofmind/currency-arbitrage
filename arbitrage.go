package main

import (
	"log"
	"time"

	"github.com/sawdustofmind/currency-arbitrage/arbalgo"
	"github.com/sawdustofmind/currency-arbitrage/exchanges"
)

func main() {
	const checkInterval = 10 * time.Second
	ticker := time.NewTicker(checkInterval)
	for {
		checkExmoOnArbitrage()
		<-ticker.C
	}
}

func checkExmoOnArbitrage() {
	tickers, err := exchanges.GetExmoTickers()
	if err != nil {
		log.Printf("error, exmo api request tickers fail, err=%v", err)
		return
	}
	shortestPaths := arbalgo.FloydWarshall(tickers.ToEdges(exchanges.ExmoComission))
	if !shortestPaths.HasCycles() {
		log.Println("info, no currency arbitrage")
		return
	}
	V := len(tickers.Currencies)
	cycles := arbalgo.FindUniqueCycles(shortestPaths, V)
	for _, cycle := range cycles {
		pricePath, price, _ := tickers.GetPricePath(cycle, exchanges.ExmoComission)
		currencyPath := tickers.GetCurrencyPath(cycle)
		profit := (price - 1) * 100
		log.Printf("info, found currency arbitrage:%s, profit:%.4f%% [%s]", currencyPath, profit, pricePath)
	}
}
