package mercury

import (
	"errors"
	"fmt"
	"math/big"

	pkgerrors "github.com/pkg/errors"
)

func ValidateBenchmarkPrice(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusBenchmarkPrice(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median benchmark price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateBid(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusBid(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median bid price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateAsk(paos []ParsedAttributedObservation, min, max *big.Int) error {
	answer := GetConsensusAsk(paos)

	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return pkgerrors.Errorf("median ask price %s is outside of allowable range (Min: %s, Max: %s)", answer, min, max)
	}

	return nil
}

func ValidateBlockValues(paos []ParsedAttributedObservation, f int, maxFinalizedBlockNumber int64) error {
	var newBlockRangePaos []ParsedAttributedObservation
	for _, pao := range paos {
		if pao.CurrentBlockNum > pao.ValidFromBlockNum {
			newBlockRangePaos = append(newBlockRangePaos, pao)
		}
	}

	if !(f+1 <= len(newBlockRangePaos)) {
		s := fmt.Sprintf("only %v/%v attributed observations have currentBlockNum > validFromBlockNum, need at least f+1 (%v/%v) to make a new report; this is most likely a duplicate report for the block range", len(newBlockRangePaos), len(paos), f+1, len(paos))
		_, currentBlockNum, err := GetConsensusCurrentBlock(paos, f)
		validFromBlockNum, err2 := GetConsensusValidFromBlock(paos, f)
		err = errors.Join(err, err2)
		if err == nil {
			err = pkgerrors.Errorf("%s; consensusCurrentBlock=%d, consensusValidFromBlock=%d", s, currentBlockNum, validFromBlockNum)
		} else {
			err = pkgerrors.Errorf("%s; could not come to consensus about block numbers: %v", s, err)
		}
		return err
	}

	_, num, err := GetConsensusCurrentBlock(paos, f)
	if err != nil {
		return pkgerrors.Wrap(err, "GetConsensusCurrentBlock failed")
	}
	if num < 0 {
		return pkgerrors.Errorf("block number must be >= 0 (got: %d)", num)
	}

	if maxFinalizedBlockNumber >= num {
		return pkgerrors.Errorf("maxFinalizedBlockNumber (%d) must be less than current block number (%d)", maxFinalizedBlockNumber, num)
	}

	validFrom, err := GetConsensusValidFromBlock(paos, f)
	if err != nil {
		return pkgerrors.Wrap(err, "GetConsensusValidFromBlock failed")
	}

	// Shouldn't be possible but leave here as a sanity check
	if validFrom > num {
		return pkgerrors.Errorf("validFromBlockNum (%d) must be less than or equal to current block number (%d)", validFrom, num)
	}
	if validFrom < 0 {
		return pkgerrors.Errorf("validFromBlockNum must be >= 0 (got: %d)", validFrom)
	}

	return nil
}
