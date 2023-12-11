package storage

import (
	"net/http"
)

var AppConfig = &StConfig{}

type StConfig struct {
	TokenString string
	HTTPClient  http.Client
	ApiRoutes   *StApiRoutes
	Credentials StLoginBody
}
