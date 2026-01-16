package arbiter

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ringpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// GRPCServer implements the Arbiter gRPC interface from chainlink-common.
type GRPCServer struct {
	ringpb.UnimplementedArbiterServer
	shardConfig ShardConfigReader
	lggr        logger.Logger
}

// NewGRPCServer creates a new gRPC server instance.
func NewGRPCServer(shardConfig ShardConfigReader, lggr logger.Logger) *GRPCServer {
	return &GRPCServer{
		shardConfig: shardConfig,
		lggr:        lggr.Named("GRPCServer"),
	}
}

// GetDesiredReplicas returns the desired number of shard replicas.
// This is called by the external scaler to determine how many shards to run.
// The desired count comes from the ShardConfig contract.
func (s *GRPCServer) GetDesiredReplicas(ctx context.Context, req *ringpb.ShardStatusRequest) (*ringpb.ArbiterResponse, error) {
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
