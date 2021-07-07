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
	null "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	_ "github.com/smartcontractkit/chainlink/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const MIGRATIONS_DIR string = "migrations"

func init() {
	goose.SetBaseFS(embedMigrations)
	goose.SetSequential(true)
	goose.SetTableName("goose_migrations")

	verbose, _ := strconv.ParseBool(os.Getenv("LOG_SQL_MIGRATIONS"))
	goose.SetVerbose(verbose)
}

// Ensure we migrated from v1 migrations to goose_migrations
func ensureMigrated(db *sql.DB) {
	sqlxDB := postgres.WrapDbWithSqlx(db)
	var names []string
	err := sqlxDB.Select(&names, `SELECT id FROM migrations`)
	if err != nil {
		// already migrated
		return
	}

	// ensure a goose migrations table exists with it's initial v0
	if _, err = goose.GetDBVersion(db); err != nil {
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
		panic("Database state is too old. Need to migrate to chainlink version 0.9.10 first before upgrading to this version")
	}

	// insert records for existing migrations
	sql := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, true);`, goose.TableName())
	err = postgres.SqlTransaction(context.Background(), db, func(tx *sqlx.Tx) error {
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

func Migrate(db *sql.DB) error {
	ensureMigrated(db)
	return goose.Up(db, MIGRATIONS_DIR)
}

func Rollback(db *sql.DB, version null.Int) error {
	ensureMigrated(db)
	if version.Valid {
		return goose.DownTo(db, MIGRATIONS_DIR, version.Int64)
	}
	return goose.Down(db, MIGRATIONS_DIR)
}

func Current(db *sql.DB) (int64, error) {
	ensureMigrated(db)
	return goose.EnsureDBVersion(db)
}

func Status(db *sql.DB) error {
	ensureMigrated(db)
	return goose.Status(db, MIGRATIONS_DIR)
}

func Create(db *sql.DB, name, migrationType string) error {
	return goose.Create(db, "core/store/migrate/migrations", name, migrationType)
}
