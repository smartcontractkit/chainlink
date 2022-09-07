package blockhashes_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	lp_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/blockhashes"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_CurrentHeight(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	lp := lp_mocks.NewLogPoller(t)
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestLogger(t)
	p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 0, 0)

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
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestLogger(t)
	h := int64(100)

	t.Run("returns expected number of hashes", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		mockBatchCallContext(ctx, t, client, []int64{92, 94, 97}, 2)

		blocks := []logpoller.LogPollerBlock{
			createLogPollerBlock(93),
			createLogPollerBlock(95),
			createLogPollerBlock(96),
			createLogPollerBlock(98),
			createLogPollerBlock(99),
			createLogPollerBlock(100),
		}

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(blocks, nil).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(100-8), startHeight)
		assert.Equal(t, 9, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns expected number of hashes when all blocks returned by lp", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		blocks := []logpoller.LogPollerBlock{
			createLogPollerBlock(92),
			createLogPollerBlock(93),
			createLogPollerBlock(94),
			createLogPollerBlock(95),
			createLogPollerBlock(96),
			createLogPollerBlock(97),
			createLogPollerBlock(98),
			createLogPollerBlock(99),
			createLogPollerBlock(100),
		}

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(blocks, nil).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(100-8), startHeight)
		assert.Equal(t, 9, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns expected number of hashes when no blocks returned by lp", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		mockBatchCallContext(ctx, t, client, []int64{92, 93, 94, 95, 96, 97, 98, 99, 100}, 5)

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(nil, errors.New("no blocks found")).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(100-8), startHeight)
		assert.Equal(t, 9, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns error when batch call returns error", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(nil, errors.New("no blocks found")).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(errors.New("network error")).Times(1)

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, "batch call context eth_getBlockByNumber: network error", err.Error())
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		lp.AssertExpectations(t)
	})

	t.Run("returns error when nil block received", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(nil, errors.New("no blocks found")).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				reqs[i].Result = nil
			}
		}).Times(5)

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		lp.AssertExpectations(t)
	})

	t.Run("returns error when empty blockhash received", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(h, nil).Once()

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{92, 93, 94, 95, 96, 97, 98, 99, 100})
		}), mock.Anything).Return(nil, errors.New("no blocks found")).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				reqs[i].Result = &evmtypes.Head{Hash: utils.EmptyHash}
			}
		}).Times(5)

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		lp.AssertExpectations(t)
	})

	t.Run("returns expected number of hashes when startHeight less than lookback", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)

		lp.On("LatestBlock", mock.Anything).Return(int64(2), nil).Once()

		mockBatchCallContext(ctx, t, client, []int64{2}, 1)

		blocks := []logpoller.LogPollerBlock{
			createLogPollerBlock(0),
			createLogPollerBlock(1),
		}

		lp.On("GetBlocks", mock.MatchedBy(func(val []uint64) bool {
			return slicesEqual(val, []uint64{0, 1, 2})
		}), mock.Anything).Return(blocks, nil).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Equal(t, 3, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
		lp.AssertExpectations(t)
	})

	t.Run("returns error when current height errors", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).Return(int64(0), errors.New("error in latest block")).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, lp, lggr, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
		lp.AssertExpectations(t)
	})
}

func mockBatchCallContext(ctx context.Context, t *testing.T, client *mocks.Client, expectedBlockNums []int64, times int) {
	var expected []*big.Int
	for _, bn := range expectedBlockNums {
		expected = append(expected, big.NewInt(bn))
	}
	client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		reqs := args.Get(1).([]rpc.BatchElem)
		for i := 0; i < len(reqs); i++ {
			blockNumString, is := reqs[i].Args[0].(string)
			assert.True(t, is)
			blockNum, ok := new(big.Int).SetString(blockNumString, 0)
			assert.True(t, ok)
			found := false
			for _, bn := range expected {
				if blockNum.Cmp(bn) == 0 {
					found = true
					reqs[i].Result = &evmtypes.Head{Hash: utils.NewHash(), Number: bn.Int64()}
					break
				}
			}
			if !found {
				assert.Fail(t, "Received unexepcted block number %d in mock", blockNum)
			}
		}
	}).Times(times)
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
