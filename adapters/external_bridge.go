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

type ExternalBridge struct {
	*models.CustomTaskType
}

func (eb *ExternalBridge) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	resp, err := http.Post(eb.URL.String(), "application/json", &bytes.Buffer{})
	if err != nil {
		return ebRunResultError("POST request", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ebRunResultError("reading response body", err)
	}

	output := models.Output{}
	err = json.Unmarshal(b, &output)
	if err != nil {
		return ebRunResultError("unmarshaling JSON", err)
	}
	return models.RunResult{Output: output}
}

func ebRunResultError(str string, err error) models.RunResult {
	return models.RunResultWithError(fmt.Errorf("ExternalBridge %v: %v", str, err))
}
