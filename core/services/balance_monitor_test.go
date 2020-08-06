package services_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"

	"github.com/pkg/errors"
)

func TestBalanceMonitor_Connect(t *testing.T) {
	var nilBigInt *big.Int

	t.Run("updates balance from nil for multiple keys", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		gethClient := new(mocks.GethClient)
		cltest.MockEthOnStore(t, store,
			eth.NewClientWith(nil, gethClient),
		)

		k0 := cltest.MustDefaultKey(t, store)
		k0Addr := k0.Address.Address()
		k1 := cltest.MustInsertRandomKey(t, store)
		k1Addr := k1.Address.Address()

		bm := services.NewBalanceMonitor(store)
		k0bal := big.NewInt(42)
		k1bal := big.NewInt(43)
		assert.Nil(t, bm.GetEthBalance(k0Addr))
		assert.Nil(t, bm.GetEthBalance(k1Addr))

		gethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)

		gethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		head := cltest.Head(0)

		// Do the thing
		bm.Connect(head)

		assert.Equal(t, k0bal, bm.GetEthBalance(k0Addr).ToInt())
		assert.Equal(t, k1bal, bm.GetEthBalance(k1Addr).ToInt())

		gethClient.AssertExpectations(t)
	})

	t.Run("handles nil head", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		gethClient := new(mocks.GethClient)
		cltest.MockEthOnStore(t, store,
			eth.NewClientWith(nil, gethClient),
		)

		k0 := cltest.MustDefaultKey(t, store)
		k0Addr := k0.Address.Address()

		bm := services.NewBalanceMonitor(store)
		k0bal := big.NewInt(42)

		gethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)

		// Do the thing
		bm.Connect(nil)

		assert.Equal(t, k0bal, bm.GetEthBalance(k0Addr).ToInt())

		gethClient.AssertExpectations(t)
	})

	t.Run("recovers on error", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		gethClient := new(mocks.GethClient)
		cltest.MockEthOnStore(t, store,
			eth.NewClientWith(nil, gethClient),
		)

		k0 := cltest.MustDefaultKey(t, store)
		k0Addr := k0.Address.Address()

		bm := services.NewBalanceMonitor(store)

		gethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(nil, errors.New("a little easter egg for the 4chan link marines error"))

		// Do the thing
		bm.Connect(nil)

		assert.Nil(t, bm.GetEthBalance(k0Addr))

		gethClient.AssertExpectations(t)
	})
}

func TestBalanceMonitor_OnNewLongestChain_UpdatesBalance(t *testing.T) {
	t.Run("updates balance for multiple keys", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		store.Config.Set("ETH_BALANCE_MONITOR_BLOCK_DELAY", 0)

		gethClient := new(mocks.GethClient)
		cltest.MockEthOnStore(t, store,
			eth.NewClientWith(nil, gethClient),
		)

		k0 := cltest.MustDefaultKey(t, store)
		k0Addr := k0.Address.Address()
		k1 := cltest.MustInsertRandomKey(t, store)
		k1Addr := k1.Address.Address()

		bm := services.NewBalanceMonitor(store)
		k0bal := big.NewInt(42)
		// Deliberately larger than a 64 bit unsigned integer to test overflow
		k1bal := big.NewInt(0)
		k1bal.SetString("19223372036854776000", 10)

		head := cltest.Head(0)

		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(head.Number)).Once().Return(k0bal, nil)
		gethClient.On("BalanceAt", mock.Anything, k1Addr, big.NewInt(head.Number)).Once().Return(k1bal, nil)

		// Do the thing
		bm.OnNewLongestChain(*head)

		assert.Equal(t, k0bal, bm.GetEthBalance(k0Addr).ToInt())
		assert.Equal(t, k1bal, bm.GetEthBalance(k1Addr).ToInt())

		// Do it again
		k0bal2 := big.NewInt(142)
		k1bal2 := big.NewInt(142)

		head = cltest.Head(1)

		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(head.Number)).Once().Return(k0bal2, nil)
		gethClient.On("BalanceAt", mock.Anything, k1Addr, big.NewInt(head.Number)).Once().Return(k1bal2, nil)

		bm.OnNewLongestChain(*head)

		assert.Equal(t, k0bal2, bm.GetEthBalance(k0Addr).ToInt())
		assert.Equal(t, k1bal2, bm.GetEthBalance(k1Addr).ToInt())

		gethClient.AssertExpectations(t)
	})

	t.Run("lags behind by ETH_BALANCE_MONITOR_BLOCK_DELAY blocks", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		store.Config.Set("ETH_BALANCE_MONITOR_BLOCK_DELAY", 2)
		k0 := cltest.MustDefaultKey(t, store)
		k0Addr := k0.Address.Address()

		gethClient := new(mocks.GethClient)
		store.EthClient = eth.NewClientWith(nil, gethClient)
		bm := services.NewBalanceMonitor(store)

		// If lagged head would be negative, just uses 0
		k0bal := big.NewInt(42)
		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(0)).Once().Return(k0bal, nil)
		head := cltest.Head(0)
		bm.OnNewLongestChain(*head)

		assert.Equal(t, k0bal, bm.GetEthBalance(k0Addr).ToInt())

		// If lagged head would be negative, just uses 0
		k0bal = big.NewInt(43)
		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(0)).Once().Return(k0bal, nil)
		head = cltest.Head(1)
		bm.OnNewLongestChain(*head)

		// If lagged head is exactly 0, uses 0
		k0bal = big.NewInt(44)
		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(0)).Once().Return(k0bal, nil)
		head = cltest.Head(2)
		bm.OnNewLongestChain(*head)

		// If lagged head is positive, uses it
		k0bal = big.NewInt(44)
		gethClient.On("BalanceAt", mock.Anything, k0Addr, big.NewInt(1)).Once().Return(k0bal, nil)
		head = cltest.Head(3)
		bm.OnNewLongestChain(*head)

		gethClient.AssertExpectations(t)
	})
}
