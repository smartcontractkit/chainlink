package sharding

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
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

func TestShardFailoverCommunicator_MultipleWorkflowsNoConflict(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))

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
		// 2F+1 = 3 matching messages needed for quorum
		for _, peer := range primaryDON.Members[:3] {
			require.NoError(t, disp.Send(peer, statusBody(t, localDON, primaryDON, peer, msg)))
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

	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))
	secondaryDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))

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
	for _, body := range sentBodies {
		assert.Equal(t, primaryDON.ID, body.CallerDonId)
	}

	// Feed 3 messages from different primary peers (2F+1 = 3 quorum).
	for i := range 3 {
		sentBodies[i].Sender = primaryDON.Members[i][:]
		secondaryComm.Receive(ctx, sentBodies[i])
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
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	handlerCalled := false
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	})

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	body := statusBody(t, localDON, primaryDON, primaryDON.Members[0], msg)
	comm.Receive(ctx, body)

	assert.False(t, handlerCalled, "handler should not be called before quorum")

	// Send second message — 2 votes, still not enough (need 2F+1 = 3).
	body.Sender = primaryDON.Members[1][:]
	comm.Receive(ctx, body)

	assert.False(t, handlerCalled, "handler should not be called with only 2 matching")

	// Third message reaches 2F+1 = 3 quorum.
	body.Sender = primaryDON.Members[2][:]
	comm.Receive(ctx, body)

	assert.True(t, handlerCalled, "handler should be called after 2F+1 quorum")
}

func statusBody(t *testing.T, localDON, primaryDON commoncap.DON, sender ragetypes.PeerID, msg *ringpb.ExecutionStatusUpdate) *remotetypes.MessageBody {
	t.Helper()
	payload, err := proto.Marshal(msg)
	require.NoError(t, err)
	return &remotetypes.MessageBody{
		CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
		Method:           remotetypes.MethodExecutionStatusUpdate,
		CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
		CapabilityDonId:  localDON.ID,
		CallerDonId:      primaryDON.ID,
		Payload:          payload,
		Sender:           sender[:],
	}
}

func TestShardFailoverCommunicator_NoRedeliveryAfterQuorum(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))

	comm := NewShardFailoverCommunicator(newInMemDispatcher(), localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	var calls atomic.Int32
	comm.RegisterHandler("wf-1", primaryDON, func(*ringpb.ExecutionStatusUpdate) { calls.Add(1) })

	msg := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	for _, sender := range primaryDON.Members[:3] {
		comm.Receive(ctx, statusBody(t, localDON, primaryDON, sender, msg))
	}
	require.Equal(t, int32(1), calls.Load())

	// replaying the same report from every primary peer is enough votes for
	// another quorum, but it must not redeliver
	for _, sender := range primaryDON.Members {
		comm.Receive(ctx, statusBody(t, localDON, primaryDON, sender, msg))
	}
	require.Equal(t, int32(1), calls.Load())
}

func TestShardFailoverCommunicator_PrunesExpiredEntries(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))

	comm := NewShardFailoverCommunicator(newInMemDispatcher(), localDON.ID, logger.Test(t))
	comm.expiryDuration = 50 * time.Millisecond
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	var calls atomic.Int32
	comm.RegisterHandler("wf-1", primaryDON, func(*ringpb.ExecutionStatusUpdate) { calls.Add(1) })
	sizes := func() (pending, delivered int) {
		comm.quorumMu.Lock()
		defer comm.quorumMu.Unlock()
		return len(comm.pending), len(comm.delivered)
	}

	// a report that never reaches quorum is dropped from pending
	partial := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	comm.Receive(ctx, statusBody(t, localDON, primaryDON, primaryDON.Members[0], partial))
	pending, _ := sizes()
	require.Equal(t, 1, pending)
	require.Eventually(t, func() bool { p, _ := sizes(); return p == 0 }, 2*time.Second, 10*time.Millisecond)

	// a delivered update is forgotten, so the same quorum can fire again
	done := makeStatusMsg("wf-1", "evt-2", "exec-2", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	deliver := func() {
		for _, sender := range primaryDON.Members[:3] {
			comm.Receive(ctx, statusBody(t, localDON, primaryDON, sender, done))
		}
	}
	deliver()
	require.Equal(t, int32(1), calls.Load())
	require.Eventually(t, func() bool { _, d := sizes(); return d == 0 }, 2*time.Second, 10*time.Millisecond)
	deliver()
	require.Equal(t, int32(2), calls.Load())
}

func TestShardFailoverCommunicator_UnregisteredWorkflowDropped(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22))
	primaryDON := makeTestDON(1, 0, makeTestPeerID(10))

	comm := NewShardFailoverCommunicator(newInMemDispatcher(), localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	msg := makeStatusMsg("wf-unknown", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	comm.Receive(ctx, statusBody(t, localDON, primaryDON, primaryDON.Members[0], msg))

	comm.quorumMu.Lock()
	defer comm.quorumMu.Unlock()
	assert.Empty(t, comm.pending, "report for an unregistered workflow must not be tracked")
	assert.Empty(t, comm.delivered)
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
	body := statusBody(t, localDON, primaryDON, primaryDON.Members[0], msg)
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
	body := statusBody(t, localDON, primaryDON, unknownRagePeerID, msg)
	comm.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called for unknown peer")
}

func TestShardFailoverCommunicator_DifferentWorkflowsDifferentPrimaryDONs(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(3, 1, makeTestPeerID(30), makeTestPeerID(31), makeTestPeerID(32), makeTestPeerID(33))
	primaryDONA := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))
	primaryDONB := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))

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
	for _, peer := range primaryDONA.Members[:3] {
		require.NoError(t, disp.Send(peer, statusBody(t, localDON, primaryDONA, peer, msgA)))
	}
	select {
	case got := <-wfAReceived:
		assert.Equal(t, "wf-A", got.WorkflowId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-A handler")
	}

	// Send from primaryDONB members for wf-B
	msgB := makeStatusMsg("wf-B", "evt-B", "exec-B", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDONB.ID)
	for _, peer := range primaryDONB.Members[:3] {
		require.NoError(t, disp.Send(peer, statusBody(t, localDON, primaryDONB, peer, msgB)))
	}
	select {
	case got := <-wfBReceived:
		assert.Equal(t, "wf-B", got.WorkflowId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for wf-B handler")
	}

	// Verify messages from primaryDONA are rejected for wf-B (different primary)
	bodyFromA := statusBody(t, localDON, primaryDONA, primaryDONA.Members[0], msgB)
	comm.Receive(ctx, bodyFromA)
	select {
	case <-wfBReceived:
		t.Fatal("wf-B should not accept messages from primaryDONA")
	default:
	}
}

func TestShardFailoverCommunicator_MixedStatusesCountedSeparately(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	// N=5, F=1, quorum = 2F+1 = 3. Different statuses have different payloads, so votes are counted separately.
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13), makeTestPeerID(14))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	var receivedStatus ringpb.ExecutionStatus
	var receivedMu sync.Mutex
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		receivedMu.Lock()
		receivedStatus = msg.Status
		receivedMu.Unlock()
	})

	msgSuccess := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, primaryDON.ID)
	msgError := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, primaryDON.ID)

	makeBody := func(msg *ringpb.ExecutionStatusUpdate, senderIdx int) *remotetypes.MessageBody {
		return statusBody(t, localDON, primaryDON, primaryDON.Members[senderIdx], msg)
	}

	// Peer 0 reports SUCCESS — 1 vote for SUCCESS (quorum = 3, not reached).
	comm.Receive(ctx, makeBody(msgSuccess, 0))
	receivedMu.Lock()
	assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED, receivedStatus,
		"handler should not be called yet — no status has 2F+1 votes")
	receivedMu.Unlock()

	// Peer 1 reports SYSTEM_ERROR — 1 vote for SYSTEM_ERROR (still not reached).
	comm.Receive(ctx, makeBody(msgError, 1))
	receivedMu.Lock()
	assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED, receivedStatus,
		"handler should not be called — neither status has 2F+1 votes")
	receivedMu.Unlock()

	// Peer 2 reports SYSTEM_ERROR — 2 votes for SYSTEM_ERROR (still not enough).
	comm.Receive(ctx, makeBody(msgError, 2))
	receivedMu.Lock()
	assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED, receivedStatus,
		"handler should not be called — SYSTEM_ERROR has only 2 votes, need 3")
	receivedMu.Unlock()

	// Peer 3 reports SYSTEM_ERROR — 3 votes for SYSTEM_ERROR (2F+1 quorum reached).
	comm.Receive(ctx, makeBody(msgError, 3))

	receivedMu.Lock()
	assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, receivedStatus,
		"handler should receive SYSTEM_ERROR — it reached 2F+1 quorum")
	receivedMu.Unlock()
}

// A faulty primary peer must not be able to change fields of the delivered
// update (e.g. TriggerEventId) by casting the vote that completes quorum.
// Every field of the update has to be part of what the 2F+1 peers agree on.
func TestShardFailoverCommunicator_FaultyPeerCannotForgeTriggerEventID(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	disp := newInMemDispatcher()
	localDON := makeTestDON(2, 1, makeTestPeerID(20), makeTestPeerID(21), makeTestPeerID(22), makeTestPeerID(23))
	// N=4, F=1, quorum = 2F+1 = 3
	primaryDON := makeTestDON(1, 1, makeTestPeerID(10), makeTestPeerID(11), makeTestPeerID(12), makeTestPeerID(13))

	comm := NewShardFailoverCommunicator(disp, localDON.ID, logger.Test(t))
	require.NoError(t, comm.Start(ctx))
	t.Cleanup(func() { _ = comm.Close() })

	var received []*ringpb.ExecutionStatusUpdate
	var receivedMu sync.Mutex
	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		receivedMu.Lock()
		received = append(received, msg)
		receivedMu.Unlock()
	})

	honest := makeStatusMsg("wf-1", "evt-1", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, primaryDON.ID)
	// Same workflow and execution ID, different trigger event ID.
	forged := makeStatusMsg("wf-1", "evt-forged", "exec-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR, primaryDON.ID)

	makeBody := func(msg *ringpb.ExecutionStatusUpdate, senderIdx int) *remotetypes.MessageBody {
		return statusBody(t, localDON, primaryDON, primaryDON.Members[senderIdx], msg)
	}

	// Two honest votes, then the faulty peer casts the third vote.
	comm.Receive(ctx, makeBody(honest, 0))
	comm.Receive(ctx, makeBody(honest, 1))
	comm.Receive(ctx, makeBody(forged, 3))

	receivedMu.Lock()
	assert.Empty(t, received, "forged vote must not complete quorum for the honest update")
	receivedMu.Unlock()

	// Third honest vote reaches quorum.
	comm.Receive(ctx, makeBody(honest, 2))

	receivedMu.Lock()
	defer receivedMu.Unlock()
	require.Len(t, received, 1)
	assert.Equal(t, "evt-1", received[0].TriggerEventId)
}
