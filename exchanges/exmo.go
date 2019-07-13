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

// ExmoCommission is commission that is being taken after any order execution
const ExmoCommission = 0.002

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

func parseExmoTickers(rawTickers map[string]exmoTicker) (*Tickers, error) {
	tickers := NewTickers(len(rawTickers))
	for key, rawTicker := range rawTickers {
		baseQuote := strings.Split(key, "_")
		if len(baseQuote) != 2 {
			return nil, fmt.Errorf("incorrect ticker, %v", key)
		}
		base, quote := baseQuote[0], baseQuote[1]

		buyPrice, err := strconv.ParseFloat(rawTicker.BuyPrice, 64)
		if err != nil {
			return nil, err
		}
		sellPrice, err := strconv.ParseFloat(rawTicker.SellPrice, 64)
		if err != nil {
			return nil, err
		}
		tickers.add(base, quote, Ticker{BuyPrice: buyPrice, SellPrice: sellPrice})
	}
	return tickers, nil
}

// GetExmoTickers do get request to tickers exmo api, then it parses result
func GetExmoTickers() (*Tickers, error) {
	rawTickers, err := requestExmoTickers()
	if err != nil {
		return nil, err
	}
	tickers, err := parseExmoTickers(rawTickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}
