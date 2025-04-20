package cache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

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
			name:          "existing key returns value",
			key:           "test-key",
			value:         "test-value",
			setupCache:    func(c cache.Cache) { _ = c.Set(context.Background(), "test-key", "test-value", 0) },
			expectedValue: "test-value",
			expectedErr:   nil,
		},
		{
			name:          "non-existent key returns error",
			key:           "non-existent",
			value:         nil,
			setupCache:    func(c cache.Cache) {},
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
		{
			name:          "empty key returns error",
			key:           "",
			value:         nil,
			setupCache:    func(c cache.Cache) {},
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
		{
			name:          "expired key returns error",
			key:           "expired-key",
			value:         "expired-value",
			setupCache:    func(c cache.Cache) { _ = c.Set(context.Background(), "expired-key", "expired-value", time.Millisecond) },
			expectedValue: "",
			expectedErr:   cache.ErrNoKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCache := cache.NewInMemoryCache()
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
				_ = c.Set(context.Background(), "test-key", "test-value", 0)
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent key",
			key:     "non-existent",
			setup:   func(c cache.Cache) {},
			wantErr: false,
		},
		{
			name:    "delete empty key",
			key:     "",
			setup:   func(c cache.Cache) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCache := cache.NewInMemoryCache()
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
	ctx := context.Background()
	const numGoroutines = 10
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", routineID, j)
				value := fmt.Sprintf("value-%d-%d", routineID, j)

				err := newCache.Set(ctx, key, value, 0)
				require.NoError(t, err)

				gotValue, err := newCache.Get(ctx, key)
				require.NoError(t, err)
				require.Equal(t, value, gotValue)

				err = newCache.Del(ctx, key)
				require.NoError(t, err)
			}
		}(i)
	}

	wg.Wait()
}
