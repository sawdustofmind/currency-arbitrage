package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
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
	shortestPaths := arbalgo.FloydWarshall(tickers.ToEdges())
	if !shortestPaths.HasCycles() {
		log.Println("info, no currency arbitrage")
		return
	}
	V := len(tickers.Currencies)
	cycles := FindUniqueCycles(shortestPaths, V)
	for _, cycle := range cycles {
		log.Println("info, found currency arbitrage,", GetPath(cycle, tickers.Currencies))
	}
}

func FindUniqueCycles(p *arbalgo.ShortestPaths, V int) [][]int {
	cycles := [][]int{}
	for i := 0; i < V; i++ {
		for j := 0; j < V; j++ {
			_, cycle, _ := p.Get(i, j)
			if len(cycle) == 0 {
				continue
			}
			contains := false
			for _, c := range cycles {
				if SameCycles(c, cycle) {
					contains = true
					break
				}
			}
			if !contains {
				cycles = append(cycles, cycle)
			}
		}
	}
	return cycles
}

func SameCycles(first []int, second []int) bool {
	firstCopy := first[:]
	sort.Ints(firstCopy)
	secondCopy := second[:]
	sort.Ints(secondCopy)
	return reflect.DeepEqual(firstCopy, secondCopy)
}

func GetPricePath(tickers exchanges.Tickers, intPath []int) (float64, string) {
	getPrice := func(i, j int) (float64, string) {
		tk := exchanges.TickerKey{Base: tickers.Currencies[i], Quote: tickers.Currencies[j]}
		ticker, ok := tickers.Instance[tk]
		if !ok {
			tk = exchanges.TickerKey{Base: tickers.Currencies[j], Quote: tickers.Currencies[i]}
			ticker, ok = tickers.Instance[tk]
			if !ok {
				panic(fmt.Sprint(tk))
			}
		}
		if ok {
			buyPrice := ticker.BuyPrice * exchanges.ExmoComission
			buyPriceStr := fmt.Sprintf("(%.8f*%.3f)", ticker.BuyPrice, exchanges.ExmoComission)
			return buyPrice, buyPriceStr
		}
		sellPrice := 1 / (ticker.SellPrice * exchanges.ExmoComission)
		sellPriceStr := fmt.Sprintf("(1/%.8f/%.3f)", ticker.SellPrice, exchanges.ExmoComission)
		return sellPrice, sellPriceStr
	}

	finalPrice := 1.0
	priceSlice := make([]string, 0, len(intPath)-1)
	for i := 0; i < len(intPath)-1; i++ {
		price, priceStr := getPrice(intPath[i], intPath[i+1])
		priceSlice = append(priceSlice, priceStr)
		finalPrice *= price
	}
	return finalPrice, strings.Join(priceSlice, "*")
}

func GetPath(intPath []int, slice []string) string {
	curList := make([]string, 0, len(intPath))
	for _, i := range intPath {
		curList = append(curList, slice[i])
	}
	return strings.Join(curList, "->")
}
