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
	if input.Status.Finished() {
		return input
	} else if input.Status.PendingBridge() {
		return resumeBridge(input)
	}
	return ba.handleNewRun(input)
}

// MinConfs specifies the number of block confirmations
// needed for the Bridge to run. This method enables the Bridge to meet the
// adapters.AdapterWithMinConfs interface.
func (ba *Bridge) MinConfs() uint64 {
	return ba.DefaultConfirmations
}

func resumeBridge(input models.RunResult) models.RunResult {
	input.Status = models.RunStatusInProgress
	return input
}

func (ba *Bridge) handleNewRun(input models.RunResult) models.RunResult {
	b, err := postToExternalAdapter(ba.URL.String(), input)
	if err != nil {
		return baRunResultError(input, "post to external adapter", err)
	}

	var brr models.BridgeRunResult
	err = json.Unmarshal(b, &brr)
	if err != nil {
		return baRunResultError(input, "unmarshaling JSON", err)
	}

	rr, err := input.Merge(brr.RunResult)
	if err != nil {
		return baRunResultError(rr, "Unable to merge received payload", err)
	}

	return rr
}

func postToExternalAdapter(url string, input models.RunResult) ([]byte, error) {
	in, err := json.Marshal(&bridgeOutgoing{input})
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(in))
	if err != nil {
		return nil, fmt.Errorf("POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("%v %v", resp.StatusCode, string(b))
		return nil, fmt.Errorf("POST response: %v", err)
	}

	return ioutil.ReadAll(resp.Body)
}

func baRunResultError(in models.RunResult, str string, err error) models.RunResult {
	return in.WithError(fmt.Errorf("ExternalBridge %v: %v", str, err))
}

type bridgeOutgoing struct {
	models.RunResult
}

func (bp bridgeOutgoing) MarshalJSON() ([]byte, error) {
	anon := struct {
		JobRunID string      `json:"id"`
		Data     models.JSON `json:"data"`
	}{
		JobRunID: bp.JobRunID,
		Data:     bp.Data,
	}
	return json.Marshal(anon)
}
