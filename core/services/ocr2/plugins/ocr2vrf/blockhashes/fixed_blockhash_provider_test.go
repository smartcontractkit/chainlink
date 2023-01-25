package blockhashes_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	lp_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/blockhashes"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_CurrentHeight(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	lp := lp_mocks.NewLogPoller(t)
	lggr := logger.TestLogger(t)
	p := blockhashes.NewFixedBlockhashProvider(lp, lggr, 8)

	t.Run("returns current height", func(t *testing.T) {
		h := int64(100)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()
		height, err := p.CurrentHeight(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), height)
		lp.AssertExpectations(t)
	})

	t.Run("returns error when log poller throws error", func(t *testing.T) {
		lp.On("LatestBlock", mock.Anything).Return(int64(0), errors.New("error in latest block")).Once()
		height, err := p.CurrentHeight(ctx)
		require.Error(t, err)
		assert.Equal(t, uint64(0), height)
		lp.AssertExpectations(t)
	})
}

func Test_OnchainVerifiableBlocks(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	h := int64(100)

	t.Run("returns expected number of hashes", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		blocks := []logpoller.LogPollerBlock{
			createLogPollerBlock(93),
			createLogPollerBlock(94),
			createLogPollerBlock(95),
			createLogPollerBlock(96),
			createLogPollerBlock(97),
			createLogPollerBlock(98),
			createLogPollerBlock(99),
			createLogPollerBlock(100),
		}

		lp.On("GetBlocksRange", ctx, mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(blocks, nil).Once()

		p := blockhashes.NewFixedBlockhashProvider(lp, lggr, 8)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(100-7), startHeight)
		assert.Equal(t, 8, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns max expected blocks", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(int64(1000), nil).Once()

		var blocks []logpoller.LogPollerBlock
		var blockHeights []uint64
		for i := (1000 - 255); i <= 1000; i++ {
			blocks = append(blocks, createLogPollerBlock(int64(i)))
			blockHeights = append(blockHeights, uint64(i))
		}

		lp.On("GetBlocksRange", ctx, mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, blockHeights)
		}), mock.Anything).Return(blocks, nil).Once()

		p := blockhashes.NewFixedBlockhashProvider(lp, lggr, 500)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(1000-255), startHeight)
		assert.Equal(t, 256, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns error when get blocks errors", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		lp.On("GetBlocksRange", ctx, mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(nil, errors.New("error in LP")).Once()

		p := blockhashes.NewFixedBlockhashProvider(lp, lggr, 8)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		assert.Equal(t, "error in LP", err.Error())
		lp.AssertExpectations(t)
	})

	t.Run("returns error when current height errors", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(int64(0), errors.New("error in latest block")).Once()

		p := blockhashes.NewFixedBlockhashProvider(lp, lggr, 8)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		lp.AssertExpectations(t)
	})
}

func createLogPollerBlock(blockNumber int64) logpoller.LogPollerBlock {
	return logpoller.LogPollerBlock{
		BlockNumber: blockNumber,
		BlockHash:   utils.NewHash(),
	}
}

func slicesEqual(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
