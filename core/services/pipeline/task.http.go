package pipeline

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type HTTPTask struct {
	BaseTask

	Method                         string          `json:"method"`
	URL                            models.WebURL   `json:"url"`
	ExtendedPath                   ExtendedPath    `json:"extendedPath"`
	Headers                        Header          `json:"headers"`
	QueryParams                    QueryParameters `json:"queryParams"`
	RequestData                    HttpRequestData `json:"requestData"`
	AllowUnrestrictedNetworkAccess bool            `json:"-"`

	config Config
}

var _ Task = (*HTTPTask)(nil)

type httpRequestConfig struct {
	timeout                        time.Duration
	maxAttempts                    uint
	sizeLimit                      int64
	allowUnrestrictedNetworkAccess bool
}

func (t *HTTPTask) Type() TaskType {
	return TaskTypeHTTP
}

func (f *HTTPTask) Run(inputs []Result) Result {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "HTTPTask requires 0 inputs")}
	}

	var contentType string
	if f.Method == "POST" {
		contentType = "application/json"
	}

	var body io.Reader
	if f.RequestData != nil {
		bs, err := json.Marshal(f.RequestData)
		if err != nil {
			return Result{Error: err}
		}
		body = bytes.NewBuffer(bs)
	}

	request, err := http.NewRequest(f.Method, f.URL.String(), body)
	if err != nil {
		return Result{Error: err}
	}

	appendExtendedPath(request, f.ExtendedPath)
	appendQueryParams(request, f.QueryParams)
	setHeaders(request, http.Header(f.Headers), contentType)
	httpConfig := httpRequestConfig{
		f.config.DefaultHTTPTimeout().Duration(),
		f.config.DefaultMaxHTTPAttempts(),
		f.config.DefaultHTTPLimit(),
		false,
	}
	httpConfig.allowUnrestrictedNetworkAccess = f.AllowUnrestrictedNetworkAccess
	resp, err := sendRequest(request, httpConfig)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: resp}
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

func (qp *QueryParameters) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), qp)
}
func (qp QueryParameters) Value() (driver.Value, error) {
	return json.Marshal(qp)
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

func (ep *ExtendedPath) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ep)
}
func (ep ExtendedPath) Value() (driver.Value, error) {
	return json.Marshal(ep)
}

type Header http.Header

func (h *Header) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), h)
}
func (h Header) Value() (driver.Value, error) {
	return json.Marshal(h)
}

type HttpRequestData map[string]interface{}

func (h *HttpRequestData) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), h)
}
func (h HttpRequestData) Value() (driver.Value, error) {
	return json.Marshal(h)
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isRestrictedIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() ||
		ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.Equal(net.IPv4bcast) ||
		ip.Equal(net.IPv4allsys) ||
		ip.Equal(net.IPv4allrouter) ||
		ip.Equal(net.IPv4zero) ||
		ip.IsMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// restrictedDialContext wraps the Dialer such that after successful connection,
// we check the IP.
// If the resolved IP is restricted, close the connection and return an error.
func restrictedDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	con, err := (&net.Dialer{
		// Defaults from GoLang standard http package
		// https://golang.org/pkg/net/http/#RoundTripper
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext(ctx, network, address)
	if err == nil {
		// If a connection could be established, ensure its not local or private
		a, _ := con.RemoteAddr().(*net.TCPAddr)

		if isRestrictedIP(a.IP) {
			defer logger.ErrorIfCalling(con.Close)
			return nil, fmt.Errorf("disallowed IP %s. Connections to local/private and multicast networks are disabled by default for security reasons. If you really want to allow this, consider using the httpgetwithunrestrictednetworkaccess or httppostwithunrestrictednetworkaccess adapter instead", a.IP.String())
		}
	}
	return con, err
}
