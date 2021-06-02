package headtracker_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockFetcher_GetBlockRange(t *testing.T) {
	t.Parallel()

	config := headtracker.NewBlockFetcherConfigWithDefaults()

	t.Run("fetches a range of blocks", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		logger := store.Config.CreateProductionLogger()

		ethClient := new(mocks.Client)

		block40 := cltest.HtBlock(40, common.Hash{})
		block41 := cltest.HtBlock(41, block40.Hash)
		block42 := cltest.HtBlock(42, block41.Hash)
		//h := cltest.HeadFromHtBlock(&block42)

		blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{block40, block41, block42})

		blockFetcher := headtracker.NewBlockFetcher(config, logger, blockClient)

		//ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		//	return len(b) == 2 &&
		//		b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true &&
		//		b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true
		//})).Return(nil).Run(func(args mock.Arguments) {
		//	elems := args.Get(1).([]rpc.BatchElem)
		//	elems[0].Result = &headtracker.Block{
		//		Number: 42,
		//		Hash:   utils.NewHash(),
		//	}
		//	elems[1].Result = &headtracker.Block{
		//		Number: 41,
		//		Hash:   utils.NewHash(),
		//	}
		//})

		blockRange, err := blockFetcher.BlockRange(context.Background(), 41, 42)
		require.NoError(t, err)

		assert.Len(t, blockRange, 2)
		assert.Len(t, blockFetcher.BlockCache(), 2)

		assert.Equal(t, int64(41), blockRange[0].Number)
		assert.Equal(t, int64(42), blockRange[1].Number)

		assert.Equal(t, block41.Hash, blockRange[0].Hash)
		assert.Equal(t, block42.Hash, blockRange[1].Hash)

		ethClient.AssertExpectations(t)
	})
}

func TestBlockFetcher_ConstructsChain(t *testing.T) {

	config := headtracker.NewBlockFetcherConfigWithDefaults()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.HtBlock(40, common.Hash{})
	block41 := cltest.HtBlock(41, block40.Hash)
	block42 := cltest.HtBlock(42, block41.Hash)
	h := cltest.HeadFromHtBlock(&block42)

	blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{block40, block41, block42})
	blockFetcher := headtracker.NewBlockFetcher(config, logger, blockClient)

	head, err := blockFetcher.Chain(context.Background(), *h)
	require.NoError(t, err)
	assert.Equal(t, 3, int(head.ChainLength()))
}
