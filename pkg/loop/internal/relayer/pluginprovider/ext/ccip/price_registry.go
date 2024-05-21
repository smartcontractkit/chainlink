package ccip

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// PriceRegistryGRPCClient implements [cciptypes.PriceRegistryReader] by wrapping a
// [ccippb.PriceRegistryReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call the PriceRegistryReader service, which
// is hosted by the relayer
type PriceRegistryGRPCClient struct {
	grpc ccippb.PriceRegistryReaderClient
	conn grpc.ClientConnInterface
}

func NewPriceRegistryGRPCClient(cc grpc.ClientConnInterface) *PriceRegistryGRPCClient {
	return &PriceRegistryGRPCClient{grpc: ccippb.NewPriceRegistryReaderClient(cc)}
}

// PriceRegistryGRPCServer implements [ccippb.PriceRegistryReaderServer] by wrapping a
// [cciptypes.PriceRegistryReader] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [PriceRegistryGRPCClient]
type PriceRegistryGRPCServer struct {
	ccippb.UnimplementedPriceRegistryReaderServer

	impl cciptypes.PriceRegistryReader
	deps []io.Closer
}

func NewPriceRegistryGRPCServer(impl cciptypes.PriceRegistryReader) *PriceRegistryGRPCServer {
	return &PriceRegistryGRPCServer{impl: impl, deps: []io.Closer{impl}}
}

// ensure the types are satisfied
var _ cciptypes.PriceRegistryReader = (*PriceRegistryGRPCClient)(nil)
var _ ccippb.PriceRegistryReaderServer = (*PriceRegistryGRPCServer)(nil)

// PriceRegistryGRPCClient implementation

// Address implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) Address(ctx context.Context) (cciptypes.Address, error) {
	resp, err := p.grpc.GetAddress(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return cciptypes.Address(resp.Address), nil
}

func (p *PriceRegistryGRPCClient) ClientConn() grpc.ClientConnInterface {
	return p.conn
}

// Close implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), p.grpc)
}

// GetFeeTokens implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetFeeTokens(ctx context.Context) ([]cciptypes.Address, error) {
	resp, err := p.grpc.GetFeeTokens(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return cciptypes.MakeAddresses(resp.FeeTokenAddresses), nil
}

// GetGasPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confirmations int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	req := &ccippb.GetGasPriceUpdatesCreatedAfterRequest{
		ChainSelector: chainSelector,
		CreatedAfter:  timestamppb.New(ts),
		Confirmations: uint64(confirmations),
	}
	resp, err := p.grpc.GetGasPriceUpdatesCreatedAfter(ctx, req)
	if err != nil {
		return nil, err
	}
	return gasPriceUpdateWithTxMetaSlice(resp.GasPriceUpdates), nil
}

// GetAllGasPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetAllGasPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	req := &ccippb.GetAllGasPriceUpdatesCreatedAfterRequest{
		CreatedAfter:  timestamppb.New(ts),
		Confirmations: uint64(confirmations),
	}
	resp, err := p.grpc.GetAllGasPriceUpdatesCreatedAfter(ctx, req)
	if err != nil {
		return nil, err
	}
	return gasPriceUpdateWithTxMetaSlice(resp.GasPriceUpdates), nil
}

// GetTokenPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.TokenPriceUpdateWithTxMeta, error) {
	req := &ccippb.GetTokenPriceUpdatesCreatedAfterRequest{
		CreatedAfter:  timestamppb.New(ts),
		Confirmations: uint64(confirmations),
	}
	resp, err := p.grpc.GetTokenPriceUpdatesCreatedAfter(ctx, req)
	if err != nil {
		return nil, err
	}
	return tokenPriceUpdateWithTxMetaSlice(resp.TokenPriceUpdates), nil
}

// GetTokenPrices implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetTokenPrices(ctx context.Context, wantedTokens []cciptypes.Address) ([]cciptypes.TokenPriceUpdate, error) {
	req := &ccippb.GetTokenPricesRequest{
		TokenAddresses: cciptypes.Addresses(wantedTokens).Strings(),
	}
	resp, err := p.grpc.GetTokenPrices(ctx, req)
	if err != nil {
		return nil, err
	}

	return tokenPriceUpdateSlice(resp.TokenPrices), nil
}

// GetTokensDecimals implements ccip.PriceRegistryReader.
func (p *PriceRegistryGRPCClient) GetTokensDecimals(ctx context.Context, tokenAddresses []cciptypes.Address) ([]uint8, error) {
	req := &ccippb.GetTokensDecimalsRequest{
		TokenAddresses: cciptypes.Addresses(tokenAddresses).Strings(),
	}
	resp, err := p.grpc.GetTokensDecimals(ctx, req)
	if err != nil {
		return nil, err
	}

	return decimals(resp.Decimals)
}

//
// PriceRegistryGRPCServer implementation
//

// Close implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) Close(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(p.deps).Close()
}

// GetAddress implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetAddress(ctx context.Context, req *emptypb.Empty) (*ccippb.GetPriceRegistryAddressResponse, error) {
	addr, err := p.impl.Address(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetPriceRegistryAddressResponse{Address: string(addr)}, nil
}

// GetFeeTokens implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetFeeTokens(ctx context.Context, req *emptypb.Empty) (*ccippb.GetFeeTokensResponse, error) {
	addrs, err := p.impl.GetFeeTokens(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetFeeTokensResponse{FeeTokenAddresses: cciptypes.Addresses(addrs).Strings()}, nil
}

// GetGasPriceUpdatesCreatedAfter implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetGasPriceUpdatesCreatedAfter(ctx context.Context, req *ccippb.GetGasPriceUpdatesCreatedAfterRequest) (*ccippb.GetGasPriceUpdatesCreatedAfterResponse, error) {
	updates, err := p.impl.GetGasPriceUpdatesCreatedAfter(ctx, req.ChainSelector, req.CreatedAfter.AsTime(), int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	return &ccippb.GetGasPriceUpdatesCreatedAfterResponse{GasPriceUpdates: gasPriceUpdateWithTxMetaSlicePB(updates)}, nil
}

// GetAllGasPriceUpdatesCreatedAfter implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetAllGasPriceUpdatesCreatedAfter(ctx context.Context, req *ccippb.GetAllGasPriceUpdatesCreatedAfterRequest) (*ccippb.GetAllGasPriceUpdatesCreatedAfterResponse, error) {
	updates, err := p.impl.GetAllGasPriceUpdatesCreatedAfter(ctx, req.CreatedAfter.AsTime(), int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	return &ccippb.GetAllGasPriceUpdatesCreatedAfterResponse{GasPriceUpdates: gasPriceUpdateWithTxMetaSlicePB(updates)}, nil
}

// GetTokenPriceUpdatesCreatedAfter implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, req *ccippb.GetTokenPriceUpdatesCreatedAfterRequest) (*ccippb.GetTokenPriceUpdatesCreatedAfterResponse, error) {
	updates, err := p.impl.GetTokenPriceUpdatesCreatedAfter(ctx, req.CreatedAfter.AsTime(), int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	return &ccippb.GetTokenPriceUpdatesCreatedAfterResponse{TokenPriceUpdates: tokenPriceUpdateWithTxMetaSlicePB(updates)}, nil
}

// GetTokenPrices implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetTokenPrices(ctx context.Context, req *ccippb.GetTokenPricesRequest) (*ccippb.GetTokenPricesResponse, error) {
	prices, err := p.impl.GetTokenPrices(ctx, cciptypes.MakeAddresses(req.TokenAddresses))
	if err != nil {
		return nil, err
	}
	return &ccippb.GetTokenPricesResponse{TokenPrices: tokenPriceUpdateSlicePB(prices)}, nil
}

// GetTokensDecimals implements ccippb.PriceRegistryReaderServer.
func (p *PriceRegistryGRPCServer) GetTokensDecimals(ctx context.Context, req *ccippb.GetTokensDecimalsRequest) (*ccippb.GetTokensDecimalsResponse, error) {
	decimals, err := p.impl.GetTokensDecimals(ctx, cciptypes.MakeAddresses(req.TokenAddresses))
	if err != nil {
		return nil, err
	}
	return decimalsPB(decimals), nil
}

// AddDep adds a dependency to the server that will be closed when the server is closed
func (p *PriceRegistryGRPCServer) AddDep(dep io.Closer) *PriceRegistryGRPCServer {
	p.deps = append(p.deps, dep)
	return p
}

func gasPriceUpdateWithTxMetaSlice(in []*ccippb.GasPriceUpdateWithTxMeta) []cciptypes.GasPriceUpdateWithTxMeta {
	out := make([]cciptypes.GasPriceUpdateWithTxMeta, len(in))
	for i, v := range in {
		out[i] = gasPriceUpdateWithTxMeta(v)
	}
	return out
}

func gasPriceUpdateWithTxMeta(in *ccippb.GasPriceUpdateWithTxMeta) cciptypes.GasPriceUpdateWithTxMeta {
	return cciptypes.GasPriceUpdateWithTxMeta{
		TxMeta:         txMeta(in.TxMeta),
		GasPriceUpdate: gasPriceUpdate(in.GasPriceUpdate),
	}
}

func gasPriceUpdate(in *ccippb.GasPriceUpdate) cciptypes.GasPriceUpdate {
	return cciptypes.GasPriceUpdate{
		GasPrice: cciptypes.GasPrice{
			Value:             in.Price.Value.Int(),
			DestChainSelector: in.Price.DestChainSelector,
		},
		TimestampUnixSec: in.UnixTimestamp.Int(),
	}
}

func tokenPriceUpdateWithTxMetaSlice(in []*ccippb.TokenPriceUpdateWithTxMeta) []cciptypes.TokenPriceUpdateWithTxMeta {
	out := make([]cciptypes.TokenPriceUpdateWithTxMeta, len(in))
	for i, v := range in {
		out[i] = tokenPriceUpdateWithTxMeta(v)
	}
	return out
}

func tokenPriceUpdateWithTxMeta(in *ccippb.TokenPriceUpdateWithTxMeta) cciptypes.TokenPriceUpdateWithTxMeta {
	return cciptypes.TokenPriceUpdateWithTxMeta{
		TxMeta:           txMeta(in.TxMeta),
		TokenPriceUpdate: tokenPriceUpdate(in.TokenPriceUpdate),
	}
}

func tokenPriceUpdate(in *ccippb.TokenPriceUpdate) cciptypes.TokenPriceUpdate {
	return cciptypes.TokenPriceUpdate{
		TokenPrice:       tokenPrice(in.Price),
		TimestampUnixSec: in.UnixTimestamp.Int(),
	}
}

func tokenPrice(in *ccippb.TokenPrice) cciptypes.TokenPrice {
	return cciptypes.TokenPrice{
		Value: in.Value.Int(),
		Token: cciptypes.Address(in.Token),
	}
}

func tokenPriceUpdateSlice(in []*ccippb.TokenPriceUpdate) []cciptypes.TokenPriceUpdate {
	out := make([]cciptypes.TokenPriceUpdate, len(in))
	for i, v := range in {
		out[i] = tokenPriceUpdate(v)
	}
	return out
}

func decimals(in []uint32) ([]uint8, error) {
	out := make([]uint8, len(in))
	for i, v := range in {
		if v > 255 {
			return nil, fmt.Errorf("decimal value %d at %d is too large", v, i)
		}
		out[i] = uint8(v)
	}
	return out, nil
}

func decimalsPB(in []uint8) *ccippb.GetTokensDecimalsResponse {
	out := make([]uint32, len(in))
	for i, v := range in {
		out[i] = uint32(v)
	}
	return &ccippb.GetTokensDecimalsResponse{Decimals: out}
}

func gasPriceUpdateWithTxMetaSlicePB(in []cciptypes.GasPriceUpdateWithTxMeta) []*ccippb.GasPriceUpdateWithTxMeta {
	out := make([]*ccippb.GasPriceUpdateWithTxMeta, len(in))
	for i, v := range in {
		out[i] = gasPriceUpdateWithTxMetaPB(v)
	}
	return out
}

func gasPriceUpdateWithTxMetaPB(in cciptypes.GasPriceUpdateWithTxMeta) *ccippb.GasPriceUpdateWithTxMeta {
	return &ccippb.GasPriceUpdateWithTxMeta{
		TxMeta:         txMetaPB(in.TxMeta),
		GasPriceUpdate: gasPriceUpdatePB(in.GasPriceUpdate),
	}
}

func gasPriceUpdatePB(in cciptypes.GasPriceUpdate) *ccippb.GasPriceUpdate {
	return &ccippb.GasPriceUpdate{
		Price:         gasPricePB(in.GasPrice),
		UnixTimestamp: pb.NewBigIntFromInt(in.TimestampUnixSec),
	}
}

func gasPricePB(in cciptypes.GasPrice) *ccippb.GasPrice {
	return &ccippb.GasPrice{
		Value:             pb.NewBigIntFromInt(in.Value),
		DestChainSelector: in.DestChainSelector,
	}
}

func tokenPriceUpdateWithTxMetaSlicePB(in []cciptypes.TokenPriceUpdateWithTxMeta) []*ccippb.TokenPriceUpdateWithTxMeta {
	out := make([]*ccippb.TokenPriceUpdateWithTxMeta, len(in))
	for i, v := range in {
		out[i] = tokenPriceUpdateWithTxMetaPB(v)
	}
	return out
}

func tokenPriceUpdateWithTxMetaPB(in cciptypes.TokenPriceUpdateWithTxMeta) *ccippb.TokenPriceUpdateWithTxMeta {
	return &ccippb.TokenPriceUpdateWithTxMeta{
		TxMeta:           txMetaPB(in.TxMeta),
		TokenPriceUpdate: tokenPriceUpdatePB(in.TokenPriceUpdate),
	}
}

func tokenPriceUpdatePB(in cciptypes.TokenPriceUpdate) *ccippb.TokenPriceUpdate {
	return &ccippb.TokenPriceUpdate{
		Price:         tokenPricePB(in.TokenPrice),
		UnixTimestamp: pb.NewBigIntFromInt(in.TimestampUnixSec),
	}
}

func tokenPricePB(in cciptypes.TokenPrice) *ccippb.TokenPrice {
	return &ccippb.TokenPrice{
		Value: pb.NewBigIntFromInt(in.Value),
		Token: string(in.Token),
	}
}

func tokenPriceUpdateSlicePB(in []cciptypes.TokenPriceUpdate) []*ccippb.TokenPriceUpdate {
	out := make([]*ccippb.TokenPriceUpdate, len(in))
	for i, v := range in {
		out[i] = tokenPriceUpdatePB(v)
	}
	return out
}
