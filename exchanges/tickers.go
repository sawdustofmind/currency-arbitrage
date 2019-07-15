package exchanges

import (
	"fmt"
	"math"
	"strings"
)

type tickerKey struct {
	Base  string
	Quote string
}

type Ticker struct {
	BuyPrice  float64
	SellPrice float64
}

// Tickers contrains itself with its currencies
type Tickers struct {
	Currencies      []string
	currencyToIndex map[string]int

	tickers map[tickerKey]Ticker
}

func NewTickers(initialSize int) *Tickers {
	return &Tickers{
		tickers:         make(map[tickerKey]Ticker),
		currencyToIndex: make(map[string]int),
	}
}

func (t *Tickers) add(base, quote string, ticker Ticker) {
	t.putCurrency(base)
	t.putCurrency(quote)
	tk := tickerKey{Base: base, Quote: quote}
	t.tickers[tk] = ticker
}

func (t *Tickers) putCurrency(currency string) {
	_, ok := t.currencyToIndex[currency]
	if !ok {
		t.Currencies = append(t.Currencies, currency)
		t.currencyToIndex[currency] = len(t.Currencies) - 1
	}
}

func (t *Tickers) Get(base, quote string) *Ticker {
	tk := tickerKey{Base: base, Quote: quote}
	ticker, ok := t.tickers[tk]
	if !ok {
		return nil
	}
	return &ticker
}

// ToEdges converts tickers to graph that can be used in arbalgo
func (t *Tickers) ToEdges(commission float64) map[int]map[int]float64 {
	edges := make(map[int]map[int]float64)
	for tk, ticker := range t.tickers {
		baseIndex := t.currencyToIndex[tk.Base]
		quoteIndex := t.currencyToIndex[tk.Quote]
		if _, ok := edges[baseIndex]; !ok {
			edges[baseIndex] = make(map[int]float64)
		}
		if _, ok := edges[quoteIndex]; !ok {
			edges[quoteIndex] = make(map[int]float64)
		}

		edges[baseIndex][quoteIndex] = -math.Log(ticker.BuyPrice * (1 - commission))
		edges[quoteIndex][baseIndex] = math.Log(ticker.SellPrice / (1 - commission))
	}
	return edges
}

// GetCurrencyPath makes currencies lexicographically ordered representation of index path
func (t *Tickers) GetCurrencyPath(path []int) string {
	path = t.arrangedCyclePath(path)
	currencies := make([]string, 0, len(path))
	for _, i := range path {
		currencies = append(currencies, t.Currencies[i])
	}
	return strings.Join(currencies, "->")
}

// GetPricePath makes price report of index path
func (t *Tickers) GetPricePath(path []int, commission float64) (string, float64, error) {
	path = t.arrangedCyclePath(path)
	finalPrice := 1.0
	reprSlice := make([]string, 0, len(path))
	for i := 0; i < len(path)-1; i++ {
		repr, price, err := t.getTickerPricePath(path[i], path[i+1], commission)
		if err != nil {
			return "", 0, err
		}
		reprSlice = append(reprSlice, repr)
		finalPrice *= price
	}
	repr, price, err := t.getTickerPricePath(path[len(path)-1], path[0], commission)
	if err != nil {
		return "", 0, err
	}
	reprSlice = append(reprSlice, repr)
	finalPrice *= price
	return strings.Join(reprSlice, "*"), finalPrice, nil
}

func (t *Tickers) getTickerPricePath(i, j int, commission float64) (string, float64, error) {
	sideBuy := true
	ticker := t.Get(t.Currencies[i], t.Currencies[j])
	if ticker == nil {
		sideBuy = false
		ticker = t.Get(t.Currencies[j], t.Currencies[i])
		if ticker == nil {
			return "", 0, fmt.Errorf("ticker beetwean %s and %s not exists", t.Currencies[i], t.Currencies[j])
		}
	}

	var repr string
	var price float64
	if sideBuy {
		price = ticker.BuyPrice * (1 - commission)
		repr = fmt.Sprintf("(%.8f*%.3f)", ticker.BuyPrice, (1 - commission))
	} else {
		price = 1 / (ticker.SellPrice * (1 - commission))
		repr = fmt.Sprintf("(1/%.8f/%.3f)", ticker.SellPrice, commission)
	}
	return repr, price, nil
}

func (t *Tickers) arrangedCyclePath(path []int) []int {
	pivotI := 0
	for i, v := range path {
		if t.Currencies[v] < t.Currencies[path[pivotI]] {
			pivotI = i
		}
	}
	cpy := make([]int, len(path))
	copy(cpy, path)
	return append(cpy[pivotI:], cpy[:pivotI]...)
}
