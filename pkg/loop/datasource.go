package loop

import (
	"context"
	"math/big"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
)

var _ median.DataSource = (*dataSourceClient)(nil)

type dataSourceClient struct {
	grpc pb.DataSourceClient
}

func newDataSourceClient(cc *grpc.ClientConn) *dataSourceClient {
	return &dataSourceClient{grpc: pb.NewDataSourceClient(cc)}
}

func (d *dataSourceClient) Observe(ctx context.Context) (*big.Int, error) {
	reply, err := d.grpc.Observe(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return reply.Value.Int(), nil
}

var _ pb.DataSourceServer = (*dataSourceServer)(nil)

type dataSourceServer struct {
	pb.UnimplementedDataSourceServer

	impl median.DataSource
}

func (d *dataSourceServer) Observe(ctx context.Context, _ *emptypb.Empty) (*pb.ObserveReply, error) {
	val, err := d.impl.Observe(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ObserveReply{Value: pb.NewBigIntFromInt(val)}, nil
}
