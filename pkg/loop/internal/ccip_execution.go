package internal

import (
	"context"
	"fmt"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipinternal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// ExecutionLOOPClient is a client is run on the core node to connect to the execution LOOP server.
type ExecutionLOOPClient struct {
	// hashicorp plugin client
	*PluginClient
	// client to base service
	*ServiceClient

	// creates new execution factory instances
	generator ccippb.ExecutionFactoryGeneratorClient
}

func NewExecutionLOOPClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *ExecutionLOOPClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "ExecutionLOOPClient")
	pc := NewPluginClient(broker, brokerCfg, conn)
	return &ExecutionLOOPClient{
		PluginClient:  pc,
		ServiceClient: NewServiceClient(pc.BrokerExt, pc),
		generator:     ccippb.NewExecutionFactoryGeneratorClient(pc),
	}
}

// NewExecutionFactory creates a new reporting plugin factory client.
// In practice this client is called by the core node.
// The reporting plugin factory client is a client to the LOOP server, which
// is run as an external process via hashicorp plugin. If the given provider is a GRPCClientConn, then the provider is proxied to the
// to the relayer, which is its own process via hashicorp plugin. If the provider is not a GRPCClientConn, then the provider is a local
// to the core node. The core must wrap the provider in a grpc server and serve it locally.
// func (c *ExecutionLOOPClient) NewExecutionFactory(ctx context.Context, provider types.CCIPExecProvider, config types.CCIPExecFactoryGeneratorConfig) (types.ReportingPluginFactory, error) {.
func (c *ExecutionLOOPClient) NewExecutionFactory(ctx context.Context, provider types.CCIPExecProvider) (types.ReportingPluginFactory, error) {
	newExecClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// TODO are there any local resources that need to be passed to the executor and started as a server?

		// the proxyable resources are the Provider,  which may or may not be local to the client process. (legacy vs loopp)
		var (
			providerID       uint32
			providerResource net.Resource
		)
		if grpcProvider, ok := provider.(GRPCClientConn); ok {
			// TODO: BCF-3061 ccip provider can create new services. the proxying needs to be augmented
			// to intercept and route to the created services. also, need to prevent leaks.
			providerID, providerResource, err = c.Serve("ExecProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			// loop client runs in the core node. if the provider is not a grpc client conn, then we are in legacy mode
			// and need to serve all the required services locally.
			providerID, providerResource, err = c.ServeNew("ExecProvider", func(s *grpc.Server) {
				registerPluginProviderServices(s, provider)
				registerCustomExecutionProviderServices(s, provider, c.BrokerExt)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerResource)

		resp, err := c.generator.NewExecutionFactory(ctx, &ccippb.NewExecutionFactoryRequest{
			ProviderServiceId: providerID,
		})
		if err != nil {
			return 0, nil, err
		}
		return resp.ExecutionFactoryServiceId, deps, nil
	}
	cc := c.NewClientConn("ExecutionFactory", newExecClientFn)
	return newReportingPluginFactoryClient(c.BrokerExt, cc), nil
}

func registerCustomExecutionProviderServices(s *grpc.Server, provider types.CCIPExecProvider, brokerExt *net.BrokerExt) {
	// register the handler for the custom methods of the provider eg NewOffRampReader
	ccippb.RegisterExecutionCustomHandlersServer(s, newExecProviderServer(provider, brokerExt))
}

// ExecutionLOOPServer is a server that runs the execution LOOP.
type ExecutionLOOPServer struct {
	ccippb.UnimplementedExecutionFactoryGeneratorServer

	*net.BrokerExt
	impl types.CCIPExecutionFactoryGenerator
}

func RegisterExecutionLOOPServer(s *grpc.Server, b net.Broker, cfg net.BrokerConfig, impl types.CCIPExecutionFactoryGenerator) error {
	ext := &net.BrokerExt{Broker: b, BrokerConfig: cfg}
	ccippb.RegisterExecutionFactoryGeneratorServer(s, newExecutionLOOPServer(impl, ext))
	return nil
}

func newExecutionLOOPServer(impl types.CCIPExecutionFactoryGenerator, b *net.BrokerExt) *ExecutionLOOPServer {
	return &ExecutionLOOPServer{impl: impl, BrokerExt: b.WithName("ExecutionLOOPServer")}
}

func (r *ExecutionLOOPServer) NewExecutionFactory(ctx context.Context, request *ccippb.NewExecutionFactoryRequest) (*ccippb.NewExecutionFactoryResponse, error) {
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			r.CloseAll(deps...)
		}
	}()

	// lookup the provider service
	providerConn, err := r.Dial(request.ProviderServiceId)
	if err != nil {
		return nil, net.ErrConnDial{Name: "ExecProvider", ID: request.ProviderServiceId, Err: err}
	}
	deps.Add(net.Resource{Closer: providerConn, Name: "ExecProvider"})
	provider := newExecProviderClient(r.BrokerExt, providerConn)

	factory, err := r.impl.NewExecutionFactory(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new execution factory: %w", err)
	}

	id, _, err := r.ServeNew("ExecutionFactory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &ServiceServer{Srv: factory})
		pb.RegisterReportingPluginFactoryServer(s, newReportingPluginFactoryServer(factory, r.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to serve new execution factory: %w", err)
	}
	return &ccippb.NewExecutionFactoryResponse{ExecutionFactoryServiceId: id}, nil
}

var (
	_ types.CCIPExecProvider = (*execProviderClient)(nil)
	_ GRPCClientConn         = (*execProviderClient)(nil)
)

type execProviderClient struct {
	*pluginProviderClient

	// must be shared with the server
	*net.BrokerExt
	grpcClient ccippb.ExecutionCustomHandlersClient
}

func newExecProviderClient(b *net.BrokerExt, conn grpc.ClientConnInterface) *execProviderClient {
	pluginProviderClient := newPluginProviderClient(b, conn)
	grpc := ccippb.NewExecutionCustomHandlersClient(conn)
	return &execProviderClient{
		pluginProviderClient: pluginProviderClient,
		BrokerExt:            b,
		grpcClient:           grpc,
	}
}

// NewCommitStoreReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewCommitStoreReader(ctx context.Context, addr cciptypes.Address) (cciptypes.CommitStoreReader, error) {
	req := ccippb.NewCommitStoreReaderRequest{Address: string(addr)}

	resp, err := e.grpcClient.NewCommitStoreReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: this works because the broker is shared and the id refers to a resource served by the broker
	commitStoreConn, err := e.BrokerExt.Dial(uint32(resp.CommitStoreReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup off ramp reader service at %d: %w", resp.CommitStoreReaderServiceId, err)
	}
	// need to wrap grpc commitStore into the desired interface
	commitStore := ccipinternal.NewCommitStoreReaderGRPCClient(commitStoreConn, e.BrokerExt)

	return commitStore, nil
}

// NewOffRampReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewOffRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OffRampReader, error) {
	req := ccippb.NewOffRampReaderRequest{Address: string(addr)}

	resp, err := e.grpcClient.NewOffRampReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: this works because the broker is shared and the id refers to a resource served by the broker
	offRampConn, err := e.BrokerExt.Dial(uint32(resp.OfframpReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup off ramp reader service at %d: %w", resp.OfframpReaderServiceId, err)
	}
	// need to wrap grpc offRamp into the desired interface
	offRamp := ccipinternal.NewOffRampReaderGRPCClient(offRampConn, e.BrokerExt)

	return offRamp, nil
}

// NewOnRampReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewOnRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OnRampReader, error) {
	req := ccippb.NewOnRampReaderRequest{Address: string(addr)}

	resp, err := e.grpcClient.NewOnRampReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	// currently this only work for an embedded relayer
	// because the broker is shared  between the core node and relayer
	// this effectively let us proxy connects to resources spawn by the embedded relay
	// by hijacking the shared broker. id refers to a resource served by the shared broker
	onRampConn, err := e.BrokerExt.Dial(uint32(resp.OnrampReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup on ramp reader service at %d: %w", resp.OnrampReaderServiceId, err)
	}
	// need to wrap grpc onRamp into the desired interface
	onRamp := ccipinternal.NewOnRampReaderClient(onRampConn)

	// how to convert resp to cciptypes.OnRampReader? i have an id and need to hydrate that into an instance of OnRampReader
	return onRamp, nil
}

// NewPriceRegistryReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error) {
	req := ccippb.NewPriceRegistryReaderRequest{Address: string(addr)}
	resp, err := e.grpcClient.NewPriceRegistryReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	priceReaderConn, err := e.BrokerExt.Dial(uint32(resp.PriceRegistryReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup price registry reader service at %d: %w", resp.PriceRegistryReaderServiceId, err)
	}
	// need to wrap grpc priceReader into the desired interface
	priceReader := ccipinternal.NewPriceRegistryGRPCClient(priceReaderConn)

	return priceReader, nil
}

// NewTokenDataReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewTokenDataReader(ctx context.Context, tokenAddress cciptypes.Address) (cciptypes.TokenDataReader, error) {
	req := ccippb.NewTokenDataRequest{Address: string(tokenAddress)}
	resp, err := e.grpcClient.NewTokenDataReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	tokenDataConn, err := e.BrokerExt.Dial(uint32(resp.TokenDataReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup token data reader service at %d: %w", resp.TokenDataReaderServiceId, err)
	}
	// need to wrap grpc tokenDataReader into the desired interface
	tokenDataReader := ccipinternal.NewTokenDataReaderGRPCClient(tokenDataConn)

	return tokenDataReader, nil
}

// NewTokenPoolBatchedReader implements types.CCIPExecProvider.
func (e *execProviderClient) NewTokenPoolBatchedReader(ctx context.Context) (cciptypes.TokenPoolBatchedReader, error) {
	resp, err := e.grpcClient.NewTokenPoolBatchedReader(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	tokenPoolConn, err := e.BrokerExt.Dial(uint32(resp.TokenPoolBatchedReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup token poll batched reader service at %d: %w", resp.TokenPoolBatchedReaderServiceId, err)
	}
	tokenPool := ccipinternal.NewTokenPoolBatchedReaderGRPCClient(tokenPoolConn)
	return tokenPool, nil
}

// SourceNativeToken implements types.CCIPExecProvider.
func (e *execProviderClient) SourceNativeToken(ctx context.Context) (cciptypes.Address, error) {
	panic("BCF-3109")
}

// execProviderServer is a server that wraps the custom methods of the [types.CCIPExecProvider]
// this is necessary because those method create new resources that need to be served by the broker
// when we are running in legacy mode.
type execProviderServer struct {
	ccippb.UnimplementedExecutionCustomHandlersServer
	// BCF-3061 this has to be a shared pointer to the same impl as the execProviderClient
	*net.BrokerExt
	impl types.CCIPExecProvider

	deps net.Resources
}

// Close implements ccippb.ExecutionCustomHandlersServer.
func (e *execProviderServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.impl.Close()
}

// NewCommitStoreReader implements ccippb.ExecutionCustomHandlersServer.
func (e *execProviderServer) NewCommitStoreReader(ctx context.Context, req *ccippb.NewCommitStoreReaderRequest) (*ccippb.NewCommitStoreReaderResponse, error) {
	reader, err := e.impl.NewCommitStoreReader(context.Background(), ccip.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	commitStoreHandler, err := ccipinternal.NewCommitStoreReaderGRPCServer(reader, e.BrokerExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp reader grpc server: %w", err)
	}
	// the id is handle to the broker, we will need it on the other sider to dial the resource
	commitStoreID, csResource, err := e.ServeNew("OffRampReader", func(s *grpc.Server) {
		ccippb.RegisterCommitStoreReaderServer(s, commitStoreHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the offRamp is closed. See comment in NewPriceRegistryReader for more details
	commitStoreHandler.AddDep(csResource)
	return &ccippb.NewCommitStoreReaderResponse{CommitStoreReaderServiceId: int32(commitStoreID)}, nil
}

var _ ccippb.ExecutionCustomHandlersServer = (*execProviderServer)(nil)

func newExecProviderServer(impl types.CCIPExecProvider, brokerExt *net.BrokerExt) *execProviderServer {
	return &execProviderServer{impl: impl, BrokerExt: brokerExt}
}

func (e *execProviderServer) NewOffRampReader(ctx context.Context, req *ccippb.NewOffRampReaderRequest) (*ccippb.NewOffRampReaderResponse, error) {
	reader, err := e.impl.NewOffRampReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	offRampHandler, err := ccipinternal.NewOffRampReaderGRPCServer(reader, e.BrokerExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp reader grpc server: %w", err)
	}
	// the id is handle to the broker, we will need it on the other sider to dial the resource
	offRampID, offRampResource, err := e.ServeNew("OffRampReader", func(s *grpc.Server) {
		ccippb.RegisterOffRampReaderServer(s, offRampHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the offRamp is closed. See comment in NewPriceRegistryReader for more details
	offRampHandler.AddDep(offRampResource)
	return &ccippb.NewOffRampReaderResponse{OfframpReaderServiceId: int32(offRampID)}, nil
}

func (e *execProviderServer) NewOnRampReader(ctx context.Context, req *ccippb.NewOnRampReaderRequest) (*ccippb.NewOnRampReaderResponse, error) {
	reader, err := e.impl.NewOnRampReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	srv := ccipinternal.NewOnRampReaderServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	onRampID, onRampResource, err := e.ServeNew("OnRampReader", func(s *grpc.Server) {
		ccippb.RegisterOnRampReaderServer(s, srv)
	})
	if err != nil {
		return nil, err
	}
	// TODO BCF-3067 LEAKS!!!
	// this dependency needs to be closed when the onramp reader is closed, which
	// should happen when the calling reporting plugin is closed/goes out of scope
	e.deps.Add(onRampResource)
	return &ccippb.NewOnRampReaderResponse{OnrampReaderServiceId: int32(onRampID)}, nil
}

func (e *execProviderServer) NewPriceRegistryReader(ctx context.Context, req *ccippb.NewPriceRegistryReaderRequest) (*ccippb.NewPriceRegistryReaderResponse, error) {
	reader, err := e.impl.NewPriceRegistryReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	priceRegistryHandler := ccipinternal.NewPriceRegistryGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	priceReaderID, spawnedServer, err := e.ServeNew("PriceRegistryReader", func(s *grpc.Server) {
		ccippb.RegisterPriceRegistryReaderServer(s, priceRegistryHandler)
	})
	if err != nil {
		return nil, err
	}
	// There is a chicken-and-egg problem here. Our broker is responsible for spawning the grpc server.
	// that server needs to be shutdown when the priceRegistry is closed. We don't have a handle to the
	// grpc server until we after we have constructed the priceRegistry, so we can't configure the shutdown
	// handler up front.
	priceRegistryHandler.AddDep(spawnedServer)
	return &ccippb.NewPriceRegistryReaderResponse{PriceRegistryReaderServiceId: int32(priceReaderID)}, nil
}

func (e *execProviderServer) NewTokenDataReader(ctx context.Context, req *ccippb.NewTokenDataRequest) (*ccippb.NewTokenDataResponse, error) {
	reader, err := e.impl.NewTokenDataReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	tokenDataHandler := ccipinternal.NewTokenDataReaderGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	tokeDataReaderID, spawnedServer, err := e.ServeNew("TokenDataReader", func(s *grpc.Server) {
		ccippb.RegisterTokenDataReaderServer(s, tokenDataHandler)
	})
	if err != nil {
		return nil, err
	}

	tokenDataHandler.AddDep(spawnedServer)
	return &ccippb.NewTokenDataResponse{TokenDataReaderServiceId: int32(tokeDataReaderID)}, nil
}

func (e *execProviderServer) NewTokenPoolBatchedReader(ctx context.Context, _ *emptypb.Empty) (*ccippb.NewTokenPoolBatchedReaderResponse, error) {
	reader, err := e.impl.NewTokenPoolBatchedReader(ctx)
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	tokenPoolHandler := ccipinternal.NewTokenPoolBatchedReaderGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	tokenPoolID, spawnedServer, err := e.ServeNew("TokenPoolBatchedReader", func(s *grpc.Server) {
		ccippb.RegisterTokenPoolBatcherReaderServer(s, tokenPoolHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the tokenPool is closed. See comment in NewPriceRegistryReader for more details
	tokenPoolHandler.AddDep(spawnedServer)
	return &ccippb.NewTokenPoolBatchedReaderResponse{TokenPoolBatchedReaderServiceId: int32(tokenPoolID)}, nil
}
