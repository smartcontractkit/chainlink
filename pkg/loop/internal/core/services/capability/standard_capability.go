package capability

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/keyvalue"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	relayersetpb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/relayerset"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayerset"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type StandardCapability interface {
	services.Service
	capabilities.BaseCapability
	Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore,
		capabilityRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
		pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error
}

type StandardCapabilityClient struct {
	*goplugin.PluginClient
	capabilitiespb.StandardCapabilityClient
	*baseCapabilityClient
	*goplugin.ServiceClient
	*net.BrokerExt

	resources []net.Resource
}

var _ StandardCapability = (*StandardCapabilityClient)(nil)

func NewStandardCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) *StandardCapabilityClient {
	return &StandardCapabilityClient{
		PluginClient:             goplugin.NewPluginClient(brokerExt.Broker, brokerExt.BrokerConfig, conn),
		ServiceClient:            goplugin.NewServiceClient(brokerExt, conn),
		StandardCapabilityClient: capabilitiespb.NewStandardCapabilityClient(conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
		BrokerExt:                brokerExt,
	}
}

func (c *StandardCapabilityClient) Initialise(ctx context.Context, config string, telemetryService core.TelemetryService,
	keyValueStore core.KeyValueStore, capabilitiesRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error {
	telemetryID, telemetryRes, err := c.ServeNew("Telemetry", func(s *grpc.Server) {
		pb.RegisterTelemetryServer(s, telemetry.NewTelemetryServer(telemetryService))
	})

	if err != nil {
		return fmt.Errorf("failed to serve new telemetry: %w", err)
	}
	var resources []net.Resource
	resources = append(resources, telemetryRes)

	keyValueStoreID, keyValueStoreRes, err := c.ServeNew("KeyValueStore", func(s *grpc.Server) {
		pb.RegisterKeyValueStoreServer(s, keyvalue.NewServer(keyValueStore))
	})
	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to serve new key value store: %w", err)
	}
	resources = append(resources, keyValueStoreRes)

	capabilitiesRegistryID, capabilityRegistryResource, err := c.ServeNew("CapabilitiesRegistry", func(s *grpc.Server) {
		pb.RegisterCapabilitiesRegistryServer(s, NewCapabilitiesRegistryServer(c.BrokerExt, capabilitiesRegistry))
	})
	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to serve new key value store: %w", err)
	}
	resources = append(resources, capabilityRegistryResource)

	errorLogID, errorLogRes, err := c.ServeNew("ErrorLog", func(s *grpc.Server) {
		pb.RegisterErrorLogServer(s, errorlog.NewServer(errorLog))
	})
	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to serve error log: %w", err)
	}
	resources = append(resources, errorLogRes)

	pipelineRunnerID, pipelineRunnerRes, err := c.ServeNew("PipelineRunner", func(s *grpc.Server) {
		pb.RegisterPipelineRunnerServiceServer(s, pipeline.NewRunnerServer(pipelineRunner))
	})
	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to serve pipeline runner: %w", err)
	}
	resources = append(resources, pipelineRunnerRes)

	relayerSetServer, relayerSetServerRes := relayerset.NewRelayerSetServer(c.Logger, relayerSet, c.BrokerExt)
	resources = append(resources, relayerSetServerRes)

	relayerSetID, relayerSetRes, err := c.ServeNew("RelayerSet", func(s *grpc.Server) {
		relayersetpb.RegisterRelayerSetServer(s, relayerSetServer)
	})
	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to serve relayer set: %w", err)
	}

	resources = append(resources, relayerSetRes)

	_, err = c.StandardCapabilityClient.Initialise(ctx, &capabilitiespb.InitialiseRequest{
		Config:           config,
		ErrorLogId:       errorLogID,
		PipelineRunnerId: pipelineRunnerID,
		TelemetryId:      telemetryID,
		CapRegistryId:    capabilitiesRegistryID,
		KeyValueStoreId:  keyValueStoreID,
		RelayerSetId:     relayerSetID,
	})

	if err != nil {
		c.CloseAll(resources...)
		return fmt.Errorf("failed to initialise standard capability: %w", err)
	}

	c.resources = resources

	return nil
}

func (c *StandardCapabilityClient) Close() error {
	c.CloseAll(c.resources...)
	return c.ServiceClient.Close()
}

type standardCapabilityServer struct {
	capabilitiespb.UnimplementedStandardCapabilityServer
	*net.BrokerExt
	impl StandardCapability

	resources []net.Resource
}

func newStandardCapabilityServer(brokerExt *net.BrokerExt, impl StandardCapability) *standardCapabilityServer {
	return &standardCapabilityServer{
		impl:      impl,
		BrokerExt: brokerExt,
	}
}

var _ capabilitiespb.StandardCapabilityServer = (*standardCapabilityServer)(nil)

func RegisterStandardCapabilityServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl StandardCapability) error {
	bext := &net.BrokerExt{
		BrokerConfig: brokerCfg,
		Broker:       broker,
	}

	capabilityServer := newStandardCapabilityServer(bext, impl)
	capabilitiespb.RegisterStandardCapabilityServer(server, capabilityServer)
	capabilitiespb.RegisterBaseCapabilityServer(server, newBaseCapabilityServer(impl))
	pb.RegisterServiceServer(server, &goplugin.ServiceServer{Srv: &resourceClosingServer{
		StandardCapability: impl,
		server:             capabilityServer,
	}})
	return nil
}

func (s *standardCapabilityServer) Initialise(ctx context.Context, request *capabilitiespb.InitialiseRequest) (*emptypb.Empty, error) {
	telemetryConn, err := s.Dial(request.TelemetryId)
	if err != nil {
		return nil, net.ErrConnDial{Name: "Telemetry", ID: request.TelemetryId, Err: err}
	}

	var resources []net.Resource
	resources = append(resources, net.Resource{Closer: telemetryConn, Name: "TelemetryConn"})
	telemetry := telemetry.NewTelemetryServiceClient(telemetryConn)

	keyValueStoreConn, err := s.Dial(request.KeyValueStoreId)
	if err != nil {
		s.CloseAll(resources...)
		return nil, net.ErrConnDial{Name: "KeyValueStore", ID: request.KeyValueStoreId, Err: err}
	}
	resources = append(resources, net.Resource{Closer: keyValueStoreConn, Name: "KeyValueStoreConn"})
	keyValueStore := keyvalue.NewClient(keyValueStoreConn)

	capabilitiesRegistryConn, err := s.Dial(request.CapRegistryId)
	if err != nil {
		s.CloseAll(resources...)
		return nil, net.ErrConnDial{Name: "CapabilitiesRegistry", ID: request.CapRegistryId, Err: err}
	}
	resources = append(resources, net.Resource{Closer: capabilitiesRegistryConn, Name: "CapabilitiesRegistryConn"})
	capabilitiesRegistry := NewCapabilitiesRegistryClient(capabilitiesRegistryConn, s.BrokerExt)

	errorLogConn, err := s.Dial(request.ErrorLogId)
	if err != nil {
		s.CloseAll(resources...)
		return nil, net.ErrConnDial{Name: "ErrorLog", ID: request.ErrorLogId, Err: err}
	}
	resources = append(resources, net.Resource{Closer: errorLogConn, Name: "ErrorLog"})
	errorLog := errorlog.NewClient(errorLogConn)

	pipelineRunnerConn, err := s.Dial(request.PipelineRunnerId)
	if err != nil {
		s.CloseAll(resources...)
		return nil, net.ErrConnDial{Name: "PipelineRunner", ID: request.PipelineRunnerId, Err: err}
	}
	resources = append(resources, net.Resource{Closer: pipelineRunnerConn, Name: "PipelineRunner"})
	pipelineRunner := pipeline.NewRunnerClient(pipelineRunnerConn)

	relayersetConn, err := s.Dial(request.RelayerSetId)
	if err != nil {
		s.CloseAll(resources...)
		return nil, net.ErrConnDial{Name: "RelayerSet", ID: request.RelayerSetId, Err: err}
	}
	resources = append(resources, net.Resource{Closer: relayersetConn, Name: "RelayerSet"})
	relayerSet := relayerset.NewRelayerSetClient(s.Logger, s.BrokerExt, relayersetConn)

	if err = s.impl.Initialise(ctx, request.Config, telemetry, keyValueStore, capabilitiesRegistry, errorLog, pipelineRunner, relayerSet); err != nil {
		s.CloseAll(resources...)
		return nil, fmt.Errorf("failed to initialise standard capability: %w", err)
	}

	s.resources = resources

	return &emptypb.Empty{}, nil
}

type resourceClosingServer struct {
	StandardCapability
	server *standardCapabilityServer
}

func (r *resourceClosingServer) Close() error {
	r.server.CloseAll(r.server.resources...)
	return r.StandardCapability.Close()
}
