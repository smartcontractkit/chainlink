package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
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

// MustSQLDB panics if there is an error getting the underlying SQL TX
func MustSQLTx(db *gorm.DB) *sql.Tx {
	sqlTx, ok := db.Statement.ConnPool.(*sql.Tx)
	if !ok {
		panic("db is not inside a tx, open a tx before calling MustSqlTx")
	}
	return sqlTx
}

func MustSQLDB(db *gorm.DB) *sql.DB {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("unable to extract underlying db connection err %v", err))
	}
	return sqlDB
}
