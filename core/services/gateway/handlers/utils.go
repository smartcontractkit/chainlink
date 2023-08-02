package handlers

import (
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	global       *rate.Limiter
	perUser      map[string]*rate.Limiter
	perUserRPS   rate.Limit
	perUserBurst int
	mu           sync.Mutex
}

func NewRateLimiter(globalRPS float64, globalBurst int, perUserRPS float64, perUserBurst int) *RateLimiter {
	return &RateLimiter{
		global:       rate.NewLimiter(rate.Limit(globalRPS), globalBurst),
		perUser:      make(map[string]*rate.Limiter),
		perUserRPS:   rate.Limit(perUserRPS),
		perUserBurst: perUserBurst,
	}
}

func (rl *RateLimiter) Allow(user string) bool {
	if !rl.global.Allow() {
		return false
	}

	rl.mu.Lock()
	userLimiter, ok := rl.perUser[user]
	if !ok {
		userLimiter = rate.NewLimiter(rl.perUserRPS, rl.perUserBurst)
		rl.perUser[user] = userLimiter
	}
	rl.mu.Unlock()

	return userLimiter.Allow()
}
