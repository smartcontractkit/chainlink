package headtracker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func Test_HeadListener_HappyPath(t *testing.T) {
	// Logic:
	// - spawn a listener instance
	// - mock SubscribeNewHead/Err/Unsubscribe to track these calls
	// - send 3 heads
	// - ask listener to stop
	// Asserts:
	// - check Connected()/ReceivingHeads() are updated
	// - 3 heads is passed to callback
	// - ethClient methods are invoked

	lggr := logger.TestLogger(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// no need to test head timeouts here
		c.EVM[0].NoNewHeadsThreshold = &models.Duration{}
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	chStop := make(chan struct{})
	hl := headtracker.NewHeadListener(lggr, ethClient, evmcfg, chStop)

	var headCount atomic.Int32
	handler := func(context.Context, *evmtypes.Head) error {
		headCount.Add(1)
		return nil
	}

	subscribeAwaiter := cltest.NewAwaiter()
	unsubscribeAwaiter := cltest.NewAwaiter()
	var chHeads chan<- *evmtypes.Head
	var chErr = make(chan error)
	var chSubErr <-chan error = chErr
	sub := evmclimocks.NewSubscription(t)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *types.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
		chHeads = args.Get(1).(chan<- *evmtypes.Head)
		subscribeAwaiter.ItHappened()
	})
	sub.On("Err").Return(chSubErr)
	sub.On("Unsubscribe").Return().Once().Run(func(mock.Arguments) {
		unsubscribeAwaiter.ItHappened()
		close(chHeads)
		close(chErr)
	})

	doneAwaiter := cltest.NewAwaiter()
	done := func() {
		doneAwaiter.ItHappened()
	}
	go hl.ListenForNewHeads(handler, done)

	subscribeAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))
	require.Eventually(t, hl.Connected, testutils.WaitTimeout(t), testutils.TestInterval)

	chHeads <- cltest.Head(0)
	chHeads <- cltest.Head(1)
	chHeads <- cltest.Head(2)

	require.True(t, hl.ReceivingHeads())

	close(chStop)
	doneAwaiter.AwaitOrFail(t)

	unsubscribeAwaiter.AwaitOrFail(t)
	require.Equal(t, int32(3), headCount.Load())
}

func Test_HeadListener_NotReceivingHeads(t *testing.T) {
	// Logic:
	// - same as Test_HeadListener_HappyPath, but
	// - send one head, make sure ReceivingHeads() is true
	// - do not send any heads within BlockEmissionIdleWarningThreshold and check ReceivingHeads() is false

	lggr := logger.TestLogger(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NoNewHeadsThreshold = models.MustNewDuration(time.Second)
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	evmcfg.BlockEmissionIdleWarningThreshold()
	chStop := make(chan struct{})
	hl := headtracker.NewHeadListener(lggr, ethClient, evmcfg, chStop)

	firstHeadAwaiter := cltest.NewAwaiter()
	handler := func(context.Context, *evmtypes.Head) error {
		firstHeadAwaiter.ItHappened()
		return nil
	}

	subscribeAwaiter := cltest.NewAwaiter()
	var chHeads chan<- *evmtypes.Head
	var chErr = make(chan error)
	var chSubErr <-chan error = chErr
	sub := evmclimocks.NewSubscription(t)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *types.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
		chHeads = args.Get(1).(chan<- *evmtypes.Head)
		subscribeAwaiter.ItHappened()
	})
	sub.On("Err").Return(chSubErr)
	sub.On("Unsubscribe").Return().Once().Run(func(_ mock.Arguments) {
		close(chHeads)
		close(chErr)
	})

	doneAwaiter := cltest.NewAwaiter()
	done := func() {
		doneAwaiter.ItHappened()
	}
	go hl.ListenForNewHeads(handler, done)

	subscribeAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))

	chHeads <- cltest.Head(0)
	firstHeadAwaiter.AwaitOrFail(t)

	require.True(t, hl.ReceivingHeads())

	time.Sleep(time.Second * 2)

	require.False(t, hl.ReceivingHeads())

	close(chStop)
	doneAwaiter.AwaitOrFail(t)
}

func Test_HeadListener_SubscriptionErr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		closeErr bool
	}{
		{"nil error", nil, false},
		{"socket error", errors.New("close 1006 (abnormal closure): unexpected EOF"), false},
		{"close Err channel", nil, true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l := logger.TestLogger(t)
			ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
			cfg := configtest.NewGeneralConfig(t, nil)
			evmcfg := evmtest.NewChainScopedConfig(t, cfg)
			chStop := make(chan struct{})
			hl := headtracker.NewHeadListener(l, ethClient, evmcfg, chStop)

			hnhCalled := make(chan *evmtypes.Head)
			hnh := func(_ context.Context, header *evmtypes.Head) error {
				hnhCalled <- header
				return nil
			}
			doneAwaiter := cltest.NewAwaiter()
			done := doneAwaiter.ItHappened

			chSubErrTest := make(chan error)
			var chSubErr <-chan error = chSubErrTest
			sub := evmclimocks.NewSubscription(t)
			// sub.Err is called twice because we enter the select loop two times: once
			// initially and once again after exactly one head has been received
			sub.On("Err").Return(chSubErr).Twice()

			subscribeAwaiter := cltest.NewAwaiter()
			var headsCh chan<- *evmtypes.Head
			// Initial subscribe
			ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *types.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
				headsCh = args.Get(1).(chan<- *evmtypes.Head)
				subscribeAwaiter.ItHappened()
			})
			go func() {
				hl.ListenForNewHeads(hnh, done)
			}()

			// Put a head on the channel to ensure we test all code paths
			subscribeAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))
			head := cltest.Head(0)
			headsCh <- head

			h := <-hnhCalled
			assert.Equal(t, head, h)

			// Expect a call to unsubscribe on error
			sub.On("Unsubscribe").Once().Run(func(_ mock.Arguments) {
				close(headsCh)
				// geth guarantees that Unsubscribe closes the errors channel
				if !test.closeErr {
					close(chSubErrTest)
				}
			})
			// Expect a resubscribe
			chSubErrTest2 := make(chan error)
			var chSubErr2 <-chan error = chSubErrTest2
			sub2 := evmclimocks.NewSubscription(t)
			sub2.On("Err").Return(chSubErr2)
			subscribeAwaiter2 := cltest.NewAwaiter()

			var headsCh2 chan<- *evmtypes.Head
			ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *types.Head")).Return(sub2, nil).Once().Run(func(args mock.Arguments) {
				headsCh2 = args.Get(1).(chan<- *evmtypes.Head)
				subscribeAwaiter2.ItHappened()
			})

			// Sending test error
			if test.closeErr {
				close(chSubErrTest)
			} else {
				chSubErrTest <- test.err
			}

			// Wait for it to resubscribe
			subscribeAwaiter2.AwaitOrFail(t, testutils.WaitTimeout(t))

			head2 := cltest.Head(1)
			headsCh2 <- head2

			h2 := <-hnhCalled
			assert.Equal(t, head2, h2)

			// Second call to unsubscribe on close
			sub2.On("Unsubscribe").Once().Run(func(_ mock.Arguments) {
				close(headsCh2)
				// geth guarantees that Unsubscribe closes the errors channel
				close(chSubErrTest2)
			})
			close(chStop)
			doneAwaiter.AwaitOrFail(t)
		})
	}
}
