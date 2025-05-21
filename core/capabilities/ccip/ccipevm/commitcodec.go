package ccipevm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_encoding_utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	ccipEncodingUtilsABI = abihelpers.MustParseABI(ccip_encoding_utils.EncodingUtilsABI)
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	isBlessed := make(map[cciptypes.ChainSelector]bool)
	for _, root := range report.BlessedMerkleRoots {
		isBlessed[root.ChainSel] = true
	}

	blessedMerkleRoots := make([]ccip_encoding_utils.InternalMerkleRoot, 0, len(report.BlessedMerkleRoots))
	unblessedMerkleRoots := make([]ccip_encoding_utils.InternalMerkleRoot, 0, len(report.UnblessedMerkleRoots))
	for _, root := range append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...) {
		imr := ccip_encoding_utils.InternalMerkleRoot{
			SourceChainSelector: uint64(root.ChainSel),
			// TODO: abi-encoded address for EVM source, figure out what to do for non-EVM.
			OnRampAddress: common.LeftPadBytes(root.OnRampAddress, 32),
			MinSeqNr:      uint64(root.SeqNumsRange.Start()),
			MaxSeqNr:      uint64(root.SeqNumsRange.End()),
			MerkleRoot:    root.MerkleRoot,
		}
		if isBl, ok := isBlessed[root.ChainSel]; ok && isBl {
			blessedMerkleRoots = append(blessedMerkleRoots, imr)
		} else {
			unblessedMerkleRoots = append(unblessedMerkleRoots, imr)
		}
	}

	rmnSignatures := make([]ccip_encoding_utils.IRMNRemoteSignature, 0, len(report.RMNSignatures))
	for _, sig := range report.RMNSignatures {
		rmnSignatures = append(rmnSignatures, ccip_encoding_utils.IRMNRemoteSignature{
			R: sig.R,
			S: sig.S,
		})
	}

	tokenPriceUpdates := make([]ccip_encoding_utils.InternalTokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		if !common.IsHexAddress(string(update.TokenID)) {
			return nil, fmt.Errorf("invalid token address: %s", update.TokenID)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tokenPriceUpdates = append(tokenPriceUpdates, ccip_encoding_utils.InternalTokenPriceUpdate{
			SourceToken: common.HexToAddress(string(update.TokenID)),
			UsdPerToken: update.Price.Int,
		})
	}

	gasPriceUpdates := make([]ccip_encoding_utils.InternalGasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gasPriceUpdates = append(gasPriceUpdates, ccip_encoding_utils.InternalGasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     update.GasPrice.Int,
		})
	}

	priceUpdates := ccip_encoding_utils.InternalPriceUpdates{
		TokenPriceUpdates: tokenPriceUpdates,
		GasPriceUpdates:   gasPriceUpdates,
	}

	commitReport := &ccip_encoding_utils.OffRampCommitReport{
		PriceUpdates:         priceUpdates,
		BlessedMerkleRoots:   blessedMerkleRoots,
		UnblessedMerkleRoots: unblessedMerkleRoots,
		RmnSignatures:        rmnSignatures,
	}

	packed, err := ccipEncodingUtilsABI.Pack("exposeCommitReport", commitReport)
	if err != nil {
		return nil, fmt.Errorf("failed to pack commit report: %w", err)
	}

	return packed[4:], nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	method, ok := ccipEncodingUtilsABI.Methods["exposeCommitReport"]
	if !ok {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("missing method exposeCommitReport")
	}

	unpacked, err := method.Inputs.Unpack(bytes)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to unpack commit report: %w", err)
	}
	if len(unpacked) != 1 {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("expected 1 argument, got %d", len(unpacked))
	}

	commitReport := *abi.ConvertType(unpacked[0], new(ccip_encoding_utils.OffRampCommitReport)).(*ccip_encoding_utils.OffRampCommitReport)

	isBlessed := make(map[uint64]bool)
	for _, root := range commitReport.BlessedMerkleRoots {
		isBlessed[root.SourceChainSelector] = true
	}

	blessedMerkleRoots := make([]cciptypes.MerkleRootChain, 0, len(commitReport.BlessedMerkleRoots))
	unblessedMerkleRoots := make([]cciptypes.MerkleRootChain, 0, len(commitReport.UnblessedMerkleRoots))
	for _, root := range append(commitReport.BlessedMerkleRoots, commitReport.UnblessedMerkleRoots...) {
		mrc := cciptypes.MerkleRootChain{
			ChainSel:      cciptypes.ChainSelector(root.SourceChainSelector),
			OnRampAddress: root.OnRampAddress,
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(root.MinSeqNr),
				cciptypes.SeqNum(root.MaxSeqNr),
			),
			MerkleRoot: root.MerkleRoot,
		}
		if isBlessed[root.SourceChainSelector] {
			blessedMerkleRoots = append(blessedMerkleRoots, mrc)
		} else {
			unblessedMerkleRoots = append(unblessedMerkleRoots, mrc)
		}
	}

	tokenPriceUpdates := make([]cciptypes.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			TokenID: cciptypes.UnknownEncodedAddress(update.SourceToken.String()),
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

	rmnSignatures := make([]cciptypes.RMNECDSASignature, 0, len(commitReport.RmnSignatures))
	for _, sig := range commitReport.RmnSignatures {
		rmnSignatures = append(rmnSignatures, cciptypes.RMNECDSASignature{
			R: sig.R,
			S: sig.S,
		})
	}

	return cciptypes.CommitPluginReport{
		BlessedMerkleRoots:   blessedMerkleRoots,
		UnblessedMerkleRoots: unblessedMerkleRoots,
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		RMNSignatures: rmnSignatures,
	}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
