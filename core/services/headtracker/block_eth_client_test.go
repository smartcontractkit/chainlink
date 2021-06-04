package headtracker_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlockEthClient_FastBlockByHash(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.NewGethBlock(40, common.Hash{})

	ethClient := new(mocks.Client)
	ethClient.On("FastBlockByHash", mock.Anything, mock.Anything).Return(block40, nil)

	blockClient := headtracker.NewBlockEthClientImpl(ethClient, logger, 2)

	block, err := blockClient.FastBlockByHash(context.Background(), block40.Hash())
	require.NoError(t, err)

	assert.Equal(t, int64(40), block.Number)
	assert.Equal(t, block40.Hash(), block.Hash)
}

func TestBlockEthClient_BlockByNumber(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.NewGethBlock(40, common.Hash{})

	ethClient := new(mocks.Client)
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Return(block40, nil)

	blockClient := headtracker.NewBlockEthClientImpl(ethClient, logger, 2)

	block, err := blockClient.BlockByNumber(context.Background(), 40)
	require.NoError(t, err)

	assert.Equal(t, int64(40), block.Number)
	assert.Equal(t, block40.Hash(), block.Hash)
}

func TestBlockEthClient_BatchGetBlocks(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.NewBlock(40, common.Hash{})
	block41 := cltest.NewBlock(41, block40.Hash)
	block42 := cltest.NewBlock(42, block41.Hash)
	block43 := cltest.NewBlock(43, block41.Hash)

	ethClient := new(mocks.Client)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = block42
		elems[1].Result = block41
	})

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2b" && b[0].Args[1] == true
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = block43
	})

	blockClient := headtracker.NewBlockEthClientImpl(ethClient, logger, 2)

	blocks, err := blockClient.FetchBlocksByNumbers(context.Background(), []int64{41, 42, 43})
	require.NoError(t, err)

	assert.Len(t, blocks, 3)
	assert.Equal(t, int64(41), blocks[41].Number)
	assert.Equal(t, int64(42), blocks[42].Number)
	assert.Equal(t, int64(43), blocks[43].Number)
}

func TestBlockEthClient_BatchReturnsFewerBlocksOnError(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	block40 := cltest.NewBlock(40, common.Hash{})
	block41 := cltest.NewBlock(41, block40.Hash)
	block42 := cltest.NewBlock(42, block41.Hash)

	ethClient := new(mocks.Client)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = block42
		elems[1].Error = errors.New("something exploded")
	})

	blockClient := headtracker.NewBlockEthClientImpl(ethClient, logger, 2)

	blocks, err := blockClient.FetchBlocksByNumbers(context.Background(), []int64{41, 42})
	require.NoError(t, err)

	assert.Len(t, blocks, 1)
	assert.Equal(t, int64(42), blocks[42].Number)
}
