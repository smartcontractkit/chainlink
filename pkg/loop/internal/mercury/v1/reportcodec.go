package v1

import (
	"context"

	"google.golang.org/grpc"

	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_v1_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v1"
)

var _ mercury_v1_types.ReportCodec = (*ReportCodecClient)(nil)

type ReportCodecClient struct {
	grpc mercury_v1_pb.ReportCodecClient
}

func NewReportCodecClient(cc grpc.ClientConnInterface) *ReportCodecClient {
	return &ReportCodecClient{grpc: mercury_v1_pb.NewReportCodecClient(cc)}
}

func (r *ReportCodecClient) BuildReport(ctx context.Context, fields mercury_v1_types.ReportFields) (ocr2plus_types.Report, error) {
	reply, err := r.grpc.BuildReport(ctx, &mercury_v1_pb.BuildReportRequest{
		ReportFields: pbReportFields(fields),
	})
	if err != nil {
		return ocr2plus_types.Report{}, err
	}
	return reply.Report, nil
}

func (r *ReportCodecClient) MaxReportLength(ctx context.Context, n int) (int, error) {
	reply, err := r.grpc.MaxReportLength(ctx, &mercury_v1_pb.MaxReportLengthRequest{})
	if err != nil {
		return 0, err
	}
	return int(reply.MaxReportLength), nil
}

func (r *ReportCodecClient) CurrentBlockNumFromReport(ctx context.Context, report ocr2plus_types.Report) (int64, error) {
	reply, err := r.grpc.CurrentBlockNumFromReport(ctx, &mercury_v1_pb.CurrentBlockNumFromReportRequest{
		Report: report,
	})
	if err != nil {
		return 0, err
	}
	return reply.CurrentBlockNum, nil
}

func pbReportFields(fields mercury_v1_types.ReportFields) *mercury_v1_pb.ReportFields {
	return &mercury_v1_pb.ReportFields{
		Timestamp:             fields.Timestamp,
		BenchmarkPrice:        pb.NewBigIntFromInt(fields.BenchmarkPrice),
		Ask:                   pb.NewBigIntFromInt(fields.Ask),
		Bid:                   pb.NewBigIntFromInt(fields.Bid),
		CurrentBlockNum:       fields.CurrentBlockNum,
		CurrentBlockHash:      fields.CurrentBlockHash,
		ValidFromBlockNum:     fields.ValidFromBlockNum,
		CurrentBlockTimestamp: fields.CurrentBlockTimestamp,
	}
}

var _ mercury_v1_pb.ReportCodecServer = (*ReportCodecServer)(nil)

type ReportCodecServer struct {
	mercury_v1_pb.UnimplementedReportCodecServer
	impl mercury_v1_types.ReportCodec
}

func NewReportCodecServer(impl mercury_v1_types.ReportCodec) *ReportCodecServer {
	return &ReportCodecServer{impl: impl}
}

func (r *ReportCodecServer) BuildReport(ctx context.Context, request *mercury_v1_pb.BuildReportRequest) (*mercury_v1_pb.BuildReportReply, error) {
	report, err := r.impl.BuildReport(ctx, reportFields(request.ReportFields))
	if err != nil {
		return nil, err
	}
	return &mercury_v1_pb.BuildReportReply{Report: report}, nil
}

func (r *ReportCodecServer) MaxReportLength(ctx context.Context, request *mercury_v1_pb.MaxReportLengthRequest) (*mercury_v1_pb.MaxReportLengthReply, error) {
	n, err := r.impl.MaxReportLength(ctx, int(request.NumOracles))
	if err != nil {
		return nil, err
	}
	return &mercury_v1_pb.MaxReportLengthReply{MaxReportLength: uint64(n)}, nil
}

func (r *ReportCodecServer) CurrentBlockNumFromReport(ctx context.Context, request *mercury_v1_pb.CurrentBlockNumFromReportRequest) (*mercury_v1_pb.CurrentBlockNumFromReportResponse, error) {
	n, err := r.impl.CurrentBlockNumFromReport(ctx, request.Report)
	if err != nil {
		return nil, err
	}
	return &mercury_v1_pb.CurrentBlockNumFromReportResponse{CurrentBlockNum: n}, nil
}

func reportFields(fields *mercury_v1_pb.ReportFields) mercury_v1_types.ReportFields {
	return mercury_v1_types.ReportFields{
		Timestamp:             fields.Timestamp,
		BenchmarkPrice:        fields.BenchmarkPrice.Int(),
		Ask:                   fields.Ask.Int(),
		Bid:                   fields.Bid.Int(),
		CurrentBlockNum:       fields.CurrentBlockNum,
		CurrentBlockHash:      fields.CurrentBlockHash,
		ValidFromBlockNum:     fields.ValidFromBlockNum,
		CurrentBlockTimestamp: fields.CurrentBlockTimestamp,
	}
}
