// Package heavyweight contains test helpers that are costly and you should
// think **real carefully** before using in your tests.
package heavyweight

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/peterldowns/pgtestdb"
	"github.com/stretchr/testify/require"

	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/store"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

// FullTestDBV2 creates a pristine DB which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types with it.
func FullTestDBV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	cfg, db := FullTestDBNoFixturesV2(t, overrideFn)
	ctx := t.Context()
	_, err := db.ExecContext(ctx, store.FixturesSQL())
	require.NoError(t, err)
	return cfg, db
}

// FullTestDBNoFixturesV2 is the same as FullTestDB, but it does not load fixtures.
func FullTestDBNoFixturesV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return prepareDB(t, true, overrideFn)
}

// FullTestDBEmptyV2 creates an empty DB (without migrations).
func FullTestDBEmptyV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return prepareDB(t, false, overrideFn)
}

func prepareDB(t testing.TB, withTemplate bool, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	tests.SkipShort(t, "FullTestDB")

	// Read env.DatabaseURL directly to get the base connection
	rawDBURL := string(env.DatabaseURL.Get())
	if rawDBURL == "" {
		t.Fatalf("you must provide a CL_DATABASE_URL environment variable")
	}

	dbURL, err := url.Parse(rawDBURL)
	require.NoError(t, err)

	migrator := testdb.Migrator(withTemplate)
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

	migrationTestDBURL := *dbURL
	migrationTestDBURL.Path = "/" + newConf.Database
	dbStr := migrationTestDBURL.String()

	gcfg := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.DriverName = pgcommon.DriverPostgres
		s.Database.URL = models.NewSecretURL((*commoncfg.URL)(&migrationTestDBURL))
		// Explicitly allow simple passwords since tests use `postgres` password
		s.Database.AllowSimplePasswords = new(true)
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	t.Cleanup(func() { os.RemoveAll(gcfg.RootDir()) })

	db, err := pg.NewConnection(t.Context(), dbStr, pgcommon.DriverPostgres, gcfg.Database())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) }) // must close before dropping

	return gcfg, db
}
