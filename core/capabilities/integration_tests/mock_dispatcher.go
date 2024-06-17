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

	"github.com/golang/protobuf/proto"
)

type receiverKey struct {
	peerID       p2ptypes.PeerID
	capabilityId string
	donId        string
}

type testAsyncMessageBroker struct {
	services.StateMachine
	t *testing.T

	nodes map[receiverKey]remotetypes.Receiver

	sendCh chan *remotetypes.MessageBody

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func (a *testAsyncMessageBroker) HealthReport() map[string]error {
	return nil
}

func (a *testAsyncMessageBroker) Name() string {
	return "testAsyncMessageBroker"
}

func newTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *testAsyncMessageBroker {
	return &testAsyncMessageBroker{
		t:      t,
		nodes:  make(map[receiverKey]remotetypes.Receiver),
		stopCh: make(services.StopChan),
		sendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
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
					receiverId := toPeerID(msg.Receiver)

					key := receiverKey{
						peerID:       receiverId,
						capabilityId: msg.CapabilityId,
						donId:        msg.CapabilityDonId,
					}

					receiver, ok := a.nodes[key]
					if !ok {
						panic("server not found for peer id")
					}

					receiver.Receive(tests.Context(a.t), msg)
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

func (a *testAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
	}
}

func (a *testAsyncMessageBroker) registerReceiverNode(nodePeerID p2ptypes.PeerID, capabilityId string, capabilityDonID string, node remotetypes.Receiver) {
	key := receiverKey{
		peerID:       nodePeerID,
		capabilityId: capabilityId,
		donId:        capabilityDonID,
	}

	//	fmt.Printf("registering receiver node: %s %s %s\n", key.peerID, key.capabilityId, key.donId)
	//  here syncer is duplciate registering the same capability
	//	if _, ok := a.nodes[key]; ok {
	//		panic(fmt.Sprintf("capability already registered: %s %s %s", key.peerID, key.capabilityId, key.donId))
	//	}

	a.nodes[key] = node
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

func (t *nodeDispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	t.broker.(*testAsyncMessageBroker).registerReceiverNode(t.callerPeerID, capabilityId, donId, receiver)
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityId string, donId string) {}
