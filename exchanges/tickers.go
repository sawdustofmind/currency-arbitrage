package exchanges

import "math"

type TickerKey struct {
	Base  string
	Quote string
}

type Ticker struct {
	BuyPrice  float64
	SellPrice float64
}

// Tickers contrains itself with its currencies
type Tickers struct {
	Currencies []string
	Instance   map[TickerKey]Ticker

	currencyToInt map[string]int
}

func (t *Tickers) getCurrencyIndex(currency string) int {
	if t.currencyToInt == nil {
		t.currencyToInt = make(map[string]int)
		for i, curr := range t.Currencies {
			t.currencyToInt[curr] = i
		}
	}
	return t.currencyToInt[currency]
}

func (t *Tickers) ToEdges() map[int]map[int]float64 {
	edges := make(map[int]map[int]float64)
	for tk, ticker := range t.Instance {
		baseIndex := t.getCurrencyIndex(tk.Base)
		quoteIndex := t.getCurrencyIndex(tk.Quote)
		if _, ok := edges[baseIndex]; !ok {
			edges[baseIndex] = make(map[int]float64)
		}
		if _, ok := edges[quoteIndex]; !ok {
			edges[quoteIndex] = make(map[int]float64)
		}

		edges[baseIndex][quoteIndex] = -math.Log(ticker.BuyPrice)
		edges[quoteIndex][baseIndex] = math.Log(ticker.SellPrice)
	}
	return edges
}
