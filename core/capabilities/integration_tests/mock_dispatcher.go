package integration_tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"google.golang.org/protobuf/proto"
)

type receiverKey struct {
	capabilityId string
	donId        uint32
}

// testAsyncMessageBroker backs the dispatchers created for each node in the test and effectively
// acts as the rageP2P network layer.
type testAsyncMessageBroker struct {
	services.StateMachine
	t *testing.T

	nodes map[p2ptypes.PeerID]*dispatcherNode

	sendCh chan *remotetypes.MessageBody

	chanBufferSize int

	stopCh services.StopChan
	wg     sync.WaitGroup
}

// NewDispatcherForNode creates a new dispatcher for a node with the given peer ID.
func (a *testAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
	}
}

func (a *testAsyncMessageBroker) HealthReport() map[string]error {
	return nil
}

func (a *testAsyncMessageBroker) Name() string {
	return "testAsyncMessageBroker"
}

func newTestAsyncMessageBroker(t *testing.T, chanBufferSize int) *testAsyncMessageBroker {
	return &testAsyncMessageBroker{
		t:              t,
		nodes:          make(map[p2ptypes.PeerID]*dispatcherNode),
		stopCh:         make(services.StopChan),
		sendCh:         make(chan *remotetypes.MessageBody, chanBufferSize),
		chanBufferSize: chanBufferSize,
	}
}

func (a *testAsyncMessageBroker) Start(ctx context.Context) error {
	return a.StartOnce("testAsyncMessageBroker", func() error {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for {
				select {
				case <-a.stopCh:
					return
				case msg := <-a.sendCh:
					peerID := toPeerID(msg.Receiver)
					node, ok := a.nodes[peerID]
					if !ok {
						panic("node not found for peer id")
					}

					node.receiveCh <- msg
				}
			}
		}()
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

type dispatcherNode struct {
	receivers map[receiverKey]remotetypes.Receiver
	receiveCh chan *remotetypes.MessageBody
}

func (a *testAsyncMessageBroker) registerReceiverNode(nodePeerID p2ptypes.PeerID, capabilityId string, capabilityDonID uint32, receiver remotetypes.Receiver) {
	key := receiverKey{
		capabilityId: capabilityId,
		donId:        capabilityDonID,
	}

	node, nodeExists := a.nodes[nodePeerID]
	if !nodeExists {
		node = &dispatcherNode{
			receivers: make(map[receiverKey]remotetypes.Receiver),
			receiveCh: make(chan *remotetypes.MessageBody, a.chanBufferSize),
		}

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for {
				select {
				case <-a.stopCh:
					return
				case msg := <-node.receiveCh:
					k := receiverKey{
						capabilityId: msg.CapabilityId,
						donId:        msg.CapabilityDonId,
					}

					r, ok := node.receivers[k]
					if !ok {
						panic("receiver not found for key")
					}

					r.Receive(tests.Context(a.t), msg)
				}
			}
		}()

		a.nodes[nodePeerID] = node
	}

	node.receivers[key] = receiver
}

func (a *testAsyncMessageBroker) Send(msg *remotetypes.MessageBody) {
	a.sendCh <- msg
}

func toPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type broker interface {
	Send(msg *remotetypes.MessageBody)
}

type nodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       broker
}

func (t *nodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	clonedMsg := proto.Clone(msgBody).(*remotetypes.MessageBody)
	clonedMsg.Version = 1
	clonedMsg.Sender = t.callerPeerID[:]
	clonedMsg.Receiver = peerID[:]
	clonedMsg.Timestamp = time.Now().UnixMilli()
	t.broker.Send(clonedMsg)
	return nil
}

func (t *nodeDispatcher) SetReceiver(capabilityId string, donId uint32, receiver remotetypes.Receiver) error {
	t.broker.(*testAsyncMessageBroker).registerReceiverNode(t.callerPeerID, capabilityId, donId, receiver)
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}
