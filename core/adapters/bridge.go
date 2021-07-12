package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Bridge adapter is responsible for connecting the task pipeline to external
// adapters, allowing for custom computations to be executed and included in runs.
type Bridge struct {
	models.BridgeType
	Params models.JSON
}

// TaskType returns the bridges defined type.
func (ba *Bridge) TaskType() models.TaskType {
	return ba.Name
}

// Perform sends a POST request containing the JSON of the input to the
// external adapter specified in the BridgeType.
//
// It records the RunResult returned to it, and optionally marks the RunResult pending.
//
// If the Perform is resumed with a pending RunResult, the RunResult is marked
// not pending and the RunResult is returned.
func (ba *Bridge) Perform(input models.RunInput, store *store.Store, _ *keystore.Master) models.RunOutput {
	if input.Status().Completed() {
		return models.NewRunOutputComplete(input.Data())
	} else if input.Status().PendingBridge() {
		return models.NewRunOutputInProgress(input.Data())
	}
	return ba.handleNewRun(input, store)
}

func (ba *Bridge) handleNewRun(input models.RunInput, store *store.Store) models.RunOutput {
	data, err := models.MergeExceptResult(input.Data(), ba.Params)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("handling data param", err))
	}
	input = input.CloneWithData(data)

	responseURL := store.Config.BridgeResponseURL()
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", input.JobRunID().String())
	}

	httpConfig := defaultHTTPConfig(store.Config)
	// URL is "safe" because it comes from the node's own database
	// Some node operators may run external adapters on their own hardware
	httpConfig.AllowUnrestrictedNetworkAccess = true

	body, err := ba.postToExternalAdapter(input, responseURL, httpConfig)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("post to external adapter", err))
	}

	return ba.responseToRunResult(body, input)
}

func (ba *Bridge) responseToRunResult(body []byte, input models.RunInput) models.RunOutput {
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
		data, err := models.Merge(ba.Params, brr.Data)
		if err != nil {
			return models.NewRunOutputError(baRunResultError("handling data param", err))
		}

		return models.NewRunOutputComplete(data)
	}

	return models.NewRunOutputCompleteWithResult(brr.Data.String(), input.ResultCollection())
}

func (ba *Bridge) postToExternalAdapter(
	input models.RunInput,
	bridgeResponseURL *url.URL,
	config utils.HTTPRequestConfig,
) ([]byte, error) {
	outgoing := bridgeOutgoing{JobRunID: input.JobRunID().String(), Data: input.Data()}
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

	httpRequest := utils.HTTPRequest{
		Request: request,
		Config:  config,
	}

	bytes, statusCode, err := httpRequest.SendRequest(context.TODO())

	if err != nil {
		return nil, err
	}

	if statusCode >= 400 {
		err = fmt.Errorf("%v %v", statusCode, string(bytes))
		return nil, fmt.Errorf("POST request: %v", err)
	}

	return bytes, nil
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
