/*
	Initialize app
*/

package app

import (
	"log"
	"os"
	"sync"
	"time"
)

// Application this is app struct and items
type Application struct {
	m                   sync.Mutex
	ELog                *log.Logger
	ILog                *log.Logger
	Port                string
	CurrentRequestCount int
	MaxRequestCount     int
	IsHeroku            bool
	OnlineUsers         map[string]*WSUser
	WSSConns            map[string]*WSConnection
	WSMessages          chan *WSMessage
}

// InitProg initialise
func InitProg() *Application {
	wd, _ := os.Getwd()
	logFile, _ := os.OpenFile(wd+"/logs/log_"+time.Now().Format("2006-01-02")+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	eLog := log.New(logFile, "\033[31m[ERROR]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog := log.New(logFile, "\033[34m[INFO]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog.Println("loggers is done!")

	return &Application{
		ELog:                eLog,
		ILog:                iLog,
		Port:                "4330",
		CurrentRequestCount: 0,
		MaxRequestCount:     1200,
		IsHeroku:            false,
		OnlineUsers:         map[string]*WSUser{},
		WSSConns:            map[string]*WSConnection{},
		WSMessages:          make(chan *WSMessage),
	}
}
