// Package heavyweight contains test helpers that are costly and you should
// think **real carefully** before using in your tests.
package heavyweight

import (
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/store"

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
	_, err := db.Exec(store.FixturesSQL())
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

func generateName() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func prepareDB(t testing.TB, withTemplate bool, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (chainlink.GeneralConfig, *sqlx.DB) {
	tests.SkipShort(t, "FullTestDB")

	gcfg := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.DriverName = pgcommon.DriverPostgres
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	require.NoError(t, os.MkdirAll(gcfg.RootDir(), 0700))
	t.Cleanup(func() { os.RemoveAll(gcfg.RootDir()) })

	migrationTestDBURL := testdb.CreateOrReplace(t, gcfg.Database().URL(), generateName(), withTemplate)
	db, err := pg.NewConnection(t.Context(), migrationTestDBURL.String(), pgcommon.DriverPostgres, gcfg.Database())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) }) // must close before dropping

	// reset with new URL
	gcfg = configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.DriverName = pgcommon.DriverPostgres
		s.Database.URL = models.NewSecretURL((*commoncfg.URL)(&migrationTestDBURL))
		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	return gcfg, db
}
