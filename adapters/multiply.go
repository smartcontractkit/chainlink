package adapters

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times interface{} `json:"times"`
}

// Perform returns the input's "value" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Get("value")
	if err != nil {
		return input.WithError(err)
	}

	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return input.WithError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}

	switch times := ma.Times.(type) {
	case int:
		res := i.Mul(i, big.NewFloat(float64(times)))
		return input.WithValue(res.String())

	case string:
		timesInt, err := strconv.Atoi(times)
		if err != nil {
			return input.WithError(fmt.Errorf("cannot parse into int: %v", times))
		}

		res := i.Mul(i, big.NewFloat(float64(timesInt)))
		return input.WithValue(res.String())

	default:
		return input.WithError(fmt.Errorf("wrong type of the multiplier: %v", times))
	}
}
