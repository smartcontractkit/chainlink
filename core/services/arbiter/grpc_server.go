package arbiter

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// GRPCServer implements the Arbiter gRPC interface from chainlink-common.
type GRPCServer struct {
	ringpb.UnimplementedArbiterServer
	shardConfig ShardConfigReader
	state       *State
	lggr        logger.Logger
}

// NewGRPCServer creates a new gRPC server instance.
func NewGRPCServer(shardConfig ShardConfigReader, state *State, lggr logger.Logger) *GRPCServer {
	return &GRPCServer{
		shardConfig: shardConfig,
		state:       state,
		lggr:        lggr.Named("GRPCServer"),
	}
}

// GetDesiredReplicas returns the desired number of shard replicas.
// This is called by the external scaler to determine how many shards to run.
// The desired count comes from the ShardConfig contract.
// The incoming shard status is stored in State for Ring OCR to query via ArbiterScaler.Status().
func (s *GRPCServer) GetDesiredReplicas(ctx context.Context, req *ringpb.ShardStatusRequest) (*ringpb.ArbiterResponse, error) {
	// Store incoming shard status in State so Ring OCR can access it via ArbiterScaler.Status()
	if s.state != nil && len(req.GetStatus()) > 0 {
		replicas := s.convertProtoStatusToReplicas(req.GetStatus())
		s.state.SetCurrentReplicas(replicas)
		s.lggr.Debugw("Updated shard status from scaler",
			"shardCount", len(replicas),
		)
	}

	// Get desired shard count from ShardConfig contract
	shardCount, err := s.shardConfig.GetDesiredShardCount(ctx)
	if err != nil {
		s.lggr.Errorw("Failed to get desired shard count",
			"error", err,
		)
		RecordRequest("GetDesiredReplicas", "INTERNAL")
		return nil, status.Error(codes.Internal, "failed to get desired shard count")
	}

	s.lggr.Debugw("GetDesiredReplicas called",
		"requestedShards", len(req.GetStatus()),
		"desiredShards", shardCount,
	)

	RecordRequest("GetDesiredReplicas", "OK")

	return &ringpb.ArbiterResponse{
		WantShards: uint32(shardCount), //nolint:gosec // G115: shard count bounded by contract
	}, nil
}

// convertProtoStatusToReplicas converts the protobuf shard status map to the internal replica format.
func (s *GRPCServer) convertProtoStatusToReplicas(protoStatus map[uint32]*ringpb.ShardStatus) map[string]ShardReplica {
	replicas := make(map[string]ShardReplica, len(protoStatus))
	for shardID, status := range protoStatus {
		replicaStatus := "UNHEALTHY"
		if status != nil && status.GetIsHealthy() {
			replicaStatus = StatusReady
		}
		replicas[strconv.FormatUint(uint64(shardID), 10)] = ShardReplica{
			Status:  replicaStatus,
			Message: "Reported by external scaler",
		}
	}
	return replicas
}
