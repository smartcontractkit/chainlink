package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
func (ba *Bridge) Perform(input models.JSON, result models.RunResult, store *store.Store) models.RunResult {
	if result.Status.Finished() {
		return result
	} else if result.Status.PendingBridge() {
		return resumeBridge(result)
	}
	ba.handleNewRun(input, &result, store.Config.BridgeResponseURL())
	return result
}

func resumeBridge(result models.RunResult) models.RunResult {
	result.Status = models.RunStatusInProgress
	return result
}

func (ba *Bridge) handleNewRun(input models.JSON, result *models.RunResult, bridgeResponseURL *url.URL) {
	if ba.Params == nil {
		ba.Params = new(models.JSON)
	}
	var err error
	if input, err = input.Merge(*ba.Params); err != nil {
		result.SetError(baRunResultError("handling data param", err))
		return
	}

	responseURL := bridgeResponseURL
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", result.CachedJobRunID)
	}
	body, err := ba.postToExternalAdapter(input, result, responseURL)
	if err != nil {
		result.SetError(baRunResultError("post to external adapter", err))
		return
	}

	err = responseToRunResult(body, result)
	if err != nil {
		result.SetError(err)
		return
	}
}

func responseToRunResult(body []byte, result *models.RunResult) error {
	var brr models.BridgeRunResult
	err := json.Unmarshal(body, &brr)
	if err != nil {
		return baRunResultError("unmarshaling JSON", err)
	}

	if brr.RunResult.Data.Exists() && !brr.RunResult.Data.IsObject() {
		result.CompleteWithResult(brr.RunResult.Data.String())
	}

	return result.Merge(brr.RunResult)
}

func (ba *Bridge) postToExternalAdapter(input models.JSON, result *models.RunResult, bridgeResponseURL *url.URL) ([]byte, error) {
	in, err := json.Marshal(&bridgeOutgoing{
		Input:       input,
		RunResult:   *result,
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
	Input       models.JSON
	ResponseURL *url.URL
}

func (bp bridgeOutgoing) MarshalJSON() ([]byte, error) {
	anon := struct {
		JobRunID    string      `json:"id"`
		Data        models.JSON `json:"data"`
		ResponseURL string      `json:"responseURL,omitempty"`
	}{
		JobRunID:    bp.RunResult.CachedJobRunID,
		Data:        bp.Input,
		ResponseURL: bp.ResponseURL.String(),
	}
	return json.Marshal(anon)
}

var zeroURL = new(url.URL)
