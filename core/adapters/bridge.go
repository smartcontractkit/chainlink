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

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
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
func (ba *Bridge) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	if input.Status().Completed() {
		return models.NewRunOutputComplete(input.Data())
	} else if input.Status().PendingBridge() {
		return models.NewRunOutputInProgress(input.Data())
	}
	meta := getMeta(store, input.JobRunID())
	return ba.handleNewRun(input, meta, store.Config.BridgeResponseURL())
}

func getMeta(store *store.Store, jobRunID *models.ID) *models.JSON {
	jobRun, err := store.ORM.FindJobRun(jobRunID)
	if err != nil {
		return nil
	} else if jobRun.RunRequest.TxHash == nil || jobRun.RunRequest.BlockHash == nil {
		return nil
	}
	meta := fmt.Sprintf(`
		{
			"initiator": {
				"transactionHash": "%s",
				"blockHash": "%s"
			}
		}`,
		jobRun.RunRequest.TxHash.Hex(),
		jobRun.RunRequest.BlockHash.Hex(),
	)
	return &models.JSON{Result: gjson.Parse(meta)}
}

func (ba *Bridge) handleNewRun(input models.RunInput, meta *models.JSON, bridgeResponseURL *url.URL) models.RunOutput {
	data, err := models.Merge(input.Data(), ba.Params)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("handling data param", err))
	}

	responseURL := bridgeResponseURL
	if *responseURL != *zeroURL {
		responseURL.Path += fmt.Sprintf("/v2/runs/%s", input.JobRunID().String())
	}

	body, err := ba.postToExternalAdapter(input, meta, responseURL)
	if err != nil {
		return models.NewRunOutputError(baRunResultError("post to external adapter", err))
	}

	input = input.CloneWithData(data)
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

	return models.NewRunOutputCompleteWithResult(brr.Data.String())
}

func (ba *Bridge) postToExternalAdapter(input models.RunInput, meta *models.JSON, bridgeResponseURL *url.URL) ([]byte, error) {
	data, err := models.Merge(input.Data(), ba.Params)
	if err != nil {
		return nil, errors.Wrap(err, "error merging bridge params with input params")
	}

	outgoing := bridgeOutgoing{JobRunID: input.JobRunID().String(), Data: data, Meta: meta}
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
	JobRunID    string       `json:"id"`
	Data        models.JSON  `json:"data"`
	Meta        *models.JSON `json:"meta,omitempty"`
	ResponseURL string       `json:"responseURL,omitempty"`
}

var zeroURL = new(url.URL)
