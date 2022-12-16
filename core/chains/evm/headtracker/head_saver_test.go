package headtracker_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/mocks"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func configureSaver(t *testing.T) (httypes.HeadSaver, headtracker.ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	htCfg := htmocks.NewConfig(t)
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
	err := saver.Save(testutils.Context(t), head)
	require.NoError(t, err)

	latest, err := saver.LatestHeadFromDB(testutils.Context(t))
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
		err := orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(i))
		require.NoError(t, err)
	}

	latestHead, err := saver.LoadFromDB(testutils.Context(t))
	require.NoError(t, err)
	require.NotNil(t, latestHead)
	require.Equal(t, int64(4), latestHead.Number)

	latestChain := saver.LatestChain()
	require.NotNil(t, latestChain)
	require.Equal(t, int64(4), latestChain.Number)
}
