package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

// HTTPGet requires a URL which is used for a GET request when the adapter is called.
type HTTPGet struct {
	URL          models.WebURL   `json:"url"`
	GET          models.WebURL   `json:"get"`
	Headers      http.Header     `json:"headers"`
	QueryParams  QueryParameters `json:"queryParams"`
	ExtendedPath ExtendedPath    `json:"extPath"`
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunResult, store *store.Store) models.RunResult {
	request, err := hga.GetRequest()
	if err != nil {
		input.SetError(err)
		return input
	}
	return sendRequest(input, request, store.Config.DefaultHTTPLimit())
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
	URL          models.WebURL   `json:"url"`
	POST         models.WebURL   `json:"post"`
	Headers      http.Header     `json:"headers"`
	QueryParams  QueryParameters `json:"queryParams"`
	Body         *string         `json:"body,omitempty"`
	ExtendedPath ExtendedPath    `json:"extPath"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunResult, store *store.Store) models.RunResult {
	request, err := hpa.GetRequest(input.Data.String())
	if err != nil {
		input.SetError(err)
		return input
	}
	return sendRequest(input, request, store.Config.DefaultHTTPLimit())
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

func sendRequest(input models.RunResult, request *http.Request, limit int64) models.RunResult {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		input.SetError(err)
		return input
	}

	defer response.Body.Close()

	source := newMaxBytesReader(response.Body, limit)
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		input.SetError(err)
		return input
	}

	responseBody := string(bytes)
	if response.StatusCode >= 400 {
		input.SetError(errors.New(responseBody))
		return input
	}

	input.CompleteWithResult(responseBody)
	return input
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

func (mbr *maxBytesReader) tooLarge() (int, error) {
	return 0, fmt.Errorf("HTTP request too large, must be less than %d bytes", mbr.limit)
}

func (mbr *maxBytesReader) Close() error {
	return mbr.rc.Close()
}

// QueryParameters are the keys and values to append to the URL
type QueryParameters url.Values

// UnmarshalJSON implements the Unmarshaler interface
func (qp *QueryParameters) UnmarshalJSON(input []byte) error {
	values := url.Values{}
	strs := []string{}
	var err error

	// input is a string like "someKey0=someVal0&someKey1=someVal1"
	if utils.IsQuoted(input) {
		var decoded string
		err := json.Unmarshal(input, &decoded)
		if err != nil {
			return fmt.Errorf("unable to unmarshal query parameters: %s", input)
		}
		strs = strings.FieldsFunc(trimQuestion(decoded), splitQueryString)

		// input is an array of strings like
		// ["someKey0", "someVal0", "someKey1", "someVal1"]
	} else {
		err = json.Unmarshal(input, &strs)
	}

	values, err = buildValues(strs)
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
