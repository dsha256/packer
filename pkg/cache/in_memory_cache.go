package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var _ ClosableCache = (*InMemoryCache)(nil)

var ErrNoKey = errors.New("key does not exist")

type cacheItem struct {
	value      any
	expiration time.Time
}

func (item *cacheItem) isExpired() bool {
	if item.expiration.IsZero() {
		return false
	}

	return time.Now().After(item.expiration)
}

type InMemoryCache struct {
	data        *sync.Map
	stopCleanup chan struct{}
	cleanupWg   sync.WaitGroup
}

func NewInMemoryCache() ClosableCache {
	return NewInMemoryCacheWithCleanup(time.Second * 10)
}

func NewInMemoryCacheWithCleanup(cleanupInterval time.Duration) ClosableCache {
	cache := &InMemoryCache{
		data:        &sync.Map{},
		stopCleanup: make(chan struct{}),
	}

	cache.cleanupWg.Add(1)
	go cache.cleanupExpiredKeys(cleanupInterval)

	return cache
}

func (cache *InMemoryCache) Get(_ context.Context, key string) (any, error) {
	value, exists := cache.data.Load(key)
	if !exists {
		return "", ErrNoKey
	}

	item, ok := value.(*cacheItem)
	if !ok {
		return "", ErrNoKey
	}

	if item.isExpired() {
		cache.data.Delete(key)

		return "", ErrNoKey
	}

	return item.value, nil
}

func (cache *InMemoryCache) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	item := &cacheItem{
		value: value,
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
	}

	cache.data.Store(key, item)

	return nil
}

func (cache *InMemoryCache) Del(_ context.Context, key string) error {
	cache.data.Delete(key)

	return nil
}

func (cache *InMemoryCache) Close() {
	close(cache.stopCleanup)
	cache.cleanupWg.Wait()
}

func (cache *InMemoryCache) cleanupExpiredKeys(interval time.Duration) {
	defer cache.cleanupWg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.data.Range(func(key, value any) bool {
				item, ok := value.(*cacheItem)
				if !ok {
					return true
				}

				if item.isExpired() {
					cache.data.Delete(key)
				}

				return true
			})
		case <-cache.stopCleanup:
			return
		}
	}
}
