package redis

import (
	"context"
	"time"

	"github.com/energimind/identity-server/pkg/hostutil"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// standaloneConnection implements a connection to a standalone redis instance.
type standaloneConnection struct {
	client *redis.Client
}

// Ensure standaloneConnection implements connection interface.
var _ connection = (*standaloneConnection)(nil)

// newStandaloneConnection creates a new standalone connection.
func newStandaloneConnection(config Config) (*standaloneConnection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            hostutil.ComposeAddress(config.Host, config.Port, defaultPort),
		Username:        config.Username,
		Password:        config.Password,
		MaxRetries:      maxRetries,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		MinIdleConns:    minIdleConns,
		ConnMaxLifetime: maxConnLifetime,
		PoolSize:        poolSize,
		PoolTimeout:     poolTimeout,
	})

	ctx, cancelFunc := context.WithTimeout(context.Background(), ioTimeout)
	defer cancelFunc()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to connect to redis (%s)", client.Options().Addr)
	}

	return &standaloneConnection{client: client}, nil
}

func (c *standaloneConnection) close() error {
	//nolint:wrapcheck
	return c.client.Close()
}

func (c *standaloneConnection) set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

func (c *standaloneConnection) get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

func (c *standaloneConnection) delete(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Del(ctx, key)
}
