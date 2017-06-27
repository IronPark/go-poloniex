// Package Poloniex is an implementation of the Poloniex API in Golang.
package poloniex

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"github.com/ironpark/coinex/trader"
	"strconv"
)

const (
	API_BASE                   = "https://poloniex.com/" // Poloniex API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                      // HTTP client timeout
)

// New return a instantiate poloniex struct
func New(apiKey, apiSecret string) *Poloniex {
	client := NewClient(apiKey, apiSecret)
	return &Poloniex{client}
}

// poloniex represent a poloniex client
type Poloniex struct {
	client *client
}

// GetTickers is used to get the ticker for all markets
func (b *Poloniex) GetTickers() (tickers map[string]Ticker, err error) {
	r, err := b.client.do("GET", "public?command=returnTicker", "", false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &tickers); err != nil {
		return
	}
	return
}

// GetVolumes is used to get the volume for all markets
func (b *Poloniex) GetVolumes() (vc VolumeCollection, err error) {
	r, err := b.client.do("GET", "public?command=return24hVolume", "", false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &vc); err != nil {
		return
	}
	return
}

func (b *Poloniex) GetCurrencies() (currencies Currencies, err error) {
	r, err := b.client.do("GET", "public?command=returnCurrencies", "", false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &currencies.Pair); err != nil {
		return
	}
	return
}

// GetOrderBook is used to get retrieve the orderbook for a given market
// market: a string literal for the market (ex: BTC_NXT). 'all' not implemented.
// cat: bid, ask or both to identify the type of orderbook to return.
// depth: how deep of an order book to retrieve
func (b *Poloniex) GetOrderBook(market, cat string, depth int) (orderBook OrderBook, err error) {
	// not implemented
	if cat != "bid" && cat != "ask" && cat != "both" {
		cat = "both"
	}
	if depth > 100 {
		depth = 100
	}
	if depth < 1 {
		depth = 1
	}

	r, err := b.client.do("GET", fmt.Sprintf("public?command=returnOrderBook&currencyPair=%s&depth=%d", strings.ToUpper(market), depth), "", false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &orderBook); err != nil {
		return
	}
	if orderBook.Error != "" {
		err = errors.New(orderBook.Error)
		return
	}
	return
}

// Returns candlestick chart data. Required GET parameters are "currencyPair",
// "period" (candlestick period in seconds; valid values are 300, 900, 1800,
// 7200, 14400, and 86400), "start", and "end". "Start" and "end" are given in
// UNIX timestamp format and used to specify the date range for the data
// returned.
func (b *Poloniex) ChartData(currencyPair string, period int, start, end time.Time) (candles []*CandleStick, err error) {
	r, err := b.client.do("GET", fmt.Sprintf(
		"/public?command=returnChartData&currencyPair=%s&period=%d&start=%d&end=%d",
		strings.ToUpper(currencyPair),
		period,
		start.Unix(),
		end.Unix(),
	), "", false)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &candles); err != nil {
		return
	}

	return
}

func (b *Poloniex) MarketHistory(currencyPair string, start, end time.Time) (trades []Trade, err error) {
	r, err := b.client.do("GET", fmt.Sprintf(
		"/public?command=returnTradeHistory&currencyPair=%s&start=%d&end=%d",
		strings.ToUpper(currencyPair),
		start.Unix(),
		end.Unix(),
	), "", false)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &trades); err != nil {
		return
	}

	return
}
//returnTradeHistory

//"available":"5.015","onOrders":"1.0025","btcValue":"0.078"
func (b *Poloniex) GetBalance() (balance []Balance, err error) {
	r, err := b.client.do("POST", "https://poloniex.com/tradingApi?command=returnCompleteBalances","", false)
	balance = []Balance{}
	if err != nil {
		return
	}
	response := make(map[string]interface{})
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}

	if response["error"] != nil {
		err = errors.New(response["error"].(string))
		return
	}

	for k, v := range response {
		values := v.(map[string]interface{})
		available, _ := strconv.ParseFloat(values["available"].(string), 64)
		onOders, _ := strconv.ParseFloat(values["onOrders"].(string), 64)
		btc, _ := strconv.ParseFloat(values["btcValue"].(string), 64)
		balance = append(balance,Balance{
			Currency:  k,
			Balance:   onOders+available,
			Available: available,
			Pending:   onOders,
			Value:     btc,
		})
	}

	return
}
