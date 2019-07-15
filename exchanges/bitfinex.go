package exchanges

import (
	"log"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

// BitfinexComission is commision that is being taken after any order execution
const BitfinexComission = 0.002

type Bitfinex struct{}

var _ Exchange = &Bitfinex{}

func (e *Bitfinex) GetName() string {
	return "Bitfinex"
}

func (e *Bitfinex) GetCommission() float64 {
	return BitfinexComission
}

// GetTickers do get request to tickers bitfinex api, then it parses result
func (e *Bitfinex) GetTickers() (*Tickers, error) {
	client := rest.NewClient()
	bookTickers, err := client.Tickers.All()
	if err != nil {
		return nil, err
	}

	tickers := NewTickers(len(*bookTickers))
	for _, ticker := range *bookTickers {
		symbol := ticker.Symbol
		if len(symbol) == 0 || symbol[0] != 't' || strings.Contains(symbol, ":") {
			continue
		}
		symbol = symbol[1:]
		if len(symbol) != 6 {
			log.Printf("warn, bitfinex non 6 letter symbol with no semicolon %q", symbol)
			continue
		}
		base, quote := symbol[:3], symbol[3:]
		tickers.add(base, quote, Ticker{BuyPrice: ticker.Bid, SellPrice: ticker.Ask})
	}
	return tickers, nil
}
