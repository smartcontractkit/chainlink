package pg

import (
	"fmt"
	"time"

	// need to make sure pgx driver is registered before opening connection
	_ "github.com/jackc/pgx/v4/stdlib"
	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type ConnectionConfig interface {
	DatabaseDefaultIdleInTxSessionTimeout() time.Duration
	DatabaseDefaultLockTimeout() time.Duration
	ORMMaxOpenConns() int
	ORMMaxIdleConns() int
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
		uri = uuid.NewV4().String()
	}

	// Initialize sql/sqlx
	db, err = sqlx.Open(string(dialect), uri)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// Set default connection options
	lockTimeout := config.DatabaseDefaultLockTimeout().Milliseconds()
	idleInTxSessionTimeout := config.DatabaseDefaultIdleInTxSessionTimeout().Milliseconds()
	stmt := fmt.Sprintf(`SET TIME ZONE 'UTC'; SET lock_timeout = %d; SET idle_in_transaction_session_timeout = %d; SET default_transaction_isolation = %q`,
		lockTimeout, idleInTxSessionTimeout, DefaultIsolation.String())
	if _, err = db.Exec(stmt); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.ORMMaxOpenConns())
	db.SetMaxIdleConns(config.ORMMaxIdleConns())

	err = disallowReplica(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func disallowReplica(db *sqlx.DB) error {
	rs, err := db.Query("SHOW session_replication_role")
	if err != nil {
		return err
	}

	// there is only one row in the result set so don't bother with a scan loop
	var val string
	if ok := rs.Next(); !ok {
		return fmt.Errorf("unable to determine session_replication_role")
	}
	err = rs.Scan(&val)
	if err != nil {
		return err
	}
	if err := rs.Close(); err != nil {
		return err
	}
	if err := rs.Err(); err != nil {
		return err
	}

	if val == "replica" {
		return fmt.Errorf("Error: invalid `session_replication_role`: %s. Refusing to connect to replica database. Writing to a replica will corrupt the database", val)
	}

	return nil
}
