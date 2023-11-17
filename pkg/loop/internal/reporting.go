package internal

import (
	"context"
	"math"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

type reportingPluginFactoryClient struct {
	*brokerExt
	*serviceClient
	grpc pb.ReportingPluginFactoryClient
}

func newReportingPluginFactoryClient(b *brokerExt, cc grpc.ClientConnInterface) *reportingPluginFactoryClient {
	return &reportingPluginFactoryClient{b.withName("ReportingPluginProviderClient"), newServiceClient(b, cc), pb.NewReportingPluginFactoryClient(cc)}
}

func (r *reportingPluginFactoryClient) NewReportingPlugin(config libocr.ReportingPluginConfig) (libocr.ReportingPlugin, libocr.ReportingPluginInfo, error) {
	ctx, cancel := r.stopCtx()
	defer cancel()

	reply, err := r.grpc.NewReportingPlugin(ctx, &pb.NewReportingPluginRequest{ReportingPluginConfig: &pb.ReportingPluginConfig{
		ConfigDigest:                            config.ConfigDigest[:],
		OracleID:                                uint32(config.OracleID),
		N:                                       uint32(config.N),
		F:                                       uint32(config.F),
		OnchainConfig:                           config.OnchainConfig,
		OffchainConfig:                          config.OffchainConfig,
		EstimatedRoundInterval:                  int64(config.EstimatedRoundInterval),
		MaxDurationQuery:                        int64(config.MaxDurationQuery),
		MaxDurationObservation:                  int64(config.MaxDurationObservation),
		MaxDurationReport:                       int64(config.MaxDurationReport),
		MaxDurationShouldAcceptFinalizedReport:  int64(config.MaxDurationShouldAcceptFinalizedReport),
		MaxDurationShouldTransmitAcceptedReport: int64(config.MaxDurationShouldTransmitAcceptedReport),
	}})
	if err != nil {
		return nil, libocr.ReportingPluginInfo{}, err
	}
	rpi := libocr.ReportingPluginInfo{
		Name:          reply.ReportingPluginInfo.Name,
		UniqueReports: reply.ReportingPluginInfo.UniqueReports,
		Limits: libocr.ReportingPluginLimits{
			MaxQueryLength:       int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxQueryLength),
			MaxObservationLength: int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxObservationLength),
			MaxReportLength:      int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxReportLength),
		},
	}
	cc, err := r.brokerExt.dial(reply.ReportingPluginID)
	if err != nil {
		return nil, libocr.ReportingPluginInfo{}, err
	}
	return newReportingPluginClient(r.brokerExt, cc), rpi, nil
}

var _ pb.ReportingPluginFactoryServer = (*reportingPluginFactoryServer)(nil)

type reportingPluginFactoryServer struct {
	pb.UnimplementedReportingPluginFactoryServer

	*brokerExt

	impl libocr.ReportingPluginFactory
}

func newReportingPluginFactoryServer(impl libocr.ReportingPluginFactory, b *brokerExt) *reportingPluginFactoryServer {
	return &reportingPluginFactoryServer{impl: impl, brokerExt: b.withName("ReportingPluginFactoryServer")}
}

func (r *reportingPluginFactoryServer) NewReportingPlugin(ctx context.Context, request *pb.NewReportingPluginRequest) (*pb.NewReportingPluginReply, error) {
	cfg := libocr.ReportingPluginConfig{
		OracleID:                                commontypes.OracleID(request.ReportingPluginConfig.OracleID),
		N:                                       int(request.ReportingPluginConfig.N),
		F:                                       int(request.ReportingPluginConfig.F),
		OnchainConfig:                           request.ReportingPluginConfig.OnchainConfig,
		OffchainConfig:                          request.ReportingPluginConfig.OffchainConfig,
		EstimatedRoundInterval:                  time.Duration(request.ReportingPluginConfig.EstimatedRoundInterval),
		MaxDurationQuery:                        time.Duration(request.ReportingPluginConfig.MaxDurationQuery),
		MaxDurationObservation:                  time.Duration(request.ReportingPluginConfig.MaxDurationObservation),
		MaxDurationReport:                       time.Duration(request.ReportingPluginConfig.MaxDurationReport),
		MaxDurationShouldAcceptFinalizedReport:  time.Duration(request.ReportingPluginConfig.MaxDurationShouldAcceptFinalizedReport),
		MaxDurationShouldTransmitAcceptedReport: time.Duration(request.ReportingPluginConfig.MaxDurationShouldTransmitAcceptedReport),
	}
	if l := len(request.ReportingPluginConfig.ConfigDigest); l != 32 {
		return nil, ErrConfigDigestLen(l)
	}
	copy(cfg.ConfigDigest[:], request.ReportingPluginConfig.ConfigDigest)

	rp, rpi, err := r.impl.NewReportingPlugin(cfg)
	if err != nil {
		return nil, err
	}

	const name = "ReportingPlugin"
	id, _, err := r.serveNew(name, func(s *grpc.Server) {
		pb.RegisterReportingPluginServer(s, &reportingPluginServer{impl: rp})
	}, resource{rp, name})
	if err != nil {
		return nil, err
	}

	return &pb.NewReportingPluginReply{ReportingPluginID: id, ReportingPluginInfo: &pb.ReportingPluginInfo{
		Name:          rpi.Name,
		UniqueReports: rpi.UniqueReports,
		ReportingPluginLimits: &pb.ReportingPluginLimits{
			MaxQueryLength:       uint64(rpi.Limits.MaxQueryLength),
			MaxObservationLength: uint64(rpi.Limits.MaxObservationLength),
			MaxReportLength:      uint64(rpi.Limits.MaxReportLength),
		},
	}}, nil
}

var _ libocr.ReportingPlugin = (*reportingPluginClient)(nil)

type reportingPluginClient struct {
	*brokerExt
	grpc pb.ReportingPluginClient
}

func newReportingPluginClient(b *brokerExt, cc grpc.ClientConnInterface) *reportingPluginClient {
	return &reportingPluginClient{b.withName("ReportingPluginClient"), pb.NewReportingPluginClient(cc)}
}

func (r *reportingPluginClient) Query(ctx context.Context, timestamp libocr.ReportTimestamp) (libocr.Query, error) {
	reply, err := r.grpc.Query(ctx, &pb.QueryRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
	})
	if err != nil {
		return nil, err
	}
	return reply.Query, nil
}

func (r *reportingPluginClient) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, query libocr.Query) (libocr.Observation, error) {
	reply, err := r.grpc.Observation(ctx, &pb.ObservationRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Query:           query,
	})
	if err != nil {
		return nil, err
	}
	return reply.Observation, nil
}

func (r *reportingPluginClient) Report(ctx context.Context, timestamp libocr.ReportTimestamp, query libocr.Query, obs []libocr.AttributedObservation) (bool, libocr.Report, error) {
	reply, err := r.grpc.Report(ctx, &pb.ReportRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Query:           query,
		Observations:    pbAttributedObservations(obs),
	})
	if err != nil {
		return false, nil, err
	}
	return reply.ShouldReport, reply.Report, nil
}

func (r *reportingPluginClient) ShouldAcceptFinalizedReport(ctx context.Context, timestamp libocr.ReportTimestamp, report libocr.Report) (bool, error) {
	reply, err := r.grpc.ShouldAcceptFinalizedReport(ctx, &pb.ShouldAcceptFinalizedReportRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Report:          report,
	})
	if err != nil {
		return false, err
	}
	return reply.ShouldAccept, nil
}

func (r *reportingPluginClient) ShouldTransmitAcceptedReport(ctx context.Context, timestamp libocr.ReportTimestamp, report libocr.Report) (bool, error) {
	reply, err := r.grpc.ShouldTransmitAcceptedReport(ctx, &pb.ShouldTransmitAcceptedReportRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Report:          report,
	})
	if err != nil {
		return false, err
	}
	return reply.ShouldTransmit, nil
}

func (r *reportingPluginClient) Close() error {
	ctx, cancel := r.stopCtx()
	defer cancel()

	_, err := r.grpc.Close(ctx, &emptypb.Empty{})
	return err
}

var _ pb.ReportingPluginServer = (*reportingPluginServer)(nil)

type reportingPluginServer struct {
	pb.UnimplementedReportingPluginServer

	impl libocr.ReportingPlugin
}

func (r *reportingPluginServer) Query(ctx context.Context, request *pb.QueryRequest) (*pb.QueryReply, error) {
	rts, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	q, err := r.impl.Query(ctx, rts)
	if err != nil {
		return nil, err
	}
	return &pb.QueryReply{Query: q}, nil
}

func (r *reportingPluginServer) Observation(ctx context.Context, request *pb.ObservationRequest) (*pb.ObservationReply, error) {
	rts, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	o, err := r.impl.Observation(ctx, rts, request.Query)
	if err != nil {
		return nil, err
	}
	return &pb.ObservationReply{Observation: o}, nil
}

func (r *reportingPluginServer) Report(ctx context.Context, request *pb.ReportRequest) (*pb.ReportReply, error) {
	rts, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	obs, err := attributedObservations(request.Observations)
	if err != nil {
		return nil, err
	}
	should, report, err := r.impl.Report(ctx, rts, request.Query, obs)
	if err != nil {
		return nil, err
	}
	return &pb.ReportReply{
		ShouldReport: should,
		Report:       report,
	}, nil
}

func (r *reportingPluginServer) ShouldAcceptFinalizedReport(ctx context.Context, request *pb.ShouldAcceptFinalizedReportRequest) (*pb.ShouldAcceptFinalizedReportReply, error) {
	rts, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	should, err := r.impl.ShouldAcceptFinalizedReport(ctx, rts, request.Report)
	if err != nil {
		return nil, err
	}
	return &pb.ShouldAcceptFinalizedReportReply{ShouldAccept: should}, nil
}

func (r *reportingPluginServer) ShouldTransmitAcceptedReport(ctx context.Context, request *pb.ShouldTransmitAcceptedReportRequest) (*pb.ShouldTransmitAcceptedReportReply, error) {
	rts, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	should, err := r.impl.ShouldTransmitAcceptedReport(ctx, rts, request.Report)
	if err != nil {
		return nil, err
	}
	return &pb.ShouldTransmitAcceptedReportReply{ShouldTransmit: should}, nil
}

func (r *reportingPluginServer) Close(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, r.impl.Close()
}

func pbReportTimestamp(ts libocr.ReportTimestamp) *pb.ReportTimestamp {
	return &pb.ReportTimestamp{
		ConfigDigest: ts.ConfigDigest[:],
		Epoch:        ts.Epoch,
		Round:        uint32(ts.Round),
	}
}

func reportTimestamp(ts *pb.ReportTimestamp) (r libocr.ReportTimestamp, err error) {
	if l := len(ts.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(r.ConfigDigest[:], ts.ConfigDigest)
	r.Epoch = ts.Epoch
	if ts.Round > math.MaxUint8 {
		err = ErrUint8Bounds{Name: "Round", U: ts.Round}
		return
	}
	r.Round = uint8(ts.Round)
	return
}

func pbAttributedObservations(obs []libocr.AttributedObservation) (r []*pb.AttributedObservation) {
	for _, o := range obs {
		r = append(r, &pb.AttributedObservation{
			Observation: o.Observation,
			Observer:    uint32(o.Observer),
		})
	}
	return
}

func attributedObservations(pbos []*pb.AttributedObservation) (r []libocr.AttributedObservation, err error) {
	for _, pbo := range pbos {
		o, err := attributedObservation(pbo)
		if err != nil {
			return nil, err
		}
		r = append(r, o)
	}
	return
}

func attributedObservation(pbo *pb.AttributedObservation) (o libocr.AttributedObservation, err error) {
	o.Observation = pbo.Observation
	if pbo.Observer > math.MaxUint8 {
		err = ErrUint8Bounds{Name: "Observer", U: pbo.Observer}
		return
	}
	o.Observer = commontypes.OracleID(pbo.Observer)
	return
}
