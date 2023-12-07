package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/models"
)

var (
	Storage         = NewStorage()
	ErrUserNotFound = errors.New("user not found")
)

type StStorage struct {
	DB *sql.DB
}

func NewStorage() *StStorage {
	return &StStorage{}
}

func InitConnDB() (*StStorage, error) {
	var err error

	Storage.DB, err = sql.Open("pgx", config.ServerConfig.DB.DBConn)
	if err != nil {
		return nil, fmt.Errorf("InitConnDB: %w", err)
	}

	return Storage, nil
}

func (s *StStorage) SaveAgent(login string, passHash []byte) (uid int64, err error) {
	var id int64
	const op = "storage.SaveAgent"

	insertAgent := `INSERT INTO users (login, pass_hash) VALUES ($1, $2) RETURNING id`

	stmt, err := s.DB.Prepare(insertAgent)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(login, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return -1, fmt.Errorf("%s: %w", op, errors.New("login already exists"))
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *StStorage) GetAgent(login string) (*models.StAgent, error) {
	const op = "storage.GetAgent"

	agent := models.NewAgent()
	agent.Login = login

	getAgent := `SELECT id, pass_hash FROM users WHERE login = $1`

	stmt, err := s.DB.Prepare(getAgent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(login).Scan(agent.ID, agent.PassHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return agent, nil
}
