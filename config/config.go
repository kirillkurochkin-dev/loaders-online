package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB PostgresConnection
}

type PostgresConnection struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	return cfg, nil
}
