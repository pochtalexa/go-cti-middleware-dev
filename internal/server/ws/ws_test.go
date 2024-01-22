package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

import (
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws/mocks"
)

var upgrader = websocket.Upgrader{}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("upgrade")
		return
	}
	defer c.Close()

	//for {
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("ReadMessage")
		//break
	}
	err = c.WriteMessage(mt, message)
	if err != nil {
		log.Error().Err(err).Msg("WriteMessage")
		//break
		//}
	}
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg map[string]interface{}) {
	t.Helper()

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err = ws.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestReadMessage(t *testing.T) {

	tests := []struct {
		name  string
		event map[string]interface{}
	}{
		{
			name: "event ok",
			event: map[string]interface{}{
				"name":  "OnClose",
				"login": "agent",
				"cid":   "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsEvent := storage.NewWsEvent()

			jsonData, _ := json.Marshal(tt.event)
			_ = json.Unmarshal(jsonData, &wsEvent)
			_ = json.Unmarshal(jsonData, &wsEvent.Body)

			s := httptest.NewServer(http.HandlerFunc(serveHTTP))
			defer s.Close()

			u := "ws" + strings.TrimPrefix(s.URL, "http")

			// Connect to the server
			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer ws.Close()

			config.Init()
			config.ServerConfig.WsConn = ws

			sendMessage(t, ws, tt.event)

			agentsInfo := mocks.NewIntAgent(t)

			agentsInfo.
				On("SetEvent", wsEvent, mock.Anything).
				Once().
				Return(nil)

			ReadMessage(agentsInfo)
		})
	}
}
