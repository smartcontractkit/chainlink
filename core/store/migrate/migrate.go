package migrate

import (
	"database/sql"
	"embed"

	// Invoke init() functions within migrations pkg.
	_ "github.com/smartcontractkit/chainlink/core/store/migrate/migrations"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const MIGRATIONS_DIR string = "migrations"

// migrate existing migrations ->
// id version_id is_applied tstamp
// 1  0          t          ...
// 2  1          t          ...

func init() {
	goose.SetBaseFS(embedMigrations)
	goose.SetSequential(true)
	// SetTableName
	// goose.SetVerbose(true) can be set to debug migrations
}

func Migrate(db *sql.DB) error {
	return goose.Up(db, MIGRATIONS_DIR)
}

func Current(db *sql.DB) (int64, error) {
	return goose.EnsureDBVersion(db)
}

func Status(db *sql.DB) error {
	return goose.Status(db, MIGRATIONS_DIR)
}

// func MigrateDown(db *sql.DB) error {
// 	return goose.Down(db, MIGRATIONS_DIR)
// }
