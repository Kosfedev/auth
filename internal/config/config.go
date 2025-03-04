package config

import (
	"time"

	"github.com/joho/godotenv"
)

// PGConfig is...
type PGConfig interface {
	DSN() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// Load is...
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
