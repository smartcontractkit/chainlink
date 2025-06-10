package ccipton

import (
	"bytes"
	"context"
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (cr *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	tpuSlice := make([]bindings.TokenPriceUpdate, len(report.PriceUpdates.TokenPriceUpdates))
	for i, tpu := range report.PriceUpdates.TokenPriceUpdates {
		tpuSlice[i] = bindings.TokenPriceUpdate{
			SourceToken: address.MustParseAddr(string(tpu.TokenID)),
			UsdPerToken: tpu.Price.Int,
		}
	}
	tpu, err := bindings.SliceToDict(tpuSlice)
	if err != nil {
		return nil, fmt.Errorf("cannot encode token price updates: %w", err)
	}

	gpuSlice := make([]bindings.GasPriceUpdate, len(report.PriceUpdates.GasPriceUpdates))
	for i, gpu := range report.PriceUpdates.GasPriceUpdates {
		gpuSlice[i] = bindings.GasPriceUpdate{
			DestChainSelector: uint64(gpu.ChainSel),
			UsdPerUnitGas:     gpu.GasPrice.Int,
		}
	}
	gpu, err := bindings.SliceToDict(gpuSlice)
	if err != nil {
		return nil, fmt.Errorf("cannot encode gas price updates: %w", err)
	}

	mkSlice := make([]bindings.MerkleRoot, len(report.BlessedMerkleRoots))
	for i, mr := range report.BlessedMerkleRoots {
		mkSlice[i] = bindings.MerkleRoot{
			SourceChainSelector: uint64(mr.ChainSel),
			OnRampAddress:       mr.OnRampAddress,
			MinSeqNr:            uint64(mr.SeqNumsRange.Start()),
			MaxSeqNr:            uint64(mr.SeqNumsRange.End()),
			MerkleRoot:          bytes.Clone(mr.MerkleRoot[:]),
		}
	}
	merkleRoots, err := bindings.SliceToDict(mkSlice)
	if err != nil {
		return nil, fmt.Errorf("cannot encode blessed merkle roots: %w", err)
	}

	unblessedMkSlice := make([]bindings.MerkleRoot, len(report.UnblessedMerkleRoots))
	for i, mr := range report.UnblessedMerkleRoots {
		unblessedMkSlice[i] = bindings.MerkleRoot{
			SourceChainSelector: uint64(mr.ChainSel),
			OnRampAddress:       mr.OnRampAddress,
			MinSeqNr:            uint64(mr.SeqNumsRange.Start()),
			MaxSeqNr:            uint64(mr.SeqNumsRange.End()),
			MerkleRoot:          bytes.Clone(mr.MerkleRoot[:]),
		}
	}
	unblessedRoots, err := bindings.SliceToDict(unblessedMkSlice)
	if err != nil {
		return nil, fmt.Errorf("cannot encode unblessed merkle roots: %w", err)
	}

	sigSlice := make([]bindings.Signature, len(report.RMNSignatures))
	for i, sig := range report.RMNSignatures {
		rmnSig64Array := make([]byte, 64)
		copy(rmnSig64Array[:32], sig.R[:])
		copy(rmnSig64Array[32:], sig.S[:])
		sigSlice[i] = bindings.Signature{
			Sig: bytes.Clone(rmnSig64Array[:]),
		}
	}
	signatures, err := bindings.SliceToDict(sigSlice)
	if err != nil {
		return nil, fmt.Errorf("cannot encode RMN signatures: %w", err)
	}

	cellReport := bindings.CommitReport{
		PriceUpdates: bindings.PriceUpdates{
			TokenPriceUpdates: tpu,
			GasPriceUpdates:   gpu,
		},
		MerkleRoot: bindings.MerkleRoots{
			BlessedMerkleRoots:   merkleRoots,
			UnblessedMerkleRoots: unblessedRoots,
		},
		RMNSignatures: signatures,
	}

	c, err := tlb.ToCell(cellReport)
	if err != nil {
		return nil, fmt.Errorf("cannot encode commit report to cell: %w", err)
	}

	// Serialize the cell to bytes
	return c.ToBOC(), nil
}

func (cr *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	c, err := cell.FromBOC(bytes)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode BOC: %w", err)
	}

	var report bindings.CommitReport
	if err := tlb.LoadFromCell(&report, c.BeginParse()); err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode commit report from cell: %w", err)
	}

	priceUpdate := report.PriceUpdates
	tpu, err := bindings.DictToSlice[bindings.TokenPriceUpdate](priceUpdate.TokenPriceUpdates)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode token price updates: %w", err)
	}

	tpuSlice := make([]cciptypes.TokenPrice, len(tpu))
	for i, update := range tpu {
		tpuSlice[i] = cciptypes.TokenPrice{
			TokenID: cciptypes.UnknownEncodedAddress(update.SourceToken.String()),
			Price:   cciptypes.NewBigInt(update.UsdPerToken),
		}
	}
	gpu, err := bindings.DictToSlice[bindings.GasPriceUpdate](priceUpdate.GasPriceUpdates)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode gas price updates: %w", err)
	}

	gpuSlice := make([]cciptypes.GasPriceChain, len(gpu))
	for i, update := range gpu {
		gpuSlice[i] = cciptypes.GasPriceChain{
			ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
			GasPrice: cciptypes.NewBigInt(update.UsdPerUnitGas),
		}
	}

	sigs, err := bindings.DictToSlice[bindings.Signature](report.RMNSignatures)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode RMN signatures: %w", err)
	}

	sigSlice := make([]cciptypes.RMNECDSASignature, len(sigs))
	for i, sig := range sigs {
		if len(sig.Sig) != 64 {
			return cciptypes.CommitPluginReport{}, fmt.Errorf("invalid RMN signature length: %d", len(sig.Sig))
		}

		var r, s [32]byte
		copy(r[:], sig.Sig[:32])
		copy(s[:], sig.Sig[32:])
		sigSlice[i] = cciptypes.RMNECDSASignature{
			R: r,
			S: s,
		}
	}

	mr := report.MerkleRoot
	bmr, err := bindings.DictToSlice[bindings.MerkleRoot](mr.BlessedMerkleRoots)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode blessed merkle roots: %w", err)
	}

	bmrSlice := make([]cciptypes.MerkleRootChain, len(bmr))
	for i, mr := range bmr {
		bmrSlice[i] = cciptypes.MerkleRootChain{
			ChainSel:      cciptypes.ChainSelector(mr.SourceChainSelector),
			OnRampAddress: mr.OnRampAddress,
			SeqNumsRange:  cciptypes.NewSeqNumRange(cciptypes.SeqNum(mr.MinSeqNr), cciptypes.SeqNum(mr.MaxSeqNr)),
			MerkleRoot:    cciptypes.Bytes32(mr.MerkleRoot),
		}
	}

	unblessedMr, err := bindings.DictToSlice[bindings.MerkleRoot](mr.UnblessedMerkleRoots)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("cannot decode unblessed merkle roots: %w", err)
	}
	unblessedMrSlice := make([]cciptypes.MerkleRootChain, len(unblessedMr))
	for i, mr := range unblessedMr {
		unblessedMrSlice[i] = cciptypes.MerkleRootChain{
			ChainSel:      cciptypes.ChainSelector(mr.SourceChainSelector),
			OnRampAddress: mr.OnRampAddress,
			SeqNumsRange:  cciptypes.NewSeqNumRange(cciptypes.SeqNum(mr.MinSeqNr), cciptypes.SeqNum(mr.MaxSeqNr)),
			MerkleRoot:    cciptypes.Bytes32(mr.MerkleRoot),
		}
	}

	return cciptypes.CommitPluginReport{
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tpuSlice,
			GasPriceUpdates:   gpuSlice,
		},
		BlessedMerkleRoots:   bmrSlice,
		UnblessedMerkleRoots: unblessedMrSlice,
		RMNSignatures:        sigSlice,
	}, nil
}
