package web

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type localRateLimiter struct {
	mu          sync.Mutex
	period      time.Duration
	limit       int64
	counters    map[string]localRateCounter
	lastCleanup time.Time
}

type localRateCounter struct {
	value      int64
	expiration time.Time
}

func newLocalRateLimiter(period time.Duration, limit int64) *localRateLimiter {
	return &localRateLimiter{
		period:      period,
		limit:       limit,
		counters:    make(map[string]localRateCounter),
		lastCleanup: time.Now(),
	}
}

func (l *localRateLimiter) increment(key string) (remaining, reset int64, reached bool) {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanup(now)
	counter, ok := l.counters[key]
	if !ok || counter.expiration.IsZero() || now.After(counter.expiration) {
		counter = localRateCounter{value: 1, expiration: now.Add(l.period)}
	} else {
		counter.value++
	}
	l.counters[key] = counter

	if counter.value <= l.limit {
		remaining = l.limit - counter.value
	} else {
		reached = true
	}
	return remaining, counter.expiration.Unix(), reached
}

func (l *localRateLimiter) cleanup(now time.Time) {
	if now.Sub(l.lastCleanup) < time.Minute {
		return
	}
	l.lastCleanup = now
	for key, counter := range l.counters {
		if counter.expiration.IsZero() || now.After(counter.expiration) {
			delete(l.counters, key)
		}
	}
}

func rateLimiter(period time.Duration, limit int64) gin.HandlerFunc {
	limiter := newLocalRateLimiter(period, limit)
	return func(c *gin.Context) {
		remaining, reset, reached := limiter.increment(c.ClientIP())
		c.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))
		if reached {
			c.String(http.StatusTooManyRequests, "Limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}
