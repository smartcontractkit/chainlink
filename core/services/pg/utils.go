package pg

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func IsSerializationAnomaly(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(errors.Cause(err).Error(), "could not serialize access due to concurrent update")
}
