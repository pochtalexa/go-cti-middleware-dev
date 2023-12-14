package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

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
	if err := dec.Decode(&tokenString); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Decode")
	}

	storage.AppConfig.TokenString = tokenString.Token

	return nil
}

func Refresh() error {

	const op = "auth.Refresh"
	var tokenString storage.StTokenString

	req, _ := http.NewRequest(http.MethodGet, storage.AppConfig.ApiRoutes.Refresh, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
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

	storage.AppConfig.TokenString = tokenString.Token

	log.Info().Str("storage.AppConfig.TokenString", storage.AppConfig.TokenString).Msg("Old")
	log.Info().Str("tokenString.Token", tokenString.Token).Msg("New")

	log.Info().Str("op", op).Msg("ok")

	return nil
}
