package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// HTTPGet requires a URL which is used for a GET request when the adapter is called.
type HTTPGet struct {
	URL                            models.WebURL   `json:"url"`
	GET                            models.WebURL   `json:"get"`
	Headers                        http.Header     `json:"headers"`
	QueryParams                    QueryParameters `json:"queryParams"`
	ExtendedPath                   ExtendedPath    `json:"extPath"`
	AllowUnrestrictedNetworkAccess bool            `json:"-"`
}

// TaskType returns the type of Adapter.
func (hga *HTTPGet) TaskType() models.TaskType {
	return TaskTypeHTTPGet
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunInput, store *store.Store, _ *keystore.Master) models.RunOutput {
	request, err := hga.GetRequest()
	if err != nil {
		return models.NewRunOutputError(err)
	}
	httpConfig := defaultHTTPConfig(store.Config)
	httpConfig.AllowUnrestrictedNetworkAccess = hga.AllowUnrestrictedNetworkAccess
	return sendRequest(input, request, httpConfig)
}

// GetURL retrieves the GET field if set otherwise returns the URL field
func (hga *HTTPGet) GetURL() string {
	if hga.GET.String() != "" {
		return hga.GET.String()
	}
	return hga.URL.String()
}

// GetRequest returns the HTTP request including query parameters and headers
func (hga *HTTPGet) GetRequest() (*http.Request, error) {
	request, err := http.NewRequest("GET", hga.GetURL(), nil)
	if err != nil {
		return nil, err
	}
	appendExtendedPath(request, hga.ExtendedPath)
	appendQueryParams(request, hga.QueryParams)
	setHeaders(request, hga.Headers, "")
	return request, nil
}

// HTTPPost requires a URL which is used for a POST request when the adapter is called.
type HTTPPost struct {
	URL                            models.WebURL   `json:"url"`
	POST                           models.WebURL   `json:"post"`
	Headers                        http.Header     `json:"headers"`
	QueryParams                    QueryParameters `json:"queryParams"`
	Body                           *string         `json:"body,omitempty"`
	ExtendedPath                   ExtendedPath    `json:"extPath"`
	AllowUnrestrictedNetworkAccess bool            `json:"-"`
}

// TaskType returns the type of Adapter.
func (hpa *HTTPPost) TaskType() models.TaskType {
	return TaskTypeHTTPPost
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunInput, store *store.Store, _ *keystore.Master) models.RunOutput {
	request, err := hpa.GetRequest(input.Data().String())
	if err != nil {
		return models.NewRunOutputError(err)
	}
	httpConfig := defaultHTTPConfig(store.Config)
	httpConfig.AllowUnrestrictedNetworkAccess = hpa.AllowUnrestrictedNetworkAccess
	return sendRequest(input, request, httpConfig)
}

// GetURL retrieves the POST field if set otherwise returns the URL field
func (hpa *HTTPPost) GetURL() string {
	if hpa.POST.String() != "" {
		return hpa.POST.String()
	}
	return hpa.URL.String()
}

// GetRequest takes the request body and returns the HTTP request including
// query parameters and headers.
//
// HTTPPost's Body parameter overrides the given argument if present.
func (hpa *HTTPPost) GetRequest(body string) (*http.Request, error) {
	if hpa.Body != nil {
		body = *hpa.Body
	}
	reqBody := bytes.NewBufferString(body)

	request, err := http.NewRequest("POST", hpa.GetURL(), reqBody)
	if err != nil {
		return nil, err
	}
	appendExtendedPath(request, hpa.ExtendedPath)
	appendQueryParams(request, hpa.QueryParams)
	setHeaders(request, hpa.Headers, "application/json")
	return request, nil
}

func appendExtendedPath(request *http.Request, extPath ExtendedPath) {
	// Remove early empty extPath entries
	extPaths := []string(extPath[:])
	for _, path := range extPaths {
		if len(path) != 0 {
			break
		}
		extPaths = extPaths[1:]
	}

	if len(extPaths) == 0 {
		return
	}

	if strings.HasPrefix(extPath[0], "/") || strings.HasSuffix(request.URL.Path, "/") {
		request.URL.Path = request.URL.Path + path.Join(extPaths...)
		return
	}

	request.URL.Path = request.URL.Path + "/" + path.Join(extPaths...)
}

func appendQueryParams(request *http.Request, queryParams QueryParameters) {
	q := request.URL.Query()
	for k, v := range queryParams {
		if !keyExists(k, q) {
			q.Add(k, v[0])
		}
	}
	request.URL.RawQuery = q.Encode()
}

func keyExists(key string, query url.Values) bool {
	_, ok := query[key]
	return ok
}

func setHeaders(request *http.Request, headers http.Header, contentType string) {
	if headers != nil {
		request.Header = headers
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
}

func sendRequest(input models.RunInput, request *http.Request, config utils.HTTPRequestConfig) models.RunOutput {
	httpRequest := utils.HTTPRequest{
		Request: request,
		Config:  config,
	}

	bytes, statusCode, _, err := httpRequest.SendRequest(context.TODO())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	responseBody := string(bytes)

	// This is either a client error caused on our end or a server error that persists even after retrying.
	// Either way, there is no way for us to complete the run with a result.
	if statusCode >= 400 {
		return models.NewRunOutputError(errors.New(responseBody))
	}

	return models.NewRunOutputCompleteWithResult(responseBody, input.ResultCollection())
}

// QueryParameters are the keys and values to append to the URL
type QueryParameters url.Values

// UnmarshalJSON implements the Unmarshaler interface
func (qp *QueryParameters) UnmarshalJSON(input []byte) error {
	var strs []string
	var err error

	// input is a string like "someKey0=someVal0&someKey1=someVal1"
	if utils.IsQuoted(input) {
		var decoded string
		unmErr := json.Unmarshal(input, &decoded)
		if unmErr != nil {
			return fmt.Errorf("unable to unmarshal query parameters: %s", input)
		}
		strs = strings.FieldsFunc(trimQuestion(decoded), splitQueryString)

		// input is an array of strings like
		// ["someKey0", "someVal0", "someKey1", "someVal1"]
	} else if err = json.Unmarshal(input, &strs); err != nil {
		return fmt.Errorf("unable to unmarshal query parameters: %s", input)
	}

	values, err := buildValues(strs)
	if err != nil {
		return fmt.Errorf("unable to build query parameters: %s", input)
	}
	*qp = QueryParameters(values)
	return err
}

func splitQueryString(r rune) bool {
	return r == '=' || r == '&'
}

func trimQuestion(input string) string {
	return strings.Replace(input, "?", "", -1)
}

func buildValues(input []string) (url.Values, error) {
	values := url.Values{}
	if len(input)%2 != 0 {
		return nil, fmt.Errorf("invalid number of parameters: %s", input)
	}
	for i := 0; i < len(input); i = i + 2 {
		values.Add(input[i], input[i+1])
	}
	return values, nil
}

// ExtendedPath is the path to append to a base URL
type ExtendedPath []string

// UnmarshalJSON implements the Unmarshaler interface
func (ep *ExtendedPath) UnmarshalJSON(input []byte) error {
	values := []string{}
	var err error
	// if input is a string like "a/b/c"
	if utils.IsQuoted(input) {
		values = strings.Split(string(utils.RemoveQuotes(input)), "/")
		// if input is an array of strings like ["a", "b", "c"]
	} else {
		err = json.Unmarshal(input, &values)
	}
	*ep = ExtendedPath(values)
	return err
}

func defaultHTTPConfig(config orm.ConfigReader) utils.HTTPRequestConfig {
	return utils.HTTPRequestConfig{
		Timeout:                        config.DefaultHTTPTimeout().Duration(),
		MaxAttempts:                    config.DefaultMaxHTTPAttempts(),
		SizeLimit:                      config.DefaultHTTPLimit(),
		AllowUnrestrictedNetworkAccess: config.DefaultHTTPAllowUnrestrictedNetworkAccess(),
	}
}
