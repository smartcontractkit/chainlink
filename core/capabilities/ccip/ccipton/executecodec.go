package ccipton

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/xssnick/tonutils-go/tlb"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
	addressCodec   AddressCodec
	extraDataCodec common.ExtraDataCodec
}

func NewExecutePluginCodecV1(extraDataCodec common.ExtraDataCodec) *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{
		addressCodec:   AddressCodec{},
		extraDataCodec: extraDataCodec,
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	// support single report and single message for now
	if len(report.ChainReports) == 0 {
		// OCR3 runs in a constant loop and will produce empty reports, so we need to handle this case
		// return an empty report, CCIP will discard it on ShouldAcceptAttestedReport/ShouldTransmitAcceptedReport
		// via validateReport before attempting to decode
		return nil, nil
	}

	tonReports := make([]bindings.ExecuteReport, 0, len(report.ChainReports))
	for _, chainReport := range report.ChainReports {
		var offChainTokenData [][]byte
		rampMessages := make([]bindings.Any2TONRampMessage, 0, len(chainReport.Messages))

		for _, msg := range chainReport.Messages {
			tokenAmounts := make([]bindings.Any2TONTokenTransfer, 0, len(msg.TokenAmounts))
			for _, tokenAmount := range msg.TokenAmounts {
				if tokenAmount.Amount.IsEmpty() {
					return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
				}

				if tokenAmount.Amount.Int.Sign() < 0 {
					return nil, fmt.Errorf("negative amount for token: %s", tokenAmount.DestTokenAddress)
				}

				if len(tokenAmount.DestTokenAddress) != 36 {
					return nil, fmt.Errorf("invalid destTokenAddress address: %v", tokenAmount.DestTokenAddress)
				}

				destExecDataDecodedMap, err := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
				if err != nil {
					return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
				}

				destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
				if err != nil {
					return nil, fmt.Errorf("extract dest gas amount: %w", err)
				}

				poolAddrCell, err := bindings.PackByteArrayToCell(tokenAmount.SourcePoolAddress)
				if err != nil {
					return nil, fmt.Errorf("pack source pool address: %w", err)
				}

				extraData, err := bindings.PackByteArrayToCell(tokenAmount.ExtraData)
				if err != nil {
					return nil, fmt.Errorf("pack extra data: %w", err)
				}

				if len(tokenAmount.DestTokenAddress) < 36 {
					return nil, fmt.Errorf("invalid dest token address length: %d", len(tokenAmount.DestTokenAddress))
				}

				destTokenAddrStr, err := e.addressCodec.AddressBytesToString(tokenAmount.DestTokenAddress)
				if err != nil {
					return nil, err
				}

				DestPoolTonAddr, err := address.ParseAddr(destTokenAddrStr)
				if err != nil {
					return nil, fmt.Errorf("invalid dest token address %s: %w", destTokenAddrStr, err)
				}

				tokenAmounts = append(tokenAmounts, bindings.Any2TONTokenTransfer{
					SourcePoolAddress: poolAddrCell,
					ExtraData:         extraData,
					DestPoolAddress:   DestPoolTonAddr,
					Amount:            tokenAmount.Amount.Int,
					DestGasAmount:     destGasAmount,
				})
			}

			header := bindings.RampMessageHeader{
				MessageID:           msg.Header.MessageID[:],
				SourceChainSelector: uint64(msg.Header.SourceChainSelector),
				DestChainSelector:   uint64(msg.Header.DestChainSelector),
				SequenceNumber:      uint64(msg.Header.SequenceNumber),
				Nonce:               msg.Header.Nonce,
			}

			tonReceiverAddrStr, err := e.addressCodec.AddressBytesToString(msg.Receiver)
			if err != nil {
				return nil, fmt.Errorf("error convert receiver address: %w", err)
			}

			tonReceiverAddr, err := address.ParseAddr(tonReceiverAddrStr)
			if err != nil {
				return nil, fmt.Errorf("invalid receiver address %s: %w", tonReceiverAddrStr, err)
			}

			var gasLimitBigInt *big.Int
			if msg.ExtraArgs != nil && len(msg.ExtraArgs) > 0 {
				extraArgsDecodeMap, err := e.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, chainReport.SourceChainSelector)
				if err != nil {
					return nil, fmt.Errorf("failed to decode extra args: %w", err)
				}

				gasLimitBigInt, err = parseExtraArgsMap(extraArgsDecodeMap)
				if err != nil {
					return nil, fmt.Errorf("parse extra args map to get gas limit: %w", err)
				}
			}

			gasLimit, err := tlb.FromNano(gasLimitBigInt, 0)
			if err != nil {
				return nil, fmt.Errorf("convert gas limit to TON cell: %w", err)
			}
			rampMsg := bindings.Any2TONRampMessage{
				Header:       header,
				Sender:       bindings.SnakeBytes(msg.Sender),
				Data:         bindings.SnakeBytes(msg.Data),
				Receiver:     tonReceiverAddr,
				GasLimit:     gasLimit, // TODO double check if this match with on-chain decimal. Note the offramp contract would not use this value base on current design.
				TokenAmounts: tokenAmounts,
			}

			rampMessages = append(rampMessages, rampMsg)
		}

		if len(chainReport.Messages) > 0 && len(chainReport.OffchainTokenData) > 0 {
			// should only have an offchain token data if there are tokens as part of the message
			offChainTokenData = chainReport.OffchainTokenData[0]
		}

		sigs := make([]bindings.Signature, 0, len(chainReport.Proofs))
		for _, proof := range chainReport.Proofs {
			sigs = append(sigs, bindings.Signature{
				Sig: proof[:],
			})
		}

		message := bindings.ExecuteReport{
			SourceChainSelector: uint64(chainReport.SourceChainSelector),
			OffChainTokenData:   offChainTokenData,
			Messages:            rampMessages,
			Proofs:              sigs,
			ProofFlagBits:       chainReport.ProofFlagBits.Int,
		}

		tonReports = append(tonReports, message)
	}

	chainedReports, err := bindings.PackArrayWithRefChaining(tonReports)
	if err != nil {
		return nil, fmt.Errorf("pack execute reports: %w", err)
	}

	return chainedReports.ToBOC(), nil
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, data []byte) (cciptypes.ExecutePluginReport, error) {
	c, err := cell.FromBOC(data)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("decode BOC: %w", err)
	}

	unpackedReports, err := bindings.UnPackArrayWithRefChaining[bindings.ExecuteReport](c)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack execute reports: %w", err)
	}

	executeReport := cciptypes.ExecutePluginReport{
		ChainReports: make([]cciptypes.ExecutePluginReportSingleChain, 0, len(unpackedReports)),
	}

	for _, tonReport := range unpackedReports {
		proofs := make([]cciptypes.Bytes32, 0, len(tonReport.Proofs))
		for _, proof := range tonReport.Proofs {
			proofs = append(proofs, cciptypes.Bytes32(proof.Sig))
		}

		messages := make([]cciptypes.Message, 0, len(tonReport.Messages))
		for _, msg := range tonReport.Messages {
			tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(msg.TokenAmounts))
			for _, tokenAmount := range msg.TokenAmounts {
				sourceTokenPoolAddr, err := bindings.UnloadCellToByteArray(tokenAmount.SourcePoolAddress)
				if err != nil {
					return executeReport, err
				}

				extraData, err := bindings.UnloadCellToByteArray(tokenAmount.ExtraData)
				if err != nil {
					return executeReport, err
				}

				destTokenAddr, err := e.addressCodec.AddressStringToBytes(tokenAmount.DestPoolAddress.String())
				if err != nil {
					return executeReport, err
				}

				// big endian encoding for dest gas amount
				destGasAmount := make([]byte, 4)
				binary.BigEndian.PutUint32(destGasAmount, tokenAmount.DestGasAmount)

				tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
					SourcePoolAddress: sourceTokenPoolAddr,
					DestTokenAddress:  destTokenAddr,
					ExtraData:         extraData,
					Amount:            cciptypes.NewBigInt(tokenAmount.Amount), // TODO double check if we need to add range check for BigInt, since TON use 256 bits
					DestExecData:      destGasAmount,
				})
			}

			receiverAddr, err := e.addressCodec.AddressStringToBytes(msg.Receiver.String())
			if err != nil {
				return executeReport, err
			}

			extraArgs := bindings.GenericExtraArgsV2{
				GasLimit:                 msg.GasLimit.Nano(),
				AllowOutOfOrderExecution: true,
			}

			extraArgsCell, err := tlb.ToCell(extraArgs)
			if err != nil {
				return cciptypes.ExecutePluginReport{}, fmt.Errorf("convert extra args to cell: %w", err)
			}

			messages = append(messages, cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           cciptypes.Bytes32(msg.Header.MessageID),
					SourceChainSelector: cciptypes.ChainSelector(msg.Header.SourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(msg.Header.DestChainSelector),
					SequenceNumber:      cciptypes.SeqNum(msg.Header.SequenceNumber),
					Nonce:               msg.Header.Nonce,
				},
				Sender:       cciptypes.UnknownAddress(msg.Sender),
				Data:         cciptypes.Bytes(msg.Data),
				Receiver:     receiverAddr,
				ExtraArgs:    extraArgsCell.ToBOC(),
				TokenAmounts: tokenAmounts,
			})
		}

		offchainTokenData := make([][][]byte, 0)
		if tonReport.OffChainTokenData != nil {
			offchainTokenData = append(offchainTokenData, tonReport.OffChainTokenData)
		}

		executeReport.ChainReports = append(executeReport.ChainReports, cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(tonReport.SourceChainSelector),
			Messages:            messages,
			OffchainTokenData:   offchainTokenData,
			Proofs:              proofs,
			ProofFlagBits:       cciptypes.NewBigInt(tonReport.ProofFlagBits),
		})
	}

	return executeReport, nil
}

// Duplicate with ccipevm, consider moving to common package
func extractDestGasAmountFromMap(input map[string]any) (uint32, error) {
	// Iterate through the expected fields in the struct
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "destgasamount":
			// Expect uint32
			if val, ok := fieldValue.(uint32); ok {
				return val, nil
			} else {
				return 0, errors.New("invalid type for destgasamount, expected uint32")
			}
		default:
		}
	}

	return 0, errors.New("invalid token message, dest gas amount not found in the DestExecDataDecoded map")
}

// Duplicate with ccipevm, consider moving to common package
func parseExtraArgsMap(input map[string]any) (*big.Int, error) {
	var outputGas *big.Int
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "gaslimit":
			if val, ok := fieldValue.(*big.Int); ok {
				outputGas = val
				return outputGas, nil
			} else {
				return nil, fmt.Errorf("unexpected type for gas limit: %T", fieldValue)
			}
		default:
			// no error here, as we only need the keys to gasLimit, other keys can be skipped without like AllowOutOfOrderExecution	etc.
		}
	}
	return outputGas, errors.New("gas limit not found in extra data map")
}
