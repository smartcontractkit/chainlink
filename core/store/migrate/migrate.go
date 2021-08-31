package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	_ "github.com/smartcontractkit/chainlink/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
	"github.com/smartcontractkit/sqlx"

	"github.com/pressly/goose/v3"
	null "gopkg.in/guregu/null.v4"
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
	var count int
	err := db.QueryRow(`SELECT count(*) FROM migrations`).Scan(&count)
	if err != nil {
		// already migrated
		return
	}

	// ensure a goose migrations table exists with it's initial v0
	if _, err = goose.GetDBVersion(db); err != nil {
		panic(err)
	}

	// insert records for existing migrations
	sql := `INSERT INTO %s (version_id, is_applied) VALUES %s;`
	valueStrings := []string{}
	for i := 1; i <= count; i++ {
		valueStrings = append(valueStrings, fmt.Sprintf("(%v, true)", strconv.FormatInt(int64(i), 10)))
	}
	sql = fmt.Sprintf(sql, goose.TableName(), strings.Join(valueStrings, ","))

	err = postgres.SqlTransaction(context.Background(), db, func(tx *sqlx.Tx) error {
		if _, err = db.Exec(sql); err != nil {
			return err
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
