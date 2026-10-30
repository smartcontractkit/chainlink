package remote_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
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

type testConfig struct {
	receiverBufferSize int
}

func (c testConfig) ReceiverBufferSize() int {
	return c.receiverBufferSize
}

func newTestConfig() testConfig {
	return testConfig{receiverBufferSize: 10000}
}

// startDispatcher starts a dispatcher with the given local peer ID. Messages
// written to the returned channel are delivered as if received from the network.
func startDispatcher(t *testing.T, cfg testConfig, localID p2ptypes.PeerID, signer *mocks.Signer, sharedPeer *mocks.SharedPeer) (remotetypes.Dispatcher, chan<- p2ptypes.Message) {
	lggr := logger.Test(t)
	recvCh := make(chan p2ptypes.Message)
	sharedPeer.On("Receive", mock.Anything).Return((<-chan p2ptypes.Message)(recvCh))
	sharedPeer.On("ID", mock.Anything).Return(localID)
	signer.EXPECT().Initialize().Return(nil)

	d, err := remote.NewDispatcher(cfg, sharedPeer, signer, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	require.NoError(t, d.Start(t.Context()))
	return d, recvCh
}

func TestDispatcher_CleanStartClose(t *testing.T) {
	t.Parallel()

	d, _ := startDispatcher(t, newTestConfig(), p2ptypes.PeerID{}, mocks.NewSigner(t), mocks.NewSharedPeer(t))
	require.NoError(t, d.Close())
}

func TestDispatcher_Receive(t *testing.T) {
	t.Parallel()

	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Sign(mock.Anything).Return(nil, errors.New("not implemented"))
	d, recvCh := startDispatcher(t, newTestConfig(), peerID2, signer, mocks.NewSharedPeer(t))

	rcv := newReceiver()
	require.NoError(t, d.SetReceiver(capID1, donID1, rcv))

	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID2, donID1, []byte("unknown capability"))
	// known capability, but the sender doesn't match the signed body
	invalid := encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte("sender mismatch"))
	invalid.Sender = peerID2
	recvCh <- invalid
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload2))

	require.Equal(t, payload1, string((<-rcv.ch).Payload))
	require.Equal(t, payload2, string((<-rcv.ch).Payload))

	d.RemoveReceiver(capID1, donID1)
	require.NoError(t, d.Close())
	require.Empty(t, rcv.ch)
}

func TestDispatcher_ReceiveForMethod(t *testing.T) {
	t.Parallel()

	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Sign(mock.Anything).Return(nil, errors.New("not implemented"))
	d, recvCh := startDispatcher(t, newTestConfig(), peerID2, signer, mocks.NewSharedPeer(t))

	methodA, methodB := "methodA", "methodB"
	rcvA, rcvB := newReceiver(), newReceiver()
	require.NoError(t, d.SetReceiverForMethod(capID1, donID1, methodA, rcvA))
	require.NoError(t, d.SetReceiverForMethod(capID1, donID1, methodB, rcvB))
	require.ErrorIs(t, d.SetReceiverForMethod(capID1, donID1, methodA, newReceiver()), remote.ErrReceiverExists)

	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID1, methodA, donID1, []byte(payload1))
	// unknown capability
	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID2, methodA, donID1, []byte(payload1))
	recvCh <- encodeAndSignForMethod(t, privKey1, peerID1, peerID2, capID1, methodB, donID1, []byte(payload2))

	require.Equal(t, payload1, string((<-rcvA.ch).Payload))
	require.Equal(t, payload2, string((<-rcvB.ch).Payload))

	d.RemoveReceiverForMethod(capID1, donID1, methodA)
	d.RemoveReceiverForMethod(capID1, donID1, methodB)
	require.NoError(t, d.Close())
	require.Empty(t, rcvA.ch)
	require.Empty(t, rcvB.ch)
}

func TestDispatcher_RespondWithError(t *testing.T) {
	t.Parallel()

	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Sign(mock.Anything).Return([]byte{1, 2, 3}, nil)
	sharedPeer := mocks.NewSharedPeer(t)
	type sent struct {
		to  p2ptypes.PeerID
		raw []byte
	}
	sendCh := make(chan sent, 10)
	sharedPeer.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		sendCh <- sent{to: args.Get(0).(p2ptypes.PeerID), raw: args.Get(1).([]byte)}
	}).Return(nil)
	d, recvCh := startDispatcher(t, newTestConfig(), peerID2, signer, sharedPeer)

	// a message that already carries an error is not answered, to avoid error loops
	errMsg := signBody(t, privKey1, &remotetypes.MessageBody{
		Sender: peerID1[:], Receiver: peerID2[:], CapabilityId: capID1, CapabilityDonId: donID1,
		Error: remotetypes.Error_INTERNAL_ERROR,
	}, peerID1)
	recvCh <- errMsg
	// unknown capability
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))

	// Sends are synchronous in the receive loop, so a first response for the
	// unknown capability proves the error message produced none.
	resp := <-sendCh
	require.Equal(t, peerID1, resp.to)
	body := decodeSentBody(t, resp.raw)
	require.Equal(t, remotetypes.Error_CAPABILITY_NOT_FOUND, body.Error)
	require.Equal(t, capID1, body.CapabilityId)
	require.Empty(t, body.Payload)

	require.NoError(t, d.Close())
	require.Empty(t, sendCh)
}

func TestDispatcher_Send(t *testing.T) {
	t.Parallel()

	_, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	signer := mocks.NewSigner(t)
	signer.EXPECT().Sign(mock.Anything).Return([]byte("signed payload"), nil)
	sharedPeer := mocks.NewSharedPeer(t)
	var rawSent []byte
	sharedPeer.On("Send", peerID1, mock.Anything).Run(func(args mock.Arguments) {
		rawSent = args.Get(1).([]byte)
	}).Return(nil)
	d, _ := startDispatcher(t, newTestConfig(), peerID2, signer, sharedPeer)

	require.NoError(t, d.Send(peerID1, &remotetypes.MessageBody{Method: "method"}))
	require.NoError(t, d.Close())

	var msg remotetypes.Message
	require.NoError(t, proto.Unmarshal(rawSent, &msg))
	require.Equal(t, []byte("signed payload"), msg.Signature)
	body := decodeSentBody(t, rawSent)
	require.Equal(t, "method", body.Method)
	require.Equal(t, peerID2[:], body.Sender)
	require.Equal(t, peerID1[:], body.Receiver)
	require.NotZero(t, body.Version)
	require.NotZero(t, body.Timestamp)
}

func TestDispatcher_DropsWhenReceiverBufferFull(t *testing.T) {
	t.Parallel()

	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	d, recvCh := startDispatcher(t, testConfig{receiverBufferSize: 1}, peerID2, mocks.NewSigner(t), mocks.NewSharedPeer(t))

	rcv := &blockingReceiver{entered: make(chan struct{}), release: make(chan struct{}), received: make(chan string, 10)}
	require.NoError(t, d.SetReceiver(capID1, donID1, rcv))
	send := func(payload string) {
		recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload))
	}

	send("handled")
	<-rcv.entered // receiver is now blocked, the buffer is empty
	send("buffered")
	send("dropped")
	// the dispatcher loop is single-threaded: once this is accepted, "dropped" has been handled
	send("dropped too")
	close(rcv.release)

	require.Equal(t, "handled", <-rcv.received)
	require.Equal(t, "buffered", <-rcv.received)
	require.NoError(t, d.Close())
	require.Empty(t, rcv.received)
}

// blockingReceiver blocks inside Receive of the first message until released.
type blockingReceiver struct {
	once     sync.Once
	entered  chan struct{}
	release  chan struct{}
	received chan string
}

func (r *blockingReceiver) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	r.once.Do(func() {
		close(r.entered)
		<-r.release
	})
	r.received <- string(msg.Payload)
}

func decodeSentBody(t *testing.T, raw []byte) *remotetypes.MessageBody {
	var msg remotetypes.Message
	require.NoError(t, proto.Unmarshal(raw, &msg))
	var body remotetypes.MessageBody
	require.NoError(t, proto.Unmarshal(msg.Body, &body))
	return &body
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

	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	d, recvCh := startDispatcher(t, newTestConfig(), peerID2, mocks.NewSigner(t), mocks.NewSharedPeer(t))

	rcv := newPanicOnFirstReceiver()
	require.NoError(t, d.SetReceiver(capID1, donID1, rcv))

	// First message triggers the panic; the goroutine must survive.
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	// Second message must still be delivered.
	recvCh <- encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload2))

	m := <-rcv.received
	require.Equal(t, payload2, string(m.Payload))

	d.RemoveReceiver(capID1, donID1)
	require.NoError(t, d.Close())
}
