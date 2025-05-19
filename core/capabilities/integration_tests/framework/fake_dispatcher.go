package framework

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"google.golang.org/protobuf/proto"
)

// FakeRageP2PNetwork backs the dispatchers created for each node in the test and effectively
// acts as the rageP2P network layer.
type FakeRageP2PNetwork struct {
	services.StateMachine
	t          *testing.T
	readyError error

	chanBufferSize int
	stopCh         services.StopChan
	wg             sync.WaitGroup

	peerIDToBrokerNode map[p2ptypes.PeerID]*brokerNode

	capabilityRegistrations map[CapabilityRegistration]bool

	mux sync.Mutex
}

func NewFakeRageP2PNetwork(ctx context.Context, t *testing.T, chanBufferSize int) *FakeRageP2PNetwork {
	network := &FakeRageP2PNetwork{
		t:                       t,
		stopCh:                  make(services.StopChan),
		chanBufferSize:          chanBufferSize,
		peerIDToBrokerNode:      make(map[p2ptypes.PeerID]*brokerNode),
		capabilityRegistrations: make(map[CapabilityRegistration]bool),
	}

	go func() {
		<-ctx.Done()
		network.SetReadyError(errors.New("context done"))
	}()

	return network
}

func (a *FakeRageP2PNetwork) Start(ctx context.Context) error {
	return a.StartOnce("FakeRageP2PNetwork", func() error {
		return nil
	})
}

func (a *FakeRageP2PNetwork) Close() error {
	return a.StopOnce("FakeRageP2PNetwork", func() error {
		close(a.stopCh)
		a.wg.Wait()
		return nil
	})
}

func (a *FakeRageP2PNetwork) Ready() error {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.readyError
}

func (a *FakeRageP2PNetwork) SetReadyError(err error) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.readyError = err
}

// NewDispatcherForNode creates a new dispatcher for a node with the given peer ID.
func (a *FakeRageP2PNetwork) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &brokerDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
		receivers:    map[key]remotetypes.Receiver{},
	}
}

func (a *FakeRageP2PNetwork) HealthReport() map[string]error {
	return nil
}

func (a *FakeRageP2PNetwork) Name() string {
	return "FakeRageP2PNetwork"
}

type CapabilityRegistration struct {
	nodePeerID      string
	capabilityID    string
	capabilityDonID uint32
}

func (a *FakeRageP2PNetwork) GetCapabilityRegistrations() map[CapabilityRegistration]bool {
	a.mux.Lock()
	defer a.mux.Unlock()

	copiedRegistrations := make(map[CapabilityRegistration]bool)
	for k, v := range a.capabilityRegistrations {
		copiedRegistrations[k] = v
	}
	return copiedRegistrations
}

func (a *FakeRageP2PNetwork) registerReceiverNode(nodePeerID p2ptypes.PeerID, capabilityID string, capabilityDonID uint32, receiver remotetypes.Receiver) {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.capabilityRegistrations[CapabilityRegistration{
		nodePeerID:      hex.EncodeToString(nodePeerID[:]),
		capabilityID:    capabilityID,
		capabilityDonID: capabilityDonID,
	}] = true

	node, ok := a.peerIDToBrokerNode[nodePeerID]
	if !ok {
		node = a.newNode()
		a.peerIDToBrokerNode[nodePeerID] = node
	}

	node.registerReceiverCh <- &registerReceiverRequest{
		receiverKey: receiverKey{
			capabilityId: capabilityID,
			donId:        capabilityDonID,
		},
		receiver: receiver,
	}
}

func (a *FakeRageP2PNetwork) Send(msg *remotetypes.MessageBody) {
	peerID := toPeerID(msg.Receiver)

	node, ok := a.getNodeForPeerID(peerID)
	if !ok {
		panic(fmt.Sprintf("node not found for peer ID %v", peerID))
	}

	node.receiveCh <- msg
}

func (a *FakeRageP2PNetwork) getNodeForPeerID(peerID types.PeerID) (*brokerNode, bool) {
	a.mux.Lock()
	defer a.mux.Unlock()
	node, ok := a.peerIDToBrokerNode[peerID]
	return node, ok
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

func (a *FakeRageP2PNetwork) newNode() *brokerNode {
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

				r.Receive(a.t.Context(), msg)
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
	Ready() error
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

	t.broker.(*FakeRageP2PNetwork).registerReceiverNode(t.callerPeerID, capabilityId, donId, receiver)
	return nil
}
func (t *brokerDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

func (t *brokerDispatcher) Start(context.Context) error { return nil }

func (t *brokerDispatcher) Close() error {
	return nil
}

func (t *brokerDispatcher) Ready() error {
	return t.broker.Ready()
}

func (t *brokerDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *brokerDispatcher) Name() string {
	return "fakeDispatcher"
}
