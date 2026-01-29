package arbiter

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

// RingArbiterClient implements ringpb.ArbiterScalerClient by calling
// the ArbiterScalerServer directly without going over gRPC.
// This is used by Ring OCR to communicate with the Arbiter in-process.
type RingArbiterClient struct {
	server ringpb.ArbiterScalerServer
	lggr   logger.Logger
}

var _ ringpb.ArbiterScalerClient = (*RingArbiterClient)(nil)

// NewRingArbiterClient creates a new RingArbiterClient.
func NewRingArbiterClient(server ringpb.ArbiterScalerServer, lggr logger.Logger) *RingArbiterClient {
	return &RingArbiterClient{
		server: server,
		lggr:   logger.Named(lggr, "RingArbiterClient"),
	}
}

// Status returns the current replica status by calling the server directly.
func (c *RingArbiterClient) Status(ctx context.Context, in *emptypb.Empty, _ ...grpc.CallOption) (*ringpb.ReplicaStatus, error) {
	return c.server.Status(ctx, in)
}

// ConsensusWantShards notifies the Arbiter about the desired shard count by calling the server directly.
func (c *RingArbiterClient) ConsensusWantShards(ctx context.Context, in *ringpb.ConsensusWantShardsRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.server.ConsensusWantShards(ctx, in)
}
