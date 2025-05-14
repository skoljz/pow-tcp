package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

type Config struct {
	Addr         string `env:"TCP_SERVER_ADDR"   env-default:":9000"`
	TargetSize   uint64 `env:"POW_TARGET_SIZE"   env-default:"8"`
	LoggingLevel string `env:"LOG_LEVEL"         env-default:"info"`
}

var DefaultLogLevel = slog.LevelInfo

func New() (*Config, error) {
	var cfg Config
	_ = cleanenv.ReadConfig(".env", &cfg)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}
	return &cfg, nil
}
