package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib" // need to make sure pgx driver is registered before opening connection
	"github.com/jmoiron/sqlx"

	commonpg "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
)

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

func NewConnection(ctx context.Context, uri string, driverName string, config ConnectionConfig) (db *sqlx.DB, err error) {
	if driverName == commonpg.DriverTxWrappedPostgres {
		if err = sqltest.RegisterTxDB(uri); err != nil {
			return nil, fmt.Errorf("failed to register %s: %w", commonpg.DriverTxWrappedPostgres, err)
		}
	}
	db, err = commonpg.DBConfig{
		IdleInTxSessionTimeout: config.DefaultIdleInTxSessionTimeout(),
		LockTimeout:            config.DefaultLockTimeout(),
		MaxOpenConns:           config.MaxOpenConns(),
		MaxIdleConns:           config.MaxIdleConns(),
	}.New(ctx, uri, driverName)
	if err != nil {
		return nil, err
	}
	setMaxMercuryConns(db, config)

	if os.Getenv("SKIP_PG_VERSION_CHECK") != "true" {
		if err = checkVersion(db, MinRequiredPGVersion); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func setMaxMercuryConns(db *sqlx.DB, config ConnectionConfig) {
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
		const errUndefinedTable = "42P01"
		var pqerr *pgconn.PgError
		if errors.As(err, &pqerr) {
			if pqerr.Code == errUndefinedTable {
				// no mercury jobs defined
				return
			}
		}
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
