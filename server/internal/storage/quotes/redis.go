package quotes

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/skoljz/pow_tcp/internal/config"
)

type RedisStorage struct {
	client *redis.Client
	key    string
	rnd    *rand.Rand
}

const cacheKey = "quotes"

func NewRedis(cfg config.RedisConfig) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &RedisStorage{
		client: client,
		key:    cacheKey,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (rdb *RedisStorage) Random() (string, error) {
	ctx := context.Background()

	length, err := rdb.client.LLen(ctx, rdb.key).Result()
	if err != nil {
		return "", fmt.Errorf("unknown error while try get lenght: %w", err)
	}
	if length == 0 {
		return "", ErrNoQuotes
	}

	idx := rdb.rnd.Int63n(length)
	return rdb.client.LIndex(ctx, rdb.key, idx).Result()
}

func (rdb *RedisStorage) Close() error {
	return rdb.client.Close()
}
