package migrations

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
)

//go:embed *.sql
var SQLFiles embed.FS

func ApplyMigrations() error {
	fsys := SQLFiles

	goose.SetBaseFS(fsys)
	goose.SetSequential(true)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(storage.Storage.DB, "."); err != nil {
		return fmt.Errorf("goose: %w", err)
	}

	return nil
}
