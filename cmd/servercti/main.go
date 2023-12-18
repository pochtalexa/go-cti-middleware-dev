package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pochtalexa/go-cti-middleware/internal/server/api"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/handlers"
	"github.com/pochtalexa/go-cti-middleware/internal/server/migrations"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO тесты
	// TODO Обработка ошибок
	// TODO обработка ответа CTI на отправленные команды
	// TODO на перспективу использовать Redis

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	config.Init()
	if err := config.ServerConfig.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}
	if config.ServerConfig.Settings.LogLevel != "info" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	st, err := storage.InitConnDB()
	if err != nil {
		log.Fatal().Err(err).Msg("ApplyMigrations")
	}
	defer st.DB.Close()

	err = migrations.ApplyMigrations()
	if err != nil {
		log.Fatal().Err(err).Msg("ApplyMigrations")
	}

	handlers.Init()

	err = cti.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("cti.Init")
	}
	defer config.ServerConfig.WsConn.Close()

	go ws.ReadMessage(storage.AgentsInfo)

	if err := cti.InitCTISess(); err != nil {
		log.Fatal().Err(err).Msg("InitCTISess")
	}

	uHTTP := config.ServerConfig.HttpAPI.Host + ":" + config.ServerConfig.HttpAPI.Port
	if err := api.RunAPI(uHTTP); err != nil {
		log.Fatal().Err(err).Msg("RunAPI")
	}
}
