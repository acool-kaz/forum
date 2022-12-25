package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"os"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.Database.DBName+cfg.Database.FKeysConstraint)
	if err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	if err = createTables(cfg.Database.MigrationsUp, db); err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	return db, nil
}

func createTables(mifrationFile string, db *sql.DB) error {
	migrationData, err := os.ReadFile(mifrationFile)
	if err != nil {
		return fmt.Errorf("create tables: read file: %w", err)
	}
	if _, err = db.Exec(string(migrationData)); err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
