package services_test

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
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGasUpdater_Start(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	store.Config.Set("GAS_UPDATER_BLOCK_DELAY", 0)
	store.Config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", 2)

	t.Run("loads initial state", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		guIface := services.NewGasUpdater(store)
		gu := services.GasUpdaterToStruct(guIface)

		h := &models.Head{Hash: cltest.NewHash(), Number: 42}
		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&services.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&services.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &services.Block{
				Number: 42,
				Hash:   cltest.NewHash(),
			}
			elems[1].Result = &services.Block{
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
	})

	t.Run("boots even if initial batch call returns nothing", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		gu := services.NewGasUpdater(store)

		h := &models.Head{Hash: cltest.NewHash(), Number: 42}
		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == int(store.Config.GasUpdaterBlockHistorySize())
		})).Return(nil)

		err := gu.Start()
		require.NoError(t, err)

		ethClient.AssertExpectations(t)
	})

	t.Run("starts anyway if fetching latest head fails", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		gu := services.NewGasUpdater(store)

		ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("something exploded"))

		err := gu.Start()
		require.NoError(t, err)

		ethClient.AssertExpectations(t)
	})
}

func TestGasUpdater_FetchBlocks(t *testing.T) {
	t.Parallel()

	t.Run("with current block height less than block delay does nothing", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		store.Config.Set("GAS_UPDATER_BLOCK_DELAY", 3)

		gu := services.GasUpdaterToStruct(services.NewGasUpdater(store))

		for i := -1; i < 3; i++ {
			head := cltest.Head(i)
			err := gu.FetchBlocks(context.Background(), *head)
			require.Error(t, err)
			require.EqualError(t, err, fmt.Sprintf("GasUpdater: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY of 3", i))
		}

		ethClient.AssertExpectations(t)
	})

	t.Run("with error retrieving blocks returns error", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

		store.Config.Set("GAS_UPDATER_BLOCK_DELAY", 3)

		gu := services.GasUpdaterToStruct(services.NewGasUpdater(store))

		err := gu.FetchBlocks(context.Background(), *cltest.Head(42))
		require.Error(t, err)
		assert.EqualError(t, err, "GasUpdater#fetchBlocks error fetching blocks with BatchCallContext: something exploded")

		ethClient.AssertExpectations(t)
	})

	t.Run("fetches heads and transactions and sets them on the gas updater instance", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		store.Config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", 3)
		store.Config.Set("GAS_UPDATER_BLOCK_DELAY", 1)

		gu := services.GasUpdaterToStruct(services.NewGasUpdater(store))

		b41 := services.Block{
			Number:       41,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(1, 2),
		}
		b42 := services.Block{
			Number:       42,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(3),
		}
		b43 := services.Block{
			Number:       43,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x28" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&services.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x29" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&services.Block{}) &&
				b[2].Method == "eth_getBlockByNumber" && b[2].Args[0] == "0x2a" && b[2].Args[1] == true && reflect.TypeOf(b[2].Result) == reflect.TypeOf(&services.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b41 // This errored block will be ignored
			elems[1].Error = errors.New("something went wrong")
			elems[2].Result = &b43
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

		b44 := services.Block{
			Number:       44,
			Hash:         cltest.NewHash(),
			Transactions: cltest.TransactionsFromGasPrices(4),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&services.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&services.Block{}) &&
				b[2].Method == "eth_getBlockByNumber" && b[2].Args[0] == "0x2b" && b[2].Args[1] == true && reflect.TypeOf(b[2].Result) == reflect.TypeOf(&services.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b42
			elems[1].Result = &b43
			elems[2].Result = &b44
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
	})
}

func TestGasUpdater_FetchBlocksAndRecalculate(t *testing.T) {
	t.Parallel()
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_BLOCK_DELAY", "0")
	config.Set("GAS_UPDATER_TRANSACTION_PERCENTILE", "35")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "3")
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient
	guIface := services.NewGasUpdater(store)
	gu := services.GasUpdaterToStruct(guIface)

	b1 := services.Block{
		Number:       1,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(1),
	}
	b2 := services.Block{
		Number:       2,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(2),
	}
	b3 := services.Block{
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

	gu.FetchBlocksAndRecalculate(context.Background(), *cltest.Head(3))

	assert.Len(t, gu.RollingBlockHistory(), 3)
	assert.Equal(t, big.NewInt(100), config.EthGasPriceDefault())
}

func TestGasUpdater_Recalculate(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	config.Set("ETH_MAX_GAS_PRICE_WEI", 100)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient
	guIface := services.NewGasUpdater(store)
	gu := services.GasUpdaterToStruct(guIface)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		require.Equal(t, big.NewInt(42), config.EthGasPriceDefault())

		blocks := []services.Block{}
		services.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))
		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())

		blocks = []services.Block{services.Block{}}
		services.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))
		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())

		blocks = []services.Block{services.Block{Transactions: []types.Transaction{}}}
		services.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))
		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
	})

	t.Run("will not set gas price higher than ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		blocks := []services.Block{
			services.Block{
				Number:       0,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9001),
			},
			services.Block{
				Number:       1,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9002),
			},
		}
		services.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(1))

		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
	})
}
