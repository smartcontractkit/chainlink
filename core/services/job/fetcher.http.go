package job

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
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/jpillora/backoff"
)

type HttpFetcher struct {
	ID                             uint64          `json:"-" gorm:"primary_key;auto_increment"`
	URL                            models.WebURL   `json:"url"`
	ExtendedPath                   ExtendedPath    `json:"extendedPath,omitempty"`
	Headers                        http.Header     `json:"headers,omitempty"`
	QueryParams                    QueryParameters `json:"queryParams,omitempty"`
	Body                           interface{}     `json:"body,omitempty"`
	AllowUnrestrictedNetworkAccess bool            `json:"-"`

	Config *orm.Config
}

type httpRequestConfig struct {
	timeout                        time.Duration
	maxAttempts                    uint
	sizeLimit                      int64
	allowUnrestrictedNetworkAccess bool
}

func (f *HttpFetcher) Fetch() (interface{}, error) {
	var contentType string
	if f.Method == "POST" {
		contentType = "application/json"
	}

	var body io.Reader
	if f.Body != nil {
		bs, err := json.Marshal(f.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBufferString(bs)
	}

	request, err := http.NewRequest(f.Method, f.URL, body)
	if err != nil {
		return nil, err
	}

	appendExtendedPath(request, f.ExtendedPath)
	appendQueryParams(request, f.QueryParams)
	setHeaders(request, f.Headers, contentType)
	httpConfig := defaultHTTPConfig(store)
	httpConfig.allowUnrestrictedNetworkAccess = f.AllowUnrestrictedNetworkAccess
	return sendRequest(request, httpConfig)
}

func (f HttpFetcher) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion HttpFetcher
	type fetcherWithType struct {
		Type FetcherType `json:"type"`
		preventInfiniteRecursion
	}
	return json.Marshal(fetcherWithType{
		FetcherTypeHttp,
		preventInfiniteRecursion(f),
	})
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

func sendRequest(request *http.Request, config httpRequestConfig) ([]byte, error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	if !config.allowUnrestrictedNetworkAccess {
		tr.DialContext = restrictedDialContext
	}
	client := &http.Client{Transport: tr}

	bs, statusCode, err := withRetry(client, request, config)
	if err != nil {
		return nil, err
	}

	// This is either a client error caused on our end or a server error that persists even after retrying.
	// Either way, there is no way for us to complete the run with a result.
	if statusCode >= 400 {
		return nil, errors.New(string(bs))
	}
	return bs, nil
}

// withRetry executes the http request in a retry. Timeout is controlled with a context
// Retry occurs if the request timeout, or there is any kind of connection or transport-layer error
// Retry also occurs on remote server 5xx errors
func withRetry(
	client *http.Client,
	originalRequest *http.Request,
	config httpRequestConfig,
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
		logger.Debugw("http fetcher error, will retry", "error", err.Error(), "attempt", bb.Attempt(), "timeout", config.timeout)
	}
}

func makeHTTPCall(
	client *http.Client,
	originalRequest *http.Request,
	config httpRequestConfig,
) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	requestWithTimeout := originalRequest.Clone(ctx)

	start := time.Now()

	r, e := client.Do(requestWithTimeout)
	if e != nil {
		return nil, 0, e
	}
	defer logger.ErrorIfCalling(r.Body.Close)

	statusCode := r.StatusCode
	elapsed := time.Since(start)
	logger.Debugw(fmt.Sprintf("http fetcher got %v in %s", statusCode, elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	source := newMaxBytesReader(r.Body, config.sizeLimit)
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		logger.Errorf("http fetcher error reading body: %v", e.Error())
		return nil, statusCode, e
	}
	elapsed = time.Since(start)
	logger.Debugw(fmt.Sprintf("http fetcher finished after %s", elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	responseBody := bytes

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
			return 0, &HTTPResponseTooLargeError{mbr.limit}
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
			return 0, &HTTPResponseTooLargeError{mbr.limit}
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

func defaultHTTPConfig(store *store.Store) httpRequestConfig {
	return httpRequestConfig{
		store.Config.DefaultHTTPTimeout().Duration(),
		store.Config.DefaultMaxHTTPAttempts(),
		store.Config.DefaultHTTPLimit(),
		false,
	}
}
