package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultPort     = "6379"
	maxRetries      = 5
	dialTimeout     = 5 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	poolSize        = 2
	minIdleConns    = 1
	maxConnLifetime = 30 * time.Minute
	poolTimeout     = 5 * time.Second
)

type connection interface {
	close() error
	set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	get(ctx context.Context, key string) *redis.StringCmd
	delete(ctx context.Context, key string) *redis.IntCmd
}
