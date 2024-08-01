package mercury

import (
	"fmt"
	"math/big"
)

// NOTE: hardcoded for now, this may need to change if we support block range on chains other than eth
const EvmHashLen = 32

// ValidateBetween checks that value is between min and max
func ValidateBetween(name string, answer *big.Int, min, max *big.Int) error {
	if answer == nil {
		return fmt.Errorf("%s: got nil value", name)
	}
	if !(min.Cmp(answer) <= 0 && answer.Cmp(max) <= 0) {
		return fmt.Errorf("%s (Value: %s) is outside of allowable range (Min: %s, Max: %s)", name, answer, min, max)
	}

	return nil
}

func ValidateValidFromTimestamp(observationTimestamp uint32, validFromTimestamp uint32) error {
	if observationTimestamp < validFromTimestamp {
		return fmt.Errorf("observationTimestamp (Value: %d) must be >= validFromTimestamp (Value: %d)", observationTimestamp, validFromTimestamp)
	}

	return nil
}

func ValidateExpiresAt(observationTimestamp uint32, expiresAt uint32) error {
	if observationTimestamp > expiresAt {
		return fmt.Errorf("expiresAt (Value: %d) must be ahead of observation timestamp (Value: %d)", expiresAt, observationTimestamp)
	}

	return nil
}

func ValidateFee(name string, answer *big.Int) error {
	return ValidateBetween(name, answer, big.NewInt(0), MaxInt192)
}
