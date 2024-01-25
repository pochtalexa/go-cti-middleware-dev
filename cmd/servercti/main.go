package main

import (
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/server/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/server/logger"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pochtalexa/go-cti-middleware/internal/server/api"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/handlers"
	"github.com/pochtalexa/go-cti-middleware/internal/server/migrations"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
)

func logPanic(multiLogger *os.File) {
	const op = "logPanic"

	if p := recover(); p != nil {
		log.Error().Str("op", op).Msg(fmt.Sprintln(p))
	}
	multiLogger.Close()
}

func main() {
	// TODO тесты
	// TODO Обработка ошибок
	// TODO обработка ответа CTI на отправленные команды
	// TODO на перспективу использовать Redis

	// TODO: добавить флаг - была ли подписка на агента по логину - актуально, когда работаем без авторизации
	// и сервер перезагрузился

	const op = "main"

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	flags.ParseFlags()

	config.Init()
	if err := config.ServerConfig.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}
	if config.ServerConfig.Settings.LogLevel != "info" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	fileLogger := logger.InitFileLogger()
	defer logPanic(fileLogger)

	useAuth := config.ServerConfig.Settings.UseAuth
	if useAuth {
		st, err := storage.InitConnDB()
		if err != nil {
			log.Fatal().Str("op", op).Err(err).Msg("InitConnDB")
		}
		defer st.DB.Close()

		err = migrations.ApplyMigrations()
		if err != nil {
			log.Fatal().Str("op", op).Err(err).Msg("ApplyMigrations")
		}
		log.Debug().Str("op", op).Bool("useAuth", useAuth).Msg("DB init - ok")
	} else {
		log.Debug().Str("op", op).Bool("useAuth", useAuth).Msg("no DB init needed")
	}

	handlers.Init()

	err := cti.Init()
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
