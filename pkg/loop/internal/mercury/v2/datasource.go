package v2

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_v2_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v2"
	mercury_common_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
)

var _ v2.DataSource = (*DataSourceClient)(nil)

type DataSourceClient struct {
	grpc mercury_v2_pb.DataSourceClient
}

func NewDataSourceClient(cc grpc.ClientConnInterface) *DataSourceClient {
	return &DataSourceClient{grpc: mercury_v2_pb.NewDataSourceClient(cc)}
}

func (d *DataSourceClient) Observe(ctx context.Context, timestamp ocr2plus_types.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (v2.Observation, error) {
	reply, err := d.grpc.Observe(ctx, &mercury_v2_pb.ObserveRequest{
		ReportTimestamp:            pb.ReportTimestampToPb(timestamp),
		FetchMaxFinalizedTimestamp: fetchMaxFinalizedTimestamp,
	})
	if err != nil {
		return v2.Observation{}, err
	}
	return observation(reply), nil
}

var _ mercury_v2_pb.DataSourceServer = (*DataSourceServer)(nil)

type DataSourceServer struct {
	mercury_v2_pb.UnimplementedDataSourceServer

	impl v2.DataSource
}

func NewDataSourceServer(impl v2.DataSource) *DataSourceServer {
	return &DataSourceServer{impl: impl}
}

func (d *DataSourceServer) Observe(ctx context.Context, request *mercury_v2_pb.ObserveRequest) (*mercury_v2_pb.ObserveResponse, error) {
	timestamp, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp, request.FetchMaxFinalizedTimestamp)
	if err != nil {
		return nil, err
	}
	return &mercury_v2_pb.ObserveResponse{Observation: pbObservation(val)}, nil
}

func observation(resp *mercury_v2_pb.ObserveResponse) v2.Observation {
	return v2.Observation{
		BenchmarkPrice:        mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.BenchmarkPrice.Int()},
		MaxFinalizedTimestamp: mercury_common_types.ObsResult[int64]{Val: resp.Observation.MaxFinalizedTimestamp},
		LinkPrice:             mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.LinkPrice.Int()},
		NativePrice:           mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.NativePrice.Int()},
	}
}

func pbObservation(obs v2.Observation) *mercury_v2_pb.Observation {
	return &mercury_v2_pb.Observation{
		BenchmarkPrice:        pb.NewBigIntFromInt(obs.BenchmarkPrice.Val),
		MaxFinalizedTimestamp: obs.MaxFinalizedTimestamp.Val,
		LinkPrice:             pb.NewBigIntFromInt(obs.LinkPrice.Val),
		NativePrice:           pb.NewBigIntFromInt(obs.NativePrice.Val),
	}
}
