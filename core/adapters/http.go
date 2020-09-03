package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/jpillora/backoff"
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

// HTTPRequestConfig holds the configurable settings for an http request
type HTTPRequestConfig struct {
	timeout                        time.Duration
	maxAttempts                    uint
	sizeLimit                      int64
	allowUnrestrictedNetworkAccess bool
}

// TaskType returns the type of Adapter.
func (hga *HTTPGet) TaskType() models.TaskType {
	return TaskTypeHTTPGet
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	request, err := hga.GetRequest(input)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	httpConfig := defaultHTTPConfig(store)
	httpConfig.allowUnrestrictedNetworkAccess = hga.AllowUnrestrictedNetworkAccess
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
func (hga *HTTPGet) GetRequest(input models.RunInput) (*http.Request, error) {
	request, err := http.NewRequest("GET", hga.GetURL(), nil)
	if err != nil {
		return nil, err
	}
	if strpkg.GetServiceMemory()[input.JobRunID().String()] != nil {
		extendedPath := strings.Split(input.Result().String(), "/")
		appendExtendedPath(request, extendedPath)
	} else {
		appendExtendedPath(request, hga.ExtendedPath)
	}
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
func (hpa *HTTPPost) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	request, err := hpa.GetRequest(input.Data().String())
	if err != nil {
		return models.NewRunOutputError(err)
	}
	httpConfig := defaultHTTPConfig(store)
	httpConfig.allowUnrestrictedNetworkAccess = hpa.AllowUnrestrictedNetworkAccess
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
	request.URL.Path = path.Join(append([]string{request.URL.Path}, []string(extPath)...)...)
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

func sendRequest(input models.RunInput, request *http.Request, config HTTPRequestConfig) models.RunOutput {
	tr := &http.Transport{
		DisableCompression: true,
	}
	if !config.allowUnrestrictedNetworkAccess {
		tr.DialContext = restrictedDialContext
	}
	client := &http.Client{Transport: tr}

	bytes, statusCode, err := withRetry(client, request, config)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	responseBody := string(bytes)

	// This is either a client error caused on our end or a server error that persists even after retrying.
	// Either way, there is no way for us to complete the run with a result.
	if statusCode >= 400 {
		return models.NewRunOutputError(errors.New(responseBody))
	}

	return models.NewRunOutputCompleteWithResult(responseBody)
}

// withRetry executes the http request in a retry. Timeout is controlled with a context
// Retry occurs if the request timeout, or there is any kind of connection or transport-layer error
// Retry also occurs on remote server 5xx errors
func withRetry(
	client *http.Client,
	originalRequest *http.Request,
	config HTTPRequestConfig,
) (responseBody []byte, statusCode int, err error) {
	bb := &backoff.Backoff{
		Min:    100,
		Max:    20 * time.Minute, // We stop retrying on the number of attempts!
		Jitter: true,
	}
	for {
		responseBody, statusCode, err = makeHTTPCall(client, originalRequest, config)
		if err == nil {
			return responseBody, statusCode, nil
		}
		if uint(bb.Attempt())+1 >= config.maxAttempts { // Stop retrying.
			return responseBody, statusCode, err
		}
		switch err.(type) {
		// There is no point in retrying a request if the response was
		// too large since it's likely that all retries will suffer the
		// same problem
		case *HTTPResponseTooLargeError:
			return responseBody, statusCode, err
		}
		// Sleep and retry.
		time.Sleep(bb.Duration())
		logger.Debugw("http adapter error, will retry", "error", err.Error(), "attempt", bb.Attempt(), "timeout", config.timeout)
	}
}

func makeHTTPCall(
	client *http.Client,
	originalRequest *http.Request,
	config HTTPRequestConfig,
) (responseBody []byte, statusCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	requestWithTimeout := originalRequest.Clone(ctx)

	start := time.Now()

	r, e := client.Do(requestWithTimeout)
	if e != nil {
		return nil, 0, e
	}
	defer logger.ErrorIfCalling(r.Body.Close)

	statusCode = r.StatusCode
	elapsed := time.Since(start)
	logger.Debugw(fmt.Sprintf("http adapter got %v in %s", statusCode, elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	source := newMaxBytesReader(r.Body, config.sizeLimit)
	bytes, e := ioutil.ReadAll(source)
	if e != nil {
		logger.Errorf("http adapter error reading body: %v", e.Error())
		return nil, statusCode, e
	}
	elapsed = time.Since(start)
	logger.Debugw(fmt.Sprintf("http adapter finished after %s", elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	responseBody = bytes

	// Retry on 5xx since this might give a different result
	if 500 <= r.StatusCode && r.StatusCode < 600 {
		return responseBody, statusCode, &RemoteServerError{responseBody, statusCode}
	}

	return responseBody, statusCode, nil
}

type RemoteServerError struct {
	responseBody []byte
	statusCode   int
}

func (e *RemoteServerError) Error() string {
	return fmt.Sprintf("remote server error: %v\nResponse body: %v", e.statusCode, string(e.responseBody))
}

// maxBytesReader is inspired by
// https://github.com/gin-contrib/size/blob/master/size.go
type maxBytesReader struct {
	rc               io.ReadCloser
	limit, remaining int64
	sawEOF           bool
}

func newMaxBytesReader(rc io.ReadCloser, limit int64) *maxBytesReader {
	return &maxBytesReader{
		rc:        rc,
		limit:     limit,
		remaining: limit,
	}
}

func (mbr *maxBytesReader) Read(p []byte) (n int, err error) {
	toRead := mbr.remaining
	if mbr.remaining == 0 {
		if mbr.sawEOF {
			return mbr.tooLarge()
		}
		// The underlying io.Reader may not return (0, io.EOF)
		// at EOF if the requested size is 0, so read 1 byte
		// instead. The io.Reader docs are a bit ambiguous
		// about the return value of Read when 0 bytes are
		// requested, and {bytes,strings}.Reader gets it wrong
		// too (it returns (0, nil) even at EOF).
		toRead = 1
	}
	if int64(len(p)) > toRead {
		p = p[:toRead]
	}
	n, err = mbr.rc.Read(p)
	if err == io.EOF {
		mbr.sawEOF = true
	}
	if mbr.remaining == 0 {
		// If we had zero bytes to read remaining (but hadn't seen EOF)
		// and we get a byte here, that means we went over our limit.
		if n > 0 {
			return mbr.tooLarge()
		}
		return 0, err
	}
	mbr.remaining -= int64(n)
	if mbr.remaining < 0 {
		mbr.remaining = 0
	}
	return
}

type HTTPResponseTooLargeError struct {
	limit int64
}

func (e *HTTPResponseTooLargeError) Error() string {
	return fmt.Sprintf("HTTP response too large, must be less than %d bytes", e.limit)
}

func (mbr *maxBytesReader) tooLarge() (int, error) {
	return 0, &HTTPResponseTooLargeError{mbr.limit}
}

func (mbr *maxBytesReader) Close() error {
	return mbr.rc.Close()
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

func defaultHTTPConfig(store *store.Store) HTTPRequestConfig {
	return HTTPRequestConfig{
		store.Config.DefaultHTTPTimeout().Duration(),
		store.Config.DefaultMaxHTTPAttempts(),
		store.Config.DefaultHTTPLimit(),
		false,
	}
}
