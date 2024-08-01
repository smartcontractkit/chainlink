package evmutil

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func SplitSignature(sig []byte) (r, s [32]byte, v byte, err error) {
	if len(sig) != 65 {
		return r, s, v, fmt.Errorf("SplitSignature: wrong size")
	}
	r = common.BytesToHash(sig[:32])
	s = common.BytesToHash(sig[32:64])
	v = sig[64]
	return r, s, v, nil
}

func RawReportContext(repctx types.ReportContext) [3][32]byte {
	rawRepctx := [3][32]byte{}
	copy(rawRepctx[0][:], repctx.ConfigDigest[:])
	binary.BigEndian.PutUint32(rawRepctx[1][32-5:32-1], repctx.Epoch)
	rawRepctx[1][31] = repctx.Round
	rawRepctx[2] = repctx.ExtraHash
	return rawRepctx
}

func ContractConfigFromConfigSetEvent(changed ocr2aggregator.OCR2AggregatorConfigSet) types.ContractConfig {
	transmitAccounts := []types.Account{}
	for _, addr := range changed.Transmitters {
		transmitAccounts = append(transmitAccounts, types.Account(addr.Hex()))
	}
	signers := []types.OnchainPublicKey{}
	for _, addr := range changed.Signers {
		addr := addr
		signers = append(signers, types.OnchainPublicKey(addr[:]))
	}
	return types.ContractConfig{
		changed.ConfigDigest,
		changed.ConfigCount,
		signers,
		transmitAccounts,
		changed.F,
		changed.OnchainConfig,
		changed.OffchainConfigVersion,
		changed.OffchainConfig,
	}
}
