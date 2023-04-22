package loop

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
)

var _ median.DataSource = (*dataSourceClient)(nil)

type dataSourceClient struct {
	grpc pb.DataSourceClient
}

func newDataSourceClient(cc *grpc.ClientConn) *dataSourceClient {
	return &dataSourceClient{grpc: pb.NewDataSourceClient(cc)}
}

func (d *dataSourceClient) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	reply, err := d.grpc.Observe(ctx, &pb.ObserveRequest{ReportTimestamp: pbReportTimestamp(timestamp)})
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

func (d *dataSourceServer) Observe(ctx context.Context, request *pb.ObserveRequest) (*pb.ObserveReply, error) {
	timestamp, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp)
	if err != nil {
		return nil, err
	}
	return &pb.ObserveReply{Value: pb.NewBigIntFromInt(val)}, nil
}
