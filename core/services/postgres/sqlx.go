package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	mapper "github.com/scylladb/go-reflectx"
	"gorm.io/gorm"
)

type Queryer interface {
	sqlx.Ext
	sqlx.ExtContext
	QueryRow(query string, args ...interface{}) *sql.Row
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

func UnwrapGorm(db *gorm.DB) Queryer {
	if tx, ok := db.Statement.ConnPool.(*sql.Tx); ok {
		// if a transaction is currently present use that instead
		mapper := reflectx.NewMapperFunc("db", mapper.CamelToSnakeASCII)
		return &sqlx.Tx{Tx: tx, Mapper: mapper}
	} else {
		return UnwrapGormDB(db)
	}
}

func SqlxTransaction(ctx context.Context, q Queryer, fc func(tx *sqlx.Tx) error, txOpts ...sql.TxOptions) (err error) {
	switch q.(type) {
	case *sqlx.Tx:
		// nested transaction: just use the outer transaction
		tx := q.(*sqlx.Tx)
		err = fc(tx)

	case *sqlx.DB:
		db := q.(*sqlx.DB)

		opts := &DefaultSqlTxOptions
		if len(txOpts) > 0 {
			opts = &txOpts[0]
		}

		tx, err := db.BeginTxx(ctx, opts)
		panicked := false

		defer func() {
			// Make sure to rollback when panic, Block error or Commit error
			if panicked || err != nil {
				if perr := tx.Rollback(); perr != nil {
					panic(perr)
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
		errors.Errorf("invalid db type")
	}

	return
}
