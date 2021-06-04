package headtracker_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockFetcher_GetBlockRange(t *testing.T) {
	t.Parallel()

	config := createConfig()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	ethClient := new(mocks.Client)

	block40 := cltest.NewBlock(40, common.Hash{})
	block41 := cltest.NewBlock(41, block40.Hash)
	block42 := cltest.NewBlock(42, block41.Hash)

	blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{*block40, *block41, *block42})
	blockFetcher := headtracker.NewBlockFetcher(config, logger, blockClient)

	blockRange, err := blockFetcher.BlockRange(context.Background(), 41, 42)
	require.NoError(t, err)

	assert.Len(t, blockRange, 2)
	assert.Len(t, blockFetcher.BlockCache(), 2)

	assert.Equal(t, int64(41), blockRange[0].Number)
	assert.Equal(t, int64(42), blockRange[1].Number)

	assert.Equal(t, block41.Hash, blockRange[0].Hash)
	assert.Equal(t, block42.Hash, blockRange[1].Hash)

	ethClient.AssertExpectations(t)
}

func TestBlockFetcher_ConstructsChain(t *testing.T) {

	config := createConfig()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.NewBlock(40, common.Hash{})
	block41 := cltest.NewBlock(41, block40.Hash)
	block42 := cltest.NewBlock(42, block41.Hash)
	h := headtracker.HeadFromBlock(*block42)

	blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{*block40, *block41, *block42})
	blockFetcher := headtracker.NewBlockFetcher(config, logger, blockClient)

	head, err := blockFetcher.Chain(context.Background(), h)
	require.NoError(t, err)
	assert.Equal(t, 3, int(head.ChainLength()))
}

func TestBlockFetcher_CreatesChainWhereSomeBlocksAreInitiallyMissing(t *testing.T) {

	config := createConfig()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block38 := cltest.NewBlock(38, common.Hash{})
	block39 := cltest.NewBlock(39, block38.Hash)
	block40 := cltest.NewBlock(40, block39.Hash)
	block41 := cltest.NewBlock(41, block40.Hash)
	block42 := cltest.NewBlock(42, block41.Hash)
	h := headtracker.HeadFromBlock(*block42)

	blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{*block38, *block39, *block40, *block41, *block42})
	blockFetcher := headtracker.NewBlockFetcher(config, logger, blockClient)

	_, err := blockFetcher.BlockRange(context.Background(), 39, 39)
	require.NoError(t, err)
	_, err = blockFetcher.BlockRange(context.Background(), 41, 42)
	require.NoError(t, err)
	head, err := blockFetcher.Chain(context.Background(), h)
	require.NoError(t, err)
	assert.Equal(t, 5, int(head.ChainLength()))
}

func createConfig() *htmocks.BlockFetcherConfig {
	config := new(htmocks.BlockFetcherConfig)
	config.On("BlockFetcherBatchSize").Return(uint32(2))
	config.On("EthFinalityDepth").Return(uint(42))
	config.On("EthHeadTrackerHistoryDepth").Return(uint(100))
	config.On("BlockBackfillDepth").Return(uint64(50))
	return config
}
