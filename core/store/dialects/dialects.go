package dialects

import (
	// need to make sure pgx driver is registered before opening connection
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DialectName is a compiler enforced type used that maps to database dialect names
type DialectName string

const (
	// Postgres represents the postgres dialect.
	Postgres DialectName = "pgx"
	// TransactionWrappedPostgres is useful for tests.
	// When the connection is opened, it starts a transaction and all
	// operations performed on the DB will be within that transaction.
	TransactionWrappedPostgres DialectName = "txdb"
)
