package remote_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

func (r *testReceiver) Receive(msg *remotetypes.MessageBody) {
	r.ch <- msg
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

	dispatcher := remote.NewDispatcher(wrapper, signer, registry, lggr)
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

	dispatcher := remote.NewDispatcher(wrapper, signer, registry, lggr)
	require.NoError(t, dispatcher.Start(ctx))

	rcv := newReceiver()
	err := dispatcher.SetReceiver(capId1, donId1, rcv)
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

	dispatcher := remote.NewDispatcher(wrapper, signer, registry, lggr)
	require.NoError(t, dispatcher.Start(ctx))

	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	responseDestPeerID := <-sendCh
	require.Equal(t, peerId1, responseDestPeerID)

	require.NoError(t, dispatcher.Close())
}
