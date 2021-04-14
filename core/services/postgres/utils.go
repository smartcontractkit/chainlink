package postgres

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const DefaultQueryTimeout = 10 * time.Second

// DefaultQueryCtx returns a context with a sensible sanity limit timeout for
// SQL queries
func DefaultQueryCtx() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultQueryTimeout)
}

func IsSerializationAnomaly(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(errors.Cause(err).Error(), "could not serialize access due to concurrent update")
}

var (
	DefaultSqlTxOptions = sql.TxOptions{
		// NOTE: This is the default level in Postgres anyway, we just make it
		// explicit here
		Isolation: sql.LevelReadCommitted,
	}
)
