package flags

import (
	"flag"
	"github.com/rs/zerolog/log"
)

var (
	ServAddr  string
	AgentName string
	Login     string
	Password  string
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

	defaultServAddr := "http://localhost:9595"
	defaultAgentName := "agent"
	defaultLogin := "test1"
	defaultPassword := "test123"

	flag.StringVar(&ServAddr, "a", defaultServAddr, "middleware api addr")
	flag.StringVar(&AgentName, "n", defaultAgentName, "agent name")
	flag.StringVar(&Login, "l", defaultLogin, "login")
	flag.StringVar(&Password, "p", defaultPassword, "password")
	flag.Parse()

	log.Info().Msg("ParseFlags - ok")
}
