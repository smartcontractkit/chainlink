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
	// Attempt to retrieve the data from the Endpoint
	response, err := http.Get(hga.Endpoint.String())
	// Return the error if present
	if err != nil {
		return models.RunResultWithError(err)
	}
	// Do not close the Body until the function returns
	defer response.Body.Close()
	// Store the raw JSON data
	bytes, err := ioutil.ReadAll(response.Body)
	// Convert to string
	body := string(bytes)
	// Return error if JSON data could not be read
	if err != nil {
		return models.RunResultWithError(err)
	}
	// Return error if there were any errors with the http request
	if response.StatusCode >= 300 {
		return models.RunResultWithError(fmt.Errorf(body))
	}
	// Return the response from the Endpoint
	return models.RunResultWithValue(body)
}
