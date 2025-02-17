package medianpoc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

var DefaultMultiplier = new(big.Int).SetInt64(1e18)

type DeviationFunctionDefinition struct {
	f median.DeviationFunc
}

func (d DeviationFunctionDefinition) Func() median.DeviationFunc {
	return d.f
}

// UnmarshalJSON  for DeviationFunctionDefinition expects a JSON object with a
// "type" field and additional fields depending on the type.
// e.g. {"type": "pendle", "expiresAt": 1739828109}
func (d *DeviationFunctionDefinition) UnmarshalJSON(data []byte) error {
	// Parse JSON object
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	// Check for type field
	typeVal, ok := obj["type"].(string)
	if !ok {
		return errors.New("missing or invalid 'type' field in deviation function definition")
	}

	switch typeVal {
	case "pendle":
		expiresAt, ok := obj["expiresAt"].(float64) // Assume its a unix TS, it can have fractions of a second
		if !ok {
			return errors.New("missing or invalid 'expiresAt' field in deviation function definition")
		}
		var multiplier *big.Int
		if multiplierStr, ok := obj["multiplier"].(string); ok { // Multiplier could be huge so we use string instead of inaccurate javascript number
			multiplier = new(big.Int)
			if _, ok := multiplier.SetString(multiplierStr, 10); !ok {
				return fmt.Errorf("invalid 'multiplier' field in deviation function definition: %s", multiplierStr)
			}
		} else {
			multiplier = DefaultMultiplier
		}

		d.f = makePendleDeviationFunc(expiresAt, SystemClock{}, multiplier)
		return nil
	default:
		return fmt.Errorf("unsupported function type in deviation function definition: %s", typeVal)
	}
}

const SecondsInYear = float64(365 * 24 * 60 * 60)

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

// makePendleDeviationFunc makes a pendle-specific deviation func
//
// NOTE: This is non-deterministic if clock.Now() is non-deterministic (the usual case)
// expiresAt expected as float64 number of seconds since epoch
func makePendleDeviationFunc(expiresAt float64, clock Clock, valMultiplier *big.Int) median.DeviationFunc {
	valMultiplierF := new(big.Float).SetInt(valMultiplier)
	return func(ctx context.Context, thresholdPPB uint64, oldVal, newVal *big.Int) (bool, error) {
		if oldVal == nil || newVal == nil {
			return false, errors.New("oldVal and newVal must be non-nil")
		}

		nowF64 := float64(clock.Now().UnixNano()) / 1e9
		// Convert expirationSeconds to years
		yearsToExpiration := (expiresAt - nowF64) / SecondsInYear

		// Compute absolute difference |oldVal - newVal|
		diff := newVal.Sub(newVal, oldVal)
		diff.Abs(diff) // Take absolute value
		// Convert big.Int to float64 for calculation
		diffFloat := new(big.Float).SetInt(diff)
		// Divide by multiplier
		diffFloat = diffFloat.Quo(diffFloat, valMultiplierF)
		diffF64, _ := diffFloat.Float64()

		// Compute logarithmic threshold
		logThreshold := math.Log(1 + float64(thresholdPPB)/1e9)

		fmt.Println("---")
		fmt.Println("left > right", (diffF64 * yearsToExpiration), logThreshold, (diffF64*yearsToExpiration) > logThreshold)
		// Return the comparison result
		return (diffF64 * yearsToExpiration) > logThreshold, nil
	}
}
