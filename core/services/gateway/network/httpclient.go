package network

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// HttpClient interfaces defines a method to send HTTP requests
// TODO: handle retries
type HttpClient interface {
	Send(ctx context.Context, req HttpRequest) (*HttpResponse, error)
}

type HttpClientConfig interface {
	MaxResponseBytes() int64
	DefaultTimeout() time.Duration
}

type HttpRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}
type HttpResponse struct {
	StatusCode int               // HTTP status code
	Headers    map[string]string // HTTP headers
	Body       []byte            // HTTP response body
}

type httpClient struct {
	client *http.Client
	config HttpClientConfig
	lggr   logger.Logger
}

// NewHttpClient creates a new HttpClient
// As of now, the client does not support TLS configuration but may be extended in the future
func NewHTTPClient(config HttpClientConfig, lggr logger.Logger) (HttpClient, error) {
	return &httpClient{
		client: &http.Client{
			Timeout:   config.DefaultTimeout(),
			Transport: http.DefaultTransport,
		},
		lggr: lggr,
	}, nil
}

func (c *httpClient) Send(ctx context.Context, req HttpRequest) (*HttpResponse, error) {
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

	reader := http.MaxBytesReader(nil, resp.Body, c.config.MaxResponseBytes())
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
	c.lggr.Debugw("received HTTP response", "statusCode", resp.StatusCode, "body", string(body), "url", req.URL, "headers", headers)

	return &HttpResponse{
		Headers:    headers,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
