package app

import "net/http"

func (app *Application) SetRoutes() http.Handler {
	appMux := http.NewServeMux()
	appMux.HandleFunc("/", app.HIndex)

	// ws
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/", app.CreateWSUser)
	appMux.Handle("/ws/", http.StripPrefix("/ws", wsMux))

	// middlewares
	muxHanlder := app.AccessLogMiddleware(appMux)
	muxHanlder = app.SecureHeaderMiddleware(muxHanlder)
	return muxHanlder
}
