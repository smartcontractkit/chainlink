package ccip

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var _ cciptypes.OnRampReader = (*OnRampReaderClient)(nil)

type OnRampReaderClient struct {
	grpc ccippb.OnRampReaderClient
}

func NewOnRampReaderClient(cc grpc.ClientConnInterface) *OnRampReaderClient {
	return &OnRampReaderClient{grpc: ccippb.NewOnRampReaderClient(cc)}
}

// Address implements ccip.OnRampReader.
func (o *OnRampReaderClient) Address() (cciptypes.Address, error) {
	resp, err := o.grpc.Address(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.Address), nil
}

// GetDynamicConfig implements ccip.OnRampReader.
func (o *OnRampReaderClient) GetDynamicConfig() (cciptypes.OnRampDynamicConfig, error) {
	resp, err := o.grpc.GetDynamicConfig(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return cciptypes.OnRampDynamicConfig{}, err
	}
	return onRampDynamicConfig(resp.DynamicConfig), nil
}

// GetSendRequestsBetweenSeqNums implements ccip.OnRampReader.
func (o *OnRampReaderClient) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, finalized bool) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
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

// RouterAddress implements ccip.OnRampReader.
func (o *OnRampReaderClient) RouterAddress() (cciptypes.Address, error) {
	resp, err := o.grpc.RouterAddress(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.RouterAddress), nil
}

// IsSourceChainHealthy returns true if the source chain is healthy.
func (o *OnRampReaderClient) IsSourceChainHealthy(ctx context.Context) (bool, error) {
	panic("unimplemented")
}

// IsSourceCursed returns true if the source chain is cursed. OnRamp communicates with the underlying RMN
// to verify if source chain was cursed or not.
func (o *OnRampReaderClient) IsSourceCursed(ctx context.Context) (bool, error) {
	panic("unimplemented")
}

// SourcePriceRegistryAddress returns the address of the current price registry configured on the onRamp.
func (o *OnRampReaderClient) SourcePriceRegistryAddress(ctx context.Context) (cciptypes.Address, error) {
	panic("unimplemented")
}

// Server

type OnRampReaderServer struct {
	ccippb.UnimplementedOnRampReaderServer

	impl cciptypes.OnRampReader
}

// mustEmbedUnimplementedOnRampReaderServer implements ccippb.OnRampReaderServer.

var _ ccippb.OnRampReaderServer = (*OnRampReaderServer)(nil)

func NewOnRampReaderServer(impl cciptypes.OnRampReader) *OnRampReaderServer {
	return &OnRampReaderServer{impl: impl}
}

// Address implements ccippb.OnRampReaderServer.
func (o *OnRampReaderServer) Address(context.Context, *emptypb.Empty) (*ccippb.OnrampAddressResponse, error) {
	addr, err := o.impl.Address()
	if err != nil {
		return nil, err
	}
	return &ccippb.OnrampAddressResponse{Address: string(addr)}, nil
}

// GetDynamicConfig implements ccippb.OnRampReaderServer.
func (o *OnRampReaderServer) GetDynamicConfig(context.Context, *emptypb.Empty) (*ccippb.GetDynamicConfigResponse, error) {
	c, err := o.impl.GetDynamicConfig()
	if err != nil {
		return nil, err
	}
	return &ccippb.GetDynamicConfigResponse{DynamicConfig: onRampDynamicConfigPB(&c)}, nil
}

// GetSendRequestsBetweenSeqNums implements ccippb.OnRampReaderServer.
func (o *OnRampReaderServer) GetSendRequestsBetweenSeqNums(ctx context.Context, req *ccippb.GetSendRequestsBetweenSeqNumsRequest) (*ccippb.GetSendRequestsBetweenSeqNumsResponse, error) {
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

// RouterAddress implements ccippb.OnRampReaderServer.
func (o *OnRampReaderServer) RouterAddress(context.Context, *emptypb.Empty) (*ccippb.RouterAddressResponse, error) {
	a, err := o.impl.RouterAddress()
	if err != nil {
		return nil, err
	}
	return &ccippb.RouterAddressResponse{RouterAddress: string(a)}, nil
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
