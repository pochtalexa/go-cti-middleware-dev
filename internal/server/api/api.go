package api

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"

	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/handlers"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
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

		useAuth := config.ServerConfig.Settings.UseAuth
		if useAuth {

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
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Header().Set("Content-Type", "application/json")

					var expiredBody = make(map[string]string)
					expiredBody["name"] = "token"
					expiredBody["data"] = "expired"

					enc := json.NewEncoder(w)
					enc.SetIndent("", "  ")
					if err := enc.Encode(expiredBody); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						log.Error().Err(err).Msg(op)
						return
					}

					log.Debug().Err(err).Msg(op)

					return

				}
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, fmt.Errorf("%s: %w", op, err).Error(), http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, fmt.Errorf("%s: token is invalid", op).Error(), http.StatusUnauthorized)
				return
			}
		}

		log.Debug().Bool("useAuth", useAuth).Msg(fmt.Sprintf("%s - ok", op))

		next.ServeHTTP(w, r)
	})
}

func RunAPI(urlStr string) error {
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

	mux.Get("/api/v1/refresh", handlers.RefreshHandler)

	log.Info().Str("Running on", urlStr).Msg("http server started")

	return http.ListenAndServe(urlStr, mux)
}
