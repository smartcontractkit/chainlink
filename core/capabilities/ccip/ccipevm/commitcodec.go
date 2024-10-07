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

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct {
	commitReportAcceptedEventInputs abi.Arguments
}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	abiParsed, err := abi.JSON(strings.NewReader(offramp.OffRampABI))
	if err != nil {
		panic(fmt.Errorf("parse multi offramp abi: %s", err))
	}
	eventInputs := abihelpers.MustGetEventInputs("CommitReportAccepted", abiParsed)
	return &CommitPluginCodecV1{commitReportAcceptedEventInputs: eventInputs}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	merkleRoots := make([]offramp.InternalMerkleRoot, 0, len(report.MerkleRoots))
	for _, root := range report.MerkleRoots {
		merkleRoots = append(merkleRoots, offramp.InternalMerkleRoot{
			SourceChainSelector: uint64(root.SourceChainSelector),
			OnRampAddress:       root.OnRampAddress,
			MinSeqNr:            uint64(root.MinSeqNr),
			MaxSeqNr:            uint64(root.MaxSeqNr),
			MerkleRoot:          root.MerkleRoot,
		})
	}

	tokenPriceUpdates := make([]offramp.InternalTokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		if !common.IsHexAddress(string(update.TokenID)) {
			return nil, fmt.Errorf("invalid token address: %s", update.TokenID)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tokenPriceUpdates = append(tokenPriceUpdates, offramp.InternalTokenPriceUpdate{
			SourceToken: common.HexToAddress(string(update.TokenID)),
			UsdPerToken: update.Price.Int,
		})
	}

	gasPriceUpdates := make([]offramp.InternalGasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gasPriceUpdates = append(gasPriceUpdates, offramp.InternalGasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     update.GasPrice.Int,
		})
	}

	priceUpdates := offramp.InternalPriceUpdates{
		TokenPriceUpdates: tokenPriceUpdates,
		GasPriceUpdates:   gasPriceUpdates,
	}

	return c.commitReportAcceptedEventInputs.PackValues([]interface{}{merkleRoots, priceUpdates})
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	unpacked, err := c.commitReportAcceptedEventInputs.Unpack(bytes)
	if err != nil {
		return cciptypes.CommitPluginReport{}, err
	}
	if len(unpacked) != 2 {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("expected 2 arguments, got %d", len(unpacked))
	}

	merkleRootsRaw := abi.ConvertType(unpacked[0], new([]offramp.InternalMerkleRoot))
	priceUpdatesRaw := abi.ConvertType(unpacked[1], new(offramp.InternalPriceUpdates))
	var commitReport offramp.OffRampCommitReportAccepted

	roots, is := merkleRootsRaw.(*[]offramp.InternalMerkleRoot)
	if !is {
		return cciptypes.CommitPluginReport{},
			fmt.Errorf("expected []InternalMerkleRoot, got %T", unpacked[0])
	}
	commitReport.MerkleRoots = *roots

	updates, is := priceUpdatesRaw.(*offramp.InternalPriceUpdates)
	if !is {
		return cciptypes.CommitPluginReport{},
			fmt.Errorf("expected InternalPriceUpdates, got %T", unpacked[1])
	}
	commitReport.PriceUpdates = *updates

	merkleRoots := make([]cciptypes.MerkleRoot, 0, len(commitReport.MerkleRoots))
	for _, root := range commitReport.MerkleRoots {
		merkleRoots = append(merkleRoots, cciptypes.MerkleRoot{
			SourceChainSelector: cciptypes.ChainSelector(root.SourceChainSelector),
			OnRampAddress:       root.OnRampAddress,
			MinSeqNr:            cciptypes.SeqNum(root.MinSeqNr),
			MaxSeqNr:            cciptypes.SeqNum(root.MaxSeqNr),
			MerkleRoot:          root.MerkleRoot,
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
