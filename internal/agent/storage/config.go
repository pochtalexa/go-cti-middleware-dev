package storage

import (
	"net/http"
	"sync"
)

var AppConfig = &StConfig{}

type StConfig struct {
	tokenString  string
	HTTPClient   http.Client
	ApiRoutes    *StApiRoutes
	Mutex        sync.RWMutex
	DisplayOkCh  chan string
	DisplayErrCh chan string
}

func (ths *StConfig) SetToken(tokenString string) error {
	ths.Mutex.Lock()
	defer ths.Mutex.Unlock()

	ths.tokenString = tokenString

	return nil
}

func (ths *StConfig) GetToken() string {
	ths.Mutex.RLock()
	defer ths.Mutex.RUnlock()

	return ths.tokenString
}

func InitDisplayCh() {
	AppConfig.DisplayOkCh = NewDisplayCh()
	AppConfig.DisplayErrCh = NewDisplayCh()
}

func NewDisplayCh() chan string {
	return make(chan string, 10)
}
