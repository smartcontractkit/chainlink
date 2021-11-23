package pg

import (
	"fmt"

	// need to make sure pgx driver is registered before opening connection
	_ "github.com/jackc/pgx/v4/stdlib"
	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type Config struct {
	Logger       logger.Logger
	MaxOpenConns int
	MaxIdleConns int
}

func NewConnection(uri string, dialect string, config Config) (db *sqlx.DB, err error) {
	if dialect == "txdb" {
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
	db, err = sqlx.Open(dialect, uri)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// Set default connection options
	stmt := fmt.Sprintf(`SET TIME ZONE 'UTC'; SET lock_timeout = %d; SET idle_in_transaction_session_timeout = %d; SET default_transaction_isolation = %q`, DefaultLockTimeout.Milliseconds(), DefaultIdleInTxSessionTimeout.Milliseconds(), DefaultIsolation.String())
	if _, err = db.Exec(stmt); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return db, nil
}
