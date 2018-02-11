package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// HttpGet requires a URL which is used for a GET request when the adapter is called.
type HttpGet struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the http URL responded without errors
// and returns the response body as the "value" field of the result.
func (hga *HttpGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	response, err := http.Get(hga.URL.String())
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
