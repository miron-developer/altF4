package app

import "net/http"

func (app *Application) SetRoutes() http.Handler {
	appMux := http.NewServeMux()
	appMux.HandleFunc("/", app.HIndex)

	// ws
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/", app.CreateWSUser)
	appMux.Handle("/ws/", http.StripPrefix("/ws", wsMux))

	// api routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", app.HApiIndex)
	apiMux.HandleFunc("/exchange-points", app.HExPoints)
	apiMux.HandleFunc("/exchange-currencies", app.HExCurrencies)
	apiMux.HandleFunc("/exchange-subscribe", app.HExSubscribe)
	apiMux.HandleFunc("/exchange-unsubscribe", app.HExUnsubscribe)
	appMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// assets get
	assets := http.FileServer(http.Dir("assets"))
	appMux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// middlewares
	muxHanlder := app.AccessLogMiddleware(appMux)
	muxHanlder = app.SecureHeaderMiddleware(muxHanlder)
	return muxHanlder
}
