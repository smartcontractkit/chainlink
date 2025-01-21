package network

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/doyensec/safeurl"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// HTTPClient interfaces defines a method to send HTTP requests
type HTTPClient interface {
	Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error)
}

type HTTPClientConfig struct {
	MaxResponseBytes uint32
	DefaultTimeout   time.Duration
	BlockedIPs       []string
	BlockedIPsCIDR   []string
	AllowedPorts     []int
	AllowedSchemes   []string

	// for testing
	allowedIPs []string
}

var (
	defaultAllowedPorts   = []int{80, 443}
	defaultAllowedSchemes = []string{"http", "https"}
)

func (c *HTTPClientConfig) ApplyDefaults() {
	if len(c.AllowedPorts) == 0 {
		c.AllowedPorts = defaultAllowedPorts
	}

	if len(c.AllowedSchemes) == 0 {
		c.AllowedSchemes = defaultAllowedSchemes
	}

	// safeurl automatically blocks internal IPs so no need
	// to set defaults here.
}

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
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
		SetTimeout(config.DefaultTimeout).
		SetAllowedIPs(config.allowedIPs...).
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

func (c *httpClient) Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()
	r, err := http.NewRequestWithContext(timeoutCtx, req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := http.MaxBytesReader(nil, resp.Body, int64(c.config.MaxResponseBytes))
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	for k, v := range resp.Header {
		// header values are usually an array of size 1
		// joining them to a single string in case array size is greater than 1
		headers[k] = strings.Join(v, ",")
	}
	c.lggr.Debugw("received HTTP response", "statusCode", resp.StatusCode, "url", req.URL, "headers", headers)

	return &HTTPResponse{
		Headers:    headers,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
