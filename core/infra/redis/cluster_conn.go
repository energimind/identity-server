package redis

import (
	"context"
	"time"

	"github.com/energimind/identity-service/pkg/hostutil"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// clusterConnection implements a connection to a redis cluster.
type clusterConnection struct {
	cluster *redis.ClusterClient
}

// Ensure clusterConnection implements connection interface.
var _ connection = (*clusterConnection)(nil)

// newClusterConnection creates a new cluster connection.
func newClusterConnection(config Config) (*clusterConnection, error) {
	addresses := hostutil.ComposeAddressList(config.Host, config.Port, defaultPort)
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           addresses,
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
		return nil, errors.Wrapf(err, "failed to connect to redis (%s)", addresses)
	}

	return &clusterConnection{cluster: client}, nil
}

func (c *clusterConnection) close() error {
	//nolint:wrapcheck
	return c.cluster.Close()
}

func (c *clusterConnection) set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return c.cluster.Set(ctx, key, value, expiration)
}

func (c *clusterConnection) get(ctx context.Context, key string) *redis.StringCmd {
	return c.cluster.Get(ctx, key)
}

func (c *clusterConnection) delete(ctx context.Context, key string) *redis.IntCmd {
	return c.cluster.Del(ctx, key)
}
