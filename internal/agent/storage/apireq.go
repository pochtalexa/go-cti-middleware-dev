package storage

import (
	"fmt"

	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
)

type StApiRoutes struct {
	Register string
	Login    string
	Refresh  string
	Events   string
	Control  string
}

func NewApiRoutes() *StApiRoutes {
	urlRegister := fmt.Sprintf("%s/api/v1/register", flags.ServAddr)
	urlLogin := fmt.Sprintf("%s/api/v1/login", flags.ServAddr)
	urlRefresh := fmt.Sprintf("%s/api/v1/refresh", flags.ServAddr)
	urlEvents := fmt.Sprintf("%s/api/v1/events/%s", flags.ServAddr, flags.Login)
	urlControl := fmt.Sprintf("%s/api/v1/control", flags.ServAddr)

	return &StApiRoutes{
		Register: urlRegister,
		Login:    urlLogin,
		Refresh:  urlRefresh,
		Events:   urlEvents,
		Control:  urlControl,
	}
}

func InitApiRoutes() {
	AppConfig.ApiRoutes = NewApiRoutes()
}
