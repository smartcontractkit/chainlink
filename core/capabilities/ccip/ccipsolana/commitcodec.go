package ccipsolana

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report ccipocr3common.CommitPluginReport) ([]byte, error) {
	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	combinedRoots := report.BlessedMerkleRoots
	combinedRoots = append(combinedRoots, report.UnblessedMerkleRoots...)
	var mr *ccip_offramp.MerkleRoot
	switch len(combinedRoots) {
	case 0:
		// price updates only, zero the root
	case 1:
		// valid
		merkleRoot := combinedRoots[0]
		mr = &ccip_offramp.MerkleRoot{
			SourceChainSelector: uint64(merkleRoot.ChainSel),
			OnRampAddress:       merkleRoot.OnRampAddress,
			MinSeqNr:            uint64(merkleRoot.SeqNumsRange.Start()),
			MaxSeqNr:            uint64(merkleRoot.SeqNumsRange.End()),
			MerkleRoot:          merkleRoot.MerkleRoot,
		}

	default:
		return nil, fmt.Errorf("unexpected merkle root length in report: %d", len(combinedRoots))
	}

	tpu := make([]ccip_offramp.TokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		token, err := solana.PublicKeyFromBase58(string(update.TokenID))
		if err != nil {
			return nil, fmt.Errorf("invalid token address: %s, %w", update.TokenID, err)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tpu = append(tpu, ccip_offramp.TokenPriceUpdate{
			SourceToken: token,
			UsdPerToken: [28]uint8(encodeBigIntToFixedLengthBE(update.Price.Int, 28)),
		})
	}

	gpu := make([]ccip_offramp.GasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gpu = append(gpu, ccip_offramp.GasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     [28]uint8(encodeBigIntToFixedLengthBE(update.GasPrice.Int, 28)),
		})
	}

	commit := ccip_offramp.CommitInput{
		MerkleRoot: mr,
		PriceUpdates: ccip_offramp.PriceUpdates{
			TokenPriceUpdates: tpu,
			GasPriceUpdates:   gpu,
		},
	}

	// Only validate if we actually have a root
	if len(combinedRoots) > 0 {
		switch len(report.RMNSignatures) {
		case 0:
			if len(report.UnblessedMerkleRoots) == 0 {
				return nil, errors.New("no RMN signature included for the blessed root")
			}
		case 1:
			if len(report.BlessedMerkleRoots) == 0 {
				return nil, errors.New("RMN signature included without a blessed root")
			}
			// R part goes into leading 32 bytes, and S part goes into the trailing 32 bytes.
			var rmnSig64Array [64]uint8
			copy(rmnSig64Array[:32], report.RMNSignatures[0].R[:])
			copy(rmnSig64Array[32:], report.RMNSignatures[0].S[:])
			commit.RmnSignatures = [][64]uint8{rmnSig64Array}
		default:
			return nil, fmt.Errorf("multiple RMNSignatures in report: %d", len(report.RMNSignatures))
		}
	}

	err := commit.MarshalWithEncoder(encoder)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (ccipocr3common.CommitPluginReport, error) {
	decoder := agbinary.NewBorshDecoder(bytes)
	commitReport := ccip_offramp.CommitInput{}
	err := commitReport.UnmarshalWithDecoder(decoder)
	if err != nil {
		return ccipocr3common.CommitPluginReport{}, err
	}

	var merkleRoots []ccipocr3common.MerkleRootChain
	if commitReport.MerkleRoot != nil {
		merkleRoots = []ccipocr3common.MerkleRootChain{
			{
				ChainSel:      ccipocr3common.ChainSelector(commitReport.MerkleRoot.SourceChainSelector),
				OnRampAddress: commitReport.MerkleRoot.OnRampAddress,
				SeqNumsRange: ccipocr3common.NewSeqNumRange(
					ccipocr3common.SeqNum(commitReport.MerkleRoot.MinSeqNr),
					ccipocr3common.SeqNum(commitReport.MerkleRoot.MaxSeqNr),
				),
				MerkleRoot: commitReport.MerkleRoot.MerkleRoot,
			},
		}
	}

	// tokenPrice and gasPrice data is big endian encoded, following EVM

	tokenPriceUpdates := make([]ccipocr3common.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, ccipocr3common.TokenPrice{
			TokenID: ccipocr3common.UnknownEncodedAddress(update.SourceToken.String()),
			Price:   decodeBEToBigInt(update.UsdPerToken[:]),
		})
	}

	gasPriceUpdates := make([]ccipocr3common.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, ccipocr3common.GasPriceChain{
			GasPrice: decodeBEToBigInt(update.UsdPerUnitGas[:]),
			ChainSel: ccipocr3common.ChainSelector(update.DestChainSelector),
		})
	}

	commitPluginReport := ccipocr3common.CommitPluginReport{
		PriceUpdates: ccipocr3common.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
	}

	if len(commitReport.RmnSignatures) == 0 {
		commitPluginReport.UnblessedMerkleRoots = merkleRoots
	} else {
		commitPluginReport.BlessedMerkleRoots = merkleRoots
		rmnSigs := make([]ccipocr3common.RMNECDSASignature, 0, len(commitReport.RmnSignatures))
		for _, sig := range commitReport.RmnSignatures {
			// Leading 32 bytes are the R part, and trailing 32 bytes are the S part
			var r [32]byte
			copy(r[:], sig[:32])
			var s [32]byte
			copy(s[:], sig[32:])
			rmnSigs = append(rmnSigs, ccipocr3common.RMNECDSASignature{
				R: r,
				S: s,
			})
		}
		commitPluginReport.RMNSignatures = rmnSigs
	}

	return commitPluginReport, nil
}

func encodeBigIntToFixedLengthBE(bi *big.Int, length int) []byte {
	// Create a fixed-length byte array
	paddedBytes := make([]byte, length)

	// Use FillBytes to fill the array with big-endian data, zero-padded
	bi.FillBytes(paddedBytes)

	return paddedBytes
}

func decodeBEToBigInt(data []byte) ccipocr3common.BigInt {
	// Use big.Int.SetBytes to construct the big.Int
	bi := new(big.Int).SetBytes(data)
	if bi.Cmp(big.NewInt(0)) == 0 {
		return ccipocr3common.NewBigInt(big.NewInt(0))
	}

	return ccipocr3common.NewBigInt(bi)
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ ccipocr3common.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
