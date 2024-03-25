package median

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

var _ median.DataSource = (*DataSourceClient)(nil)

type DataSourceClient struct {
	grpc pb.DataSourceClient
}

func NewDataSourceClient(cc grpc.ClientConnInterface) *DataSourceClient {
	return &DataSourceClient{grpc: pb.NewDataSourceClient(cc)}
}

func (d *DataSourceClient) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	reply, err := d.grpc.Observe(ctx, &pb.ObserveRequest{
		ReportTimestamp: pb.ReportTimestampToPb(timestamp),
	})
	if err != nil {
		return nil, err
	}
	return reply.Value.Int(), nil
}

var _ pb.DataSourceServer = (*DataSourceServer)(nil)

type DataSourceServer struct {
	pb.UnimplementedDataSourceServer

	impl median.DataSource
}

func NewDataSourceServer(impl median.DataSource) *DataSourceServer {
	return &DataSourceServer{impl: impl}
}

func (d *DataSourceServer) Observe(ctx context.Context, request *pb.ObserveRequest) (*pb.ObserveReply, error) {
	timestamp, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp)
	if err != nil {
		return nil, err
	}
	return &pb.ObserveReply{Value: pb.NewBigIntFromInt(val)}, nil
}
