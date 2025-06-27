package ccipton

import (
	"context"
	"encoding/base64"
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/binding"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/xssnick/tonutils-go/address"
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

				//_, err := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
				//if err != nil {
				//	return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
				//}

				// TODO pending implementation of dest gas amount extraction, waiting for router contract design
				//destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
				//if err != nil {
				//	return nil, err
				//}

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

				destTokenTonAddr, err := convertBase64ToAddress(tokenAmount.DestTokenAddress)
				if err != nil {
					return nil, fmt.Errorf("error convert dest token address: %w", err)
				}

				tokenAmounts = append(tokenAmounts, binding.Any2TONTokenTransfer{
					SourcePoolAddress: poolAddrCell,
					ExtraData:         extraData,
					DestPoolAddress:   destTokenTonAddr,
					Amount:            tokenAmount.Amount.Int,
					//DestGasAmount:     destGasAmount, // TODO pending implementation
				})

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

				tonReceiverAddr, err := convertBase64ToAddress(msg.Receiver)
				if err != nil {
					return nil, fmt.Errorf("error convert receiver address: %w", err)
				}

				rampMsg := binding.Any2TONRampMessage{
					Header:       header,
					Sender:       senderAddr,
					Data:         dataCell,
					Receiver:     tonReceiverAddr,
					GasLimit:     make([]byte, 32), // TODO: implement gas limit handling with extra data codec
					TokenAmounts: tokenAmountsDict,
				}

				rampMessages = append(rampMessages, rampMsg)
			}
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
				tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{

				})
			}
			messages = append(messages, cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           cciptypes.Bytes32(msg.Header.MessageID),
					SourceChainSelector: cciptypes.ChainSelector(msg.Header.SourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(msg.Header.DestChainSelector),
					SequenceNumber:      cciptypes.SeqNum(msg.Header.SequenceNumber),
					Nonce:               msg.Header.Nonce,
				},
				Sender:       ,
				Data:         ,
				Receiver:     ,
				// ExtraArgs: ,// TODO: implement gas limit handling with extra data codec
				TokenAmounts: ,
			})
		}
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
