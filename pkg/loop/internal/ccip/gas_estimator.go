package ccip

import (
	"context"
	"fmt"
	"math/big"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// CommitGasEstimatorGRPCClient implements [cciptypes.GasEstimatorCommit] by wrapping a
// [ccippb.CommitGasEstimatorReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call the CommitGasEstimatorReader service, which
// is hosted by the relayer
type CommitGasEstimatorGRPCClient struct {
	client ccippb.GasPriceEstimatorCommitClient
}

func NewCommitGasEstimatorGRPCClient(cc grpc.ClientConnInterface) *CommitGasEstimatorGRPCClient {
	return &CommitGasEstimatorGRPCClient{client: ccippb.NewGasPriceEstimatorCommitClient(cc)}
}

// CommitGasEstimatorGRPCServer implements [ccippb.CommitGasEstimatorReaderServer] by wrapping a
// [cciptypes.GasEstimatorCommit] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [CommitGasEstimatorGRPCClient]
type CommitGasEstimatorGRPCServer struct {
	ccippb.UnimplementedGasPriceEstimatorCommitServer

	impl cciptypes.GasPriceEstimatorCommit
}

func NewCommitGasEstimatorGRPCServer(impl cciptypes.GasPriceEstimatorCommit) *CommitGasEstimatorGRPCServer {
	return &CommitGasEstimatorGRPCServer{impl: impl}
}

var _ cciptypes.GasPriceEstimatorCommit = (*CommitGasEstimatorGRPCClient)(nil)
var _ ccippb.GasPriceEstimatorCommitServer = (*CommitGasEstimatorGRPCServer)(nil)

// DenoteInUSD implements ccip.GasPriceEstimatorCommit.
func (c *CommitGasEstimatorGRPCClient) DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	resp, err := c.client.DenoteInUSD(context.Background(), &ccippb.DenoteInUSDRequest{
		P:                  pb.NewBigIntFromInt(p),
		WrappedNativePrice: pb.NewBigIntFromInt(wrappedNativePrice),
	})
	if err != nil {
		return nil, err
	}
	return resp.UsdPrice.Int(), nil
}

// Deviates implements ccip.GasPriceEstimatorCommit.
func (c *CommitGasEstimatorGRPCClient) Deviates(p1 *big.Int, p2 *big.Int) (bool, error) {
	resp, err := c.client.Deviates(context.Background(), &ccippb.DeviatesRequest{
		P1: pb.NewBigIntFromInt(p1),
		P2: pb.NewBigIntFromInt(p2),
	})
	if err != nil {
		return false, err
	}
	return resp.Deviates, nil
}

// GetGasPrice implements ccip.GasPriceEstimatorCommit.
func (c *CommitGasEstimatorGRPCClient) GetGasPrice(ctx context.Context) (*big.Int, error) {
	resp, err := c.client.GetGasPrice(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.GasPrice.Int(), nil
}

// Median implements ccip.GasPriceEstimatorCommit.
func (c *CommitGasEstimatorGRPCClient) Median(gasPrices []*big.Int) (*big.Int, error) {
	resp, err := c.client.Median(context.Background(), &ccippb.MedianRequest{
		GasPrices: bigIntSlicePB(gasPrices),
	})
	if err != nil {
		return nil, err
	}
	return resp.GasPrice.Int(), nil
}

// Server implementation

// DenoteInUSD implements ccippb.GasPriceEstimatorCommitServer.
func (c *CommitGasEstimatorGRPCServer) DenoteInUSD(ctx context.Context, req *ccippb.DenoteInUSDRequest) (*ccippb.DenoteInUSDResponse, error) {
	usd, err := c.impl.DenoteInUSD(req.P.Int(), req.WrappedNativePrice.Int())
	if err != nil {
		return nil, err
	}
	return &ccippb.DenoteInUSDResponse{UsdPrice: pb.NewBigIntFromInt(usd)}, nil
}

// Deviates implements ccippb.GasPriceEstimatorCommitServer.
func (c *CommitGasEstimatorGRPCServer) Deviates(ctx context.Context, req *ccippb.DeviatesRequest) (*ccippb.DeviatesResponse, error) {
	deviates, err := c.impl.Deviates(req.P1.Int(), req.P2.Int())
	if err != nil {
		return nil, err
	}
	return &ccippb.DeviatesResponse{Deviates: deviates}, nil
}

// GetGasPrice implements ccippb.GasPriceEstimatorCommitServer.
func (c *CommitGasEstimatorGRPCServer) GetGasPrice(ctx context.Context, req *emptypb.Empty) (*ccippb.GetGasPriceResponse, error) {
	gasPrice, err := c.impl.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetGasPriceResponse{GasPrice: pb.NewBigIntFromInt(gasPrice)}, nil
}

// Median implements ccippb.GasPriceEstimatorCommitServer.
func (c *CommitGasEstimatorGRPCServer) Median(ctx context.Context, req *ccippb.MedianRequest) (*ccippb.MedianResponse, error) {
	gasPrice, err := c.impl.Median(bigIntSlice(req.GasPrices))
	if err != nil {
		return nil, err
	}
	return &ccippb.MedianResponse{GasPrice: pb.NewBigIntFromInt(gasPrice)}, nil
}

// ExecGasEstimatorGRPCClient implements [cciptypes.GasEstimatorExec] by wrapping a
// [ccippb.ExecGasEstimatorReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call the ExecGasEstimatorReader service, which
// is hosted by the relayer
type ExecGasEstimatorGRPCClient struct {
	client ccippb.GasPriceEstimatorExecClient
}

func NewExecGasEstimatorGRPCClient(cc grpc.ClientConnInterface) *ExecGasEstimatorGRPCClient {
	return &ExecGasEstimatorGRPCClient{client: ccippb.NewGasPriceEstimatorExecClient(cc)}
}

// ExecGasEstimatorGRPCServer implements [ccippb.ExecGasEstimatorReaderServer] by wrapping a
// [cciptypes.GasEstimatorExec] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [ExecGasEstimatorGRPCClient]
type ExecGasEstimatorGRPCServer struct {
	ccippb.UnimplementedGasPriceEstimatorExecServer

	impl cciptypes.GasPriceEstimatorExec
}

func NewExecGasEstimatorGRPCServer(impl cciptypes.GasPriceEstimatorExec) *ExecGasEstimatorGRPCServer {
	return &ExecGasEstimatorGRPCServer{impl: impl}
}

// ensure interfaces are implemented
var _ cciptypes.GasPriceEstimatorExec = (*ExecGasEstimatorGRPCClient)(nil)
var _ ccippb.GasPriceEstimatorExecServer = (*ExecGasEstimatorGRPCServer)(nil)

// DenoteInUSD implements ccip.GasPriceEstimatorExec.
func (e *ExecGasEstimatorGRPCClient) DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	resp, err := e.client.DenoteInUSD(context.Background(), &ccippb.DenoteInUSDRequest{
		P:                  pb.NewBigIntFromInt(p),
		WrappedNativePrice: pb.NewBigIntFromInt(wrappedNativePrice),
	})
	if err != nil {
		return nil, err
	}
	return resp.UsdPrice.Int(), nil
}

// EstimateMsgCostUSD implements ccip.GasPriceEstimatorExec.
func (e *ExecGasEstimatorGRPCClient) EstimateMsgCostUSD(p *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	msgPB := evm2EVMOnRampCCIPSendRequestedWithMeta(&msg)
	resp, err := e.client.EstimateMsgCostUSD(context.Background(), &ccippb.EstimateMsgCostUSDRequest{
		P:                  pb.NewBigIntFromInt(p),
		WrappedNativePrice: pb.NewBigIntFromInt(wrappedNativePrice),
		Msg:                msgPB,
	})
	if err != nil {
		return nil, err
	}
	return resp.UsdCost.Int(), nil
}

// GetGasPrice implements ccip.GasPriceEstimatorExec.
func (e *ExecGasEstimatorGRPCClient) GetGasPrice(ctx context.Context) (*big.Int, error) {
	resp, err := e.client.GetGasPrice(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.GasPrice.Int(), nil
}

// Median implements ccip.GasPriceEstimatorExec.
func (e *ExecGasEstimatorGRPCClient) Median(gasPrices []*big.Int) (*big.Int, error) {
	resp, err := e.client.Median(context.Background(), &ccippb.MedianRequest{
		GasPrices: bigIntSlicePB(gasPrices),
	})
	if err != nil {
		return nil, err
	}
	return resp.GasPrice.Int(), nil
}

// DenoteInUSD implements ccippb.GasPriceEstimatorExecServer.
func (e *ExecGasEstimatorGRPCServer) DenoteInUSD(ctx context.Context, req *ccippb.DenoteInUSDRequest) (*ccippb.DenoteInUSDResponse, error) {
	usd, err := e.impl.DenoteInUSD(req.P.Int(), req.WrappedNativePrice.Int())
	if err != nil {
		return nil, err
	}
	return &ccippb.DenoteInUSDResponse{UsdPrice: pb.NewBigIntFromInt(usd)}, nil
}

// EstimateMsgCostUSD implements ccippb.GasPriceEstimatorExecServer.
func (e *ExecGasEstimatorGRPCServer) EstimateMsgCostUSD(ctx context.Context, req *ccippb.EstimateMsgCostUSDRequest) (*ccippb.EstimateMsgCostUSDResponse, error) {
	msg, err := evm2EVMOnRampCCIPSendRequestedWithMetaPB(req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to convert evm2evm msg: %w", err)
	}
	cost, err := e.impl.EstimateMsgCostUSD(req.P.Int(), req.WrappedNativePrice.Int(), *msg)
	if err != nil {
		return nil, err
	}
	return &ccippb.EstimateMsgCostUSDResponse{UsdCost: pb.NewBigIntFromInt(cost)}, nil
}

// GetGasPrice implements ccippb.GasPriceEstimatorExecServer.
func (e *ExecGasEstimatorGRPCServer) GetGasPrice(ctx context.Context, req *emptypb.Empty) (*ccippb.GetGasPriceResponse, error) {
	price, err := e.impl.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetGasPriceResponse{GasPrice: pb.NewBigIntFromInt(price)}, nil
}

// Median implements ccippb.GasPriceEstimatorExecServer.
func (e *ExecGasEstimatorGRPCServer) Median(ctx context.Context, req *ccippb.MedianRequest) (*ccippb.MedianResponse, error) {
	median, err := e.impl.Median(bigIntSlice(req.GasPrices))
	if err != nil {
		return nil, err
	}
	return &ccippb.MedianResponse{GasPrice: pb.NewBigIntFromInt(median)}, nil
}

func bigIntSlicePB(in []*big.Int) []*pb.BigInt {
	out := make([]*pb.BigInt, len(in))
	for i, v := range in {
		out[i] = pb.NewBigIntFromInt(v)
	}
	return out
}

func bigIntSlice(in []*pb.BigInt) []*big.Int {
	out := make([]*big.Int, len(in))
	for i, v := range in {
		out[i] = v.Int()
	}
	return out
}

func evm2EVMOnRampCCIPSendRequestedWithMeta(in *cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) *ccippb.EVM2EVMOnRampCCIPSendRequestedWithMeta {
	return &ccippb.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EvmToEvmMsg:    evm2EVMMessagePB(in.EVM2EVMMessage),
		BlockTimestamp: timestamppb.New(in.BlockTimestamp),
		Executed:       in.Executed,
		Finalized:      in.Finalized,
		LogIndex:       uint64(in.LogIndex),
		TxHash:         in.TxHash,
	}
}

func evm2EVMOnRampCCIPSendRequestedWithMetaPB(pb *ccippb.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, error) {
	msg, err := evm2EVMMessage(pb.EvmToEvmMsg)
	if err != nil {
		return nil, err
	}
	return &cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: msg,
		BlockTimestamp: pb.BlockTimestamp.AsTime(),
		Executed:       pb.Executed,
		Finalized:      pb.Finalized,
		LogIndex:       uint(pb.LogIndex),
		TxHash:         pb.TxHash,
	}, nil
}
