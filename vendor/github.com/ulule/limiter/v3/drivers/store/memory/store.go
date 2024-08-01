package memory

import (
	"context"
	"strings"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/common"
)

// Store is the in-memory store.
type Store struct {
	// Prefix used for the key.
	Prefix string
	// cache used to store values in-memory.
	cache *CacheWrapper
}

// NewStore creates a new instance of memory store with defaults.
func NewStore() limiter.Store {
	return NewStoreWithOptions(limiter.StoreOptions{
		Prefix:          limiter.DefaultPrefix,
		CleanUpInterval: limiter.DefaultCleanUpInterval,
	})
}

// NewStoreWithOptions creates a new instance of memory store with options.
func NewStoreWithOptions(options limiter.StoreOptions) limiter.Store {
	return &Store{
		Prefix: options.Prefix,
		cache:  NewCache(options.CleanUpInterval),
	}
}

// Get returns the limit for given identifier.
func (store *Store) Get(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Increment(store.getCacheKey(key), 1, rate.Period)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// Increment increments the limit by given count & returns the new limit value for given identifier.
func (store *Store) Increment(ctx context.Context, key string, count int64, rate limiter.Rate) (limiter.Context, error) {
	newCount, expiration := store.cache.Increment(store.getCacheKey(key), count, rate.Period)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, newCount)
	return lctx, nil
}

// Peek returns the limit for given identifier, without modification on current values.
func (store *Store) Peek(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Get(store.getCacheKey(key), rate.Period)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// Reset returns the limit for given identifier.
func (store *Store) Reset(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Reset(store.getCacheKey(key), rate.Period)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// getCacheKey returns the full path for an identifier.
func (store *Store) getCacheKey(key string) string {
	buffer := strings.Builder{}
	buffer.WriteString(store.Prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}
