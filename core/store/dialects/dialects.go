package dialects

import (
	// need to make sure pgx driver is registered before opening connection
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
)

// DialectName is a compiler enforced type used that maps to database dialect names
type DialectName = pg.Driver

const (
	// Postgres represents the postgres dialect.
	Postgres DialectName = pg.DriverPostgres
	// TransactionWrappedPostgres is useful for tests.
	// When the connection is opened, it starts a transaction and all
	// operations performed on the DB will be within that transaction.
	TransactionWrappedPostgres DialectName = pg.DriverTxWrappedPostgres
)
