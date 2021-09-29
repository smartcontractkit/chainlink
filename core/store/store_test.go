package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_UpsertLatestNodeVersion(t *testing.T) {
	t.Parallel()

	t.Run("static version unset", func(t *testing.T) {
		store := cltest.NewStore(t)
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
		store := cltest.NewStore(t)
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
	store := cltest.NewStore(t)
	assert.NoError(t, store.Start())
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s := cltest.NewStore(t)

	assert.NoError(t, s.Close())
}
