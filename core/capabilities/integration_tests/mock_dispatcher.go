package integration_tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"google.golang.org/protobuf/proto"
)

// testAsyncMessageBroker backs the dispatchers created for each node in the test and effectively
// acts as the rageP2P network layer.
type testAsyncMessageBroker struct {
	services.StateMachine
	t *testing.T

	chanBufferSize int
	stopCh         services.StopChan
	wg             sync.WaitGroup

	peerIDToBrokerNode map[p2ptypes.PeerID]*brokerNode

	mux sync.Mutex
}

func newTestAsyncMessageBroker(t *testing.T, chanBufferSize int) *testAsyncMessageBroker {
	return &testAsyncMessageBroker{
		t:                  t,
		stopCh:             make(services.StopChan),
		chanBufferSize:     chanBufferSize,
		peerIDToBrokerNode: make(map[p2ptypes.PeerID]*brokerNode),
	}
}

func (a *testAsyncMessageBroker) Start(ctx context.Context) error {
	return a.StartOnce("testAsyncMessageBroker", func() error {
		return nil
	})
}

func (a *testAsyncMessageBroker) Close() error {
	return a.StopOnce("testAsyncMessageBroker", func() error {
		close(a.stopCh)
		a.wg.Wait()
		return nil
	})
}

// NewDispatcherForNode creates a new dispatcher for a node with the given peer ID.
func (a *testAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &brokerDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
		receivers:    map[key]remotetypes.Receiver{},
	}
}

func (a *testAsyncMessageBroker) HealthReport() map[string]error {
	return nil
}

func (a *testAsyncMessageBroker) Name() string {
	return "testAsyncMessageBroker"
}

func (a *testAsyncMessageBroker) registerReceiverNode(nodePeerID p2ptypes.PeerID, capabilityId string, capabilityDonID uint32, receiver remotetypes.Receiver) {
	a.mux.Lock()
	defer a.mux.Unlock()

	node, ok := a.peerIDToBrokerNode[nodePeerID]
	if !ok {
		node = a.newNode()
		a.peerIDToBrokerNode[nodePeerID] = node
	}

	node.registerReceiverCh <- &registerReceiverRequest{
		receiverKey: receiverKey{
			capabilityId: capabilityId,
			donId:        capabilityDonID,
		},
		receiver: receiver,
	}
}

func (a *testAsyncMessageBroker) Send(msg *remotetypes.MessageBody) {
	peerID := toPeerID(msg.Receiver)
	node, ok := a.peerIDToBrokerNode[peerID]
	if !ok {
		panic(fmt.Sprintf("node not found for peer ID %v", peerID))
	}

	node.receiveCh <- msg
}

type brokerNode struct {
	registerReceiverCh chan *registerReceiverRequest
	receiveCh          chan *remotetypes.MessageBody
}

type receiverKey struct {
	capabilityId string
	donId        uint32
}

type registerReceiverRequest struct {
	receiverKey
	receiver remotetypes.Receiver
}

func (a *testAsyncMessageBroker) newNode() *brokerNode {
	n := &brokerNode{
		receiveCh:          make(chan *remotetypes.MessageBody, a.chanBufferSize),
		registerReceiverCh: make(chan *registerReceiverRequest, a.chanBufferSize),
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		receivers := make(map[receiverKey]remotetypes.Receiver)
		for {
			select {
			case <-a.stopCh:
				return
			case msg := <-n.receiveCh:
				k := receiverKey{
					capabilityId: msg.CapabilityId,
					donId:        msg.CapabilityDonId,
				}

				r, ok := receivers[k]
				if !ok {
					panic(fmt.Sprintf("receiver not found for key %+v", k))
				}

				r.Receive(tests.Context(a.t), msg)
			case reg := <-n.registerReceiverCh:
				receivers[reg.receiverKey] = reg.receiver
			}
		}
	}()
	return n
}

func toPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type broker interface {
	Send(msg *remotetypes.MessageBody)
}

type brokerDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       broker

	receivers map[key]remotetypes.Receiver
	mu        sync.Mutex
}

type key struct {
	capId string
	donId uint32
}

func (t *brokerDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	clonedMsg := proto.Clone(msgBody).(*remotetypes.MessageBody)
	clonedMsg.Version = 1
	clonedMsg.Sender = t.callerPeerID[:]
	clonedMsg.Receiver = peerID[:]
	clonedMsg.Timestamp = time.Now().UnixMilli()
	t.broker.Send(clonedMsg)
	return nil
}

func (t *brokerDispatcher) SetReceiver(capabilityId string, donId uint32, receiver remotetypes.Receiver) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key{capabilityId, donId}
	_, ok := t.receivers[k]
	if ok {
		return fmt.Errorf("%w: receiver already exists for capability %s and don %d", remote.ErrReceiverExists, capabilityId, donId)
	}
	t.receivers[k] = receiver

	t.broker.(*testAsyncMessageBroker).registerReceiverNode(t.callerPeerID, capabilityId, donId, receiver)
	return nil
}
func (t *brokerDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

func (t *brokerDispatcher) Start(context.Context) error { return nil }

func (t *brokerDispatcher) Close() error {
	return nil
}

func (t *brokerDispatcher) Ready() error {
	return nil
}

func (t *brokerDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *brokerDispatcher) Name() string {
	return "mockDispatcher"
}
