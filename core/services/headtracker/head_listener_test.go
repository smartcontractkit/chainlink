package headtracker_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
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
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	stopCh := make(chan struct{})
	hl := headtracker.NewHeadListener(lggr, ethClient, evmcfg, stopCh)

	var headCount atomic.Int32
	handler := func(ctx context.Context, header *eth.Head) error {
		headCount.Inc()
		return nil
	}

	subscribeAwaiter := cltest.NewAwaiter()
	unsubscribeAwaiter := cltest.NewAwaiter()
	var headsCh chan<- *eth.Head
	var errCh = make(chan error)
	var subErrCh <-chan error = errCh
	sub := new(ethmocks.Subscription)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *eth.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
		headsCh = args.Get(1).(chan<- *eth.Head)
		subscribeAwaiter.ItHappened()
	})
	sub.On("Err").Return(subErrCh)
	sub.On("Unsubscribe").Return().Once().Run(func(args mock.Arguments) {
		unsubscribeAwaiter.ItHappened()
		close(headsCh)
		close(errCh)
	})

	doneAwaiter := cltest.NewAwaiter()
	done := func() {
		doneAwaiter.ItHappened()
	}
	go hl.ListenForNewHeads(handler, done)

	subscribeAwaiter.AwaitOrFail(t)
	require.True(t, hl.Connected())

	headsCh <- cltest.Head(0)
	headsCh <- cltest.Head(1)
	headsCh <- cltest.Head(2)

	require.True(t, hl.ReceivingHeads())

	close(stopCh)
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
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	cfg := cltest.NewTestGeneralConfig(t)
	idleDuration := time.Second
	cfg.Overrides.GlobalBlockEmissionIdleWarningThreshold = &idleDuration
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	evmcfg.BlockEmissionIdleWarningThreshold()
	stopCh := make(chan struct{})
	hl := headtracker.NewHeadListener(lggr, ethClient, evmcfg, stopCh)

	firstHeadAwaiter := cltest.NewAwaiter()
	handler := func(ctx context.Context, header *eth.Head) error {
		firstHeadAwaiter.ItHappened()
		return nil
	}

	subscribeAwaiter := cltest.NewAwaiter()
	var headsCh chan<- *eth.Head
	var errCh = make(chan error)
	var subErrCh <-chan error = errCh
	sub := new(ethmocks.Subscription)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *eth.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
		headsCh = args.Get(1).(chan<- *eth.Head)
		subscribeAwaiter.ItHappened()
	})
	sub.On("Err").Return(subErrCh)
	sub.On("Unsubscribe").Return().Once().Run(func(_ mock.Arguments) {
		close(headsCh)
		close(errCh)
	})

	doneAwaiter := cltest.NewAwaiter()
	done := func() {
		doneAwaiter.ItHappened()
	}
	go hl.ListenForNewHeads(handler, done)

	subscribeAwaiter.AwaitOrFail(t)

	headsCh <- cltest.Head(0)
	firstHeadAwaiter.AwaitOrFail(t)

	require.True(t, hl.ReceivingHeads())

	time.Sleep(idleDuration * 2)

	require.False(t, hl.ReceivingHeads())

	close(stopCh)
	doneAwaiter.AwaitOrFail(t)
}

func Test_HeadListener_ResubscribesIfWSClosed(t *testing.T) {
	l := logger.TestLogger(t)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	chStop := make(chan struct{})
	hl := headtracker.NewHeadListener(l, ethClient, evmcfg, chStop)

	hnhCalled := make(chan *eth.Head)
	hnh := func(ctx context.Context, header *eth.Head) error {
		hnhCalled <- header
		return nil
	}
	doneAwaiter := cltest.NewAwaiter()
	done := func() {
		doneAwaiter.ItHappened()
	}

	chSubErrTest := make(chan error)
	var chSubErr <-chan error = chSubErrTest
	sub := new(ethmocks.Subscription)
	// sub.Err is called twice because we enter the select loop two times: once
	// initially and once again after exactly one head has been received
	sub.On("Err").Return(chSubErr).Twice()

	subscribeAwaiter := cltest.NewAwaiter()
	var headsCh chan<- *eth.Head
	// Initial subscribe
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *eth.Head")).Return(sub, nil).Once().Run(func(args mock.Arguments) {
		headsCh = args.Get(1).(chan<- *eth.Head)
		subscribeAwaiter.ItHappened()
	})
	go func() {
		hl.ListenForNewHeads(hnh, done)
	}()

	// Put a head on the channel to ensure we test all code paths
	subscribeAwaiter.AwaitOrFail(t)
	head := cltest.Head(0)
	headsCh <- head

	h := <-hnhCalled
	assert.Equal(t, head, h)

	// Expect a call to unsubscribe on error
	sub.On("Unsubscribe").Once().Run(func(_ mock.Arguments) {
		close(headsCh)
		// geth guarantees that Unsubscribe closes the errors channel
		close(chSubErrTest)
	})
	// Expect a resubscribe
	chSubErrTest2 := make(chan error)
	var chSubErr2 <-chan error = chSubErrTest2
	sub2 := new(ethmocks.Subscription)
	sub2.On("Err").Return(chSubErr2)
	subscribeAwaiter2 := cltest.NewAwaiter()

	var headsCh2 chan<- *eth.Head
	ethClient.On("SubscribeNewHead", mock.Anything, mock.AnythingOfType("chan<- *eth.Head")).Return(sub2, nil).Once().Run(func(args mock.Arguments) {
		headsCh2 = args.Get(1).(chan<- *eth.Head)
		subscribeAwaiter2.ItHappened()
	})

	// Simulate websocket error/close
	chSubErrTest <- errors.New("close 1006 (abnormal closure): unexpected EOF")

	// Wait for it to resubscribe
	subscribeAwaiter2.AwaitOrFail(t)

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
}
