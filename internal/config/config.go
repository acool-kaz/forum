package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Handler struct {
		Addr         string `json:"addr"`
		WriteTimeout int    `json:"write_timeout"`
		ReadTimeout  int    `json:"read_timeout"`
	} `json:"handler"`
	Database struct {
		DBName          string `json:"db_name"`
		MigrationsUp    string `json:"migrations_up"`
		FKeysConstraint string `json:"f_keys_constraint"`
	} `json:"database"`
}

func InitConfig(path string) (*Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("init configs: open: %w", err)
	}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("init configs: decode: %w", err)
	}
	return &cfg, nil
}
