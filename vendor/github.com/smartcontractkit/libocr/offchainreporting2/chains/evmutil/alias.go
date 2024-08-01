// alias for offchainreporting2plus/chains/evmutil
package evmutil

import (
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func SplitSignature(sig []byte) (r, s [32]byte, v byte, err error) {
	return evmutil.SplitSignature(sig)
}

func RawReportContext(repctx types.ReportContext) [3][32]byte {
	return evmutil.RawReportContext(repctx)
}

func ContractConfigFromConfigSetEvent(changed ocr2aggregator.OCR2AggregatorConfigSet) types.ContractConfig {
	return evmutil.ContractConfigFromConfigSetEvent(changed)
}

type EVMOffchainConfigDigester = evmutil.EVMOffchainConfigDigester
