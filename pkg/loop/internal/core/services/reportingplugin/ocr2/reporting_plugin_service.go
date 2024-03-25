package ocr2

import (
	"context"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ReportingPluginClient = (*ReportingPluginServiceClient)(nil)

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

func (m *ReportingPluginServiceClient) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	grpcProvider grpc.ClientConnInterface,
	pipelineRunner types.PipelineRunnerService,
	telemetryService types.TelemetryService,
	errorLog types.ErrorLog,
) (types.ReportingPluginFactory, error) {
	cc := m.NewClientConn("ReportingPluginServiceFactory", func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		providerID, providerRes, err := m.Serve("PluginProvider", proxy.NewProxy(grpcProvider))
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		pipelineRunnerID, pipelineRunnerRes, err := m.ServeNew("PipelineRunner", func(s *grpc.Server) {
			pb.RegisterPipelineRunnerServiceServer(s, pipeline.NewRunnerServer(pipelineRunner))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(pipelineRunnerRes)

		telemetryID, telemetryRes, err := m.ServeNew("Telemetry", func(s *grpc.Server) {
			pb.RegisterTelemetryServer(s, telemetry.NewTelemetryServer(telemetryService))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(telemetryRes)

		errorLogID, errorLogRes, err := m.ServeNew("ErrorLog", func(s *grpc.Server) {
			pb.RegisterErrorLogServer(s, errorlog.NewServer(errorLog))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(errorLogRes)

		reply, err := m.reportingPluginService.NewReportingPluginFactory(ctx, &pb.NewReportingPluginFactoryRequest{
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
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ID, nil, nil
	})
	return NewReportingPluginFactoryClient(m.PluginClient.BrokerExt, cc), nil
}

var _ pb.ReportingPluginServiceServer = (*reportingPluginServiceServer)(nil)

type reportingPluginServiceServer struct {
	pb.UnimplementedReportingPluginServiceServer

	*net.BrokerExt
	impl types.ReportingPluginClient
}

func RegisterReportingPluginServiceServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl types.ReportingPluginClient) error {
	pb.RegisterReportingPluginServiceServer(server, newReportingPluginServiceServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func newReportingPluginServiceServer(b *net.BrokerExt, gp types.ReportingPluginClient) *reportingPluginServiceServer {
	return &reportingPluginServiceServer{BrokerExt: b.WithName("ReportingPluginService"), impl: gp}
}

func (m *reportingPluginServiceServer) NewReportingPluginFactory(ctx context.Context, request *pb.NewReportingPluginFactoryRequest) (*pb.NewReportingPluginFactoryReply, error) {
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

	config := types.ReportingPluginServiceConfig{
		ProviderType:  request.ReportingPluginServiceConfig.ProviderType,
		PluginConfig:  request.ReportingPluginServiceConfig.PluginConfig,
		PluginName:    request.ReportingPluginServiceConfig.PluginName,
		Command:       request.ReportingPluginServiceConfig.Command,
		TelemetryType: request.ReportingPluginServiceConfig.TelemetryType,
	}

	factory, err := m.impl.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetry, errorLog)
	if err != nil {
		m.CloseAll(providerRes, errorLogRes, pipelineRunnerRes, telemetryRes)
		return nil, err
	}

	id, _, err := m.ServeNew("ReportingPluginProvider", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		pb.RegisterReportingPluginFactoryServer(s, NewReportingPluginFactoryServer(factory, m.BrokerExt))
	}, providerRes, errorLogRes, pipelineRunnerRes, telemetryRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewReportingPluginFactoryReply{ID: id}, nil
}
