package adapters

import (
	"encoding/json"
	"math/big"

	"chainlink/core/utils"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times *big.Float `json:"-"`
}

type jsonMultiply struct {
	Times *utils.BigFloat `json:"times,omitempty"`
}

// MarshalJSON implements the json.Marshal interface.
func (ma Multiply) MarshalJSON() ([]byte, error) {
	jsonObj := jsonMultiply{Times: (*utils.BigFloat)(ma.Times)}
	return json.Marshal(jsonObj)
}

// UnmarshalJSON implements the json.Unmarshal interface.
func (ma *Multiply) UnmarshalJSON(buf []byte) error {
	var jsonObj jsonMultiply
	err := json.Unmarshal(buf, &jsonObj)
	if err != nil {
		return err
	}
	ma.Times = jsonObj.Times.Value()
	return nil
}
