package app

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// exported consts
const (
	// WSMessage types
	WSM_ORDER_BOOK  = 1  // when first get
	WSM_UPDATE_BOOK = 2  // updater
	WSM_ERROR_ALL   = -1 // if error was

	// api's
	WSAPI_BINANCE_ID = 100 // binance api id

	// for ws ping pong work & connections
	WSC_WRITE_WAIT  = 10 * time.Second
	WSC_PONG_WAIT   = 60 * time.Second
	WSC_PING_PERIOD = (WSC_PONG_WAIT * 9) / 10
)

// WSMessage one message from ws connection to users
type WSMessage struct {
	MsgType     int         `json:"msgType"`   // ws message type or api type
	AddresserID string      `json:"addresser"` // who sended, if = -1, then send server
	ReceiverID  string      `json:"receiver"`  // who get, if = -2, then get all, -1 server
	Body        interface{} `json:"body"`      // message body
}

// WSUser is one ws connection user
type WSUser struct {
	Conn *websocket.Conn
	ID   string
}

// WSConnection is one ws connection to api
type WSConnection struct {
	Conn *websocket.Conn
	ID   string
	API  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func (app *Application) findUserByID(finder string) *WSUser {
	for _, user := range app.OnlineUsers {
		if user.ID == finder {
			return user
		}
	}
	return nil
}

// WSWork work with channels
func (app *Application) WSWork() {
	for {
		msg := <-app.WSMessages

		// if send to all
		if msg.ReceiverID == "all" {
			for _, v := range app.OnlineUsers {
				go v.Conn.WriteJSON(msg)
			}
			continue
		}

		receiver := app.findUserByID(msg.ReceiverID)
		if receiver == nil {
			continue
		}

		// futher acting
		go func() {
			if e := receiver.Conn.WriteJSON(msg); e != nil {
				app.ELog.Println("something wrong: " + e.Error())
			}
		}()
	}
}

// HandleUserMsg handle received msg from front user
func (user *WSUser) HandleUserMsg(app *Application) {
	user.Conn.SetPongHandler(
		func(string) error {
			user.Conn.SetReadDeadline(time.Now().Add(WSC_PONG_WAIT))
			return nil
		},
	)

	for {
		msg := &WSMessage{}
		if e := user.Conn.ReadJSON(msg); e != nil {
			app.m.Lock()
			delete(app.OnlineUsers, user.ID)
			app.m.Unlock()

			app.ELog.Println(e)
			user.Conn.Close()
			return
		}
		msg.AddresserID = user.ID
		app.WSMessages <- msg
	}
}

// HandleConnMsg handle received msg api
func (conn *WSConnection) HandleConnMsg(app *Application) {
	conn.Conn.SetPingHandler(func(appData string) error {
		return conn.Conn.WriteMessage(websocket.PongMessage, nil)
	})

	for {
		appMsg := &WSMessage{}
		apiMsg := createWsMessageByAPI(conn.API)
		if e := conn.Conn.ReadJSON(apiMsg); e != nil {
			app.m.Lock()
			delete(app.WSSConns, conn.ID)
			app.m.Unlock()

			app.ELog.Println(e)
			conn.Conn.Close()
			return
		}
		appMsg.Body = apiMsg
		app.WSMessages <- appMsg
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
	if apiID == 1 {
		return WSBinanceMsg{}
	}
	return nil
}
