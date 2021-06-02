package services_test

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"

	"github.com/pkg/errors"
)

var nilBigInt *big.Int

func TestBalanceMonitor_Connect(t *testing.T) {
	t.Run("updates balance from nil for multiple keys", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		ethClient := new(mocks.Client)
		defer ethClient.AssertExpectations(t)
		store.EthClient = ethClient

		_, k0Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
		_, k1Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

		bm := services.NewBalanceMonitor(store, ethKeyStore)
		defer bm.Close()

		k0bal := big.NewInt(42)
		k1bal := big.NewInt(43)
		assert.Nil(t, bm.GetEthBalance(k0Addr))
		assert.Nil(t, bm.GetEthBalance(k1Addr))

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		head := cltest.Head(0)

		// Do the thing
		bm.Connect(head)

		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal))
		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k1Addr).ToInt()
		}).Should(gomega.Equal(k1bal))
	})

	t.Run("handles nil head", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		ethClient := new(mocks.Client)
		defer ethClient.AssertExpectations(t)
		store.EthClient = ethClient

		_, k0Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

		bm := services.NewBalanceMonitor(store, ethKeyStore)
		defer bm.Close()
		k0bal := big.NewInt(42)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)

		// Do the thing
		bm.Connect(nil)

		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal))
	})

	t.Run("recovers on error", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		ethClient := new(mocks.Client)
		defer ethClient.AssertExpectations(t)
		store.EthClient = ethClient

		_, k0Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

		bm := services.NewBalanceMonitor(store, ethKeyStore)
		defer bm.Close()

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).
			Once().
			Return(nil, errors.New("a little easter egg for the 4chan link marines error"))

		// Do the thing
		bm.Connect(nil)

		gomega.NewGomegaWithT(t).Consistently(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.BeNil())
	})
}

func TestBalanceMonitor_OnNewLongestChain_UpdatesBalance(t *testing.T) {
	t.Run("updates balance for multiple keys", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

		ethClient := new(mocks.Client)
		defer ethClient.AssertExpectations(t)
		store.EthClient = ethClient

		_, k0Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
		_, k1Addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

		bm := services.NewBalanceMonitor(store, ethKeyStore)
		defer bm.Close()
		k0bal := big.NewInt(42)
		// Deliberately larger than a 64 bit unsigned integer to test overflow
		k1bal := big.NewInt(0)
		k1bal.SetString("19223372036854776000", 10)

		head := cltest.Head(0)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		// Do the thing
		bm.OnNewLongestChain(context.TODO(), *head)

		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal))
		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k1Addr).ToInt()
		}).Should(gomega.Equal(k1bal))

		// Do it again
		k0bal2 := big.NewInt(142)
		k1bal2 := big.NewInt(142)

		head = cltest.Head(1)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal2, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal2, nil)

		bm.OnNewLongestChain(context.TODO(), *head)

		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal2))
		gomega.NewGomegaWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k1Addr).ToInt()
		}).Should(gomega.Equal(k1bal2))
	})
}

func TestBalanceMonitor_FewerRPCCallsWhenBehind(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	ethClient := new(mocks.Client)
	ethClient.AssertExpectations(t)
	store.EthClient = ethClient

	bm := services.NewBalanceMonitor(store, ethKeyStore)

	head := cltest.Head(0)

	// Only expect this twice, even though 10 heads will come in
	mockUnblocker := make(chan time.Time)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).
		WaitUntil(mockUnblocker).
		Once().
		Return(big.NewInt(42), nil)
	// This second call is Maybe because the SleeperTask may not have started
	// before we call `OnNewLongestChain` 10 times, in which case it's only
	// executed once
	var callCount int32
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).
		Run(func(mock.Arguments) { atomic.AddInt32(&callCount, 1) }).
		Maybe().
		Return(big.NewInt(42), nil)

	// Do the thing multiple times
	for i := 0; i < 10; i++ {
		bm.OnNewLongestChain(context.TODO(), *head)
	}

	// Unblock the first mock
	cltest.CallbackOrTimeout(t, "FewerRPCCallsWhenBehind unblock BalanceAt", func() {
		mockUnblocker <- time.Time{}
	})

	bm.Close()

	// Make sure the BalanceAt mock wasn't called more than once
	assert.LessOrEqual(t, atomic.LoadInt32(&callCount), int32(1))
}
