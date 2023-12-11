package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type IntStorage interface {
	SaveAgent(login string, passHash []byte) (uid int64, err error)
	GetAgent(login string) (*storage.StAgent, error)
}

// NewToken генерируем JWT токен для агента
func NewToken(agent *storage.StAgent, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := storage.NewClaims()

	claims.ID = agent.ID
	claims.Login = agent.Login
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен, используя секретный ключ
	tokenString, err := token.SignedString([]byte(config.ServerConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("NewToken: %w", err)
	}

	return tokenString, nil
}

func RegisterNewUser(login string, pass string, curStorage IntStorage) (int64, error) {
	// op (operation) - имя текущей функции и пакета. Такую метку удобно
	// добавлять в логи и в текст ошибок, чтобы легче было искать хвосты
	const op = "auth.RegisterNewUser"

	// Создаём локальный объект логгера с доп. полями, содержащими полезную инфу
	log := log.With().Str("op", op).Str("login", login).Logger()

	log.Debug().Msg("registering user")

	// Генерируем хэш и соль для пароля
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем пользователя в БД
	id, err := curStorage.SaveAgent(login, passHash)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Login checks if user with given credentials exists in the system and returns access token.
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func Login(login string, password string, curStorage IntStorage) (string, error) {
	const op = "auth.Login"

	log := log.With().Str("op", op).Str("login", login).Logger()
	log.Debug().Msg("attempting to login user")

	agent, err := curStorage.GetAgent(login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error().Err(err).Msg("empty select")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error().Err(err).Msg("failed to get user")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(agent.PassHash, []byte(password)); err != nil {
		log.Error().Err(err).Msg("invalid credentials")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Debug().Msg("user logged in successfully")

	// Создаём токен авторизации
	token, err := NewToken(agent, config.ServerConfig.TokenTTL)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
