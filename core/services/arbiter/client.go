package arbiter

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"
)

// LocalArbiterScalerClient implements pb.ArbiterScalerClient by calling
// the ArbiterScalerServer directly without going over gRPC.
// This avoids network overhead when the Arbiter is running in the same process.
type LocalArbiterScalerClient struct {
	server pb.ArbiterScalerServer
	lggr   logger.Logger
}

var _ pb.ArbiterScalerClient = (*LocalArbiterScalerClient)(nil)

// NewLocalArbiterScalerClient creates a new LocalArbiterScalerClient.
func NewLocalArbiterScalerClient(server pb.ArbiterScalerServer, lggr logger.Logger) *LocalArbiterScalerClient {
	return &LocalArbiterScalerClient{
		server: server,
		lggr:   logger.Named(lggr, "LocalArbiterScalerClient"),
	}
}

// Status returns the current replica status by calling the server directly.
func (c *LocalArbiterScalerClient) Status(ctx context.Context, in *emptypb.Empty, _ ...grpc.CallOption) (*pb.ReplicaStatus, error) {
	return c.server.Status(ctx, in)
}

// ConsensusWantShards notifies the Arbiter about the desired shard count by calling the server directly.
func (c *LocalArbiterScalerClient) ConsensusWantShards(ctx context.Context, in *pb.ConsensusWantShardsRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.server.ConsensusWantShards(ctx, in)
}
