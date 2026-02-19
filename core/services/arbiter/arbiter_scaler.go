package arbiter

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

// RingArbiterHandler implements the ArbiterScalerServer interface from chainlink-common.
// This is the inbound handler for Ring OCR → Arbiter communication.
// It provides:
//   - Status(): Returns routable shard count and per-shard health for Ring OCR routing decisions
//   - ConsensusWantShards(): Receives the Ring consensus decision about desired shard count
type RingArbiterHandler struct {
	ringpb.UnimplementedArbiterScalerServer
	state *State
	lggr  logger.Logger
}

// NewRingArbiterHandler creates a new RingArbiterHandler.
func NewRingArbiterHandler(state *State, lggr logger.Logger) *RingArbiterHandler {
	return &RingArbiterHandler{
		state: state,
		lggr:  logger.Named(lggr, "RingArbiterHandler"),
	}
}

// Status returns the current replica status for Ring OCR routing.
// Returns only READY shards count and per-shard health status.
// This is called by the Ring plugin to determine which shards can receive traffic.
func (h *RingArbiterHandler) Status(ctx context.Context, _ *emptypb.Empty) (*ringpb.ReplicaStatus, error) {
	routable := h.state.GetRoutableShards()

	h.lggr.Debugw("Status requested",
		"readyShards", routable.ReadyCount,
		"totalShards", len(routable.ShardInfo),
	)

	// Convert internal shard health to protobuf ShardStatus
	shardStatus := make(map[uint32]*ringpb.ShardStatus, len(routable.ShardInfo))
	for shardID, health := range routable.ShardInfo {
		shardStatus[shardID] = &ringpb.ShardStatus{
			IsHealthy: health.IsHealthy,
		}
	}

	// TODO: Rename WantShards to ReadyShards in protobuf (breaking change)
	// The field name "WantShards" is misleading - it actually represents
	// the number of shards ready for routing, not what Ring "wants".
	return &ringpb.ReplicaStatus{
		WantShards: uint32(routable.ReadyCount), //nolint:gosec // G115: replica count bounded
		Status:     shardStatus,
	}, nil
}

// ConsensusWantShards is called by the Ring consensus to report the desired number of shards.
// The consensus has agreed on how many shards the system should have.
func (h *RingArbiterHandler) ConsensusWantShards(ctx context.Context, req *ringpb.ConsensusWantShardsRequest) (*emptypb.Empty, error) {
	nShards := req.GetNShards()

	if nShards == 0 {
		h.lggr.Warnw("Consensus reported 0 shards, this may indicate a problem",
			"nShards", nShards,
		)
	}

	h.lggr.Infow("Consensus wants shards",
		"nShards", nShards,
	)

	// Update the state with the consensus's desired shard count
	// This informs the Arbiter what the Ring consensus has agreed upon
	h.state.SetConsensusWantShards(int(nShards))

	return &emptypb.Empty{}, nil
}
