package pg

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // need to make sure pgx driver is registered before opening connection
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type ConnectionConfig interface {
	DefaultIdleInTxSessionTimeout() time.Duration
	DefaultLockTimeout() time.Duration
	MaxOpenConns() int
	MaxIdleConns() int
}

func NewConnection(uri string, dialect dialects.DialectName, config ConnectionConfig) (db *sqlx.DB, err error) {
	if dialect == dialects.TransactionWrappedPostgres {
		// Dbtx uses the uri as a unique identifier for each transaction. Each ORM
		// should be encapsulated in it's own transaction, and thus needs its own
		// unique id.
		//
		// We can happily throw away the original uri here because if we are using
		// txdb it should have already been set at the point where we called
		// txdb.Register
		uri = uuid.New().String()
	}

	// Initialize sql/sqlx
	db, err = sqlx.Open(string(dialect), uri)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// Set default connection options
	lockTimeout := config.DefaultLockTimeout().Milliseconds()
	idleInTxSessionTimeout := config.DefaultIdleInTxSessionTimeout().Milliseconds()
	stmt := fmt.Sprintf(`SET TIME ZONE 'UTC'; SET lock_timeout = %d; SET idle_in_transaction_session_timeout = %d; SET default_transaction_isolation = %q`,
		lockTimeout, idleInTxSessionTimeout, DefaultIsolation.String())
	if _, err = db.Exec(stmt); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns())
	db.SetMaxIdleConns(config.MaxIdleConns())

	return db, disallowReplica(db)
}

func disallowReplica(db *sqlx.DB) error {
	var val string
	err := db.Get(&val, "SHOW session_replication_role")
	if err != nil {
		return err
	}

	if val == "replica" {
		return fmt.Errorf("invalid `session_replication_role`: %s. Refusing to connect to replica database. Writing to a replica will corrupt the database", val)
	}

	return nil
}
