package ccipsolana

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

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
	if len(report.ChainReports) != 1 {
		return nil, fmt.Errorf("unexpected chain report length: %d", len(report.ChainReports))
	}

	chainReport := report.ChainReports[0]
	if len(chainReport.Messages) > 1 {
		return nil, fmt.Errorf("unexpected report message length: %d", len(chainReport.Messages))
	}

	var message ccip_router.Any2SVMRampMessage
	var offChainTokenData [][]byte
	if len(chainReport.Messages) > 0 {
		// currently only allow executing one message at a time
		msg := chainReport.Messages[0]
		tokenAmounts := make([]ccip_router.Any2SVMTokenTransfer, 0, len(msg.TokenAmounts))
		for _, tokenAmount := range msg.TokenAmounts {
			if tokenAmount.Amount.IsEmpty() {
				return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
			}

			if len(tokenAmount.DestTokenAddress) != solana.PublicKeyLength {
				return nil, fmt.Errorf("invalid destTokenAddress address: %v", tokenAmount.DestTokenAddress)
			}

			tokenAmounts = append(tokenAmounts, ccip_router.Any2SVMTokenTransfer{
				SourcePoolAddress: tokenAmount.SourcePoolAddress,
				DestTokenAddress:  solana.PublicKeyFromBytes(tokenAmount.DestTokenAddress),
				ExtraData:         tokenAmount.ExtraData,
				Amount:            ccip_router.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(tokenAmount.Amount.Int, 32))},
				DestGasAmount:     bytesToUint32LE(tokenAmount.DestExecData),
			})
		}

		var extraArgs ccip_router.SVMExtraArgs
		extraArgs, err := parseExtraArgsMap(msg.ExtraArgsDecoded)
		if err != nil {
			return nil, fmt.Errorf("invalid extra args map: %w", err)
		}

		if len(msg.Receiver) != solana.PublicKeyLength {
			return nil, fmt.Errorf("invalid receiver address: %v", msg.Receiver)
		}

		message = ccip_router.Any2SVMRampMessage{
			Header: ccip_router.RampMessageHeader{
				MessageId:           msg.Header.MessageID,
				SourceChainSelector: uint64(msg.Header.SourceChainSelector),
				DestChainSelector:   uint64(msg.Header.DestChainSelector),
				SequenceNumber:      uint64(msg.Header.SequenceNumber),
				Nonce:               msg.Header.Nonce,
			},
			Sender:        msg.Sender,
			Data:          msg.Data,
			LogicReceiver: solana.PublicKeyFromBytes(msg.Receiver),
			// TODO how to get the TokenReceiver ?
			//TokenReceiver: solana.PublicKeyFromBytes(msg.Receiver),
			TokenAmounts: tokenAmounts,
			ExtraArgs:    extraArgs,
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

	solanaReport := ccip_router.ExecutionReportSingleChain{
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
	executeReport := ccip_router.ExecutionReportSingleChain{}
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
			Receiver:       executeReport.Message.LogicReceiver.Bytes(),
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

func parseExtraArgsMap(input map[string]any) (ccip_router.SVMExtraArgs, error) {
	// Parse input map into SolanaExtraArgs
	var out ccip_router.SVMExtraArgs

	// Iterate through the expected fields in the struct
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "computeunits":
			// Expect uint32
			if v, ok := fieldValue.(uint32); ok {
				out.ComputeUnits = v
			} else {
				return out, errors.New("invalid type for ComputeUnits, expected uint32")
			}
		case "iswritablebitmap":
			// Expect uint64
			if v, ok := fieldValue.(uint64); ok {
				out.IsWritableBitmap = v
			} else {
				return out, errors.New("invalid type for IsWritableBitmap, expected uint64")
			}
		case "accounts":
			// Expect [][32]byte
			if v, ok := fieldValue.([][32]byte); ok {
				accounts := make([]solana.PublicKey, len(v))
				for i, val := range v {
					accounts[i] = solana.PublicKeyFromBytes(val[:])
				}
				out.Accounts = accounts
			} else {
				return out, errors.New("invalid type for Accounts, expected [][32]byte")
			}
		}
	}

	return out, nil
}

func bytesToUint32LE(b []byte) uint32 {
	if len(b) < 4 {
		var padded [4]byte
		copy(padded[:len(b)], b) // Pad from the right for little-endian
		return binary.LittleEndian.Uint32(padded[:])
	}

	return binary.LittleEndian.Uint32(b)
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
