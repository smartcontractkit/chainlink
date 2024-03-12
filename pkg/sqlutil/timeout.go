package sqlutil

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// TimeoutHook returns a [QueryHook] which adds the defaultTimeout to each context.Context,
// unless [WithoutDefaultTimeout] has been applied to bypass intentionally.
func TimeoutHook(defaultTimeout time.Duration) QueryHook {
	return func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
		if wo := ctx.Value(ctxKeyWithoutDefaultTimeout{}); wo == nil {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
			defer cancel()
		}

		return do(ctx)
	}
}

type ctxKeyWithoutDefaultTimeout struct{}

// WithoutDefaultTimeout makes a [context.Context] exempt from the default timeout normally applied by a [TimeoutHook].
func WithoutDefaultTimeout(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyWithoutDefaultTimeout{}, struct{}{})
}
