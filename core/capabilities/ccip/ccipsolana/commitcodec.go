package ccipsolana

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
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
	combinedRoots := append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...)
	if len(combinedRoots) != 1 {
		return nil, fmt.Errorf("unexpected merkle root length in report: %d", len(combinedRoots))
	}

	merkleRoot := combinedRoots[0]
	mr := ccip_router.MerkleRoot{
		SourceChainSelector: uint64(merkleRoot.ChainSel),
		OnRampAddress:       merkleRoot.OnRampAddress,
		MinSeqNr:            uint64(merkleRoot.SeqNumsRange.Start()),
		MaxSeqNr:            uint64(merkleRoot.SeqNumsRange.End()),
		MerkleRoot:          merkleRoot.MerkleRoot,
	}

	tpu := make([]ccip_router.TokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		token, err := solana.PublicKeyFromBase58(string(update.TokenID))
		if err != nil {
			return nil, fmt.Errorf("invalid token address: %s, %w", update.TokenID, err)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		tpu = append(tpu, ccip_router.TokenPriceUpdate{
			SourceToken: token,
			UsdPerToken: [28]uint8(encodeBigIntToFixedLengthLE(update.Price.Int, 28)),
		})
	}

	gpu := make([]ccip_router.GasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gpu = append(gpu, ccip_router.GasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     [28]uint8(encodeBigIntToFixedLengthLE(update.GasPrice.Int, 28)),
		})
	}

	commit := ccip_router.CommitInput{
		MerkleRoot: mr,
		PriceUpdates: ccip_router.PriceUpdates{
			TokenPriceUpdates: tpu,
			GasPriceUpdates:   gpu,
		},
	}

	if len(report.RMNSignatures) > 1 {
		return nil, fmt.Errorf("Multiple RMNSignatures in report: %d", len(report.RMNSignatures))
	} else if len(report.RMNSignatures) == 1 {
		if report.BlessedMerkleRoots == nil {
			return nil, fmt.Errorf("RMN signature included without a blessed root")
		}
		var rmnSig64Array [64]uint8
		copy(rmnSig64Array[:32], report.RMNSignatures[0].R[:])
		copy(rmnSig64Array[32:], report.RMNSignatures[0].S[:])
		commit.RmnSignatures = [][64]uint8{rmnSig64Array}
	} else if report.UnblessedMerkleRoots == nil {
		return nil, fmt.Errorf("No RMN signature included for the blessed root")
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
			Price:   decodeLEToBigInt(update.UsdPerToken[:]),
		})
	}

	gasPriceUpdates := make([]cciptypes.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, cciptypes.GasPriceChain{
			GasPrice: decodeLEToBigInt(update.UsdPerUnitGas[:]),
			ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
		})
	}

	commitPluginReport := cciptypes.CommitPluginReport{
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
	}

	if len(commitReport.RmnSignatures) == 0 {
		commitPluginReport.UnblessedMerkleRoots = merkleRoots
	} else {
		commitPluginReport.BlessedMerkleRoots = merkleRoots
		rmnSigs := make([]cciptypes.RMNECDSASignature, 0, len(commitReport.RmnSignatures))
		for _, sig := range commitReport.RmnSignatures {
			var r [32]byte
			copy(r[:], sig[:32])
			var s [32]byte
			copy(s[:], sig[32:])
			rmnSigs = append(rmnSigs, cciptypes.RMNECDSASignature{
				R: r,
				S: s,
			})
		}
		commitPluginReport.RMNSignatures = rmnSigs
	}

	return commitPluginReport, nil
}

func encodeBigIntToFixedLengthLE(bi *big.Int, length int) []byte {
	// Create a fixed-length byte array
	paddedBytes := make([]byte, length)

	// Use FillBytes to fill the array with big-endian data, zero-padded
	bi.FillBytes(paddedBytes)

	// Reverse the array for little-endian encoding
	for i, j := 0, len(paddedBytes)-1; i < j; i, j = i+1, j-1 {
		paddedBytes[i], paddedBytes[j] = paddedBytes[j], paddedBytes[i]
	}

	return paddedBytes
}

func decodeLEToBigInt(data []byte) cciptypes.BigInt {
	// Reverse the byte array to convert it from little-endian to big-endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	// Use big.Int.SetBytes to construct the big.Int
	bi := new(big.Int).SetBytes(data)
	if bi.Int64() == 0 {
		return cciptypes.NewBigInt(big.NewInt(0))
	}

	return cciptypes.NewBigInt(bi)
}

func encodedRMNSignature(signature cciptypes.Bytes32) [64]uint8 {
	rmnSig64 := make([]uint8, 0, 64)

	// RMN signature on SVM onchain is 64 bytes long.
	// So left pad 32 bytes with 0s, and then 32 bytes of signature
	rmnSig64 = append(rmnSig64[:32], signature[:]...)
	var rmnSig64Array [64]uint8
	copy(rmnSig64Array[:], rmnSig64[:64])
	return rmnSig64Array
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
