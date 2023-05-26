package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/stretchr/testify/require"
)

func lease(c *chainlink.Config, s *chainlink.Secrets) {
	t := true
	c.Database.Lock.Enabled = &t
	c.Database.Lock.LeaseDuration = models.MustNewDuration(10 * time.Second)
	c.Database.Lock.LeaseRefreshInterval = models.MustNewDuration(time.Second)
}

func TestLockedDB_HappyPath(t *testing.T) {
	testutils.SkipShortDB(t)
	config := configtest.NewGeneralConfig(t, lease)
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	err := ldb.Open(testutils.Context(t))
	require.NoError(t, err)
	require.NotNil(t, ldb.DB())

	err = ldb.Close()
	require.NoError(t, err)
	require.Nil(t, ldb.DB())
}

func TestLockedDB_ContextCancelled(t *testing.T) {
	testutils.SkipShortDB(t)
	config := configtest.NewGeneralConfig(t, lease)
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	cancel()
	err := ldb.Open(ctx)
	require.Error(t, err)
	require.Nil(t, ldb.DB())
}

func TestLockedDB_OpenTwice(t *testing.T) {
	testutils.SkipShortDB(t)
	config := configtest.NewGeneralConfig(t, lease)
	lggr := logger.TestLogger(t)
	ldb := pg.NewLockedDB(config, lggr)

	err := ldb.Open(testutils.Context(t))
	require.NoError(t, err)
	require.Panics(t, func() {
		_ = ldb.Open(testutils.Context(t))
	})

	_ = ldb.Close()
}

func TestLockedDB_TwoInstances(t *testing.T) {
	testutils.SkipShortDB(t)
	config := configtest.NewGeneralConfig(t, lease)
	lggr := logger.TestLogger(t)

	ldb1 := pg.NewLockedDB(config, lggr)
	err := ldb1.Open(testutils.Context(t))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, ldb1.Close())
	}()

	// second instance would wait for locks to be released,
	// hence we use some timeout
	ctx, cancel := context.WithTimeout(testutils.Context(t), config.LeaseLockDuration())
	defer cancel()
	ldb2 := pg.NewLockedDB(config, lggr)
	err = ldb2.Open(ctx)
	require.Error(t, err)
}

func TestOpenUnlockedDB(t *testing.T) {
	testutils.SkipShortDB(t)
	config := configtest.NewGeneralConfig(t, nil)

	db1, err1 := pg.OpenUnlockedDB(config)
	require.NoError(t, err1)
	require.NotNil(t, db1)

	// should not block the second connection
	db2, err2 := pg.OpenUnlockedDB(config)
	require.NoError(t, err2)
	require.NotNil(t, db2)

	require.NoError(t, db1.Close())
	require.NoError(t, db2.Close())
}
