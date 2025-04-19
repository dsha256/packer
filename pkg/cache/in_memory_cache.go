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

var ErrNoKey = errors.New("key does not exist")

type InMemoryCache struct {
	data *sync.Map
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{data: &sync.Map{}}
}

func (cache *InMemoryCache) Get(_ context.Context, key string) (any, error) {
	value, exists := cache.data.Load(key)
	if !exists {
		return "", ErrNoKey
	}

	return value, nil
}

func (cache *InMemoryCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	cache.data.Store(key, value)

	if expiration < 1 {
		return nil
	}

	go func() {
		<-time.After(expiration)
		if err := cache.Del(ctx, key); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to delete key %s: %s", key, err))
		}
	}()

	return nil
}

func (cache *InMemoryCache) Del(_ context.Context, key string) error {
	cache.data.Delete(key)

	return nil
}
