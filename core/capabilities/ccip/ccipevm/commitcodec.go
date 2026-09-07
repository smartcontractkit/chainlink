package ccipevm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_encoding_utils"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/utils/abihelpers"
)

var ccipEncodingUtilsABI = abihelpers.MustParseABI(ccip_encoding_utils.EncodingUtilsABI)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report ccipocr3.CommitPluginReport) ([]byte, error) {
	isBlessed := make(map[ccipocr3.ChainSelector]bool)
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

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (ccipocr3.CommitPluginReport, error) {
	method, ok := ccipEncodingUtilsABI.Methods["exposeCommitReport"]
	if !ok {
		return ccipocr3.CommitPluginReport{}, errors.New("missing method exposeCommitReport")
	}

	unpacked, err := method.Inputs.Unpack(bytes)
	if err != nil {
		return ccipocr3.CommitPluginReport{}, fmt.Errorf("failed to unpack commit report: %w", err)
	}
	if len(unpacked) != 1 {
		return ccipocr3.CommitPluginReport{}, fmt.Errorf("expected 1 argument, got %d", len(unpacked))
	}

	commitReport := *abi.ConvertType(unpacked[0], new(ccip_encoding_utils.OffRampCommitReport)).(*ccip_encoding_utils.OffRampCommitReport)

	isBlessed := make(map[uint64]bool)
	for _, root := range commitReport.BlessedMerkleRoots {
		isBlessed[root.SourceChainSelector] = true
	}

	blessedMerkleRoots := make([]ccipocr3.MerkleRootChain, 0, len(commitReport.BlessedMerkleRoots))
	unblessedMerkleRoots := make([]ccipocr3.MerkleRootChain, 0, len(commitReport.UnblessedMerkleRoots))
	for _, root := range append(commitReport.BlessedMerkleRoots, commitReport.UnblessedMerkleRoots...) {
		mrc := ccipocr3.MerkleRootChain{
			ChainSel:      ccipocr3.ChainSelector(root.SourceChainSelector),
			OnRampAddress: root.OnRampAddress,
			SeqNumsRange: ccipocr3.NewSeqNumRange(
				ccipocr3.SeqNum(root.MinSeqNr),
				ccipocr3.SeqNum(root.MaxSeqNr),
			),
			MerkleRoot: root.MerkleRoot,
		}
		if isBlessed[root.SourceChainSelector] {
			blessedMerkleRoots = append(blessedMerkleRoots, mrc)
		} else {
			unblessedMerkleRoots = append(unblessedMerkleRoots, mrc)
		}
	}

	tokenPriceUpdates := make([]ccipocr3.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, ccipocr3.TokenPrice{
			TokenID: ccipocr3.UnknownEncodedAddress(update.SourceToken.String()),
			Price:   ccipocr3.NewBigInt(big.NewInt(0).Set(update.UsdPerToken)),
		})
	}

	gasPriceUpdates := make([]ccipocr3.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, ccipocr3.GasPriceChain{
			GasPrice: ccipocr3.NewBigInt(big.NewInt(0).Set(update.UsdPerUnitGas)),
			ChainSel: ccipocr3.ChainSelector(update.DestChainSelector),
		})
	}

	rmnSignatures := make([]ccipocr3.RMNECDSASignature, 0, len(commitReport.RmnSignatures))
	for _, sig := range commitReport.RmnSignatures {
		rmnSignatures = append(rmnSignatures, ccipocr3.RMNECDSASignature{
			R: sig.R,
			S: sig.S,
		})
	}

	return ccipocr3.CommitPluginReport{
		BlessedMerkleRoots:   blessedMerkleRoots,
		UnblessedMerkleRoots: unblessedMerkleRoots,
		PriceUpdates: ccipocr3.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		RMNSignatures: rmnSignatures,
	}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ ccipocr3.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
