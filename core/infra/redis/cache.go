package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/energimind/identity-service/core/domain/cache"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const ioTimeout = 10 * time.Second

// Cache implements the cache interface.
//
// It uses Redis as the underlying storage.
//
// It implements the cache.Cache interface.
type Cache struct {
	conn      connection
	namespace string
}

// Ensure service implements the cache.Cache interface.
var _ cache.Cache = (*Cache)(nil)

// NewCache creates a new cache. The cache is in the connected state.
func NewCache(config Config) (*Cache, error) {
	if config.Namespace == "" {
		return nil, errors.New("missing redis namespace")
	}

	var conn connection

	if config.Standalone {
		sConn, err := newStandaloneConnection(config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create standalone redis conn")
		}

		conn = sConn
	} else {
		cConn, err := newClusterConnection(config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create cluster redis conn")
		}

		conn = cConn
	}

	return &Cache{
		conn:      conn,
		namespace: config.Namespace,
	}, nil
}

// Stop stops the cache.
func (c *Cache) Stop() {
	_ = c.conn.close()
}

// Put implements the cache.Cache interface.
func (c *Cache) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return NewCacheError("failed to marshal value: %v", err)
	}

	if sErr := c.conn.set(ctx, c.fqn(key), b, ttl).Err(); sErr != nil {
		return NewCacheError("failed to set key value: %s", sErr)
	}

	return nil
}

// Get implements the cache.Cache interface.
func (c *Cache) Get(ctx context.Context, key string, receiver any) (bool, error) {
	cmd := c.conn.get(ctx, c.fqn(key))

	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, NewCacheError("failed to get key value: %v", err)
	}

	b, err := cmd.Bytes()
	if err != nil {
		return false, NewCacheError("failed to get key value: %v", err)
	}

	if uErr := json.Unmarshal(b, receiver); uErr != nil {
		return false, NewCacheError("failed to unmarshal value: %v", uErr)
	}

	return true, nil
}

// Delete implements the cache.Cache interface.
func (c *Cache) Delete(ctx context.Context, key string) error {
	if sErr := c.conn.delete(ctx, c.fqn(key)).Err(); sErr != nil {
		return NewCacheError("failed to delete key: %v", sErr)
	}

	return nil
}

func (c *Cache) fqn(key string) string {
	if c.namespace != "" {
		return fmt.Sprintf("%s.%s", c.namespace, key)
	}

	return key
}
