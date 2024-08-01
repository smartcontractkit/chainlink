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
	mercurypb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercuryv1pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v1"
	mercuryv2pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v2"
	mercuryv3pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v3"
	mercuryv4pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury/v4"
	mercuryprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury"
	mercury_v1_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1"
	mercury_v2_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2"
	mercury_v3_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3"
	mercury_v4_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v4"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	mercuryv1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	mercuryv4 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"
)

type AdapterClient struct {
	*goplugin.PluginClient
	*goplugin.ServiceClient

	mercury mercurypb.MercuryAdapterClient
}

func NewMercuryAdapterClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *AdapterClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "MercuryAdapterClient")
	pc := goplugin.NewPluginClient(broker, brokerCfg, conn)
	return &AdapterClient{
		PluginClient:  pc,
		ServiceClient: goplugin.NewServiceClient(pc.BrokerExt, pc),
		mercury:       mercurypb.NewMercuryAdapterClient(pc),
	}
}

func (c *AdapterClient) NewMercuryV1Factory(ctx context.Context,
	provider types.MercuryProvider, dataSource mercuryv1.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercuryv1pb.RegisterDataSourceServer(s, mercury_v1_internal.NewDataSourceServer(dataSource))
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
				mercuryprovider.RegisterProviderServicesV1(s, provider)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV1Factory(ctx, &mercurypb.NewMercuryV1FactoryRequest{
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
	provider types.MercuryProvider, dataSource mercuryv2.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercuryv2pb.RegisterDataSourceServer(s, mercury_v2_internal.NewDataSourceServer(dataSource))
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
				mercuryprovider.RegisterProviderServicesV2(s, provider)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV2Factory(ctx, &mercurypb.NewMercuryV2FactoryRequest{
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
	provider types.MercuryProvider, dataSource mercuryv3.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercuryv3pb.RegisterDataSourceServer(s, mercury_v3_internal.NewDataSourceServer(dataSource))
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
				mercuryprovider.RegisterProviderServicesV3(s, provider)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV3Factory(ctx, &mercurypb.NewMercuryV3FactoryRequest{
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

func (c *AdapterClient) NewMercuryV4Factory(ctx context.Context,
	provider types.MercuryProvider, dataSource mercuryv4.DataSource,
) (types.MercuryPluginFactory, error) {
	// every time a new client is created, we have to ensure that all the external dependencies are satisfied.
	// at this layer of the stack, all of those dependencies are other gRPC services.
	// some of those services are hosted in the same process as the client itself and others may be remote.
	newMercuryClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// the local resources for mercury are the DataSource
		dataSourceID, dsRes, err := c.ServeNew("DataSource", func(s *grpc.Server) {
			mercuryv4pb.RegisterDataSourceServer(s, mercury_v4_internal.NewDataSourceServer(dataSource))
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
				mercuryprovider.RegisterProviderServicesV4(s, provider)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		reply, err := c.mercury.NewMercuryV4Factory(ctx, &mercurypb.NewMercuryV4FactoryRequest{
			MercuryProviderID: providerID,
			DataSourceV4ID:    dataSourceID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.MercuryV4FactoryID, deps, nil
	}

	cc := c.NewClientConn("MercuryV4Factory", newMercuryClientFn)
	return NewPluginFactoryClient(c.PluginClient.BrokerExt, cc), nil
}

var _ mercurypb.MercuryAdapterServer = (*AdapterServer)(nil)

type AdapterServer struct {
	mercurypb.UnimplementedMercuryAdapterServer

	*net.BrokerExt
	impl types.PluginMercury
}

func RegisterMercuryAdapterServer(s *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl types.PluginMercury) error {
	mercurypb.RegisterMercuryAdapterServer(s, NewMercuryAdapterServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func NewMercuryAdapterServer(b *net.BrokerExt, impl types.PluginMercury) *AdapterServer {
	return &AdapterServer{BrokerExt: b.WithName("MercuryAdapter"), impl: impl}
}

func (ms *AdapterServer) NewMercuryV1Factory(ctx context.Context, req *mercurypb.NewMercuryV1FactoryRequest) (*mercurypb.NewMercuryV1FactoryReply, error) {
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
	provider := mercuryprovider.NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV1Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV1Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV1Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercurypb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v1 factory server: %w", err)
	}

	return &mercurypb.NewMercuryV1FactoryReply{MercuryV1FactoryID: id}, nil
}

func (ms *AdapterServer) NewMercuryV2Factory(ctx context.Context, req *mercurypb.NewMercuryV2FactoryRequest) (*mercurypb.NewMercuryV2FactoryReply, error) {
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
	provider := mercuryprovider.NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV2Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV2Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV2Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercurypb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v2 factory server: %w", err)
	}

	return &mercurypb.NewMercuryV2FactoryReply{MercuryV2FactoryID: id}, nil
}

func (ms *AdapterServer) NewMercuryV3Factory(ctx context.Context, req *mercurypb.NewMercuryV3FactoryRequest) (*mercurypb.NewMercuryV3FactoryReply, error) {
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
	provider := mercuryprovider.NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV3Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV3Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV3Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercurypb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v3 factory server: %w", err)
	}

	return &mercurypb.NewMercuryV3FactoryReply{MercuryV3FactoryID: id}, nil
}

func (ms *AdapterServer) NewMercuryV4Factory(ctx context.Context, req *mercurypb.NewMercuryV4FactoryRequest) (*mercurypb.NewMercuryV4FactoryReply, error) {
	// declared so we can clean up open resources
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			ms.CloseAll(deps...)
		}
	}()

	dsConn, err := ms.Dial(req.DataSourceV4ID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "DataSourceV3", ID: req.DataSourceV4ID, Err: err}
	}
	dsRes := net.Resource{Closer: dsConn, Name: "DataSourceV4"}
	deps.Add(dsRes)
	ds := mercury_v4_internal.NewDataSourceClient(dsConn)

	providerConn, err := ms.Dial(req.MercuryProviderID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "MercuryProvider", ID: req.MercuryProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "MercuryProvider"}
	deps.Add(providerRes)
	provider := mercuryprovider.NewProviderClient(ms.BrokerExt, providerConn)
	factory, err := ms.impl.NewMercuryV4Factory(ctx, provider, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create MercuryV4Factory: %w", err)
	}

	id, _, err := ms.ServeNew("MercuryV4Factory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		mercurypb.RegisterMercuryPluginFactoryServer(s, newMercuryPluginFactoryServer(factory, ms.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new v4 factory server: %w", err)
	}

	return &mercurypb.NewMercuryV4FactoryReply{MercuryV4FactoryID: id}, nil
}
