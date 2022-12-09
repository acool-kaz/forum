package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"os"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	if err = createTables(cfg, db); err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	return db, nil
}

func createTables(cfg *config.Config, db *sql.DB) error {
	migrationData, err := os.ReadFile(cfg.Database.MigrationsUp)
	if err != nil {
		return fmt.Errorf("create tables: read file: %w", err)
	}
	if _, err = db.Exec(string(migrationData)); err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
