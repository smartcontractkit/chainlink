package pg

import (
	"context"
	"time"
)

const DefaultQueryTimeout = 10 * time.Second

// DefaultQueryCtx returns a context with a sensible sanity limit timeout for SQL queries
func DefaultQueryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultQueryTimeout)
}

// DefaultQueryCtxWithParent returns a context with a sensible sanity limit timeout for
// SQL queries with the given parent context
func DefaultQueryCtxWithParent(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultQueryTimeout)
}
