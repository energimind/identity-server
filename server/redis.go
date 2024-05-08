package server

import (
	"fmt"

	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/core/config"
	"github.com/energimind/identity-server/core/infra/redis"
)

func connectRedis(cfg config.RedisConfig, closer *closer) (*redis.Cache, error) {
	cache, err := redis.NewCache(redis.Config{
		Host:       cfg.Host,
		Port:       cfg.Port,
		Username:   cfg.Username,
		Password:   cfg.Password,
		Namespace:  cfg.Namespace,
		Standalone: cfg.Standalone,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis cache: %w", err)
	}

	closer.add(cache.Stop)

	slog.Info().Msg("Connected to redis")

	return cache, nil
}
