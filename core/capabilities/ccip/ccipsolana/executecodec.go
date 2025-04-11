package ccipsolana

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
	extraDataCodec common.ExtraDataCodec
}

func NewExecutePluginCodecV1(extraDataCodec common.ExtraDataCodec) *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{
		extraDataCodec: extraDataCodec,
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	if len(report.ChainReports) == 0 {
		// OCR3 runs in a constant loop and will produce empty reports, so we need to handle this case
		// return an empty report, CCIP will discard it on ShouldAcceptAttestedReport/ShouldTransmitAcceptedReport
		// via validateReport before attempting to decode
		return nil, nil
	}

	if len(report.ChainReports) != 1 {
		return nil, fmt.Errorf("unexpected chain report length: %d", len(report.ChainReports))
	}

	chainReport := report.ChainReports[0]
	if len(chainReport.Messages) > 1 {
		return nil, fmt.Errorf("unexpected report message length: %d", len(chainReport.Messages))
	}

	var message ccip_offramp.Any2SVMRampMessage
	var offChainTokenData [][]byte
	if len(chainReport.Messages) > 0 {
		// currently only allow executing one message at a time
		msg := chainReport.Messages[0]
		tokenAmounts := make([]ccip_offramp.Any2SVMTokenTransfer, 0, len(msg.TokenAmounts))
		for _, tokenAmount := range msg.TokenAmounts {
			if tokenAmount.Amount.IsEmpty() {
				return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
			}

			if tokenAmount.Amount.Int.Sign() < 0 {
				return nil, fmt.Errorf("negative amount for token: %s", tokenAmount.DestTokenAddress)
			}

			if len(tokenAmount.DestTokenAddress) != solana.PublicKeyLength {
				return nil, fmt.Errorf("invalid destTokenAddress address: %v", tokenAmount.DestTokenAddress)
			}

			destExecDataDecodedMap, err := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
			if err != nil {
				return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
			}

			destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
			if err != nil {
				return nil, err
			}

			if solana.PublicKeyLength != len(tokenAmount.DestTokenAddress) {
				return nil, fmt.Errorf("invalid DestTokenAddress length: %d", len(tokenAmount.DestTokenAddress))
			}
			tokenAmounts = append(tokenAmounts, ccip_offramp.Any2SVMTokenTransfer{
				SourcePoolAddress: tokenAmount.SourcePoolAddress,
				DestTokenAddress:  solana.PublicKeyFromBytes(tokenAmount.DestTokenAddress),
				ExtraData:         tokenAmount.ExtraData,
				Amount:            ccip_offramp.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(tokenAmount.Amount.Int, 32))},
				DestGasAmount:     destGasAmount,
			})
		}

		extraDataDecodedMap, err := e.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, chainReport.SourceChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to decode extra args: %w", err)
		}

		ed, err := parseExtraDataMap(extraDataDecodedMap)
		if err != nil {
			return nil, fmt.Errorf("invalid extra args map: %w", err)
		}

		if len(msg.Receiver) != solana.PublicKeyLength {
			return nil, fmt.Errorf("invalid receiver address: %v", msg.Receiver)
		}

		message = ccip_offramp.Any2SVMRampMessage{
			Header: ccip_offramp.RampMessageHeader{
				MessageId:           msg.Header.MessageID,
				SourceChainSelector: uint64(msg.Header.SourceChainSelector),
				DestChainSelector:   uint64(msg.Header.DestChainSelector),
				SequenceNumber:      uint64(msg.Header.SequenceNumber),
				Nonce:               msg.Header.Nonce,
			},
			Sender:        msg.Sender,
			Data:          msg.Data,
			TokenReceiver: ed.tokenReceiver,
			TokenAmounts:  tokenAmounts,
			ExtraArgs:     ed.extraArgs,
		}

		// should only have an offchain token data if there are tokens as part of the message
		if len(chainReport.OffchainTokenData) > 0 {
			offChainTokenData = chainReport.OffchainTokenData[0]
		}
	}

	solanaProofs := make([][32]byte, 0, len(chainReport.Proofs))
	for _, proof := range chainReport.Proofs {
		solanaProofs = append(solanaProofs, proof)
	}

	solanaReport := ccip_offramp.ExecutionReportSingleChain{
		SourceChainSelector: uint64(chainReport.SourceChainSelector),
		Message:             message,
		OffchainTokenData:   offChainTokenData,
		Proofs:              solanaProofs,
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
	executeReport := ccip_offramp.ExecutionReportSingleChain{}
	err := executeReport.UnmarshalWithDecoder(decoder)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack encoded report: %w", err)
	}

	tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(executeReport.Message.TokenAmounts))
	for _, tokenAmount := range executeReport.Message.TokenAmounts {
		destData := make([]byte, 4)
		binary.LittleEndian.PutUint32(destData, tokenAmount.DestGasAmount)

		tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
			SourcePoolAddress: tokenAmount.SourcePoolAddress,
			DestTokenAddress:  tokenAmount.DestTokenAddress.Bytes(),
			ExtraData:         tokenAmount.ExtraData,
			Amount:            decodeLEToBigInt(tokenAmount.Amount.LeBytes[:]),
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
			Receiver:       executeReport.Message.TokenReceiver.Bytes(),
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

func extractDestGasAmountFromMap(input map[string]any) (uint32, error) {
	// Search for the gas fields
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "destgasamount":
			// Expect uint32
			if v, ok := fieldValue.(uint32); ok {
				return v, nil
			} else {
				return 0, errors.New("invalid type for destgasamount, expected uint32")
			}
		default:

		}
	}

	return 0, errors.New("invalid token message, dest gas amount not found in the DestExecDataDecoded map")
}

func encodeBigIntToFixedLengthLE(bi *big.Int, length int) []byte {
	// Create a fixed-length byte array
	paddedBytes := make([]byte, length)

	// Use FillBytes to fill the array with big-endian data, zero-padded
	bi.FillBytes(paddedBytes)

	// Reverse the array for little-endian encoding
	for i, j := 0, len(paddedBytes)-1; i < j; i, j = i+1, j-1 {
		paddedBytes[i], paddedBytes[j] = paddedBytes[j], paddedBytes[i]
	}

	return paddedBytes
}

func decodeLEToBigInt(data []byte) cciptypes.BigInt {
	// Avoid modifying original data
	buf := make([]byte, len(data))
	copy(buf, data)

	// Reverse the byte array to convert it from little-endian to big-endian
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	// Use big.Int.SetBytes to construct the big.Int
	bi := new(big.Int).SetBytes(buf)
	if bi.Cmp(big.NewInt(0)) == 0 {
		return cciptypes.NewBigInt(big.NewInt(0))
	}

	return cciptypes.NewBigInt(bi)
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
