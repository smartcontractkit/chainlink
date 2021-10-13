package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx/reflectx"
	mapper "github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"
	"gorm.io/gorm"
)

// AllowUnknownQueryerTypeInTransaction can be set by tests to allow a mock to be passed as a Queryer
var AllowUnknownQueryerTypeInTransaction bool

//go:generate mockery --name Queryer --output ./mocks/ --case=underscore
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

func UnwrapGormDB(db *gorm.DB) *sqlx.DB {
	d, err := db.DB()
	if err != nil {
		panic(err)
	}
	return WrapDbWithSqlx(d)
}

func TryUnwrapGormDB(db *gorm.DB) *sqlx.DB {
	if db == nil {
		return nil
	}
	return UnwrapGormDB(db)
}

func UnwrapGorm(db *gorm.DB) Queryer {
	if tx, ok := db.Statement.ConnPool.(*sql.Tx); ok {
		// if a transaction is currently present use that instead
		mapper := reflectx.NewMapperFunc("db", mapper.CamelToSnakeASCII)
		txx := sqlx.NewTx(tx, db.Dialector.Name())
		txx.Mapper = mapper
		return txx
	}
	return UnwrapGormDB(db)
}

func SqlxTransactionWithDefaultCtx(q Queryer, fc func(q Queryer) error, txOpts ...sql.TxOptions) (err error) {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	return SqlxTransaction(ctx, q, fc, txOpts...)
}

func SqlxTransaction(ctx context.Context, q Queryer, fc func(q Queryer) error, txOpts ...sql.TxOptions) (err error) {
	switch db := q.(type) {
	case *sqlx.Tx:
		// nested transaction: just use the outer transaction
		err = fc(db)

	case *sqlx.DB:
		opts := &DefaultSqlTxOptions
		if len(txOpts) > 0 {
			opts = &txOpts[0]
		}

		var tx *sqlx.Tx
		tx, err = db.BeginTxx(ctx, opts)
		if err != nil {
			return errors.Wrap(err, "failed to begin tx")
		}
		panicked := false

		defer func() {
			// Attempt to rollback when panic, Block error or Commit error
			if panicked || err != nil {
				if perr := tx.Rollback(); perr != nil {
					err = errors.Wrapf(perr, "additional error encountered attempting to rollback transaction on panic or error: panicked: %v, original err: %v ", panicked, err)
				}
			}
		}()

		_, err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds()))
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}

		panicked = true
		err = fc(tx)
		panicked = false

		if err == nil {
			err = errors.WithStack(tx.Commit())
		}
	default:
		if AllowUnknownQueryerTypeInTransaction {
			err = fc(q)
		} else {
			err = errors.Errorf("invalid db type")
		}
	}

	return
}
