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
		accessOrigin := "http://localhost:3000"
		if app.IsHeroku {
			accessOrigin = "https://wnet-sn.herokuapp.com"
		}
		w.Header().Set("Access-Control-Allow-Origin", accessOrigin)
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

	user := &WSUser{Conn: conn, ID: StringWithCharset(8)}
	app.m.Lock()
	app.OnlineUsers[user.ID] = user
	app.m.Unlock()

	go app.OnlineUsers[user.ID].HandleUserMsg(app)
	go app.OnlineUsers[user.ID].Pinger()
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
