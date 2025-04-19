package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var _ Cache = (*InMemoryCache)(nil)

var (
	ErrValueNotStringerType = errors.New("value can not be casted into a string")
	ErrNoKey                = errors.New("key does not exist")
)

type InMemoryCache struct {
	data *sync.Map
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{data: &sync.Map{}}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (any, error) {
	value, exists := c.data.Load(key)
	if !exists {
		return "", ErrNoKey
	}

	return value, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.data.Store(key, value)

	if expiration < 1 {
		return nil
	}

	go func() {
		<-time.After(expiration)
		if err := c.Del(ctx, key); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to delete key %s: %s", key, err))
		}
	}()

	return nil
}

func (c *InMemoryCache) Del(_ context.Context, key string) error {
	c.data.Delete(key)

	return nil
}
