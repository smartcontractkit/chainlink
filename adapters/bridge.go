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

type Bridge struct {
	*models.BridgeType
}

func (ba *Bridge) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	resp, err := http.Post(ba.URL.String(), "application/json", &bytes.Buffer{})
	if err != nil {
		return baRunResultError("POST request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("%v %v", resp.StatusCode, string(b))
		return baRunResultError("POST reponse", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return baRunResultError("reading response body", err)
	}

	output := models.Output{}
	err = json.Unmarshal(b, &output)
	if err != nil {
		return baRunResultError("unmarshaling JSON", err)
	}
	return models.RunResult{Output: output}
}

func baRunResultError(str string, err error) models.RunResult {
	return models.RunResultWithError(fmt.Errorf("ExternalBridge %v: %v", str, err))
}
