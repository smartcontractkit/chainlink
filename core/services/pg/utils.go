package pg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/smartcontractkit/chainlink/core/config/parse"
)

func init() {
	s := os.Getenv("DATABASE_DEFAULT_QUERY_TIMEOUT")
	if s != "" {
		t, err := parse.Duration(s)
		if err != nil {
			panic(fmt.Sprintf("DATABASE_DEFAULT_QUERY_TIMEOUT value of %s is not a valid duration", s))
		}
		DefaultQueryTimeout = t.(time.Duration)
	}
	s = os.Getenv("DATABASE_DEFAULT_LOCK_TIMEOUT")
	if s != "" {
		t, err := parse.Duration(s)
		if err != nil {
			panic(fmt.Sprintf("DATABASE_DEFAULT_LOCK_TIMEOUT value of %s is not a valid duration", s))
		}
		DefaultLockTimeout = t.(time.Duration)
	}
	s = os.Getenv("DATABASE_DEFAULT_IDLE_IN_TX_SESSION_TIMEOUT")
	if s != "" {
		t, err := parse.Duration(s)
		if err != nil {
			panic(fmt.Sprintf("DATABASE_DEFAULT_IDLE_IN_TX_SESSION_TIMEOUT value of %s is not a valid duration", s))
		}
		DefaultIdleInTxSessionTimeout = t.(time.Duration)
	}
}

var (
	// DefaultQueryTimeout is a reasonable upper bound for how long a SQL query should take
	DefaultQueryTimeout = 10 * time.Second
	// DefaultLockTimeout controls the max time we will wait for any kind of database lock.
	// It's good to set this to _something_ because waiting for locks forever is really bad.
	DefaultLockTimeout = 15 * time.Second
	// DefaultIdleInTxSessionTimeout controls the max time we leave a transaction open and idle.
	// It's good to set this to _something_ because leaving transactions open forever is really bad.
	DefaultIdleInTxSessionTimeout = 1 * time.Hour
)

// DefaultQueryCtx returns a context with a sensible sanity limit timeout for SQL queries
func DefaultQueryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultQueryTimeout)
}

// DefaultQueryCtxWithParent returns a context with a sensible sanity limit timeout for
// SQL queries with the given parent context
func DefaultQueryCtxWithParent(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultQueryTimeout)
}
