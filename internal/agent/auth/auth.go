package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
)

func GetAuthorization() error {
	const op = "auth.GetAuthorization"
	authCounter := 0

	tickerAuth := time.NewTicker(time.Millisecond * 1000)
	for range tickerAuth.C {
		if authCounter >= 10 {
			tickerAuth.Stop()
			return fmt.Errorf("%s: can not authorize", op)
		}

		if err := Login(); err != nil {
			log.Error().Str("op", op).Err(err).Msg("Login")
			storage.AppConfig.DisplayErrCh <- fmt.Sprintf("login error. attempt %d from 10", authCounter)
			authCounter++
			continue
		}
		break

	}
	tickerAuth.Stop()

	return nil
}

func Login() error {
	const op = "auth.Login"

	var tokenString storage.StTokenString

	buf := bytes.Buffer{}
	body := storage.StLoginBody{
		Login:    flags.Login,
		Password: flags.Password,
	}

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Login, &buf)
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%s: %w, StatusCode - %d", op, err, res.StatusCode)
		}
		bodyString := string(bodyBytes)

		log.Error().Str("op", op).
			Str("StatusCode", strconv.Itoa(res.StatusCode)).
			Str("body", bodyString).
			Msg("Do")
		return fmt.Errorf("%s: StatusCode - %d", op, res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)
	if err = dec.Decode(&tokenString); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Decode")
	}

	if err = storage.AppConfig.SetToken(tokenString.Token); err != nil {
		log.Error().Str("op", op).Err(err).Msg("SetToken")
	}

	return nil
}

func Refresh() error {
	const op = "auth.Refresh"

	var tokenString storage.StTokenString

	req, _ := http.NewRequest(http.MethodGet, storage.AppConfig.ApiRoutes.Refresh, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.GetToken()))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%s: %w, StatusCode - %d", op, err, res.StatusCode)
		}
		bodyString := string(bodyBytes)

		log.Error().Str("op", op).
			Str("StatusCode", strconv.Itoa(res.StatusCode)).
			Str("body", bodyString).
			Msg("Do")
		return fmt.Errorf("%s: StatusCode - %d", op, res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&tokenString); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Decode")
	}

	log.Info().Str("op", op).Str("storage.AppConfig.tokenString", storage.AppConfig.GetToken()).Msg("Old")

	if err = storage.AppConfig.SetToken(tokenString.Token); err != nil {
		log.Error().Str("op", op).Err(err).Msg("SetToken")
	}

	log.Info().Str("op", op).Str("tokenString.Token", tokenString.Token).Msg("New")
	log.Info().Str("op", op).Msg("ok")

	return nil
}
