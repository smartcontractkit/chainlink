package limiter

import "time"

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "limiter"

	// DefaultMaxRetry is the default maximum number of key retries under
	// race condition (mainly used with database-based stores).
	DefaultMaxRetry = 3

	// DefaultCleanUpInterval is the default time duration for cleanup.
	DefaultCleanUpInterval = 30 * time.Second
)
