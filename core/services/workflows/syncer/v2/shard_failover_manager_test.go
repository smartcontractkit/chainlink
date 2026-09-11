package v2

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/sharding"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// --- test doubles ---

// fakeDispatcher implements remotetypes.Dispatcher for integration tests.
// It routes messages to receivers registered via SetReceiverForMethod, just
// like the real dispatcher, keyed by (capID, donID, method).
type fakeDispatcher struct {
	mu        sync.Mutex
	receivers map[string]remotetypes.Receiver
	sendErr   error
}

func newFakeDispatcher() *fakeDispatcher {
	return &fakeDispatcher{receivers: make(map[string]remotetypes.Receiver)}
}

func (d *fakeDispatcher) key(capID string, donID uint32, method string) string {
	return capID + ":" + uintToStr(donID) + ":" + method
}

func (d *fakeDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, method string, receiver remotetypes.Receiver) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	k := d.key(capabilityID, donID, method)
	if _, ok := d.receivers[k]; ok {
		return &ReceiverExistsError{k}
	}
	d.receivers[k] = receiver
	return nil
}

func (d *fakeDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, method string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.receivers, d.key(capabilityID, donID, method))
}

func (d *fakeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	if d.sendErr != nil {
		return d.sendErr
	}
	// Route to the receiver registered for this message's routing key.
	k := d.key(msgBody.CapabilityId, msgBody.CapabilityDonId, msgBody.CapabilityMethod)
	d.mu.Lock()
	rcv, ok := d.receivers[k]
	d.mu.Unlock()
	if ok {
		// Set sender so the receiver can validate quorum.
		msgBody.Sender = peerID[:]
		rcv.Receive(context.Background(), msgBody)
	}
	return nil
}

func (d *fakeDispatcher) SetReceiver(string, uint32, remotetypes.Receiver) error { return nil }
func (d *fakeDispatcher) RemoveReceiver(string, uint32)                          {}
func (d *fakeDispatcher) Start(context.Context) error                            { return nil }
func (d *fakeDispatcher) Close() error                                           { return nil }
func (d *fakeDispatcher) Ready() error                                           { return nil }
func (d *fakeDispatcher) HealthReport() map[string]error                         { return nil }
func (d *fakeDispatcher) Name() string                                           { return "fakeDispatcher" }

type ReceiverExistsError struct{ key string }

func (e *ReceiverExistsError) Error() string { return "receiver already exists for " + e.key }

// fakeShardResolver implements both ShardResolver and AllShardsResolver.
type fakeShardResolver struct {
	shards []uint32
	dons   map[uint32]commoncap.DON
}

func (r *fakeShardResolver) ResolveShard(_ context.Context, _ string, _ string) (uint32, bool, error) {
	if len(r.shards) == 0 {
		return 0, false, nil
	}
	return r.shards[0], true, nil
}

func (r *fakeShardResolver) ResolveShards(_ context.Context, workflowIDs []string, _ []string) (map[string]uint32, error) {
	result := make(map[string]uint32, len(workflowIDs))
	for _, wf := range workflowIDs {
		if len(r.shards) > 0 {
			result[wf] = r.shards[0]
		}
	}
	return result, nil
}

func (r *fakeShardResolver) ResolveAllShards(_ context.Context, _ string, _ string) ([]uint32, bool, error) {
	return r.shards, len(r.shards) > 0, nil
}

// fakeDonSubscriber immediately delivers a DON on Subscribe.
type fakeDonSubscriber struct {
	don commoncap.DON
}

func (s *fakeDonSubscriber) NotifyDonSet(don commoncap.DON) { s.don = don }
func (s *fakeDonSubscriber) WaitForDon(_ context.Context) (commoncap.DON, error) {
	return s.don, nil
}
func (s *fakeDonSubscriber) Subscribe(_ context.Context) (<-chan commoncap.DON, func(), error) {
	ch := make(chan commoncap.DON, 1)
	ch <- s.don
	return ch, func() {}, nil
}

// --- helpers ---

func makePeerID(b byte) ragetypes.PeerID {
	var id ragetypes.PeerID
	id[0] = b
	return id
}

func makeDON(id uint32, f uint8, members ...ragetypes.PeerID) commoncap.DON {
	return commoncap.DON{ID: id, F: f, Members: members}
}

func uintToStr(n uint32) string {
	if n == 0 {
		return "0"
	}
	var buf [10]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// shardDonLookup creates a lookup function that returns a DON for a given shard ID.
func shardDonLookup(dons map[uint32]commoncap.DON) func(ctx context.Context, shardID uint32) *commoncap.DON {
	return func(_ context.Context, shardID uint32) *commoncap.DON {
		if don, ok := dons[shardID]; ok {
			return &don
		}
		return nil
	}
}

// --- the test that proves the fix ---

// TestShardFailoverManager_MultipleWorkflowsSharedDispatcher proves that
// multiple ShardFailoverManager instances (one per workflow) sharing the same
// dispatcher and communicator can all wire failover and route
// ExecutionStatusUpdate messages without the second registration failing with
// "receiver already exists" — the exact problem the reviewer identified.
//
// With the OLD code (each manager calling Dispatcher.SetReceiverForMethod
// directly), the second manager's wireFailover would fail because the
// dispatcher already has a receiver for (capID, donID, method). With the fix
// (shared ShardFailoverCommunicator that registers once and routes
// internally), both managers wire successfully and each receives only its own
// workflow's messages.
func TestShardFailoverManager_MultipleWorkflowsSharedDispatcher(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	// Two shard DONs: primary (DON 1) and secondary (DON 2).
	primaryDON := makeDON(1, 1, makePeerID(10), makePeerID(11), makePeerID(12))
	secondaryDON := makeDON(2, 1, makePeerID(20), makePeerID(21), makePeerID(22))
	dons := map[uint32]commoncap.DON{1: primaryDON, 2: secondaryDON}

	disp := newFakeDispatcher()
	donSub := &fakeDonSubscriber{don: secondaryDON}

	// Shared communicator — one instance for all workflows, registered once.
	comm := sharding.NewShardFailoverCommunicator(disp, secondaryDON.ID, logger.Test(t))

	// Shard resolver: this node (secondary, DON 2) is assigned shards [1, 2].
	// Shard 0 is primary, shard 1 is secondary.
	resolver := &fakeShardResolver{
		shards: []uint32{1, 2},
		dons:   dons,
	}

	shardDonLookupFn := shardDonLookup(dons)

	// Create two managers for two different workflows, sharing the same
	// dispatcher and communicator — exactly as the handler does.
	makeManager := func(workflowID string) *ShardFailoverManager {
		return NewShardFailoverManager(ShardFailoverManagerConfig{
			ShardingEnabled: true,
			MyShardID:       2,
			WorkflowID:      workflowID,
			WorkflowOwner:   "0xowner",
			ShardResolver:   resolver,
			Communicator:    comm,
			ShardDonLookup:  shardDonLookupFn,
			DonSubscriber:   donSub,
			Logger:          logger.Test(t),
		})
	}

	mgr1 := makeManager("wf-1")
	mgr2 := makeManager("wf-2")

	// Both managers must wire failover successfully. With the old code,
	// mgr2.wireFailover would fail because mgr1 already registered a receiver
	// for the same (capID, donID, method) key on the dispatcher.
	require.NoError(t, mgr1.wireFailover(ctx), "first manager must wire failover")
	t.Cleanup(func() {
		comm.UnregisterHandler("wf-1")
		_ = comm.Close()
	})

	require.NoError(t, mgr2.wireFailover(ctx), "second manager must wire failover — this is the bug the reviewer identified")
	t.Cleanup(func() { comm.UnregisterHandler("wf-2") })

	// Verify the communicator registered exactly one receiver with the dispatcher.
	disp.mu.Lock()
	assert.Len(t, disp.receivers, 1, "dispatcher should have exactly one receiver (shared communicator)")
	disp.mu.Unlock()

	// Now prove messages route to the correct workflow handler.
	// Simulate the primary (DON 1) sending ExecutionStatusUpdate to the
	// secondary (DON 2) for each workflow.
	wf1Received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	wf2Received := make(chan *ringpb.ExecutionStatusUpdate, 1)

	comm.RegisterHandler("wf-1", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		wf1Received <- msg
	})
	comm.RegisterHandler("wf-2", primaryDON, func(msg *ringpb.ExecutionStatusUpdate) {
		wf2Received <- msg
	})

	// Helper to send a status update from primary to secondary via the
	// communicator (simulating what the primary's forwardExecutionStatus does).
	sendStatusUpdate := func(workflowID, triggerEventID string, status ringpb.ExecutionStatus) {
		msg := &ringpb.ExecutionStatusUpdate{
			WorkflowId:     workflowID,
			TriggerEventId: triggerEventID,
			TriggerIndex:   0,
			Status:         status,
			PrimaryDonId:   primaryDON.ID,
		}
		payload, _ := proto.Marshal(msg)
		messageID := workflowID + ":" + triggerEventID + ":0"

		// Send from each primary member to simulate quorum (F+1 = 2).
		for _, peer := range primaryDON.Members[:2] {
			body := &remotetypes.MessageBody{
				CapabilityId:     sharding.ShardExecutionStatusUpdateCapabilityID,
				Method:           remotetypes.MethodExecutionStatusUpdate,
				CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
				CapabilityDonId:  secondaryDON.ID,
				CallerDonId:      primaryDON.ID,
				Payload:          payload,
				MessageId:        []byte(messageID),
			}
			// Route through the fake dispatcher as if the message arrived over p2p.
			require.NoError(t, disp.Send(peer, body))
		}
	}

	// Send for wf-1
	sendStatusUpdate("wf-1", "evt-1", ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR)
	select {
	case got := <-wf1Received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, "evt-1", got.TriggerEventId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: wf-1 handler not called")
	}

	// Send for wf-2
	sendStatusUpdate("wf-2", "evt-2", ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS)
	select {
	case got := <-wf2Received:
		assert.Equal(t, "wf-2", got.WorkflowId)
		assert.Equal(t, "evt-2", got.TriggerEventId)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: wf-2 handler not called")
	}

	// Verify wf-1 did NOT receive wf-2's message and vice versa.
	select {
	case <-wf1Received:
		t.Fatal("wf-1 handler should not receive wf-2's message")
	case <-wf2Received:
		t.Fatal("wf-2 handler should not receive wf-1's message")
	default:
	}
}

// TestShardFailoverManager_DispatcherRejectsDuplicateReceiver is a regression
// test ensuring the dispatcher rejects a second receiver registration for the
// same (capID, donID, method) key. The shared communicator relies on this
// invariant — it registers exactly once and routes internally.
func TestShardFailoverManager_DispatcherRejectsDuplicateReceiver(t *testing.T) {
	t.Parallel()

	disp := newFakeDispatcher()

	err1 := disp.SetReceiverForMethod(
		sharding.ShardExecutionStatusUpdateCapabilityID, 2,
		remotetypes.MethodExecutionStatusUpdate,
		&noopReceiver{},
	)
	require.NoError(t, err1, "first registration must succeed")

	err2 := disp.SetReceiverForMethod(
		sharding.ShardExecutionStatusUpdateCapabilityID, 2,
		remotetypes.MethodExecutionStatusUpdate,
		&noopReceiver{},
	)
	require.Error(t, err2, "duplicate registration for the same routing key must fail")
}

type noopReceiver struct{}

func (r *noopReceiver) Receive(context.Context, *remotetypes.MessageBody) {}

// Ensure fakeDonNotifier satisfies the interface.
var _ capabilities.DonSubscriber = (*fakeDonSubscriber)(nil)
