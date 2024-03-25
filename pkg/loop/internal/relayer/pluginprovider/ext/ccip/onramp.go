package ccip

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var _ cciptypes.OnRampReader = (*OnRampReaderGRPCClient)(nil)

type OnRampReaderGRPCClient struct {
	grpc ccippb.OnRampReaderClient
}

func NewOnRampReaderGRPCClient(cc grpc.ClientConnInterface) *OnRampReaderGRPCClient {
	return &OnRampReaderGRPCClient{grpc: ccippb.NewOnRampReaderClient(cc)}
}

// Address implements ccip.OnRampReader.
func (o *OnRampReaderGRPCClient) Address(ctx context.Context) (cciptypes.Address, error) {
	resp, err := o.grpc.Address(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.Address), nil
}

func (o *OnRampReaderGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), o.grpc)
}

// GetDynamicConfig implements ccip.OnRampReader.
func (o *OnRampReaderGRPCClient) GetDynamicConfig(ctx context.Context) (cciptypes.OnRampDynamicConfig, error) {
	resp, err := o.grpc.GetDynamicConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.OnRampDynamicConfig{}, err
	}
	return onRampDynamicConfig(resp.DynamicConfig), nil
}

// GetSendRequestsBetweenSeqNums implements ccip.OnRampReader.
func (o *OnRampReaderGRPCClient) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, finalized bool) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
	resp, err := o.grpc.GetSendRequestsBetweenSeqNums(ctx, &ccippb.GetSendRequestsBetweenSeqNumsRequest{
		SeqNumMin: seqNumMin,
		SeqNumMax: seqNumMax,
		Finalized: finalized,
	})
	if err != nil {
		return nil, err
	}
	return evm2EVMMessageWithTxMetaSlice(resp.SendRequests)
}

// IsSourceChainHealthy returns true if the source chain is healthy.
func (o *OnRampReaderGRPCClient) IsSourceChainHealthy(ctx context.Context) (bool, error) {
	resp, err := o.grpc.IsSourceChainHealthy(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.IsHealthy, nil
}

// IsSourceCursed returns true if the source chain is cursed. OnRamp communicates with the underlying RMN
// to verify if source chain was cursed or not.
func (o *OnRampReaderGRPCClient) IsSourceCursed(ctx context.Context) (bool, error) {
	resp, err := o.grpc.IsSourceCursed(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.IsCursed, nil
}

// RouterAddress implements ccip.OnRampReader.
func (o *OnRampReaderGRPCClient) RouterAddress(ctx context.Context) (cciptypes.Address, error) {
	resp, err := o.grpc.RouterAddress(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.RouterAddress), nil
}

// SourcePriceRegistryAddress returns the address of the current price registry configured on the onRamp.
func (o *OnRampReaderGRPCClient) SourcePriceRegistryAddress(ctx context.Context) (cciptypes.Address, error) {
	resp, err := o.grpc.SourcePriceRegistryAddress(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.PriceRegistryAddress), nil
}

// Server

type OnRampReaderGRPCServer struct {
	ccippb.UnimplementedOnRampReaderServer

	impl cciptypes.OnRampReader
	deps []io.Closer
}

var _ ccippb.OnRampReaderServer = (*OnRampReaderGRPCServer)(nil)

func NewOnRampReaderGRPCServer(impl cciptypes.OnRampReader) *OnRampReaderGRPCServer {
	return &OnRampReaderGRPCServer{impl: impl}
}

// AddDep adds a dependency to the server that will be closed when the server is closed.
func (o *OnRampReaderGRPCServer) AddDep(dep io.Closer) *OnRampReaderGRPCServer {
	o.deps = append(o.deps, dep)
	return o
}

// Address implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) Address(ctx context.Context, _ *emptypb.Empty) (*ccippb.OnrampAddressResponse, error) {
	addr, err := o.impl.Address(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.OnrampAddressResponse{Address: string(addr)}, nil
}

func (o *OnRampReaderGRPCServer) Close(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(o.deps).Close()
}

// GetDynamicConfig implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) GetDynamicConfig(ctx context.Context, _ *emptypb.Empty) (*ccippb.GetDynamicConfigResponse, error) {
	c, err := o.impl.GetDynamicConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetDynamicConfigResponse{DynamicConfig: onRampDynamicConfigPB(&c)}, nil
}

// GetSendRequestsBetweenSeqNums implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) GetSendRequestsBetweenSeqNums(ctx context.Context, req *ccippb.GetSendRequestsBetweenSeqNumsRequest) (*ccippb.GetSendRequestsBetweenSeqNumsResponse, error) {
	sendRequests, err := o.impl.GetSendRequestsBetweenSeqNums(ctx, req.SeqNumMin, req.SeqNumMax, req.Finalized)
	if err != nil {
		return nil, err
	}
	sendRequestsPB, err := evm2EVMMessageWithTxMetaSlicePB(sendRequests)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetSendRequestsBetweenSeqNumsResponse{SendRequests: sendRequestsPB}, nil
}

// IsSourceChainHealthy implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) IsSourceChainHealthy(ctx context.Context, _ *emptypb.Empty) (*ccippb.IsSourceChainHealthyResponse, error) {
	isHealthy, err := o.impl.IsSourceChainHealthy(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.IsSourceChainHealthyResponse{IsHealthy: isHealthy}, nil
}

// IsSourceCursed implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) IsSourceCursed(ctx context.Context, _ *emptypb.Empty) (*ccippb.IsSourceCursedResponse, error) {
	isCursed, err := o.impl.IsSourceCursed(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.IsSourceCursedResponse{IsCursed: isCursed}, nil
}

// RouterAddress implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) RouterAddress(ctx context.Context, _ *emptypb.Empty) (*ccippb.RouterAddressResponse, error) {
	a, err := o.impl.RouterAddress(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.RouterAddressResponse{RouterAddress: string(a)}, nil
}

// SourcePriceRegistryAddress implements ccippb.OnRampReaderServer.
func (o *OnRampReaderGRPCServer) SourcePriceRegistryAddress(ctx context.Context, _ *emptypb.Empty) (*ccippb.SourcePriceRegistryAddressResponse, error) {
	a, err := o.impl.SourcePriceRegistryAddress(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.SourcePriceRegistryAddressResponse{PriceRegistryAddress: string(a)}, nil
}

func onRampDynamicConfig(config *ccippb.OnRampDynamicConfig) cciptypes.OnRampDynamicConfig {
	return cciptypes.OnRampDynamicConfig{
		Router:                            cciptypes.Address(config.Router),
		MaxNumberOfTokensPerMsg:           uint16(config.MaxNumberOfTokensPerMsg),
		DestGasOverhead:                   config.DestGasOverhead,
		DestGasPerPayloadByte:             uint16(config.DestGasPerByte),
		DestDataAvailabilityOverheadGas:   config.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    uint16(config.DestGasPerDataAvailabilityByte),
		DestDataAvailabilityMultiplierBps: uint16(config.DestDataAvailabilityMultiplierBps),
		PriceRegistry:                     cciptypes.Address(config.PriceRegistry),
		MaxDataBytes:                      config.MaxDataBytes,
		MaxPerMsgGasLimit:                 config.MaxPerMsgGasLimit,
	}
}

func onRampDynamicConfigPB(config *cciptypes.OnRampDynamicConfig) *ccippb.OnRampDynamicConfig {
	return &ccippb.OnRampDynamicConfig{
		Router:                            string(config.Router),
		MaxNumberOfTokensPerMsg:           uint32(config.MaxNumberOfTokensPerMsg),
		DestGasOverhead:                   config.DestGasOverhead,
		DestGasPerByte:                    uint32(config.DestGasPerPayloadByte),
		DestDataAvailabilityOverheadGas:   config.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    uint32(config.DestGasPerDataAvailabilityByte),
		DestDataAvailabilityMultiplierBps: uint32(config.DestDataAvailabilityMultiplierBps),
		PriceRegistry:                     string(config.PriceRegistry),
		MaxDataBytes:                      config.MaxDataBytes,
		MaxPerMsgGasLimit:                 config.MaxPerMsgGasLimit,
	}
}

func evm2EVMMessageWithTxMetaSlice(messages []*ccippb.EVM2EVMMessageWithTxMeta) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
	res := make([]cciptypes.EVM2EVMMessageWithTxMeta, len(messages))
	for i, m := range messages {
		decodedMsg, err := evm2EVMMessage(m.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to convert grpc message (%v) to evm2evm message: %w", m.Message, err)
		}
		res[i] = cciptypes.EVM2EVMMessageWithTxMeta{
			TxMeta:         txMeta(m.TxMeta),
			EVM2EVMMessage: decodedMsg,
		}
	}
	return res, nil
}

func evm2EVMMessageWithTxMetaSlicePB(messages []cciptypes.EVM2EVMMessageWithTxMeta) ([]*ccippb.EVM2EVMMessageWithTxMeta, error) {
	res := make([]*ccippb.EVM2EVMMessageWithTxMeta, len(messages))
	for i, m := range messages {
		decodedMsg := evm2EVMMessagePB(m.EVM2EVMMessage)
		res[i] = &ccippb.EVM2EVMMessageWithTxMeta{
			TxMeta:  txMetaPB(m.TxMeta),
			Message: decodedMsg,
		}
	}
	return res, nil
}

func evm2EVMMessage(message *ccippb.EVM2EVMMessage) (cciptypes.EVM2EVMMessage, error) {
	msgID, err := hash(message.MessageId)
	if err != nil {
		return cciptypes.EVM2EVMMessage{}, fmt.Errorf("failed to convert message id (%v): %w", message.MessageId, err)
	}

	return cciptypes.EVM2EVMMessage{
		SequenceNumber:      message.SequenceNumber,
		GasLimit:            message.GasLimit.Int(),
		Nonce:               message.Nonce,
		MessageID:           msgID,
		SourceChainSelector: message.SourceChainSelector,
		Sender:              cciptypes.Address(message.Sender),
		Receiver:            cciptypes.Address(message.Receiver),
		Strict:              message.Strict,
		FeeToken:            cciptypes.Address(message.FeeToken),
		FeeTokenAmount:      message.FeeTokenAmount.Int(),
		Data:                message.Data,
		TokenAmounts:        tokenAmountSlice(message.TokenAmounts),
		SourceTokenData:     message.SourceTokenData,
	}, nil
}

func evm2EVMMessagePB(message cciptypes.EVM2EVMMessage) *ccippb.EVM2EVMMessage {
	return &ccippb.EVM2EVMMessage{
		SequenceNumber:      message.SequenceNumber,
		GasLimit:            pb.NewBigIntFromInt(message.GasLimit),
		Nonce:               message.Nonce,
		MessageId:           message.MessageID[:],
		SourceChainSelector: message.SourceChainSelector,
		Sender:              string(message.Sender),
		Receiver:            string(message.Receiver),
		Strict:              message.Strict,
		FeeToken:            string(message.FeeToken),
		FeeTokenAmount:      pb.NewBigIntFromInt(message.FeeTokenAmount),
		Data:                message.Data,
		TokenAmounts:        tokenAmountSlicePB(message.TokenAmounts),
		SourceTokenData:     message.SourceTokenData,
	}
}

func tokenAmountSlice(tokenAmounts []*ccippb.TokenAmount) []cciptypes.TokenAmount {
	res := make([]cciptypes.TokenAmount, len(tokenAmounts))
	for i, t := range tokenAmounts {
		res[i] = cciptypes.TokenAmount{
			Token:  cciptypes.Address(t.Token),
			Amount: t.Amount.Int(),
		}
	}
	return res
}

func tokenAmountSlicePB(tokenAmounts []cciptypes.TokenAmount) []*ccippb.TokenAmount {
	res := make([]*ccippb.TokenAmount, len(tokenAmounts))
	for i, t := range tokenAmounts {
		res[i] = &ccippb.TokenAmount{
			Token:  string(t.Token),
			Amount: pb.NewBigIntFromInt(t.Amount),
		}
	}
	return res
}

func txMeta(meta *ccippb.TxMeta) cciptypes.TxMeta {
	return cciptypes.TxMeta{
		BlockTimestampUnixMilli: meta.BlockTimestampUnixMilli,
		BlockNumber:             meta.BlockNumber,
		TxHash:                  meta.TxHash,
		LogIndex:                meta.LogIndex,
	}
}

func txMetaPB(meta cciptypes.TxMeta) *ccippb.TxMeta {
	return &ccippb.TxMeta{
		BlockTimestampUnixMilli: meta.BlockTimestampUnixMilli,
		BlockNumber:             meta.BlockNumber,
		TxHash:                  meta.TxHash,
		LogIndex:                meta.LogIndex,
	}
}

func hash(h []byte) (cciptypes.Hash, error) {
	var res cciptypes.Hash
	if len(h) != 32 {
		return res, fmt.Errorf("hash length is not 32 bytes")
	}
	copy(res[:], h)
	return res, nil
}
