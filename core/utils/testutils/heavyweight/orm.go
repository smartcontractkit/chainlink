// Package heavyweight contains test helpers that are costly and you should
// think **real carefully** before using in your tests.
package heavyweight

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

// FullTestDBV2 creates a pristine DB which runs in a separate database than the normal
// unit tests, so you can do things like use other Postgres connection types with it.
func FullTestDBV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return KindFixtures.PrepareDB(t, overrideFn)
}

// FullTestDBNoFixturesV2 is the same as FullTestDB, but it does not load fixtures.
func FullTestDBNoFixturesV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return KindTemplate.PrepareDB(t, overrideFn)
}

// FullTestDBEmptyV2 creates an empty DB (without migrations).
func FullTestDBEmptyV2(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	return KindEmpty.PrepareDB(t, overrideFn)
}

func generateName() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

type Kind int

const (
	KindEmpty Kind = iota
	KindTemplate
	KindFixtures
)

func (c Kind) PrepareDB(t testing.TB, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	testutils.SkipShort(t, "FullTestDB")

	gcfg := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Dialect = dialects.Postgres
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	migrationTestDBURL, err := testdb.CreateOrReplace(gcfg.Database().URL(), generateName(), c != KindEmpty)
	require.NoError(t, err)
	db, err := pg.NewConnection(migrationTestDBURL, dialects.Postgres, gcfg.Database())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close()) // must close before dropping
		require.NoError(t, testdb.Drop(*testutils.MustParseURL(t, migrationTestDBURL)))
		os.RemoveAll(gcfg.RootDir())
	})

	gcfg = configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Dialect = dialects.Postgres
		s.Database.URL = models.MustSecretURL(migrationTestDBURL)
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	if c == KindFixtures {
		_, filename, _, ok := runtime.Caller(1)
		if !ok {
			t.Fatal("could not get runtime.Caller(1)")
		}
		filepath := path.Join(path.Dir(filename), "../../../store/fixtures/fixtures.sql")
		fixturesSQL, err := os.ReadFile(filepath)
		require.NoError(t, err)
		_, err = db.Exec(string(fixturesSQL))
		require.NoError(t, err)
	}

	return gcfg, db
}
