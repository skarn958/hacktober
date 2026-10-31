package db

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func Open(path string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)", path)

	database, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	database.SetMaxOpenConns(1)

	if _, err := database.Exec(schema); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}
	return database, nil
}
