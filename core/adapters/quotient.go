package adapters

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Quotient holds the Dividend.
type Quotient struct {
	Dividend *big.Float `json:"-"`
}

// TaskType returns the type of Adapter.
func (q *Quotient) TaskType() models.TaskType {
	return TaskTypeQuotient
}

type jsonQuotient struct {
	Dividend *utils.BigFloat `json:"dividend,omitempty"`
}

// MarshalJSON implements the json.Marshal interface.
func (q Quotient) MarshalJSON() ([]byte, error) {
	jsonObj := jsonQuotient{(*utils.BigFloat)(q.Dividend)}
	return json.Marshal(jsonObj)
}

// UnmarshalJSON implements the json.Unqrshal interface.
func (q *Quotient) UnmarshalJSON(buf []byte) error {
	var jsonObj jsonQuotient
	err := json.Unmarshal(buf, &jsonObj)
	if err != nil {
		return err
	}
	q.Dividend = jsonObj.Dividend.Value()
	return nil
}

// Perform returns result of dividend / divisor were divisor is
// the input's "result" field.
//
// For example, if input value is "2.5", and the adapter's "dividend" value
// is "1", the result's value will be "0.4".
func (q *Quotient) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	val := input.Result()
	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.NewRunOutputError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}
	if i.Cmp(big.NewFloat(0)) == 0 {
		return models.NewRunOutputError(fmt.Errorf("cannot divide by zero"))
	}
	if q.Dividend != nil {
		i = new(big.Float).Quo(q.Dividend, i)
	}
	return models.NewRunOutputCompleteWithResult(i.String(), input.ResultCollection())
}
