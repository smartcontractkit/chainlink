package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/smartcontractkit/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const MIGRATIONS_DIR string = "migrations"

func init() {
	goose.SetBaseFS(embedMigrations)
	goose.SetSequential(true)
	goose.SetTableName("goose_migrations")
	// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
	var logMigrations string
	if v1, v2 := os.Getenv("LOG_SQL_MIGRATIONS"), os.Getenv("CL_LOG_SQL_MIGRATIONS"); v1 != "" && v2 != "" {
		if v1 != v2 {
			panic("you may only set one of LOG_SQL_MIGRATIONS and CL_LOG_SQL_MIGRATIONS environment variables, not both")
		}
	} else if v1 == "" {
		logMigrations = v2
	} else if v2 == "" {
		logMigrations = v1
	}
	verbose, _ := strconv.ParseBool(logMigrations)
	goose.SetVerbose(verbose)
}

// Ensure we migrated from v1 migrations to goose_migrations
func ensureMigrated(db *sql.DB, lggr logger.Logger) {
	sqlxDB := pg.WrapDbWithSqlx(db)
	var names []string
	err := sqlxDB.Select(&names, `SELECT id FROM migrations`)
	if err != nil {
		// already migrated
		return
	}
	err = pg.SqlTransaction(context.Background(), db, lggr, func(tx *sqlx.Tx) error {
		// ensure that no legacy job specs are present: we _must_ bail out early if
		// so because otherwise we run the risk of dropping working jobs if the
		// user has not read the release notes
		return migrations.CheckNoLegacyJobs(tx.Tx)
	})
	if err != nil {
		panic(err)
	}

	// Look for the squashed migration. If not present, the db needs to be migrated on an earlier release first
	found := false
	for _, name := range names {
		if name == "1611847145" {
			found = true
		}
	}
	if !found {
		panic("Database state is too old. Need to migrate to chainlink version 0.9.10 first before upgrading to this version. This upgrade is NOT REVERSIBLE, so it is STRONGLY RECOMMENDED that you take a database backup before continuing.")
	}

	// ensure a goose migrations table exists with it's initial v0
	if _, err = goose.GetDBVersion(db); err != nil {
		panic(err)
	}

	// insert records for existing migrations
	//nolint
	sql := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, true);`, goose.TableName())
	err = pg.SqlTransaction(context.Background(), db, lggr, func(tx *sqlx.Tx) error {
		for _, name := range names {
			var id int64
			// the first migration doesn't follow the naming convention
			if name == "1611847145" {
				id = 1
			} else {
				idx := strings.Index(name, "_")
				if idx < 0 {
					// old migration we don't care about
					continue
				}

				id, err = strconv.ParseInt(name[:idx], 10, 64)
				if err == nil && id <= 0 {
					return errors.New("migration IDs must be greater than zero")
				}
			}

			if _, err = db.Exec(sql, id); err != nil {
				return err
			}
		}

		_, err = db.Exec("DROP TABLE migrations;")
		return err
	})
	if err != nil {
		panic(err)
	}
}

func Migrate(db *sql.DB, lggr logger.Logger) error {
	ensureMigrated(db, lggr)
	// WithAllowMissing is necessary when upgrading from 0.10.14 since it
	// includes out-of-order migrations
	return goose.Up(db, MIGRATIONS_DIR, goose.WithAllowMissing())
}

func Rollback(db *sql.DB, lggr logger.Logger, version null.Int) error {
	ensureMigrated(db, lggr)
	if version.Valid {
		return goose.DownTo(db, MIGRATIONS_DIR, version.Int64)
	}
	return goose.Down(db, MIGRATIONS_DIR)
}

func Current(db *sql.DB, lggr logger.Logger) (int64, error) {
	ensureMigrated(db, lggr)
	return goose.EnsureDBVersion(db)
}

func Status(db *sql.DB, lggr logger.Logger) error {
	ensureMigrated(db, lggr)
	return goose.Status(db, MIGRATIONS_DIR)
}

func Create(db *sql.DB, name, migrationType string) error {
	return goose.Create(db, "core/store/migrate/migrations", name, migrationType)
}
