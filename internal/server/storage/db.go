package storage

import (
	"database/sql"
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
)

var DBconn *sql.DB

func InitConnDB(appConfig *config.Config) (*sql.DB, error) {
	var err error

	DBconn, err = sql.Open("pgx", appConfig.DB.DBConn)
	if err != nil {
		return nil, fmt.Errorf("InitConnDB: %w", err)
	}

	return DBconn, nil
}

func SaveUser(login string, passHash []byte) (uid int64, err error) {

}
