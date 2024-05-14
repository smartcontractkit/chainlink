package monitor_test

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/monitor"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var nilBigInt *big.Int

func newEthClientMock(t *testing.T) *evmclimocks.Client {
	mockEth := evmclimocks.NewClient(t)
	mockEth.On("ConfiguredChainID").Maybe().Return(big.NewInt(0))
	return mockEth
}

func TestBalanceMonitor_Start(t *testing.T) {
	t.Parallel()

	t.Run("updates balance from nil for multiple keys", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		ethClient := newEthClientMock(t)
		_, k1Addr := cltest.MustInsertRandomKey(t, ethKeyStore)
		_, k0Addr := cltest.MustInsertRandomKey(t, ethKeyStore)

		bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))

		k0bal := big.NewInt(42)
		k1bal := big.NewInt(43)
		assert.Nil(t, bm.GetEthBalance(k0Addr))
		assert.Nil(t, bm.GetEthBalance(k1Addr))

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		servicetest.RunHealthy(t, bm)

		gomega.NewWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal))
		gomega.NewWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k1Addr).ToInt()
		}).Should(gomega.Equal(k1bal))
	})

	t.Run("handles nil head", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		ethClient := newEthClientMock(t)

		_, k0Addr := cltest.MustInsertRandomKey(t, ethKeyStore)

		bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))
		k0bal := big.NewInt(42)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)

		servicetest.RunHealthy(t, bm)

		gomega.NewWithT(t).Eventually(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.Equal(k0bal))
	})

	t.Run("cancelled context", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		ethClient := newEthClientMock(t)

		_, k0Addr := cltest.MustInsertRandomKey(t, ethKeyStore)

		bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))
		ctxCancelledAwaiter := cltest.NewAwaiter()

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			select {
			case <-time.After(testutils.WaitTimeout(t)):
			case <-ctx.Done():
				ctxCancelledAwaiter.ItHappened()
			}
		}).Return(nil, nil)

		ctx, cancel := context.WithCancel(testutils.Context(t))
		go func() {
			<-time.After(time.Second)
			cancel()
		}()
		assert.NoError(t, bm.Start(ctx))

		ctxCancelledAwaiter.AwaitOrFail(t)
	})

	t.Run("recovers on error", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		ethClient := newEthClientMock(t)

		_, k0Addr := cltest.MustInsertRandomKey(t, ethKeyStore)

		bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).
			Once().
			Return(nil, pkgerrors.New("a little easter egg for the 4chan link marines error"))

		servicetest.RunHealthy(t, bm)

		gomega.NewWithT(t).Consistently(func() *big.Int {
			return bm.GetEthBalance(k0Addr).ToInt()
		}).Should(gomega.BeNil())
	})
}

func TestBalanceMonitor_OnNewLongestChain_UpdatesBalance(t *testing.T) {
	t.Parallel()

	t.Run("updates balance for multiple keys", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		ethClient := newEthClientMock(t)

		_, k0Addr := cltest.MustInsertRandomKey(t, ethKeyStore)
		_, k1Addr := cltest.MustInsertRandomKey(t, ethKeyStore)

		bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))
		k0bal := big.NewInt(42)
		// Deliberately larger than a 64 bit unsigned integer to test overflow
		k1bal := big.NewInt(0)
		k1bal.SetString("19223372036854776000", 10)

		head := cltest.Head(0)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		servicetest.RunHealthy(t, bm)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal, nil)

		// Do the thing
		bm.OnNewLongestChain(testutils.Context(t), head)

		<-bm.WorkDone()
		assert.Equal(t, k0bal, bm.GetEthBalance(k0Addr).ToInt())
		assert.Equal(t, k1bal, bm.GetEthBalance(k1Addr).ToInt())

		// Do it again
		k0bal2 := big.NewInt(142)
		k1bal2 := big.NewInt(142)

		head = cltest.Head(1)

		ethClient.On("BalanceAt", mock.Anything, k0Addr, nilBigInt).Once().Return(k0bal2, nil)
		ethClient.On("BalanceAt", mock.Anything, k1Addr, nilBigInt).Once().Return(k1bal2, nil)

		bm.OnNewLongestChain(testutils.Context(t), head)

		<-bm.WorkDone()
		assert.Equal(t, k0bal2, bm.GetEthBalance(k0Addr).ToInt())
		assert.Equal(t, k1bal2, bm.GetEthBalance(k1Addr).ToInt())
	})
}

func TestBalanceMonitor_FewerRPCCallsWhenBehind(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	cltest.MustInsertRandomKey(t, ethKeyStore)

	ethClient := newEthClientMock(t)

	bm := monitor.NewBalanceMonitor(ethClient, ethKeyStore, logger.Test(t))
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(big.NewInt(1), nil)
	servicetest.RunHealthy(t, bm)

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
	var callCount atomic.Int32
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).
		Run(func(mock.Arguments) { callCount.Add(1) }).
		Maybe().
		Return(big.NewInt(42), nil)

	// Do the thing multiple times
	for i := 0; i < 10; i++ {
		bm.OnNewLongestChain(testutils.Context(t), head)
	}

	// Unblock the first mock
	cltest.CallbackOrTimeout(t, "FewerRPCCallsWhenBehind unblock BalanceAt", func() {
		mockUnblocker <- time.Time{}
	})

	// Make sure the BalanceAt mock wasn't called more than once
	assert.LessOrEqual(t, callCount.Load(), int32(1))
}

func Test_ApproximateFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		want      float64
		wantError bool
	}{
		{"zero", "0", 0, false},
		{"small", "1", 0.000000000000000001, false},
		{"rounding", "12345678901234567890", 12.345678901234567, false},
		{"large", "123456789012345678901234567890", 123456789012.34567, false},
		{"extreme", "1234567890123456789012345678901234567890123456789012345678901234567890", 1.2345678901234568e+51, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eth := assets.NewEth(0)
			eth.SetString(test.input, 10)
			float, err := monitor.ApproximateFloat64(eth)
			require.NoError(t, err)
			require.Equal(t, test.want, float)
		})
	}
}
