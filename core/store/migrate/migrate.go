package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
)

//go:embed migrations/*.sql migrations/*.go
var embedMigrations embed.FS

// go:embed relayers/**/*.tmpl.sql
var embedRelayerMigrations embed.FS

const RELAYER_TEMPLATE_DIR string = "relayers"

const MIGRATIONS_DIR string = "migrations"

func init() {
	setupCoreMigrations()
	logMigrations := os.Getenv("CL_LOG_SQL_MIGRATIONS")
	verbose, _ := strconv.ParseBool(logMigrations)
	goose.SetVerbose(verbose)
}

func setupCoreMigrations() {
	goose.SetBaseFS(embedMigrations)
	goose.SetSequential(true)
	goose.SetTableName("goose_migrations")
}

// Ensure we migrated from v1 migrations to goose_migrations
func ensureMigrated(ctx context.Context, db *sql.DB) error {
	sqlxDB := pg.WrapDbWithSqlx(db)
	var names []string
	err := sqlxDB.SelectContext(ctx, &names, `SELECT id FROM migrations`)
	if err != nil {
		// already migrated
		return nil
	}
	// ensure that no legacy job specs are present: we _must_ bail out early if
	// so because otherwise we run the risk of dropping working jobs if the
	// user has not read the release notes
	err = migrations.CheckNoLegacyJobs(ctx, db)
	if err != nil {
		return err
	}

	// Look for the squashed migration. If not present, the db needs to be migrated on an earlier release first
	found := false
	for _, name := range names {
		if name == "1611847145" {
			found = true
		}
	}
	if !found {
		return pkgerrors.New("database state is too old. Need to migrate to chainlink version 0.9.10 first before upgrading to this version. This upgrade is NOT REVERSIBLE, so it is STRONGLY RECOMMENDED that you take a database backup before continuing")
	}

	// ensure a goose migrations table exists with it's initial v0
	if _, err = goose.GetDBVersionContext(ctx, db); err != nil {
		return err
	}

	// insert records for existing migrations
	//nolint
	sql := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, true);`, goose.TableName())
	return sqlutil.TransactDataSource(ctx, sqlxDB, nil, func(tx sqlutil.DataSource) error {
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
					return pkgerrors.New("migration IDs must be greater than zero")
				}
			}

			if _, err = tx.ExecContext(ctx, sql, id); err != nil {
				return err
			}
		}

		_, err = tx.ExecContext(ctx, "DROP TABLE migrations;")
		return err
	})
}

type OptionalMigration struct {
	Type     string
	Template string
	Schema   string
}

func (o OptionalMigration) validate() error {
	if o.Type != "relayer" {
		return fmt.Errorf("unknown migration type: %s", o.Type)
	}
	if o.Template != "evm" {
		return fmt.Errorf("unknown migration template: %s", o.Template)
	}
	return nil
}

func (o OptionalMigration) rootDir() string {
	return fmt.Sprintf("%s/%s", RELAYER_TEMPLATE_DIR, o.Template)
}

func Migrate(ctx context.Context, db *sql.DB, opts ...OptionalMigration) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	// WithAllowMissing is necessary when upgrading from 0.10.14 since it
	// includes out-of-order migrations
	err := goose.Up(db, MIGRATIONS_DIR, goose.WithAllowMissing())
	if err != nil {
		return fmt.Errorf("failed to do core database migration: %w", err)
	}
	// Reset to default OS filesystem so that goose can find the generated migrations
	if len(opts) > 0 {
		defer goose.SetBaseFS(embedMigrations)
		goose.SetBaseFS(nil)
		goose.ResetGlobalMigrations()
	}
	err = ensureMigrated(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to ensure migration: %w", err)
	}
	c, err := Current(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}
	fmt.Printf("Current migration version: %d\n", c)
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)
	for _, opt := range opts {
		if err := opt.validate(); err != nil {
			return fmt.Errorf("invalid migration option: %w", err)
		}
		root := opt.rootDir()
		d := filepath.Join(tmpDir, opt.Template, opt.Schema)
		migrations, err := generateMigrations(root, d, RelayerDB{Schema: opt.Schema})
		if err != nil {
			return fmt.Errorf("failed to generate migrations for opt %v: %w", opt, err)
		}
		fmt.Printf("Generated migrations: %v\n", migrations)

		err = goose.Up(db, d, goose.WithNoVersioning())
		if err != nil {
			return fmt.Errorf("failed to do %s database migration: %w", opt.Type, err)
		}
	}
	return nil
}

// MigrateSubsytstem migrates a subsystem of the chainlink database.
// It generates migrations based on the template for the subsystem and applies them to the database.
func MigrateSubsytstem(ctx context.Context, db *sql.DB, opt OptionalMigration) error {

	if err := opt.validate(); err != nil {
		return fmt.Errorf("invalid migration option: %w", err)
	}

	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	defer setupCoreMigrations()
	goose.SetBaseFS(nil)
	goose.ResetGlobalMigrations()
	goose.SetTableName(fmt.Sprintf("goose_migration_%s_%s", opt.Template, opt.Schema))

	root := opt.rootDir()
	d := filepath.Join(tmpDir, opt.Template, opt.Schema)
	migrations, err := generateMigrations(root, d, RelayerDB{Schema: opt.Schema})
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", opt, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)

	err = goose.Up(db, d, goose.WithNoVersioning())
	if err != nil {
		return fmt.Errorf("failed to do %s database migration: %w", opt.Type, err)
	}

	return nil
}

func Rollback(ctx context.Context, db *sql.DB, version null.Int) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	if version.Valid {
		return goose.DownTo(db, MIGRATIONS_DIR, version.Int64)
	}
	return goose.Down(db, MIGRATIONS_DIR)
}

func RollbackSubsystem(ctx context.Context, db *sql.DB, version null.Int, opt OptionalMigration) error {
	if err := opt.validate(); err != nil {
		return fmt.Errorf("invalid migration option: %w", err)
	}

	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	defer setupCoreMigrations()
	goose.SetBaseFS(nil)
	goose.ResetGlobalMigrations()
	goose.SetTableName(fmt.Sprintf("goose_migration_%s_%s", opt.Template, opt.Schema))

	// TODO: should these be saved somewhere? if so where, if not if the db itself?
	root := opt.rootDir()
	d := filepath.Join(tmpDir, opt.Template, opt.Schema)
	migrations, err := generateMigrations(root, d, RelayerDB{Schema: opt.Schema})
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", opt, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)

	if version.Valid {
		return goose.DownTo(db, d, version.Int64)
	}
	return goose.Down(db, d)
}

func Current(ctx context.Context, db *sql.DB) (int64, error) {
	if err := ensureMigrated(ctx, db); err != nil {
		return -1, err
	}
	return goose.EnsureDBVersion(db)
}

func Status(ctx context.Context, db *sql.DB) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	return goose.Status(db, MIGRATIONS_DIR)
}

func Create(db *sql.DB, name, migrationType string) error {
	return goose.Create(db, "core/store/migrate/migrations", name, migrationType)
}

// SetMigrationENVVars is used to inject values from config to goose migrations via env.
func SetMigrationENVVars(generalConfig chainlink.GeneralConfig) error {
	if generalConfig.EVMEnabled() {
		err := os.Setenv(env.EVMChainIDNotNullMigration0195, generalConfig.EVMConfigs()[0].ChainID.String())
		if err != nil {
			panic(pkgerrors.Wrap(err, "failed to set migrations env variables"))
		}
	}
	return nil
}
