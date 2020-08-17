package services_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGasUpdater_OnNewLongestChain_whenDisabledDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "false")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	gu := services.NewGasUpdater(store)
	head := cltest.Head(0)

	gu.OnNewLongestChain(*head)

	// No mock calls
	ethClient.AssertExpectations(t)
}

func TestGasUpdater_OnNewLongestChain_WithCurrentBlockHeightLessThanBlockDelayDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "3")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	gu := services.NewGasUpdater(store)

	for i := -1; i < 3; i++ {
		head := cltest.Head(i)
		gu.OnNewLongestChain(*head)
	}

	// No mock calls
	ethClient.AssertExpectations(t)
}

func TestGasUpdater_OnNewLongestChain_WithErrorRetrievingBlockDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	gu := services.NewGasUpdater(store)

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(0)).Return(nil, errors.New("foo"))

	head := cltest.Head(3)

	gu.OnNewLongestChain(*head)
	ethClient.AssertExpectations(t)
}

func TestGasUpdater_OnNewLongestChain_AddsBlockToBlockHistory(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	gu := services.NewGasUpdater(store)

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(0)).Return(cltest.BlockWithTransactions(), nil)
	head := cltest.Head(3)
	gu.OnNewLongestChain(*head)

	// Empty blocks are not added
	assert.Len(t, gu.RollingBlockHistory(), 0)

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(1)).Return(cltest.BlockWithTransactions(20000), nil)
	head = cltest.Head(4)
	gu.OnNewLongestChain(*head)

	// Blocks with transactions are added
	assert.Len(t, gu.RollingBlockHistory(), 1)

	ethClient.AssertExpectations(t)
}

func TestGasUpdater_OnNewLongestChain_DoesNotOverflowBlockHistory(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "3")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "5")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient
	gu := services.NewGasUpdater(store)

	for i := 0; i < 5; i++ {
		ethClient.On("BlockByNumber", mock.Anything, big.NewInt(int64(i))).Return(cltest.BlockWithTransactions(42), nil)
		head := cltest.Head(i + 3)
		gu.OnNewLongestChain(*head)
		assert.Len(t, gu.RollingBlockHistory(), i+1)
	}

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(5)).Return(cltest.BlockWithTransactions(42), nil)
	head := cltest.Head(8)
	gu.OnNewLongestChain(*head)

	assert.Len(t, gu.RollingBlockHistory(), 5)
}

func TestGasUpdater_OnNewLongestChain_SetsGlobalGasPriceWhenHistoryFull(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "0")
	config.Set("GAS_UPDATER_TRANSACTION_PERCENTILE", "35")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "3")
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient
	gu := services.NewGasUpdater(store)

	for i := 0; i < 3; i++ {
		ethClient.On("BlockByNumber", mock.Anything, big.NewInt(int64(i))).Return(cltest.BlockWithTransactions(int64((1+i)*100)), nil)
		head := cltest.Head(i)
		gu.OnNewLongestChain(*head)
		assert.Len(t, gu.RollingBlockHistory(), i+1)
		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
	}

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(3)).Return(cltest.BlockWithTransactions(200, 300, 100, 100, 100, 100), nil)
	head := cltest.Head(3)
	gu.OnNewLongestChain(*head)

	assert.Len(t, gu.RollingBlockHistory(), 3)
	assert.Equal(t, big.NewInt(100), config.EthGasPriceDefault())
}

func TestGasUpdater_OnNewLongestChain_WillNotSetGasHigherThanEthMaxGasPriceWei(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "0")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "1")
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	config.Set("ETH_MAX_GAS_PRICE_WEI", 100)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient
	gu := services.NewGasUpdater(store)

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(0)).Return(cltest.BlockWithTransactions(9001), nil)
	head := cltest.Head(0)
	gu.OnNewLongestChain(*head)
	assert.Len(t, gu.RollingBlockHistory(), 1)
	assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())

	ethClient.On("BlockByNumber", mock.Anything, big.NewInt(1)).Return(cltest.BlockWithTransactions(9002), nil)
	head = cltest.Head(1)
	gu.OnNewLongestChain(*head)

	assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
}
