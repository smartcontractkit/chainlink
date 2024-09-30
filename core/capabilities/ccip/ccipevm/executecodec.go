package ccipevm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
	executeReportMethodInputs abi.Arguments
}

func NewExecutePluginCodecV1() *ExecutePluginCodecV1 {
	abiParsed, err := abi.JSON(strings.NewReader(offramp.OffRampABI))
	if err != nil {
		panic(fmt.Errorf("parse multi offramp abi: %s", err))
	}
	methodInputs := abihelpers.MustGetMethodInputs("manuallyExecute", abiParsed)
	if len(methodInputs) == 0 {
		panic("no inputs found for method: manuallyExecute")
	}

	return &ExecutePluginCodecV1{
		executeReportMethodInputs: methodInputs[:1],
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	evmReport := make([]offramp.InternalExecutionReportSingleChain, 0, len(report.ChainReports))

	for _, chainReport := range report.ChainReports {
		if chainReport.ProofFlagBits.IsEmpty() {
			return nil, fmt.Errorf("proof flag bits are empty")
		}

		evmProofs := make([][32]byte, 0, len(chainReport.Proofs))
		for _, proof := range chainReport.Proofs {
			evmProofs = append(evmProofs, proof)
		}

		evmMessages := make([]offramp.InternalAny2EVMRampMessage, 0, len(chainReport.Messages))
		for _, message := range chainReport.Messages {
			receiver := common.BytesToAddress(message.Receiver)

			tokenAmounts := make([]offramp.InternalRampTokenAmount, 0, len(message.TokenAmounts))
			for _, tokenAmount := range message.TokenAmounts {
				if tokenAmount.Amount.IsEmpty() {
					return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
				}

				tokenAmounts = append(tokenAmounts, offramp.InternalRampTokenAmount{
					SourcePoolAddress: tokenAmount.SourcePoolAddress,
					DestTokenAddress:  tokenAmount.DestTokenAddress,
					ExtraData:         tokenAmount.ExtraData,
					Amount:            tokenAmount.Amount.Int,
				})
			}

			gasLimit, err := decodeExtraArgsV1V2(message.ExtraArgs)
			if err != nil {
				return nil, fmt.Errorf("decode extra args to get gas limit: %w", err)
			}

			evmMessages = append(evmMessages, offramp.InternalAny2EVMRampMessage{
				Header: offramp.InternalRampMessageHeader{
					MessageId:           message.Header.MessageID,
					SourceChainSelector: uint64(message.Header.SourceChainSelector),
					DestChainSelector:   uint64(message.Header.DestChainSelector),
					SequenceNumber:      uint64(message.Header.SequenceNumber),
					Nonce:               message.Header.Nonce,
				},
				Sender:       message.Sender,
				Data:         message.Data,
				Receiver:     receiver,
				GasLimit:     gasLimit,
				TokenAmounts: tokenAmounts,
			})
		}

		evmChainReport := offramp.InternalExecutionReportSingleChain{
			SourceChainSelector: uint64(chainReport.SourceChainSelector),
			Messages:            evmMessages,
			OffchainTokenData:   chainReport.OffchainTokenData,
			Proofs:              evmProofs,
			ProofFlagBits:       chainReport.ProofFlagBits.Int,
		}
		evmReport = append(evmReport, evmChainReport)
	}

	return e.executeReportMethodInputs.PackValues([]interface{}{&evmReport})
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	unpacked, err := e.executeReportMethodInputs.Unpack(encodedReport)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack encoded report: %w", err)
	}
	if len(unpacked) != 1 {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpacked report is empty")
	}

	evmReportRaw := abi.ConvertType(unpacked[0], new([]offramp.InternalExecutionReportSingleChain))
	evmReportPtr, is := evmReportRaw.(*[]offramp.InternalExecutionReportSingleChain)
	if !is {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("got an unexpected report type %T", unpacked[0])
	}
	if evmReportPtr == nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("evm report is nil")
	}

	evmReport := *evmReportPtr
	executeReport := cciptypes.ExecutePluginReport{
		ChainReports: make([]cciptypes.ExecutePluginReportSingleChain, 0, len(evmReport)),
	}

	for _, evmChainReport := range evmReport {
		proofs := make([]cciptypes.Bytes32, 0, len(evmChainReport.Proofs))
		for _, proof := range evmChainReport.Proofs {
			proofs = append(proofs, proof)
		}

		messages := make([]cciptypes.Message, 0, len(evmChainReport.Messages))
		for _, evmMessage := range evmChainReport.Messages {
			tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(evmMessage.TokenAmounts))
			for _, tokenAmount := range evmMessage.TokenAmounts {
				tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
					SourcePoolAddress: tokenAmount.SourcePoolAddress,
					DestTokenAddress:  tokenAmount.DestTokenAddress,
					ExtraData:         tokenAmount.ExtraData,
					Amount:            cciptypes.NewBigInt(tokenAmount.Amount),
				})
			}

			message := cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           evmMessage.Header.MessageId,
					SourceChainSelector: cciptypes.ChainSelector(evmMessage.Header.SourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(evmMessage.Header.DestChainSelector),
					SequenceNumber:      cciptypes.SeqNum(evmMessage.Header.SequenceNumber),
					Nonce:               evmMessage.Header.Nonce,
					MsgHash:             cciptypes.Bytes32{}, // <-- todo: info not available, but not required atm
					OnRamp:              cciptypes.Bytes{},   // <-- todo: info not available, but not required atm
				},
				Sender:         evmMessage.Sender,
				Data:           evmMessage.Data,
				Receiver:       evmMessage.Receiver.Bytes(),
				ExtraArgs:      cciptypes.Bytes{},  // <-- todo: info not available, but not required atm
				FeeToken:       cciptypes.Bytes{},  // <-- todo: info not available, but not required atm
				FeeTokenAmount: cciptypes.BigInt{}, // <-- todo: info not available, but not required atm
				TokenAmounts:   tokenAmounts,
			}
			messages = append(messages, message)
		}

		chainReport := cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(evmChainReport.SourceChainSelector),
			Messages:            messages,
			OffchainTokenData:   evmChainReport.OffchainTokenData,
			Proofs:              proofs,
			ProofFlagBits:       cciptypes.NewBigInt(evmChainReport.ProofFlagBits),
		}

		executeReport.ChainReports = append(executeReport.ChainReports, chainReport)
	}

	return executeReport, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
