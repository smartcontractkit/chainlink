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
}

// StoreOptions are options for store.
type StoreOptions struct {
	// Prefix is the prefix to use for the key.
	Prefix string

	// MaxRetry is the maximum number of retry under race conditions.
	MaxRetry int

	// CleanUpInterval is the interval for cleanup.
	CleanUpInterval time.Duration
}
