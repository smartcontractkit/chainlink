package mercury

import (
	"context"
	"fmt"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercury_v1_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v1"
	mercury_v2_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v2"
	mercury_v3_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2"

	mercury_v1_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1"
	mercury_v2_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2"
	mercury_v3_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercury_v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercury_v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

type AdapterClient struct {
	*goplugin.PluginClient
	*goplugin.ServiceClient

	mercury mercury_pb.MercuryAdapterClient
}

func NewMercuryAdapterClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *AdapterClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "MercuryAdapterClient")
	pc := goplugin.NewPluginClient(broker, brokerCfg, conn)
	return &AdapterClient{
		PluginClient:  pc,
		ServiceClient: goplugin.NewServiceClient(pc.BrokerExt, pc),
		mercury:       mercury_pb.NewMercuryAdapterClient(pc),
	}
}

func (c *AdapterClient) NewMercuryV1Factory(ctx context.Context,
	provider types.MercuryProvider, dataSource mercury_v1.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercury_v1_pb.RegisterDataSourceServer(s, mercury_v1_internal.NewDataSourceServer(dataSource))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(dsRes)

		// the proxyable resources for mercury are the Provider,  which may or may not be local to the client process. (legacy vs loopp)
		var (
			providerID  uint32
			providerRes net.Resource
		)
		if grpcProvider, ok := provider.(goplugin.GRPCClientConn); ok {
			providerID, providerRes, err = c.Serve("MercuryProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			providerID, providerRes, err = c.ServeNew("MercuryProvider", func(s *grpc.Server) {
				ocr2.RegisterPluginProviderServices(s, provider)
				registerVersionAgnosticServices(s, provider)

				mercury_pb.RegisterReportCodecV1Server(s, NewReportCodecV1Server(s, provider.ReportCodecV1()))
				mercury_pb.RegisterReportCodecV2Server(s, mercury_pb.UnimplementedReportCodecV2Server{})
				mercury_pb.RegisterReportCodecV3Server(s, mercury_pb.UnimplementedReportCodecV3Server{})
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV1Factory(ctx, &mercury_pb.NewMercuryV1FactoryRequest{
			MercuryProviderID: providerID,
			DataSourceV1ID:    dataSourceID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.MercuryV1FactoryID, deps, nil
	}

	cc := c.NewClientConn("MercuryV3Factory", newMercuryClientFn)
	return NewPluginFactoryClient(c.PluginClient.BrokerExt, cc), nil
}

func (c *AdapterClient) NewMercuryV2Factory(ctx context.Context,
	provider types.MercuryProvider, dataSource mercury_v2.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercury_v2_pb.RegisterDataSourceServer(s, mercury_v2_internal.NewDataSourceServer(dataSource))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(dsRes)

		// the proxyable resources for mercury are the Provider,  which may or may not be local to the client process. (legacy vs loopp)
		var (
			providerID  uint32
			providerRes net.Resource
		)
		if grpcProvider, ok := provider.(goplugin.GRPCClientConn); ok {
			providerID, providerRes, err = c.Serve("MercuryProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			providerID, providerRes, err = c.ServeNew("MercuryProvider", func(s *grpc.Server) {
				ocr2.RegisterPluginProviderServices(s, provider)
				registerVersionAgnosticServices(s, provider)

				mercury_pb.RegisterReportCodecV2Server(s, NewReportCodecV2Server(s, provider.ReportCodecV2()))

				mercury_pb.RegisterReportCodecV1Server(s, mercury_pb.UnimplementedReportCodecV1Server{})
				mercury_pb.RegisterReportCodecV3Server(s, mercury_pb.UnimplementedReportCodecV3Server{})
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV2Factory(ctx, &mercury_pb.NewMercuryV2FactoryRequest{
			MercuryProviderID: providerID,
			DataSourceV2ID:    dataSourceID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.MercuryV2FactoryID, deps, nil
	}

	cc := c.NewClientConn("MercuryV2Factory", newMercuryClientFn)
	return NewPluginFactoryClient(c.PluginClient.BrokerExt, cc), nil
}

func (c *AdapterClient) NewMercuryV3Factory(ctx context.Context,
	provider types.MercuryProvider, dataSource mercury_v3.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercury_v3_pb.RegisterDataSourceServer(s, mercury_v3_internal.NewDataSourceServer(dataSource))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(dsRes)

		// the proxyable resources for mercury are the Provider,  which may or may not be local to the client process. (legacy vs loopp)
		var (
			providerID  uint32
			providerRes net.Resource
		)
		// loop mode; proxy to the relayer
		if grpcProvider, ok := provider.(goplugin.GRPCClientConn); ok {
			providerID, providerRes, err = c.Serve("MercuryProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			// legacy mode; serve the provider locally in the client process (ie the core node)
			providerID, providerRes, err = c.ServeNew("MercuryProvider", func(s *grpc.Server) {
				ocr2.RegisterPluginProviderServices(s, provider)
				registerVersionAgnosticServices(s, provider)

				mercury_pb.RegisterReportCodecV3Server(s, NewReportCodecV3Server(s, provider.ReportCodecV3()))
				// don't register the other codecs, as they are not used in v3
				mercury_pb.RegisterReportCodecV1Server(s, mercury_pb.UnimplementedReportCodecV1Server{})
				mercury_pb.RegisterReportCodecV2Server(s, mercury_pb.UnimplementedReportCodecV2Server{})
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV3Factory(ctx, &mercury_pb.NewMercuryV3FactoryRequest{
			MercuryProviderID: providerID,
			DataSourceV3ID:    dataSourceID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.MercuryV3FactoryID, deps, nil
	}

	cc := c.NewClientConn("MercuryV3Factory", newMercuryClientFn)
	return NewPluginFactoryClient(c.PluginClient.BrokerExt, cc), nil
}

var _ mercury_pb.MercuryAdapterServer = (*AdapterServer)(nil)

type AdapterServer struct {
	mercury_pb.UnimplementedMercuryAdapterServer

	*net.BrokerExt
	impl types.PluginMercury
}

func RegisterMercuryAdapterServer(s *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl types.PluginMercury) error {
	mercury_pb.RegisterMercuryAdapterServer(s, NewMercuryAdapterServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func NewMercuryAdapterServer(b *net.BrokerExt, impl types.PluginMercury) *AdapterServer {
	return &AdapterServer{BrokerExt: b.WithName("MercuryAdapter"), impl: impl}
}

func (ms *AdapterServer) NewMercuryV1Factory(ctx context.Context, req *mercury_pb.NewMercuryV1FactoryRequest) (*mercury_pb.NewMercuryV1FactoryReply, error) {
	// declared so we can clean up open resources
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			ms.CloseAll(deps...)
		}
	}()

	dsConn, err := ms.Dial(req.DataSourceV1ID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "DataSourceV1", ID: req.DataSourceV1ID, Err: err}
	}
	dsRes := net.Resource{Closer: dsConn, Name: "DataSourceV1"}
	deps.Add(dsRes)
	ds := mercury_v1_internal.NewDataSourceClient(dsConn)

	providerConn, err := ms.Dial(req.MercuryProviderID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "MercuryProvider", ID: req.MercuryProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "MercuryProvider"}
	deps.Add(providerRes)
	provider := NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV1Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV1Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV1Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercury_pb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v1 factory server: %w", err)
	}

	return &mercury_pb.NewMercuryV1FactoryReply{MercuryV1FactoryID: id}, nil
}

func (ms *AdapterServer) NewMercuryV2Factory(ctx context.Context, req *mercury_pb.NewMercuryV2FactoryRequest) (*mercury_pb.NewMercuryV2FactoryReply, error) {
	// declared so we can clean up open resources
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			ms.CloseAll(deps...)
		}
	}()

	dsConn, err := ms.Dial(req.DataSourceV2ID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "DataSourceV2", ID: req.DataSourceV2ID, Err: err}
	}
	dsRes := net.Resource{Closer: dsConn, Name: "DataSourceV2"}
	deps.Add(dsRes)
	ds := mercury_v2_internal.NewDataSourceClient(dsConn)

	providerConn, err := ms.Dial(req.MercuryProviderID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "MercuryProvider", ID: req.MercuryProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "MercuryProvider"}
	deps.Add(providerRes)
	provider := NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV2Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV2Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV2Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercury_pb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v2 factory server: %w", err)
	}

	return &mercury_pb.NewMercuryV2FactoryReply{MercuryV2FactoryID: id}, nil
}

func (ms *AdapterServer) NewMercuryV3Factory(ctx context.Context, req *mercury_pb.NewMercuryV3FactoryRequest) (*mercury_pb.NewMercuryV3FactoryReply, error) {
	// declared so we can clean up open resources
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			ms.CloseAll(deps...)
		}
	}()

	dsConn, err := ms.Dial(req.DataSourceV3ID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "DataSourceV3", ID: req.DataSourceV3ID, Err: err}
	}
	dsRes := net.Resource{Closer: dsConn, Name: "DataSourceV3"}
	deps.Add(dsRes)
	ds := mercury_v3_internal.NewDataSourceClient(dsConn)

	providerConn, err := ms.Dial(req.MercuryProviderID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "MercuryProvider", ID: req.MercuryProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "MercuryProvider"}
	deps.Add(providerRes)
	provider := NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV3Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV3Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV3Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercury_pb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v3 factory server: %w", err)
	}

	return &mercury_pb.NewMercuryV3FactoryReply{MercuryV3FactoryID: id}, nil
}

var (
	_ types.MercuryProvider = (*ProviderClient)(nil)
	// in practice, inherited from pluginProviderClient.
	_ goplugin.GRPCClientConn = (*ProviderClient)(nil)
)

type ProviderClient struct {
	*ocr2.PluginProviderClient
	reportCodecV3      mercury_v3.ReportCodec
	reportCodecV2      mercury_v2.ReportCodec
	reportCodecV1      mercury_v1.ReportCodec
	onchainConfigCodec mercury.OnchainConfigCodec
	serverFetcher      mercury.ServerFetcher
	chainReader        types.ChainReader
	mercuryChainReader mercury.ChainReader
}

func NewProviderClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *ProviderClient {
	m := &ProviderClient{PluginProviderClient: ocr2.NewPluginProviderClient(b.WithName("MercuryProviderClient"), cc)}

	m.reportCodecV1 = NewReportCodecV1Client(mercury_v1_internal.NewReportCodecClient(cc))
	m.reportCodecV2 = NewReportCodecV2Client(mercury_v2_internal.NewReportCodecClient(cc))
	m.reportCodecV3 = NewReportCodecV3Client(mercury_v3_internal.NewReportCodecClient(cc))

	m.onchainConfigCodec = NewOnchainConfigCodecClient(cc)
	m.serverFetcher = NewServerFetcherClient(cc)
	m.mercuryChainReader = NewChainReaderClient(cc)

	m.chainReader = chainreader.NewClient(b, cc)
	return m
}

func (m *ProviderClient) ReportCodecV3() mercury_v3.ReportCodec {
	return m.reportCodecV3
}

func (m *ProviderClient) ReportCodecV2() mercury_v2.ReportCodec {
	return m.reportCodecV2
}

func (m *ProviderClient) ReportCodecV1() mercury_v1.ReportCodec {
	return m.reportCodecV1
}

func (m *ProviderClient) OnchainConfigCodec() mercury.OnchainConfigCodec {
	return m.onchainConfigCodec
}

func (m *ProviderClient) ChainReader() types.ChainReader {
	return m.chainReader
}

func (m *ProviderClient) MercuryChainReader() mercury.ChainReader {
	return m.mercuryChainReader
}

func (m *ProviderClient) MercuryServerFetcher() mercury.ServerFetcher {
	return m.serverFetcher
}

func registerVersionAgnosticServices(s *grpc.Server, provider types.MercuryProvider) {
	mercury_pb.RegisterOnchainConfigCodecServer(s, NewOnchainConfigCodecServer(provider.OnchainConfigCodec()))
	mercury_pb.RegisterServerFetcherServer(s, NewServerFetcherServer(provider.MercuryServerFetcher()))
	mercury_pb.RegisterMercuryChainReaderServer(s, NewChainReaderServer(provider.MercuryChainReader()))
}

func RegisterProviderServices(s *grpc.Server, provider types.MercuryProvider) {
	registerVersionAgnosticServices(s, provider)

	mercury_pb.RegisterReportCodecV1Server(s, NewReportCodecV1Server(s, provider.ReportCodecV1()))
	mercury_pb.RegisterReportCodecV2Server(s, NewReportCodecV2Server(s, provider.ReportCodecV2()))
	mercury_pb.RegisterReportCodecV3Server(s, NewReportCodecV3Server(s, provider.ReportCodecV3()))
}
