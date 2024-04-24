package testutils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

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

func NewSqlxDB(t testing.TB) *sqlx.DB {
	SkipShortDB(t)
	db, err := sqlx.Open(string(TransactionWrappedPostgres), uuid.New().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}

// SkipShortDB skips tb during -short runs, and notes the DB dependency.
func SkipShortDB(tb testing.TB) {
	SkipShort(tb, "DB dependency")
}

// SkipShort skips tb during -short runs, and notes why.
func SkipShort(tb testing.TB, why string) {
	if testing.Short() {
		tb.Skipf("skipping: %s", why)
	}
}
