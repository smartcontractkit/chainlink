package ccipevm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "EVM2EVMMultiOffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct {
	commitReportAcceptedEventInputs abi.Arguments
}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	abiParsed, err := abi.JSON(strings.NewReader(evm_2_evm_multi_offramp.EVM2EVMMultiOffRampABI))
	if err != nil {
		panic(fmt.Errorf("parse multi offramp abi: %s", err))
	}
	eventInputs := abihelpers.MustGetEventInputs("CommitReportAccepted", abiParsed)
	return &CommitPluginCodecV1{commitReportAcceptedEventInputs: eventInputs}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	merkleRoots := make([]evm_2_evm_multi_offramp.EVM2EVMMultiOffRampMerkleRoot, 0, len(report.MerkleRoots))
	for _, root := range report.MerkleRoots {
		merkleRoots = append(merkleRoots, evm_2_evm_multi_offramp.EVM2EVMMultiOffRampMerkleRoot{
			SourceChainSelector: uint64(root.ChainSel),
			Interval: evm_2_evm_multi_offramp.EVM2EVMMultiOffRampInterval{
				Min: uint64(root.SeqNumsRange.Start()),
				Max: uint64(root.SeqNumsRange.End()),
			},
			MerkleRoot: root.MerkleRoot,
		})
	}

	tokenPriceUpdates := make([]evm_2_evm_multi_offramp.InternalTokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		if !common.IsHexAddress(string(update.TokenID)) {
			return nil, fmt.Errorf("invalid token address: %s", update.TokenID)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tokenPriceUpdates = append(tokenPriceUpdates, evm_2_evm_multi_offramp.InternalTokenPriceUpdate{
			SourceToken: common.HexToAddress(string(update.TokenID)),
			UsdPerToken: update.Price.Int,
		})
	}

	gasPriceUpdates := make([]evm_2_evm_multi_offramp.InternalGasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gasPriceUpdates = append(gasPriceUpdates, evm_2_evm_multi_offramp.InternalGasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     update.GasPrice.Int,
		})
	}

	evmReport := evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReport{
		PriceUpdates: evm_2_evm_multi_offramp.InternalPriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		MerkleRoots: merkleRoots,
	}

	return c.commitReportAcceptedEventInputs.PackValues([]interface{}{evmReport})
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	unpacked, err := c.commitReportAcceptedEventInputs.Unpack(bytes)
	if err != nil {
		return cciptypes.CommitPluginReport{}, err
	}
	if len(unpacked) != 1 {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("expected 1 argument, got %d", len(unpacked))
	}

	commitReportRaw := abi.ConvertType(unpacked[0], new(evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReport))
	commitReport, is := commitReportRaw.(*evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReport)
	if !is {
		return cciptypes.CommitPluginReport{},
			fmt.Errorf("expected EVM2EVMMultiOffRampCommitReport, got %T", unpacked[0])
	}

	merkleRoots := make([]cciptypes.MerkleRootChain, 0, len(commitReport.MerkleRoots))
	for _, root := range commitReport.MerkleRoots {
		merkleRoots = append(merkleRoots, cciptypes.MerkleRootChain{
			ChainSel: cciptypes.ChainSelector(root.SourceChainSelector),
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(root.Interval.Min),
				cciptypes.SeqNum(root.Interval.Max),
			),
			MerkleRoot: root.MerkleRoot,
		})
	}

	tokenPriceUpdates := make([]cciptypes.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			TokenID: types.Account(update.SourceToken.String()),
			Price:   cciptypes.NewBigInt(big.NewInt(0).Set(update.UsdPerToken)),
		})
	}

	gasPriceUpdates := make([]cciptypes.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, cciptypes.GasPriceChain{
			GasPrice: cciptypes.NewBigInt(big.NewInt(0).Set(update.UsdPerUnitGas)),
			ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
		})
	}

	return cciptypes.CommitPluginReport{
		MerkleRoots: merkleRoots,
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
	}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
