package flags

import (
	"flag"
	"github.com/rs/zerolog/log"
)

var (
	ServAddr   string
	Login      string
	Password   string
	UrlControl string
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ParseFlags() {

	defaultLogin := "agent"
	defaultPassword := "123"

	flag.StringVar(&ServAddr, "a", "http://localhost:9595", "middleware api addr")
	flag.StringVar(&Login, "l", defaultLogin, "login")
	flag.StringVar(&Password, "p", defaultPassword, "password")
	flag.Parse()

	UrlControl = ServAddr + "/api/v1/control"

	log.Info().Msg("ParseFlags - ok")
}
