package ratelimiter

import (
	"errors"
	"fmt"

	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type Config struct {
	GlobalRPS      float64 `json:"globalRPS"`
	GlobalBurst    int     `json:"globalBurst"`
	PerSenderRPS   float64 `json:"perSenderRPS"`
	PerSenderBurst int     `json:"perSenderBurst"`
}

func NewRateLimiter(cfg Config, f limits.Factory) (limits.RateLimiter, error) {
	if cfg.GlobalRPS <= 0.0 || cfg.PerSenderRPS <= 0.0 {
		return nil, errors.New("RPS values must be positive")
	}
	if cfg.GlobalBurst <= 0 || cfg.PerSenderBurst <= 0 {
		return nil, errors.New("burst values must be positive")
	}

	globalLimit := cresettings.Default.WorkflowTriggerRateLimit // make a copy
	globalLimit.DefaultValue = config.Rate{Limit: rate.Limit(cfg.GlobalRPS), Burst: cfg.GlobalBurst}
	global, err := f.MakeRateLimiter(globalLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create global rate limiter: %w", err)
	}
	ownerLimit := cresettings.Default.PerOwner.WorkflowTriggerRateLimit // make a copy
	ownerLimit.DefaultValue = config.Rate{Limit: rate.Limit(cfg.PerSenderRPS), Burst: cfg.PerSenderBurst}
	owner, err := f.MakeRateLimiter(ownerLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner rate limiter: %w", err)
	}
	workflow, err := f.MakeRateLimiter(cresettings.Default.PerWorkflow.TriggerRateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow rate limiter: %w", err)
	}
	return limits.MultiRateLimiter{workflow, owner, global}, nil
}
