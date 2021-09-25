package app

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

// WSBinanceMsg one message from binance
type WSBinanceMsg struct {
	Code   int    `json:"code"`   // error code
	Msg    string `json:"msg"`    // error msg
	Result string `json:"result"` // result
	ID     int    `json:"id"`     // msg id
}

const BINANCE_API = "https://api.binance.com/api/v3"
const SYMBOLS = "/exchangeInfo"
const WSS = "wss://stream.binance.com:9443/ws"
const DEPTH = "/SYMBOL@depth"

func CreateWSConn(address string) (*WSConnection, error) {
	ctx, cc := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cc()

	conn, _, e := websocket.DefaultDialer.DialContext(ctx, address, nil)
	if e != nil {
		return nil, e
	}

	return &WSConnection{Conn: conn, ID: StringWithCharset(4)}, nil
}

func (app *Application) ConnToBinance() {
	conn, e := CreateWSConn(WSS)
	if e != nil {
		app.ELog.Fatal(e)
		return
	}

	conn.API = WSAPI_BINANCE_ID
	app.m.Lock()
	app.WSSConns[conn.ID] = conn
	app.m.Unlock()

	go app.WSSConns[conn.ID].HandleConnMsg(app)
}
