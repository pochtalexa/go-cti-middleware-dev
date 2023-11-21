package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
)

var Storage = NewStorage()

type StStorage struct {
	DB *sql.DB
}

func NewStorage() *StStorage {
	return &StStorage{}
}

func InitConnDB(appConfig *config.Config) (*StStorage, error) {
	var err error

	Storage.DB, err = sql.Open("pgx", appConfig.DB.DBConn)
	if err != nil {
		return nil, fmt.Errorf("InitConnDB: %w", err)
	}

	return Storage, nil
}

func (s *StStorage) SaveUser(login string, passHash []byte) (uid int64, err error) {
	var id int64
	const op = "storage.SaveUser"

	insertUser := `INSERT INTO users (login, pass_hash) VALUES ($1, $2) RETURNING id`

	stmt, err := s.DB.Prepare(insertUser)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(login, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return -1, fmt.Errorf("%s: %w", op, errors.New("login already exists"))
		} else {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}
