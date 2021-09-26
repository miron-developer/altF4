package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// exported consts
const (
	// WSMessage types
	WSM_SEND_USERID  = 1        // send user his id
	WSM_UPDATE_BOOK  = 11       // updater
	WSM_ERROR_ALL    = -10      // if error was
	WSM_ERROR_SUBS   = -11      // if subscribe error
	WSM_SERVER_TOKEN = "server" // server error on wsmessage receiver/addresser id
	WSM_RECEIVE_ALL  = "all"    // send to all person

	// api's
	WSAPI_BINANCE_ID = 100 // binance api id

	// for ws ping pong work & connections
	WSC_WRITE_WAIT  = 10 * time.Second
	WSC_PONG_WAIT   = 60 * time.Second
	WSC_PING_PERIOD = (WSC_PONG_WAIT * 9) / 10
)

// WSMessage one message from ws connection to users
type WSMessage struct {
	MsgType        int         `json:"msgType"`
	AddresserToken string      `json:"addresser"`
	ReceiverToken  string      `json:"receiver"`
	Body           interface{} `json:"body"`
}

// WSConnection is one ws connection to api
type WSConnection struct {
	Conn           *websocket.Conn
	Token          string
	ID             int
	API            int
	UnsubsFunction func(id int, symbol string) interface{}
}

// WSUser is one ws connection user
type WSUser struct {
	Conn              *websocket.Conn
	Token             string
	ListeningApiToken string
	ListeningSymbol   string
	ID                int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func (app *Application) findByID(finder, searchPlace int) interface{} {
	if searchPlace == 1 {
		for _, user := range app.OnlineUsers {
			if user.ID == finder {
				return user
			}
		}
		return nil
	}

	for _, conn := range app.WSSConns {
		if conn.ID == finder {
			return conn
		}
	}
	return nil
}

// get api token by api type
func (app *Application) GetApiToken(apiID int) string {
	for _, v := range app.WSSConns {
		if v.API == apiID {
			return v.Token
		}
	}
	return ""
}

// CreateWSConn create one api ws connection
func (app *Application) CreateWSConn(address string) (*WSConnection, error) {
	ctx, cc := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cc()

	conn, _, e := websocket.DefaultDialer.DialContext(ctx, address, nil)
	if e != nil {
		return nil, e
	}

	// generating unical id
	id := -1
	for {
		id = RandomInt()
		if exist := app.findByID(id, 2); exist == nil {
			break
		}
	}

	return &WSConnection{Conn: conn, ID: id, Token: StringWithCharset(4)}, nil
}

// HandleUserMsg handle received msg from front user
func (user *WSUser) HandleUserMsg(app *Application) {
	user.Conn.SetPongHandler(
		func(string) error {
			user.Conn.SetReadDeadline(time.Now().Add(WSC_PONG_WAIT))
			return nil
		},
	)
}

// HandleConnMsg handle received msg api
func (conn *WSConnection) HandleConnMsg(app *Application) {
	conn.Conn.SetPingHandler(func(appData string) error {
		return conn.Conn.WriteMessage(websocket.PongMessage, nil)
	})

	for {
		appMsg := &WSMessage{AddresserToken: WSM_SERVER_TOKEN}
		apiMsg := createWsMessageByAPI(conn.API)
		if e := conn.Conn.ReadJSON(&apiMsg); e != nil {
			app.m.Lock()
			delete(app.WSSConns, conn.Token)
			app.m.Unlock()

			app.ELog.Println(e)
			conn.Conn.Close()
			return
		}
		appMsg.Body = apiMsg

		// TODO: create unical point
		mapMsg, ok := apiMsg.(map[string]interface{})
		if !ok {
			continue
		}

		symbol, ok := mapMsg["s"].(string)
		if symbol != "" && ok {
			appMsg.MsgType = WSM_UPDATE_BOOK
			go func() {
				for _, v := range app.OnlineUsers {
					if v.ListeningSymbol == symbol {
						v.Conn.WriteJSON(appMsg)
					}
				}
			}()
		}
	}
}

// Pinger server ping every pingPeriod
func (user *WSUser) Pinger() {
	ticker := time.NewTicker(WSC_PING_PERIOD)
	for {
		<-ticker.C
		if e := user.Conn.WriteMessage(websocket.PingMessage, nil); e != nil {
			return
		}
	}
}

// create message by api id
func createWsMessageByAPI(apiID int) interface{} {
	if apiID == 100 {
		return WSBinanceMsg{}
	}
	return nil
}
