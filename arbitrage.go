package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sawdustofmind/currency-arbitrage/arbalgo"
	"github.com/sawdustofmind/currency-arbitrage/exchanges"
)

const checkInterval = 5 * time.Second

type exchangeInfo struct {
	Route string `json:"route,omitempty"`
	Name  string `json:"name,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("PORT must be set as first command-line argument")
	}
	port := os.Args[1]
	log.Println("info, starting server at", port)

	es := []exchanges.Exchange{&exchanges.Exmo{}, &exchanges.Binance{}}
	exchangesInfo := make([]exchangeInfo, 0, len(es))
	for _, ex := range es {
		store := ArbitrageHistoryStore{}
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			entries := store.Get()
			if len(entries) == 0 {
				_, _ = w.Write([]byte("[]"))
				return
			}
			bytes, err := json.Marshal(entries)
			if err != nil {
				bytes, _ = json.Marshal(err.Error())
				w.WriteHeader(400)
			}
			_, _ = w.Write(bytes)
		}
		route := "/history/" + strings.ToLower(ex.GetName())
		http.HandleFunc(route, handler)
		exchangesInfo = append(exchangesInfo, exchangeInfo{Route: route, Name: ex.GetName()})
		startChecking(ex, &store, checkInterval)
	}
	http.HandleFunc("/exchanges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		bytes, _ := json.Marshal(exchangesInfo)
		w.Write(bytes)
	})
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func startChecking(exchange exchanges.Exchange, store *ArbitrageHistoryStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			checkOnArbitrage(exchange, store)
			<-ticker.C
		}
	}()
}

func checkOnArbitrage(exchange exchanges.Exchange, store *ArbitrageHistoryStore) {
	now := time.Now()
	tickers, err := exchange.GetTickers()
	if err != nil {
		log.Printf("error, %v api request tickers fail, err=%v", exchange.GetName(), err)
		return
	}
	shortestPaths := arbalgo.FloydWarshall(tickers.ToEdges(exchange.GetCommission()))
	if !shortestPaths.HasCycles() {
		log.Printf("info, %s - no currency arbitrage", exchange.GetName())
		return
	}
	V := len(tickers.Currencies)
	cycles := arbalgo.FindUniqueCycles(shortestPaths, V)
	for _, cycle := range cycles {
		pricePath, price, _ := tickers.GetPricePath(cycle, exchanges.ExmoCommission)
		currencyPath := tickers.GetCurrencyPath(cycle)
		profit := (price - 1) * 100

		log.Printf("info, %s - found currency arbitrage:%s, profit:%.4f%% [%s]", exchange.GetName(), currencyPath, profit, pricePath)
		entry := ArbitrageHistoryEntry{
			Time:   now,
			Cycle:  currencyPath,
			Report: pricePath,
			Profit: profit,
		}
		store.Add(&entry)
	}
}
