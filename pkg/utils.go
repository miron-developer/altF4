package app

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// ExchangeSymbol general symbol struct
type ExchangeSymbol struct {
	Symbol                   string        `json:"symbol"`
	Status                   string        `json:"status"`
	BaseAsset                string        `json:"baseAsset"`
	BaseAssetPrecision       int           `json:"baseAssetPrecision"`
	BaseCommissionPrecision  int           `json:"baseCommissionPrecision"`
	QuoteAsset               string        `json:"quoteAsset"`
	QuoteAssetPrecision      int           `json:"quoteAssetPrecision"`
	QuoteCommissionPrecision int           `json:"quoteCommissionPrecision"`
	Filters                  []interface{} `json:"filters"`
}

// OrderBook general order book struct
type OrderBook struct {
	LastUpdatedID int           `json:"lastUpdateId"`
	Bids          []interface{} `json:"bids"`
	Asks          []interface{} `json:"asks"`
}

func RandomInt() int {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(999999)
}

func StringWithCharset(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// write in log each request
func logingReq(r *http.Request) string {
	return fmt.Sprintf("%v %v: '%v'\n", r.RemoteAddr, r.Method, r.URL)
}

func (app *Application) GetCurrenciesFromID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	point, e := strconv.Atoi(r.FormValue("point"))
	if r.FormValue("point") == "" || e != nil {
		return nil, errors.New("wrong point")
	}

	user, ok := app.OnlineUsers[r.FormValue("token")]
	if !ok {
		return nil, errors.New("wrong user token")
	}

	// choose by point
	if point == WSAPI_BINANCE_ID {
		user.ListeningApiToken = app.GetApiToken(WSAPI_BINANCE_ID)
		return app.GetBinanceCurrencies()
	}
	return nil, errors.New("wrong point")
}

func (app *Application) SubscribeCurrencie(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	point, e := strconv.Atoi(r.FormValue("point"))
	if r.FormValue("point") == "" || e != nil {
		return nil, errors.New("wrong point")
	}

	// choose by point
	if point == WSAPI_BINANCE_ID {
		return app.BinanceSubscribeCurrencie(w, r)
	}
	return nil, errors.New("wrong point")
}

func (app *Application) UnsubscribeCurrencie(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "GET" {
		return nil, errors.New("wrong method")
	}

	point, e := strconv.Atoi(r.FormValue("point"))
	if r.FormValue("point") == "" || e != nil {
		return nil, errors.New("wrong point")
	}

	// choose by point
	if point == WSAPI_BINANCE_ID {
		return app.BinanceUnsubscribeCurrencie(w, r)
	}
	return nil, errors.New("wrong point")
}

func (app *Application) CheckPerMin() {
	timer := time.NewTicker(1 * time.Minute)
	for {
		// manage timer
		<-timer.C
		timer.Reset(1 * time.Minute)

		// change conf app
		app.CurrentRequestCount = 0
	}
}
