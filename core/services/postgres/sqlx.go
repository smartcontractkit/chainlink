package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx/reflectx"
	"github.com/pkg/errors"
	mapper "github.com/scylladb/go-reflectx"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/sqlx"
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

func SqlxTransactionWithDefaultCtx(q Queryer, lggr logger.Logger, fc func(q Queryer) error, txOpts ...TxOptions) (err error) {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	return SqlxTransaction(ctx, q, lggr, fc, txOpts...)
}

func SqlxTransaction(ctx context.Context, q Queryer, lggr logger.Logger, fc func(q Queryer) error, txOpts ...TxOptions) (err error) {
	switch db := q.(type) {
	case *sqlx.Tx:
		// nested transaction: just use the outer transaction
		err = fc(db)
	case *sqlx.DB:
		err = sqlxTransactionQ(ctx, db, lggr, fc, txOpts...)
	default:
		if AllowUnknownQueryerTypeInTransaction {
			err = fc(q)
		} else {
			err = errors.Errorf("invalid db type")
		}
	}

	return
}
