package headtracker_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

func configureInMemorySaver(t *testing.T) *headtracker.EvmInMemoryHeadSaver {
	lggr := logger.TestLogger(t)
	htCfg := htmocks.NewConfig(t)
	htCfg.On("EvmHeadTrackerHistoryDepth").Return(uint32(6))
	htCfg.On("EvmFinalityDepth").Return(uint32(1))
	return headtracker.NewEvmInMemoryHeadSaver(htCfg, lggr)
}

func TestInMemoryHeadSaver_Save_Happy(t *testing.T) {
	t.Parallel()
	saver := configureInMemorySaver(t)

	t.Run("happy path, saving heads", func(t *testing.T) {
		head := cltest.Head(1)
		err := saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		latest := saver.LatestChain()
		require.NoError(t, err)
		require.Equal(t, int64(1), latest.Number)

		latest = saver.LatestChain()
		require.NotNil(t, latest)
		require.Equal(t, int64(1), latest.Number)

		latest = saver.Chain(head.Hash)
		require.NotNil(t, latest)
		require.Equal(t, int64(1), latest.Number)

		// Add more heads
		head = cltest.Head(2)
		err = saver.Save(testutils.Context(t), head)
		require.NoError(t, err)
		head = cltest.Head(3)
		err = saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		latest = saver.LatestChain()
		require.Equal(t, int64(3), latest.Number)
	})

	t.Run("saving head with same block number", func(t *testing.T) {
		head := cltest.Head(1)
		err := saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		latest := saver.LatestChain()
		require.NoError(t, err)
		require.Equal(t, int64(1), latest.Number)

		head = cltest.Head(1)
		err = saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		latest = saver.LatestChain()
		require.NoError(t, err)
		require.Equal(t, int64(1), latest.Number)

		// Log the heads in HeadsNumber
		hs := saver.HeadsNumber
		t.Log(hs)
		// check hs.HeadsNumber to make sure there are 2 heads
		require.Equal(t, 2, len(saver.HeadsNumber))
	})

}
