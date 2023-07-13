package common

import (
	"errors"
	"sync"

	"golang.org/x/time/rate"
)

// Wrapper around Go's rate.Limiter that supports both global and a per-user rate limiting.
type RateLimiter struct {
	global  *rate.Limiter
	perUser map[string]*rate.Limiter
	config  RateLimiterConfig
	mu      sync.Mutex
}

type RateLimiterConfig struct {
	GlobalRPS    float64 `json:"globalRPS"`
	GlobalBurst  int     `json:"globalBurst"`
	PerUserRPS   float64 `json:"perUserRPS"`
	PerUserBurst int     `json:"perUserBurst"`
}

func NewRateLimiter(config RateLimiterConfig) (*RateLimiter, error) {
	if config.GlobalRPS <= 0.0 || config.PerUserRPS <= 0.0 {
		return nil, errors.New("RPS values must be positive")
	}
	if config.GlobalBurst <= 0 || config.PerUserBurst <= 0 {
		return nil, errors.New("burst values must be positive")
	}
	return &RateLimiter{
		global:  rate.NewLimiter(rate.Limit(config.GlobalRPS), config.GlobalBurst),
		perUser: make(map[string]*rate.Limiter),
		config:  config,
	}, nil
}

func (rl *RateLimiter) Allow(user string) bool {
	rl.mu.Lock()
	userLimiter, ok := rl.perUser[user]
	if !ok {
		userLimiter = rate.NewLimiter(rate.Limit(rl.config.PerUserRPS), rl.config.PerUserBurst)
		rl.perUser[user] = userLimiter
	}
	rl.mu.Unlock()

	return userLimiter.Allow() && rl.global.Allow()
}
