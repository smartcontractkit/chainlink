package ccipsolana

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
}

func NewExecutePluginCodecV1() *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	if len(report.ChainReports) == 0 || len(report.ChainReports) > 1 {
		return nil, fmt.Errorf("unexpected chain report length: %d", len(report.ChainReports))
	}

	var message ccip_router.Any2SolanaRampMessage
	chainReport := report.ChainReports[0]
	if len(chainReport.Messages) > 0 {
		// currently only allow committing one message at a time
		msg := chainReport.Messages[0]
		tokenAmounts := make([]ccip_router.Any2SolanaTokenTransfer, 0, len(msg.TokenAmounts))
		for _, tokenAmount := range msg.TokenAmounts {
			if tokenAmount.Amount.IsEmpty() {
				return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
			}

			DestTokenAddress, err := solana.PublicKeyFromBase58(string(tokenAmount.DestTokenAddress))
			if err != nil {
				return nil, fmt.Errorf("invalid destTokenAddress address: %w", err)
			}

			tokenAmounts = append(tokenAmounts, ccip_router.Any2SolanaTokenTransfer{
				SourcePoolAddress: tokenAmount.SourcePoolAddress,
				DestTokenAddress:  DestTokenAddress,
				ExtraData:         tokenAmount.ExtraData,
				Amount:            bigIntToBytes32(tokenAmount.Amount),
				DestGasAmount:     bytesToUint32(tokenAmount.DestExecData),
			})
		}

		var extraArgs ccip_router.SolanaExtraArgs
		decoder := agbinary.NewBorshDecoder(msg.ExtraArgs)
		err := extraArgs.UnmarshalWithDecoder(decoder)
		if err != nil {
			return nil, fmt.Errorf("invalid extra arguments: %w", err)
		}

		receiver, err := solana.PublicKeyFromBase58(string(msg.Receiver))
		if err != nil {
			return nil, fmt.Errorf("invalid receiver address: %s, %w", string(msg.Receiver), err)
		}

		message = ccip_router.Any2SolanaRampMessage{
			Header: ccip_router.RampMessageHeader{
				MessageId:           msg.Header.MessageID,
				SourceChainSelector: uint64(msg.Header.SourceChainSelector),
				DestChainSelector:   uint64(msg.Header.DestChainSelector),
				SequenceNumber:      uint64(msg.Header.SequenceNumber),
				Nonce:               msg.Header.Nonce,
			},
			Sender:       msg.Sender,
			Data:         msg.Data,
			Receiver:     receiver,
			TokenAmounts: tokenAmounts,
			ExtraArgs:    extraArgs,
		}
	}

	var offchainTokenData [][]byte
	if len(chainReport.OffchainTokenData) > 0 {
		offchainTokenData = chainReport.OffchainTokenData[0]
	}

	solanaProofs := make([][32]byte, 0, len(chainReport.Proofs))
	for _, proof := range chainReport.Proofs {
		solanaProofs = append(solanaProofs, proof)
	}

	solanaReport := ccip_router.ExecutionReportSingleChain{
		SourceChainSelector: uint64(chainReport.SourceChainSelector),
		Message:             message,
		OffchainTokenData:   offchainTokenData,
		Proofs:              solanaProofs,
		// TODO: add TokenIndexes in the cciptypes.ExecutePluginReport
	}

	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	err := solanaReport.MarshalWithEncoder(encoder)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	decoder := agbinary.NewBorshDecoder(encodedReport)
	executeReport := ccip_router.ExecutionReportSingleChain{}
	err := executeReport.UnmarshalWithDecoder(decoder)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack encoded report: %w", err)
	}

	tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(executeReport.Message.TokenAmounts))
	for _, tokenAmount := range executeReport.Message.TokenAmounts {
		destData := make([]byte, 4)
		binary.BigEndian.PutUint32(destData, tokenAmount.DestGasAmount)

		tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
			SourcePoolAddress: tokenAmount.SourcePoolAddress,
			DestTokenAddress:  cciptypes.UnknownAddress(tokenAmount.DestTokenAddress.String()),
			ExtraData:         tokenAmount.ExtraData,
			Amount:            priceHelper(tokenAmount.Amount[:]),
			DestExecData:      destData,
		})
	}

	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	err = executeReport.Message.ExtraArgs.MarshalWithEncoder(encoder)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack encoded report: %w", err)
	}

	messages := []cciptypes.Message{
		{
			Header: cciptypes.RampMessageHeader{
				MessageID:           executeReport.Message.Header.MessageId,
				SourceChainSelector: cciptypes.ChainSelector(executeReport.Message.Header.SourceChainSelector),
				DestChainSelector:   cciptypes.ChainSelector(executeReport.Message.Header.DestChainSelector),
				SequenceNumber:      cciptypes.SeqNum(executeReport.Message.Header.SequenceNumber),
				Nonce:               executeReport.Message.Header.Nonce,
				MsgHash:             cciptypes.Bytes32{},        // todo: info not available, but not required atm
				OnRamp:              cciptypes.UnknownAddress{}, // todo: info not available, but not required atm
			},
			Sender:         executeReport.Message.Sender,
			Data:           executeReport.Message.Data,
			Receiver:       cciptypes.UnknownAddress(executeReport.Message.Receiver.String()),
			ExtraArgs:      buf.Bytes(),
			FeeToken:       cciptypes.UnknownAddress{}, // <-- todo: info not available, but not required atm
			FeeTokenAmount: cciptypes.BigInt{},         // <-- todo: info not available, but not required atm
			TokenAmounts:   tokenAmounts,
		},
	}

	offchainTokenData := make([][][]byte, 0, 1)
	if executeReport.OffchainTokenData != nil {
		offchainTokenData = append(offchainTokenData, executeReport.OffchainTokenData)
	}

	proofs := make([]cciptypes.Bytes32, 0, len(executeReport.Proofs))
	for _, proof := range executeReport.Proofs {
		proofs = append(proofs, proof)
	}

	chainReport := cciptypes.ExecutePluginReportSingleChain{
		SourceChainSelector: cciptypes.ChainSelector(executeReport.SourceChainSelector),
		Messages:            messages,
		OffchainTokenData:   offchainTokenData,
		Proofs:              proofs,
	}

	report := cciptypes.ExecutePluginReport{
		ChainReports: []cciptypes.ExecutePluginReportSingleChain{chainReport},
	}

	return report, nil
}

func bytesToUint32(b []byte) uint32 {
	if len(b) < 4 {
		var padded [4]byte
		copy(padded[4-len(b):], b) // Pad from the right for big-endian
		return binary.BigEndian.Uint32(padded[:])
	}

	return binary.BigEndian.Uint32(b)
}

func bigIntToBytes32(n cciptypes.BigInt) [32]uint8 {
	var b [32]uint8
	raw := n.Bytes()
	copy(b[32-len(raw):], raw) // Right-align and zero-pad
	return b
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
