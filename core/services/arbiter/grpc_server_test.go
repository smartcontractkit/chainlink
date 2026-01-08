package arbiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/smartcontractkit/chainlink/v2/core/services/arbiter/proto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockDecisionEngine is a mock implementation of DecisionEngine for testing.
type mockDecisionEngine struct {
	approvedCount int
	err           error
}

func (m *mockDecisionEngine) ComputeApprovedCount(ctx context.Context, desiredCount int) (int, error) {
	return m.approvedCount, m.err
}

func TestGRPCServer_SubmitScaleIntent_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{approvedCount: 5}

	server := NewGRPCServer(state, mockDecision, lggr)

	req := &pb.ScaleIntentRequest{
		CurrentReplicas: map[string]*pb.ShardReplica{
			"shard-0": {Status: pb.ReleaseStatus_RELEASE_STATUS_READY, Message: "Running"},
			"shard-1": {Status: pb.ReleaseStatus_RELEASE_STATUS_READY, Message: "Running"},
		},
		DesiredReplicaCount: 5,
		Reason:              "high CPU utilization",
	}

	resp, err := server.SubmitScaleIntent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)

	// Verify state was updated
	spec := state.GetScalingSpec()
	assert.Equal(t, 2, spec.CurrentReplicaCount)
	assert.Equal(t, 5, spec.DesiredReplicaCount)
	assert.Equal(t, 5, spec.ApprovedReplicaCount)
	assert.Equal(t, "high CPU utilization", spec.LastScalingReason)
}

func TestGRPCServer_SubmitScaleIntent_InvalidArgument(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{approvedCount: 1}

	server := NewGRPCServer(state, mockDecision, lggr)

	tests := []struct {
		name    string
		request *pb.ScaleIntentRequest
	}{
		{
			name: "desired count is zero",
			request: &pb.ScaleIntentRequest{
				CurrentReplicas:     map[string]*pb.ShardReplica{},
				DesiredReplicaCount: 0,
				Reason:              "test",
			},
		},
		{
			name: "desired count is negative",
			request: &pb.ScaleIntentRequest{
				CurrentReplicas:     map[string]*pb.ShardReplica{},
				DesiredReplicaCount: -1,
				Reason:              "test",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.SubmitScaleIntent(context.Background(), tc.request)

			require.Error(t, err)
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

func TestGRPCServer_SubmitScaleIntent_DecisionEngineError(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{
		err: errors.New("contract read failed"),
	}

	server := NewGRPCServer(state, mockDecision, lggr)

	req := &pb.ScaleIntentRequest{
		CurrentReplicas:     map[string]*pb.ShardReplica{},
		DesiredReplicaCount: 5,
		Reason:              "test",
	}

	resp, err := server.SubmitScaleIntent(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetScalingSpec(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{approvedCount: 5}

	// Set up some state
	replicas := map[string]ShardReplica{
		"shard-0": {Status: "READY", Message: "Running"},
		"shard-1": {Status: "READY", Message: "Running"},
		"shard-2": {Status: "INSTALLING", Message: "In progress"},
	}
	state.Update(replicas, 10, "scale-up request")
	state.SetApprovedCount(8)

	server := NewGRPCServer(state, mockDecision, lggr)

	resp, err := server.GetScalingSpec(context.Background(), &pb.GetScalingSpecRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int32(3), resp.CurrentReplicaCount)
	assert.Equal(t, int32(10), resp.DesiredReplicaCount)
	assert.Equal(t, int32(8), resp.ApprovedReplicaCount)
	assert.Equal(t, "scale-up request", resp.LastScalingReason)
}

func TestGRPCServer_GetScalingSpec_InitialState(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{}

	server := NewGRPCServer(state, mockDecision, lggr)

	resp, err := server.GetScalingSpec(context.Background(), &pb.GetScalingSpecRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int32(0), resp.CurrentReplicaCount)
	assert.Equal(t, int32(1), resp.DesiredReplicaCount)
	assert.Equal(t, int32(1), resp.ApprovedReplicaCount)
	assert.Equal(t, "Initial state", resp.LastScalingReason)
}

func TestGRPCServer_HealthCheck(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{}

	server := NewGRPCServer(state, mockDecision, lggr)

	resp, err := server.HealthCheck(context.Background(), &pb.HealthCheckRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)
}

func TestGRPCServer_SubmitScaleIntent_WithMetrics(t *testing.T) {
	lggr := logger.TestLogger(t)
	state := NewState()
	mockDecision := &mockDecisionEngine{approvedCount: 3}

	server := NewGRPCServer(state, mockDecision, lggr)

	req := &pb.ScaleIntentRequest{
		CurrentReplicas: map[string]*pb.ShardReplica{
			"shard-0": {
				Status:  pb.ReleaseStatus_RELEASE_STATUS_READY,
				Message: "Running",
				Metrics: map[string]float64{
					"cpu_usage":    0.75,
					"memory_usage": 0.60,
				},
			},
		},
		DesiredReplicaCount: 3,
		Reason:              "metrics-based scaling",
	}

	resp, err := server.SubmitScaleIntent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)

	spec := state.GetScalingSpec()
	assert.Equal(t, 1, spec.CurrentReplicaCount)
	assert.Equal(t, 3, spec.DesiredReplicaCount)
	assert.Equal(t, 3, spec.ApprovedReplicaCount)
}

func TestGRPCServer_SubmitScaleIntent_AllReleaseStatuses(t *testing.T) {
	lggr := logger.TestLogger(t)

	statuses := []pb.ReleaseStatus{
		pb.ReleaseStatus_RELEASE_STATUS_INSTALLING,
		pb.ReleaseStatus_RELEASE_STATUS_INSTALL_SUCCEEDED,
		pb.ReleaseStatus_RELEASE_STATUS_INSTALL_FAILED,
		pb.ReleaseStatus_RELEASE_STATUS_READY,
		pb.ReleaseStatus_RELEASE_STATUS_DEGRADED,
		pb.ReleaseStatus_RELEASE_STATUS_UNKNOWN,
	}

	for _, s := range statuses {
		t.Run(s.String(), func(t *testing.T) {
			state := NewState()
			mockDecision := &mockDecisionEngine{approvedCount: 1}
			server := NewGRPCServer(state, mockDecision, lggr)

			req := &pb.ScaleIntentRequest{
				CurrentReplicas: map[string]*pb.ShardReplica{
					"shard-0": {
						Status:  s,
						Message: "test",
					},
				},
				DesiredReplicaCount: 1,
				Reason:              "status test",
			}

			resp, err := server.SubmitScaleIntent(context.Background(), req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "ok", resp.Status)
		})
	}
}
