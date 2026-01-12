package arbiter

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"
)

// ArbiterScalerHandler implements the ArbiterScalerServer interface from chainlink-common.
// This allows the Ring consensus to communicate with the Arbiter about shard scaling.
type ArbiterScalerHandler struct {
	pb.UnimplementedArbiterScalerServer
	state *State
	lggr  logger.Logger
}

// NewArbiterScalerHandler creates a new ArbiterScalerHandler.
func NewArbiterScalerHandler(state *State, lggr logger.Logger) *ArbiterScalerHandler {
	return &ArbiterScalerHandler{
		state: state,
		lggr:  logger.Named(lggr, "ArbiterScalerHandler"),
	}
}

// Status returns the current replica status for Ring OCR routing.
// Returns only READY shards count and per-shard health status.
// This is called by the Ring plugin to determine which shards can receive traffic.
func (h *ArbiterScalerHandler) Status(ctx context.Context, _ *emptypb.Empty) (*pb.ReplicaStatus, error) {
	routable := h.state.GetRoutableShards()

	h.lggr.Debugw("Status requested",
		"readyShards", routable.ReadyCount,
		"totalShards", len(routable.ShardInfo),
	)

	// Convert internal shard health to protobuf ShardStatus
	shardStatus := make(map[uint32]*pb.ShardStatus, len(routable.ShardInfo))
	for shardID, health := range routable.ShardInfo {
		shardStatus[shardID] = &pb.ShardStatus{
			IsHealthy: health.IsHealthy,
		}
	}

	// TODO: Rename WantShards to ReadyShards in protobuf (breaking change)
	// The field name "WantShards" is misleading - it actually represents
	// the number of shards ready for routing, not what Ring "wants".
	return &pb.ReplicaStatus{
		WantShards: uint32(routable.ReadyCount), //nolint:gosec // G115: replica count bounded
		Status:     shardStatus,
	}, nil
}

// ConsensusWantShards is called by the Ring consensus to report the desired number of shards.
// The consensus has agreed on how many shards the system should have.
func (h *ArbiterScalerHandler) ConsensusWantShards(ctx context.Context, req *pb.ConsensusWantShardsRequest) (*emptypb.Empty, error) {
	h.lggr.Infow("Consensus wants shards",
		"nShards", req.GetNShards(),
	)

	// Update the state with the consensus's desired shard count
	// This informs the Arbiter what the Ring consensus has agreed upon
	h.state.SetConsensusWantShards(int(req.GetNShards()))

	return &emptypb.Empty{}, nil
}
