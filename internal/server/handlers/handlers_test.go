package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

import (
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws/mocks"
)

func TestControlHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ControlHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestEventsHandler(t *testing.T) {
	tests := []struct {
		name   string
		login  string
		events storage.StAgentEvents
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAgentsInfo := mocks.NewIntAgent(t)
			AgentsInfo = testAgentsInfo

			// TODO: описать вызовы AgentsInfo

			mux := chi.NewRouter()
			mux.Use(middleware.Logger)
			mux.Post("/api/v1/events/{login}", EventsHandler)

			url := fmt.Sprintf("/api/v1/events/{%s}", tt.login)

			reqBody, _ := json.Marshal(tt.events)

			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()

		})
	}
}

func TestLoginHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoginHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RefreshHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestRegisterUserHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterUserHandler(tt.args.w, tt.args.r)
		})
	}
}
