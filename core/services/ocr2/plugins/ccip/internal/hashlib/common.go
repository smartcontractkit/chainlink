package hashlib

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// BytesOfBytesKeccak will compute a keccak256 hash of the provided bytes of bytes slice
func BytesOfBytesKeccak(b [][]byte) ([32]byte, error) {
	if len(b) == 0 {
		return [32]byte{}, nil
	}

	h := utils.Keccak256Fixed(b[0])
	for _, v := range b[1:] {
		h = utils.Keccak256Fixed(append(h[:], v...))
	}
	return h, nil
}

// LeavesFromIntervals Extracts the hashed leaves from a given set of logs
func LeavesFromIntervals(
	lggr logger.Logger,
	interval commit_store.CommitStoreInterval,
	hasher LeafHasherInterface[[32]byte],
	sendReqs []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested],
) ([][32]byte, error) {
	var seqNrs []uint64
	for _, req := range sendReqs {
		seqNrs = append(seqNrs, req.Data.Message.SequenceNumber)
	}

	if !ccipcalc.ContiguousReqs(lggr, interval.Min, interval.Max, seqNrs) {
		return nil, errors.Errorf("do not have full range [%v, %v] have %v", interval.Min, interval.Max, seqNrs)
	}
	var leaves [][32]byte

	for _, sendReq := range sendReqs {
		hash, err2 := hasher.HashLeaf(sendReq.Data.Raw)
		if err2 != nil {
			return nil, err2
		}
		leaves = append(leaves, hash)
	}

	return leaves, nil
}
