package ccipton

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ag_binary "github.com/gagliardetto/binary"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/binding"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	cell "github.com/xssnick/tonutils-go/tvm/cell"
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
	// support single report and single message for now
	if len(report.ChainReports) == 0 {
		// OCR3 runs in a constant loop and will produce empty reports, so we need to handle this case
		// return an empty report, CCIP will discard it on ShouldAcceptAttestedReport/ShouldTransmitAcceptedReport
		// via validateReport before attempting to decode
		return nil, nil
	}

	tonReports := make([]binding.ExecuteReport, 0, len(report.ChainReports))
	for _, chainReport := range report.ChainReports {
		var offChainTokenData *cell.Cell
		var err error
		rampMessages := make([]binding.Any2TONRampMessage, 0, len(chainReport.Messages))

		for _, msg := range chainReport.Messages {
			tokenAmounts := make([]binding.Any2TONTokenTransfer, 0, len(msg.TokenAmounts))
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

				poolAddrCell, err := binding.PackByteArrayToCell(tokenAmount.SourcePoolAddress)
				if err != nil {
					return nil, fmt.Errorf("pack source pool address: %w", err)
				}

				extraData, err := binding.PackByteArrayToCell(tokenAmount.ExtraData)
				if err != nil {
					return nil, fmt.Errorf("pack extra data: %w", err)
				}

				if len(tokenAmount.DestTokenAddress) < 36 {
					return nil, fmt.Errorf("invalid dest token address length: %d", len(tokenAmount.DestTokenAddress))
				}

				// TODO consider using address codec ?
				destTokenTonAddr, err := convertBase64ToAddress(tokenAmount.DestTokenAddress)
				if err != nil {
					return nil, fmt.Errorf("error convert dest token address: %w", err)
				}

				tokenAmounts = append(tokenAmounts, binding.Any2TONTokenTransfer{
					SourcePoolAddress: poolAddrCell,
					ExtraData:         extraData,
					DestPoolAddress:   destTokenTonAddr,
					Amount:            tokenAmount.Amount.Int,
					DestGasAmount:     destGasAmount,
				})
			}

			tokenAmountsDict, err := binding.PackArrayWithRefChaining(tokenAmounts)
			if err != nil {
				return nil, fmt.Errorf("pack token amounts: %w", err)
			}

			header := binding.RampMessageHeader{
				MessageID:           msg.Header.MessageID[:],
				SourceChainSelector: uint64(msg.Header.SourceChainSelector),
				DestChainSelector:   uint64(msg.Header.DestChainSelector),
				SequenceNumber:      uint64(msg.Header.SequenceNumber),
				Nonce:               msg.Header.Nonce,
			}

			senderAddr, err := binding.PackByteArrayToCell(msg.Sender)
			if err != nil {
				return nil, fmt.Errorf("pack sender address: %w", err)
			}

			dataCell, err := binding.PackByteArrayToCell(msg.Data)
			if err != nil {
				return nil, fmt.Errorf("pack data: %w", err)
			}

			// TODO consider using address codec ?
			tonReceiverAddr, err := convertBase64ToAddress(msg.Receiver)
			if err != nil {
				return nil, fmt.Errorf("error convert receiver address: %w", err)
			}

			extraArgsDecodeMap, err := e.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, chainReport.SourceChainSelector)
			if err != nil {
				return nil, fmt.Errorf("failed to decode extra args: %w", err)
			}

			gasLimitBigInt, err := parseExtraArgsMap(extraArgsDecodeMap)
			if err != nil {
				return nil, fmt.Errorf("parse extra args map to get gas limit: %w", err)
			}

			rampMsg := binding.Any2TONRampMessage{
				Header:       header,
				Sender:       senderAddr,
				Data:         dataCell,
				Receiver:     tonReceiverAddr,
				GasLimit:     tlb.FromNanoTON(gasLimitBigInt), // TODO double check if this match with on-chain decimal
				TokenAmounts: tokenAmountsDict,
			}

			rampMessages = append(rampMessages, rampMsg)
		}

		packedRampMsgsCell, err := binding.PackArrayWithRefChaining(rampMessages)
		if err != nil {
			return nil, fmt.Errorf("pack ramp messages: %w", err)
		}

		if len(chainReport.Messages) > 0 && len(chainReport.OffchainTokenData) > 0 {
			// should only have an offchain token data if there are tokens as part of the message
			offChainTokenData, err = binding.Pack2DByteArrayToCell(chainReport.OffchainTokenData[0])
			if err != nil {
				return nil, fmt.Errorf("pack offchain token data: %w", err)
			}
		}

		sigs := make([]binding.Signature, 0, len(chainReport.Proofs))
		for _, proof := range chainReport.Proofs {
			sigs = append(sigs, binding.Signature{
				Sig: proof[:],
			})
		}

		sigCell, err := binding.PackArrayWithStaticType(sigs)
		if err != nil {
			return nil, fmt.Errorf("pack signatures: %w", err)
		}

		message := binding.ExecuteReport{
			SourceChainSelector: uint64(chainReport.SourceChainSelector),
			OffChainTokenData:   offChainTokenData,
			Messages:            packedRampMsgsCell,
			Proofs:              sigCell,
			ProofFlagBits:       chainReport.ProofFlagBits.Int,
		}

		tonReports = append(tonReports, message)
	}

	chainedReports, err := binding.PackArrayWithRefChaining(tonReports)
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

	unpackedReports, err := binding.UnPackArrayWithRefChaining[binding.ExecuteReport](c)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("unpack execute reports: %w", err)
	}

	executeReport := cciptypes.ExecutePluginReport{
		ChainReports: make([]cciptypes.ExecutePluginReportSingleChain, 0, len(unpackedReports)),
	}

	for _, tonReport := range unpackedReports {
		signatures, err := binding.UnpackArrayWithStaticType[binding.Signature](tonReport.Proofs)
		if err != nil {
			return cciptypes.ExecutePluginReport{}, err
		}

		proofs := make([]cciptypes.Bytes32, 0, len(signatures))
		for _, proof := range signatures {
			proofs = append(proofs, cciptypes.Bytes32(proof.Sig))
		}

		unpackedMsgs, err := binding.UnPackArrayWithRefChaining[binding.Any2TONRampMessage](tonReport.Messages)
		if err != nil {
			return executeReport, fmt.Errorf("unpack ramp messages: %w", err)
		}

		messages := make([]cciptypes.Message, 0, len(unpackedMsgs))
		for _, msg := range unpackedMsgs {
			tonTokenAmounts, err := binding.UnPackArrayWithRefChaining[binding.Any2TONTokenTransfer](msg.TokenAmounts)
			if err != nil {
				return executeReport, fmt.Errorf("unpack token amounts: %w", err)
			}

			tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(tonTokenAmounts))
			for _, tokenAmount := range tonTokenAmounts {
				sourceTokenPoolAddr, err := binding.UnloadCellToByteArray(tokenAmount.SourcePoolAddress)
				if err != nil {
					return executeReport, err
				}

				extraData, err := binding.UnloadCellToByteArray(tokenAmount.ExtraData)
				if err != nil {
					return executeReport, err
				}

				// TODO consider using address codec ?
				destTokenAddr, err := convertAddressToBase64(tokenAmount.DestPoolAddress)
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
					Amount:            cciptypes.NewBigInt(tokenAmount.Amount),
					DestExecData:      destGasAmount,
				})
			}

			senderAddr, err := binding.UnloadCellToByteArray(msg.Sender)
			if err != nil {
				return executeReport, fmt.Errorf("unload sender address: %w", err)
			}

			receiverAddr, err := convertAddressToBase64(msg.Receiver)
			if err != nil {
				return executeReport, fmt.Errorf("convert receiver address: %w", err)
			}

			msgData, err := binding.UnloadCellToByteArray(msg.Data)
			if err != nil {
				return executeReport, fmt.Errorf("unload message data: %w", err)
			}

			// TODO make sure generic
			extraArgs := binding.GenericExtraArgsV2{
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
				Sender:       senderAddr,
				Data:         msgData,
				Receiver:     receiverAddr,
				ExtraArgs:    extraArgsCell.ToBOC(),
				TokenAmounts: tokenAmounts,
			})
		}

		offchainTokenData := make([][][]byte, 0)
		// TODO check if TON will support multiple offchain token data, then change binding to 3DByteArray
		if tonReport.OffChainTokenData != nil {
			offchainData, err := binding.Unpack2DByteArrayFromCell(tonReport.OffChainTokenData)
			if err != nil {
				return executeReport, fmt.Errorf("unload offchain token data: %w", err)
			}
			offchainTokenData = append(offchainTokenData, offchainData)
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

// Convert the raw address bytes to a TON address.
func convertBase64ToAddress(rawAddr []byte) (*address.Address, error) {
	addrStr := base64.RawURLEncoding.EncodeToString(rawAddr)
	tonAddr, err := address.ParseAddr(addrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TON address %s: %w", addrStr, err)
	}

	return tonAddr, nil
}

func convertAddressToBase64(addr *address.Address) ([]byte, error) {
	if addr == nil {
		return nil, fmt.Errorf("nil address")
	}

	addrStr := addr.String()
	if len(addrStr) == 0 {
		return nil, fmt.Errorf("empty address string")
	}
	return base64.RawURLEncoding.DecodeString(addrStr)
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

// TODO could be duplicate from ccipevm, consider moving to common package
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
				// when source chain is svm, the gas limit is an ag_binary.Uint128 struct instead of *big.Int
				if val, ok := fieldValue.(ag_binary.Uint128); ok {
					outputGas = val.BigInt()
					return outputGas, nil
				}
				return nil, fmt.Errorf("unexpected type for gas limit: %T", fieldValue)
			}
		default:
			// no error here, as we only need the keys to gasLimit, other keys can be skipped without like AllowOutOfOrderExecution	etc.
		}
	}
	return outputGas, errors.New("gas limit not found in extra data map")
}
