package ratelimiter

import (
	"errors"
	"fmt"

	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type Config struct {
	GlobalRPS      float64 `json:"globalRPS"`
	GlobalBurst    int     `json:"globalBurst"`
	PerSenderRPS   float64 `json:"perSenderRPS"`
	PerSenderBurst int     `json:"perSenderBurst"`
}

func NewRateLimiter(config Config, f limits.Factory) (limits.RateLimiter, error) {
	if config.GlobalRPS <= 0.0 || config.PerSenderRPS <= 0.0 {
		return nil, errors.New("RPS values must be positive")
	}
	if config.GlobalBurst <= 0 || config.PerSenderBurst <= 0 {
		return nil, errors.New("burst values must be positive")
	}

	// TODO cresettings
	globalRate := settings.Rate(rate.Limit(config.GlobalRPS), config.GlobalBurst)
	globalRate.Scope = settings.ScopeGlobal
	global, err := f.NewRateLimiter(globalRate)
	if err != nil {
		return nil, fmt.Errorf("failed to create global rate limiter: %w", err)
	}
	ownerRate := settings.Rate(rate.Limit(config.PerSenderRPS), config.PerSenderBurst)
	ownerRate.Scope = settings.ScopeOwner
	owner, err := f.NewRateLimiter(ownerRate)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner rate limiter: %w", err)
	}
	return limits.MultiRateLimiter{owner, global}, nil
}
