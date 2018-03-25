package adapters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// HTTPGet requires a URL which is used for a GET request when the adapter is called.
type HTTPGet struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
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

// HTTPPost requires a URL which is used for a POST request when the adapter is called.
type HTTPPost struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	reqBody := bytes.NewBufferString(input.Data.String())
	response, err := http.Post(hpa.URL.String(), "application/json", reqBody)
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
