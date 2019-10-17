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
	Params models.JSON
}

// Perform sends a POST request containing the JSON of the input to the
// external adapter specified in the BridgeType.
//
// It records the RunResult returned to it, and optionally marks the RunResult pending.
//
// If the Perform is resumed with a pending RunResult, the RunResult is marked
// not pending and the RunResult is returned.
func (ba *Bridge) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	if input.Status.Finished() {
		return models.RunOutput{
			Data:         input.Data,
			Status:       input.Status,
			ErrorMessage: input.ErrorMessage,
		}
	} else if input.Status.PendingBridge() {
		return models.NewRunOutputInProgress(input.Data)
	}
	return ba.handleNewRun(input, store.Config.BridgeResponseURL())
}

func (ba *Bridge) handleNewRun(input models.RunInput, bridgeResponseURL *url.URL) models.RunOutput {
	var err error
	if input.Data, err = input.Data.Merge(ba.Params); err != nil {
		return models.NewRunOutputError(baRunResultError("handling data param", err))
	}

	responseURL := bridgeResponseURL
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", input.JobRunID.String())
	}

	body, err := ba.postToExternalAdapter(input, responseURL)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("post to external adapter", err))
	}

	return responseToRunResult(body, input)
}

func responseToRunResult(body []byte, input models.RunInput) models.RunOutput {
	var brr models.BridgeRunResult
	err := json.Unmarshal(body, &brr)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("unmarshaling JSON", err))
	}

	if brr.HasError() {
		return models.NewRunOutputError(brr.GetError())
	}

	if brr.ExternalPending {
		return models.NewRunOutputPendingBridge()
	}

	if brr.Data.IsObject() {
		return models.NewRunOutputComplete(brr.Data)
	}

	return models.NewRunOutputCompleteWithResult(brr.Data.String())
}

func (ba *Bridge) postToExternalAdapter(input models.RunInput, bridgeResponseURL *url.URL) ([]byte, error) {
	outgoing := bridgeOutgoing{JobRunID: input.JobRunID.String(), Data: input.Data}
	if bridgeResponseURL != nil {
		outgoing.ResponseURL = bridgeResponseURL.String()
	}
	in, err := json.Marshal(&outgoing)
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
	JobRunID    string      `json:"id"`
	Data        models.JSON `json:"data"`
	ResponseURL string      `json:"responseURL,omitempty"`
}

var zeroURL = new(url.URL)
