package adapters

import (
	"bytes"
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

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HttpGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	response, err := http.Get(hga.URL.String())
	if err != nil {
		return input.WithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(body))
	}

	return input.WithValue(body)
}

// HttpPost requires a URL which is used for a POST request when the adapter is called.
type HttpPost struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hga *HttpPost) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	reqBody := bytes.NewBufferString(input.Data.String())
	response, err := http.Post(hga.URL.String(), "application/json", reqBody)
	if err != nil {
		return input.WithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(body))
	}

	return input.WithValue(body)
}
