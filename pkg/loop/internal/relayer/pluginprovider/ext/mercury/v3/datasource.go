package v3

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_v3_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v3"
	mercury_common_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

var _ v3.DataSource = (*DataSourceClient)(nil)

type DataSourceClient struct {
	grpc mercury_v3_pb.DataSourceClient
}

func NewDataSourceClient(cc grpc.ClientConnInterface) *DataSourceClient {
	return &DataSourceClient{grpc: mercury_v3_pb.NewDataSourceClient(cc)}
}

func (d *DataSourceClient) Observe(ctx context.Context, timestamp ocr2plus_types.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (v3.Observation, error) {
	reply, err := d.grpc.Observe(ctx, &mercury_v3_pb.ObserveRequest{
		ReportTimestamp:           pb.ReportTimestampToPb(timestamp),
		FetchMaxFinalizedBlockNum: fetchMaxFinalizedTimestamp,
	})
	if err != nil {
		return v3.Observation{}, err
	}
	return observation(reply), nil
}

var _ mercury_v3_pb.DataSourceServer = (*DataSourceServer)(nil)

type DataSourceServer struct {
	mercury_v3_pb.UnimplementedDataSourceServer

	impl v3.DataSource
}

func NewDataSourceServer(impl v3.DataSource) *DataSourceServer {
	return &DataSourceServer{impl: impl}
}

func (d *DataSourceServer) Observe(ctx context.Context, request *mercury_v3_pb.ObserveRequest) (*mercury_v3_pb.ObserveResponse, error) {
	timestamp, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp, request.FetchMaxFinalizedBlockNum)
	if err != nil {
		return nil, err
	}
	return &mercury_v3_pb.ObserveResponse{Observation: pbObservation(val)}, nil
}

func observation(resp *mercury_v3_pb.ObserveResponse) v3.Observation {
	return v3.Observation{
		BenchmarkPrice:        mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.BenchmarkPrice.Int()},
		Bid:                   mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.Bid.Int()},
		Ask:                   mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.Ask.Int()},
		MaxFinalizedTimestamp: mercury_common_types.ObsResult[int64]{Val: resp.Observation.MaxFinalizedTimestamp},
		LinkPrice:             mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.LinkPrice.Int()},
		NativePrice:           mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.NativePrice.Int()},
	}
}

func pbObservation(obs v3.Observation) *mercury_v3_pb.Observation {
	return &mercury_v3_pb.Observation{
		BenchmarkPrice:        pb.NewBigIntFromInt(obs.BenchmarkPrice.Val),
		Bid:                   pb.NewBigIntFromInt(obs.Bid.Val),
		Ask:                   pb.NewBigIntFromInt(obs.Ask.Val),
		MaxFinalizedTimestamp: obs.MaxFinalizedTimestamp.Val,
		LinkPrice:             pb.NewBigIntFromInt(obs.LinkPrice.Val),
		NativePrice:           pb.NewBigIntFromInt(obs.NativePrice.Val),
	}
}
