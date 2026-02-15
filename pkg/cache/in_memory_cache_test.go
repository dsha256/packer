package cache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dsha256/packer/pkg/cache"
)

func TestInMemoryCache_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value         any
		expectedValue any
		expectedErr   error
		setupCache    func(c cache.Cache)
		name          string
		key           string
	}{
		{
			name:  "existing key returns value",
			key:   "test-key",
			value: "test-value",
			setupCache: func(c cache.Cache) {
				err := c.Set(context.Background(), "test-key", "test-value", 0)
				require.NoError(t, err)
			},
			expectedValue: "test-value",
			expectedErr:   nil,
		},
		{
			name:          "non-existent key returns error",
			key:           "non-existent",
			value:         nil,
			setupCache:    func(_ cache.Cache) {},
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
		{
			name:          "empty key returns error",
			key:           "",
			value:         nil,
			setupCache:    func(_ cache.Cache) {},
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
		{
			name:  "expired key returns error",
			key:   "expired-key",
			value: "expired-value",
			setupCache: func(c cache.Cache) {
				err := c.Set(context.Background(), "expired-key", "expired-value", time.Millisecond)
				require.NoError(t, err)
			},
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCache := cache.NewInMemoryCache()
			defer closeCache(newCache)

			tt.setupCache(newCache)

			if tt.name == "expired key returns error" {
				time.Sleep(time.Millisecond * 2)
			}

			value, err := newCache.Get(context.Background(), tt.key)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, value)
		})
	}
}

func TestInMemoryCache_Set(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value      any
		name       string
		key        string
		expiration time.Duration
		wantErr    bool
	}{
		{
			name:       "set value with no expiration",
			key:        "test-key",
			value:      "test-value",
			expiration: 0,
			wantErr:    false,
		},
		{
			name:       "set value with expiration",
			key:        "expiring-key",
			value:      "expiring-value",
			expiration: time.Millisecond,
			wantErr:    false,
		},
		{
			name:       "set empty key",
			key:        "",
			value:      "empty-key-value",
			expiration: 0,
			wantErr:    false,
		},
		{
			name:       "set nil value",
			key:        "nil-value-key",
			value:      nil,
			expiration: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCache := cache.NewInMemoryCache()
			defer closeCache(newCache)

			err := newCache.Set(context.Background(), tt.key, tt.value, tt.expiration)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			if tt.expiration > 0 {
				time.Sleep(tt.expiration + time.Millisecond)
				_, err = newCache.Get(context.Background(), tt.key)
				require.ErrorIs(t, err, cache.ErrNoKey)

				return
			}

			value, err := newCache.Get(context.Background(), tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.value, value)
		})
	}
}

func TestInMemoryCache_Del(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup   func(c cache.Cache)
		name    string
		key     string
		wantErr bool
	}{
		{
			name: "delete existing key",
			key:  "test-key",
			setup: func(c cache.Cache) {
				err := c.Set(context.Background(), "test-key", "test-value", 0)
				require.NoError(t, err)
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent key",
			key:     "non-existent",
			setup:   func(_ cache.Cache) {},
			wantErr: false,
		},
		{
			name:    "delete empty key",
			key:     "",
			setup:   func(_ cache.Cache) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCache := cache.NewInMemoryCache()
			defer closeCache(newCache)
			tt.setup(newCache)

			err := newCache.Del(context.Background(), tt.key)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			_, err = newCache.Get(context.Background(), tt.key)
			require.ErrorIs(t, err, cache.ErrNoKey)
		})
	}
}

func TestInMemoryCache_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	newCache := cache.NewInMemoryCache()
	defer closeCache(newCache)
	ctx := context.Background()
	const numGoroutines = 10
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for index := range numGoroutines {
		go func(routineID int) {
			defer wg.Done()

			for j := range operationsPerGoroutine {
				key := fmt.Sprintf("key-%d-%d", routineID, j)
				value := fmt.Sprintf("value-%d-%d", routineID, j)

				err := newCache.Set(ctx, key, value, 0)
				assert.NoError(t, err)

				gotValue, err := newCache.Get(ctx, key)
				assert.NoError(t, err)
				assert.Equal(t, value, gotValue)

				err = newCache.Del(ctx, key)
				assert.NoError(t, err)
			}
		}(index)
	}

	wg.Wait()
}

func closeCache(c cache.Cache) {
	if closable, ok := c.(cache.Closable); ok {
		closable.Close()
	}
}

func TestInMemoryCache_PassiveExpiration(t *testing.T) {
	t.Parallel()

	testCache := cache.NewInMemoryCacheWithCleanup(time.Hour)
	defer closeCache(testCache)

	ctx := context.Background()

	err := testCache.Set(ctx, "short-lived", "value", 50*time.Millisecond)
	require.NoError(t, err)

	val, err := testCache.Get(ctx, "short-lived")
	require.NoError(t, err)
	require.Equal(t, "value", val)

	time.Sleep(100 * time.Millisecond)

	_, err = testCache.Get(ctx, "short-lived")
	require.ErrorIs(t, err, cache.ErrNoKey, "expired key should be removed on Get (passive expiration)")
}

func TestInMemoryCache_ActiveExpiration(t *testing.T) {
	t.Parallel()

	testCache := cache.NewInMemoryCacheWithCleanup(50 * time.Millisecond)
	defer closeCache(testCache)

	ctx := context.Background()

	for i := range 10 {
		key := fmt.Sprintf("key-%d", i)
		err := testCache.Set(ctx, key, fmt.Sprintf("value-%d", i), 30*time.Millisecond)
		require.NoError(t, err)
	}

	for i := range 10 {
		key := fmt.Sprintf("key-%d", i)
		val, err := testCache.Get(ctx, key)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("value-%d", i), val)
	}

	time.Sleep(150 * time.Millisecond)

	for i := range 10 {
		key := fmt.Sprintf("key-%d", i)
		_, err := testCache.Get(ctx, key)
		require.ErrorIs(t, err, cache.ErrNoKey, "key %s should be removed by active cleanup", key)
	}
}

func TestInMemoryCache_MixedExpiration(t *testing.T) {
	t.Parallel()

	testCache := cache.NewInMemoryCacheWithCleanup(50 * time.Millisecond)
	defer closeCache(testCache)

	ctx := context.Background()

	err := testCache.Set(ctx, "no-expiry", "persistent", 0)
	require.NoError(t, err)

	err = testCache.Set(ctx, "short-expiry", "temporary", 50*time.Millisecond)
	require.NoError(t, err)

	err = testCache.Set(ctx, "long-expiry", "lasting", 500*time.Millisecond)
	require.NoError(t, err)

	_, err = testCache.Get(ctx, "no-expiry")
	require.NoError(t, err)
	_, err = testCache.Get(ctx, "short-expiry")
	require.NoError(t, err)
	_, err = testCache.Get(ctx, "long-expiry")
	require.NoError(t, err)

	time.Sleep(150 * time.Millisecond)

	_, err = testCache.Get(ctx, "short-expiry")
	require.ErrorIs(t, err, cache.ErrNoKey)

	val, err := testCache.Get(ctx, "no-expiry")
	require.NoError(t, err)
	require.Equal(t, "persistent", val)

	val, err = testCache.Get(ctx, "long-expiry")
	require.NoError(t, err)
	require.Equal(t, "lasting", val)

	time.Sleep(500 * time.Millisecond)

	_, err = testCache.Get(ctx, "long-expiry")
	require.ErrorIs(t, err, cache.ErrNoKey)

	val, err = testCache.Get(ctx, "no-expiry")
	require.NoError(t, err)
	require.Equal(t, "persistent", val)
}

func TestInMemoryCache_UpdateExpiration(t *testing.T) {
	t.Parallel()

	testCache := cache.NewInMemoryCacheWithCleanup(50 * time.Millisecond)
	defer closeCache(testCache)

	ctx := context.Background()

	err := testCache.Set(ctx, "key", "value1", 100*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(60 * time.Millisecond)

	err = testCache.Set(ctx, "key", "value2", 200*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(60 * time.Millisecond)

	val, err := testCache.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, "value2", val, "key should have new value with extended expiration")

	time.Sleep(150 * time.Millisecond)

	_, err = testCache.Get(ctx, "key")
	require.ErrorIs(t, err, cache.ErrNoKey, "key should expire after new expiration time")
}

func TestInMemoryCache_CloseStopsCleanup(t *testing.T) {
	t.Parallel()

	c := cache.NewInMemoryCacheWithCleanup(10 * time.Millisecond)

	done := make(chan struct{})
	go func() {
		closeCache(c)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Close() did not complete in time - cleanup goroutine may not have stopped")
	}
}
