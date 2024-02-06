package server

import (
	"fmt"

	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/redis"
)

func connectToRedis(cfg *config.RedisConfig) (*redis.Cache, error) {
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

	return cache, nil
}

func disconnectFromRedis(cache *redis.Cache) {
	cache.Stop()
}
