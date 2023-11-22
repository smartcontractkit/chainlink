package mercury_v1 //nolint:revive

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury"
)

// ValidateCurrentBlock sanity checks number and hash
func ValidateCurrentBlock(rf ReportFields) error {
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
		return errors.Errorf("invalid length for hash; expected %d (got: %d)", mercury.EvmHashLen, len(rf.CurrentBlockHash))
	}

	return nil
}
