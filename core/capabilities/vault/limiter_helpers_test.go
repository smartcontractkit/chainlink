package vault

import (
	"context"

	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// ownerOverrideLimiter is a BoundLimiter that applies different size limits based on the owner in context.
// This lets tests verify that the validator correctly threads the owner through context to the limiter.
type ownerOverrideLimiter struct {
	defaultBound pkgconfig.Size
	overrides    map[string]pkgconfig.Size
}

func (o *ownerOverrideLimiter) Close() error { return nil }
func (o *ownerOverrideLimiter) Limit(ctx context.Context) (pkgconfig.Size, error) {
	return o.boundFor(ctx), nil
}
func (o *ownerOverrideLimiter) Check(ctx context.Context, n pkgconfig.Size) error {
	bound := o.boundFor(ctx)
	if n > bound {
		return limits.ErrorBoundLimited[pkgconfig.Size]{Limit: bound, Amount: n}
	}
	return nil
}
func (o *ownerOverrideLimiter) boundFor(ctx context.Context) pkgconfig.Size {
	if override, ok := o.overrides[contexts.CREValue(ctx).Owner]; ok {
		return override
	}
	return o.defaultBound
}
