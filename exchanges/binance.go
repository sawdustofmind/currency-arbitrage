package exchanges

import (
	"context"
	"log"
	"strconv"

	"github.com/adshao/go-binance"
)

// BinanceComission is commision that is being taken after any order execution
const BinanceComission = 0.001

type Binance struct{}

var _ Exchange = &Binance{}

func (e *Binance) GetName() string {
	return "Binance"
}

func (e *Binance) GetCommission() float64 {
	return BinanceComission
}

// GetTickers do get request to tickers binance api, then it parses result
func (e *Binance) GetTickers() (*Tickers, error) {
	client := binance.NewClient("", "")
	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	symbolinfo := make(map[string]*binance.Symbol)
	for i := range exchangeInfo.Symbols {
		si := exchangeInfo.Symbols[i]
		symbolinfo[si.Symbol] = &si
	}

	bookTickers, err := client.NewListBookTickersService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	tickers := NewTickers(len(bookTickers))
	for _, ticker := range bookTickers {
		si, ok := symbolinfo[ticker.Symbol]
		if !ok {
			log.Printf("warn, Binance - symbol %s missed in exchange info, skipped", ticker.Symbol)
			continue
		}
		base, quote := si.BaseAsset, si.QuoteAsset
		buyPrice, err := strconv.ParseFloat(ticker.BidPrice, 64)
		if err != nil {
			return nil, err
		}
		sellPrice, err := strconv.ParseFloat(ticker.AskPrice, 64)
		if err != nil {
			return nil, err
		}

		if buyPrice < 1e-6 || sellPrice < 1e-6 {
			continue
		}
		tickers.add(base, quote, Ticker{BuyPrice: buyPrice, SellPrice: sellPrice})
	}
	return tickers, nil
}
