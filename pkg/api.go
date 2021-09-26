package app

import (
	"encoding/json"
	"net/http"
)

type API_RESPONSE struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// DoJS do json and write it
func DoJS(w http.ResponseWriter, data interface{}) {
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "Application/json")
	w.Write(js)
}

// SendErrorJSON send to front error
func SendErrorJSON(w http.ResponseWriter, data API_RESPONSE, err string) {
	data.Data = err
	data.Code = 401
	DoJS(w, data)
}

// HApi general handler from api
func HApi(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request) (interface{}, error)) {
	data := API_RESPONSE{
		Data: "",
		Code: 200,
	}

	datas, e := f(w, r)
	if e != nil {
		SendErrorJSON(w, data, e.Error())
		return
	}
	data.Data = datas
	DoJS(w, data)
}
