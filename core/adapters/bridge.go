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
func (ba *Bridge) Perform(input models.RunResult, store *store.Store) models.RunOutput {
	if input.Status.Finished() {
		return models.RunOutput{
			Data:         input.Data,
			Status:       input.Status,
			ErrorMessage: input.ErrorMessage,
		}
	} else if input.Status.PendingBridge() {
		return models.NewRunOutputInProgress()
	}
	return ba.handleNewRun(input, store.Config.BridgeResponseURL())
}

func (ba *Bridge) handleNewRun(input models.RunResult, bridgeResponseURL *url.URL) models.RunOutput {
	if ba.Params == nil {
		ba.Params = new(models.JSON)
	}
	var err error
	if input.Data, err = input.Data.Merge(*ba.Params); err != nil {
		return models.NewRunOutputError(baRunResultError("handling data param", err))
	}

	responseURL := bridgeResponseURL
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", input.CachedJobRunID)
	}

	body, err := ba.postToExternalAdapter(input, responseURL)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("post to external adapter", err))
	}

	return responseToRunResult(body, input)
}

func responseToRunResult(body []byte, input models.RunResult) models.RunOutput {
	var brr models.BridgeRunResult
	err := json.Unmarshal(body, &brr)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("unmarshaling JSON", err))
	}

	status := brr.Status
	// FIXME: This code is bizarre to me, what is the intention? It looks like it
	// has a bunch of edge cases...
	var data models.JSON
	if brr.Data.Exists() && !brr.Data.IsObject() {
		data, _ = data.Add("result", brr.Data.String())
		status = models.RunStatusCompleted
	}
	data, _ = data.Merge(brr.Data)

	return models.RunOutput{
		Status:       status,
		ErrorMessage: brr.ErrorMessage,
		Data:         data,
	}
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
