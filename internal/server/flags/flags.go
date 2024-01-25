package flags

import (
	"flag"

	"github.com/rs/zerolog/log"
)

var (
	CfgPath string
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

	defaultCfgPath := "config.toml"

	flag.StringVar(&CfgPath, "cfg", defaultCfgPath, "toml config path")
	flag.Parse()

	log.Info().Msg("ParseFlags - ok")
}
