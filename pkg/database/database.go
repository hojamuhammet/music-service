package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func NewDatabase(addr string) (*sql.DB, error) {
	if addr == "" {
		return nil, fmt.Errorf("database address is empty")
	}

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetLogger(log.Default())

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func CloseDatabase(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	err := db.Close()
	if err != nil {
		return fmt.Errorf("failed to close the database connection: %w", err)
	}
	return nil
}
