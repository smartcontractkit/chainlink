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
	Select(dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Get(dest any, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExec(query string, arg any) (sql.Result, error)
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
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
