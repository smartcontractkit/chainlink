package adapters

import (
	"errors"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/tidwall/gjson"
)

// Condition adapter type takes an Operator and a Value field to
// compare to the previous task's Result.
type Condition struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

var (
	ErrResultNotNumber      = errors.New("The result was not a number")
	ErrValueNotNumber       = errors.New("The value was not a number")
	ErrOperatorNotSpecified = errors.New("Operator not specified")
	ErrValueNotSpecified    = errors.New("Value not specified")
)

// Perform uses the Operator to check the run's result against the
// specified Value.
func (c *Condition) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	prevResult := input.Result()

	if c.Value == "" {
		input.SetError(ErrValueNotSpecified)
		return input
	}

	switch c.Operator {
	case "eq":
		if c.Value == prevResult.String() {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	case "neq":
		if c.Value != prevResult.String() {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	case "gt":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			input.SetError(err)
			return input
		}
		if desired < value {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	case "gte":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			input.SetError(err)
			return input
		}
		if desired <= value {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	case "lt":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			input.SetError(err)
			return input
		}
		if desired > value {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	case "lte":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			input.SetError(err)
			return input
		}
		if desired >= value {
			input.CompleteWithResult(true)
		} else {
			input.CompleteWithResult(false)
		}
	default:
		input.SetError(ErrOperatorNotSpecified)
	}

	return input
}

func getValues(result gjson.Result, d string) (float64, float64, error) {
	value, err := strconv.ParseFloat(result.String(), 64)
	if err != nil {
		return 0, 0, ErrResultNotNumber
	}
	desired, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0, 0, ErrValueNotNumber
	}
	return value, desired, nil
}
