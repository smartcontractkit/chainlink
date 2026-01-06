package e2etesting

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type TestAsyncMessageBroker struct {
	services.Service
	eng *services.Engine
	t   *testing.T

	nodes map[p2ptypes.PeerID]remotetypes.Receiver

	sendCh chan *remotetypes.MessageBody

	messageFilter atomic.Value // stores func(msg *remotetypes.MessageBody) bool
}

func NewTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *TestAsyncMessageBroker {
	b := &TestAsyncMessageBroker{
		t:      t,
		nodes:  make(map[p2ptypes.PeerID]remotetypes.Receiver),
		sendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
	}
	b.Service, b.eng = services.Config{
		Name:  "TestAsyncMessageBroker",
		Start: b.start,
	}.NewServiceEngine(logger.Test(t))
	return b
}

// SetMessageFilter sets the message filter function in a thread-safe way.
func (a *TestAsyncMessageBroker) SetMessageFilter(filter func(msg *remotetypes.MessageBody) bool) {
	a.messageFilter.Store(filter)
}

func (a *TestAsyncMessageBroker) start(ctx context.Context) error {
	fmt.Print("TestAsyncMessageBroker started\n")
	a.eng.Go(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.sendCh:
				receiverID := ToPeerID(msg.Receiver)

				receiver, ok := a.nodes[receiverID]
				if !ok {
					panic("server not found for peer id")
				}

				receiver.Receive(a.t.Context(), msg)
			}
		}
	})
	return nil
}

func (a *TestAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
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
	filterVal := a.messageFilter.Load()
	if filterVal != nil {
		if filter := filterVal.(func(msg *remotetypes.MessageBody) bool); !filter(msg) {
			// Drop the message
			return
		}
	}

	// Clone the message to simulate real behaviour when sent across the network
	cloned := proto.Clone(msg).(*remotetypes.MessageBody)
	a.sendCh <- cloned
}

func ToPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type broker interface {
	Send(msg *remotetypes.MessageBody)
}

type nodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       broker
}

func (t *nodeDispatcher) Name() string {
	return "nodeDispatcher"
}

func (t *nodeDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *nodeDispatcher) Close() error {
	return nil
}

func (t *nodeDispatcher) Ready() error {
	return nil
}

func (t *nodeDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *nodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.Send(msgBody)
	return nil
}

func (t *nodeDispatcher) SetReceiver(capabilityID string, donID uint32, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (t *nodeDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, methodName string, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, methodName string) {
}

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

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID)
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
