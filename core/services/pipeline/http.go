package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
)

var (
	// Client represents a HTTP Client
	Client *http.Client
	// UnrestrictedClient represents a HTTP Client with no Transport restrictions
	UnrestrictedClient *http.Client
)

func newDefaultTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// There are certain classes of vulnerabilities that open up when
	// compression is enabled. For simplicity, we disable compression
	// to cut off this class of attacks.
	// https://www.cyberis.co.uk/2013/08/vulnerabilities-that-just-wont-die.html
	t.DisableCompression = true
	return t
}

func init() {
	tr := newDefaultTransport()
	tr.DialContext = restrictedDialContext
	Client = &http.Client{Transport: tr}

	unrestrictedTr := newDefaultTransport()
	UnrestrictedClient = &http.Client{Transport: unrestrictedTr}
}

// HTTPRequest holds the request and config struct for a http request
type HTTPRequest struct {
	Request *http.Request
	Config  HTTPRequestConfig
	Logger  logger.Logger
}

// HTTPRequestConfig holds the configurable settings for a http request
type HTTPRequestConfig struct {
	SizeLimit                      int64
	AllowUnrestrictedNetworkAccess bool
}

// SendRequest sends a HTTPRequest,
// returns a body, status code, and error.
func (h *HTTPRequest) SendRequest() (responseBody []byte, statusCode int, headers http.Header, err error) {
	var client *http.Client
	if h.Config.AllowUnrestrictedNetworkAccess {
		client = UnrestrictedClient
	} else {
		client = Client
	}
	start := time.Now()

	r, err := client.Do(h.Request)
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
