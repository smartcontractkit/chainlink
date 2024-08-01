package v1

import (
	"fmt"

	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"

	"github.com/smartcontractkit/chainlink-data-streams/mercury"
)

// ValidateCurrentBlock sanity checks number and hash
func ValidateCurrentBlock(rf v1.ReportFields) error {
	if rf.ValidFromBlockNum < 0 {
		return fmt.Errorf("validFromBlockNum must be >= 0 (got: %d)", rf.ValidFromBlockNum)
	}
	if rf.CurrentBlockNum < 0 {
		return fmt.Errorf("currentBlockNum must be >= 0 (got: %d)", rf.ValidFromBlockNum)
	}
	if rf.ValidFromBlockNum > rf.CurrentBlockNum {
		return fmt.Errorf("validFromBlockNum (Value: %d) must be less than or equal to CurrentBlockNum (Value: %d)", rf.ValidFromBlockNum, rf.CurrentBlockNum)
	}
	// NOTE: hardcoded ethereum hash
	if len(rf.CurrentBlockHash) != mercury.EvmHashLen {
		return fmt.Errorf("invalid length for hash; expected %d (got: %d)", mercury.EvmHashLen, len(rf.CurrentBlockHash))
	}

	return nil
}
