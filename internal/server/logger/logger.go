package logger

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
)

func InitFileLogger() *os.File {
	file, err := os.OpenFile(
		config.ServerConfig.Settings.LogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("initFileLogger")
	}

	log.Logger = log.Output(file)

	log.Info().Msg("FileLogger initiated")

	return file
}
