package logger

import (
	"os"

	"github.com/rs/zerolog/log"
)

func InitFileLogger() *os.File {
	fileLogger, err := os.OpenFile(
		"clientcti.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("initMultiLogger")
	}

	log.Logger = log.Output(fileLogger)

	log.Info().Msg("FileLogger initiated")

	return fileLogger
}
