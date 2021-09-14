package postgres

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/config"

	"github.com/smartcontractkit/sqlx"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(uri string, dialect string, cfg config.GeneralConfig, shutdownSignal gracefulpanic.Signal) (db *sqlx.DB, gormDB *gorm.DB, err error) {
	originalURI := uri
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
	// NOTE: SetConsumerName was already called in config.DatabaseURL(), we don't need to do it here

	newLogger := NewLogWrapper(logger.Default, cfg.LogSQLStatements(), time.Second)

	// Initialize sql/sqlx
	db, err = sqlx.Open(dialect, uri)
	if err != nil {
		return nil, nil, err
	}
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// Initialize gorm (using the same connection)
	gormDB, err = gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: db.DB,
		DSN:  originalURI,
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to open %s for gorm DB conn %v", uri, db)
	}

	// Set connection options
	if _, err = db.Exec(`SET TIME ZONE 'UTC'`); err != nil {
		return nil, nil, err
	}
	db.SetMaxOpenConns(cfg.ORMMaxOpenConns())
	db.SetMaxIdleConns(cfg.ORMMaxIdleConns())

	return db, gormDB, nil
}