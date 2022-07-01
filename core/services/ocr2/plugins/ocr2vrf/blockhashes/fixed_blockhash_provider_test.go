package blockhashes_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/blockhashes"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_FixedBlockhashProvider(t *testing.T) {
	t.Parallel()

	client := evmtest.NewEthClientMockWithDefaultChain(t)

	p := blockhashes.NewFixedBlockhashProvider(client, 0, 0)
	ctx := testutils.Context(t)

	t.Run("returns current height", func(t *testing.T) {
		h := &evmtypes.Head{Number: 100}
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()
		height, err := p.CurrentHeight(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), height)
	})

	t.Run("returns error when negative block number", func(t *testing.T) {
		h := &evmtypes.Head{Number: -10}
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()
		height, err := p.CurrentHeight(ctx)
		require.Error(t, err)
		assert.Equal(t, uint64(0), height)
	})
}

func Test_OnchainVerifiableBlocks(t *testing.T) {
	t.Parallel()

	client := evmtest.NewEthClientMockWithDefaultChain(t)
	ctx := testutils.Context(t)
	h := &evmtypes.Head{Number: 100}

	t.Run("returns expected number of hashes", func(t *testing.T) {
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				reqs[i].Result = &evmtypes.Head{Hash: utils.NewHash()}
			}
		}).Times(5)

		p := blockhashes.NewFixedBlockhashProvider(client, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint64(100-8), startHeight)
		assert.Equal(t, 9, len(hashes))
		for _, hash := range hashes {
			assert.NotEmpty(t, hash)
		}
	})

	t.Run("returns error when underlying batch call returns error", func(t *testing.T) {
		client = evmtest.NewEthClientMockWithDefaultChain(t)
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()

		e := errors.New("network error")
		client.On("BatchCallContext", ctx, mock.Anything).Return(e).Once()

		p := blockhashes.NewFixedBlockhashProvider(client, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, "batch call context eth_getBlockByNumber: network error", err.Error())
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
	})

	t.Run("returns error when nil block received", func(t *testing.T) {
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				reqs[i].Result = nil
			}
		}).Times(5)

		p := blockhashes.NewFixedBlockhashProvider(client, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
	})

	t.Run("returns error when empty blockhash received", func(t *testing.T) {
		client.On("HeadByNumber", ctx, mock.MatchedBy(func(val *big.Int) bool {
			return val == nil
		})).Return(h, nil).Once()

		client.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				reqs[i].Result = &evmtypes.Head{Hash: utils.EmptyHash}
			}
		}).Times(5)

		p := blockhashes.NewFixedBlockhashProvider(client, 8, 2)
		startHeight, hashes, err := p.OnchainVerifiableBlocks(ctx)

		require.Error(t, err)
		assert.Equal(t, "missing block hash", err.Error())
		assert.Equal(t, uint64(0), startHeight)
		assert.Nil(t, hashes)
	})
}
