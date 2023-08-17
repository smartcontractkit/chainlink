package mercury_v1

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
)

// ValidateCurrentBlock sanity checks number and hash
func ValidateCurrentBlock(paos []ParsedAttributedObservation, f int, validFromBlockNum int64) error {
	if validFromBlockNum < 0 {
		return fmt.Errorf("validFromBlockNum must be >= 0 (got: %d)", validFromBlockNum)
	}
	var newBlockRangePaos []ParsedAttributedObservation
	for _, pao := range paos {
		blockNum, valid := pao.GetCurrentBlockNum()
		if valid && blockNum >= validFromBlockNum {
			newBlockRangePaos = append(newBlockRangePaos, pao)
		}
	}

	if len(newBlockRangePaos) < f+1 {
		s := fmt.Sprintf("only %v/%v attributed observations have currentBlockNum >= validFromBlockNum, need at least f+1 (%v/%v) to make a new report", len(newBlockRangePaos), len(paos), f+1, len(paos))
		_, num, _, err := GetConsensusCurrentBlock(paos, f)
		if err == nil {
			return fmt.Errorf("%s; consensusCurrentBlock=%d, validFromBlockNum=%d", s, num, validFromBlockNum)
		}
		return fmt.Errorf("%s; GetConsensusCurrentBlock failed: %w", s, err)
	}
	hash, num, _, err := GetConsensusCurrentBlock(newBlockRangePaos, f)
	if err != nil {
		return fmt.Errorf("GetConsensusCurrentBlock failed: %w", err)
	}

	if num < 0 {
		return errors.Errorf("block number must be >= 0 (got: %d)", num)
	}

	// NOTE: hardcoded ethereum hash
	if len(hash) != mercury.EvmHashLen {
		return errors.Errorf("invalid length for hash; expected %d (got: %d)", mercury.EvmHashLen, len(hash))
	}

	if validFromBlockNum > num {
		// should be impossible actually due to filtering above, but here for sanity check
		return errors.Errorf("validFromBlockNum (%d) must be less than current block number (%d)", validFromBlockNum, num)
	}

	return nil
}
