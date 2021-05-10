package gasupdater_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/gasupdater"
	gumocks "github.com/smartcontractkit/chainlink/core/services/gasupdater/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGasUpdater_Start(t *testing.T) {
	t.Parallel()

	config := new(gumocks.Config)

	var batchSize uint32 = 0
	var blockDelay uint16 = 0
	var historySize uint16 = 2
	var ethFinalityDepth uint = 42
	var percentile uint16 = 35

	config.On("GasUpdaterBatchSize").Return(batchSize)
	config.On("GasUpdaterBlockDelay").Return(blockDelay)
	config.On("GasUpdaterBlockHistorySize").Return(historySize)
	config.On("EthFinalityDepth").Return(ethFinalityDepth)
	config.On("GasUpdaterTransactionPercentile").Return(percentile)

	t.Run("loads initial state", func(t *testing.T) {
		ethClient := new(mocks.Client)

		guIface := gasupdater.NewGasUpdater(ethClient, config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		h := &models.Head{Hash: cltest.NewHash(), Number: 42}
		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&models.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&models.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &models.Block{
				Number: 42,
				Hash:   cltest.NewHash(),
			}
			elems[1].Result = &models.Block{
				Number: 41,
				Hash:   cltest.NewHash(),
			}
		})

		err := guIface.Start()
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 2)
		assert.Equal(t, int(gu.RollingBlockHistory()[0].Number), 41)
		assert.Equal(t, int(gu.RollingBlockHistory()[1].Number), 42)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("boots even if initial batch call returns nothing", func(t *testing.T) {
		ethClient := new(mocks.Client)

		gu := gasupdater.NewGasUpdater(ethClient, config)

		h := &models.Head{Hash: cltest.NewHash(), Number: 42}
		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == int(historySize)
		})).Return(nil)

		err := gu.Start()
		require.NoError(t, err)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("starts anyway if fetching latest head fails", func(t *testing.T) {
		ethClient := new(mocks.Client)

		gu := gasupdater.NewGasUpdater(ethClient, config)

		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("something exploded"))

		err := gu.Start()
		require.NoError(t, err)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})
}

func TestGasUpdater_FetchBlocks(t *testing.T) {
	t.Parallel()

	t.Run("with history size of 0, errors", func(t *testing.T) {
		ethClient := new(mocks.Client)
		config := new(gumocks.Config)
		gu := gasupdater.GasUpdaterToStruct(gasupdater.NewGasUpdater(ethClient, config))

		var blockDelay uint16 = 3
		var historySize uint16 = 0
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)

		head := cltest.Head(42)
		err := gu.FetchBlocks(context.Background(), *head)
		require.Error(t, err)
		require.EqualError(t, err, "GasUpdater: history size must be > 0, got: 0")
	})

	t.Run("with current block height less than block delay does nothing", func(t *testing.T) {
		ethClient := new(mocks.Client)
		config := new(gumocks.Config)
		gu := gasupdater.GasUpdaterToStruct(gasupdater.NewGasUpdater(ethClient, config))

		var blockDelay uint16 = 3
		var historySize uint16 = 1
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)

		for i := -1; i < 3; i++ {
			head := cltest.Head(i)
			err := gu.FetchBlocks(context.Background(), *head)
			require.Error(t, err)
			require.EqualError(t, err, fmt.Sprintf("GasUpdater: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY=3", i))
		}

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("with error retrieving blocks returns error", func(t *testing.T) {
		ethClient := new(mocks.Client)
		config := new(gumocks.Config)
		gu := gasupdater.GasUpdaterToStruct(gasupdater.NewGasUpdater(ethClient, config))

		var blockDelay uint16 = 3
		var historySize uint16 = 3
		var batchSize uint32 = 0
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)
		config.On("GasUpdaterBatchSize").Return(batchSize)

		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

		err := gu.FetchBlocks(context.Background(), *cltest.Head(42))
		require.Error(t, err)
		assert.EqualError(t, err, "GasUpdater#fetchBlocks error fetching blocks with BatchCallContext: something exploded")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("batch fetches heads and transactions and sets them on the gas updater instance", func(t *testing.T) {
		ethClient := new(mocks.Client)
		config := new(gumocks.Config)
		gu := gasupdater.GasUpdaterToStruct(gasupdater.NewGasUpdater(ethClient, config))

		var blockDelay uint16 = 1
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)
		// Test batching
		config.On("GasUpdaterBatchSize").Return(batchSize)

		b41 := models.Block{
			Number:       41,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(1, 2),
		}
		b42 := models.Block{
			Number:       42,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(3),
		}
		b43 := models.Block{
			Number:       43,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x28" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&models.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x29" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&models.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b41 // This errored block will be ignored
			elems[1].Error = errors.New("something went wrong")
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&models.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b43
		})

		err := gu.FetchBlocks(context.Background(), *cltest.Head(43))
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 2)
		assert.Equal(t, 41, int(gu.RollingBlockHistory()[0].Number))
		assert.Equal(t, 43, int(gu.RollingBlockHistory()[1].Number))
		assert.Len(t, gu.RollingBlockHistory()[0].Transactions, 2)
		assert.Len(t, gu.RollingBlockHistory()[1].Transactions, 0)

		ethClient.AssertExpectations(t)

		// On new fetch, rolls over the history and drops the old heads

		b44 := models.Block{
			Number:       44,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(4),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&models.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&models.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b42
			elems[1].Result = &b43
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2b" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&models.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b44
		})

		err = gu.FetchBlocks(context.Background(), *cltest.Head(44))
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 3)
		assert.Equal(t, 42, int(gu.RollingBlockHistory()[0].Number))
		assert.Equal(t, 43, int(gu.RollingBlockHistory()[1].Number))
		assert.Equal(t, 44, int(gu.RollingBlockHistory()[2].Number))
		assert.Len(t, gu.RollingBlockHistory()[0].Transactions, 1)
		assert.Len(t, gu.RollingBlockHistory()[1].Transactions, 0)
		assert.Len(t, gu.RollingBlockHistory()[2].Transactions, 1)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})
}

func TestGasUpdater_FetchBlocksAndRecalculate(t *testing.T) {
	t.Parallel()

	ethClient := new(mocks.Client)
	config := new(gumocks.Config)

	config.On("GasUpdaterBlockDelay").Return(uint16(0))
	config.On("GasUpdaterTransactionPercentile").Return(uint16(35))
	config.On("GasUpdaterBlockHistorySize").Return(uint16(3))
	config.On("EthMaxGasPriceWei").Return(big.NewInt(1000))
	config.On("GasUpdaterBatchSize").Return(uint32(0))

	guIface := gasupdater.NewGasUpdater(ethClient, config)
	gu := gasupdater.GasUpdaterToStruct(guIface)

	b1 := models.Block{
		Number:       1,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(1),
	}
	b2 := models.Block{
		Number:       2,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(2),
	}
	b3 := models.Block{
		Number:       3,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(200, 300, 100, 100, 100, 100),
	}

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 3 &&
			b[0].Args[0] == "0x1" &&
			b[1].Args[0] == "0x2" &&
			b[2].Args[0] == "0x3"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b1
		elems[1].Result = &b2
		elems[2].Result = &b3
	})

	config.On("SetEthGasPriceDefault", big.NewInt(100)).Return(nil)

	gu.FetchBlocksAndRecalculate(context.Background(), *cltest.Head(3))

	assert.Len(t, gu.RollingBlockHistory(), 3)

	config.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestGasUpdater_Recalculate(t *testing.T) {
	t.Parallel()

	ethClient := new(mocks.Client)
	config := new(gumocks.Config)

	maxGasPrice := big.NewInt(100)
	config.On("EthMaxGasPriceWei").Return(maxGasPrice)
	config.On("GasUpdaterTransactionPercentile").Return(uint16(35))

	guIface := gasupdater.NewGasUpdater(ethClient, config)
	gu := gasupdater.GasUpdaterToStruct(guIface)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		blocks := []models.Block{}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))

		blocks = []models.Block{models.Block{}}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))

		blocks = []models.Block{models.Block{Transactions: []types.Transaction{}}}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))
	})

	t.Run("sets gas price to ETH_MAX_GAS_PRICE_WEI if the calculation would otherwise exceed it", func(t *testing.T) {
		blocks := []models.Block{
			models.Block{
				Number:       0,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9001),
			},
			models.Block{
				Number:       1,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9002),
			},
		}

		config.On("SetEthGasPriceDefault", maxGasPrice).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(1))
	})

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
}
