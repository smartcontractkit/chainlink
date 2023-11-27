package headtracker_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

type headTrackerConfig struct {
	historyDepth uint32
}

func (h *headTrackerConfig) HistoryDepth() uint32 {
	return h.historyDepth
}

func (h *headTrackerConfig) SamplingInterval() time.Duration {
	return time.Duration(0)
}

func (h *headTrackerConfig) MaxBufferSize() uint32 {
	return uint32(0)
}

type config struct {
	finalityDepth                     uint32
	blockEmissionIdleWarningThreshold time.Duration
}

func (c *config) FinalityDepth() uint32 { return c.finalityDepth }
func (c *config) BlockEmissionIdleWarningThreshold() time.Duration {
	return c.blockEmissionIdleWarningThreshold
}

func configureSaver(t *testing.T) (httypes.HeadSaver, headtracker.ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.Test(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	htCfg := &config{finalityDepth: uint32(1)}
	orm := headtracker.NewORM(db, lggr, cfg.Database(), cltest.FixtureChainID)
	saver := headtracker.NewHeadSaver(lggr, orm, htCfg, &headTrackerConfig{historyDepth: 6})
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

func TestHeadSaver_Load(t *testing.T) {
	t.Parallel()

	saver, orm := configureSaver(t)

	for i := 0; i < 5; i++ {
		err := orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(i))
		require.NoError(t, err)
	}

	latestHead, err := saver.Load(testutils.Context(t))
	require.NoError(t, err)
	require.NotNil(t, latestHead)
	require.Equal(t, int64(4), latestHead.Number)

	latestChain := saver.LatestChain()
	require.NotNil(t, latestChain)
	require.Equal(t, int64(4), latestChain.Number)
}
