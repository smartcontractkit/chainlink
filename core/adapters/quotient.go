package adapters

import (
	"fmt"
	"math/big"
	"strconv"

	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

// Dividend represents x where x / y.
type Dividend float64

// UnmarshalJSON implements json.Unmarshaler.
func (n *Dividend) UnmarshalJSON(input []byte) error {
	input = utils.RemoveQuotes(input)
	dividend, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return fmt.Errorf("cannot parse into float: %s", input)
	}

	*n = Dividend(dividend)

	return nil
}

// Quotient holds the Dividend.
type Quotient struct {
	Dividend *Dividend `json:"dividend"`
}

// Perform returns result of dividend / divisor were divisor is
// the input's "result" field.
//
// For example, if input value is "2.5", and the adapter's "dividend" value
// is "1", the result's value will be "0.4".
func (q *Quotient) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	fmt.Println(q.Dividend)
	val := input.Result()
	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.NewRunOutputError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}
	if i.Cmp(big.NewFloat(0)) == 0 {
		return models.NewRunOutputError(fmt.Errorf("cannot divide by zero"))
	}
	if q.Dividend != nil {
		i = new(big.Float).Quo(big.NewFloat(float64(*q.Dividend)), i)
	}
	return models.NewRunOutputCompleteWithResult(i.String())
}
