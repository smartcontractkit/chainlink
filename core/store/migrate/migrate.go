package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
)

//go:embed migrations/*.sql migrations/*.go
var embedMigrations embed.FS

const MIGRATIONS_DIR string = "migrations"

func NewProvider(ctx context.Context, db *sql.DB) (*goose.Provider, error) {
	store, err := database.NewStore(goose.DialectPostgres, "goose_migrations")
	if err != nil {
		return nil, err
	}

	goMigrations := []*goose.Migration{
		migrations.Migration36,
		migrations.Migration54,
		migrations.Migration56,
		migrations.Migration195,
	}

	logMigrations := os.Getenv("CL_LOG_SQL_MIGRATIONS")
	verbose, _ := strconv.ParseBool(logMigrations)

	fys, err := fs.Sub(embedMigrations, MIGRATIONS_DIR)
	if err != nil {
		return nil, fmt.Errorf("failed to get sub filesystem for embedded migration dir: %w", err)
	}
	// hack to work around global go migrations
	// https: //github.com/pressly/goose/issues/782
	goose.ResetGlobalMigrations()
	p, err := goose.NewProvider("", db, fys,
		goose.WithStore(store),
		goose.WithGoMigrations(goMigrations...),
		goose.WithVerbose(verbose))
	if err != nil {
		return nil, fmt.Errorf("failed to create goose provider: %w", err)
	}

	err = ensureMigrated(ctx, db, p, store.Tablename())
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Ensure we migrated from v1 migrations to goose_migrations
// TODO remove this for v3
func ensureMigrated(ctx context.Context, db *sql.DB, p *goose.Provider, providerTableName string) error {
	todo, err := p.HasPending(ctx)
	if !todo && err == nil {
		return nil
	}
	sqlxDB := pg.WrapDbWithSqlx(db)
	var names []string
	err = sqlxDB.SelectContext(ctx, &names, `SELECT id FROM migrations`)
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
	if _, err = p.GetDBVersion(ctx); err != nil {
		return err
	}

	// insert records for existing migrations
	sql := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, true);`, providerTableName)
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

func Migrate(ctx context.Context, db *sql.DB) error {
	provider, err := NewProvider(ctx, db)
	if err != nil {
		return err
	}
	_, err = provider.Up(ctx)
	return err
}

func Rollback(ctx context.Context, db *sql.DB, version null.Int) error {
	provider, err := NewProvider(ctx, db)
	if err != nil {
		return err
	}
	if version.Valid {
		_, err = provider.DownTo(ctx, version.Int64)
	} else {
		_, err = provider.Down(ctx)
	}
	return err
}

func Current(ctx context.Context, db *sql.DB) (int64, error) {
	provider, err := NewProvider(ctx, db)
	if err != nil {
		return -1, err
	}
	return provider.GetDBVersion(ctx)
}

func Status(ctx context.Context, db *sql.DB) error {
	provider, err := NewProvider(ctx, db)
	if err != nil {
		return err
	}
	migrations, err := provider.Status(ctx)
	if err != nil {
		return err
	}
	for _, m := range migrations {
		fmt.Printf("version:%d, path:%s, type:%s, state:%s, appliedAt: %s \n", m.Source.Version, m.Source.Path, m.Source.Type, m.State, m.AppliedAt.String())
	}
	return nil
}

func Create(db *sql.DB, name, migrationType string) error {
	return goose.Create(db, "core/store/migrate/migrations", name, migrationType)
}

// SetMigrationENVVars is used to inject values from config to goose migrations via env.
func SetMigrationENVVars(generalConfig toml.EVMConfigs) error {
	if generalConfig.Enabled() {
		err := os.Setenv(env.EVMChainIDNotNullMigration0195, generalConfig[0].ChainID.String())
		if err != nil {
			panic(pkgerrors.Wrap(err, "failed to set migrations env variables"))
		}
	}
	return nil
}
