package services_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/stretchr/testify/assert"
)

func TestGasUpdater_OnNewHead_whenDisabledDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "false")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)
	head := cltest.Head(0)

	gu.OnNewHead(head)

	// No mock calls
	txm.AssertExpectations(t)
}

func TestGasUpdater_OnNewHead_WithCurrentBlockHeightLessThanBlockDelayDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "3")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	for i := -1; i < 3; i++ {
		head := cltest.Head(i)
		gu.OnNewHead(head)
	}

	// No mock calls
	txm.AssertExpectations(t)
}

func TestGasUpdater_OnNewHead_WithErrorRetrievingBlockDoesNothing(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	txm.On("GetBlockByNumber", "0x0").Return(cltest.EmptyBlock(), errors.New("foo"))

	head := cltest.Head(3)

	gu.OnNewHead(head)
	txm.AssertExpectations(t)
}

func TestGasUpdater_OnNewHead_AddsBlockToBlockHistory(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	txm.On("GetBlockByNumber", "0x0").Return(cltest.BlockWithTransactions(), nil)
	head := cltest.Head(3)
	gu.OnNewHead(head)

	// Empty blocks are not added
	assert.Len(t, gu.RollingBlockHistory(), 0)

	txm.On("GetBlockByNumber", "0x1").Return(cltest.BlockWithTransactions(20000), nil)
	head = cltest.Head(4)
	gu.OnNewHead(head)

	// Blocks with transactions are added
	assert.Len(t, gu.RollingBlockHistory(), 1)

	txm.AssertExpectations(t)
}

func TestGasUpdater_OnNewHead_DoesNotOverflowBlockHistory(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "3")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "5")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	for i := 0; i < 5; i++ {
		txm.On("GetBlockByNumber", fmt.Sprintf("0x%v", i)).Return(cltest.BlockWithTransactions(42), nil)
		head := cltest.Head(i + 3)
		gu.OnNewHead(head)
		assert.Len(t, gu.RollingBlockHistory(), i+1)
	}

	txm.On("GetBlockByNumber", "0x5").Return(cltest.BlockWithTransactions(42), nil)
	head := cltest.Head(8)
	gu.OnNewHead(head)

	assert.Len(t, gu.RollingBlockHistory(), 5)
}

func TestGasUpdater_OnNewHead_SetsGlobalGasPriceWhenHistoryFull(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "0")
	config.Set("GAS_UPDATER_TRANSACTION_PERCENTILE", "35")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "3")
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	for i := 0; i < 3; i++ {
		txm.On("GetBlockByNumber", fmt.Sprintf("0x%v", i)).Return(cltest.BlockWithTransactions(uint64((1+i)*100)), nil)
		head := cltest.Head(i)
		gu.OnNewHead(head)
		assert.Len(t, gu.RollingBlockHistory(), i+1)
		assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
	}

	txm.On("GetBlockByNumber", "0x3").Return(cltest.BlockWithTransactions(200, 300, 100, 100, 100, 100), nil)
	head := cltest.Head(3)
	gu.OnNewHead(head)

	assert.Len(t, gu.RollingBlockHistory(), 3)
	assert.Equal(t, big.NewInt(100), config.EthGasPriceDefault())
}

func TestGasUpdater_OnNewHead_WillNotSetGasHigherThanEthMaxGasPriceWei(t *testing.T) {
	config, _ := cltest.NewConfig(t)
	config.Set("GAS_UPDATER_ENABLED", "true")
	config.Set("GAS_UPDATER_BLOCK_DELAY", "0")
	config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "1")
	config.Set("ETH_GAS_PRICE_DEFAULT", 42)
	config.Set("ETH_MAX_GAS_PRICE_WEI", 100)
	store, cleanup := cltest.NewStoreWithConfig(config)
	config.SetRuntimeStore(store.ORM)
	defer cleanup()
	txm := new(mocks.TxManager)
	store.TxManager = txm
	gu := services.NewGasUpdater(store)

	txm.On("GetBlockByNumber", "0x0").Return(cltest.BlockWithTransactions(9001), nil)
	head := cltest.Head(0)
	gu.OnNewHead(head)
	assert.Len(t, gu.RollingBlockHistory(), 1)
	assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())

	txm.On("GetBlockByNumber", "0x1").Return(cltest.BlockWithTransactions(9002), nil)
	head = cltest.Head(1)
	gu.OnNewHead(head)

	assert.Equal(t, big.NewInt(42), config.EthGasPriceDefault())
}
