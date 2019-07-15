package exchanges

type Exchange interface {
	GetName() string
	GetCommission() float64
	GetTickers() (*Tickers, error)
}
