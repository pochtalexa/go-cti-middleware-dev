package auth

import (
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(login string, pass string) (int64, error) {
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
	id, err := storage.Storage.SaveUser(login, passHash)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
