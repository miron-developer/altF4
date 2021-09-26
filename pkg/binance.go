package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// WSBinanceMsg one message from binance
type WSBinanceMsg struct {
	Code   int    `json:"code"`   // error code
	Msg    string `json:"msg"`    // error msg
	Result string `json:"result"` // result
	ID     int    `json:"id"`     // msg id
}

// BinanceSymbol describe one symbol info
type BinanceSymbol struct {
	Symbol                   string `json:"symbol"`
	Status                   string `json:"status"`
	BaseAsset                string `json:"baseAsset"`
	BaseAssetPrecision       int    `json:"baseAssetPrecision"`
	BaseCommissionPrecision  int    `json:"baseCommissionPrecision"`
	QuoteAsset               string `json:"quoteAsset"`
	QuoteAssetPrecision      int    `json:"quoteAssetPrecision"`
	QuoteCommissionPrecision int    `json:"quoteCommissionPrecision"`
	// Filters                  []interface{} `json:"filters"`
}

// BinanceExchanges describe response from "https://api.binance.com/api/v3/exchangeInfo"
type BinanceExchanges struct {
	TimeZone   string          `json:"timezone"`
	ServerTime int             `json:"serverTime"`
	Symbols    []BinanceSymbol `json:"symbols"`
}

// BinanceOrderBook describe response from "https://api.binance.com/api/v3/depth?symbol=SYMBOL&limit=LIMIT"
type BinanceOrderBook struct {
	LastUpdatedID int           `json:"lastUpdateId"`
	Bids          []interface{} `json:"bids"`
	Asks          []interface{} `json:"asks"`
}

// BinanceSubsUnsubs describe message to subscribe and unsubscribe user
type BinanceSubsUnsubs struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ID     int         `json:"id"`

	// on depth update
	EventTime     int        `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int        `json:"U"`
	LastUpdateID  int        `json:"u"`
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}

const BINANCE_API = "https://api.binance.com/api/v3"
const SYMBOLS = "/exchangeInfo"
const ORDER_BOOK = "/depth?symbol=SYMBOL&limit=LIMIT"
const WSS = "wss://stream.binance.com:9443/ws"

func BinanceUnsubcribeMsg(id int, symbol string) interface{} {
	return &BinanceSubsUnsubs{
		Method: "UNSUBSCRIBE",
		ID:     id,
		Params: []interface{}{strings.ToLower(symbol) + "@depth"},
	}
}

func (app *Application) ConnToBinance() {
	conn, e := app.CreateWSConn(WSS)
	if e != nil {
		app.ELog.Fatal(e)
		return
	}

	conn.API = WSAPI_BINANCE_ID
	conn.UnsubsFunction = BinanceUnsubcribeMsg
	app.m.Lock()
	app.WSSConns[conn.Token] = conn
	app.m.Unlock()

	go app.WSSConns[conn.Token].HandleConnMsg(app)
}

// response example
// {
//     "timezone": "UTC",
//     "serverTime": 1632586812823,
//     "rateLimits": [
//         {
//             "rateLimitType": "REQUEST_WEIGHT",
//             "interval": "MINUTE",
//             "intervalNum": 1,
//             "limit": 1200
//         },
//         {
//             "rateLimitType": "ORDERS",
//             "interval": "SECOND",
//             "intervalNum": 10,
//             "limit": 50
//         },
//         {
//             "rateLimitType": "ORDERS",
//             "interval": "DAY",
//             "intervalNum": 1,
//             "limit": 160000
//         },
//         {
//             "rateLimitType": "RAW_REQUESTS",
//             "interval": "MINUTE",
//             "intervalNum": 5,
//             "limit": 6100
//         }
//     ],
//     "exchangeFilters": [],
//     "symbols": [{
// 		"symbol": "ETHBTC",
// 		"status": "TRADING",
// 		"baseAsset": "ETH",
// 		"baseAssetPrecision": 8,
// 		"quoteAsset": "BTC",
// 		"quoteAssetPrecision": 8,
// 		"baseCommissionPrecision": 8,
// 		"quoteCommissionPrecision": 8,
// 		"filters": [
// 			{
// 				"filterType": "PRICE_FILTER",
// 				"minPrice": "0.00000100",
// 				"maxPrice": "922327.00000000",
// 				"tickSize": "0.00000100"
// 			},
// 			{
// 				"filterType": "PERCENT_PRICE",
// 				"multiplierUp": "5",
// 				"multiplierDown": "0.2",
// 				"avgPriceMins": 5
// 			},
// 			{
// 				"filterType": "LOT_SIZE",
// 				"minQty": "0.00010000",
// 				"maxQty": "100000.00000000",
// 				"stepSize": "0.00010000"
// 			},
// 			{
// 				"filterType": "MIN_NOTIONAL",
// 				"minNotional": "0.00010000",
// 				"applyToMarket": true,
// 				"avgPriceMins": 5
// 			},
// 			{
// 				"filterType": "ICEBERG_PARTS",
// 				"limit": 10
// 			},
// 			{
// 				"filterType": "MARKET_LOT_SIZE",
// 				"minQty": "0.00000000",
// 				"maxQty": "952.00626854",
// 				"stepSize": "0.00000000"
// 			},
// 			{
// 				"filterType": "MAX_NUM_ORDERS",
// 				"maxNumOrders": 200
// 			},
// 			{
// 				"filterType": "MAX_NUM_ALGO_ORDERS",
// 				"maxNumAlgoOrders": 5
// 			}
// 		],
// 		"permissions": [
// 			"SPOT",
// 			"MARGIN"
// 		]
// 	}]
// }

// GetBinanceCurrencies get Binance exchange info
func (app *Application) GetBinanceCurrencies() (interface{}, error) {
	resp, e := http.Get(BINANCE_API + SYMBOLS)
	if e != nil {
		app.ELog.Println("Binance get error:", e)
		return nil, errors.New("500: binance api is not avalable")
	}

	r, e := io.ReadAll(resp.Body)
	if e != nil {
		app.ELog.Println("Binance read body error:", e)
		return nil, errors.New("500: binance api is not avalable")
	}

	data := &BinanceExchanges{}
	if e := json.Unmarshal(r, data); e != nil {
		app.ELog.Println("Binance unmarshal error:", e)
		return nil, errors.New("500: server error")
	}

	// common binance symbol to general symbol
	symbols := []*ExchangeSymbol{}
	for _, v := range data.Symbols {
		symbol := &ExchangeSymbol{}
		if e := FillStructFromArr(symbol, MakeArrFromStruct(v)); e != nil {
			app.ELog.Println("Binance common error:", e)
			return nil, errors.New("500: server error, try later")
		}
		symbols = append(symbols, symbol)
	}

	// here exist second way
	// the 2nd way is just send data to front w/out unmarshal
	// but its not work, when we have many exchange points
	// cuz we must to common them
	return symbols, nil
}

// response example
// {
//     "lastUpdateId": 13911799197,
// 	"bids": [
// 		[
// 			1,
// 			2,
// 		],
// 		[
// 			1,
// 			2,
// 		]
// 	],
// 	"asks": [
// 		[
// 			1,
// 			2,
// 		],
// 		[
// 			1,
// 			2,
// 		]
// 	]
// }

// BinanceSubscribeCurrencie get order book
func (app *Application) BinanceSubscribeCurrencie(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	symbol := r.FormValue("symbol")
	if symbol == "" {
		return nil, errors.New("400: wrong symbol")
	}

	limit := r.FormValue("limit")
	if limit == "" {
		limit = "20"
	}

	resp, e := http.Get(BINANCE_API + strings.ReplaceAll(strings.ReplaceAll(ORDER_BOOK, "SYMBOL", symbol), "LIMIT", limit))
	if e != nil {
		app.ELog.Println("Binance subscribe/get error:", e)
		return nil, errors.New("500: binance api is not avalable")
	}

	res, e := io.ReadAll(resp.Body)
	if e != nil {
		app.ELog.Println("Binance subscribe/read body error:", e)
		return nil, errors.New("500: binance api is not avalable")
	}

	data := &BinanceOrderBook{}
	if e := json.Unmarshal(res, data); e != nil {
		app.ELog.Println("Binance subscribe/unmarshal error:", e)
		return nil, errors.New("500: server error")
	}

	// common binance book to general book
	orderBook := &OrderBook{}
	if e := FillStructFromArr(orderBook, MakeArrFromStruct(*data)); e != nil {
		app.ELog.Println("Binance common error:", e)
		return nil, errors.New("500: server error, try later")
	}

	go func() {
		user, ok := app.OnlineUsers[r.FormValue("token")]
		if !ok {
			return
		}

		// get connection
		var conn *WSConnection
		for _, v := range app.WSSConns {
			if v.API == WSAPI_BINANCE_ID {
				conn = v
			}
		}

		// subscribe user to get depths
		if e := conn.Conn.WriteJSON(&BinanceSubsUnsubs{
			ID:     user.ID,
			Method: "SUBSCRIBE",
			Params: []interface{}{strings.ToLower(symbol) + "@depth"},
		}); e != nil {
			app.ELog.Println("Binance subscribe error: ", e)
			user.Conn.WriteJSON(&WSMessage{MsgType: WSM_ERROR_SUBS, AddresserToken: WSM_SERVER_TOKEN, ReceiverToken: user.Token})
			return
		}

		user.ListeningSymbol = symbol
	}()

	// here exist second way
	// the 2nd way is just send data to front w/out unmarshal
	// but its not work, when we have many exchange points
	// cuz we must to common them
	return orderBook, nil
}

func (app *Application) BinanceUnsubscribeCurrencie(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	symbol := r.FormValue("symbol")
	if symbol == "" {
		return nil, errors.New("400: wrong symbol")
	}

	user, ok := app.OnlineUsers[r.FormValue("token")]
	if !ok {
		return nil, errors.New("wrong user token")
	}

	// get connection
	var conn *WSConnection
	for _, v := range app.WSSConns {
		if v.API == WSAPI_BINANCE_ID {
			conn = v
		}
	}

	// subscribe user to get depths
	if e := conn.Conn.WriteJSON(BinanceUnsubcribeMsg(user.ID, symbol)); e != nil {
		app.ELog.Println("Binance unsubscribe error: ", e)
		return nil, errors.New("binance unsubscribe error, try later")
	}

	user.ListeningSymbol = ""

	return nil, nil
}
