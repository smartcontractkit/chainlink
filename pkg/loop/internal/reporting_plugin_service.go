package internal

import (
	"context"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ types.ReportingPluginClient = (*ReportingPluginServiceClient)(nil)

type ReportingPluginServiceClient struct {
	*pluginClient
	*serviceClient

	reportingPluginService pb.ReportingPluginServiceClient
}

func NewReportingPluginServiceClient(broker Broker, brokerCfg BrokerConfig, conn *grpc.ClientConn) *ReportingPluginServiceClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "ReportingPluginServiceClient")
	pc := newPluginClient(broker, brokerCfg, conn)
	return &ReportingPluginServiceClient{pluginClient: pc, reportingPluginService: pb.NewReportingPluginServiceClient(pc), serviceClient: newServiceClient(pc.brokerExt, pc)}
}

func (m *ReportingPluginServiceClient) NewReportingPluginFactory(ctx context.Context, config types.ReportingPluginServiceConfig, grpcProvider grpc.ClientConnInterface, pipelineRunner types.PipelineRunnerService, errorLog types.ErrorLog) (types.ReportingPluginFactory, error) {
	cc := m.newClientConn("ReportingPluginServiceFactory", func(ctx context.Context) (id uint32, deps resources, err error) {
		providerID, providerRes, err := m.serve("PluginProvider", proxy.NewProxy(grpcProvider))
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		errorLogID, errorLogRes, err := m.serveNew("ErrorLog", func(s *grpc.Server) {
			pb.RegisterErrorLogServer(s, &errorLogServer{impl: errorLog})
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(errorLogRes)

		pipelineRunnerID, pipelineRunnerRes, err := m.serveNew("PipelineRunner", func(s *grpc.Server) {
			pb.RegisterPipelineRunnerServiceServer(s, &pipelineRunnerServiceServer{impl: pipelineRunner})
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(pipelineRunnerRes)

		reply, err := m.reportingPluginService.NewReportingPluginFactory(ctx, &pb.NewReportingPluginFactoryRequest{
			ReportingPluginServiceConfig: &pb.ReportingPluginServiceConfig{
				ProviderType: config.ProviderType,
				Command:      config.Command,
				PluginName:   config.PluginName,
				PluginConfig: config.PluginConfig,
			},
			ProviderID:       providerID,
			ErrorLogID:       errorLogID,
			PipelineRunnerID: pipelineRunnerID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ID, nil, nil
	})
	return newReportingPluginFactoryClient(m.pluginClient.brokerExt, cc), nil
}

var _ pb.ReportingPluginServiceServer = (*reportingPluginServiceServer)(nil)

type reportingPluginServiceServer struct {
	pb.UnimplementedReportingPluginServiceServer

	*brokerExt
	impl types.ReportingPluginClient
}

func RegisterReportingPluginServiceServer(server *grpc.Server, broker Broker, brokerCfg BrokerConfig, impl types.ReportingPluginClient) error {
	pb.RegisterReportingPluginServiceServer(server, newReportingPluginServiceServer(&brokerExt{broker, brokerCfg}, impl))
	return nil
}

func newReportingPluginServiceServer(b *brokerExt, gp types.ReportingPluginClient) *reportingPluginServiceServer {
	return &reportingPluginServiceServer{brokerExt: b.withName("ReportingPluginService"), impl: gp}
}

func (m *reportingPluginServiceServer) NewReportingPluginFactory(ctx context.Context, request *pb.NewReportingPluginFactoryRequest) (*pb.NewReportingPluginFactoryReply, error) {
	errorLogConn, err := m.dial(request.ErrorLogID)
	if err != nil {
		return nil, ErrConnDial{Name: "ErrorLog", ID: request.ErrorLogID, Err: err}
	}
	errorLogRes := resource{errorLogConn, "ErrorLog"}
	errorLog := newErrorLogClient(errorLogConn)

	providerConn, err := m.dial(request.ProviderID)
	if err != nil {
		m.closeAll(errorLogRes)
		return nil, ErrConnDial{Name: "PluginProvider", ID: request.ProviderID, Err: err}
	}
	providerRes := resource{providerConn, "PluginProvider"}

	pipelineRunnerConn, err := m.dial(request.PipelineRunnerID)
	if err != nil {
		m.closeAll(errorLogRes, providerRes)
		return nil, ErrConnDial{Name: "PipelineRunner", ID: request.PipelineRunnerID, Err: err}
	}
	pipelineRunnerRes := resource{pipelineRunnerConn, "PipelineRunner"}
	pipelineRunner := newPipelineRunnerClient(pipelineRunnerConn)

	config := types.ReportingPluginServiceConfig{
		ProviderType: request.ReportingPluginServiceConfig.ProviderType,
		PluginConfig: request.ReportingPluginServiceConfig.PluginConfig,
		PluginName:   request.ReportingPluginServiceConfig.PluginName,
		Command:      request.ReportingPluginServiceConfig.Command,
	}

	factory, err := m.impl.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, errorLog)
	if err != nil {
		m.closeAll(providerRes, errorLogRes, pipelineRunnerRes)
		return nil, err
	}

	id, _, err := m.serveNew("ReportingPluginProvider", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: factory})
		pb.RegisterReportingPluginFactoryServer(s, newReportingPluginFactoryServer(factory, m.brokerExt))
	}, providerRes, errorLogRes, pipelineRunnerRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewReportingPluginFactoryReply{ID: id}, nil
}
