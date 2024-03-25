package v3

import (
	"context"

	"google.golang.org/grpc"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_v3_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v3"
	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

var _ mercury_v3_types.ReportCodec = (*ReportCodecClient)(nil)

type ReportCodecClient struct {
	grpc mercury_v3_pb.ReportCodecClient
}

func NewReportCodecClient(cc grpc.ClientConnInterface) *ReportCodecClient {
	return &ReportCodecClient{grpc: mercury_v3_pb.NewReportCodecClient(cc)}
}

func (r *ReportCodecClient) BuildReport(ctx context.Context, fields mercury_v3_types.ReportFields) (ocr2plus_types.Report, error) {
	reply, err := r.grpc.BuildReport(ctx, &mercury_v3_pb.BuildReportRequest{
		ReportFields: pbReportFields(fields),
	})
	if err != nil {
		return ocr2plus_types.Report{}, err
	}
	return reply.Report, nil
}

func (r *ReportCodecClient) MaxReportLength(ctx context.Context, n int) (int, error) {
	reply, err := r.grpc.MaxReportLength(ctx, &mercury_v3_pb.MaxReportLengthRequest{})
	if err != nil {
		return 0, err
	}
	return int(reply.MaxReportLength), nil
}

func (r *ReportCodecClient) ObservationTimestampFromReport(ctx context.Context, report ocr2plus_types.Report) (uint32, error) {
	reply, err := r.grpc.ObservationTimestampFromReport(ctx, &mercury_v3_pb.ObservationTimestampFromReportRequest{
		Report: report,
	})
	if err != nil {
		return 0, err
	}
	return reply.Timestamp, nil
}

func pbReportFields(fields mercury_v3_types.ReportFields) *mercury_v3_pb.ReportFields {
	return &mercury_v3_pb.ReportFields{
		ValidFromTimestamp: fields.ValidFromTimestamp,
		Timestamp:          fields.Timestamp,
		NativeFee:          pb.NewBigIntFromInt(fields.NativeFee),
		LinkFee:            pb.NewBigIntFromInt(fields.LinkFee),
		ExpiresAt:          fields.ExpiresAt,
		BenchmarkPrice:     pb.NewBigIntFromInt(fields.BenchmarkPrice),
		Ask:                pb.NewBigIntFromInt(fields.Ask),
		Bid:                pb.NewBigIntFromInt(fields.Bid),
	}
}

var _ mercury_v3_pb.ReportCodecServer = (*ReportCodecServer)(nil)

type ReportCodecServer struct {
	mercury_v3_pb.UnimplementedReportCodecServer
	impl mercury_v3_types.ReportCodec
}

func NewReportCodecServer(impl mercury_v3_types.ReportCodec) *ReportCodecServer {
	return &ReportCodecServer{impl: impl}
}

func (r *ReportCodecServer) BuildReport(ctx context.Context, request *mercury_v3_pb.BuildReportRequest) (*mercury_v3_pb.BuildReportReply, error) {
	report, err := r.impl.BuildReport(ctx, reportFields(request.ReportFields))
	if err != nil {
		return nil, err
	}
	return &mercury_v3_pb.BuildReportReply{Report: report}, nil
}

func (r *ReportCodecServer) MaxReportLength(ctx context.Context, request *mercury_v3_pb.MaxReportLengthRequest) (*mercury_v3_pb.MaxReportLengthReply, error) {
	n, err := r.impl.MaxReportLength(ctx, int(request.NumOracles))
	if err != nil {
		return nil, err
	}
	return &mercury_v3_pb.MaxReportLengthReply{MaxReportLength: uint64(n)}, nil
}

func (r *ReportCodecServer) ObservationTimestampFromReport(ctx context.Context, request *mercury_v3_pb.ObservationTimestampFromReportRequest) (*mercury_v3_pb.ObservationTimestampFromReportReply, error) {
	timestamp, err := r.impl.ObservationTimestampFromReport(ctx, request.Report)
	if err != nil {
		return nil, err
	}
	return &mercury_v3_pb.ObservationTimestampFromReportReply{Timestamp: timestamp}, nil
}

func reportFields(fields *mercury_v3_pb.ReportFields) mercury_v3_types.ReportFields {
	return mercury_v3_types.ReportFields{
		ValidFromTimestamp: fields.ValidFromTimestamp,
		Timestamp:          fields.Timestamp,
		NativeFee:          fields.NativeFee.Int(),
		LinkFee:            fields.LinkFee.Int(),
		ExpiresAt:          fields.ExpiresAt,
		BenchmarkPrice:     fields.BenchmarkPrice.Int(),
		Ask:                fields.Ask.Int(),
		Bid:                fields.Bid.Int(),
	}
}
