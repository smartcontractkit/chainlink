package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/plugins/relayer/evm"
)

func TestMigrate(t *testing.T) {

	t.Run("core migration with optional relayer migration", func(t *testing.T) {
		_, db := heavyweight.FullTestDBEmptyV2(t, nil)

		ctx := testutils.Context(t)
		cfg := evm.Cfg{
			Schema:  "evm_42",
			ChainID: 42,
		}
		// the evm migrations only work if the core migrations have been run
		// because we are moving existing tables
		err := evm.Migrate(ctx, db.DB, cfg)
		require.Error(t, err)
		err = migrate.Migrate(ctx, db.DB)
		require.NoError(t, err)

		err = evm.Migrate(ctx, db.DB, cfg)
		require.NoError(t, err)

		v2, err := evm.Current(ctx, db.DB, cfg)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, int64(2), v2)

		err = evm.Rollback(ctx, db.DB, null.IntFrom(0), cfg)
		require.NoError(t, err)

		v2, err = evm.Current(ctx, db.DB, cfg)
		require.NoError(t, err)
		assert.Equal(t, int64(0), v2)
	})
}
