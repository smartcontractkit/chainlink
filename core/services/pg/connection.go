package pg

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // need to make sure pgx driver is registered before opening connection
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"

	"github.com/XSAM/otelsql"
)

// NOTE: This is the default level in Postgres anyway, we just make it
// explicit here
const defaultIsolation = sql.LevelReadCommitted

var MinRequiredPGVersion = 110000

func init() {
	// from: https://www.postgresql.org/support/versioning/
	now := time.Now()
	if now.Year() > 2023 {
		MinRequiredPGVersion = 120000
	} else if now.Year() > 2024 {
		MinRequiredPGVersion = 130000
	} else if now.Year() > 2025 {
		MinRequiredPGVersion = 140000
	} else if now.Year() > 2026 {
		MinRequiredPGVersion = 150000
	} else if now.Year() > 2027 {
		MinRequiredPGVersion = 160000
	}
}

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
	sqldb, err := otelsql.Open(string(dialect), uri,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithTracerProvider(otel.GetTracerProvider()),
		otelsql.WithSQLCommenter(true),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
			OmitConnPrepare:      true,
			OmitRows:             true,
			OmitConnectorConnect: true,
			OmitConnQuery:        false,
		}),
	)
	if err != nil {
		return nil, err
	}
	db = sqlx.NewDb(sqldb, string(dialect))
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// Set default connection options
	lockTimeout := config.DefaultLockTimeout().Milliseconds()
	idleInTxSessionTimeout := config.DefaultIdleInTxSessionTimeout().Milliseconds()
	stmt := fmt.Sprintf(`SET TIME ZONE 'UTC'; SET lock_timeout = %d; SET idle_in_transaction_session_timeout = %d; SET default_transaction_isolation = %q`,
		lockTimeout, idleInTxSessionTimeout, defaultIsolation.String())
	if _, err = db.Exec(stmt); err != nil {
		return nil, err
	}
	setMaxConns(db, config)

	if os.Getenv("SKIP_PG_VERSION_CHECK") != "true" {
		if err := checkVersion(db, MinRequiredPGVersion); err != nil {
			return nil, err
		}
	}

	return db, disallowReplica(db)
}

func setMaxConns(db *sqlx.DB, config ConnectionConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns())
	db.SetMaxIdleConns(config.MaxIdleConns())

	// HACK: In the case of mercury jobs, one conn is needed per job for good
	// performance. Most nops will forget to increase the defaults to account
	// for this so we detect it here instead.
	//
	// This problem will be solved by replacing mercury with parallel
	// compositions (llo plugin).
	//
	// See: https://smartcontract-it.atlassian.net/browse/MERC-3654
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM ocr2_oracle_specs WHERE plugin_type = 'mercury'`); err != nil {
		log.Printf("Error checking mercury jobs: %s", err.Error())
		return
	}
	if cnt > config.MaxOpenConns() {
		log.Printf("Detected %d mercury jobs, increasing max open connections from %d to %d", cnt, config.MaxOpenConns(), cnt)
		db.SetMaxOpenConns(cnt)
	}
	if cnt > config.MaxIdleConns() {
		log.Printf("Detected %d mercury jobs, increasing max idle connections from %d to %d", cnt, config.MaxIdleConns(), cnt)
		db.SetMaxIdleConns(cnt)
	}
}

type Getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

func checkVersion(db Getter, minVersion int) error {
	var version int
	if err := db.Get(&version, "SHOW server_version_num"); err != nil {
		log.Printf("Error getting server version, skipping Postgres version check: %s", err.Error())
		return nil
	}
	if version < 10000 {
		log.Printf("Unexpectedly small version, skipping Postgres version check (you are running: %d)", version)
		return nil
	}
	if version < minVersion {
		return fmt.Errorf("The minimum required Postgres server version is %d, you are running: %d, which is EOL (see: https://www.postgresql.org/support/versioning/). It is recommended to upgrade your Postgres server. To forcibly override this check, set SKIP_PG_VERSION_CHECK=true", minVersion/10000, version/10000)
	}
	return nil
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
