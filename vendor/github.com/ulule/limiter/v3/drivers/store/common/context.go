package common

import (
	"time"

	"github.com/ulule/limiter/v3"
)

// GetContextFromState generate a new limiter.Context from given state.
func GetContextFromState(now time.Time, rate limiter.Rate, expiration time.Time, count int64) limiter.Context {
	limit := rate.Limit
	remaining := int64(0)
	reached := true

	if count <= limit {
		remaining = limit - count
		reached = false
	}

	reset := expiration.Unix()

	return limiter.Context{
		Limit:     limit,
		Remaining: remaining,
		Reset:     reset,
		Reached:   reached,
	}
}
