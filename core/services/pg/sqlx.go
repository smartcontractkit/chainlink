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
	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return db
}
