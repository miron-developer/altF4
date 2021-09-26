package app

import (
	"errors"
	"net/http"
	"os"
	"text/template"
)

// SecureHeaderMiddleware set secure header option
func (app *Application) SecureHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cross-origin-resource-policy", "cross-origin")
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// AccessLogMiddleware logging request
func (app *Application) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.CurrentRequestCount < app.MaxRequestCount {
			app.CurrentRequestCount++
			app.ILog.Printf(logingReq(r))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "service is overloaded", 529)
			app.ELog.Println(errors.New("rate < curl"))
		}
	})
}

/* ----------------------------------------------- Websocket ---------------------------------------------- */

// CreateWSUser create one WSUser
func (app *Application) CreateWSUser(w http.ResponseWriter, r *http.Request) {
	conn, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		app.ELog.Println(e, r.RemoteAddr)
		return
	}

	// generating unical id
	id := -1
	for {
		id = RandomInt()
		if exist := app.findByID(id, 1); exist == nil {
			break
		}
	}

	user := &WSUser{Conn: conn, ID: id, Token: StringWithCharset(8)}
	app.m.Lock()
	app.OnlineUsers[user.Token] = user
	app.m.Unlock()

	// send user token
	user.Conn.WriteJSON(&WSMessage{MsgType: WSM_SEND_USERID, AddresserToken: WSM_SERVER_TOKEN, ReceiverToken: user.Token, Body: user.Token})

	go app.OnlineUsers[user.Token].HandleUserMsg(app)
	go app.OnlineUsers[user.Token].Pinger()
}

// HIndex handle all GETs
func (app *Application) HIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		return
	}

	wd, _ := os.Getwd()
	t, e := template.ParseFiles(wd + "/assets/index.html")
	if e != nil {
		http.Error(w, "can't load this page", 500)
		app.ELog.Println(e)
		return
	}
	t.Execute(w, nil)
}

/* ----------------------------------------------- API ---------------------------------------------- */

// HApiIndex for handle '/api/'
func (app *Application) HApiIndex(w http.ResponseWriter, r *http.Request) {
	type Route struct {
		Path        string            `json:"route"`
		Description string            `json:"description"`
		Params      map[string]string `json:"params"`
		Children    []Route           `json:"children"`
	}

	data := API_RESPONSE{
		Code: 200,
		Data: Route{
			Path:        "/",
			Description: "Api possible routes",
			Children: []Route{
				{Path: "/exchange-points", Description: "get all possible exchange points"},
				{Path: "/exchange-currencies", Description: "get all exchange currencies from selected exchange point", Params: map[string]string{"point": "from which point do you want to see currencies"}},
			},
		},
	}

	DoJS(w, data)
}

// HExPoints for handle '/api/exchange-points'
func (app *Application) HExPoints(w http.ResponseWriter, r *http.Request) {
	HApi(w, r, func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if r.Method == "POST" {
			return nil, errors.New("wrong method")
		}
		return []map[string]string{{"id": "100", "name": "Binance"}}, nil
	})
}

// HExCurrencies for handle '/api/exchange-currencies'
func (app *Application) HExCurrencies(w http.ResponseWriter, r *http.Request) {
	HApi(w, r, app.GetCurrenciesFromID)
}

// HExSubscribe for handle '/api/exchange-subscribe'
func (app *Application) HExSubscribe(w http.ResponseWriter, r *http.Request) {
	HApi(w, r, app.SubscribeCurrencie)
}

// HExUnsubscribe for handle '/api/exchange-unsubscribe'
func (app *Application) HExUnsubscribe(w http.ResponseWriter, r *http.Request) {
	HApi(w, r, app.UnsubscribeCurrencie)
}
