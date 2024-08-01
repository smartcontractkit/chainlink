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

func RegisterExecutionProviderServices(s *grpc.Server, provider types.CCIPExecProvider, brokerExt *net.BrokerExt) {
	ocr2.RegisterPluginProviderServices(s, provider)
	// register the handler for the custom methods of the provider eg NewOffRampReader
	ccippb.RegisterExecutionCustomHandlersServer(s, NewExecProviderServer(provider, brokerExt))
}

var (
	_ types.CCIPExecProvider  = (*ExecProviderClient)(nil)
	_ goplugin.GRPCClientConn = (*ExecProviderClient)(nil)
)

type ExecProviderClient struct {
	*ocr2.PluginProviderClient

	// must be shared with the server
	*net.BrokerExt
	grpcClient ccippb.ExecutionCustomHandlersClient
}

func NewExecProviderClient(b *net.BrokerExt, conn grpc.ClientConnInterface) *ExecProviderClient {
	pluginProviderClient := ocr2.NewPluginProviderClient(b, conn)
	grpc := ccippb.NewExecutionCustomHandlersClient(conn)
	return &ExecProviderClient{
		PluginProviderClient: pluginProviderClient,
		BrokerExt:            b,
		grpcClient:           grpc,
	}
}

// GetTransactionStatus implements types.CCIPExecProvider.
func (e *ExecProviderClient) GetTransactionStatus(ctx context.Context, transactionID string) (types.TransactionStatus, error) {
	// unlike the other methods, this one does not create a new resource, so we do not
	// need the broker to serve it. we can just call the grpc method directly.
	resp, err := e.grpcClient.GetTransactionStatus(ctx, &ccippb.GetTransactionStatusRequest{TransactionId: transactionID})
	if err != nil {
		return 0, err
	}
	return types.TransactionStatus(resp.TransactionStatus), nil
}

// NewCommitStoreReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewCommitStoreReader(ctx context.Context, addr cciptypes.Address) (cciptypes.CommitStoreReader, error) {
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
	commitStore := NewCommitStoreReaderGRPCClient(e.BrokerExt, commitStoreConn)

	return commitStore, nil
}

// NewOffRampReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewOffRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OffRampReader, error) {
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
	offRamp := NewOffRampReaderGRPCClient(e.BrokerExt, offRampConn)

	return offRamp, nil
}

// NewOnRampReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewOnRampReader(ctx context.Context, addr cciptypes.Address, srcChainSelector uint64, dstChainSelector uint64) (cciptypes.OnRampReader, error) {
	req := ccippb.NewOnRampReaderRequest{Address: string(addr), SourceChainSelector: srcChainSelector, DestChainSelector: dstChainSelector}

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
	onRamp := NewOnRampReaderGRPCClient(onRampConn)

	// how to convert resp to cciptypes.OnRampReader? i have an id and need to hydrate that into an instance of OnRampReader
	return onRamp, nil
}

// NewPriceRegistryReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error) {
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
	priceReader := NewPriceRegistryGRPCClient(priceReaderConn)

	return priceReader, nil
}

// NewTokenDataReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewTokenDataReader(ctx context.Context, tokenAddress cciptypes.Address) (cciptypes.TokenDataReader, error) {
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
	tokenDataReader := NewTokenDataReaderGRPCClient(tokenDataConn)

	return tokenDataReader, nil
}

// NewTokenPoolBatchedReader implements types.CCIPExecProvider.
func (e *ExecProviderClient) NewTokenPoolBatchedReader(ctx context.Context, offRampAddress cciptypes.Address, srcChainSelector uint64) (cciptypes.TokenPoolBatchedReader, error) {
	req := ccippb.NewTokenPoolBatchedReaderRequest{Address: string(offRampAddress), SourceChainSelector: srcChainSelector}
	resp, err := e.grpcClient.NewTokenPoolBatchedReader(ctx, &req)
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: make this work for proxied relayer
	tokenPoolConn, err := e.BrokerExt.Dial(uint32(resp.TokenPoolBatchedReaderServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup token poll batched reader service at %d: %w", resp.TokenPoolBatchedReaderServiceId, err)
	}
	tokenPool := NewTokenPoolBatchedReaderGRPCClient(tokenPoolConn)
	return tokenPool, nil
}

// SourceNativeToken implements types.CCIPExecProvider.
func (e *ExecProviderClient) SourceNativeToken(ctx context.Context, addr cciptypes.Address) (cciptypes.Address, error) {
	// unlike the other methods, this one does not create a new resource, so we do not
	// need the broker to serve it. we can just call the grpc method directly.
	resp, err := e.grpcClient.SourceNativeToken(ctx, &ccippb.SourceNativeTokenRequest{SourceRouterAddress: string(addr)})
	if err != nil {
		return "", err
	}
	return cciptypes.Address(resp.NativeTokenAddress), nil
}

// Close implements types.CCIPExecProvider.
func (e *ExecProviderClient) Close() error {
	return shutdownGRPCServer(context.Background(), e.grpcClient)
}

// ExecProviderServer is a server that wraps the custom methods of the [types.CCIPExecProvider]
// this is necessary because those method create new resources that need to be served by the broker
// when we are running in legacy mode.
type ExecProviderServer struct {
	ccippb.UnimplementedExecutionCustomHandlersServer
	// BCF-3061 this has to be a shared pointer to the same impl as the execProviderClient
	*net.BrokerExt
	impl types.CCIPExecProvider
}

// Close implements ccippb.ExecutionCustomHandlersServer.
func (e *ExecProviderServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.impl.Close()
}

// NewCommitStoreReader implements ccippb.ExecutionCustomHandlersServer.
func (e *ExecProviderServer) NewCommitStoreReader(ctx context.Context, req *ccippb.NewCommitStoreReaderRequest) (*ccippb.NewCommitStoreReaderResponse, error) {
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
	commitStoreID, csResource, err := e.ServeNew("CommitStoreReader", func(s *grpc.Server) {
		ccippb.RegisterCommitStoreReaderServer(s, commitStoreHandler)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the offRamp is closed. See comment in NewPriceRegistryReader for more details
	commitStoreHandler.AddDep(csResource)
	return &ccippb.NewCommitStoreReaderResponse{CommitStoreReaderServiceId: int32(commitStoreID)}, nil
}

var _ ccippb.ExecutionCustomHandlersServer = (*ExecProviderServer)(nil)

func NewExecProviderServer(impl types.CCIPExecProvider, brokerExt *net.BrokerExt) *ExecProviderServer {
	return &ExecProviderServer{impl: impl, BrokerExt: brokerExt}
}

func (e *ExecProviderServer) GetTransactionStatus(ctx context.Context, req *ccippb.GetTransactionStatusRequest) (*ccippb.GetTransactionStatusResponse, error) {
	ts, err := e.impl.GetTransactionStatus(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetTransactionStatusResponse{TransactionStatus: int32(ts)}, nil
}

func (e *ExecProviderServer) NewOffRampReader(ctx context.Context, req *ccippb.NewOffRampReaderRequest) (*ccippb.NewOffRampReaderResponse, error) {
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

func (e *ExecProviderServer) NewOnRampReader(ctx context.Context, req *ccippb.NewOnRampReaderRequest) (*ccippb.NewOnRampReaderResponse, error) {
	reader, err := e.impl.NewOnRampReader(ctx, cciptypes.Address(req.Address), req.SourceChainSelector, req.DestChainSelector)
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	srv := NewOnRampReaderGRPCServer(reader)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	onRampID, onRampResource, err := e.ServeNew("OnRampReader", func(s *grpc.Server) {
		ccippb.RegisterOnRampReaderServer(s, srv)
	})
	if err != nil {
		return nil, err
	}
	// ensure the grpc server is closed when the onRamp is closed. See comment in NewPriceRegistryReader for more details
	srv.AddDep(onRampResource)
	return &ccippb.NewOnRampReaderResponse{OnrampReaderServiceId: int32(onRampID)}, nil
}

func (e *ExecProviderServer) NewPriceRegistryReader(ctx context.Context, req *ccippb.NewPriceRegistryReaderRequest) (*ccippb.NewPriceRegistryReaderResponse, error) {
	reader, err := e.impl.NewPriceRegistryReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	priceRegistryHandler := NewPriceRegistryGRPCServer(reader)
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

func (e *ExecProviderServer) NewTokenDataReader(ctx context.Context, req *ccippb.NewTokenDataRequest) (*ccippb.NewTokenDataResponse, error) {
	reader, err := e.impl.NewTokenDataReader(ctx, cciptypes.Address(req.Address))
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	tokenDataHandler := NewTokenDataReaderGRPCServer(reader)
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

func (e *ExecProviderServer) NewTokenPoolBatchedReader(ctx context.Context, req *ccippb.NewTokenPoolBatchedReaderRequest) (*ccippb.NewTokenPoolBatchedReaderResponse, error) {
	reader, err := e.impl.NewTokenPoolBatchedReader(ctx, cciptypes.Address(req.Address), req.SourceChainSelector)
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	tokenPoolHandler := NewTokenPoolBatchedReaderGRPCServer(reader)
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

func (e *ExecProviderServer) SourceNativeToken(ctx context.Context, req *ccippb.SourceNativeTokenRequest) (*ccippb.SourceNativeTokenResponse, error) {
	addr, err := e.impl.SourceNativeToken(ctx, cciptypes.Address(req.SourceRouterAddress))
	if err != nil {
		return nil, err
	}
	return &ccippb.SourceNativeTokenResponse{NativeTokenAddress: string(addr)}, nil
}
