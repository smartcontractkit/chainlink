package ccip

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

func RegisterCommitProviderServices(s *grpc.Server, provider types.CCIPCommitProvider, brokerExt *net.BrokerExt) {
	ocr2.RegisterPluginProviderServices(s, provider)
	// register the handler for the custom methods of the provider eg NewOffRampReader
	ccippb.RegisterCommitCustomHandlersServer(s, NewCommitProviderServer(provider, brokerExt))
}

var (
	_ types.CCIPCommitProvider = (*CommitProviderClient)(nil)
	_ goplugin.GRPCClientConn  = (*CommitProviderClient)(nil)
)

type CommitProviderClient struct {
	*ocr2.PluginProviderClient

	// must be shared with the server
	*net.BrokerExt
	grpcClient ccippb.CommitCustomHandlersClient
}

func NewCommitProviderClient(b *net.BrokerExt, conn grpc.ClientConnInterface) *CommitProviderClient {
	pluginProviderClient := ocr2.NewPluginProviderClient(b, conn)
	grpc := ccippb.NewCommitCustomHandlersClient(conn)
	return &CommitProviderClient{
		PluginProviderClient: pluginProviderClient,
		BrokerExt:            b,
		grpcClient:           grpc,
	}
}

// NewCommitStoreReader implements types.CCIPCommitProvider.
func (e *CommitProviderClient) NewCommitStoreReader(ctx context.Context, addr cciptypes.Address) (cciptypes.CommitStoreReader, error) {
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
	return NewCommitStoreReaderGRPCClient(e.BrokerExt, commitStoreConn), nil
}

// NewOffRampReader implements types.CCIPCommitProvider.
func (e *CommitProviderClient) NewOffRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OffRampReader, error) {
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
	return NewOffRampReaderGRPCClient(e.BrokerExt, offRampConn), nil
}

// NewOnRampReader implements types.CCIPCommitProvider.
func (e *CommitProviderClient) NewOnRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OnRampReader, error) {
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
	return NewOnRampReaderGRPCClient(onRampConn), nil
}

// NewPriceGetter implements types.CCIPCommitProvider.
func (e *CommitProviderClient) NewPriceGetter(ctx context.Context) (cciptypes.PriceGetter, error) {
	resp, err := e.grpcClient.NewPriceGetter(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	priceGetterConn, err := e.BrokerExt.Dial(uint32(resp.PriceGetterServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup price getter service at %d: %w", resp.PriceGetterServiceId, err)
	}
	return NewPriceGetterGRPCClient(priceGetterConn), nil
}

// NewPriceRegistryReader implements types.CCIPCommitProvider.
func (e *CommitProviderClient) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error) {
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
	return NewPriceRegistryGRPCClient(priceReaderConn), nil
}

// SourceNativeToken implements types.CCIPCommitProvider.
func (e *CommitProviderClient) SourceNativeToken(ctx context.Context) (cciptypes.Address, error) {
	// unlike the other methods, this one does not create a new resource, so we do not
	// need the broker to serve it. we can just call the grpc method directly.
	resp, err := e.grpcClient.SourceNativeToken(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return cciptypes.Address(resp.NativeTokenAddress), nil
}

func (e *CommitProviderClient) Close() error {
	_, err := e.grpcClient.Close(context.Background(), &emptypb.Empty{})
	return err
}

// CommitProviderServer is a server that wraps the custom methods of the [types.CCIPCommitProvider]
// this is necessary because those method create new resources that need to be served by the broker
// when we are running in legacy mode.
type CommitProviderServer struct {
	ccippb.UnimplementedCommitCustomHandlersServer
	// BCF-3061 this has to be a shared pointer to the same impl as the execProviderClient
	*net.BrokerExt
	impl types.CCIPCommitProvider
}

// Close implements ccippb.CommitCustomHandlersServer.
func (e *CommitProviderServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.impl.Close()
}

// NewCommitStoreReader implements ccippb.CommitCustomHandlersServer.
func (e *CommitProviderServer) NewCommitStoreReader(ctx context.Context, req *ccippb.NewCommitStoreReaderRequest) (*ccippb.NewCommitStoreReaderResponse, error) {
	reader, err := e.impl.NewCommitStoreReader(context.Background(), ccip.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	commitStoreHandler, err := NewCommitStoreReaderGRPCServer(reader, e.BrokerExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp reader grpc server: %w", err)
	}
	// the id is handle to the broker, we will need it on the other sider to dial the resource
	commitStoreID, csResource, err := e.ServeNew(e.formatSubserviceName("CommitStoreReader"), func(s *grpc.Server) {
		ccippb.RegisterCommitStoreReaderServer(s, commitStoreHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the offRamp is closed. See comment in NewPriceRegistryReader for more details
	commitStoreHandler.AddDep(csResource)
	return &ccippb.NewCommitStoreReaderResponse{CommitStoreReaderServiceId: int32(commitStoreID)}, nil
}

var _ ccippb.CommitCustomHandlersServer = (*CommitProviderServer)(nil)

func NewCommitProviderServer(impl types.CCIPCommitProvider, brokerExt *net.BrokerExt) *CommitProviderServer {
	return &CommitProviderServer{impl: impl, BrokerExt: brokerExt}
}

func (e *CommitProviderServer) NewOffRampReader(ctx context.Context, req *ccippb.NewOffRampReaderRequest) (*ccippb.NewOffRampReaderResponse, error) {
	reader, err := e.impl.NewOffRampReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	offRampHandler, err := NewOffRampReaderGRPCServer(reader, e.BrokerExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp reader grpc server: %w", err)
	}
	// the id is handle to the broker, we will need it on the other sider to dial the resource
	offRampID, offRampResource, err := e.ServeNew(e.formatSubserviceName("OffRampReader"), func(s *grpc.Server) {
		ccippb.RegisterOffRampReaderServer(s, offRampHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the offRamp is closed. See comment in NewPriceRegistryReader for more details
	offRampHandler.AddDep(offRampResource)
	return &ccippb.NewOffRampReaderResponse{OfframpReaderServiceId: int32(offRampID)}, nil
}

func (e *CommitProviderServer) NewOnRampReader(ctx context.Context, req *ccippb.NewOnRampReaderRequest) (*ccippb.NewOnRampReaderResponse, error) {
	reader, err := e.impl.NewOnRampReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	srv := NewOnRampReaderGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	onRampID, onRampResource, err := e.ServeNew(e.formatSubserviceName("OnRampReader"), func(s *grpc.Server) {
		ccippb.RegisterOnRampReaderServer(s, srv)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the onRamp is closed. See comment in NewPriceRegistryReader for more details
	srv.AddDep(onRampResource)
	return &ccippb.NewOnRampReaderResponse{OnrampReaderServiceId: int32(onRampID)}, nil
}

func (e *CommitProviderServer) NewPriceGetter(ctx context.Context, _ *emptypb.Empty) (*ccippb.NewPriceGetterResponse, error) {
	priceGetter, err := e.impl.NewPriceGetter(ctx)
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	priceGetterHandler := NewPriceGetterGRPCServer(priceGetter)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	priceGetterID, spawnedServer, err := e.ServeNew(e.formatSubserviceName("PriceGetter"), func(s *grpc.Server) {
		ccippb.RegisterPriceGetterServer(s, priceGetterHandler)
	})
	if err != nil {
		return nil, err
	}
	// There is a chicken-and-egg problem here. Our broker is responsible for spawning the grpc server.
	// that server needs to be shutdown when the priceGetter is closed. We don't have a handle to the
	// grpc server until we after we have constructed the priceGetter, so we can't configure the shutdown
	// handler up front.
	priceGetterHandler.AddDep(spawnedServer)
	return &ccippb.NewPriceGetterResponse{PriceGetterServiceId: int32(priceGetterID)}, nil
}

func (e *CommitProviderServer) NewPriceRegistryReader(ctx context.Context, req *ccippb.NewPriceRegistryReaderRequest) (*ccippb.NewPriceRegistryReaderResponse, error) {
	reader, err := e.impl.NewPriceRegistryReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	priceRegistryHandler := NewPriceRegistryGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	priceReaderID, spawnedServer, err := e.ServeNew(e.formatSubserviceName("PriceRegistryReader"), func(s *grpc.Server) {
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

func (e *CommitProviderServer) SourceNativeToken(ctx context.Context, _ *emptypb.Empty) (*ccippb.SourceNativeTokenResponse, error) {
	addr, err := e.impl.SourceNativeToken(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.SourceNativeTokenResponse{NativeTokenAddress: string(addr)}, nil
}

func (e *CommitProviderServer) formatSubserviceName(serviceName string) string {
	return fmt.Sprintf("CommitProvider.%s", serviceName)
}
