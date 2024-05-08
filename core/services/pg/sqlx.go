package pg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
)

// Queryer is deprecated. Use sqlutil.DataSource instead
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
	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return db
}
