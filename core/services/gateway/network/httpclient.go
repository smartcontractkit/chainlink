package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/doyensec/safeurl"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// HTTPClient interfaces defines a method to send HTTP requests
type HTTPClient interface {
	Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error)
}

type HTTPClientConfig struct {
	MaxResponseBytes uint32
	DefaultTimeout   time.Duration

	// An HTTPRequest may override the DefaultTimeout, but is capped by
	// maxRequestDuration.
	maxRequestDuration time.Duration
	BlockedIPs         []string
	BlockedIPsCIDR     []string
	AllowedPorts       []int
	AllowedSchemes     []string
	AllowedIPs         []string
	AllowedIPsCIDR     []string
	AllowedMethods     []string
	BlockedHeaders     []string
}

var (
	defaultAllowedPorts   = []int{80, 443}
	defaultAllowedSchemes = []string{"http", "https"}
	defaultAllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	defaultBlockedHeaders = []string{
		"host",              // target host is set in the http client
		"content-length",    // length is computed from actual body to ensure integrity
		"transfer-encoding", // http client manages encoding based on actual content
		"user-agent",        // gateway controls its own identification to backend services
		"upgrade",           // prevents protocol upgrade attacks
		"expect",            // prevents 100-continue exploitation
		"connection",        // external developers cannot control connection behavior or persistence
		"keep-alive",        // gateway manages its own connection pooling and timeouts
		"te",                // blocks attempts to manipulate how request bodies are processed
		"trailer",           // blocks delayed header injection after request body
		"x-forwarded-for",   // prevents IP spoofing
		"x-forwarded-host",  // prevents host header spoofing
		"x-forwarded-proto", // prevents protocol spoofing
		"x-real-ip",         // prevents IP address spoofing
	}
	defaultMaxResponseBytes   = uint32(26.4 * utils.KB)
	defaultMaxRequestDuration = 60 * time.Second
	defaultTimeout            = 5 * time.Second
	ErrHTTPSend               = errors.New("failed to send HTTP request")
	ErrHTTPRead               = errors.New("failed to read HTTP response body")
	ErrInvalidMethod          = errors.New("HTTP method not allowed")
	ErrBlockedHeader          = errors.New("HTTP header not allowed")
)

func (c *HTTPClientConfig) ApplyDefaults() {
	if len(c.AllowedPorts) == 0 {
		c.AllowedPorts = defaultAllowedPorts
	}

	if len(c.AllowedSchemes) == 0 {
		c.AllowedSchemes = defaultAllowedSchemes
	}

	if len(c.AllowedMethods) == 0 {
		c.AllowedMethods = defaultAllowedMethods
	}

	if len(c.BlockedHeaders) == 0 {
		c.BlockedHeaders = defaultBlockedHeaders
	}

	if c.MaxResponseBytes == 0 {
		c.MaxResponseBytes = defaultMaxResponseBytes
	}

	if c.DefaultTimeout == 0 {
		c.DefaultTimeout = defaultTimeout
	}

	c.maxRequestDuration = defaultMaxRequestDuration

	// safeurl automatically blocks internal IPs so no need
	// to set defaults here.
}

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration

	// Maximum number of bytes to read from the response body.  If 0, the default value is used.
	// Does not override a request specific value gte 0.
	MaxResponseBytes uint32
}

type HTTPResponse struct {
	StatusCode int               // HTTP status code
	Headers    map[string]string // HTTP headers
	Body       []byte            // HTTP response body
}

type httpClient struct {
	client *safeurl.WrappedClient
	config HTTPClientConfig
	lggr   logger.Logger
}

// NewHTTPClient creates a new NewHTTPClient
// As of now, the client does not support TLS configuration but may be extended in the future
func NewHTTPClient(config HTTPClientConfig, lggr logger.Logger) (HTTPClient, error) {
	config.ApplyDefaults()
	safeConfig := safeurl.
		GetConfigBuilder().
		SetAllowedIPs(config.AllowedIPs...).
		SetAllowedIPsCIDR(config.AllowedIPsCIDR...).
		SetAllowedPorts(config.AllowedPorts...).
		SetAllowedSchemes(config.AllowedSchemes...).
		SetBlockedIPs(config.BlockedIPs...).
		SetBlockedIPsCIDR(config.BlockedIPsCIDR...).
		SetCheckRedirect(disableRedirects).
		Build()

	return &httpClient{
		config: config,
		client: safeurl.Client(safeConfig),
		lggr:   lggr,
	}, nil
}

func disableRedirects(req *http.Request, via []*http.Request) error {
	return errors.New("redirects are not allowed")
}

func (c *httpClient) validateMethod(method string) error {
	methodUpper := strings.ToUpper(method)
	for _, allowedMethod := range c.config.AllowedMethods {
		if strings.ToUpper(allowedMethod) == methodUpper {
			return nil
		}
	}
	return fmt.Errorf("HTTP method not allowed: %s", method)
}

func (c *httpClient) validateHeaders(headers map[string]string) error {
	for headerName := range headers {
		headerNameLower := strings.ToLower(headerName)
		for _, blockedHeader := range c.config.BlockedHeaders {
			if strings.ToLower(blockedHeader) == headerNameLower {
				return fmt.Errorf("HTTP header not allowed: %s", headerName)
			}
		}
	}
	return nil
}

// Send executes an http request that is always time limited by at least the
// default timeout.  Override the default timeout with a non-zero duration by
// passing a Timeout value on the request.
func (c *httpClient) Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
	if err := c.validateMethod(req.Method); err != nil {
		return nil, err
	}
	if err := c.validateHeaders(req.Headers); err != nil {
		return nil, err
	}

	to := req.Timeout
	if to == 0 {
		to = c.config.DefaultTimeout
	}

	if to > c.config.maxRequestDuration {
		to = c.config.maxRequestDuration
	}

	c.lggr.Debugw("sending HTTP request with timeout", "request timeout", to)

	timeoutCtx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	r, err := http.NewRequestWithContext(timeoutCtx, req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		r.Header.Add(k, v)
	}

	resp, err := c.client.Do(r)
	if err != nil {
		c.lggr.Errorw("failed to send HTTP request", "url", req.URL, "err", err)
		return nil, errors.Join(err, ErrHTTPSend)
	}
	defer resp.Body.Close()

	n := maxReadBytes(readSize{defaultSize: c.config.MaxResponseBytes, requestSize: req.MaxResponseBytes})
	c.lggr.Debugw("max bytes to read from HTTP response", "bytes", n)

	reader := http.MaxBytesReader(nil, resp.Body, int64(n))
	body, err := io.ReadAll(reader)
	if err != nil {
		c.lggr.Errorw("failed to read HTTP response body", "url", req.URL, "err", err)
		return nil, errors.Join(err, ErrHTTPRead)
	}
	headers := make(map[string]string)
	for k, v := range resp.Header {
		// header values are usually an array of size 1
		// joining them to a single string in case array size is greater than 1
		headers[k] = strings.Join(v, ",")
	}
	c.lggr.Debugw("received HTTP response", "statusCode", resp.StatusCode)

	return &HTTPResponse{
		Headers:    headers,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

type readSize struct {
	defaultSize uint32
	requestSize uint32
}

func maxReadBytes(sizes readSize) uint32 {
	if sizes.requestSize == 0 {
		return sizes.defaultSize
	}
	return minUint32(sizes.defaultSize, sizes.requestSize)
}

func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
