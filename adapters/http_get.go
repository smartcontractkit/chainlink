package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// HttpGet holds the endpoint for the actual URL of the service or
// external adapter which will return the JSON object
type HttpGet struct {
	Endpoint models.WebURL `json:"endpoint"`
}

// Perform ensures that the http Endpoint responded without errors
// and returns the JSON result if successful
func (hga *HttpGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	response, err := http.Get(hga.Endpoint.String())
	if err != nil {
		return models.RunResultWithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return models.RunResultWithError(err)
	}

	if response.StatusCode >= 300 {
		return models.RunResultWithError(fmt.Errorf(body))
	}

	return models.RunResultWithValue(body)
}
