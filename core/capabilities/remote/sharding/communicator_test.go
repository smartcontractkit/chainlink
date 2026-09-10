package sharding

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	dispatchermocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// --- in-memory dispatcher for integration tests ---

// inMemDispatcher simulates the real dispatcher's routing logic for testing.
// It routes messages based on {CapabilityId, CapabilityDonId, CapabilityMethod}
// just like the real dispatcher, but without P2P/signing.
type inMemDispatcher struct {
	mu         sync.Mutex
	receivers  map[string]remotetypes.Receiver
	sentBodies []*remotetypes.MessageBody
}

func newInMemDispatcher() *inMemDispatcher {
	return &inMemDispatcher{receivers: make(map[string]remotetypes.Receiver)}
}

func (d *inMemDispatcher) receiverKey(capID string, donID uint32, method string) string {
	return fmt.Sprintf("%s:%d:%s", capID, donID, method)
}

func (d *inMemDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, method string, receiver remotetypes.Receiver) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	k := d.receiverKey(capabilityID, donID, method)
	if _, ok := d.receivers[k]; ok {
		return fmt.Errorf("receiver already exists for %s", k)
	}
	d.receivers[k] = receiver
	return nil
}

func (d *inMemDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, method string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.receivers, d.receiverKey(capabilityID, donID, method))
}

func (d *inMemDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	d.mu.Lock()
	d.sentBodies = append(d.sentBodies, msgBody)
	d.mu.Unlock()
	// Simulate the message arriving at the destination node's dispatcher.
	k := d.receiverKey(msgBody.CapabilityId, msgBody.CapabilityDonId, msgBody.CapabilityMethod)
	d.mu.Lock()
	rcv, ok := d.receivers[k]
	d.mu.Unlock()
	if ok {
		rcv.Receive(context.Background(), msgBody)
	}
	return nil
}

func (d *inMemDispatcher) getSentBodies() []*remotetypes.MessageBody {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]*remotetypes.MessageBody{}, d.sentBodies...)
}

// Stub out the rest of the Dispatcher interface.
func (d *inMemDispatcher) SetReceiver(string, uint32, remotetypes.Receiver) error { return nil }
func (d *inMemDispatcher) RemoveReceiver(string, uint32)                          {}
func (d *inMemDispatcher) Start(context.Context) error                            { return nil }
func (d *inMemDispatcher) Close() error                                           { return nil }
func (d *inMemDispatcher) Ready() error                                           { return nil }
func (d *inMemDispatcher) HealthReport() map[string]error                         { return nil }
func (d *inMemDispatcher) Name() string                                           { return "inMemDispatcher" }

// --- helpers ---

func makeTestPeerID(b byte) ragetypes.PeerID {
	var id ragetypes.PeerID
	id[0] = b
	return id
}

func makeTestDON(id uint32, f uint8, members ...ragetypes.PeerID) commoncap.DON {
	return commoncap.DON{
		ID:      id,
		F:       f,
		Members: members,
	}
}

func makeStatusMsg(workflowID, triggerEventID string, status ringpb.ExecutionStatus, primaryDonID uint32) *ringpb.ExecutionStatusUpdate {
	return &ringpb.ExecutionStatusUpdate{
		WorkflowId:     workflowID,
		TriggerEventId: triggerEventID,
		TriggerIndex:   0,
		Status:         status,
		PrimaryDonId:   primaryDonID,
	}
}

// --- tests ---

func TestShardFailoverCommunicator_SendSetsCapabilityID(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	mockDisp := dispatchermocks.NewDispatcher(t)
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))
	secondaryDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))

	var sentBodies []*remotetypes.MessageBody
	mockDisp.On("Send", mock.Anything, mock.MatchedBy(func(body *remotetypes.MessageBody) bool {
		sentBodies = append(sentBodies, body)
		return true
	})).Return(nil)

	comm := NewShardFailoverCommunicator(mockDisp, primaryDON.ID, logger.Test(t))
	comm.SetShardDons(primaryDON, secondaryDON)

	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	comm.Send(ctx, msg)

	require.Len(t, sentBodies, len(secondaryDON.Members), "should send to all secondary members")
	for _, body := range sentBodies {
		assert.Equal(t, ShardExecutionStatusUpdateCapabilityID, body.CapabilityId,
			"CapabilityId must be set for dispatcher routing")
		assert.Equal(t, remotetypes.MethodExecutionStatusUpdate, body.CapabilityMethod,
			"CapabilityMethod must be set for dispatcher routing")
		assert.Equal(t, secondaryDON.ID, body.CapabilityDonId,
			"CapabilityDonId must match the peer DON ID")
		assert.Equal(t, primaryDON.ID, body.CallerDonId,
			"CallerDonId must be the local DON ID")
	}
}

func TestShardFailoverCommunicator_MultipleWorkflowsNoConflict(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	peerDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	comm.SetShardDons(peerDON, commoncap.DON{})
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	wf1Received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	wf2Received := make(chan *ringpb.ExecutionStatusUpdate, 1)

	comm.RegisterHandler("wf-1", func(msg *ringpb.ExecutionStatusUpdate) {
		wf1Received <- msg
	})
	comm.RegisterHandler("wf-2", func(msg *ringpb.ExecutionStatusUpdate) {
		wf2Received <- msg
	})

	// Simulate messages from the peer shard (primary) reaching the dispatcher.
	// Each message from a different peer, both routed to the same communicator.
	msg1 := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)
	msg2 := makeStatusMsg("wf-2", "evt-2", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)

	payload1, _ := proto.Marshal(msg1)
	payload2, _ := proto.Marshal(msg2)

	// Send from peer members to simulate quorum (F+1 = 2 messages needed)
	for _, peer := range peerDON.Members[:2] {
		require.NoError(t, disp.Send(peer, &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			CapabilityDonId:  localDON.ID,
			CallerDonId:      peerDON.ID,
			Payload:          payload1,
			Sender:           peer[:],
		}))
	}

	select {
	case got := <-wf1Received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, "evt-1", got.TriggerEventId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-1 handler")
	}

	for _, peer := range peerDON.Members[:2] {
		require.NoError(t, disp.Send(peer, &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			CapabilityDonId:  localDON.ID,
			CallerDonId:      peerDON.ID,
			Payload:          payload2,
			Sender:           peer[:],
		}))
	}

	select {
	case got := <-wf2Received:
		assert.Equal(t, "wf-2", got.WorkflowId)
		assert.Equal(t, "evt-2", got.TriggerEventId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-2 handler")
	}
}

func TestShardFailoverCommunicator_EndToEndSendReceive(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()

	// Primary node: DON 1, Secondary node: DON 2
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))
	secondaryDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))

	// Primary's communicator (for sending)
	primaryComm := NewShardFailoverCommunicator(disp, primaryDON.ID, logger.Test(t))
	primaryComm.SetShardDons(primaryDON, secondaryDON)
	require.NoError(t, primaryComm.Start(ctx))
	t.Cleanup(func() { _ = primaryComm.Close() })

	// Secondary's communicator (for receiving)
	secondaryComm := NewShardFailoverCommunicator(disp, secondaryDON.ID, logger.Test(t))
	secondaryComm.SetShardDons(primaryDON, secondaryDON)
	require.NoError(t, secondaryComm.Start(ctx))
	t.Cleanup(func() { _ = secondaryComm.Close() })

	// Register handler on secondary for workflow "wf-1"
	received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	secondaryComm.RegisterHandler("wf-1", func(msg *ringpb.ExecutionStatusUpdate) {
		received <- msg
	})

	// Primary sends a status update
	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, primaryDON.ID)
	primaryComm.Send(ctx, msg)

	// The inMemDispatcher delivers to each secondary member, but only
	// the secondaryComm (registered with secondaryDON.ID) receives it.
	// Quorum = F+1 = 2. We need 2 deliveries from different primary peers.
	// The primaryComm.Send sends to all secondary members, but the inMemDispatcher
	// routes based on {CapabilityId, CapabilityDonId, CapabilityMethod}.
	// The CapabilityDonId is secondaryDON.ID, so it routes to secondaryComm.
	// However, the inMemDispatcher.Send is called once per secondary member,
	// each with the same Sender (primaryComm doesn't set Sender; the real
	// dispatcher sets Sender on Send). Let's simulate this properly:
	sentBodies := disp.getSentBodies()
	require.Len(t, sentBodies, len(secondaryDON.Members))

	// Feed the sent messages into the secondaryComm as if they came from
	// different primary peers (for quorum).
	for i, body := range sentBodies {
		body.Sender = primaryDON.Members[i][:]
		secondaryComm.Receive(ctx, body)
	}

	select {
	case got := <-received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, "evt-1", got.TriggerEventId)
		assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, got.Status)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for handler invocation")
	}
}

func TestShardFailoverCommunicator_QuorumRequired(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	peerDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	comm.SetShardDons(peerDON, commoncap.DON{})
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})

	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)
	payload, _ := proto.Marshal(msg)

	// Send only 1 message (quorum = F+1 = 2, so 1 is not enough)
	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      peerDON.ID,
		Payload:          payload,
		Sender:           peerDON.Members[0][:],
	}
	comm.Receive(ctx, body)

	// Give some time to ensure handler is not called
	time.Sleep(100 * time.Millisecond)
	assert.False(t, handlerCalled, "handler should not be called before quorum")

	// Send second message from a different peer (reaches quorum)
	body.Sender = peerDON.Members[1][:]
	comm.Receive(ctx, body)

	// Handler should now be called (synchronously in Receive)
	assert.True(t, handlerCalled, "handler should be called after quorum")
}

func TestShardFailoverCommunicator_UnregisteredWorkflowDropped(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	peerDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	comm.SetShardDons(peerDON, commoncap.DON{})
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	// No handler registered for "wf-unknown"
	msg := makeStatusMsg("wf-unknown", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      peerDON.ID,
		Payload:          payload,
		Sender:           peerDON.Members[0][:],
	}
	// Should not panic, just log and drop
	assert.NotPanics(t, func() {
		comm.Receive(ctx, body)
	})
}

func TestShardFailoverCommunicator_UnregisterHandler(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 0, makeTestPeerID(20))
	peerDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	comm.SetShardDons(peerDON, commoncap.DON{})
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})
	comm.UnregisterHandler("wf-1")

	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      peerDON.ID,
		Payload:          payload,
		Sender:           peerDON.Members[0][:],
	}
	comm.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called after unregister")
}

func TestShardFailoverCommunicator_IgnoresUnknownPeer(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 0, makeTestPeerID(20))
	peerDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	comm.SetShardDons(peerDON, commoncap.DON{})
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})

	unknownPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	unknownRagePeerID, err := ragetypes.PeerIDFromPublicKey(unknownPubKey)
	require.NoError(t, err)

	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, peerDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      peerDON.ID,
		Payload:          payload,
		Sender:           unknownRagePeerID[:],
	}
	comm.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called for unknown peer")
}

func TestShardFailoverCommunicator_SetShardDons_UpdatesDynamically(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	mockDisp := dispatchermocks.NewDispatcher(t)
	mockDisp.On("SetReceiverForMethod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDisp.On("RemoveReceiverForMethod", mock.Anything, mock.Anything, mock.Anything).Return()
	mockDisp.On("Send", mock.Anything, mock.Anything).Return(nil)

	localDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))
	secondaryDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	anotherSecondaryDON := makeTestDON(3, 1, makeTestPeerID(30), makeTestPeerID(31), makeTestPeerID(32))

	comm := NewShardFailoverCommunicator(mockDisp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	// Initial shard DONs
	comm.SetShardDons(localDON, secondaryDON)

	msg := makeStatusMsg("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, localDON.ID)
	comm.Send(ctx, msg)
	mockDisp.AssertNumberOfCalls(t, "Send", len(secondaryDON.Members))

	// Dynamically update to a different secondary DON (simulates shard reassignment)
	comm.SetShardDons(localDON, anotherSecondaryDON)
	comm.Send(ctx, msg)
	mockDisp.AssertNumberOfCalls(t, "Send", len(secondaryDON.Members)+len(anotherSecondaryDON.Members))
}
