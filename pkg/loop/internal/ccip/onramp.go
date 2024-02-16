package ccip

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var _ cciptypes.OnRampReader = (*OnRampReaderClient)(nil)

type OnRampReaderClient struct {
	grpc ccippb.OnRampReaderClient
}

func NewOnRampReaderClient(grpc ccippb.OnRampReaderClient) *OnRampReaderClient {
	return &OnRampReaderClient{grpc: grpc}
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
	panic("unimplemented")
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

func txMeta(meta *ccippb.TxMeta) cciptypes.TxMeta {
	return cciptypes.TxMeta{
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
