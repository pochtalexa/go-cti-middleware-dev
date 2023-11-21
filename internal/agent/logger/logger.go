package logger

import (
	"github.com/rs/zerolog/log"
	"os"
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

	//writers := io.MultiWriter(os.Stdout, fileLogger)
	//log.Logger = log.Output(writers)
	log.Logger = log.Output(fileLogger)

	log.Info().Msg("FileLogger initiated")

	return fileLogger
}
