package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type httpClientConfig interface {
	DatabaseURL() url.URL
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

// HTTPRequest holds the request and config struct for a http request
type HTTPRequest struct {
	Client  *http.Client
	Request *http.Request
	Config  HTTPRequestConfig
	Logger  logger.Logger
}

// HTTPRequestConfig holds the configurable settings for a http request
type HTTPRequestConfig struct {
	SizeLimit int64
}

// SendRequest sends a HTTPRequest,
// returns a body, status code, and error.
func (h *HTTPRequest) SendRequest() (responseBody []byte, statusCode int, headers http.Header, err error) {
	start := time.Now()

	r, err := h.Client.Do(h.Request)
	if err != nil {
		h.Logger.Warnw("http adapter got error", "error", err)
		return nil, 0, nil, err
	}
	defer h.Logger.ErrorIfClosing(r.Body, "SendRequest response body")

	statusCode = r.StatusCode
	elapsed := time.Since(start)
	h.Logger.Debugw(fmt.Sprintf("http adapter got %v in %s", statusCode, elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	source := http.MaxBytesReader(nil, r.Body, h.Config.SizeLimit)
	bytes, err := io.ReadAll(source)
	if err != nil {
		h.Logger.Errorw("http adapter error reading body", "error", err)
		return nil, statusCode, nil, err
	}
	elapsed = time.Since(start)
	h.Logger.Debugw(fmt.Sprintf("http adapter finished after %s", elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	responseBody = bytes

	return responseBody, statusCode, r.Header, nil
}
