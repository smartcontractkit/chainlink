package ocr3

import (
	"context"
	"fmt"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/keyvalue"
	relayersetpb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/relayerset"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayerset"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/capability"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/validation"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ocr3pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
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
	config core.ReportingPluginServiceConfig,
	grpcProvider grpc.ClientConnInterface,
	pipelineRunner core.PipelineRunnerService,
	telemetryService core.TelemetryService,
	errorLog core.ErrorLog,
	capRegistry core.CapabilitiesRegistry,
	keyValueStore core.KeyValueStore,
	relayerSet core.RelayerSet,
) (core.OCR3ReportingPluginFactory, error) {
	cc := o.NewClientConn("ReportingPluginServiceFactory", func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		providerID, providerRes, err := o.Serve("PluginProvider", proxy.NewProxy(grpcProvider))
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		pipelineRunnerID, pipelineRunnerRes, err := o.ServeNew("PipelineRunner", func(s *grpc.Server) {
			pb.RegisterPipelineRunnerServiceServer(s, pipeline.NewRunnerServer(pipelineRunner))
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
			pb.RegisterErrorLogServer(s, errorlog.NewServer(errorLog))
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

		keyValueStoreID, keyValueStoreRes, err := o.ServeNew("KeyValueStore", func(s *grpc.Server) {
			pb.RegisterKeyValueStoreServer(s, keyvalue.NewServer(keyValueStore))
		})
		if err != nil {
			return 0, nil, fmt.Errorf("failed to serve KeyValueStore: %w", err)
		}
		deps.Add(keyValueStoreRes)

		relayerSetServer, relayerSetServerRes := relayerset.NewRelayerSetServer(o.Logger, relayerSet, o.BrokerExt)

		relayerSetID, relayerSetRes, err := o.ServeNew("RelayerSet", func(s *grpc.Server) {
			relayersetpb.RegisterRelayerSetServer(s, relayerSetServer)
		})

		if err != nil {
			return 0, nil, fmt.Errorf("failed to serve new relayer set: %w", err)
		}

		deps.Add(relayerSetRes)
		deps.Add(relayerSetServerRes)

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
			KeyValueStoreID:  keyValueStoreID,
			RelayerSetID:     relayerSetID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ID, nil, nil
	})
	return newReportingPluginFactoryClient(o.PluginClient.BrokerExt, cc), nil
}

func (o *ReportingPluginServiceClient) NewValidationService(ctx context.Context) (core.ValidationService, error) {
	cc := o.NewClientConn("ValidationService", func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		reply, err := o.reportingPluginService.NewValidationService(ctx, &pb.ValidationServiceRequest{})
		if err != nil {
			return 0, nil, err
		}
		return reply.ID, nil, nil
	})
	return validation.NewValidationServiceClient(o.PluginClient.BrokerExt, cc), nil
}

var _ pb.ReportingPluginServiceServer = (*reportingPluginServiceServer)(nil)

type reportingPluginServiceServer struct {
	pb.UnimplementedReportingPluginServiceServer

	*net.BrokerExt
	impl core.OCR3ReportingPluginClient
}

func (m reportingPluginServiceServer) NewValidationService(ctx context.Context, request *pb.ValidationServiceRequest) (*pb.ValidationServiceResponse, error) {
	service, err := m.impl.NewValidationService(ctx)
	if err != nil {
		return nil, err
	}

	id, _, err := m.ServeNew("ValidationService", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: service})
		pb.RegisterValidationServiceServer(s, validation.NewValidationServiceServer(service, m.BrokerExt))
	})
	if err != nil {
		return nil, err
	}

	return &pb.ValidationServiceResponse{ID: id}, nil
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

	keyValueStoreConn, err := m.Dial(request.KeyValueStoreID)
	if err != nil {
		m.CloseAll(errorLogRes, providerRes, pipelineRunnerRes, telemetryRes, capRegistryRes)
		return nil, net.ErrConnDial{Name: "KeyValueStore", ID: request.KeyValueStoreID, Err: err}
	}
	keyValueStoreRes := net.Resource{Closer: keyValueStoreConn, Name: "KeyValueStore"}
	keyValueStore := keyvalue.NewClient(keyValueStoreConn)

	config := core.ReportingPluginServiceConfig{
		ProviderType:  request.ReportingPluginServiceConfig.ProviderType,
		PluginConfig:  request.ReportingPluginServiceConfig.PluginConfig,
		PluginName:    request.ReportingPluginServiceConfig.PluginName,
		Command:       request.ReportingPluginServiceConfig.Command,
		TelemetryType: request.ReportingPluginServiceConfig.TelemetryType,
	}

	relayersetConn, err := m.Dial(request.RelayerSetID)
	if err != nil {
		m.CloseAll(errorLogRes, providerRes, pipelineRunnerRes, telemetryRes, keyValueStoreRes)
		return nil, net.ErrConnDial{Name: "RelayerSet", ID: request.RelayerSetID, Err: err}
	}
	relayerSetRes := net.Resource{Closer: relayersetConn, Name: "RelayerSet"}
	relayerSet := relayerset.NewRelayerSetClient(m.Logger, m.BrokerExt, relayersetConn)

	factory, err := m.impl.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetry, errorLog, capRegistry, keyValueStore,
		relayerSet)
	if err != nil {
		m.CloseAll(providerRes, errorLogRes, pipelineRunnerRes, telemetryRes, capRegistryRes, relayerSetRes)
		return nil, err
	}

	id, _, err := m.ServeNew("ReportingPluginProvider", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		ocr3pb.RegisterReportingPluginFactoryServer(s, newReportingPluginFactoryServer(factory, m.BrokerExt))
	}, providerRes, errorLogRes, pipelineRunnerRes, telemetryRes, capRegistryRes, keyValueStoreRes, relayerSetRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewReportingPluginFactoryReply{ID: id}, nil
}

func RegisterReportingPluginServiceServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl core.OCR3ReportingPluginClient) error {
	pb.RegisterReportingPluginServiceServer(server, newReportingPluginServiceServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func newReportingPluginServiceServer(b *net.BrokerExt, gp core.OCR3ReportingPluginClient) *reportingPluginServiceServer {
	return &reportingPluginServiceServer{BrokerExt: b.WithName("OCR3ReportingPluginService"), impl: gp}
}
