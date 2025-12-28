package ratelimiter

import (
	"errors"
	"sync"

	"golang.org/x/time/rate"
)

// Wrapper around Go's rate.Limiter that supports both global and a per-sender rate limiting.
type RateLimiter struct {
	global    *rate.Limiter
	perSender map[string]*rate.Limiter
	config    Config
	mu        sync.Mutex
}

type Config struct {
	GlobalRPS      float64 `json:"globalRPS"`
	GlobalBurst    int     `json:"globalBurst"`
	PerSenderRPS   float64 `json:"perSenderRPS"`
	PerSenderBurst int     `json:"perSenderBurst"`
}

func NewRateLimiter(cfg Config) (*RateLimiter, error) {
	if cfg.GlobalRPS <= 0.0 || cfg.PerSenderRPS <= 0.0 {
		return nil, errors.New("RPS values must be positive")
	}
	if cfg.GlobalBurst <= 0 || cfg.PerSenderBurst <= 0 {
		return nil, errors.New("burst values must be positive")
	}

	return &RateLimiter{
		global:    rate.NewLimiter(rate.Limit(cfg.GlobalRPS), cfg.GlobalBurst),
		perSender: make(map[string]*rate.Limiter),
		config:    cfg,
	}, nil
}

func (rl *RateLimiter) Allow(sender string) (senderAllow bool, globalAllow bool) {
	rl.mu.Lock()
	senderLimiter, ok := rl.perSender[sender]
	if !ok {
		senderLimiter = rate.NewLimiter(rate.Limit(rl.config.PerSenderRPS), rl.config.PerSenderBurst)
		rl.perSender[sender] = senderLimiter
	}
	rl.mu.Unlock()

	senderAllow = senderLimiter.Allow()
	globalAllow = rl.global.Allow()
	return senderAllow, globalAllow
}
