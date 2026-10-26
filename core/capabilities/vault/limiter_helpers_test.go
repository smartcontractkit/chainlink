package vault

import (
	"context"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// ownerOverrideLimiter is a BoundLimiter that applies different size limits based on the owner in context.
// This lets tests verify that the validator correctly threads the owner through context to the limiter.
type ownerOverrideLimiter struct {
	defaultBound commonconfig.Size
	overrides    map[string]commonconfig.Size
}

func (o *ownerOverrideLimiter) Close() error { return nil }
func (o *ownerOverrideLimiter) Limit(ctx context.Context) (commonconfig.Size, error) {
	return o.boundFor(ctx), nil
}

func (o *ownerOverrideLimiter) Check(ctx context.Context, n commonconfig.Size) error {
	bound := o.boundFor(ctx)
	if n > bound {
		return limits.ErrorBoundLimited[commonconfig.Size]{Limit: bound, Amount: n}
	}
	return nil
}

func (o *ownerOverrideLimiter) boundFor(ctx context.Context) commonconfig.Size {
	if override, ok := o.overrides[contexts.CREValue(ctx).Owner]; ok {
		return override
	}
	return o.defaultBound
}
