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
func (ba *Bridge) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	in, err := json.Marshal(&input.Data)
	if err != nil {
		return baRunResultError("marshaling request body", err)
	}

	resp, err := http.Post(ba.URL.String(), "application/json", bytes.NewBuffer(in))
	if err != nil {
		return baRunResultError("POST request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("%v %v", resp.StatusCode, string(b))
		return baRunResultError("POST reponse", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return baRunResultError("reading response body", err)
	}

	rr := models.RunResult{}
	err = json.Unmarshal(b, &rr)
	if err != nil {
		return baRunResultError("unmarshaling JSON", err)
	}
	return rr
}

func baRunResultError(str string, err error) models.RunResult {
	return models.RunResultWithError(fmt.Errorf("ExternalBridge %v: %v", str, err))
}
