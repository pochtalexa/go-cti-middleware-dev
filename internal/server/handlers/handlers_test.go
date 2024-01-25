package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/pochtalexa/go-cti-middleware/internal/server/auth"
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

func getRandRid() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// generate a random 8-digit ID
	id := ""
	for i := 0; i < 8; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}

	return id
}

// TestLoginHandler - тест не требуется, т.к. проверяем в auth_test.go
// TestRegisterUserHandler - тест не требуется, т.к. проверяем в auth_test.go

func TestControlHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name   string
		login  string
		status string
		want   want
	}{
		{
			name:   "ChangeUserState",
			login:  "agent",
			status: "normal",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(serveHTTP))
			defer ts.Close()

			tu := "ws" + strings.TrimPrefix(ts.URL, "http")

			// Connect to the server
			tws, _, err := websocket.DefaultDialer.Dial(tu, nil)
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer tws.Close()

			config.Init()
			config.ServerConfig.WsConn = tws

			mux := chi.NewRouter()
			mux.Use(middleware.Logger)
			mux.Post("/api/v1/control", ControlHandler)

			url := fmt.Sprintf("/api/v1/control")

			buf := bytes.Buffer{}

			body := storage.NewWsCommand()
			body.Name = tt.name
			body.Rid = getRandRid()
			body.Login = tt.login
			body.State = tt.status

			enc := json.NewEncoder(&buf)
			enc.SetIndent("", "  ")
			if err = enc.Encode(body); err != nil {
				t.Fatalf("%v", err)
			}

			request := httptest.NewRequest(http.MethodPost, url, &buf)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

		})
	}
}

func TestEventsHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name   string
		login  string
		events storage.StAgentEvents
		want   want
	}{
		{
			name:  "OnClose",
			login: "agent",
			events: storage.StAgentEvents{
				OnClose: storage.StOnClose{
					Name:  "OnClose",
					Login: "agent",
					Cid:   123,
				},
			},
			want: want{statusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mutex sync.RWMutex
			testAgentsInfo := mocks.NewIntAgent(t)
			AgentsInfo = testAgentsInfo

			testAgentsInfo.
				On("IsUpdated", tt.login).
				Once().
				Return(true, true)

			testAgentsInfo.
				On("GetMutex", tt.login).
				Once().
				Return(&mutex, true)

			testAgentsInfo.
				On("GetEvents", tt.login).
				Once().
				Return(tt.events, true)

			testAgentsInfo.
				On("SetUpdated", tt.login, false).
				Once().
				Return(nil)

			testAgentsInfo.
				On("DropAgentEvents", tt.login).
				Once().
				Return(nil)

			mux := chi.NewRouter()
			mux.Use(middleware.Logger)
			mux.Get("/api/v1/events/{login}", EventsHandler)

			url := fmt.Sprintf("/api/v1/events/%s", tt.login)

			request := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name     string
		useAuth  bool
		login    string
		id       int64
		tokenTTL int64
		time     time.Duration
		want     want
	}{
		{
			name:     "Refresh ok",
			useAuth:  true,
			login:    "agent",
			id:       123,
			tokenTTL: 25,
			time:     time.Second,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:     "Refresh err#1",
			useAuth:  true,
			login:    "agent",
			id:       123,
			tokenTTL: 60,
			time:     time.Second,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Init()
			config.ServerConfig.Settings.UseAuth = tt.useAuth
			config.ServerConfig.TokenTTL = time.Duration(tt.tokenTTL) * tt.time

			agent := storage.NewAgent()
			agent.Login = tt.login
			agent.ID = tt.id

			token, err := auth.NewToken(agent, config.ServerConfig.TokenTTL)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to generate token")
			}

			mux := chi.NewRouter()
			mux.Use(middleware.Logger)
			mux.Get("/api/v1/refresh", RefreshHandler)

			url := fmt.Sprintf("/api/v1/refresh")

			request := httptest.NewRequest(http.MethodGet, url, nil)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
