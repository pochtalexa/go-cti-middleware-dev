package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"slices"
)

type IntAgent interface {
	SetEvent(event *storage.StWsEvent, message []byte) error
	DropAgentEvents(login string)
}

func SendCommand(c *websocket.Conn, body []byte) error {
	const op = "ws.SendCommand"

	err := c.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return fmt.Errorf("wsSendMessage: %w", err)
	}

	log.Info().Str("op", op).Str("message body", string(body)).Msg("WS SendCommand")

	return nil
}

// TODO ошибки отдавать в канал

// ReadMessage WS goroutine
func ReadMessage(c *websocket.Conn, agentsInfo IntAgent) {
	const op = "ws.ReadMessage"

	for {
		wsEvent := storage.NewWsEvent()
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Str("op", op).Err(err).Msg("ReadMessage")
		}

		if err = json.Unmarshal(message, &wsEvent); err != nil {
			log.Error().Str("op", op).Err(err).Msg("Unmarshal wsEvent")
		}

		if err = json.Unmarshal(message, &wsEvent.Body); err != nil {
			log.Error().Err(err).Str("op", op).Msg("Unmarshal wsEvent.Body")
		}

		if err = agentsInfo.SetEvent(wsEvent, message); err != nil {
			log.Error().Err(err).Str("op", op).Msg("SetEvent")
		}

		if slices.Contains(wsEvent.ErrorNames, wsEvent.Name) {
			log.Error().Str("message", fmt.Sprintln(wsEvent.Body)).Msg("wsReadMessage")
		} else {
			log.Info().Str("event name", wsEvent.Name).Str("message body", fmt.Sprintln(wsEvent.Body)).Msg("wsReadMessage")
		}
		log.Debug().Str("op", op).Str("message", wsEvent.Name).Msg("name")
		log.Debug().Str("op", op).Str("message", wsEvent.Login).Msg("login")
	}
}
