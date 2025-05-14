package quotes

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/skoljz/pow_tcp/internal/config"
)

type Storage interface {
	Random() (string, error)
	Close() error
}

var ErrNoQuotes = errors.New("quotes: storage is empty")

func NewStorage(cfg *config.Config, log *slog.Logger) (Storage, error) {
	if cfg.RedisConfig.Addr != "" {
		rdb, err := NewRedis(cfg.RedisConfig)
		if err == nil {
			log.Info("using Redis storage", "addr", cfg.RedisConfig.Addr)
			return rdb, nil
		}
		log.Warn("redis unavailable, fallback to memory", "error", err)
	}

	mem, err := NewInMemory(cfg.StorageFile)
	if err != nil {
		return nil, fmt.Errorf("in-memory storage: %w", err)
	}
	log.Info("using in-memory storage", "file", cfg.StorageFile)
	return mem, nil
}
