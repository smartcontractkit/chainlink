package remote_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

func TestDispatcher_CleanStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	peer := mocks.NewPeer(t)
	recvCh := make(<-chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return(recvCh)
	peer.On("ID", mock.Anything).Return(p2ptypes.PeerID{})
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
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
	}, wrapper, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_Receive(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	privKey1, peerId1 := newKeyPair(t)
	_, peerId2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerId2)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.On("Sign", mock.Anything).Return(nil, errors.New("not implemented"))
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
	}, wrapper, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	rcv := newReceiver()
	err = dispatcher.SetReceiver(capId1, donId1, rcv)
	require.NoError(t, err)

	// supported capability
	recvCh <- encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerId1, peerId2, capId2, donId1, []byte(payload1))
	// sender doesn't match
	invalid := encodeAndSign(t, privKey1, peerId1, peerId2, capId2, donId1, []byte(payload1))
	invalid.Sender = peerId2
	recvCh <- invalid
	// supported capability again
	recvCh <- encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload2))

	m := <-rcv.ch
	require.Equal(t, payload1, string(m.Payload))
	m = <-rcv.ch
	require.Equal(t, payload2, string(m.Payload))

	dispatcher.RemoveReceiver(capId1, donId1)
	require.NoError(t, dispatcher.Close())
}

func TestDispatcher_RespondWithError(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	privKey1, peerId1 := newKeyPair(t)
	_, peerId2 := newKeyPair(t)

	peer := mocks.NewPeer(t)
	recvCh := make(chan p2ptypes.Message)
	peer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	peer.On("ID", mock.Anything).Return(peerId2)
	sendCh := make(chan p2ptypes.PeerID)
	peer.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		peerID := args.Get(0).(p2ptypes.PeerID)
		sendCh <- peerID
	}).Return(nil)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	signer := mocks.NewSigner(t)
	signer.On("Sign", mock.Anything).Return([]byte{}, nil)
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
	}, wrapper, signer, registry, lggr)
	require.NoError(t, err)
	require.NoError(t, dispatcher.Start(ctx))

	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	responseDestPeerID := <-sendCh
	require.Equal(t, peerId1, responseDestPeerID)

	require.NoError(t, dispatcher.Close())
}
