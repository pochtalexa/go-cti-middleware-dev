package main

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/httpconf"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/logger"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/pgui"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	fileLogger := logger.InitFileLogger()
	defer fileLogger.Close()

	httpconf.Init()
	flags.ParseFlags()

	go pgui.Init()

	url := flags.ServAddr + "/api/v1/events/agent"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// по таймеру запрашиваем новые метрики
	for range time.Tick(time.Second * 1) {
		tempAgentEvents := storage.NewAgentEvents()

		res, err := httpconf.HTTPClient.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("httpClient.Do")
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
			log.Fatal().Err(err).Msg("Decode")
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
