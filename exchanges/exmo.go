package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ExmoURL is the url to exmo api
const ExmoURL = "https://api.exmo.me/v1/ticker/"

type exmoTicker struct {
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
}

func requestExmoTickers() (map[string]exmoTicker, error) {
	rawTickers := make(map[string]exmoTicker)
	resp, err := http.Get(ExmoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&rawTickers)
	if err != nil {
		return nil, err
	}
	return rawTickers, nil
}

func parseExmoTickers(rawTickers map[string]exmoTicker) (Tickers, error) {
	tickers := Tickers{}
	tickers.Instance = make(map[TickerKey]Ticker, len(rawTickers))
	currencies := make(map[string]struct{})
	for key, rawTicker := range rawTickers {
		baseQuote := strings.Split(key, "_")
		if len(baseQuote) != 2 {
			return Tickers{}, fmt.Errorf("incorrect ticker, %v", key)
		}
		base, quote := baseQuote[0], baseQuote[1]

		buyPrice, err := strconv.ParseFloat(rawTicker.BuyPrice, 64)
		if err != nil {
			return Tickers{}, err
		}
		sellPrice, err := strconv.ParseFloat(rawTicker.SellPrice, 64)
		if err != nil {
			return Tickers{}, err
		}
		tk := TickerKey{Base: base, Quote: quote}
		tickers.Instance[tk] = Ticker{BuyPrice: buyPrice, SellPrice: sellPrice}
		currencies[base] = struct{}{}
		currencies[quote] = struct{}{}
	}
	currenciesSlice := make([]string, 0, len(currencies))
	for cur := range currencies {
		currenciesSlice = append(currenciesSlice, cur)
	}
	tickers.Currencies = currenciesSlice
	return tickers, nil
}

// GetExmoTickers do get request to tickers exmo api, then it parses result
func GetExmoTickers() (Tickers, error) {
	rawTickers, err := requestExmoTickers()
	if err != nil {
		return Tickers{}, err
	}
	tickers, err := parseExmoTickers(rawTickers)
	if err != nil {
		return Tickers{}, err
	}
	return tickers, nil
}
