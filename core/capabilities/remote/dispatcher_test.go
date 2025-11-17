package remote_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
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
	sendToSharedPeer   bool
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
	return c.sendToSharedPeer
}

func TestDispatcher_CleanStartClose(t *testing.T) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	peer := mocks.NewPeer(t)
	recvCh := make(<-chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return(recvCh)
	peer.On("ID", mock.Anything).Return(p2ptypes.PeerID{})
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(false), wrapper, nil, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_Receive(t *testing.T) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerID2)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return(nil, errors.New("not implemented"))
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(false), wrapper, nil, signer, registry, lggr)
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
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerID2)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
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
	}, wrapper, nil, signer, registry, lggr)
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
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerID2)
	sendCh := make(chan p2ptypes.PeerID)
	peer.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		peerID := args.Get(0).(p2ptypes.PeerID)
		sendCh <- peerID
	}).Return(nil)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return([]byte{1, 2, 3}, nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(false), wrapper, nil, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	responseDestPeerID := <-sendCh
	require.Equal(t, peerID1, responseDestPeerID)

	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_ReceiveFromBothPeers(t *testing.T) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerID2)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeerRecvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(sharedPeerRecvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(false), wrapper, sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	rcv := newReceiver()
	err = dispatcher.SetReceiver(capID1, donID1, rcv)
	require.NoError(t, err)

	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	sharedPeerRecvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload2))
	close(sharedPeerRecvCh) // make sure Dispatcher handles SharedPeer shutdown gracefully

	m := <-rcv.ch
	require.Equal(t, payload1, string(m.Payload))
	m = <-rcv.ch
	require.Equal(t, payload2, string(m.Payload))

	dispatcher.RemoveReceiver(capID1, donID1)
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_SendToSharedPeer(t *testing.T) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	_, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerID2)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Initialize().Return(nil)
	signer.EXPECT().Sign(mock.Anything).Return([]byte("signed payload"), nil)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeerRecvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(sharedPeerRecvCh))
	sharedPeer.On("ID", mock.Anything).Return(peerID2)
	sharedPeer.On("Send", mock.Anything, mock.Anything).Return(nil)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	dispatcher, err := remote.NewDispatcher(newTestConfig(true), wrapper, sharedPeer, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	require.NoError(t, dispatcher.Send(peerID1, &remotetypes.MessageBody{}))
	// mocks expect Sign() and Send()

	require.NoError(t, dispatcher.Close())
}

func newTestConfig(sendToSharedPeer bool) testConfig {
	return testConfig{
		supportedVersion:   1,
		receiverBufferSize: 10000,
		rateLimit: testRateLimitConfig{
			globalRPS:   800.0,
			globalBurst: 100,
			rps:         10.0,
			burst:       50,
		},
		sendToSharedPeer: sendToSharedPeer,
	}
}
