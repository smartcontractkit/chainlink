package pg

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"strconv"
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
	// LongQueryTimeout is a bigger upper bound for how long a SQL query should take
	LongQueryTimeout = 1 * time.Minute
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

var _ driver.Valuer = Limit(-1)

// Limit is a helper driver.Valuer for LIMIT queries which uses nil/NULL for negative values.
type Limit int

func (l Limit) String() string {
	if l < 0 {
		return "NULL"
	}
	return strconv.Itoa(int(l))
}

func (l Limit) Value() (driver.Value, error) {
	if l < 0 {
		return nil, nil
	}
	return l, nil
}
