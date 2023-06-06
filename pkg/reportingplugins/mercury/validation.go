package mercury

import (
	"fmt"
	"math/big"

	pkgerrors "github.com/pkg/errors"
)

// NOTE: hardcoded for now, this may need to change if we support block range on chains other than eth
const evmHashLen = 32

func ValidateBenchmarkPrice(paos []ParsedAttributedObservation, f int, min, max *big.Int) error {
	answer, err := GetConsensusBenchmarkPrice(paos, f)
	if err != nil {
		return err
	}

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median benchmark price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateBid(paos []ParsedAttributedObservation, f int, min, max *big.Int) error {
	answer, err := GetConsensusBid(paos, f)
	if err != nil {
		return err
	}

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median bid price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateAsk(paos []ParsedAttributedObservation, f int, min, max *big.Int) error {
	answer, err := GetConsensusAsk(paos, f)
	if err != nil {
		return err
	}

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median ask price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateCurrentBlock(paos []ParsedAttributedObservation, f int, validFromBlockNum int64) error {
	if validFromBlockNum < 0 {
		return fmt.Errorf("validFromBlockNum must be >= 0 (got: %d)", validFromBlockNum)
	}
	var newBlockRangePaos []ParsedAttributedObservation
	for _, pao := range paos {
		if pao.CurrentBlockValid && pao.CurrentBlockNum >= validFromBlockNum {
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
		return pkgerrors.Errorf("block number must be >= 0 (got: %d)", num)
	}

	// NOTE: hardcoded ethereum hash
	if len(hash) != evmHashLen {
		return pkgerrors.Errorf("invalid length for hash; expected %d (got: %d)", evmHashLen, len(hash))
	}

	if validFromBlockNum > num {
		// should be impossible actually due to filtering above, but here for sanity check
		return pkgerrors.Errorf("validFromBlockNum (%d) must be less than current block number (%d)", validFromBlockNum, num)
	}

	return nil
}
