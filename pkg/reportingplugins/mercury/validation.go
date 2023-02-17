package mercury

import (
	"math/big"

	"github.com/pkg/errors"
)

func ValidateBenchmarkPrice(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusBenchmarkPrice(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return errors.Errorf("median benchmark price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateBid(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusBid(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return errors.Errorf("median bid price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateAsk(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusAsk(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return errors.Errorf("median ask price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateBlockValues(paos []ParsedAttributedObservation, f int, maxFinalizedBlockNumber int64) error {
	_, num, err := GetConsensusCurrentBlock(paos, f)
	if err != nil {
		return errors.Wrap(err, "GetConsensusCurrentBlock failed")
	}
	if num < 0 {
		return errors.Errorf("block number must be >= 0 (got: %d)", num)
	}

	if maxFinalizedBlockNumber >= num {
		return errors.Errorf("maxFinalizedBlockNumber (%d) must be less than current block number (%d)", maxFinalizedBlockNumber, num)
	}

	validFrom, err := GetConsensusValidFromBlock(paos, f)
	if err != nil {
		return errors.Wrap(err, "GetConsensusValidFromBlock failed")
	}

	if validFrom > num {
		return errors.Errorf("validFromBlockNum (%d) must be less than or equal to current block number (%d)", validFrom, num)
	}
	if validFrom < 0 {
		return errors.Errorf("validFromBlockNum must be >= 0 (got: %d)", validFrom)
	}

	return nil
}
