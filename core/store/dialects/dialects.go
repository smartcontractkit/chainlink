package dialects

// DialectName is a compiler enforced type used that maps to gorm's dialect
// names.
type DialectName string

const (
	// Postgres represents the postgres dialect.
	Postgres DialectName = "pgx"
	// TransactionWrappedPostgres is useful for tests.
	// When the connection is opened, it starts a transaction and all
	// operations performed on the DB will be within that transaction.
	//
	TransactionWrappedPostgres DialectName = "txdb"
	// PostgresWithoutLock represents the postgres dialect but it does not
	// wait for a lock to connect. Intended to be used for read only access, or in tests.
	PostgresWithoutLock DialectName = "postgresWithoutLock"
)
