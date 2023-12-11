package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/logger"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/pgui"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// TODO: добавить логику авторизации
func main() {
	fileLogger := logger.InitFileLogger()
	defer fileLogger.Close()

	flags.ParseFlags()
	storage.InitApiRoutes()

	if err := auth.Login(); err != nil {
		log.Fatal().Err(err).Msg("Login")
	}

	go pgui.Init()

	// по таймеру запрашиваем новые метрики
	for range time.Tick(time.Second * 1) {
		const op = "main loop"

		tempAgentEvents := storage.NewAgentEvents()

		req, _ := http.NewRequest(http.MethodGet, storage.AppConfig.ApiRoutes.Events, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
		res, err := storage.AppConfig.HTTPClient.Do(req)
		if err != nil {
			log.Error().Str("op", op).Err(err).Msg("httpClient.Do")
			pgui.FooterSetText("Connection error", tcell.ColorRed)
			continue
		}
		defer res.Body.Close()
		pgui.FooterSetText("Connected", tcell.ColorGreen)

		// нет новых данных
		if res.StatusCode == http.StatusNoContent {
			continue
		}

		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&tempAgentEvents); err != nil {
			log.Fatal().Str("op", op).Err(err).Msg("Decode")
		}

		// обновляем только те events которые получили от API
		storage.AgentEvents.UpdateEvents(tempAgentEvents)

		result, _ := storage.AgentEvents.ToString("UserState")
		pgui.UserState.SetText(result)

		result, _ = storage.AgentEvents.ToString("NewCall")
		pgui.NewCall.SetText(result)

		result, _ = storage.AgentEvents.ToString("CallStatus")
		pgui.CallStatus.SetText(result)

		//log.Info().Str("resp", fmt.Sprintln(resp)).Msg("")
		//log.Info().Str("resp[state]", fmt.Sprintln(resp["state"])).Msg("")

	}

}
