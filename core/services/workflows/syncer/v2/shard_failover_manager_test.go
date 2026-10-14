package v2

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/sharding"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
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
	primaryDON := makeDON(1, 1, makePeerID(10), makePeerID(11), makePeerID(12), makePeerID(13))
	secondaryDON := makeDON(2, 1, makePeerID(20), makePeerID(21), makePeerID(22), makePeerID(23))
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
			MyShardIndex:    2,
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

		// Send from each primary member to simulate quorum (2F+1 = 3).
		for _, peer := range primaryDON.Members[:3] {
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

// newTestShardFailoverManager wraps engine in a manager that is never started,
// so the sharding collaborators are left nil.
func newTestShardFailoverManager(t *testing.T, engine v2.WorkflowEngine) *ShardFailoverManager {
	t.Helper()
	m := NewShardFailoverManager(ShardFailoverManagerConfig{
		ShardingEnabled: true,
		WorkflowID:      "wf-1",
		WorkflowOwner:   "0xowner",
		Logger:          logger.Test(t),
	})
	m.SetEngine(engine)
	return m
}

// TestShardFailoverManager_DelegatesIsCoordinated covers the classification of
// the wrapper a sharded node registers, which must follow the engine it wraps.
func TestShardFailoverManager_DelegatesIsCoordinated(t *testing.T) {
	t.Parallel()

	t.Run("wrapping a legacy engine reports legacy", func(t *testing.T) {
		t.Parallel()
		m := newTestShardFailoverManager(t, &mockEngine{})
		require.False(t, m.IsCoordinated())
	})

	t.Run("wrapping a coordinated engine reports coordinated", func(t *testing.T) {
		t.Parallel()
		m := newTestShardFailoverManager(t, &fakeCoordinatedDrainableEngine{})
		require.True(t, m.IsCoordinated())
	})
}

const (
	testMyDONID    = uint32(1)
	testOtherDONID = uint32(2)
)

// recordingEngine is a coordinated v2.WorkflowEngine that records the events
// delivered to ExecuteTrigger. Only the methods the coordinator path touches
// are implemented.
type recordingEngine struct {
	v2.WorkflowEngine
	subs    []*sdkpb.TriggerSubscription
	subErr  error
	execErr error

	mu       sync.Mutex
	executed []string
}

func (e *recordingEngine) IsCoordinated() bool { return true }

func (e *recordingEngine) Subscribe(context.Context) ([]*sdkpb.TriggerSubscription, error) {
	return e.subs, e.subErr
}

func (e *recordingEngine) ExecuteTrigger(_ context.Context, event triggers.CoordinatedEvent) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executed = append(e.executed, event.Event.Event.ID)
	return e.execErr
}

// executedIDs returns the IDs of the events delivered so far, in order.
func (e *recordingEngine) executedIDs() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.executed...)
}

// erroringShardResolver fails every lookup, as an unreachable orchestrator would.
type erroringShardResolver struct{}

func (erroringShardResolver) ResolveShard(context.Context, string, string) (uint32, bool, error) {
	return 0, false, errors.New("orchestrator unavailable")
}

func (erroringShardResolver) ResolveShards(context.Context, []string, []string) (map[string]uint32, error) {
	return nil, errors.New("orchestrator unavailable")
}

var _ shardownership.ShardResolver = erroringShardResolver{}

var (
	// ownerResolver makes this node (testMyDONID) the primary shard.
	ownerResolver = &fakeShardResolver{shards: []uint32{testMyDONID, testOtherDONID}}
	// secondaryResolver makes this node the secondary shard: listed, but not the owner.
	secondaryResolver = &fakeShardResolver{shards: []uint32{testOtherDONID, testMyDONID}}
)

// newAdmissionTestManager wraps engine in a manager for testMyDONID that is
// never started. cfg supplies the sharding and failover settings under test.
func newAdmissionTestManager(t *testing.T, engine v2.WorkflowEngine, cfg ShardFailoverManagerConfig) *ShardFailoverManager {
	t.Helper()
	cfg.MyDONID = testMyDONID
	cfg.WorkflowID = "wf-1"
	cfg.WorkflowOwner = "0xowner"
	cfg.Logger = logger.Test(t)
	m := NewShardFailoverManager(cfg)
	m.SetEngine(engine)
	return m
}

func testCoordinatedEvent(id string) triggers.CoordinatedEvent {
	return triggers.CoordinatedEvent{
		WorkflowID:   "wf-1",
		TriggerCapID: "cron-trigger@1.0.0",
		Event:        commoncap.TriggerResponse{Event: commoncap.TriggerEvent{ID: id}},
	}
}

// TestShardFailoverManager_Subscribe covers trigger registration on the
// coordinated path: the coordinator subscribes through the manager, which must
// hand back the wrapped engine's subscriptions.
func TestShardFailoverManager_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("returns the engine's subscriptions", func(t *testing.T) {
		t.Parallel()
		want := []*sdkpb.TriggerSubscription{{Id: "cron-trigger@1.0.0", Method: "Trigger"}}
		m := newTestShardFailoverManager(t, &recordingEngine{subs: want})

		got, err := m.Subscribe(t.Context())
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("passes through the engine's error", func(t *testing.T) {
		t.Parallel()
		wantErr := errors.New("subscribe failed")
		m := newTestShardFailoverManager(t, &recordingEngine{subErr: wantErr})

		_, err := m.Subscribe(t.Context())
		require.ErrorIs(t, err, wantErr)
	})
}

// TestShardFailoverManager_ExecuteTrigger covers admission on the coordinated
// path. The coordinated engine never calls OnTriggerAdmission, so the manager
// must check shard ownership itself before delivering. It only decides: an
// admission error is returned for the coordinator to act on (ACK a denial,
// leave a cached event un-ACKed for replay).
func TestShardFailoverManager_ExecuteTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cfg          ShardFailoverManagerConfig
		execErr      error
		wantErr      error
		wantExecuted bool
		wantCached   bool
	}{
		{
			name:         "sharding disabled delivers",
			cfg:          ShardFailoverManagerConfig{ShardingEnabled: false, ShardResolver: secondaryResolver},
			wantExecuted: true,
		},
		{
			name:         "owner delivers",
			cfg:          ShardFailoverManagerConfig{ShardingEnabled: true, ShardResolver: ownerResolver},
			wantExecuted: true,
		},
		{
			name: "secondary with failover enabled caches",
			cfg: ShardFailoverManagerConfig{
				ShardingEnabled: true, ShardResolver: secondaryResolver, FailoverGate: limits.NewGateLimiter(true),
			},
			wantErr:    triggers.ErrAdmissionCached,
			wantCached: true,
		},
		{
			name: "non-owner with failover disabled is denied",
			cfg: ShardFailoverManagerConfig{
				ShardingEnabled: true, ShardResolver: secondaryResolver, FailoverGate: limits.NewGateLimiter(false),
			},
			wantErr: triggers.ErrAdmissionDenied,
		},
		{
			name:    "non-owner without a failover gate is denied",
			cfg:     ShardFailoverManagerConfig{ShardingEnabled: true, ShardResolver: secondaryResolver},
			wantErr: triggers.ErrAdmissionDenied,
		},
		{
			name:    "orchestrator error is denied",
			cfg:     ShardFailoverManagerConfig{ShardingEnabled: true, ShardResolver: erroringShardResolver{}},
			wantErr: triggers.ErrAdmissionDenied,
		},
		{
			name:         "engine error passes through",
			cfg:          ShardFailoverManagerConfig{ShardingEnabled: true, ShardResolver: ownerResolver},
			execErr:      v2.ErrEngineDraining,
			wantErr:      v2.ErrEngineDraining,
			wantExecuted: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			engine := &recordingEngine{execErr: tc.execErr}
			m := newAdmissionTestManager(t, engine, tc.cfg)

			err := m.ExecuteTrigger(t.Context(), testCoordinatedEvent("evt-1"))
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			var wantIDs []string
			if tc.wantExecuted {
				wantIDs = []string{"evt-1"}
			}
			assert.Equal(t, wantIDs, engine.executedIDs(), "events delivered to the engine")
			_, gotCached := m.cache["evt-1"]
			assert.Equal(t, tc.wantCached, gotCached, "event cached for failover")
		})
	}
}

// TestShardFailoverManager_CachedEventOnStatusUpdate covers failover on the
// coordinated path: an event the secondary cached in ExecuteTrigger is dropped
// when the primary reports an outcome, and replayed on a system error. Replay
// goes straight to the wrapped engine — the secondary is still not the owner,
// so routing it back through admission would cache it again.
func TestShardFailoverManager_CachedEventOnStatusUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       ringpb.ExecutionStatus
		wantExecuted []string
		wantCached   bool
	}{
		{name: "success drops the event", status: ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS},
		{name: "user error drops the event", status: ringpb.ExecutionStatus_EXECUTION_STATUS_USER_ERROR},
		{
			name:         "system error replays the event",
			status:       ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR,
			wantExecuted: []string{"evt-1"},
		},
		{
			name:       "unspecified keeps the event cached",
			status:     ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED,
			wantCached: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			engine := &recordingEngine{}
			m := newAdmissionTestManager(t, engine, ShardFailoverManagerConfig{
				ShardingEnabled: true, ShardResolver: secondaryResolver, FailoverGate: limits.NewGateLimiter(true),
			})
			require.ErrorIs(t, m.ExecuteTrigger(t.Context(), testCoordinatedEvent("evt-1")), triggers.ErrAdmissionCached)

			m.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{
				WorkflowId:     "wf-1",
				ExecutionId:    "exec-1",
				TriggerEventId: "evt-1",
				Status:         tc.status,
			})

			assert.Equal(t, tc.wantExecuted, engine.executedIDs(), "events delivered to the engine")
			_, gotCached := m.cache["evt-1"]
			assert.Equal(t, tc.wantCached, gotCached, "event cached for failover")
		})
	}
}
