package poloniex

type Balance struct {
	Currency  string
	Balance   float64
	Available float64
	Pending   float64
	Value     float64
}