package ccipton

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

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
		addr, err := address.ParseAddr(string(tpu.TokenID))
		if err != nil {
			return nil, fmt.Errorf("cannot parse token address %s: %w", tpu.TokenID, err)
		}

		if tpu.Price.IsEmpty() {
			return nil, fmt.Errorf("empty token price for token %s", tpu.TokenID)
		}
		tpuSlice[i] = bindings.TokenPriceUpdate{
			SourceToken: addr,
			UsdPerToken: tpu.Price.Int,
		}
	}

	gpuSlice := make([]bindings.GasPriceUpdate, len(report.PriceUpdates.GasPriceUpdates))
	for i, gpu := range report.PriceUpdates.GasPriceUpdates {
		if gpu.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain selector %d", gpu.ChainSel)
		}
		gpuSlice[i] = bindings.GasPriceUpdate{
			DestChainSelector: uint64(gpu.ChainSel),
			UsdPerUnitGas:     gpu.GasPrice.Int,
		}
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

	sigSlice := make([]bindings.Signature, len(report.RMNSignatures))
	for i, sig := range report.RMNSignatures {
		rmnSig64Array := make([]byte, 64)
		copy(rmnSig64Array[:32], sig.R[:])
		copy(rmnSig64Array[32:], sig.S[:])
		sigSlice[i] = bindings.Signature{
			Sig: bytes.Clone(rmnSig64Array[:]),
		}
	}

	cellReport := bindings.CommitReport{
		PriceUpdates: bindings.PriceUpdates{
			TokenPriceUpdates: tpuSlice,
			GasPriceUpdates:   gpuSlice,
		},
		MerkleRoot: bindings.MerkleRoots{
			BlessedMerkleRoots:   mkSlice,
			UnblessedMerkleRoots: unblessedMkSlice,
		},
		RMNSignatures: sigSlice,
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
	var tpuSlice []cciptypes.TokenPrice
	if len(priceUpdate.TokenPriceUpdates) > 0 {
		tpuSlice = make([]cciptypes.TokenPrice, len(priceUpdate.TokenPriceUpdates))
		for i, update := range priceUpdate.TokenPriceUpdates {
			var tokenPrice *big.Int
			if update.UsdPerToken != nil && update.UsdPerToken.Sign() != 0 {
				tokenPrice = update.UsdPerToken
			} else if update.UsdPerToken != nil {
				tokenPrice = big.NewInt(0)
			}
			tpuSlice[i] = cciptypes.TokenPrice{
				TokenID: cciptypes.UnknownEncodedAddress(update.SourceToken.String()),
				Price:   cciptypes.NewBigInt(tokenPrice),
			}
		}
	}

	var gpuSlice []cciptypes.GasPriceChain
	if len(priceUpdate.GasPriceUpdates) > 0 {
		gpuSlice = make([]cciptypes.GasPriceChain, len(priceUpdate.GasPriceUpdates))
		for i, update := range priceUpdate.GasPriceUpdates {
			var gasPrice *big.Int
			if update.UsdPerUnitGas != nil && update.UsdPerUnitGas.Sign() != 0 {
				gasPrice = update.UsdPerUnitGas
			} else if update.UsdPerUnitGas != nil {
				gasPrice = big.NewInt(0)
			}
			gpuSlice[i] = cciptypes.GasPriceChain{
				ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
				GasPrice: cciptypes.NewBigInt(gasPrice),
			}
		}
	}

	var sigSlice []cciptypes.RMNECDSASignature
	if len(report.RMNSignatures) > 0 {
		sigSlice = make([]cciptypes.RMNECDSASignature, len(report.RMNSignatures))
		for i, sig := range report.RMNSignatures {
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
	}

	mr := report.MerkleRoot
	var bmrSlice []cciptypes.MerkleRootChain
	if len(mr.BlessedMerkleRoots) > 0 {
		bmrSlice = make([]cciptypes.MerkleRootChain, len(mr.BlessedMerkleRoots))
		for i, mr := range mr.BlessedMerkleRoots {
			bmrSlice[i] = cciptypes.MerkleRootChain{
				ChainSel:      cciptypes.ChainSelector(mr.SourceChainSelector),
				OnRampAddress: mr.OnRampAddress,
				SeqNumsRange:  cciptypes.NewSeqNumRange(cciptypes.SeqNum(mr.MinSeqNr), cciptypes.SeqNum(mr.MaxSeqNr)),
				MerkleRoot:    cciptypes.Bytes32(mr.MerkleRoot),
			}
		}
	}

	var unblessedMrSlice []cciptypes.MerkleRootChain
	if len(mr.UnblessedMerkleRoots) > 0 {
		unblessedMrSlice = make([]cciptypes.MerkleRootChain, len(mr.UnblessedMerkleRoots))
		for i, mr := range mr.UnblessedMerkleRoots {
			unblessedMrSlice[i] = cciptypes.MerkleRootChain{
				ChainSel:      cciptypes.ChainSelector(mr.SourceChainSelector),
				OnRampAddress: mr.OnRampAddress,
				SeqNumsRange:  cciptypes.NewSeqNumRange(cciptypes.SeqNum(mr.MinSeqNr), cciptypes.SeqNum(mr.MaxSeqNr)),
				MerkleRoot:    cciptypes.Bytes32(mr.MerkleRoot),
			}
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
