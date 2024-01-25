package flags

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	ServScheme string
	ServAddr   string
	Login      string
	Password   string
	Register   bool
	Polling    int
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

	defaultScheme := "http"
	defaultServAddr := "localhost:9595"
	defaultLogin := "agent"
	defaultPassword := "123"
	defaultPolling := 1

	flag.StringVar(&ServScheme, "s", defaultScheme, "middleware api scheme")
	flag.StringVar(&ServAddr, "a", defaultServAddr, "middleware api addr")
	flag.StringVar(&Login, "l", defaultLogin, "login")
	flag.StringVar(&Password, "p", defaultPassword, "password")
	flag.IntVar(&Polling, "plg", defaultPolling, "polling time")
	flag.BoolVar(&Register, "r", false, "use registration")
	flag.Parse()

	ServAddr = fmt.Sprintf("%s://%s", ServScheme, ServAddr)

	log.Info().Msg("ParseFlags - ok")
}
