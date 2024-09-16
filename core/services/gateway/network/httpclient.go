package network

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// HttpClient interfaces defines a method to send HTTP requests
// handles retries and timeouts
type HttpClient interface {
	Send(ctx context.Context, req HttpRequest) (*HttpResponse, error)
}

type HttpRequest struct {
	Method     string
	URL        string
	Headers    map[string]string
	Body       []byte
	Timeout    time.Duration
	RetryCount uint8
}

type HttpResponse struct {
	StatusCode int               // HTTP status code
	Headers    map[string]string // HTTP headers
	Body       []byte            // Base64-encoded binary body
}

type httpClient struct {
	client *http.Client
	lggr   logger.Logger
}

func NewHttpClient(tlsConfig *tls.Config, defaultTimeout time.Duration, lggr logger.Logger) (HttpClient, error) {
	transport := http.DefaultTransport
	if tlsConfig != nil {
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}
	return &httpClient{
		client: &http.Client{
			Timeout:   defaultTimeout,
			Transport: transport,
		},
		lggr: lggr,
	}, nil
}

func (c *httpClient) Send(ctx context.Context, req HttpRequest) (*HttpResponse, error) {
	r, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range req.Headers {
		r.Header.Add(k, v)
	}
	retryCount := uint8(0)
	for {
		timeoutCtx, cancel := context.WithTimeout(ctx, req.Timeout)
		defer cancel()
		r = r.WithContext(timeoutCtx)
		resp, err := c.client.Do(r)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		headers := make(map[string]string)
		for k, v := range resp.Header {
			// header values are usually an array of size 1
			// joining them to a single string in case array size is greater than 1
			headers[k] = strings.Join(v, ",")
		}
		l := logger.With(c.lggr, "statusCode", resp.StatusCode, "body", string(body), "url", req.URL, "attempt", retryCount+1, "headers", headers)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			l.Debugw("received successful response from user")
			return &HttpResponse{
				StatusCode: resp.StatusCode,
				Headers:    headers,
				Body:       body,
			}, nil
		} else {
			l.Warnw("received unsuccessful response from user")
			if retryCount >= req.RetryCount {
				return &HttpResponse{
					StatusCode: resp.StatusCode,
					Headers:    headers,
					Body:       body,
				}, nil
			}
			retryCount++
		}
	}
}
