package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type Closable interface {
	Close()
}

type ClosableCache interface {
	Cache
	Closable
}
