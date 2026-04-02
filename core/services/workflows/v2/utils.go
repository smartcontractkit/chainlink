package v2

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func gateEnabledOrError(ctx context.Context, gate limits.GateLimiter) (bool, error) {
	if gate == nil {
		return false, errors.New("gate is nil")
	}

	return gate.Limit(ctx)
}
