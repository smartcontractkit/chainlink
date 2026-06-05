package testdb

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"testing"

	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/common"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
)

// New provisions a new isolated test database (via pgtestdb) and returns the connected URL.
// The returned URL will have the provisioned DB role credentials (not superuser), so tests
// run exactly as they would in a real deployment environment without excess privileges.
func New(t testing.TB, withTemplate bool) *url.URL {
	t.Helper()

	rawDBURL := string(env.DatabaseURL.Get())
	if rawDBURL == "" {
		t.Fatalf("you must provide a CL_DATABASE_URL environment variable")
	}

	dbURL, err := url.Parse(rawDBURL)
	require.NoError(t, err)

	migrator := Migrator(withTemplate)
	conf := pgtestdb.Config{
		DriverName:                pgcommon.DriverPostgres,
		User:                      dbURL.User.Username(),
		Host:                      dbURL.Hostname(),
		Port:                      dbURL.Port(),
		Database:                  strings.TrimLeft(dbURL.Path, "/"),
		Options:                   dbURL.RawQuery,
		ForceTerminateConnections: true,
	}
	if pass, ok := dbURL.User.Password(); ok {
		conf.Password = pass
	}

	newConf := pgtestdb.Custom(t, conf, migrator)
	newURLStr := newConf.URL()
	newURL, err := url.Parse(newURLStr)
	require.NoError(t, err)
	return newURL
}

type migrator struct {
	withTemplate bool
}

func (m *migrator) Hash() (string, error) {
	if !m.withTemplate {
		return "empty", nil
	}
	// Hash the embedded migration files so template databases rebuild if schemas change.
	return common.HashDirs(migrate.EmbedMigrations, "*.sql", migrate.MigrationsDir)
}

func (m *migrator) Migrate(ctx context.Context, db *sql.DB, config pgtestdb.Config) error {
	if !m.withTemplate {
		return nil
	}
	// Note: We do not call SetMigrationENVVars because it is only strictly needed for goose up-to
	// specific versions or custom seeder data. If needed, you can explicitly configure goose here.
	return migrate.Migrate(ctx, db)
}

func Migrator(withTemplate bool) pgtestdb.Migrator {
	return &migrator{withTemplate: withTemplate}
}
