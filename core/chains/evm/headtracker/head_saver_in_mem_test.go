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
	htCfg.On("EvmHeadTrackerHistoryDepth").Return(uint32(1))
	htCfg.On("EvmFinalityDepth").Return(uint32(1))
	return headtracker.NewEvmInMemoryHeadSaver(htCfg, lggr)
}

func TestInMemoryHeadSaver_Save(t *testing.T) {
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

	t.Run("saving heads with same block number", func(t *testing.T) {
		head := cltest.Head(4)
		err := saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		head = cltest.Head(4)
		err = saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		head = cltest.Head(4)
		err = saver.Save(testutils.Context(t), head)
		require.NoError(t, err)

		latest := saver.LatestChain()
		require.NoError(t, err)
		require.Equal(t, int64(4), latest.Number)

		headsWithSameNumber := len(saver.HeadByNumber(4))
		require.Equal(t, 3, headsWithSameNumber)
	})
}

func TestInMemoryHeadSaver_TrimOldHeads(t *testing.T) {
	t.Parallel()
	saver := configureInMemorySaver(t)

	t.Run("happy path, trimming old heads", func(t *testing.T) {

		// Save heads with block numbers 1, 2, 3, and 4
		for i := 1; i <= 4; i++ {
			head := cltest.Head(i)
			err := saver.Save(testutils.Context(t), head)
			require.NoError(t, err)
		}

		// Trim old heads, keeping only the last two (block numbers 3 and 4)
		saver.TrimOldHeads(3)

		// Check that the correct heads remain
		require.Equal(t, 2, len(saver.Heads))
		require.Equal(t, 1, len(saver.HeadByNumber(3)))
		require.Equal(t, 1, len(saver.HeadByNumber(4)))
		require.Equal(t, 0, len(saver.HeadByNumber(1)))

		// Check that the latest head is correct
		latest := saver.LatestChain()
		require.Equal(t, int64(4), latest.Number)

		// Clear All Heads
		saver.TrimOldHeads(6)
		require.Equal(t, 0, len(saver.Heads))
		require.Equal(t, 0, len(saver.HeadsNumber))
	})

	t.Run("error path, block number lower than highest chain", func(t *testing.T) {
		// Save heads with block numbers 1, 2, 3, and 4
		for i := 1; i <= 4; i++ {
			head := cltest.Head(i)
			err := saver.Save(testutils.Context(t), head)
			require.NoError(t, err)
		}

		saver.TrimOldHeads(1)

		// Check that no heads are removed
		require.Equal(t, 4, len(saver.Heads))
		require.Equal(t, 4, len(saver.HeadsNumber))

		// Check that the latest head remains the same
		latest := saver.LatestChain()
		require.Equal(t, int64(4), latest.Number)
	})
}
