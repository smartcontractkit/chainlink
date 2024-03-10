package common

import (
	"errors"
	"sync"

	"golang.org/x/time/rate"
)

// Wrapper around Go's rate.Limiter that supports both global and a per-sender rate limiting.
type RateLimiter struct {
	global    *rate.Limiter
	perSender map[string]*rate.Limiter
	config    RateLimiterConfig
	mu        sync.Mutex
}

type RateLimiterConfig struct {
	GlobalRPS      float64 `json:"globalRPS"`
	GlobalBurst    int     `json:"globalBurst"`
	PerSenderRPS   float64 `json:"perSenderRPS"`
	PerSenderBurst int     `json:"perSenderBurst"`
}

func NewRateLimiter(config RateLimiterConfig) (*RateLimiter, error) {
	if config.GlobalRPS <= 0.0 || config.PerSenderRPS <= 0.0 {
		return nil, errors.New("RPS values must be positive")
	}
	if config.GlobalBurst <= 0 || config.PerSenderBurst <= 0 {
		return nil, errors.New("burst values must be positive")
	}
	return &RateLimiter{
		global:    rate.NewLimiter(rate.Limit(config.GlobalRPS), config.GlobalBurst),
		perSender: make(map[string]*rate.Limiter),
		config:    config,
	}, nil
}

func (rl *RateLimiter) Allow(sender string) bool {
	rl.mu.Lock()
	senderLimiter, ok := rl.perSender[sender]
	if !ok {
		senderLimiter = rate.NewLimiter(rate.Limit(rl.config.PerSenderRPS), rl.config.PerSenderBurst)
		rl.perSender[sender] = senderLimiter
	}
	rl.mu.Unlock()

	return senderLimiter.Allow() && rl.global.Allow()
}
