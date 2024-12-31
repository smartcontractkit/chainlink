package ccipsolana

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	mr := ccip_router.MerkleRoot{}

	if len(report.MerkleRoots) != 0 {
		mr = ccip_router.MerkleRoot{
			SourceChainSelector: uint64(report.MerkleRoots[0].ChainSel),
			OnRampAddress:       report.MerkleRoots[0].OnRampAddress,
			MinSeqNr:            uint64(report.MerkleRoots[0].SeqNumsRange.Start()),
			MaxSeqNr:            uint64(report.MerkleRoots[0].SeqNumsRange.End()),
			MerkleRoot:          report.MerkleRoots[0].MerkleRoot,
		}
	}

	tpu := make([]ccip_router.TokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		token, err := solana.PublicKeyFromBase58(string(update.TokenID))
		if err != nil {
			return nil, fmt.Errorf("invalid token address: %s, %v", update.TokenID, err)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tpu = append(tpu, ccip_router.TokenPriceUpdate{
			SourceToken: token,
			UsdPerToken: common.To28BytesBE(update.Price.Int.Uint64()),
		})
	}

	gpu := make([]ccip_router.GasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		gpu = append(gpu, ccip_router.GasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     common.To28BytesBE(update.GasPrice.Int.Uint64()),
		})
	}

	commit := ccip_router.CommitInput{
		MerkleRoot: mr,
		PriceUpdates: ccip_router.PriceUpdates{
			TokenPriceUpdates: tpu,
			GasPriceUpdates:   gpu,
		},
	}

	err := commit.MarshalWithEncoder(encoder)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	decoder := agbinary.NewBorshDecoder(bytes)
	commitReport := ccip_router.CommitInput{}
	err := commitReport.UnmarshalWithDecoder(decoder)
	if err != nil {
		return cciptypes.CommitPluginReport{}, err
	}

	merkleRoots := []cciptypes.MerkleRootChain{
		{
			ChainSel:      cciptypes.ChainSelector(commitReport.MerkleRoot.SourceChainSelector),
			OnRampAddress: commitReport.MerkleRoot.OnRampAddress,
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(commitReport.MerkleRoot.MinSeqNr),
				cciptypes.SeqNum(commitReport.MerkleRoot.MaxSeqNr),
			),
			MerkleRoot: commitReport.MerkleRoot.MerkleRoot,
		},
	}

	tokenPriceUpdates := make([]cciptypes.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			TokenID: cciptypes.UnknownEncodedAddress(update.SourceToken.String()),
			Price:   priceHelper(update.UsdPerToken[:]),
		})
	}

	gasPriceUpdates := make([]cciptypes.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, cciptypes.GasPriceChain{
			GasPrice: priceHelper(update.UsdPerUnitGas[:]),
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

func priceHelper(input []byte) cciptypes.BigInt {
	var tokenPrice cciptypes.BigInt
	price := new(big.Int).SetBytes(input)
	if price.Int64() == 0 {
		tokenPrice = cciptypes.NewBigInt(big.NewInt(0))
	} else {
		tokenPrice = cciptypes.NewBigInt(price)
	}

	return tokenPrice
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
