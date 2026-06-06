package testdb

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // registers pgx driver for pgcommon.DriverPostgres
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/common"
	"github.com/stretchr/testify/require"

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
	require.NotEmpty(t, rawDBURL, "CL_DATABASE_URL environment variable is required")
	dbURL, err := url.Parse(rawDBURL)
	require.NoError(t, err)
	require.NotEmpty(t, dbURL.String(), "CL_DATABASE_URL environment variable is required")

	migrator := migratorConfig(withTemplate)
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

// templateHash is computed once per process from embedded migration files.
// Migrations are compile-time constants (go:embed), so the hash is safe to cache.
var (
	templateHashOnce sync.Once
	templateHash     string
	errTemplateHash  error
)

func (m *migrator) Hash() (string, error) {
	if !m.withTemplate {
		return "empty", nil
	}
	templateHashOnce.Do(func() {
		h1, err := common.HashDirs(migrate.EmbedMigrations, "*.sql", migrate.MigrationsDir)
		if err != nil {
			errTemplateHash = err
			return
		}
		h2, err := common.HashDirs(migrate.EmbedMigrations, "*.go", migrate.MigrationsDir)
		if err != nil {
			errTemplateHash = err
			return
		}
		templateHash = common.NewRecursiveHash(
			common.Field("sql", h1),
			common.Field("go", h2),
		).String()
	})
	return templateHash, errTemplateHash
}

func (m *migrator) Migrate(ctx context.Context, db *sql.DB, config pgtestdb.Config) error {
	if !m.withTemplate {
		return nil
	}
	// Note: We do not call SetMigrationENVVars because it is only strictly needed for goose up-to
	// specific versions or custom seeder data. If needed, you can explicitly configure goose here.
	return migrate.Migrate(ctx, db)
}

func migratorConfig(withTemplate bool) pgtestdb.Migrator {
	return &migrator{withTemplate: withTemplate}
}
