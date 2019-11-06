// +build !sgx_enclave

package adapters

import (
	"fmt"
	"math/big"
	"strconv"

	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

// Multiplier represents the number to multiply by in Multiply adapter.
type Multiplier float64

// UnmarshalJSON implements json.Unmarshaler.
func (m *Multiplier) UnmarshalJSON(input []byte) error {
	input = utils.RemoveQuotes(input)
	times, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return fmt.Errorf("cannot parse into float: %s", input)
	}

	*m = Multiplier(times)

	return nil
}

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times *Multiplier `json:"times"`
}

// Perform returns the input's "result" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	val := input.Result()
	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.NewRunOutputError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}

	if ma.Times != nil {
		i.Mul(i, big.NewFloat(float64(*ma.Times)))
	}
	return models.NewRunOutputCompleteWithResult(i.String())
}
