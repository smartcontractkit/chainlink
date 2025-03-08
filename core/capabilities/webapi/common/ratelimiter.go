package common

import (
	"sync"

	"golang.org/x/time/rate"
)

const (
	DefaultGlobalRPS     = 100.0
	DefaultGlobalBurst   = 100
	DefaultWorkflowRPS   = 5.0
	DefaultWorkflowBurst = 50
)

// Wrapper around Go's rate.Limiter that supports both global and a per-workflow rate limiting.
type RateLimiter struct {
	global      *rate.Limiter
	perWorkflow map[string]*rate.Limiter
	config      RateLimiterConfig
	mu          sync.Mutex
}

type RateLimiterConfig struct {
	GlobalRPS        float64 `json:"globalRPS"`
	GlobalBurst      int     `json:"globalBurst"`
	PerWorkflowRPS   float64 `json:"perWorkflowRPS"`
	PerWorkflowBurst int     `json:"perWorkflowBurst"`
}

func NewRateLimiter(config RateLimiterConfig) (*RateLimiter, error) {
	if config.GlobalBurst <= 0 {
		config.GlobalBurst = DefaultGlobalBurst
	}

	if config.GlobalRPS <= 0.0 {
		config.GlobalRPS = DefaultGlobalRPS
	}

	if config.PerWorkflowBurst <= 0 {
		config.PerWorkflowBurst = DefaultWorkflowBurst
	}

	if config.PerWorkflowRPS <= 0.0 {
		config.PerWorkflowRPS = DefaultWorkflowRPS
	}

	return &RateLimiter{
		global:      rate.NewLimiter(rate.Limit(config.GlobalRPS), config.GlobalBurst),
		perWorkflow: make(map[string]*rate.Limiter),
		config:      config,
	}, nil
}

// Allow checks that the workflow is not rate limited.
func (rl *RateLimiter) Allow(id string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	wfLimiter, ok := rl.perWorkflow[id]
	if !ok {
		wfLimiter = rate.NewLimiter(rate.Limit(rl.config.PerWorkflowRPS), rl.config.PerWorkflowBurst)
		rl.perWorkflow[id] = wfLimiter
	}

	return wfLimiter.Allow() && rl.global.Allow()
}
