package exchanges

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExmo(t *testing.T) {
	rawTickers := make(map[string]exmoTicker)
	err := json.Unmarshal([]byte(ExmoTickersResponceExample), &rawTickers)
	require.NoError(t, err)
	tickers, err := parseExmoTickers(rawTickers)
	require.NoError(t, err)
	assert.ElementsMatch(t, ExmoTickersExample.Currencies, tickers.Currencies)
	assert.Equal(t, ExmoTickersExample.Instance, tickers.Instance)
}

var ExmoTickersExample = Tickers{
	Currencies: []string{"USD", "RUB", "BTC", "ETH"},
	Instance: map[TickerKey]Ticker{
		{Base: "BTC", Quote: "RUB"}: {BuyPrice: 341999.000037, SellPrice: 342771.74790796},
		{Base: "USD", Quote: "RUB"}: {BuyPrice: 65.25000001, SellPrice: 65.4},
		{Base: "ETH", Quote: "BTC"}: {BuyPrice: 0.03181388, SellPrice: 0.03191428},
		{Base: "ETH", Quote: "RUB"}: {BuyPrice: 10899.05638749, SellPrice: 10940.50783749},
		{Base: "ETH", Quote: "USD"}: {BuyPrice: 166.83422943, SellPrice: 167.6272268},
		{Base: "BTC", Quote: "USD"}: {BuyPrice: 5227.76352776, SellPrice: 5239.97378516},
	},
}

const ExmoTickersResponceExample = `{
    "BTC_USD": {
        "buy_price": "5227.76352776",
        "sell_price": "5239.97378516",
        "last_trade": "5227.76352776",
        "high": "5260",
        "low": "5069.8",
        "avg": "5164.14063257",
        "vol": "783.09376492",
        "vol_curr": "4093829.02310608",
        "updated": 1555479940
    },
    "BTC_RUB": {
        "buy_price": "341999.000037",
        "sell_price": "342771.74790796",
        "last_trade": "342300",
        "high": "343894.8",
        "low": "331059",
        "avg": "337285.54487217",
        "vol": "713.85742492",
        "vol_curr": "244353396.55026183",
        "updated": 1555479940
    },
    "USD_RUB": {
        "buy_price": "65.25000001",
        "sell_price": "65.4",
        "last_trade": "65.4",
        "high": "65.6",
        "low": "65.06443053",
        "avg": "65.3472876",
        "vol": "120986.36746048",
        "vol_curr": "7912508.431916",
        "updated": 1555479934
    },
    "ETH_BTC": {
        "buy_price": "0.03181388",
        "sell_price": "0.03191428",
        "last_trade": "0.03180398",
        "high": "0.03232373",
        "low": "0.0316929",
        "avg": "0.03200206",
        "vol": "5149.29429394",
        "vol_curr": "163.76805273",
        "updated": 1555479939
    },
    "ETH_RUB": {
        "buy_price": "10899.05638749",
        "sell_price": "10940.50783749",
        "last_trade": "10900",
        "high": "11000",
        "low": "10610",
        "avg": "10790.1199925",
        "vol": "1358.94731421",
        "vol_curr": "14812525.7249454",
        "updated": 1555479939
    },
    "ETH_USD": {
        "buy_price": "166.83422943",
        "sell_price": "167.6272268",
        "last_trade": "167.5264",
        "high": "168.40590116",
        "low": "162.06268213",
        "avg": "165.51862843",
        "vol": "4680.11949399",
        "vol_curr": "784043.57039835",
        "updated": 1555479940
    }
}`
