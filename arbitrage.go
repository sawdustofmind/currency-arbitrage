package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sawdustofmind/currency-arbitrage/arbalgo"
	"github.com/sawdustofmind/currency-arbitrage/exchanges"
)

const checkInterval = 5 * time.Second

func main() {
	if len(os.Args) < 2 {
		log.Fatal("PORT must be set as first command-line argument")
	}
	port := os.Args[1]
	log.Println("info, starting server at", port)

	store := ArbitrageHistoryStore{}
	startExmoChecking(&store, checkInterval)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		entries := store.Get()
		if len(entries) == 0 {
			w.Write([]byte("[]"))
			return
		}
		bytes, err := json.Marshal(entries)
		if err != nil {
			bytes, _ = json.Marshal(err.Error())
			w.WriteHeader(400)
		}
		w.Write(bytes)
	}
	http.HandleFunc("/history/exmo", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func startExmoChecking(store *ArbitrageHistoryStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			checkExmoOnArbitrage(store)
			<-ticker.C
		}
	}()
}

func checkExmoOnArbitrage(store *ArbitrageHistoryStore) {
	now := time.Now()
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
		entry := ArbitrageHistoryEntry{
			Time:   now,
			Cycle:  currencyPath,
			Path:   cycle,
			Report: pricePath,
			Profit: profit,
		}
		store.Add(&entry)
	}
}
