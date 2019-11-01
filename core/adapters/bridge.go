package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"chainlink/core/store"
	"chainlink/core/store/models"
)

// Bridge adapter is responsible for connecting the task pipeline to external
// adapters, allowing for custom computations to be executed and included in runs.
type Bridge struct {
	*models.BridgeType
	Params *models.JSON
}

// Perform sends a POST request containing the JSON of the input RunResult to
// the external adapter specified in the BridgeType.
// It records the RunResult returned to it, and optionally marks the RunResult pending.
//
// If the Perform is resumed with a pending RunResult, the RunResult is marked
// not pending and the RunResult is returned.
func (ba *Bridge) Perform(input models.RunResult, store *store.Store) models.RunResult {
	if input.Status.Finished() {
		return input
	} else if input.Status.PendingBridge() {
		return resumeBridge(input)
	}
	return ba.handleNewRun(input, store.Config.BridgeResponseURL())
}

func resumeBridge(input models.RunResult) models.RunResult {
	input.Status = models.RunStatusInProgress
	return input
}

func (ba *Bridge) handleNewRun(input models.RunResult, bridgeResponseURL *url.URL) models.RunResult {
	if ba.Params == nil {
		ba.Params = new(models.JSON)
	}
	var err error
	if input.Data, err = input.Data.Merge(*ba.Params); err != nil {
		return models.RunResultError(baRunResultError("handling data param", err))
	}

	responseURL := bridgeResponseURL
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", input.CachedJobRunID)
	}

	body, err := ba.postToExternalAdapter(input, responseURL)
	if err != nil {
		return models.RunResultError(baRunResultError("post to external adapter", err))
	}

	return responseToRunResult(body, input)
}

func responseToRunResult(body []byte, input models.RunResult) models.RunResult {
	var output models.RunResult

	var brr models.BridgeRunResult
	err := json.Unmarshal(body, &brr)
	if err != nil {
		return models.RunResultError(baRunResultError("unmarshaling JSON", err))
	} else if brr.HasError() {
		return brr.RunResult
	}

	if brr.RunResult.Data.Exists() && !brr.RunResult.Data.IsObject() {
		output.CompleteWithResult(brr.RunResult.Data.String())
	}

	err = output.Merge(brr.RunResult)
	if err != nil {
		output.SetError(err)
	}

	return output
}

func (ba *Bridge) postToExternalAdapter(input models.RunResult, bridgeResponseURL *url.URL) ([]byte, error) {
	in, err := json.Marshal(&bridgeOutgoing{
		RunResult:   input,
		ResponseURL: bridgeResponseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %v", err)
	}

	request, err := http.NewRequest("POST", ba.URL.String(), bytes.NewBuffer(in))
	if err != nil {
		return nil, fmt.Errorf("building outgoing bridge http post: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+ba.BridgeType.OutgoingToken)
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
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

func baRunResultError(str string, err error) error {
	return fmt.Errorf("ExternalBridge %v: %v", str, err)
}

type bridgeOutgoing struct {
	models.RunResult
	ResponseURL *url.URL
}

func (bp bridgeOutgoing) MarshalJSON() ([]byte, error) {
	anon := struct {
		JobRunID    *models.ID  `json:"id"`
		Data        models.JSON `json:"data"`
		ResponseURL string      `json:"responseURL,omitempty"`
	}{
		JobRunID:    bp.CachedJobRunID,
		Data:        bp.Data,
		ResponseURL: bp.ResponseURL.String(),
	}
	return json.Marshal(anon)
}

var zeroURL = new(url.URL)
