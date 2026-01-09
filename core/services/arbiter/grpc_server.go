package arbiter

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// TODO: Update this import path once proto is moved
	pb "github.com/smartcontractkit/chainlink/v2/core/services/arbiter/proto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// GRPCServer implements the ArbiterService gRPC interface.
type GRPCServer struct {
	pb.UnimplementedArbiterServiceServer
	state    *State
	decision DecisionEngine
	lggr     logger.Logger
}

// NewGRPCServer creates a new gRPC server instance.
func NewGRPCServer(state *State, decision DecisionEngine, lggr logger.Logger) *GRPCServer {
	return &GRPCServer{
		state:    state,
		decision: decision,
		lggr:     lggr.Named("GRPCServer"),
	}
}

// SubmitScaleIntent handles the SubmitScaleIntent RPC.
func (s *GRPCServer) SubmitScaleIntent(ctx context.Context, req *pb.ScaleIntentRequest) (*pb.ScaleIntentResponse, error) {
	// Validate request
	if req.DesiredReplicaCount < 1 {
		RecordRequest("SubmitScaleIntent", "INVALID_ARGUMENT")
		return nil, status.Error(codes.InvalidArgument, "DesiredReplicaCount must be at least 1")
	}

	// Convert protobuf types to internal types
	currentReplicas := make(map[string]ShardReplica)
	for name, replica := range req.CurrentReplicas {
		currentReplicas[name] = ShardReplica{
			Status:  protoStatusToString(replica.Status),
			Message: replica.Message,
			Metrics: replica.Metrics,
		}
	}

	// Update state
	s.state.Update(currentReplicas, int(req.DesiredReplicaCount), req.Reason)

	// Update metrics
	SetCurrentShardCount(len(currentReplicas))
	SetDesiredShardCount(int(req.DesiredReplicaCount))

	// Compute approved count using decision engine
	approved, err := s.decision.ComputeApprovedCount(ctx, int(req.DesiredReplicaCount))
	if err != nil {
		s.lggr.Errorw("failed to compute approved count",
			"error", err,
			"desiredCount", req.DesiredReplicaCount,
		)
		RecordRequest("SubmitScaleIntent", "INTERNAL")
		return nil, status.Error(codes.Internal, "failed to compute approved count")
	}

	// Update state and metrics with approved count
	s.state.SetApprovedCount(approved)
	SetApprovedShardCount(approved)

	RecordRequest("SubmitScaleIntent", "OK")

	s.lggr.Infow("Processed scale intent",
		"currentReplicasCount", len(currentReplicas),
		"desiredReplicasCount", req.DesiredReplicaCount,
		"approvedReplicasCount", approved,
		"reason", req.Reason,
	)

	return &pb.ScaleIntentResponse{Status: "ok"}, nil
}

// GetScalingSpec handles the GetScalingSpec RPC.
func (s *GRPCServer) GetScalingSpec(ctx context.Context, req *pb.GetScalingSpecRequest) (*pb.ScalingSpecResponse, error) {
	spec := s.state.GetScalingSpec()

	RecordRequest("GetScalingSpec", "OK")

	return &pb.ScalingSpecResponse{
		CurrentReplicaCount:  int32(spec.CurrentReplicaCount),  //nolint:gosec // G115: replica count bounded
		DesiredReplicaCount:  int32(spec.DesiredReplicaCount),  //nolint:gosec // G115: replica count bounded
		ApprovedReplicaCount: int32(spec.ApprovedReplicaCount), //nolint:gosec // G115: replica count bounded
		LastScalingReason:    spec.LastScalingReason,
	}, nil
}

// HealthCheck handles the HealthCheck RPC.
func (s *GRPCServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	RecordRequest("HealthCheck", "OK")
	return &pb.HealthCheckResponse{Status: "ok"}, nil
}

// protoStatusToString converts protobuf ReleaseStatus to string.
func protoStatusToString(status pb.ReleaseStatus) string {
	// Convert RELEASE_STATUS_INSTALLING -> INSTALLING
	name := status.String()
	return strings.TrimPrefix(name, "RELEASE_STATUS_")
}
