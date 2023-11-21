package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pochtalexa/go-cti-middleware/internal/server/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// ControlHandler принимем команду по http API и вызваем соотвествующий медот CTI API
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := storage.NewWsCommand()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Decode")
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
		if err := cti.ChageStatus(cti.Conn, "agent", reqBody.State); err != nil {
			log.Error().Err(err).Msg("call ChageStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Answer":
		if err := cti.Answer(cti.Conn, "agent", reqBody.Cid); err != nil {
			log.Error().Err(err).Msg("call ChageStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Hangup":
		if err := cti.Hangup(cti.Conn, "agent", reqBody.Cid); err != nil {
			log.Error().Err(err).Msg("call ChageStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "Mute":
		if err := cti.Mute(cti.Conn, "agent", reqBody.Cid, reqBody.On); err != nil {
			log.Error().Err(err).Msg("call ChageStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Info().Str("reqBody", fmt.Sprint(reqBody)).Msg("reqBody")
}

// EventsHandler запрос на получение текущих events
func EventsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO добавить авторизацию

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
		log.Error().Err(errorText).Msg("")
		return
	}
	if !updated {
		errorText := fmt.Errorf("no updated data for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNoContent)
		log.Debug().Err(errorText).Msg("StatusNoContent")
		return
	}

	resBody, ok := storage.AgentsInfo.Events[login]
	if !ok {
		errorText := fmt.Errorf("no key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("")
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusOK)
	log.Info().Msg("GetEventsHandler - ok")

	return
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.RegisterUserHandler"

	reqBody := storage.NewRegister()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	id, err := auth.RegisterNewUser(reqBody.Login, reqBody.Password)
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
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg(op)
		return
	}

	log.Debug().Str("id", strconv.FormatInt(id, 10)).Msg(op)
	return
}
