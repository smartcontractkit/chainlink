package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Bridge adapter is responsible for connecting the task pipeline to external
// adapters, allowing for custom computations to be executed and included in runs.
type Bridge struct {
	models.BridgeType
}

// Perform sends a POST request containing the JSON of the input RunResult to
// the external adapter specified in the BridgeType.
// It records the RunResult returned to it, and optionally marks the RunResult pending.
//
// If the Perform is resumed with a pending RunResult, the RunResult is marked
// not pending and the RunResult is returned.
func (ba *Bridge) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	if input.Pending {
		return markNotPending(input)
	}
	return ba.handleNewRun(input)
}

func markNotPending(input models.RunResult) models.RunResult {
	input.Pending = false
	return input
}

func (ba *Bridge) handleNewRun(input models.RunResult) models.RunResult {
	in, err := json.Marshal(&bridgePayload{input})
	if err != nil {
		return baRunResultError(input, "marshaling request body", err)
	}

	resp, err := http.Post(ba.URL.String(), "application/json", bytes.NewBuffer(in))
	if err != nil {
		return baRunResultError(input, "POST request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("%v %v", resp.StatusCode, string(b))
		return baRunResultError(input, "POST response", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return baRunResultError(input, "reading response body", err)
	}

	rr := models.RunResult{}
	err = json.Unmarshal(b, &rr)
	if err != nil {
		return baRunResultError(input, "unmarshaling JSON", err)
	}
	return rr
}

func baRunResultError(in models.RunResult, str string, err error) models.RunResult {
	return in.WithError(fmt.Errorf("ExternalBridge %v: %v", str, err))
}

type bridgePayload struct {
	models.RunResult
}

func (bp bridgePayload) MarshalJSON() ([]byte, error) {
	anon := struct {
		JobRunID string      `json:"id"`
		Data     models.JSON `json:"data"`
	}{
		JobRunID: bp.JobRunID,
		Data:     bp.Data,
	}
	return json.Marshal(anon)
}
