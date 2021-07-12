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
	// HACK: This must be the string 'cloudsqlpostgres' because of an absolutely
	// horrible design in gorm. We need gorm to enable postgres-specific
	// features for the txdb driver, but it can only do that if the dialect is
	// called "postgres" or "cloudsqlpostgres".
	//
	// Since "postgres" is already taken, "cloudsqlpostgres" is our only
	// remaining option
	//
	// See: https://github.com/jinzhu/gorm/blob/master/dialect_postgres.go#L15
	TransactionWrappedPostgres DialectName = "cloudsqlpostgres"
	// PostgresWithoutLock represents the postgres dialect but it does not
	// wait for a lock to connect. Intended to be used for read only access, or in tests.
	PostgresWithoutLock DialectName = "postgresWithoutLock"
)
