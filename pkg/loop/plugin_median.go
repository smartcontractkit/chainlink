package loop

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// PluginMedianName is the name for [PluginMedian]/[NewGRPCPluginMedian].
const PluginMedianName = "median"

type PluginMedian interface {
	NewMedianPluginFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog ErrorLog) (libocr.ReportingPluginFactory, error)
}

type ErrorLog interface {
	SaveError(ctx context.Context, msg string) error
}

func PluginMedianHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  0,
		MagicCookieKey:   "magic-key-median-todo",
		MagicCookieValue: "magic-value-median-todo",
	}
}

func PluginMedianClientConfig(lggr logger.Logger) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig: PluginMedianHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			PluginMedianName: NewGRPCPluginMedian(nil, lggr),
		},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}

type grpcPluginMedian struct {
	plugin.NetRPCUnsupportedPlugin

	impl PluginMedian
	lggr logger.Logger
}

func NewGRPCPluginMedian(mp PluginMedian, lggr logger.Logger) plugin.Plugin {
	return &grpcPluginMedian{impl: mp, lggr: lggr}
}

func (p *grpcPluginMedian) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	pb.RegisterPluginMedianServer(server, newPluginMedianServer(&lggrBroker{p.lggr, broker}, p.impl))
	return nil
}

// GRPCClient implements [plugin.GRPCPlugin] and returns a [PluginRelayer].
func (p *grpcPluginMedian) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return newMedianPluginClient(p.lggr, broker, conn), nil
}

var _ PluginMedian = (*pluginMedianClient)(nil)

type pluginMedianClient struct {
	*lggrBroker

	grpc pb.PluginMedianClient
}

func newMedianPluginClient(lggr logger.Logger, broker *plugin.GRPCBroker, conn *grpc.ClientConn) *pluginMedianClient {
	lggr = logger.Named(lggr, "MedianReportingClient")
	return &pluginMedianClient{lggrBroker: &lggrBroker{lggr, broker}, grpc: pb.NewPluginMedianClient(conn)}
}

func (m *pluginMedianClient) NewMedianPluginFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog ErrorLog) (libocr.ReportingPluginFactory, error) {
	dsSrv := grpc.NewServer()
	pb.RegisterDataSourceServer(dsSrv, &dataSourceServer{impl: dataSource})
	dataSourceID, err := m.serve(dsSrv, "DataSource")
	if err != nil {
		return nil, err
	}
	juelsSrv := grpc.NewServer()
	pb.RegisterDataSourceServer(juelsSrv, &dataSourceServer{impl: juelsPerFeeCoin})
	juelsPerFeeCoinDataSourceID, err := m.serve(juelsSrv, "JuelsPerFeeCoinDataSource")
	if err != nil {
		dsSrv.Stop()
		return nil, err
	}

	providerSrv := grpc.NewServer()
	pb.RegisterServiceServer(providerSrv, &serviceServer{srv: provider, stop: func() {
		time.AfterFunc(time.Second, providerSrv.GracefulStop)
	}})
	pb.RegisterOffchainConfigDigesterServer(providerSrv, &offchainConfigDigesterServer{impl: provider.OffchainConfigDigester()})
	pb.RegisterContractConfigTrackerServer(providerSrv, &contractConfigTrackerServer{impl: provider.ContractConfigTracker()})
	pb.RegisterContractTransmitterServer(providerSrv, &contractTransmitterServer{impl: provider.ContractTransmitter()})
	pb.RegisterReportCodecServer(providerSrv, &reportCodecServer{impl: provider.ReportCodec()})
	pb.RegisterMedianContractServer(providerSrv, &medianContractServer{impl: provider.MedianContract()})
	pb.RegisterOnchainConfigCodecServer(providerSrv, &onchainConfigCodecServer{impl: provider.OnchainConfigCodec()})
	providerID, err := m.serve(providerSrv, "MedianProvider")
	if err != nil {
		dsSrv.Stop()
		juelsSrv.Stop()
		return nil, err
	}

	errorLogSrv := grpc.NewServer()
	pb.RegisterErrorLogServer(errorLogSrv, &errorLogServer{impl: errorLog})
	errorLogID, err := m.serve(errorLogSrv, "ErrorLog")
	if err != nil {
		dsSrv.Stop()
		juelsSrv.Stop()
		providerSrv.Stop()
		return nil, err
	}

	reply, err := m.grpc.NewMedianPluginFactory(ctx, &pb.NewMedianPluginFactoryRequest{
		MedianProviderID:            providerID,
		DataSourceID:                dataSourceID,
		JuelsPerFeeCoinDataSourceID: juelsPerFeeCoinDataSourceID,
		ErrorLogID:                  errorLogID,
	})
	if err != nil {
		dsSrv.Stop()
		juelsSrv.Stop()
		providerSrv.Stop()
		errorLogSrv.Stop()
		return nil, err
	}
	id := reply.ReportingPluginFactoryID
	cc, err := m.broker.Dial(id)
	if err != nil {
		dsSrv.Stop()
		juelsSrv.Stop()
		providerSrv.Stop()
		errorLogSrv.Stop()
		return nil, ErrConnDial{Name: "MedianPluginFactory", ID: id, Err: err}
	}
	//TODO client should close everything
	return newReportingPluginFactoryClient(m.lggrBroker, cc), nil
}

var _ pb.PluginMedianServer = (*pluginMedianServer)(nil)

type pluginMedianServer struct {
	pb.UnimplementedPluginMedianServer

	*lggrBroker
	impl PluginMedian
}

func newPluginMedianServer(lb *lggrBroker, mp PluginMedian) *pluginMedianServer {
	return &pluginMedianServer{lggrBroker: lb.named("PluginMedian"), impl: mp}
}

func (m *pluginMedianServer) NewMedianPluginFactory(ctx context.Context, request *pb.NewMedianPluginFactoryRequest) (*pb.NewMedianPluginFactoryReply, error) {
	dsConn, err := m.broker.Dial(request.DataSourceID)
	if err != nil {
		return nil, ErrConnDial{Name: "DataSource", ID: request.DataSourceID, Err: err}
	}
	dsRes := resource{dsConn, "DataSource"}
	dataSource := newDataSourceClient(dsConn)

	juelsConn, err := m.broker.Dial(request.JuelsPerFeeCoinDataSourceID)
	if err != nil {
		m.closeAll(dsRes)
		return nil, ErrConnDial{Name: "JuelsPerFeeCoinDataSource", ID: request.JuelsPerFeeCoinDataSourceID, Err: err}
	}
	juelsRes := resource{juelsConn, "JuelsPerFeeCoinDataSource"}
	juelsPerFeeCoin := newDataSourceClient(juelsConn)

	providerConn, err := m.broker.Dial(request.MedianProviderID)
	if err != nil {
		m.closeAll(dsRes, juelsRes)
		return nil, ErrConnDial{Name: "MedianProvider", ID: request.MedianProviderID, Err: err}
	}
	providerRes := resource{providerConn, "MedianProvider"}
	provider := newMedianProviderClient(m.lggrBroker, providerConn)

	errorLogConn, err := m.broker.Dial(request.ErrorLogID)
	if err != nil {
		m.closeAll(dsRes, juelsRes, providerRes)
		return nil, ErrConnDial{Name: "ErrorLog", ID: request.ErrorLogID, Err: err}
	}
	errorLogRes := resource{errorLogConn, "ErrorLog"}
	errorLog := newErrorLogClient(errorLogConn)

	factory, err := m.impl.NewMedianPluginFactory(ctx, provider, dataSource, juelsPerFeeCoin, errorLog)
	if err != nil {
		m.closeAll(dsRes, juelsRes, providerRes, errorLogRes)
		return nil, err
	}

	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &serviceServer{srv: provider, stop: func() {
		time.AfterFunc(time.Second, s.GracefulStop)
		m.closeAll(dsRes, juelsRes, providerRes, errorLogRes)
	}})
	pb.RegisterReportingPluginFactoryServer(s, newReportingPluginFactoryServer(factory, m.lggrBroker))
	id, err := m.serve(s, "ReportingPluginProvider")
	if err != nil {
		m.closeAll(dsRes, juelsRes, providerRes, errorLogRes)
		return nil, err
	}
	return &pb.NewMedianPluginFactoryReply{ReportingPluginFactoryID: id}, nil
}

type medianProviderClient struct {
	*configProviderClient
	contractTransmitter libocr.ContractTransmitter
	reportCodec         median.ReportCodec
	medianContract      median.MedianContract
	onchainConfigCodec  median.OnchainConfigCodec
}

func newMedianProviderClient(lb *lggrBroker, cc *grpc.ClientConn) *medianProviderClient {
	m := &medianProviderClient{configProviderClient: newConfigProviderClient(lb.named("MedianProviderClient"), cc)}
	m.contractTransmitter = &contractTransmitterClient{pb.NewContractTransmitterClient(m.cc)}
	m.reportCodec = &reportCodecClient{pb.NewReportCodecClient(m.cc)}
	m.medianContract = &medianContractClient{pb.NewMedianContractClient(m.cc)}
	m.onchainConfigCodec = &onchainConfigCodecClient{pb.NewOnchainConfigCodecClient(m.cc)}
	return m
}

func (m *medianProviderClient) ContractTransmitter() libocr.ContractTransmitter {
	return m.contractTransmitter
}

func (m *medianProviderClient) ReportCodec() median.ReportCodec {
	return m.reportCodec
}

func (m *medianProviderClient) MedianContract() median.MedianContract {
	return m.medianContract
}

func (m *medianProviderClient) OnchainConfigCodec() median.OnchainConfigCodec {
	return m.onchainConfigCodec
}

var _ median.ReportCodec = (*reportCodecClient)(nil)

type reportCodecClient struct {
	grpc pb.ReportCodecClient
}

func (r *reportCodecClient) BuildReport(observations []median.ParsedAttributedObservation) (report libocr.Report, err error) {
	var req pb.BuildReportRequest
	for _, o := range observations {
		req.Observations = append(req.Observations, &pb.ParsedAttributedObservation{
			Timestamp:       o.Timestamp,
			Value:           pb.NewBigIntFromInt(o.Value),
			JulesPerFeeCoin: pb.NewBigIntFromInt(o.JuelsPerFeeCoin),
			Observer:        uint32(o.Observer),
		})
	}
	var reply *pb.BuildReportReply
	reply, err = r.grpc.BuildReport(context.TODO(), &req)
	if err != nil {
		return
	}
	report = reply.Report
	return
}

func (r *reportCodecClient) MedianFromReport(report libocr.Report) (*big.Int, error) {
	reply, err := r.grpc.MedianFromReport(context.TODO(), &pb.MedianFromReportRequest{Report: report})
	if err != nil {
		return nil, err
	}
	return reply.Median.Int(), nil
}

func (r *reportCodecClient) MaxReportLength(n int) int {
	reply, err := r.grpc.MaxReportLength(context.TODO(), &pb.MaxReportLengthRequest{N: int64(n)})
	if err != nil {
		panic(err) //TODO retry https://smartcontract-it.atlassian.net/browse/BCF-2112
	}
	return int(reply.Max)
}

var _ pb.ReportCodecServer = (*reportCodecServer)(nil)

type reportCodecServer struct {
	pb.UnimplementedReportCodecServer
	impl median.ReportCodec
}

func (r *reportCodecServer) BuildReport(ctx context.Context, request *pb.BuildReportRequest) (*pb.BuildReportReply, error) {
	var obs []median.ParsedAttributedObservation
	for _, o := range request.Observations {

		val, jpfc := o.Value.Int(), o.JulesPerFeeCoin.Int()
		if o.Observer > math.MaxUint8 {
			return nil, fmt.Errorf("expected uint8 Observer (max %d) but got %d", math.MaxUint8, o.Observer)
		}
		obs = append(obs, median.ParsedAttributedObservation{
			Timestamp:       o.Timestamp,
			Value:           val,
			JuelsPerFeeCoin: jpfc,
			Observer:        commontypes.OracleID(o.Observer),
		})
	}
	report, err := r.impl.BuildReport(obs)
	if err != nil {
		return nil, err
	}
	return &pb.BuildReportReply{Report: report}, nil
}

func (r *reportCodecServer) MedianFromReport(ctx context.Context, request *pb.MedianFromReportRequest) (*pb.MedianFromReportReply, error) {
	m, err := r.impl.MedianFromReport(request.Report)
	if err != nil {
		return nil, err
	}
	return &pb.MedianFromReportReply{Median: pb.NewBigIntFromInt(m)}, nil
}

func (r *reportCodecServer) MaxReportLength(ctx context.Context, request *pb.MaxReportLengthRequest) (*pb.MaxReportLengthReply, error) {
	return &pb.MaxReportLengthReply{Max: int64(r.impl.MaxReportLength(int(request.N)))}, nil
}

var _ median.MedianContract = (*medianContractClient)(nil)

type medianContractClient struct {
	grpc pb.MedianContractClient
}

func (m *medianContractClient) LatestTransmissionDetails(ctx context.Context) (configDigest libocr.ConfigDigest, epoch uint32, round uint8, latestAnswer *big.Int, latestTimestamp time.Time, err error) {
	var reply *pb.LatestTransmissionDetailsReply
	reply, err = m.grpc.LatestTransmissionDetails(ctx, &pb.LatestTransmissionDetailsRequest{})
	if err != nil {
		return
	}
	if l := len(reply.ConfigDigest); l != 32 {
		err = fmt.Errorf("expected ConfigDigest length 32 but got %d", l)
		return
	}
	copy(configDigest[:], reply.ConfigDigest)
	epoch = reply.Epoch
	if reply.Round > math.MaxUint8 {
		err = fmt.Errorf("expected uint8 Round (max %d) but got %d", math.MaxUint8, reply.Round)
		return
	}
	round = uint8(reply.Round)
	latestAnswer = reply.LatestAnswer.Int()
	latestTimestamp = reply.LatestTimestamp.AsTime()
	return
}

func (m *medianContractClient) LatestRoundRequested(ctx context.Context, lookback time.Duration) (configDigest libocr.ConfigDigest, epoch uint32, round uint8, err error) {
	reply, err := m.grpc.LatestRoundRequested(ctx, &pb.LatestRoundRequestedRequest{Lookback: int64(lookback)})
	if err != nil {
		return
	}
	if l := len(reply.ConfigDigest); l != 32 {
		err = fmt.Errorf("expected ConfigDigest length 32 but got %d", l)
		return
	}
	copy(configDigest[:], reply.ConfigDigest)
	epoch = reply.Epoch
	if reply.Round > math.MaxUint8 {
		err = fmt.Errorf("expected uint8 Round (max %d) but got %d", math.MaxUint8, reply.Round)
		return
	}
	round = uint8(reply.Round)
	return
}

var _ pb.MedianContractServer = (*medianContractServer)(nil)

type medianContractServer struct {
	pb.UnimplementedMedianContractServer
	impl median.MedianContract
}

func (m *medianContractServer) LatestTransmissionDetails(ctx context.Context, _ *pb.LatestTransmissionDetailsRequest) (*pb.LatestTransmissionDetailsReply, error) {
	digest, epoch, round, latestAnswer, latestTimestamp, err := m.impl.LatestTransmissionDetails(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.LatestTransmissionDetailsReply{
		ConfigDigest:    digest[:],
		Epoch:           epoch,
		Round:           uint32(round),
		LatestAnswer:    pb.NewBigIntFromInt(latestAnswer),
		LatestTimestamp: timestamppb.New(latestTimestamp),
	}, nil
}

func (m *medianContractServer) LatestRoundRequested(ctx context.Context, request *pb.LatestRoundRequestedRequest) (*pb.LatestRoundRequestedReply, error) {
	digest, epoch, round, err := m.impl.LatestRoundRequested(ctx, time.Duration(request.Lookback))
	if err != nil {
		return nil, err
	}

	return &pb.LatestRoundRequestedReply{
		ConfigDigest: digest[:],
		Epoch:        epoch,
		Round:        uint32(round),
	}, nil
}

var _ median.OnchainConfigCodec = (*onchainConfigCodecClient)(nil)

type onchainConfigCodecClient struct {
	grpc pb.OnchainConfigCodecClient
}

func (o *onchainConfigCodecClient) Encode(config median.OnchainConfig) ([]byte, error) {
	req := &pb.EncodeRequest{OnchainConfig: &pb.OnchainConfig{
		Min: pb.NewBigIntFromInt(config.Min),
		Max: pb.NewBigIntFromInt(config.Max),
	}}
	reply, err := o.grpc.Encode(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	return reply.Encoded, nil
}

func (o *onchainConfigCodecClient) Decode(bytes []byte) (oc median.OnchainConfig, err error) {
	var reply *pb.DecodeReply
	reply, err = o.grpc.Decode(context.TODO(), &pb.DecodeRequest{Encoded: bytes})
	if err != nil {
		return
	}
	oc.Min, oc.Max = reply.OnchainConfig.Min.Int(), reply.OnchainConfig.Max.Int()
	return
}

var _ pb.OnchainConfigCodecServer = (*onchainConfigCodecServer)(nil)

type onchainConfigCodecServer struct {
	pb.UnimplementedOnchainConfigCodecServer
	impl median.OnchainConfigCodec
}

func (o *onchainConfigCodecServer) Encode(ctx context.Context, request *pb.EncodeRequest) (*pb.EncodeReply, error) {
	min, max := request.OnchainConfig.Min.Int(), request.OnchainConfig.Max.Int()
	b, err := o.impl.Encode(median.OnchainConfig{Max: max, Min: min})
	if err != nil {
		return nil, err
	}
	return &pb.EncodeReply{Encoded: b}, nil
}

func (o *onchainConfigCodecServer) Decode(ctx context.Context, request *pb.DecodeRequest) (*pb.DecodeReply, error) {
	oc, err := o.impl.Decode(request.Encoded)
	if err != nil {
		return nil, err
	}
	return &pb.DecodeReply{OnchainConfig: &pb.OnchainConfig{
		Min: pb.NewBigIntFromInt(oc.Min),
		Max: pb.NewBigIntFromInt(oc.Max),
	}}, nil
}
