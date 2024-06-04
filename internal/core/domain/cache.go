package domain

import (
	"context"
	"time"
)

// Cache defines a cache interface.
type Cache interface {
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, receiver any) (bool, error)
	Delete(ctx context.Context, key string) error
}
