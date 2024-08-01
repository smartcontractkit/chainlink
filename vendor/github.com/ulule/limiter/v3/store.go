package limiter

import (
	"context"
	"time"
)

// Store is the common interface for limiter stores.
type Store interface {
	// Get returns the limit for given identifier.
	Get(ctx context.Context, key string, rate Rate) (Context, error)
	// Peek returns the limit for given identifier, without modification on current values.
	Peek(ctx context.Context, key string, rate Rate) (Context, error)
	// Reset resets the limit to zero for given identifier.
	Reset(ctx context.Context, key string, rate Rate) (Context, error)
	// Increment increments the limit by given count & gives back the new limit for given identifier
	Increment(ctx context.Context, key string, count int64, rate Rate) (Context, error)
}

// StoreOptions are options for store.
type StoreOptions struct {
	// Prefix is the prefix to use for the key.
	Prefix string

	// MaxRetry is the maximum number of retry under race conditions on redis store.
	// Deprecated: this option is no longer required since all operations are atomic now.
	MaxRetry int

	// CleanUpInterval is the interval for cleanup (run garbage collection) on stale entries on memory store.
	// Setting this to a low value will optimize memory consumption, but will likely
	// reduce performance and increase lock contention.
	// Setting this to a high value will maximum throughput, but will increase the memory footprint.
	CleanUpInterval time.Duration
}
