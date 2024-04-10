package ccip

import (
	"context"
	"io"
	"math/big"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// PriceGetterGRPCClient implements [cciptypes.PriceGetter] by wrapping a
// [ccippb.PriceGetterGRPCClient] grpc client.
// This client is used by a ReportingPlugin to call the PriceGetter service, which
// is hosted by the relayer
type PriceGetterGRPCClient struct {
	grpc ccippb.PriceGetterClient
}

func NewPriceGetterGRPCClient(cc grpc.ClientConnInterface) *PriceGetterGRPCClient {
	return &PriceGetterGRPCClient{grpc: ccippb.NewPriceGetterClient(cc)}
}

// PriceGetterGRPCServer implements [ccippb.PriceGetterServer] by wrapping a
// [cciptypes.PriceGetter] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [PriceGetterGRPCClient]
type PriceGetterGRPCServer struct {
	ccippb.UnimplementedPriceGetterServer

	impl cciptypes.PriceGetter
	deps []io.Closer
}

func NewPriceGetterGRPCServer(impl cciptypes.PriceGetter) *PriceGetterGRPCServer {
	return &PriceGetterGRPCServer{impl: impl, deps: []io.Closer{impl}}
}

// ensure the types are satisfied
var _ cciptypes.PriceGetter = (*PriceGetterGRPCClient)(nil)
var _ ccippb.PriceGetterServer = (*PriceGetterGRPCServer)(nil)

// FilterConfiguredTokens implements ccip.PriceGetter.
func (p *PriceGetterGRPCClient) FilterConfiguredTokens(ctx context.Context, tokens []cciptypes.Address) (configured []cciptypes.Address, unconfigured []cciptypes.Address, err error) {
	configured = []cciptypes.Address{}
	unconfigured = []cciptypes.Address{}

	// convert the format
	requestedTokens := make([]string, len(tokens))
	for i, t := range tokens {
		requestedTokens[i] = string(t)
	}

	resp, err := p.grpc.FilterConfiguredTokens(ctx, &ccippb.FilterConfiguredTokensRequest{Tokens: requestedTokens})
	if err != nil {
		return nil, nil, err
	}
	for _, confTk := range resp.Configured {
		if len(confTk) > 0 {
			configured = append(configured, ccip.Address(confTk))
		}
	}
	for _, unconfTk := range resp.Unconfigured {
		if len(unconfTk) > 0 {
			unconfigured = append(unconfigured, ccip.Address(unconfTk))
		}
	}

	return configured, unconfigured, nil
}

// TokenPricesUSD implements ccip.PriceGetter.
func (p *PriceGetterGRPCClient) TokenPricesUSD(ctx context.Context, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	// convert the format
	requestedTokens := make([]string, len(tokens))
	for i, t := range tokens {
		requestedTokens[i] = string(t)
	}

	resp, err := p.grpc.TokenPricesUSD(ctx, &ccippb.TokenPricesRequest{Tokens: requestedTokens})
	if err != nil {
		return nil, err
	}
	prices := make(map[cciptypes.Address]*big.Int, len(resp.Prices))
	for addr, p := range resp.Prices {
		prices[ccip.Address(addr)] = p.Int()
	}
	return prices, nil
}

func (p *PriceGetterGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), p.grpc)
}

// FilterConfiguredTokens implements ccippb.PriceGetterServer.
func (p *PriceGetterGRPCServer) FilterConfiguredTokens(ctx context.Context, req *ccippb.FilterConfiguredTokensRequest) (*ccippb.FilterConfiguredTokensResponse, error) {
	tokenAddresses := make([]cciptypes.Address, len(req.Tokens))
	for i, t := range req.Tokens {
		tokenAddresses[i] = cciptypes.Address(t)
	}

	configuredTokens, unconfiguredTokens, err := p.impl.FilterConfiguredTokens(ctx, tokenAddresses)
	if err != nil {
		return nil, err
	}

	convertedConfTks := make([]string, len(configuredTokens))
	for _, confTk := range configuredTokens {
		convertedConfTks = append(convertedConfTks, string(confTk))
	}
	convertedUnconfTks := make([]string, len(unconfiguredTokens))
	for _, unconfTk := range unconfiguredTokens {
		convertedUnconfTks = append(convertedUnconfTks, string(unconfTk))
	}

	return &ccippb.FilterConfiguredTokensResponse{Configured: convertedConfTks, Unconfigured: convertedUnconfTks}, nil
}

// TokenPricesUSD implements ccippb.PriceGetterServer.
func (p *PriceGetterGRPCServer) TokenPricesUSD(ctx context.Context, req *ccippb.TokenPricesRequest) (*ccippb.TokenPricesResponse, error) {
	tokenAddresses := make([]cciptypes.Address, len(req.Tokens))
	for i, t := range req.Tokens {
		tokenAddresses[i] = cciptypes.Address(t)
	}

	prices, err := p.impl.TokenPricesUSD(ctx, tokenAddresses)
	if err != nil {
		return nil, err
	}

	convertedPrices := make(map[string]*pb.BigInt, len(prices))
	for addr, p := range prices {
		convertedPrices[string(addr)] = pb.NewBigIntFromInt(p)
	}
	return &ccippb.TokenPricesResponse{Prices: convertedPrices}, nil
}

func (p *PriceGetterGRPCServer) AddDep(closer io.Closer) *PriceGetterGRPCServer {
	p.deps = append(p.deps, closer)
	return p
}

func (p *PriceGetterGRPCServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(p.deps).Close()
}
