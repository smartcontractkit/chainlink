package pg

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	mapper "github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//go:generate mockery --quiet --name Queryer --output ./mocks/ --case=underscore
type Queryer interface {
	sqlx.Ext
	sqlx.ExtContext
	sqlx.Preparer
	sqlx.PreparerContext
	sqlx.Queryer
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

func WrapDbWithSqlx(rdb *sql.DB) *sqlx.DB {
	db := sqlx.NewDb(rdb, "postgres")
	db.MapperFunc(mapper.CamelToSnakeASCII)
	return db
}

func SqlxTransaction(ctx context.Context, q Queryer, lggr logger.Logger, fc func(q Queryer) error, txOpts ...TxOptions) (err error) {
	switch db := q.(type) {
	case *sqlx.Tx:
		// nested transaction: just use the outer transaction
		err = fc(db)
	case *sqlx.DB:
		err = sqlxTransactionQ(ctx, db, lggr, fc, txOpts...)
	case Q:
		err = sqlxTransactionQ(ctx, db.db, lggr, fc, txOpts...)
	default:
		err = errors.Errorf("invalid db type: %T", q)
	}

	return
}
