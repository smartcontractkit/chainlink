package headtracker_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
)

func configureSaver(t *testing.T) (httypes.HeadSaver, headtracker.ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	htCfg := new(htmocks.Config)
	htCfg.Test(t)
	htCfg.On("EvmHeadTrackerHistoryDepth").Return(uint32(6))
	htCfg.On("EvmFinalityDepth").Return(uint32(1))
	orm := headtracker.NewORM(db, lggr, cfg, cltest.FixtureChainID)
	saver := headtracker.NewHeadSaver(lggr, orm, htCfg)
	return saver, orm
}

func TestHeadSaver_Save(t *testing.T) {
	t.Parallel()

	saver, _ := configureSaver(t)

	head := cltest.Head(1)
	err := saver.Save(context.TODO(), head)
	require.NoError(t, err)

	latest, err := saver.LatestHeadFromDB(context.TODO())
	require.NoError(t, err)
	require.Equal(t, int64(1), latest.Number)

	latest = saver.LatestChain()
	require.NotNil(t, latest)
	require.Equal(t, int64(1), latest.Number)

	latest = saver.Chain(head.Hash)
	require.NotNil(t, latest)
	require.Equal(t, int64(1), latest.Number)
}

func TestHeadSaver_LoadFromDB(t *testing.T) {
	t.Parallel()

	saver, orm := configureSaver(t)

	for i := 0; i < 5; i++ {
		err := orm.IdempotentInsertHead(context.TODO(), cltest.Head(i))
		require.NoError(t, err)
	}

	latestHead, err := saver.LoadFromDB(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, latestHead)
	require.Equal(t, int64(4), latestHead.Number)

	latestChain := saver.LatestChain()
	require.NotNil(t, latestChain)
	require.Equal(t, int64(4), latestChain.Number)
}
