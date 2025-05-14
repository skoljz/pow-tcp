package config

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr            string        `env:"TCP_SERVER_ADDR" env-default:":9000"`
	PowComplexity   uint8         `env:"POW_COMPLEXITY"   env-default:"20"`
	StorageFile     string        `env:"STORAGE_FILE"    env-default:"quotes.txt"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"    env-default:"30s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"   env-default:"30s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	LoggingLevel    string        `env:"LOG_LEVEL"       env-default:"info"`

	RedisConfig RedisConfig
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

func New() (*Config, error) {
	var cfg Config

	_ = cleanenv.ReadConfig(".env", &cfg)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	return &cfg, nil
}

func (c *Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LoggingLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
