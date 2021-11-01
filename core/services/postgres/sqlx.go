package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx/reflectx"
	mapper "github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"
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
		txx := sqlx.NewTx(tx, db.Dialector.Name())
		txx.Mapper = mapper
		return txx
	}
	return UnwrapGormDB(db)
}

func SqlxTransaction(ctx context.Context, q Queryer, fc func(tx *sqlx.Tx) error, txOpts ...sql.TxOptions) (err error) {
	switch db := q.(type) {
	case *sqlx.Tx:
		// nested transaction: just use the outer transaction
		err = fc(db)
	case *sqlx.DB:
		err = sqlxTransaction(ctx, db, fc, txOpts...)
	default:
		err = errors.Errorf("invalid db type")
	}

	return
}
