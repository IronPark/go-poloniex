package poloniex

type Trade struct {
	GlobalTradeID int     `json:"globalTradeID"`
	TradeID       int     `json:"tradeID"`
	Date          jTime  `json:"date"`
	Type          string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
}
