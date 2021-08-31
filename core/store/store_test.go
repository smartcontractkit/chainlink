package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/migrate"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_SquashMigrationUpgrade(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrationssquash", false)
	defer cleanup()
	db := orm.DB

	// Latest migrations should work fine.
	static.Version = "0.9.11"
	err := migrate.Migrate(postgres.UnwrapGormDB(db).DB)
	require.NoError(t, err)
	err = store.CheckSquashUpgrade(db)
	require.NoError(t, err)
	static.Version = "unset"
}

func TestStore_UpsertLatestNodeVersion(t *testing.T) {
	t.Parallel()

	t.Run("static version unset", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		t.Cleanup(cleanup)
		verORM := versioning.NewORM(postgres.WrapDbWithSqlx(
			postgres.MustSQLDB(store.DB)),
		)

		// Test node version unset
		ver, err := verORM.FindLatestNodeVersion()
		require.NoError(t, err)
		require.NotNil(t, ver)
		require.Contains(t, ver.Version, "random")
	})

	t.Run("static version set", func(t *testing.T) {
		static.Version = "0.9.11"
		store, cleanup := cltest.NewStore(t)
		t.Cleanup(cleanup)
		verORM := versioning.NewORM(postgres.WrapDbWithSqlx(
			postgres.MustSQLDB(store.DB)),
		)

		ver, err := verORM.FindLatestNodeVersion()
		require.NoError(t, err)
		require.NotNil(t, ver)
		require.Equal(t, "0.9.11", ver.Version)

		static.Version = "unset"
	})
}

func TestStore_Start(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	assert.NoError(t, store.Start())
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	defer cleanup()

	assert.NoError(t, s.Close())
}
