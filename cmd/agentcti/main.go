package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-retry"
)

import (
	"github.com/pochtalexa/go-cti-middleware/internal/agent/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/logger"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/pgui"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
)

func logPanic(fileLogger *os.File) {
	const op = "logPanic"
	if p := recover(); p != nil {
		log.Error().Str("op", op).Msg(fmt.Sprintln(p))
	}
	fileLogger.Close()
}

func footerSetText() {
	for {
		select {
		case msg := <-storage.AppConfig.DisplayOkCh:
			pgui.FooterSetText(msg, tcell.ColorGreen)
		case msg := <-storage.AppConfig.DisplayErrCh:
			pgui.FooterSetText(msg, tcell.ColorRed)
		}
	}
}

// TODO: добавить логику авторизации и регистрации по флагам
func main() {
	fileLogger := logger.InitFileLogger()
	defer logPanic(fileLogger)

	var netErr net.Error
	b := retry.NewFibonacci(1 * time.Second)
	ctx := context.Background()

	flags.ParseFlags()

	storage.InitApiRoutes()
	storage.InitDisplayCh()

	go pgui.Init()

	go footerSetText()

	// пробуем авторизоваться
	tickerAuth := time.NewTicker(time.Millisecond * 1000)
	defer tickerAuth.Stop()
	authCounter := 0
	for range tickerAuth.C {
		if authCounter >= 10 {
			log.Fatal().Err(errors.New("can not authorise")).Msg("auth")
		}

		if err := auth.Login(); err != nil {
			log.Error().Err(err).Msg("Login")
			storage.AppConfig.DisplayErrCh <- fmt.Sprintf("login error. attempt %d from 10", authCounter)
			authCounter++
			continue
		}
		break
	}
	tickerAuth.Stop()

	// по таймеру запрашиваем новые метрики
	tickerMain := time.NewTicker(time.Millisecond * 5000)
	defer tickerMain.Stop()
	for range tickerMain.C {
		const op = "main loop"

		tempAgentEvents := storage.NewAgentEvents()

		req, _ := http.NewRequest(http.MethodGet, storage.AppConfig.ApiRoutes.Events, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.GetToken()))

		log.Debug().Str("storage.AppConfig.tokenString", storage.AppConfig.GetToken()).Msg("Events")

		err := retry.Do(ctx, retry.WithMaxRetries(3, b), func(ctx context.Context) error {
			res, err := storage.AppConfig.HTTPClient.Do(req)
			if err != nil {
				if errors.As(err, &netErr) ||
					netErr.Timeout() ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "connection reset by peer") {
					return retry.RetryableError(err)
				}
				return err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				// нет новых данных
				if res.StatusCode == http.StatusNoContent {
					return fmt.Errorf("%s: StatusNoContent", op)
				}
				if res.StatusCode == http.StatusUnauthorized {
					if err = auth.Refresh(); err != nil {
						return fmt.Errorf("%s:Refresh: %w", op, err)
					}
				}
			}

			dec := json.NewDecoder(res.Body)
			if err = dec.Decode(&tempAgentEvents); err != nil {
				return retry.RetryableError(fmt.Errorf("%w. jsonDecodeError", err))
			}

			return nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "StatusNoContent") {
				continue
			} else if strings.Contains(err.Error(), "jsonDecodeError") {
				log.Fatal().Str("op", op).Err(err).Msg("Decode")
			} else {
				log.Error().Str("op", op).Err(err).Msg("httpClient.Do")
				storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: connection error. err: %s", op, err)
				continue
			}
		}

		storage.AppConfig.DisplayOkCh <- fmt.Sprintf("%s: connected", op)

		// обновляем только те events которые получили от API
		storage.AgentEvents.UpdateEvents(tempAgentEvents)

		result, _ := storage.AgentEvents.ToString("UserState")
		pgui.UserState.SetText(result)

		result, _ = storage.AgentEvents.ToString("NewCall")
		pgui.NewCall.SetText(result)

		result, _ = storage.AgentEvents.ToString("CallStatus")
		pgui.CallStatus.SetText(result)
	}

}
