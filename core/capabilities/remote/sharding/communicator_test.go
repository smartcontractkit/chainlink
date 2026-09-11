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
	return commoncap.DON{ID: id, F: f, Members: members}
}

func makeStatusMsg(workflowID, triggerEventID, executionID string, status ringpb.ExecutionStatus, primaryDonID uint32) *ringpb.ExecutionStatusUpdate {
	return &ringpb.ExecutionStatusUpdate{
		WorkflowId:     workflowID,
		ExecutionId:    executionID,
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

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	comm.Send(ctx, msg, secondaryDON)

	require.Len(t, sentBodies, len(secondaryDON.Members), "should send to all secondary members")
	for _, body := range sentBodies {
		assert.Equal(t, ShardExecutionStatusUpdateCapabilityID, body.CapabilityId)
		assert.Equal(t, remotetypes.MethodExecutionStatusUpdate, body.CapabilityMethod)
		assert.Equal(t, secondaryDON.ID, body.CapabilityDonId)
		assert.Equal(t, primaryDON.ID, body.CallerDonId)
	}
}

func TestShardFailoverCommunicator_MultipleWorkflowsNoConflict(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	wf1Received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	wf2Received := make(chan *ringpb.ExecutionStatusUpdate, 1)

	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		wf1Received <- msg
	})
	comm.RegisterHandler("wf-2", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		wf2Received <- msg
	})

	sendStatusUpdate := func(workflowID, triggerEventID, executionID string) {
		msg := makeStatusMsg(workflowID, triggerEventID, executionID, ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
		payload, _ := proto.Marshal(msg)
		for _, peer := range primaryDON.Members[:2] {
			require.NoError(t, disp.Send(peer, &remotetypes.MessageBody{
				CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
				Method:           remotetypes.MethodExecutionStatusUpdate,
				CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
				CapabilityDonId:  localDON.ID,
				CallerDonId:      primaryDON.ID,
				Payload:          payload,
				Sender:           peer[:],
			}))
		}
	}

	sendStatusUpdate("wf-1", "evt-1", "exec-1")
	select {
	case got := <-wf1Received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, "exec-1", got.ExecutionId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-1 handler")
	}

	sendStatusUpdate("wf-2", "evt-2", "exec-2")
	select {
	case got := <-wf2Received:
		assert.Equal(t, "wf-2", got.WorkflowId)
		assert.Equal(t, "exec-2", got.ExecutionId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-2 handler")
	}
}

func TestShardFailoverCommunicator_EndToEndSendReceive(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()

	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))
	secondaryDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))

	primaryComm := NewShardFailoverCommunicator(disp, primaryDON.ID, logger.Test(t))
	require.NoError(t, primaryComm.Start(ctx))
	t.Cleanup(func() { _ = primaryComm.Close() })

	secondaryComm := NewShardFailoverCommunicator(disp, secondaryDON.ID, logger.Test(t))
	require.NoError(t, secondaryComm.Start(ctx))
	t.Cleanup(func() { _ = secondaryComm.Close() })

	received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	secondaryComm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		received <- msg
	})

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, primaryDON.ID)
	primaryComm.Send(ctx, msg, secondaryDON)

	sentBodies := disp.getSentBodies()
	require.Len(t, sentBodies, len(secondaryDON.Members))

	for i, body := range sentBodies {
		body.Sender = primaryDON.Members[i][:]
		secondaryComm.Receive(ctx, body)
	}

	select {
	case got := <-received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, "exec-1", got.ExecutionId)
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
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDON.ID,
		Payload:          payload,
		Sender:           primaryDON.Members[0][:],
	}
	comm.Receive(ctx, body)

	time.Sleep(100 * time.Millisecond)
	assert.False(t, handlerCalled, "handler should not be called before quorum")

	body.Sender = primaryDON.Members[1][:]
	comm.Receive(ctx, body)

	assert.True(t, handlerCalled, "handler should be called after quorum")
}

func TestShardFailoverCommunicator_UnregisteredWorkflowDropped(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	primaryDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	msg := makeStatusMsg("wf-unknown", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDON.ID,
		Payload:          payload,
		Sender:           primaryDON.Members[0][:],
	}
	assert.NotPanics(t, func() {
		comm.Receive(ctx, body)
	})
}

func TestShardFailoverCommunicator_UnregisterHandler(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 0, makeTestPeerID(20))
	primaryDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})
	comm.UnregisterHandler("wf-1")

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDON.ID,
		Payload:          payload,
		Sender:           primaryDON.Members[0][:],
	}
	comm.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called after unregister")
}

func TestShardFailoverCommunicator_IgnoresUnknownPeer(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 0, makeTestPeerID(20))
	primaryDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})

	unknownPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	unknownRagePeerID, err := ragetypes.PeerIDFromPublicKey(unknownPubKey)
	require.NoError(t, err)

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	payload, _ := proto.Marshal(msg)

	body := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDON.ID,
		Payload:          payload,
		Sender:           unknownRagePeerID[:],
	}
	comm.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called for unknown peer")
}

func TestShardFailoverCommunicator_DifferentWorkflowsDifferentPrimaryDONs(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(3, 1, makeTestPeerID(30), makeTestPeerID(31), makeTestPeerID(32))
	primaryDONA := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12))
	primaryDONB := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	wfAReceived := make(chan *ringpb.ExecutionStatusUpdate, 1)
	wfBReceived := make(chan *ringpb.ExecutionStatusUpdate, 1)

	// Two workflows with different primary DONs — this is the key
	// scenario the reviewer identified: primary/secondary is per-workflow.
	comm.RegisterHandler("wf-A", primaryDONA, func(msg *ringpb.ExecutionStatusUpdate) {
		wfAReceived <- msg
	})
	comm.RegisterHandler("wf-B", primaryDONB, func(msg *ringpb.ExecutionStatusUpdate) {
		wfBReceived <- msg
	})

	// Send from primaryDONA members for wf-A
	msgA := makeStatusMsg("wf-A", "evt-A", "exec-A", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDONA.ID)
	payloadA, _ := proto.Marshal(msgA)
	for _, peer := range primaryDONA.Members[:2] {
		require.NoError(t, disp.Send(peer, &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			CapabilityDonId:  localDON.ID,
			CallerDonId:      primaryDONA.ID,
			Payload:          payloadA,
			Sender:           peer[:],
		}))
	}
	select {
	case got := <-wfAReceived:
		assert.Equal(t, "wf-A", got.WorkflowId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-A handler")
	}

	// Send from primaryDONB members for wf-B
	msgB := makeStatusMsg("wf-B", "evt-B", "exec-B", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDONB.ID)
	payloadB, _ := proto.Marshal(msgB)
	for _, peer := range primaryDONB.Members[:2] {
		require.NoError(t, disp.Send(peer, &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			CapabilityDonId:  localDON.ID,
			CallerDonId:      primaryDONB.ID,
			Payload:          payloadB,
			Sender:           peer[:],
		}))
	}
	select {
	case got := <-wfBReceived:
		assert.Equal(t, "wf-B", got.WorkflowId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-B handler")
	}

	// Verify messages from primaryDONA are rejected for wf-B (different primary)
	bodyFromA := &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDONA.ID,
		Payload:          payloadB,
		Sender:           primaryDONA.Members[0][:],
	}
	comm.Receive(ctx, bodyFromA)
	select {
	case <-wfBReceived:
		t.Fatal("wf-B should not accept messages from primaryDONA")
	default:
	}
}
