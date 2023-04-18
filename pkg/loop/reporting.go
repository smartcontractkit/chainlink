package loop

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
)

type reportingPluginFactoryClient struct {
	*lggrBroker
	grpc pb.ReportingPluginFactoryClient
}

func newReportingPluginFactoryClient(lb *lggrBroker, cc *grpc.ClientConn) *reportingPluginFactoryClient {
	return &reportingPluginFactoryClient{lb.named("ReportingPluginProviderClient"), pb.NewReportingPluginFactoryClient(cc)}
}

func (r *reportingPluginFactoryClient) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	reply, err := r.grpc.NewReportingPlugin(context.TODO(), &pb.NewReportingPluginRequest{ReportingPluginConfig: &pb.ReportingPluginConfig{
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
		return nil, types.ReportingPluginInfo{}, err
	}
	rpConn, err := r.broker.Dial(reply.ReportingPluginID)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, ErrConnDial{Name: "ReportingPlugin", ID: reply.ReportingPluginID, Err: err}
	}
	rpi := types.ReportingPluginInfo{
		Name:          reply.ReportingPluginInfo.Name,
		UniqueReports: reply.ReportingPluginInfo.UniqueReports,
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxQueryLength),
			MaxObservationLength: int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxObservationLength),
			MaxReportLength:      int(reply.ReportingPluginInfo.ReportingPluginLimits.MaxReportLength),
		},
	}
	return newReportingPluginClient(rpConn), rpi, nil
}

var _ pb.ReportingPluginFactoryServer = (*reportingPluginFactoryServer)(nil)

type reportingPluginFactoryServer struct {
	pb.UnimplementedReportingPluginFactoryServer

	*lggrBroker

	impl types.ReportingPluginFactory
}

func newReportingPluginFactoryServer(impl types.ReportingPluginFactory, lb *lggrBroker) *reportingPluginFactoryServer {
	return &reportingPluginFactoryServer{impl: impl, lggrBroker: lb.named("ReportingPluginFactoryServer")}
}

func (r *reportingPluginFactoryServer) NewReportingPlugin(ctx context.Context, request *pb.NewReportingPluginRequest) (*pb.NewReportingPluginReply, error) {
	cfg := types.ReportingPluginConfig{
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
	s := grpc.NewServer()
	pb.RegisterReportingPluginServer(s, &reportingPluginServer{impl: rp, stop: func() {
		time.AfterFunc(time.Second, s.GracefulStop)
	}})
	const name = "ReportingPlugin"
	id, err := r.serve(s, name, resource{rp, name})
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

var _ types.ReportingPlugin = (*reportingPluginClient)(nil)

type reportingPluginClient struct {
	grpc  pb.ReportingPluginClient
	close func() error
}

func newReportingPluginClient(cc *grpc.ClientConn) *reportingPluginClient {
	return &reportingPluginClient{pb.NewReportingPluginClient(cc), cc.Close}
}

func (r *reportingPluginClient) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	reply, err := r.grpc.Query(ctx, &pb.QueryRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
	})
	if err != nil {
		return nil, err
	}
	return reply.Query, nil
}

func (r *reportingPluginClient) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	reply, err := r.grpc.Observation(ctx, &pb.ObservationRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Query:           query,
	})
	if err != nil {
		return nil, err
	}
	return reply.Observation, nil
}

func (r *reportingPluginClient) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
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

func (r *reportingPluginClient) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	reply, err := r.grpc.ShouldAcceptFinalizedReport(ctx, &pb.ShouldAcceptFinalizedReportRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
		Report:          report,
	})
	if err != nil {
		return false, err
	}
	return reply.ShouldAccept, nil
}

func (r *reportingPluginClient) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
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
	_, err := r.grpc.Close(context.TODO(), &emptypb.Empty{})
	err = errors.Join(err, r.close())
	return err
}

var _ pb.ReportingPluginServer = (*reportingPluginServer)(nil)

type reportingPluginServer struct {
	pb.UnimplementedReportingPluginServer

	impl types.ReportingPlugin
	stop func()
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
	r.stop()
	return &emptypb.Empty{}, r.impl.Close()
}

func pbReportTimestamp(ts types.ReportTimestamp) *pb.ReportTimestamp {
	return &pb.ReportTimestamp{
		ConfigDigest: ts.ConfigDigest[:],
		Epoch:        ts.Epoch,
		Round:        uint32(ts.Round),
	}
}

func reportTimestamp(ts *pb.ReportTimestamp) (r types.ReportTimestamp, err error) {
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

func pbAttributedObservations(obs []types.AttributedObservation) (r []*pb.AttributedObservation) {
	for _, o := range obs {
		r = append(r, &pb.AttributedObservation{
			Observation: o.Observation,
			Observer:    uint32(o.Observer),
		})
	}
	return
}

func attributedObservations(pbos []*pb.AttributedObservation) (r []types.AttributedObservation, err error) {
	for _, pbo := range pbos {
		o, err := attributedObservation(pbo)
		if err != nil {
			return nil, err
		}
		r = append(r, o)
	}
	return
}

func attributedObservation(pbo *pb.AttributedObservation) (o types.AttributedObservation, err error) {
	o.Observation = pbo.Observation
	if pbo.Observer > math.MaxUint8 {
		err = ErrUint8Bounds{Name: "Observer", U: pbo.Observer}
		return
	}
	o.Observer = commontypes.OracleID(pbo.Observer)
	return
}
