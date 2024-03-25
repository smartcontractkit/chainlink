package ocr3

import (
	"context"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/capability"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ocr3pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ReportingPluginServiceClient struct {
	*goplugin.PluginClient
	*goplugin.ServiceClient

	reportingPluginService pb.ReportingPluginServiceClient
}

func NewReportingPluginServiceClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *ReportingPluginServiceClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "ReportingPluginServiceClient")
	pc := goplugin.NewPluginClient(broker, brokerCfg, conn)
	return &ReportingPluginServiceClient{PluginClient: pc, reportingPluginService: pb.NewReportingPluginServiceClient(pc), ServiceClient: goplugin.NewServiceClient(pc.BrokerExt, pc)}
}

func (o *ReportingPluginServiceClient) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	grpcProvider grpc.ClientConnInterface,
	pipelineRunner types.PipelineRunnerService,
	telemetryService types.TelemetryService,
	errorLog types.ErrorLog,
	capRegistry types.CapabilitiesRegistry,
) (types.OCR3ReportingPluginFactory, error) {
	cc := o.NewClientConn("ReportingPluginServiceFactory", func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		providerID, providerRes, err := o.Serve("PluginProvider", proxy.NewProxy(grpcProvider))
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		pipelineRunnerID, pipelineRunnerRes, err := o.ServeNew("PipelineRunner", func(s *grpc.Server) {
			pb.RegisterPipelineRunnerServiceServer(s, &pipeline.RunnerServer{Impl: pipelineRunner})
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(pipelineRunnerRes)

		telemetryID, telemetryRes, err := o.ServeNew("Telemetry", func(s *grpc.Server) {
			pb.RegisterTelemetryServer(s, telemetry.NewTelemetryServer(telemetryService))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(telemetryRes)

		errorLogID, errorLogRes, err := o.ServeNew("ErrorLog", func(s *grpc.Server) {
			pb.RegisterErrorLogServer(s, &errorlog.Server{Impl: errorLog})
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(errorLogRes)

		capRegistryID, capRegistryRes, err := o.ServeNew("CapRegistry", func(s *grpc.Server) {
			pb.RegisterCapabilitiesRegistryServer(s, capability.NewCapabilitiesRegistryServer(o.BrokerExt, capRegistry))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(capRegistryRes)

		reply, err := o.reportingPluginService.NewReportingPluginFactory(ctx, &pb.NewReportingPluginFactoryRequest{
			ReportingPluginServiceConfig: &pb.ReportingPluginServiceConfig{
				ProviderType:  config.ProviderType,
				Command:       config.Command,
				PluginName:    config.PluginName,
				TelemetryType: config.TelemetryType,
				PluginConfig:  config.PluginConfig,
			},
			ProviderID:       providerID,
			ErrorLogID:       errorLogID,
			PipelineRunnerID: pipelineRunnerID,
			TelemetryID:      telemetryID,
			CapRegistryID:    capRegistryID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ID, nil, nil
	})
	return newReportingPluginFactoryClient(o.PluginClient.BrokerExt, cc), nil
}

var _ pb.ReportingPluginServiceServer = (*reportingPluginServiceServer)(nil)

type reportingPluginServiceServer struct {
	pb.UnimplementedReportingPluginServiceServer

	*net.BrokerExt
	impl types.OCR3ReportingPluginClient
}

func (m reportingPluginServiceServer) NewReportingPluginFactory(ctx context.Context, request *pb.NewReportingPluginFactoryRequest) (*pb.NewReportingPluginFactoryReply, error) {
	errorLogConn, err := m.Dial(request.ErrorLogID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "ErrorLog", ID: request.ErrorLogID, Err: err}
	}
	errorLogRes := net.Resource{Closer: errorLogConn, Name: "ErrorLog"}
	errorLog := errorlog.NewClient(errorLogConn)

	providerConn, err := m.Dial(request.ProviderID)
	if err != nil {
		m.CloseAll(errorLogRes)
		return nil, net.ErrConnDial{Name: "PluginProvider", ID: request.ProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "PluginProvider"}

	pipelineRunnerConn, err := m.Dial(request.PipelineRunnerID)
	if err != nil {
		m.CloseAll(errorLogRes, providerRes)
		return nil, net.ErrConnDial{Name: "PipelineRunner", ID: request.PipelineRunnerID, Err: err}
	}
	pipelineRunnerRes := net.Resource{Closer: pipelineRunnerConn, Name: "PipelineRunner"}
	pipelineRunner := pipeline.NewRunnerClient(pipelineRunnerConn)

	telemetryConn, err := m.Dial(request.TelemetryID)
	if err != nil {
		m.CloseAll(errorLogRes, providerRes, pipelineRunnerRes)
		return nil, net.ErrConnDial{Name: "Telemetry", ID: request.TelemetryID, Err: err}
	}
	telemetryRes := net.Resource{Closer: telemetryConn, Name: "Telemetry"}
	telemetry := telemetry.NewTelemetryServiceClient(telemetryConn)

	capRegistryConn, err := m.Dial(request.CapRegistryID)
	if err != nil {
		m.CloseAll(errorLogRes, providerRes, pipelineRunnerRes, telemetryRes)
		return nil, net.ErrConnDial{Name: "CapabilitiesRegistry", ID: request.CapRegistryID, Err: err}
	}
	capRegistryRes := net.Resource{Closer: capRegistryConn, Name: "CapabilitiesRegistry"}
	capRegistry := capability.NewCapabilitiesRegistryClient(capRegistryConn, m.BrokerExt)

	config := types.ReportingPluginServiceConfig{
		ProviderType:  request.ReportingPluginServiceConfig.ProviderType,
		PluginConfig:  request.ReportingPluginServiceConfig.PluginConfig,
		PluginName:    request.ReportingPluginServiceConfig.PluginName,
		Command:       request.ReportingPluginServiceConfig.Command,
		TelemetryType: request.ReportingPluginServiceConfig.TelemetryType,
	}

	factory, err := m.impl.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetry, errorLog, capRegistry)
	if err != nil {
		m.CloseAll(providerRes, errorLogRes, pipelineRunnerRes, telemetryRes, capRegistryRes)
		return nil, err
	}

	id, _, err := m.ServeNew("ReportingPluginProvider", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		ocr3pb.RegisterReportingPluginFactoryServer(s, newReportingPluginFactoryServer(factory, m.BrokerExt))
	}, providerRes, errorLogRes, pipelineRunnerRes, telemetryRes, capRegistryRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewReportingPluginFactoryReply{ID: id}, nil
}

func RegisterReportingPluginServiceServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl types.OCR3ReportingPluginClient) error {
	pb.RegisterReportingPluginServiceServer(server, newReportingPluginServiceServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func newReportingPluginServiceServer(b *net.BrokerExt, gp types.OCR3ReportingPluginClient) *reportingPluginServiceServer {
	return &reportingPluginServiceServer{BrokerExt: b.WithName("OCR3ReportingPluginService"), impl: gp}
}
