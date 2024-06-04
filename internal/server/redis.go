package server

import (
	"context"
	"fmt"
	"net"

	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/internal/config"
	"github.com/energimind/identity-server/internal/core/infra/redis"
)

func connectRedis(ctx context.Context, cfg config.RedisConfig, closer *closer) (*redis.Cache, error) {
	connect := func() (*redis.Cache, error) {
		cache, err := redis.NewCache(ctx, redisDriverConfig(cfg))
		if err != nil {
			return nil, fmt.Errorf("failed to create redis cache: %w", err)
		}

		closer.add(cache.Stop)

		slog.Info().Str("address", net.JoinHostPort(cfg.Host, cfg.Port)).Msg("Connected to redis")

		return cache, nil
	}

	return withDelayReporter[*redis.Cache](ctx, connect, "redis-connect")
}

func redisDriverConfig(cfg config.RedisConfig) redis.Config {
	return redis.Config{
		Host:       cfg.Host,
		Port:       cfg.Port,
		Username:   cfg.Username,
		Password:   cfg.Password,
		Namespace:  cfg.Namespace,
		Standalone: cfg.Standalone,
	}
}
