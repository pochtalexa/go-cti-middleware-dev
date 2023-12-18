package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pochtalexa/go-cti-middleware/internal/server/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ControlHandler принимем команду по http API и вызваем соотвествующий медот CTI API
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ControlHandler"

	reqBody := storage.NewWsCommand()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Str("op", op).Err(err).Msg("Decode")
		return
	}

	if reqBody.Name == "" {
		errorText := fmt.Errorf("no key 'name' in request")
		log.Error().Err(errorText).Msg("ControlHandler")
		http.Error(w, errorText.Error(), http.StatusBadRequest)
		return
	}

	switch reqBody.Name {
	case "ChangeUserState":
		if err := cti.ChangeStatus(reqBody.Rid, reqBody.Login, reqBody.State); err != nil {
			log.Error().Err(err).Msg("call ChangeStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Answer":
		if err := cti.Answer(reqBody.Rid, reqBody.Login, reqBody.Cid); err != nil {
			log.Error().Err(err).Msg("call Answer")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Hangup":
		if err := cti.Hangup(reqBody.Rid, reqBody.Login, reqBody.Cid); err != nil {
			log.Error().Err(err).Msg("call Hangup")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Close":
		if err := cti.Close(reqBody.Rid, reqBody.Login, reqBody.Cid); err != nil {
			log.Error().Err(err).Msg("call Close")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Mute":
		if err := cti.Mute(reqBody.Rid, reqBody.Login, reqBody.Cid, reqBody.On); err != nil {
			log.Error().Err(err).Msg("call Mute")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Info().Str("reqBody", fmt.Sprint(reqBody)).Msg(op)
}

// EventsHandler запрос на получение текущих events
func EventsHandler(w http.ResponseWriter, r *http.Request) {
	// проверяем что логин в запросе не пустой
	login := chi.URLParam(r, "login")
	if login == "" {
		errorText := fmt.Errorf("no login in request")
		http.Error(w, errorText.Error(), http.StatusBadRequest)
		log.Error().Err(errorText).Msg("parse url")
		return
	}

	// проверяем что есть обновленные данные
	updated, ok := storage.AgentsInfo.Updated[login]
	if !ok {
		errorText := fmt.Errorf("no key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("EventsHandler")
		return
	}
	if !updated {
		errorText := fmt.Errorf("no updated data for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNoContent)
		log.Debug().Str("updated", errorText.Error()).Msg("StatusNoContent")
		return
	}

	mutex, ok := storage.AgentsInfo.Mutex[login]
	if !ok {
		errorText := fmt.Errorf("no mutex key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("EventsHandler")
		return
	}
	mutex.RLock()
	resBody, ok := storage.AgentsInfo.Events[login]
	if !ok {
		errorText := fmt.Errorf("no events key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("resBody")
		return
	}
	mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("Encode")
		return
	}

	storage.AgentsInfo.Updated[login] = false

	// очищаем хранилище после отправки
	storage.AgentsInfo.DropAgentEvents(login)

	log.Info().Msg("GetEventsHandler - ok")

	return
}

// RegisterUserHandler регистрация новго агента CTI
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.RegisterUserHandler"

	reqBody := storage.NewCredentials()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	id, err := auth.RegisterNewUser(reqBody.Login, reqBody.Password, storage.Storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg(op)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resBody := storage.NewRegisterSuccess()
	resBody.Id = id

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	log.Debug().Str("id", strconv.FormatInt(id, 10)).Msg(op)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.LoginHandler"

	reqBody := storage.NewCredentials()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	token, err := auth.Login(reqBody.Login, reqBody.Password, storage.Storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg(op)
		return
	}

	if err := cti.AttachUser(reqBody.Login); err != nil {
		log.Fatal().Str("op", op).Err(err).Msg("ws AttachUser")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resBody := storage.NewTokenBody()
	resBody.Token = token

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	log.Debug().Str("login", reqBody.Login).Msg(fmt.Sprintf("login success: %s", op))

	return
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.RefreshHandler"

	if !config.ServerConfig.Settings.UseAuth {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "no refresh needed", http.StatusBadRequest)
		return
	}

	tokenField := r.Header.Get("Authorization")
	tokenSlice := strings.Split(tokenField, " ")
	if tokenField == "" || tokenSlice[0] != "Bearer" || len(tokenSlice) != 2 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, errors.New("no token provided").Error(), http.StatusBadRequest)
		return
	}

	tokenString := tokenSlice[1]
	claims := storage.NewClaims()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, errors.New("unauthorized").Error(), http.StatusUnauthorized)
			return "", errors.New("unauthorized")
		}

		return []byte(config.ServerConfig.Secret), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Errorf("%s: %w", op, err).Error(), http.StatusUnauthorized)
		return
	}

	//if !token.Valid - не проверям т.к. если просрочен - становится не валидным

	// не обновляем ранее указанного времени
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		http.Error(w, fmt.Errorf("%s: not expired", op).Error(), http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(config.ServerConfig.TokenTTL)
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(config.ServerConfig.Secret))
	if err != nil {
		http.Error(w, fmt.Errorf("%s: refresh error", op).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resBody := storage.NewTokenBody()
	resBody.Token = tokenString

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	log.Debug().Str("id", strconv.FormatInt(claims.ID, 10)).Msg(op)
	return
}
