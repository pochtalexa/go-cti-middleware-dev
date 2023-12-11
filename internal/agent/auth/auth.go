package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

func Login() error {
	type stTokenString struct {
		Token string `json:"token"`
	}

	const op = "auth.Login"
	var tokenString stTokenString

	buf := bytes.Buffer{}
	body := storage.AppConfig.Credentials

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
