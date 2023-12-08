package api

import (
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/handlers"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

// проверяем, что клиент отправил серверу сжатые данные в формате gzip
func checkGzipEncoding(r *http.Request) bool {

	encodingSlice := r.Header.Values("Content-Encoding")
	encodingsStr := strings.Join(encodingSlice, ",")
	encodings := strings.Split(encodingsStr, ",")

	log.Debug().Str("encodingsStr", encodingsStr).Msg("checkGzipEncoding")

	for _, el := range encodings {
		if el == "gzip" {
			return true
		}
	}

	return false
}

func GzipDecompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if checkGzipEncoding(r) {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gzipReader
			defer gzipReader.Close()
		}

		log.Debug().Msg("GzipDecompression passed")

		next.ServeHTTP(w, r)
	})
}

func checkCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.checkCredentials"

		if !config.ServerConfig.Settings.UseAuth {
			next.ServeHTTP(w, r)
		}

		tokenField := r.Header.Get("Authorization")
		tokenSlice := strings.Split(tokenField, " ")
		if tokenField == "" || tokenSlice[0] != "Bearer" || len(tokenSlice) != 2 {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, errors.New("no token provided").Error(), http.StatusBadRequest)
			return
		}

		tokenString := tokenSlice[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, errors.New("unauthorized").Error(), http.StatusUnauthorized)
				return "", errors.New("unauthorized")
			}

			return config.ServerConfig.Secret, nil
		})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Errorf("%s: %w", op, err).Error(), http.StatusUnauthorized)
			return
		}

		var agent storage.StAgent
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			agent.ID = claims["uid"].(int64)
			agent.Login = claims["login"].(string)
			agent.TokenTTL = claims["exp"].(int64)
		} else {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Errorf("%s: invalid token", op).Error(), http.StatusUnauthorized)
			return
		}

		if time.Now().Unix() <= agent.TokenTTL {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Errorf("%s: expired", op).Error(), http.StatusUnauthorized)
			return
		}

		log.Debug().Msg(fmt.Sprintf("%s - ok", op))

		next.ServeHTTP(w, r)
	})
}

func RunAPI(urlStr string) error {
	//logger := httplog.NewLogger("httplog", httplog.Options{
	//	LogLevel: slog.LevelDebug,
	//	//JSON:             true,
	//	Concise:          false,
	//	RequestHeaders:   true,
	//	ResponseHeaders:  true,
	//	MessageFieldName: "msg",
	//	//LevelFieldName:   "severity",
	//	TimeFieldFormat: time.RFC3339,
	//	Tags: map[string]string{
	//		"version": "v1.0",
	//		"env":     "dev",
	//	},

	//QuietDownRoutes: []string{
	//	"/",
	//	"/ping",
	//},
	//QuietDownPeriod: 10 * time.Second,
	//SourceFieldName: "source",
	//})

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(GzipDecompression)
	mux.Use(middleware.Compress(flate.DefaultCompression, "application/json", "text/html"))
	//mux.Use(middleware.Logger)
	//mux.Use(httplog.RequestLogger(logger))

	mux.Post("/api/v1/register", handlers.RegisterUserHandler)
	mux.Post("/api/v1/login", handlers.LoginHandler)

	mux.Route("/api/v1/control", func(r chi.Router) {
		r.Use(checkCredentials)
		r.Post("/", handlers.ControlHandler)
	})

	mux.Route("/api/v1/events/{login}", func(r chi.Router) {
		r.Use(checkCredentials)
		r.Get("/", handlers.EventsHandler)
	})

	log.Info().Str("Running on", urlStr).Msg("httpconf server started")

	return http.ListenAndServe(urlStr, mux)
}
