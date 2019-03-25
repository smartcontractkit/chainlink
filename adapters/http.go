package adapters

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// HTTPGet requires a URL which is used for a GET request when the adapter is called.
type HTTPGet struct {
	URL models.WebURL `json:"url"`
	GET models.WebURL `json:"get"`
	Headers http.Header `json:"headers"`
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", hga.GetURL(), nil)
	if err != nil {
		input.SetError(err)
		return input
	}
	if hga.Headers != nil {
		request.Header = hga.Headers
	}
	response, err := client.Do(request)
	if err != nil {
		input.SetError(err)
		return input
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		input.SetError(err)
		return input
	}

	if response.StatusCode >= 400 {
		input.SetError(errors.New(body))
		return input
	}

	input.ApplyResult(body)
	return input
}

// GetURL retrieves the GET field if set otherwise returns the URL field
func (hga *HTTPGet) GetURL() string {
	if hga.GET.String() != "" {
		return hga.GET.String()
	}
	return hga.URL.String()
}

// HTTPPost requires a URL which is used for a POST request when the adapter is called.
type HTTPPost struct {
	URL  models.WebURL `json:"url"`
	POST models.WebURL `json:"post"`
	Headers http.Header `json:"headers"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	reqBody := bytes.NewBufferString(input.Data.String())
	request, err := http.NewRequest("POST", hpa.GetURL(), reqBody)
	if err != nil {
		input.SetError(err)
		return input
	}
	if hpa.Headers != nil {
		request.Header = hpa.Headers
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		input.SetError(err)
		return input
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		input.SetError(err)
		return input
	}

	if response.StatusCode >= 400 {
		input.SetError(errors.New(body))
		return input
	}

	input.ApplyResult(body)
	return input
}

// GetURL retrieves the POST field if set otherwise returns the URL field
func (hpa *HTTPPost) GetURL() string {
	if hpa.POST.String() != "" {
		return hpa.POST.String()
	}
	return hpa.URL.String()
}
