package testutils

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/test-go/testify/require"
)

const (
	WorkflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	WorkflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
)

type TestAsyncMessageBroker struct {
	services.Service
	eng *services.Engine
	t   *testing.T

	nodes map[p2ptypes.PeerID]remotetypes.Receiver

	SendCh chan *remotetypes.MessageBody
}

func NewTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *TestAsyncMessageBroker {
	b := &TestAsyncMessageBroker{
		t:      t,
		nodes:  make(map[p2ptypes.PeerID]remotetypes.Receiver),
		SendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
	}
	b.Service, b.eng = services.Config{
		Name:  "testAsyncMessageBroker",
		Start: b.start,
	}.NewServiceEngine(logger.TestLogger(t))
	return b
}

func (a *TestAsyncMessageBroker) start(ctx context.Context) error {
	a.eng.Go(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.SendCh:
				receiverId := ToPeerID(msg.Receiver)

				receiver, ok := a.nodes[receiverId]
				if !ok {
					panic("server not found for peer id")
				}

				receiver.Receive(tests.Context(a.t), msg)
			}
		}
	})
	return nil
}

func (a *TestAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &NodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
	}
}

func (a *TestAsyncMessageBroker) RegisterReceiverNode(nodePeerID p2ptypes.PeerID, node remotetypes.Receiver) {
	if _, ok := a.nodes[nodePeerID]; ok {
		panic("node already registered")
	}

	a.nodes[nodePeerID] = node
}

func (a *TestAsyncMessageBroker) Send(msg *remotetypes.MessageBody) {
	a.SendCh <- msg
}

func ToPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type Broker interface {
	Send(msg *remotetypes.MessageBody)
}

type NodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       Broker
}

func (t *NodeDispatcher) Name() string {
	return "nodeDispatcher"
}

func (t *NodeDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *NodeDispatcher) Close() error {
	return nil
}

func (t *NodeDispatcher) Ready() error {
	return nil
}

func (t *NodeDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *NodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.Send(msgBody)
	return nil
}

func (t *NodeDispatcher) SetReceiver(capabilityId string, donId uint32, receiver remotetypes.Receiver) error {
	return nil
}
func (t *NodeDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

func NewP2PPeerID(t *testing.T) p2ptypes.PeerID {
	id := p2ptypes.PeerID{}
	require.NoError(t, id.UnmarshalText([]byte(NewPeerID())))
	return id
}

func NewPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(Libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID[:])
}

func Libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
