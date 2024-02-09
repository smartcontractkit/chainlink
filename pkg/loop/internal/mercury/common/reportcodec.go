package common

import (
	"context"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"

	mercury_v1_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v1"
	mercury_v2_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v2"
	mercury_v3_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v3"
	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercury_v1_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v1"
	mercury_v2_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v2"
	mercury_v3_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v3"
	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

// The point of this is to translate between the well-versioned gRPC api in [pkg/loop/internal/pb/mercury] and the
// mercury provider [pkg/types/provider_mercury.go] which is not versioned.

// reportCodecV3Server implements mercury_pb.ReportCodecV3Server by wrapping [mercury_v3_internal.ReportCodecServer]
type reportCodecV3Server struct {
	mercury_pb.UnimplementedReportCodecV3Server

	impl *mercury_v3_internal.ReportCodecServer
}

var _ mercury_pb.ReportCodecV3Server = (*reportCodecV3Server)(nil)

// NewReportCodecV3Server returns a new instance of [mercury_pb.ReportCodecV3Server] which wraps [mercury_v3_internal.ReportCodecServer]
func NewReportCodecV3Server(s *grpc.Server, rc mercury_v3_types.ReportCodec) mercury_pb.ReportCodecV3Server {
	internalServer := mercury_v3_internal.NewReportCodecServer(rc)
	mercury_v3_pb.RegisterReportCodecServer(s, internalServer)
	return &reportCodecV3Server{impl: internalServer}
}

func (r *reportCodecV3Server) BuildReport(ctx context.Context, request *mercury_v3_pb.BuildReportRequest) (*mercury_v3_pb.BuildReportReply, error) {
	return r.impl.BuildReport(ctx, request)
}

func (r *reportCodecV3Server) MaxReportLength(ctx context.Context, request *mercury_v3_pb.MaxReportLengthRequest) (*mercury_v3_pb.MaxReportLengthReply, error) {
	return r.impl.MaxReportLength(ctx, request)
}

func (r *reportCodecV3Server) ObservationTimestampFromReport(ctx context.Context, request *mercury_v3_pb.ObservationTimestampFromReportRequest) (*mercury_v3_pb.ObservationTimestampFromReportReply, error) {
	return r.impl.ObservationTimestampFromReport(ctx, request)
}

var _ mercury_v3_types.ReportCodec = (*reportCodecV3Client)(nil)

type reportCodecV3Client struct {
	impl *mercury_v3_internal.ReportCodecClient
}

var _ mercury_v3_types.ReportCodec = (*reportCodecV3Client)(nil)

func NewReportCodecV3Client(impl *mercury_v3_internal.ReportCodecClient) mercury_v3_types.ReportCodec {
	return &reportCodecV3Client{impl: impl}
}

func (r *reportCodecV3Client) BuildReport(fields mercury_v3_types.ReportFields) (ocr2plus_types.Report, error) {
	return r.impl.BuildReport(fields)
}

func (r *reportCodecV3Client) MaxReportLength(n int) (int, error) {
	return r.impl.MaxReportLength(n)
}

func (r *reportCodecV3Client) ObservationTimestampFromReport(report ocr2plus_types.Report) (uint32, error) {
	return r.impl.ObservationTimestampFromReport(report)
}

// reportCodecV2Server implements mercury_pb.ReportCodecV2Server by wrapping [mercury_v2_internal.ReportCodecServer]
type reportCodecV2Server struct {
	mercury_pb.UnimplementedReportCodecV2Server

	impl *mercury_v2_internal.ReportCodecServer
}

var _ mercury_pb.ReportCodecV2Server = (*reportCodecV2Server)(nil)

// NewReportCodecV2Server returns a new instance of [mercury_pb.ReportCodecV2Server] which wraps [mercury_v2_internal.ReportCodecServer]
func NewReportCodecV2Server(s *grpc.Server, rc mercury_v2_types.ReportCodec) mercury_pb.ReportCodecV2Server {
	internalServer := mercury_v2_internal.NewReportCodecServer(rc)
	mercury_v2_pb.RegisterReportCodecServer(s, internalServer)
	return &reportCodecV2Server{impl: internalServer}
}

func (r *reportCodecV2Server) BuildReport(ctx context.Context, request *mercury_v2_pb.BuildReportRequest) (*mercury_v2_pb.BuildReportReply, error) {
	return r.impl.BuildReport(ctx, request)
}

func (r *reportCodecV2Server) MaxReportLength(ctx context.Context, request *mercury_v2_pb.MaxReportLengthRequest) (*mercury_v2_pb.MaxReportLengthReply, error) {
	return r.impl.MaxReportLength(ctx, request)
}

func (r *reportCodecV2Server) ObservationTimestampFromReport(ctx context.Context, request *mercury_v2_pb.ObservationTimestampFromReportRequest) (*mercury_v2_pb.ObservationTimestampFromReportReply, error) {
	return r.impl.ObservationTimestampFromReport(ctx, request)
}

type reportCodecV2Client struct {
	impl *mercury_v2_internal.ReportCodecClient
}

var _ mercury_v2_types.ReportCodec = (*reportCodecV2Client)(nil)

func NewReportCodecV2Client(impl *mercury_v2_internal.ReportCodecClient) mercury_v2_types.ReportCodec {
	return &reportCodecV2Client{impl: impl}
}

func (r *reportCodecV2Client) BuildReport(fields mercury_v2_types.ReportFields) (ocr2plus_types.Report, error) {
	return r.impl.BuildReport(fields)
}

func (r *reportCodecV2Client) MaxReportLength(n int) (int, error) {
	return r.impl.MaxReportLength(n)
}

func (r *reportCodecV2Client) ObservationTimestampFromReport(report ocr2plus_types.Report) (uint32, error) {
	return r.impl.ObservationTimestampFromReport(report)
}

// reportCodecV1Server implements mercury_pb.ReportCodecV1Server by wrapping [mercury_v1_internal.ReportCodecServer]
type reportCodecV1Server struct {
	mercury_pb.UnimplementedReportCodecV1Server

	impl *mercury_v1_internal.ReportCodecServer
}

var _ mercury_pb.ReportCodecV1Server = (*reportCodecV1Server)(nil)

// NewReportCodecV1Server returns a new instance of [mercury_pb.ReportCodecV1Server] which wraps [mercury_v1_internal.ReportCodecServer]
func NewReportCodecV1Server(s *grpc.Server, rc mercury_v1_types.ReportCodec) mercury_pb.ReportCodecV1Server {
	internalServer := mercury_v1_internal.NewReportCodecServer(rc)
	mercury_v1_pb.RegisterReportCodecServer(s, internalServer)
	return &reportCodecV1Server{impl: internalServer}
}

func (r *reportCodecV1Server) BuildReport(ctx context.Context, request *mercury_v1_pb.BuildReportRequest) (*mercury_v1_pb.BuildReportReply, error) {
	return r.impl.BuildReport(ctx, request)
}

func (r *reportCodecV1Server) MaxReportLength(ctx context.Context, request *mercury_v1_pb.MaxReportLengthRequest) (*mercury_v1_pb.MaxReportLengthReply, error) {
	return r.impl.MaxReportLength(ctx, request)
}

func (r *reportCodecV1Server) CurrentBlockNumFromReport(ctx context.Context, request *mercury_v1_pb.CurrentBlockNumFromReportRequest) (*mercury_v1_pb.CurrentBlockNumFromReportResponse, error) {
	return r.impl.CurrentBlockNumFromReport(ctx, request)
}

type reportCodecV1Client struct {
	impl *mercury_v1_internal.ReportCodecClient
}

var _ mercury_v1_types.ReportCodec = (*reportCodecV1Client)(nil)

func NewReportCodecV1Client(impl *mercury_v1_internal.ReportCodecClient) mercury_v1_types.ReportCodec {
	return &reportCodecV1Client{impl: impl}
}

func (r *reportCodecV1Client) BuildReport(fields mercury_v1_types.ReportFields) (ocr2plus_types.Report, error) {
	return r.impl.BuildReport(fields)
}

func (r *reportCodecV1Client) MaxReportLength(n int) (int, error) {
	return r.impl.MaxReportLength(n)
}

func (r *reportCodecV1Client) CurrentBlockNumFromReport(report ocr2plus_types.Report) (int64, error) {
	return r.impl.CurrentBlockNumFromReport(report)
}
