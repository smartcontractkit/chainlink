package ccipton

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
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

	if len(report.ChainReports) != 1 {
		return nil, fmt.Errorf("unexpected chain report length: %d", len(report.ChainReports))
	}

	chainReport := report.ChainReports[0]
	if len(chainReport.Messages) > 1 {
		return nil, fmt.Errorf("unexpected report message length: %d", len(chainReport.Messages))
	}

	var offChainTokenData *cell.Cell
	var err error

	if len(chainReport.Messages) > 0 {
		msg := chainReport.Messages[0]
		tokenAmounts := make([]binding.Any2TONTokenTransfer, 0, len(msg.TokenAmounts))
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

			_, err := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
			if err != nil {
				return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
			}

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

			tokenAmounts = append(tokenAmounts, binding.Any2TONTokenTransfer{
				SourcePoolAddress: poolAddrCell,
				ExtraData:         extraData,
				DestPoolAddress:   address.NewAddress(tokenAmount.DestTokenAddress[0], tokenAmount.DestTokenAddress[1], tokenAmount.DestTokenAddress[2:]),
				Amount:            tokenAmount.Amount.Int,
				//DestGasAmount:     destGasAmount,
			})
		}

		offChainTokenData, err = binding.Pack2DByteArrayToCell(chainReport.OffchainTokenData[0])
		if err != nil {
			return nil, fmt.Errorf("pack offchain token data: %w", err)
		}
	}

	message := binding.ExecuteReport{
		SourceChainSelector: uint64(chainReport.SourceChainSelector),
		OffChainTokenData:   offChainTokenData,
	}

	cell, err := tlb.ToCell(message)
	if err != nil {
		return nil, fmt.Errorf("convert message to cell: %w", err)
	}

	return cell.ToBOC(), nil
}
