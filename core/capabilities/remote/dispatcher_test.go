package remote_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

type testReceiver struct {
	ch chan *remotetypes.MessageBody
}

func newReceiver() *testReceiver {
	return &testReceiver{
		ch: make(chan *remotetypes.MessageBody, 100),
	}
}

func (r *testReceiver) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	r.ch <- msg
}

type testRateLimitConfig struct {
	globalRPS   float64
	globalBurst int
	rps         float64
	burst       int
}

func (c testRateLimitConfig) GlobalRPS() float64 {
	return c.globalRPS
}

func (c testRateLimitConfig) GlobalBurst() int {
	return c.globalBurst
}

func (c testRateLimitConfig) PerSenderRPS() float64 {
	return c.rps
}

func (c testRateLimitConfig) PerSenderBurst() int {
	return c.burst
}

type testConfig struct {
	supportedVersion   int
	receiverBufferSize int
	rateLimit          testRateLimitConfig
}

func (c testConfig) SupportedVersion() int {
	return c.supportedVersion
}

func (c testConfig) ReceiverBufferSize() int {
	return c.receiverBufferSize
}

func (c testConfig) RateLimit() config.DispatcherRateLimit {
	return c.rateLimit
}

func (c testConfig) SendToSharedPeer() bool {
	return true
}

func TestDispatcher_CleanStartClose(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("Receive", mock.Anything).Return(make(<-chan p2ptypes.Message))
	sharedPeer.On("ID", mock.Anything).Return(p2ptypes.PeerID{})
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(), sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_Receive(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	sharedPeer := mocks.NewSharedPeer(t)
	recvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return(nil, errors.New("not implemented"))
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(), sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	rcv := newReceiver()
	err = dispatcher.SetReceiver(capID1, donID1, rcv)
	require.NoError(t, err)

	// supported capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID2, donID1, []byte(payload1))
	// sender doesn't match
	invalid := encodeAndSign(t, privKey1, peerID1, peerID2, capID2, donID1, []byte(payload1))
	invalid.Sender = peerID2
	recvCh <- invalid
	// supported capability again
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload2))

	m := <-rcv.ch
	require.Equal(t, payload1, string(m.Payload))
	m = <-rcv.ch
	require.Equal(t, payload2, string(m.Payload))

	dispatcher.RemoveReceiver(capID1, donID1)
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_ReceiveForMethod(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	sharedPeer := mocks.NewSharedPeer(t)
	recvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return(nil, errors.New("not implemented"))
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(testConfig{
		supportedVersion:   1,
		receiverBufferSize: 10000,
		rateLimit: testRateLimitConfig{
			globalRPS:   800.0,
			globalBurst: 100,
			rps:         10.0,
			burst:       50,
		},
	}, sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	methodA, methodB := "methodA", "methodB"
	rcvA, rcvB := newReceiver(), newReceiver()
	require.NoError(t, dispatcher.SetReceiverForMethod(capID1, donID1, methodA, rcvA))
	require.NoError(t, dispatcher.SetReceiverForMethod(capID1, donID1, methodB, rcvB))

	// supported capability / methodA
	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID1, methodA, donID1, []byte(payload1))
	// unknown capability
	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID2, methodA, donID1, []byte(payload1))
	// supported capability / methodB
	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID1, methodB, donID1, []byte(payload2))

	m := <-rcvA.ch
	require.Equal(t, payload1, string(m.Payload))
	m = <-rcvB.ch
	require.Equal(t, payload2, string(m.Payload))

	dispatcher.RemoveReceiverForMethod(capID1, donID1, methodA)
	dispatcher.RemoveReceiverForMethod(capID1, donID1, methodB)
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_RespondWithError(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	sharedPeer := mocks.NewSharedPeer(t)
	recvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	sendCh := make(chan p2ptypes.PeerID)
	sharedPeer.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		peerID := args.Get(0).(p2ptypes.PeerID)
		sendCh <- peerID
	}).Return(nil)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return([]byte{1, 2, 3}, nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(), sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	responseDestPeerID := <-sendCh
	require.Equal(t, peerID1, responseDestPeerID)

	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_Send(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	_, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return([]byte("signed payload"), nil)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeerRecvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(sharedPeerRecvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	sharedPeer.On("Send", mock.Anything, mock.Anything).Return(nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(), sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	require.NoError(t, dispatcher.Send(peerID1, &remotetypes.MessageBody{}))
	// mocks expect Sign() and Send()

	require.NoError(t, dispatcher.Close())
}

// panicOnFirstReceiver panics on the first Receive call and records all
// subsequent messages so tests can assert the goroutine survived.
type panicOnFirstReceiver struct {
	mu       sync.Mutex
	calls    int
	received chan *remotetypes.MessageBody
}

func newPanicOnFirstReceiver() *panicOnFirstReceiver {
	return &panicOnFirstReceiver{received: make(chan *remotetypes.MessageBody, 10)}
}

func (r *panicOnFirstReceiver) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	r.mu.Lock()
	r.calls++
	call := r.calls
	r.mu.Unlock()
	if call == 1 {
		panic("deliberate panic from receiver")
	}
	r.received <- msg
}

// TestDispatcher_ReceiverPanicDoesNotKillLoop verifies that a panic inside a
// receiver's Receive() method is caught by the recover() wrapper in the
// dispatcher goroutine and does not prevent subsequent messages from being
// delivered to the same receiver.
func TestDispatcher_ReceiverPanicDoesNotKillLoop(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := t.Context()
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	sharedPeer := mocks.NewSharedPeer(t)
	recvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(), sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	rcv := newPanicOnFirstReceiver()
	err = dispatcher.SetReceiver(capID1, donID1, rcv)
	require.NoError(t, err)

	// First message triggers the panic; the goroutine must survive.
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	// Second message must still be delivered.
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload2))

	m := <-rcv.received
	require.Equal(t, payload2, string(m.Payload))

	dispatcher.RemoveReceiver(capID1, donID1)
	require.NoError(t, dispatcher.Close())
}

func newTestConfig() testConfig {
	return testConfig{
		supportedVersion:   1,
		receiverBufferSize: 10000,
		rateLimit: testRateLimitConfig{
			globalRPS:   800.0,
			globalBurst: 100,
			rps:         10.0,
			burst:       50,
		},
	}
}
