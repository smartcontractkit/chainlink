package models

import (
	"encoding/json"
	"errors"

	null "gopkg.in/guregu/null.v4"
)

// BridgeRunResult handles the parsing of RunResults from external adapters.
type BridgeRunResult struct {
	Data            JSON        `json:"data"`
	Status          RunStatus   `json:"status"`
	ErrorMessage    null.String `json:"error"`
	ExternalPending bool        `json:"pending"`
	AccessToken     string      `json:"accessToken"`
}

// UnmarshalJSON parses the given input and updates the BridgeRunResult in the
// external adapter format.
func (brr *BridgeRunResult) UnmarshalJSON(input []byte) error {
	// XXX: This indirection prevents an infinite regress during json.Unmarshal
	type biAlias BridgeRunResult
	var anon biAlias
	err := json.Unmarshal(input, &anon)
	*brr = BridgeRunResult(anon)

	if brr.Status == RunStatusErrored || brr.ErrorMessage.Valid {
		brr.Status = RunStatusErrored
	} else if brr.ExternalPending || brr.Status.PendingBridge() {
		brr.Status = RunStatusPendingBridge
	} else {
		brr.Status = RunStatusCompleted
	}

	return err
}

// HasError returns true if the status is errored or the error message is set
func (brr BridgeRunResult) HasError() bool {
	return brr.Status == RunStatusErrored || brr.ErrorMessage.Valid
}

// GetError returns the error of a BridgeRunResult if it is present.
func (brr BridgeRunResult) GetError() error {
	if brr.HasError() {
		return errors.New(brr.ErrorMessage.ValueOrZero())
	}
	return nil
}
