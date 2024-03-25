package v1

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_v1_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v1"
	mercury_common_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
)

var _ v1.DataSource = (*DataSourceClient)(nil)

type DataSourceClient struct {
	grpc mercury_v1_pb.DataSourceClient
}

func NewDataSourceClient(cc grpc.ClientConnInterface) *DataSourceClient {
	return &DataSourceClient{grpc: mercury_v1_pb.NewDataSourceClient(cc)}
}

func (d *DataSourceClient) Observe(ctx context.Context, timestamp ocr2plus_types.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (v1.Observation, error) {
	reply, err := d.grpc.Observe(ctx, &mercury_v1_pb.ObserveRequest{
		ReportTimestamp: pb.ReportTimestampToPb(timestamp),
	})
	if err != nil {
		return v1.Observation{}, err
	}
	return observation(reply), nil
}

var _ mercury_v1_pb.DataSourceServer = (*DataSourceServer)(nil)

type DataSourceServer struct {
	mercury_v1_pb.UnimplementedDataSourceServer

	impl v1.DataSource
}

func NewDataSourceServer(impl v1.DataSource) *DataSourceServer {
	return &DataSourceServer{impl: impl}
}

func (d *DataSourceServer) Observe(ctx context.Context, request *mercury_v1_pb.ObserveRequest) (*mercury_v1_pb.ObserveResponse, error) {
	timestamp, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp, request.FetchMaxFinalizedBlockNum)
	if err != nil {
		return nil, err
	}
	return &mercury_v1_pb.ObserveResponse{Observation: pbObservation(val)}, nil
}

func observation(resp *mercury_v1_pb.ObserveResponse) v1.Observation {
	return v1.Observation{
		BenchmarkPrice:          mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.BenchmarkPrice.Int()},
		Bid:                     mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.Bid.Int()},
		Ask:                     mercury_common_types.ObsResult[*big.Int]{Val: resp.Observation.Ask.Int()},
		CurrentBlockNum:         mercury_common_types.ObsResult[int64]{Val: resp.Observation.CurrentBlockNum},
		CurrentBlockHash:        mercury_common_types.ObsResult[[]byte]{Val: resp.Observation.CurrentBlockHash},
		CurrentBlockTimestamp:   mercury_common_types.ObsResult[uint64]{Val: resp.Observation.CurrentBlockTimestamp},
		LatestBlocks:            blocks(resp.Observation.LatestBlocks),
		MaxFinalizedBlockNumber: mercury_common_types.ObsResult[int64]{Val: resp.Observation.MaxFinalizedBlockNumber},
	}
}

func blocks(blocks []*mercury_v1_pb.Block) []v1.Block {
	var ret []v1.Block
	for _, b := range blocks {
		ret = append(ret, block(b))
	}
	return ret
}

func block(pb *mercury_v1_pb.Block) v1.Block {
	return v1.Block{
		Num:  pb.Number,
		Hash: string(pb.Hash),
		Ts:   pb.Timestamp,
	}
}

func pbObservation(obs v1.Observation) *mercury_v1_pb.Observation {
	return &mercury_v1_pb.Observation{
		BenchmarkPrice:          pb.NewBigIntFromInt(obs.BenchmarkPrice.Val),
		Bid:                     pb.NewBigIntFromInt(obs.Bid.Val),
		Ask:                     pb.NewBigIntFromInt(obs.Ask.Val),
		CurrentBlockNum:         obs.CurrentBlockNum.Val,
		CurrentBlockHash:        obs.CurrentBlockHash.Val,
		CurrentBlockTimestamp:   obs.CurrentBlockTimestamp.Val,
		LatestBlocks:            pbBlocks(obs.LatestBlocks),
		MaxFinalizedBlockNumber: obs.MaxFinalizedBlockNumber.Val,
	}
}

func pbBlocks(blocks []v1.Block) []*mercury_v1_pb.Block {
	var ret []*mercury_v1_pb.Block
	for _, b := range blocks {
		ret = append(ret, pbBlock(b))
	}
	return ret
}

func pbBlock(b v1.Block) *mercury_v1_pb.Block {
	return &mercury_v1_pb.Block{
		Number:    b.Num,
		Hash:      []byte(b.Hash),
		Timestamp: b.Ts,
	}
}
