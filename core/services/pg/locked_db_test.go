package pg_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestLockedDB_HappyPath(t *testing.T) {
	testutils.SkipShortDB(t)
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.DatabaseLockingMode = null.StringFrom("dual")
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	err := ldb.Open(context.Background())
	require.NoError(t, err)
	require.NotNil(t, ldb.DB())

	err = ldb.Close()
	require.NoError(t, err)
	require.Nil(t, ldb.DB())
}

func TestLockedDB_ContextCancelled(t *testing.T) {
	testutils.SkipShortDB(t)
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.DatabaseLockingMode = null.StringFrom("dual")
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ldb.Open(ctx)
	require.Error(t, err)
	require.Nil(t, ldb.DB())
}

func TestLockedDB_OpenTwice(t *testing.T) {
	testutils.SkipShortDB(t)
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.DatabaseLockingMode = null.StringFrom("lease")
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	err := ldb.Open(context.Background())
	require.NoError(t, err)
	require.Panics(t, func() {
		_ = ldb.Open(context.Background())
	})

	_ = ldb.Close()
}

func TestLockedDB_TwoInstances(t *testing.T) {
	testutils.SkipShortDB(t)
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.DatabaseLockingMode = null.StringFrom("dual")
	lggr := logger.TestLogger(t)

	ldb1 := pg.NewLockedDB(config, lggr)
	err := ldb1.Open(context.Background())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, ldb1.Close())
	}()

	// second instance would wait for locks to be released,
	// hence we use some timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.LeaseLockDuration())
	defer cancel()
	ldb2 := pg.NewLockedDB(config, lggr)
	err = ldb2.Open(ctx)
	require.Error(t, err)
}

func TestOpenUnlockedDB(t *testing.T) {
	testutils.SkipShortDB(t)
	config := cltest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)

	db1, err1 := pg.OpenUnlockedDB(config, lggr)
	require.NoError(t, err1)
	require.NotNil(t, db1)

	// should not block the second connection
	db2, err2 := pg.OpenUnlockedDB(config, lggr)
	require.NoError(t, err2)
	require.NotNil(t, db2)

	require.NoError(t, db1.Close())
	require.NoError(t, db2.Close())
}
