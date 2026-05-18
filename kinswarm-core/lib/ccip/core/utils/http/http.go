package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type httpClientConfig interface {
	URL() url.URL // DatabaseURL
}

// NewRestrictedHTTPClient returns a secure HTTP Client (queries to certain
// local addresses are blocked)
func NewRestrictedHTTPClient(cfg httpClientConfig, lggr logger.Logger) *http.Client {
	tr := newDefaultTransport()
	tr.DialContext = makeRestrictedDialContext(cfg, lggr)
	return &http.Client{Transport: tr}
}

// NewUnrestrictedClient returns a HTTP Client with no Transport restrictions
func NewUnrestrictedHTTPClient() *http.Client {
	unrestrictedTr := newDefaultTransport()
	return &http.Client{Transport: unrestrictedTr}
}

func newDefaultTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// There are certain classes of vulnerabilities that open up when
	// compression is enabled. For simplicity, we disable compression
	// to cut off this class of attacks.
	// https://www.cyberis.co.uk/2013/08/vulnerabilities-that-just-wont-die.html
	t.DisableCompression = true
	return t
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPRequest holds the request and config struct for a http request
type HTTPRequest struct {
	Client  HTTPClient
	Request *http.Request
	Config  HTTPRequestConfig
	Logger  logger.Logger
}

// HTTPRequestConfig holds the configurable settings for a http request
type HTTPRequestConfig struct {
	// SizeLimit in bytes
	SizeLimit int64
}

// SendRequestReader allows for streaming the body directly and does not read
// it all into memory
//
// CALLER IS RESPONSIBLE FOR CLOSING RETURNED RESPONSE BODY
func (h *HTTPRequest) SendRequestReader() (responseBody io.ReadCloser, statusCode int, headers http.Header, err error) {
	start := time.Now()
	r, err := h.Client.Do(h.Request)
	if err != nil {
		logger.Sugared(h.Logger).Tracew("http adapter got error", "err", err)
		return nil, 0, nil, err
	}

	statusCode = r.StatusCode
	elapsed := time.Since(start)
	logger.Sugared(h.Logger).Tracew(fmt.Sprintf("http adapter got %v in %s", statusCode, elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	source := http.MaxBytesReader(nil, r.Body, h.Config.SizeLimit)

	return source, statusCode, r.Header, nil
}

// SendRequest sends a HTTPRequest,
// returns a body, status code, and error.
func (h *HTTPRequest) SendRequest() (responseBody []byte, statusCode int, headers http.Header, err error) {
	start := time.Now()

	source, statusCode, headers, err := h.SendRequestReader()
	if err != nil {
		return nil, statusCode, headers, err
	}
	defer logger.Sugared(h.Logger).ErrorIfFn(source.Close, "Error closing SendRequest response body")
	bytes, err := io.ReadAll(source)
	elapsed := time.Since(start)
	logger.Sugared(h.Logger).Tracew(fmt.Sprintf("http adapter finished after %s", elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	responseBody = bytes

	return responseBody, statusCode, headers, nil
}
