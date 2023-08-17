package mercury

import (
	"math"
	"math/big"

	pkgerrors "github.com/pkg/errors"
)

// NOTE: hardcoded for now, this may need to change if we support block range on chains other than eth
const EvmHashLen = 32

// ValidateBenchmarkPrice checks that value is between min and max
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

// ValidateBid checks that value is between min and max
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

// ValidateAsk checks that value is between min and max
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

func ValidateValidFromTimestamp(observationTimestamp uint32, validFromTimestamp uint32) error {
	if observationTimestamp < validFromTimestamp {
		return pkgerrors.Errorf("observationTimestamp (%d) must be >= validFromTimestamp (%d)", observationTimestamp, validFromTimestamp)
	}

	return nil
}

func ValidateExpiresAt(observationTimestamp uint32, expirationWindow uint32) error {
	if int64(observationTimestamp)+int64(expirationWindow) > math.MaxUint32 {
		return pkgerrors.Errorf("timestamp %d + expiration window %d overflows uint32", observationTimestamp, expirationWindow)
	}

	return nil
}
