package ring

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

type mockArbiter struct {
	status *ringpb.ReplicaStatus
}

func (m *mockArbiter) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ringpb.ReplicaStatus, error) {
	if m.status != nil {
		return m.status, nil
	}
	return &ringpb.ReplicaStatus{}, nil
}

func (m *mockArbiter) ConsensusWantShards(ctx context.Context, req *ringpb.ConsensusWantShardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

var twoHealthyShards = []map[uint32]*ringpb.ShardStatus{
	{0: {IsHealthy: true}, 1: {IsHealthy: true}},
	{0: {IsHealthy: true}, 1: {IsHealthy: true}},
	{0: {IsHealthy: true}, 1: {IsHealthy: true}},
}

func toShardStatus(m map[uint32]bool) map[uint32]*ringpb.ShardStatus {
	result := make(map[uint32]*ringpb.ShardStatus, len(m))
	for k, v := range m {
		result[k] = &ringpb.ShardStatus{IsHealthy: v}
	}
	return result
}

func TestPlugin_Outcome(t *testing.T) {
	t.Run("WithMultiNodeObservations", func(t *testing.T) {
		lggr := logger.Test(t)
		store := NewStore()
		store.SetAllShardHealth(map[uint32]bool{0: true, 1: true, 2: true})

		config := ocr3types.ReportingPluginConfig{
			N: 4, F: 1,
			OffchainConfig:                          []byte{},
			MaxDurationObservation:                  0,
			MaxDurationShouldAcceptAttestedReport:   0,
			MaxDurationShouldTransmitAcceptedReport: 0,
		}

		plugin, err := NewPlugin(store, &mockArbiter{}, config, lggr, nil)
		require.NoError(t, err)

		ctx := t.Context()
		intialSeqNr := uint64(42)
		outcomeCtx := ocr3types.OutcomeContext{SeqNr: intialSeqNr}

		// Observations from 4 NOPs reporting health, workflows, and wantShards=3
		observations := []struct {
			name        string
			shardStatus map[uint32]*ringpb.ShardStatus
			workflows   []string
			wantShards  uint32
		}{
			{
				name:        "NOP 0",
				shardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 2: true}),
				workflows:   []string{"wf-A", "wf-B", "wf-C"},
				wantShards:  3,
			},
			{
				name:        "NOP 1",
				shardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 2: true}),
				workflows:   []string{"wf-B", "wf-C", "wf-D"},
				wantShards:  3,
			},
			{
				name:        "NOP 2",
				shardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 2: false}), // shard 2 unhealthy
				workflows:   []string{"wf-A", "wf-C"},
				wantShards:  3,
			},
			{
				name:        "NOP 3",
				shardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 2: true}),
				workflows:   []string{"wf-A", "wf-B", "wf-D"},
				wantShards:  3,
			},
		}

		// Build attributed observations
		aos := make([]types.AttributedObservation, 0)
		for idx, obs := range observations {
			pbObs := &ringpb.Observation{
				ShardStatus: obs.shardStatus,
				WorkflowIds: obs.workflows,
				Now:         timestamppb.Now(),
				WantShards:  obs.wantShards,
			}
			rawObs, marshalErr := proto.Marshal(pbObs)
			require.NoError(t, marshalErr)

			aos = append(aos, types.AttributedObservation{
				Observation: rawObs,
				Observer:    commontypes.OracleID(idx), //nolint:gosec // G115: idx bounded by observations slice
			})
		}

		// Execute Outcome phase
		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)
		require.NotNil(t, outcome)

		// Verify outcome
		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Check consensus results
		require.NotNil(t, outcomeProto.State)
		// When bootstrapping without PreviousOutcome, we use wantShards from observations (3)
		// Since consensus wantShards (3) equals bootstrap shards, no transition needed - ID stays the same
		require.Equal(t, intialSeqNr, outcomeProto.State.Id, "ID should match SeqNr (no transition needed)")
		t.Logf("Outcome - ID: %d, HealthyShards: %v", outcomeProto.State.Id, outcomeProto.State.GetRoutableShards())
		t.Logf("Workflows assigned: %d", len(outcomeProto.Routes))

		// Verify all workflows are assigned
		expectedWorkflows := map[string]bool{"wf-A": true, "wf-B": true, "wf-C": true, "wf-D": true}
		require.Len(t, outcomeProto.Routes, len(expectedWorkflows))
		for wf := range expectedWorkflows {
			route, exists := outcomeProto.Routes[wf]
			require.True(t, exists, "workflow %s should be assigned", wf)
			require.LessOrEqual(t, route.Shard, uint32(2), "shard should be healthy (0-2)")
			t.Logf("  %s → shard %d", wf, route.Shard)
		}

		// Verify determinism: run again, should get same assignments
		outcome2, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto2 := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome2, outcomeProto2)
		require.NoError(t, err)

		// Same workflows → same shards
		for wf, route1 := range outcomeProto.Routes {
			route2, exists := outcomeProto2.Routes[wf]
			require.True(t, exists)
			require.Equal(t, route1.Shard, route2.Shard, "workflow %s should assign to same shard", wf)
		}
	})
}

func TestPlugin_StateTransitions(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()

	config := ocr3types.ReportingPluginConfig{
		N: 4, F: 1,
	}

	// Use short time to sync for testing
	plugin, err := NewPlugin(store, &mockArbiter{}, config, lggr, &ConsensusConfig{
		BatchSize:  100,
		TimeToSync: 1 * time.Second,
	})
	require.NoError(t, err)

	ctx := t.Context()
	now := time.Now()

	// Test 1: Initial state with no previous outcome
	t.Run("initial_state", func(t *testing.T) {
		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           1,
			PreviousOutcome: nil,
		}

		// Only 1 healthy shard in observations with wantShards=1
		aos := makeObservationsWithWantShards(t, []map[uint32]*ringpb.ShardStatus{
			{0: {IsHealthy: true}},
			{0: {IsHealthy: true}},
			{0: {IsHealthy: true}},
		}, []string{"wf-1"}, now, 1)

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Should be in stable state with min shard count
		require.NotNil(t, outcomeProto.State.GetRoutableShards())
		require.Equal(t, uint32(1), outcomeProto.State.GetRoutableShards())
		t.Logf("Initial state: %d routable shards", outcomeProto.State.GetRoutableShards())
	})

	// Test 2: Transition triggered when wantShards changes
	t.Run("transition_triggered", func(t *testing.T) {
		// Start with 1 shard in stable state
		priorOutcome := &ringpb.Outcome{
			State: &ringpb.RoutingState{
				Id: 1,
				State: &ringpb.RoutingState_RoutableShards{
					RoutableShards: 1,
				},
			},
			Routes: map[string]*ringpb.WorkflowRoute{},
		}
		priorBytes, err := proto.Marshal(priorOutcome)
		require.NoError(t, err)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           2,
			PreviousOutcome: priorBytes,
		}

		// Observations show 2 healthy shards and wantShards=2
		aos := makeObservationsWithWantShards(t, twoHealthyShards, []string{"wf-1"}, now, 2)

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Should transition to Transition state
		transition := outcomeProto.State.GetTransition()
		require.NotNil(t, transition, "should be in transition state")
		require.Equal(t, uint32(2), transition.WantShards, "want 2 shards")
		require.Equal(t, uint32(1), transition.LastStableCount, "was at 1 shard")
		require.True(t, transition.ChangesSafeAfter.AsTime().After(now), "safety period should be in future")
		t.Logf("Transition: %d → %d, safe after %v", transition.LastStableCount, transition.WantShards, transition.ChangesSafeAfter.AsTime())
	})

	// Test 3: Stay in transition during safety period
	t.Run("stay_in_transition", func(t *testing.T) {
		safeAfter := now.Add(1 * time.Hour)
		priorOutcome := &ringpb.Outcome{
			State: &ringpb.RoutingState{
				Id: 2,
				State: &ringpb.RoutingState_Transition{
					Transition: &ringpb.Transition{
						WantShards:       2,
						LastStableCount:  1,
						ChangesSafeAfter: timestamppb.New(safeAfter),
					},
				},
			},
			Routes: map[string]*ringpb.WorkflowRoute{},
		}
		priorBytes, err := proto.Marshal(priorOutcome)
		require.NoError(t, err)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           3,
			PreviousOutcome: priorBytes,
		}

		// Still showing 2 healthy shards with wantShards=2, but safety period not elapsed
		aos := makeObservationsWithWantShards(t, twoHealthyShards, []string{"wf-1"}, now, 2)

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Should still be in transition state
		transition := outcomeProto.State.GetTransition()
		require.NotNil(t, transition, "should still be in transition")
		require.Equal(t, uint32(2), transition.WantShards)
		t.Logf("Still in transition, waiting for safety period")
	})

	// Test 4: Complete transition after safety period
	t.Run("complete_transition", func(t *testing.T) {
		safeAfter := now.Add(-1 * time.Second) // Safety period already passed
		priorOutcome := &ringpb.Outcome{
			State: &ringpb.RoutingState{
				Id: 2,
				State: &ringpb.RoutingState_Transition{
					Transition: &ringpb.Transition{
						WantShards:       2,
						LastStableCount:  1,
						ChangesSafeAfter: timestamppb.New(safeAfter),
					},
				},
			},
			Routes: map[string]*ringpb.WorkflowRoute{},
		}
		priorBytes, err := proto.Marshal(priorOutcome)
		require.NoError(t, err)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           3,
			PreviousOutcome: priorBytes,
		}

		aos := makeObservationsWithWantShards(t, twoHealthyShards, []string{"wf-1"}, now, 2)

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Should now be in stable state with 2 shards
		require.NotNil(t, outcomeProto.State.GetRoutableShards(), "should be in stable state")
		require.Equal(t, uint32(2), outcomeProto.State.GetRoutableShards())
		require.Equal(t, uint64(3), outcomeProto.State.Id, "state ID should increment")
		t.Logf("Transition complete: now at %d routable shards", outcomeProto.State.GetRoutableShards())
	})

	// Test 5: Stay stable when wantShards matches current
	t.Run("stay_stable", func(t *testing.T) {
		priorOutcome := &ringpb.Outcome{
			State: &ringpb.RoutingState{
				Id: 3,
				State: &ringpb.RoutingState_RoutableShards{
					RoutableShards: 2,
				},
			},
			Routes: map[string]*ringpb.WorkflowRoute{},
		}
		priorBytes, err := proto.Marshal(priorOutcome)
		require.NoError(t, err)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           4,
			PreviousOutcome: priorBytes,
		}

		// Same 2 healthy shards with wantShards=2
		aos := makeObservationsWithWantShards(t, twoHealthyShards, []string{"wf-1"}, now, 2)

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Should stay in stable state, ID unchanged
		require.NotNil(t, outcomeProto.State.GetRoutableShards())
		require.Equal(t, uint32(2), outcomeProto.State.GetRoutableShards())
		require.Equal(t, uint64(3), outcomeProto.State.Id, "state ID should not change when stable")
		t.Logf("Staying stable at %d routable shards", outcomeProto.State.GetRoutableShards())
	})
}

func makeObservationsWithWantShards(t *testing.T, shardStatuses []map[uint32]*ringpb.ShardStatus, workflows []string, now time.Time, wantShards uint32) []types.AttributedObservation {
	aos := make([]types.AttributedObservation, 0, len(shardStatuses))
	for i, status := range shardStatuses {
		pbObs := &ringpb.Observation{
			ShardStatus: status,
			WorkflowIds: workflows,
			Now:         timestamppb.New(now),
			WantShards:  wantShards,
		}
		rawObs, marshalErr := proto.Marshal(pbObs)
		require.NoError(t, marshalErr)

		aos = append(aos, types.AttributedObservation{
			Observation: rawObs,
			Observer:    commontypes.OracleID(i), //nolint:gosec // G115: i bounded by shardStatuses slice
		})
	}
	return aos
}

func TestPlugin_NewPlugin_NilArbiter(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	config := ocr3types.ReportingPluginConfig{N: 4, F: 1}

	_, err := NewPlugin(store, nil, config, lggr, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "RingOCR arbiterScaler is required")
}

func TestPlugin_getHealthyShards(t *testing.T) {
	tests := []struct {
		name  string
		votes map[uint32]int // shardID -> vote count
		f     int
		want  int
	}{
		{"all healthy", map[uint32]int{0: 2, 1: 2, 2: 2}, 1, 3},
		{"some unhealthy", map[uint32]int{0: 2, 1: 1, 2: 2}, 1, 2},
		{"none healthy", map[uint32]int{0: 1, 1: 1}, 1, 0},
		{"higher F threshold", map[uint32]int{0: 3, 1: 2, 2: 3}, 2, 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := &Plugin{
				store:  NewStore(),
				config: ocr3types.ReportingPluginConfig{F: tc.f},
			}
			got := plugin.getHealthyShards(tc.votes)
			require.Len(t, got, tc.want)
		})
	}
}

func TestPlugin_NoHealthyShardsFallbackToShardZero(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()

	// Set all shards unhealthy - store starts in transition state
	store.SetAllShardHealth(map[uint32]bool{0: false, 1: false, 2: false})

	config := ocr3types.ReportingPluginConfig{
		N: 4, F: 1,
	}

	arbiter := &mockArbiter{}
	plugin, err := NewPlugin(store, arbiter, config, lggr, &ConsensusConfig{
		BatchSize:  100,
		TimeToSync: 1 * time.Second,
	})
	require.NoError(t, err)

	transmitter := NewTransmitter(lggr, store, nil, arbiter, "test-account")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Start a goroutine that requests allocation (will block waiting for OCR)
	resultCh := make(chan uint32)
	errCh := make(chan error, 1)
	go func() {
		shard, shardErr := store.GetShardForWorkflow(ctx, "workflow-123")
		if shardErr != nil {
			errCh <- shardErr
			return
		}
		resultCh <- shard
	}()

	// Give goroutine time to enqueue request
	time.Sleep(10 * time.Millisecond)

	// Verify request is pending for OCR consensus
	pending := store.GetPendingAllocations()
	require.Contains(t, pending, "workflow-123")

	// Simulate OCR round with observations showing no healthy shards
	// The pending allocation "workflow-123" should be included in observation
	now := time.Now()
	aos := make([]types.AttributedObservation, 3)
	for i := 0; i < 3; i++ {
		pbObs := &ringpb.Observation{
			ShardStatus: toShardStatus(map[uint32]bool{0: false, 1: false, 2: false}),
			WorkflowIds: []string{"workflow-123"},
			Now:         timestamppb.New(now),
		}
		rawObs, marshalErr := proto.Marshal(pbObs)
		require.NoError(t, marshalErr)
		aos[i] = types.AttributedObservation{
			Observation: rawObs,
			Observer:    commontypes.OracleID(i), //nolint:gosec // G115: i bounded by loop
		}
	}

	// Use a previous outcome in steady state so we can test the fallback
	priorOutcome := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		},
		Routes: map[string]*ringpb.WorkflowRoute{},
	}
	priorBytes, err := proto.Marshal(priorOutcome)
	require.NoError(t, err)

	outcomeCtx := ocr3types.OutcomeContext{
		SeqNr:           2,
		PreviousOutcome: priorBytes,
	}

	// Run plugin Outcome phase
	outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
	require.NoError(t, err)

	// Transmit the outcome (applies routes to store)
	reports, err := plugin.Reports(ctx, 2, outcome)
	require.NoError(t, err)
	require.Len(t, reports, 1)

	err = transmitter.Transmit(ctx, types.ConfigDigest{}, 2, reports[0].ReportWithInfo, nil)
	require.NoError(t, err)

	// Blocked goroutine should now receive result from OCR - should be shard 0 (fallback)
	select {
	case shard := <-resultCh:
		require.Equal(t, uint32(0), shard, "should fallback to shard 0 when no healthy shards")
	case recvErr := <-errCh:
		t.Fatalf("unexpected error: %v", recvErr)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("allocation was not fulfilled by OCR")
	}

	// Verify the outcome assigned workflow-123 to shard 0
	outcomeProto := &ringpb.Outcome{}
	err = proto.Unmarshal(outcome, outcomeProto)
	require.NoError(t, err)

	route, exists := outcomeProto.Routes["workflow-123"]
	require.True(t, exists, "workflow-123 should be in routes")
	require.Equal(t, uint32(0), route.Shard, "workflow-123 should be assigned to shard 0 (fallback)")
}

func TestPlugin_ObservationQuorum(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	config := ocr3types.ReportingPluginConfig{N: 4, F: 1}
	plugin, err := NewPlugin(store, &mockArbiter{}, config, lggr, nil)
	require.NoError(t, err)

	ctx := context.Background()
	outctx := ocr3types.OutcomeContext{}

	t.Run("quorum_reached", func(t *testing.T) {
		// Need 2F+1 = 3 observations for quorum with N=4, F=1
		aos := make([]types.AttributedObservation, 3)
		for i := range aos {
			aos[i] = types.AttributedObservation{Observer: commontypes.OracleID(i)} //nolint:gosec // G115: i bounded by slice
		}

		quorum, qErr := plugin.ObservationQuorum(ctx, outctx, nil, aos)
		require.NoError(t, qErr)
		require.True(t, quorum)
	})

	t.Run("quorum_not_reached", func(t *testing.T) {
		// Only 2 observations - not enough for quorum
		aos := make([]types.AttributedObservation, 2)
		for i := range aos {
			aos[i] = types.AttributedObservation{Observer: commontypes.OracleID(i)} //nolint:gosec // G115: i bounded by slice
		}

		quorum, qErr := plugin.ObservationQuorum(ctx, outctx, nil, aos)
		require.NoError(t, qErr)
		require.False(t, quorum)
	})

	t.Run("exact_quorum", func(t *testing.T) {
		// Exactly 2F+1 = 3 observations
		aos := make([]types.AttributedObservation, 3)
		for i := range aos {
			aos[i] = types.AttributedObservation{Observer: commontypes.OracleID(i)} //nolint:gosec // G115: i bounded by slice
		}

		quorum, qErr := plugin.ObservationQuorum(ctx, outctx, nil, aos)
		require.NoError(t, qErr)
		require.True(t, quorum)
	})

	t.Run("all_observations", func(t *testing.T) {
		// All N=4 observations
		aos := make([]types.AttributedObservation, 4)
		for i := range aos {
			aos[i] = types.AttributedObservation{Observer: commontypes.OracleID(i)} //nolint:gosec // G115: i bounded by slice
		}

		quorum, qErr := plugin.ObservationQuorum(ctx, outctx, nil, aos)
		require.NoError(t, qErr)
		require.True(t, quorum)
	})
}

func TestPlugin_ShardOrchestratorIntegration(t *testing.T) {
	lggr := logger.Test(t)

	// Create both stores
	ringStore := NewStore()
	orchestratorStore := shardorchestrator.NewStore(lggr)

	// Initialize ring store with healthy shards
	ringStore.SetAllShardHealth(map[uint32]bool{0: true, 1: true, 2: true})

	config := ocr3types.ReportingPluginConfig{
		N: 4, F: 1,
	}

	arbiter := &mockArbiter{}
	plugin, err := NewPlugin(ringStore, arbiter, config, lggr, &ConsensusConfig{
		BatchSize:  100,
		TimeToSync: 1 * time.Second,
	})
	require.NoError(t, err)

	// Create transmitter with both stores
	transmitter := NewTransmitter(lggr, ringStore, orchestratorStore, arbiter, "test-account")

	ctx := context.Background()
	now := time.Now()

	t.Run("initial_workflow_assignments", func(t *testing.T) {
		// Create observations with workflows
		workflows := []string{"wf-A", "wf-B", "wf-C"}
		aos := makeObservationsWithWantShards(t, []map[uint32]*ringpb.ShardStatus{
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
		}, workflows, now, 3)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           1,
			PreviousOutcome: nil,
		}

		// Generate outcome
		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, aos)
		require.NoError(t, err)

		// Generate report and transmit
		reports, err := plugin.Reports(ctx, 1, outcome)
		require.NoError(t, err)
		require.Len(t, reports, 1)

		err = transmitter.Transmit(ctx, types.ConfigDigest{}, 1, reports[0].ReportWithInfo, nil)
		require.NoError(t, err)

		// Verify ring store was updated
		for _, wf := range workflows {
			shard, err := ringStore.GetShardForWorkflow(ctx, wf)
			require.NoError(t, err)
			require.LessOrEqual(t, shard, uint32(2), "workflow should be assigned to valid shard")
			t.Logf("Ring store: %s → shard %d", wf, shard)
		}

		// Verify orchestrator store was updated with correct state
		for _, wf := range workflows {
			mapping, err := orchestratorStore.GetWorkflowMapping(ctx, wf)
			require.NoError(t, err)
			require.Equal(t, wf, mapping.WorkflowID)
			require.LessOrEqual(t, mapping.NewShardID, uint32(2))
			require.Equal(t, uint32(0), mapping.OldShardID, "initial assignment should have oldShardID=0")
			require.Equal(t, shardorchestrator.StateSteady, mapping.TransitionState, "initial assignment should be steady")
			t.Logf("Orchestrator store: %s → shard %d (state: %s)", wf, mapping.NewShardID, mapping.TransitionState.String())
		}

		// Verify version tracking
		version := orchestratorStore.GetMappingVersion()
		require.Equal(t, uint64(1), version, "version should increment after first update")
	})

	t.Run("workflow_transition_detected", func(t *testing.T) {
		// First, establish a baseline with workflows distributed across 3 shards
		// Use wantShards=3 to ensure workflows actually get assigned to shard 2
		baselineAos := makeObservationsWithWantShards(t, []map[uint32]*ringpb.ShardStatus{
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}, 2: {IsHealthy: true}},
		}, []string{"wf-A", "wf-B", "wf-C", "wf-D", "wf-E"}, now, 3)

		baselineOutcome, err := plugin.Outcome(ctx, ocr3types.OutcomeContext{SeqNr: 2}, nil, baselineAos)
		require.NoError(t, err)

		baselineReports, err := plugin.Reports(ctx, 2, baselineOutcome)
		require.NoError(t, err)

		err = transmitter.Transmit(ctx, types.ConfigDigest{}, 2, baselineReports[0].ReportWithInfo, nil)
		require.NoError(t, err)

		// Parse baseline to see which workflows were on shard 2
		baselineProto := &ringpb.Outcome{}
		err = proto.Unmarshal(baselineOutcome, baselineProto)
		require.NoError(t, err)

		workflowsOnShard2 := []string{}
		for wfID, route := range baselineProto.Routes {
			if route.Shard == 2 {
				workflowsOnShard2 = append(workflowsOnShard2, wfID)
			}
			t.Logf("Baseline: %s on shard %d", wfID, route.Shard)
		}
		require.NotEmpty(t, workflowsOnShard2, "at least one workflow should be on shard 2 for this test")

		// Now scale down to 2 shards - workflows on shard 2 MUST move
		transitionAos := makeObservationsWithWantShards(t, []map[uint32]*ringpb.ShardStatus{
			{0: {IsHealthy: true}, 1: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}},
			{0: {IsHealthy: true}, 1: {IsHealthy: true}},
		}, []string{"wf-A", "wf-B", "wf-C", "wf-D", "wf-E"}, now, 2)

		outcomeCtx := ocr3types.OutcomeContext{
			SeqNr:           3,
			PreviousOutcome: baselineOutcome,
		}

		outcome, err := plugin.Outcome(ctx, outcomeCtx, nil, transitionAos)
		require.NoError(t, err)

		reports, err := plugin.Reports(ctx, 3, outcome)
		require.NoError(t, err)

		err = transmitter.Transmit(ctx, types.ConfigDigest{}, 3, reports[0].ReportWithInfo, nil)
		require.NoError(t, err)

		// Verify orchestrator store shows transition state for workflows that moved from shard 2
		outcomeProto := &ringpb.Outcome{}
		err = proto.Unmarshal(outcome, outcomeProto)
		require.NoError(t, err)

		// Workflows that were on shard 2 must have moved and should show TransitionState
		for _, wfID := range workflowsOnShard2 {
			mapping, err := orchestratorStore.GetWorkflowMapping(ctx, wfID)
			require.NoError(t, err)

			newRoute := outcomeProto.Routes[wfID]
			require.NotEqual(t, uint32(2), newRoute.Shard, "workflow should have moved from shard 2")
			require.Equal(t, shardorchestrator.StateTransitioning, mapping.TransitionState,
				"workflow %s moved from shard 2 to shard %d, should be transitioning", wfID, newRoute.Shard)
			require.Equal(t, uint32(2), mapping.OldShardID, "should track old shard")
			require.Equal(t, newRoute.Shard, mapping.NewShardID, "should track new shard")
			t.Logf("Workflow %s transitioned: shard 2 → %d", wfID, newRoute.Shard)
		}

		// Verify version incremented
		version := orchestratorStore.GetMappingVersion()
		require.Equal(t, uint64(3), version, "version should increment after update")
	})
}
