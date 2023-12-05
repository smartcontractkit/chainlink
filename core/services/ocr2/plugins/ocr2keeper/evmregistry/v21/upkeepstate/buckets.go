package upkeepstate

import (
	"context"
	"sync"
	"time"
)

// tokenBuckets provides a thread-safe token buckets for rate limiting.
// All buckets are cleaned up every cleanupInterval, to ensure the number of buckets
// is kept to a minimum.
// This implementation is specialized for the use case of the upkeep state store,
// where we want to avoid checking the same workID more than workIDRateLimit times
// in a workIDRatePeriod, while keeping memory footprint low.
type tokenBuckets struct {
	mutex sync.RWMutex

	cleanupInterval time.Duration
	maxTokens       uint32

	buckets map[string]uint32
}

func newTokenBuckets(maxTokens uint32, cleanupInterval time.Duration) *tokenBuckets {
	return &tokenBuckets{
		cleanupInterval: cleanupInterval,
		maxTokens:       maxTokens,
		buckets:         make(map[string]uint32),
	}
}

// Start starts the cleanup goroutine that runs every t.cleanupInterval.
func (t *tokenBuckets) Start(ctx context.Context) {
	ticker := time.NewTicker(t.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.clean()
		}
	}
}

// accept returns true if the bucket has enough tokens to accept the request,
// otherwise it returns false. It also updates the bucket with the updated number of tokens.
func (t *tokenBuckets) Accept(key string, tokens uint32) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	bucket, ok := t.buckets[key]
	if !ok {
		bucket := t.maxTokens
	}
	if bucket < tokens {
		return false
	}
	t.buckets[key] = bucket - tokens

	return true
}

func (t *tokenBuckets) clean() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.buckets = make(map[string]uint32)
}
