package internal

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercurypb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
)

type mercuryPluginFactoryClient struct {
	*brokerExt
	*serviceClient
	grpc mercurypb.MercuryPluginFactoryClient
}

func newMercuryPluginFactoryClient(b *brokerExt, cc grpc.ClientConnInterface) *mercuryPluginFactoryClient {
	return &mercuryPluginFactoryClient{b.withName("MercuryPluginProviderClient"), newServiceClient(b, cc), mercurypb.NewMercuryPluginFactoryClient(cc)}
}

func (r *mercuryPluginFactoryClient) NewMercuryPlugin(config ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	ctx, cancel := r.stopCtx()
	defer cancel()

	response, err := r.grpc.NewMercuryPlugin(ctx, &mercurypb.NewMercuryPluginRequest{MercuryPluginConfig: &mercurypb.MercuryPluginConfig{
		ConfigDigest:           config.ConfigDigest[:],
		OracleID:               uint32(config.OracleID),
		N:                      uint32(config.N),
		F:                      uint32(config.F),
		OnchainConfig:          config.OnchainConfig,
		OffchainConfig:         config.OffchainConfig,
		EstimatedRoundInterval: int64(config.EstimatedRoundInterval),
		MaxDurationObservation: int64(config.MaxDurationObservation),
	}})
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}
	rpi := ocr3types.MercuryPluginInfo{
		Name: response.MercuryPluginInfo.Name,
		Limits: ocr3types.MercuryPluginLimits{
			MaxObservationLength: int(response.MercuryPluginInfo.MercuryPluginLimits.MaxObservationLength),
			MaxReportLength:      int(response.MercuryPluginInfo.MercuryPluginLimits.MaxReportLength),
		},
	}
	cc, err := r.brokerExt.dial(response.MercuryPluginID)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}
	return newMercuryPluginClient(r.brokerExt, cc), rpi, nil
}

var _ mercurypb.MercuryPluginFactoryServer = (*mercuryPluginFactoryServer)(nil)

type mercuryPluginFactoryServer struct {
	mercurypb.UnimplementedMercuryPluginFactoryServer

	*brokerExt

	impl ocr3types.MercuryPluginFactory
}

func newMercuryPluginFactoryServer(impl ocr3types.MercuryPluginFactory, b *brokerExt) *mercuryPluginFactoryServer {
	return &mercuryPluginFactoryServer{impl: impl, brokerExt: b.withName("MercuryPluginFactoryServer")}
}

func (r *mercuryPluginFactoryServer) NewMercuryPlugin(ctx context.Context, request *mercurypb.NewMercuryPluginRequest) (*mercurypb.NewMercuryPluginResponse, error) {
	cfg := ocr3types.MercuryPluginConfig{
		ConfigDigest:           libocr.ConfigDigest(request.MercuryPluginConfig.ConfigDigest),
		OracleID:               commontypes.OracleID(request.MercuryPluginConfig.OracleID),
		N:                      int(request.MercuryPluginConfig.N),
		F:                      int(request.MercuryPluginConfig.F),
		OnchainConfig:          request.MercuryPluginConfig.OnchainConfig,
		OffchainConfig:         request.MercuryPluginConfig.OffchainConfig,
		EstimatedRoundInterval: time.Duration(request.MercuryPluginConfig.EstimatedRoundInterval),
		MaxDurationObservation: time.Duration(request.MercuryPluginConfig.MaxDurationObservation),
	}
	if l := len(request.MercuryPluginConfig.ConfigDigest); l != 32 {
		return nil, pb.ErrConfigDigestLen(l)
	}
	copy(cfg.ConfigDigest[:], request.MercuryPluginConfig.ConfigDigest)

	rp, rpi, err := r.impl.NewMercuryPlugin(cfg)
	if err != nil {
		return nil, err
	}

	const mercuryname = "MercuryPlugin"
	id, _, err := r.serveNew(mercuryname, func(s *grpc.Server) {
		mercurypb.RegisterMercuryPluginServer(s, &mercuryPluginServer{impl: rp})
	}, resource{rp, mercuryname})
	if err != nil {
		return nil, err
	}

	return &mercurypb.NewMercuryPluginResponse{MercuryPluginID: id, MercuryPluginInfo: &mercurypb.MercuryPluginInfo{
		Name: rpi.Name,
		MercuryPluginLimits: &mercurypb.MercuryPluginLimits{
			MaxObservationLength: uint64(rpi.Limits.MaxObservationLength),
			MaxReportLength:      uint64(rpi.Limits.MaxReportLength),
		},
	}}, nil
}

var _ ocr3types.MercuryPlugin = (*mercuryPluginClient)(nil)

type mercuryPluginClient struct {
	*brokerExt
	grpc mercurypb.MercuryPluginClient
}

func newMercuryPluginClient(b *brokerExt, cc grpc.ClientConnInterface) *mercuryPluginClient {
	return &mercuryPluginClient{b.withName("MercuryPluginClient"), mercurypb.NewMercuryPluginClient(cc)}
}

func (r *mercuryPluginClient) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, previous libocr.Report) (libocr.Observation, error) {
	response, err := r.grpc.Observation(ctx, &mercurypb.ObservationRequest{
		ReportTimestamp: pb.ReportTimestampToPb(timestamp),
		PreviousReport:  previous,
	})
	if err != nil {
		return nil, err
	}
	return response.Observation, nil
}

// TODO: BCF-2887 plumb context through
func (r *mercuryPluginClient) Report(timestamp libocr.ReportTimestamp, previousReport libocr.Report, obs []libocr.AttributedObservation) (bool, libocr.Report, error) {
	response, err := r.grpc.Report(context.TODO(), &mercurypb.ReportRequest{
		ReportTimestamp: pb.ReportTimestampToPb(timestamp),
		PreviousReport:  previousReport,
		Observations:    mercurypbAttributedObservations(obs),
	})
	if err != nil {
		return false, nil, err
	}
	return response.ShouldReport, response.Report, nil
}

func (r *mercuryPluginClient) Close() error {
	ctx, cancel := r.stopCtx()
	defer cancel()

	_, err := r.grpc.Close(ctx, &emptypb.Empty{})
	return err
}

var _ mercurypb.MercuryPluginServer = (*mercuryPluginServer)(nil)

type mercuryPluginServer struct {
	mercurypb.UnimplementedMercuryPluginServer

	impl ocr3types.MercuryPlugin
}

func (r *mercuryPluginServer) Observation(ctx context.Context, request *mercurypb.ObservationRequest) (*mercurypb.ObservationResponse, error) {
	rts, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	o, err := r.impl.Observation(ctx, rts, request.PreviousReport)
	if err != nil {
		return nil, err
	}
	return &mercurypb.ObservationResponse{Observation: o}, nil
}

func (r *mercuryPluginServer) Report(ctx context.Context, request *mercurypb.ReportRequest) (*mercurypb.ReportResponse, error) {
	rts, err := pb.ReportTimestampFromPb(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	obs, err := mercuryattributedObservations(request.Observations)
	if err != nil {
		return nil, err
	}
	// TODO: BCF-2887 plumb context through
	should, report, err := r.impl.Report(rts, request.PreviousReport, obs)
	if err != nil {
		return nil, err
	}
	return &mercurypb.ReportResponse{
		ShouldReport: should,
		Report:       report,
	}, nil
}

func (r *mercuryPluginServer) Close(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, r.impl.Close()
}

func mercurypbAttributedObservations(obs []libocr.AttributedObservation) []*mercurypb.AttributedObservation {
	ret := make([]*mercurypb.AttributedObservation, len(obs))
	for i, o := range obs {
		ret[i] = &mercurypb.AttributedObservation{
			Observation: o.Observation,
			Observer:    uint32(o.Observer),
		}
	}
	return ret
}

func mercuryattributedObservations(obs []*mercurypb.AttributedObservation) ([]libocr.AttributedObservation, error) {
	ret := make([]libocr.AttributedObservation, len(obs))
	for i, o := range obs {
		ret[i] = libocr.AttributedObservation{
			Observation: o.Observation,
			Observer:    commontypes.OracleID(o.Observer),
		}
	}
	return ret, nil
}
