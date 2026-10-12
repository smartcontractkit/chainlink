package ccipevm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/utils/abihelpers"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0"
type ExecutePluginCodecV1 struct {
	executeReportMethodInputs abi.Arguments
	extraDataCodec            ccipocr3.ExtraDataCodecBundle
}

func NewExecutePluginCodecV1(extraDataCodec ccipocr3.ExtraDataCodecBundle) *ExecutePluginCodecV1 {
	abiParsed, err := abi.JSON(strings.NewReader(offramp.OffRampABI))
	if err != nil {
		panic(fmt.Errorf("parse multi offramp abi: %w", err))
	}
	methodInputs := abihelpers.MustGetMethodInputs("manuallyExecute", abiParsed)
	if len(methodInputs) == 0 {
		panic("no inputs found for method: manuallyExecute")
	}

	return &ExecutePluginCodecV1{
		executeReportMethodInputs: methodInputs[:1],
		extraDataCodec:            extraDataCodec,
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report ccipocr3.ExecutePluginReport) ([]byte, error) {
	evmReport := make([]offramp.InternalExecutionReport, 0, len(report.ChainReports))

	for _, chainReport := range report.ChainReports {
		sourceChainFamily, err := chainsel.GetSelectorFamily(uint64(chainReport.SourceChainSelector))
		if err != nil {
			return nil, fmt.Errorf("get source chain family: %w", err)
		}

		if chainReport.ProofFlagBits.IsEmpty() {
			return nil, errors.New("proof flag bits are empty")
		}

		evmProofs := make([][32]byte, 0, len(chainReport.Proofs))
		for _, proof := range chainReport.Proofs {
			evmProofs = append(evmProofs, proof)
		}

		evmMessages := make([]offramp.InternalAny2EVMRampMessage, 0, len(chainReport.Messages))
		for _, message := range chainReport.Messages {
			receiver := common.BytesToAddress(message.Receiver)
			tokenAmounts := make([]offramp.InternalAny2EVMTokenTransfer, 0, len(message.TokenAmounts))
			for _, tokenAmount := range message.TokenAmounts {
				if tokenAmount.Amount.IsEmpty() {
					return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
				}

				if tokenAmount.Amount.Sign() < 0 {
					return nil, fmt.Errorf("negative amount for token: %s", tokenAmount.DestTokenAddress)
				}

				destExecDataDecodedMap, err := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
				if err != nil {
					return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
				}

				destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
				if err != nil {
					return nil, fmt.Errorf("decode dest gas amount: %w", err)
				}

				var sourcePoolAddr []byte
				if sourceChainFamily == chainsel.FamilyEVM {
					// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
					// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
					sourcePoolAddr, err = abiEncodeAddress(common.BytesToAddress(tokenAmount.SourcePoolAddress))
					if err != nil {
						return nil, fmt.Errorf("abi encode source pool address: %w", err)
					}
				} else {
					sourcePoolAddr = tokenAmount.SourcePoolAddress
				}

				tokenAmounts = append(tokenAmounts, offramp.InternalAny2EVMTokenTransfer{
					SourcePoolAddress: sourcePoolAddr,
					DestTokenAddress:  common.BytesToAddress(tokenAmount.DestTokenAddress),
					ExtraData:         tokenAmount.ExtraData,
					Amount:            tokenAmount.Amount.Int,
					DestGasAmount:     destGasAmount,
				})
			}

			decodedExtraArgsMap, err := e.extraDataCodec.DecodeExtraArgs(message.ExtraArgs, chainReport.SourceChainSelector)
			if err != nil {
				return nil, err
			}

			gasLimit, err := parseExtraArgsMap(decodedExtraArgsMap)
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
				Sender:       common.LeftPadBytes(message.Sender, 32), // todo: make it chain-agnostic
				Data:         message.Data,
				Receiver:     receiver,
				GasLimit:     gasLimit,
				TokenAmounts: tokenAmounts,
			})
		}

		evmChainReport := offramp.InternalExecutionReport{
			SourceChainSelector: uint64(chainReport.SourceChainSelector),
			Messages:            evmMessages,
			OffchainTokenData:   chainReport.OffchainTokenData,
			Proofs:              evmProofs,
			ProofFlagBits:       chainReport.ProofFlagBits.Int,
		}
		evmReport = append(evmReport, evmChainReport)
	}

	return e.executeReportMethodInputs.PackValues([]any{&evmReport})
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (ccipocr3.ExecutePluginReport, error) {
	unpacked, err := e.executeReportMethodInputs.Unpack(encodedReport)
	if err != nil {
		return ccipocr3.ExecutePluginReport{}, fmt.Errorf("unpack encoded report: %w", err)
	}
	if len(unpacked) != 1 {
		return ccipocr3.ExecutePluginReport{}, errors.New("unpacked report is empty")
	}

	evmReportRaw := abi.ConvertType(unpacked[0], new([]offramp.InternalExecutionReport))
	evmReportPtr, is := evmReportRaw.(*[]offramp.InternalExecutionReport)
	if !is {
		return ccipocr3.ExecutePluginReport{}, fmt.Errorf("got an unexpected report type %T", unpacked[0])
	}
	if evmReportPtr == nil {
		return ccipocr3.ExecutePluginReport{}, errors.New("evm report is nil")
	}

	evmReport := *evmReportPtr
	executeReport := ccipocr3.ExecutePluginReport{
		ChainReports: make([]ccipocr3.ExecutePluginReportSingleChain, 0, len(evmReport)),
	}

	for _, evmChainReport := range evmReport {
		proofs := make([]ccipocr3.Bytes32, 0, len(evmChainReport.Proofs))
		for _, proof := range evmChainReport.Proofs {
			proofs = append(proofs, proof)
		}

		messages := make([]ccipocr3.Message, 0, len(evmChainReport.Messages))
		for _, evmMessage := range evmChainReport.Messages {
			tokenAmounts := make([]ccipocr3.RampTokenAmount, 0, len(evmMessage.TokenAmounts))
			for _, tokenAmount := range evmMessage.TokenAmounts {
				destData, err := abiEncodeUint32(tokenAmount.DestGasAmount)
				if err != nil {
					return ccipocr3.ExecutePluginReport{}, fmt.Errorf("abi encode dest gas amount: %w", err)
				}
				tokenAmounts = append(tokenAmounts, ccipocr3.RampTokenAmount{
					// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
					// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
					// its unclear as of writing how we will handle non-EVM chains and their addresses.
					// e.g, will we encode them as bytes or bytes32?
					SourcePoolAddress: common.BytesToAddress(tokenAmount.SourcePoolAddress).Bytes(),
					// TODO: should this be abi-encoded?
					DestTokenAddress: tokenAmount.DestTokenAddress.Bytes(),
					ExtraData:        tokenAmount.ExtraData,
					Amount:           ccipocr3.NewBigInt(tokenAmount.Amount),
					DestExecData:     destData,
				})
			}

			message := ccipocr3.Message{
				Header: ccipocr3.RampMessageHeader{
					MessageID:           evmMessage.Header.MessageId,
					SourceChainSelector: ccipocr3.ChainSelector(evmMessage.Header.SourceChainSelector),
					DestChainSelector:   ccipocr3.ChainSelector(evmMessage.Header.DestChainSelector),
					SequenceNumber:      ccipocr3.SeqNum(evmMessage.Header.SequenceNumber),
					Nonce:               evmMessage.Header.Nonce,
					MsgHash:             ccipocr3.Bytes32{},        // todo: info not available, but not required atm
					OnRamp:              ccipocr3.UnknownAddress{}, // todo: info not available, but not required atm
				},
				Sender:         evmMessage.Sender,
				Data:           evmMessage.Data,
				Receiver:       evmMessage.Receiver.Bytes(),
				ExtraArgs:      ccipocr3.Bytes{},          // <-- todo: info not available, but not required atm
				FeeToken:       ccipocr3.UnknownAddress{}, // <-- todo: info not available, but not required atm
				FeeTokenAmount: ccipocr3.BigInt{},         // <-- todo: info not available, but not required atm
				TokenAmounts:   tokenAmounts,
			}
			messages = append(messages, message)
		}

		chainReport := ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: ccipocr3.ChainSelector(evmChainReport.SourceChainSelector),
			Messages:            messages,
			OffchainTokenData:   evmChainReport.OffchainTokenData,
			Proofs:              proofs,
			ProofFlagBits:       ccipocr3.NewBigInt(evmChainReport.ProofFlagBits),
		}

		executeReport.ChainReports = append(executeReport.ChainReports, chainReport)
	}

	return executeReport, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ ccipocr3.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
