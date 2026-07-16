package network

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptrace"
	"slices"
	"strings"

	"time"

	"github.com/docker/go-connections/nat"

	"github.com/doyensec/safeurl"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
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
	// AllowedPortRanges accepts port ranges in "start-end" format (e.g. "8000-9000").
	// Expanded ports are merged into AllowedPorts during construction.
	AllowedPortRanges []string
	AllowedSchemes    []string
	AllowedIPs        []string
	AllowedIPsCIDR    []string
	AllowedMethods    []string
	BlockedHeaders    []string

	// Mtls, when set, configures the client to present the supplied client
	// certificate for mutual TLS.
	Mtls *gateway.MtlsAuth

	// ConcurrencyLimiter, when set together with Mtls, bounds the number of
	// in-flight mTLS requests. The limiter is acquired on the request's
	// (capped) context so waiters self-evict at the request timeout.
	ConcurrencyLimiter limits.ResourcePoolLimiter[int]
}

// merge returns a copy of c with any set fields from override applied on top.
// A field in override is only applied when it holds a non-zero value, so the
// static base config supplies defaults that the dynamic config can selectively
// override.
func (c HTTPClientConfig) merge(override HTTPClientConfig) HTTPClientConfig {
	merged := c
	if override.MaxResponseBytes != 0 {
		merged.MaxResponseBytes = override.MaxResponseBytes
	}
	if override.DefaultTimeout != 0 {
		merged.DefaultTimeout = override.DefaultTimeout
	}
	if override.maxRequestDuration != 0 {
		merged.maxRequestDuration = override.maxRequestDuration
	}
	if len(override.BlockedIPs) > 0 {
		merged.BlockedIPs = override.BlockedIPs
	}
	if len(override.BlockedIPsCIDR) > 0 {
		merged.BlockedIPsCIDR = override.BlockedIPsCIDR
	}
	if len(override.AllowedPorts) > 0 {
		merged.AllowedPorts = override.AllowedPorts
	}
	if len(override.AllowedPortRanges) > 0 {
		merged.AllowedPortRanges = override.AllowedPortRanges
	}
	if len(override.AllowedSchemes) > 0 {
		merged.AllowedSchemes = override.AllowedSchemes
	}
	if len(override.AllowedIPs) > 0 {
		merged.AllowedIPs = override.AllowedIPs
	}
	if len(override.AllowedIPsCIDR) > 0 {
		merged.AllowedIPsCIDR = override.AllowedIPsCIDR
	}
	if len(override.AllowedMethods) > 0 {
		merged.AllowedMethods = override.AllowedMethods
	}
	if len(override.BlockedHeaders) > 0 {
		merged.BlockedHeaders = override.BlockedHeaders
	}
	if override.Mtls != nil {
		merged.Mtls = override.Mtls
	}
	if override.ConcurrencyLimiter != nil {
		merged.ConcurrencyLimiter = override.ConcurrencyLimiter
	}
	return merged
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
	ErrBlockedRequest         = errors.New("blocked request")
	ErrHTTPSend               = errors.New("failed to send HTTP request")
	ErrHTTPRead               = errors.New("failed to read HTTP response body")
)

// expandPortRanges parses all port range strings and returns the combined list of individual ports.
// Each range string is parsed by nat.ParsePortRange (e.g. "8000-9000" or a single port "443").
func expandPortRanges(ranges []string) ([]int, error) {
	var ports []int
	for _, r := range ranges {
		start, end, err := nat.ParsePortRange(r)
		if err != nil {
			return nil, fmt.Errorf("invalid port range %q: %w", r, err)
		}
		if start < 1 {
			return nil, fmt.Errorf("port range %q: start port must be >= 1", r)
		}
		for p := start; p <= end; p++ {
			ports = append(ports, int(p)) //nolint:gosec // port values are validated above to be >= 1 and within uint16 range
		}
	}
	return ports, nil
}

func (c *HTTPClientConfig) ApplyDefaults() {
	if len(c.AllowedPorts) == 0 {
		c.AllowedPorts = slices.Clone(defaultAllowedPorts)
	}
	if len(c.AllowedSchemes) == 0 {
		c.AllowedSchemes = slices.Clone(defaultAllowedSchemes)
	}
	if len(c.AllowedMethods) == 0 {
		c.AllowedMethods = slices.Clone(defaultAllowedMethods)
	}
	if len(c.BlockedHeaders) == 0 {
		c.BlockedHeaders = slices.Clone(defaultBlockedHeaders)
	}
	if c.MaxResponseBytes == 0 {
		c.MaxResponseBytes = defaultMaxResponseBytes
	}
	if c.DefaultTimeout == 0 {
		c.DefaultTimeout = defaultTimeout
	}
	c.maxRequestDuration = defaultMaxRequestDuration
	// safeurl automatically blocks internal IPs so no need to set defaults here.
}

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string // request headers (deprecated: use MultiHeaders when multiple values per key are needed)
	// MultiHeaders holds multiple values per header name; when set, Headers is ignored for the outgoing request.
	MultiHeaders map[string][]string
	Body         []byte
	Timeout      time.Duration

	// Maximum number of bytes to read from the response body.  If 0, the default value is used.
	// Does not override a request specific value gte 0.
	MaxResponseBytes uint32
}

type HTTPResponse struct {
	StatusCode   int                 // HTTP status code
	Headers      map[string]string   // HTTP headers (deprecated: use MultiHeaders, contains first value only for backward compatibility)
	MultiHeaders map[string][]string // HTTP headers with all values preserved
	Body         []byte              // HTTP response body
}

// requestToNetHeader builds net/http.Header from req. Uses MultiHeaders when set, otherwise Headers.
func requestToNetHeader(req HTTPRequest) http.Header {
	out := make(http.Header)
	if len(req.MultiHeaders) > 0 {
		for k, values := range req.MultiHeaders {
			for _, v := range values {
				out.Add(k, v)
			}
		}
		return out
	}
	for k, v := range req.Headers {
		out.Add(k, v)
	}
	return out
}

// responseHeadersFromNetHeader builds Headers (comma-joined) and MultiHeaders from net/http.Header. Skips keys with no values.
func responseHeadersFromNetHeader(h http.Header) (map[string]string, map[string][]string) {
	headers := make(map[string]string, len(h))
	multiHeaders := make(map[string][]string, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		multiHeaders[k] = slices.Clone(v)
		headers[k] = strings.Join(v, ",")
	}
	return headers, multiHeaders
}

// httpDoer is the subset of the HTTP client used by httpClient. It is satisfied
// by *safeurl.WrappedClient and by concurrencyLimitedClient, which decorates it.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// concurrencyLimitedClient bounds the number of in-flight requests delegated to
// the underlying client. The slot is acquired on the request's context, so a
// waiter self-evicts when that context (carrying the capped request timeout)
// expires rather than blocking indefinitely.
type concurrencyLimitedClient struct {
	client  httpDoer
	limiter limits.ResourcePoolLimiter[int]
}

func (c *concurrencyLimitedClient) Do(req *http.Request) (*http.Response, error) {
	free, err := c.limiter.Wait(req.Context(), 1)
	if err != nil {
		return nil, fmt.Errorf("mtls concurrency limit exceeded: %w", ErrBlockedRequest)
	}
	defer free()
	return c.client.Do(req)
}

type httpClient struct {
	client  httpDoer
	config  HTTPClientConfig
	lggr    logger.Logger
	metrics *httpClientMetrics
}

// NewHTTPClient creates a new HTTPClient. For mTLS support, set Mtls on the
// config, or use NewHTTPClientFactory to merge a dynamic config carrying it.
func NewHTTPClient(config HTTPClientConfig, lggr logger.Logger) (HTTPClient, error) {
	if len(config.AllowedPortRanges) > 0 {
		expanded, err := expandPortRanges(config.AllowedPortRanges)
		if err != nil {
			return nil, fmt.Errorf("invalid AllowedPortRanges: %w", err)
		}
		config.AllowedPorts = append(config.AllowedPorts, expanded...)
	}
	config.ApplyDefaults()

	dt, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("could not coerce http.DefaultTransport to *http.Transport")
	}
	defaultTransport := dt.Clone()

	safeConfigBuilder := safeurl.
		GetConfigBuilder().
		SetAllowedIPs(config.AllowedIPs...).
		SetAllowedIPsCIDR(config.AllowedIPsCIDR...).
		SetAllowedPorts(config.AllowedPorts...).
		SetAllowedSchemes(config.AllowedSchemes...).
		SetBlockedIPs(config.BlockedIPs...).
		SetBlockedIPsCIDR(config.BlockedIPsCIDR...).
		SetCheckRedirect(disableRedirects).
		SetTransport(defaultTransport)

	var client httpDoer

	if config.Mtls != nil {
		// Defence-in-depth protection against accidental reuse
		// of the HTTP client leading to auth'd connections leaking across
		// users.
		defaultTransport.DisableKeepAlives = true
		defaultTransport.TLSHandshakeTimeout = 10 * time.Second

		cert, err := tls.X509KeyPair(config.Mtls.Certificate, config.Mtls.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MtlsAuth into KeyPair: %w", err)
		}

		defaultTransport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		safeConfigBuilder.SetTransport(defaultTransport)

		if config.ConcurrencyLimiter == nil {
			return nil, errors.New("mtls requires a ConcurrencyLimiter")
		}
		client = &concurrencyLimitedClient{
			client:  safeurl.Client(safeConfigBuilder.Build()),
			limiter: config.ConcurrencyLimiter,
		}
	} else {
		client = safeurl.Client(safeConfigBuilder.Build())
	}

	metrics, err := newHTTPClientMetrics()
	if err != nil {
		return nil, err
	}

	return &httpClient{
		config:  config,
		client:  client,
		lggr:    lggr,
		metrics: metrics,
	}, nil
}

// NewHTTPClientFactory returns a factory that builds an HTTPClient by merging
// the supplied dynamic config on top of the static base config. The dynamic
// config typically carries per-request settings such as Mtls.
func NewHTTPClientFactory(config HTTPClientConfig, lggr logger.Logger) HTTPClientFactory {
	return func(dynamic HTTPClientConfig) (HTTPClient, error) {
		return NewHTTPClient(config.merge(dynamic), lggr)
	}
}

// HTTPClientFactory builds an HTTPClient from a dynamic config that is merged
// onto the factory's static base config. Only fields holding a non-zero value
// in the dynamic config override the base; zero-valued fields fall through to
// the base, so the dynamic config can add or override settings but cannot unset
// one already populated by the base.
type HTTPClientFactory func(config HTTPClientConfig) (HTTPClient, error)

func disableRedirects(req *http.Request, via []*http.Request) error {
	return &redirectsDisabledError{}
}

type redirectsDisabledError struct{}

func (e *redirectsDisabledError) Error() string { return "redirects are not allowed" }

// isBlockedRequest checks if an error is caused by blocked/invalid input (e.g., blocked IP, invalid scheme, blocked headers)
// It checks for safeurl typed errors.
func isBlockedRequest(err error) bool {
	if err == nil {
		return false
	}

	// Check safeurl typed errors - use errors.As for type checking
	var (
		ipv6Err              *safeurl.IPv6BlockedError
		portErr              *safeurl.AllowedPortError
		schemeErr            *safeurl.AllowedSchemeError
		invalidHostErr       *safeurl.InvalidHostError
		ipErr                *safeurl.AllowedIPError
		redirectsDisabledErr *redirectsDisabledError
	)

	return errors.As(err, &ipv6Err) ||
		errors.As(err, &portErr) ||
		errors.As(err, &schemeErr) ||
		errors.As(err, &invalidHostErr) ||
		errors.As(err, &ipErr) ||
		errors.As(err, &redirectsDisabledErr)
}

func (c *httpClient) validateMethod(method string) error {
	isAllowed := func(allowed string) bool {
		return strings.EqualFold(allowed, method)
	}
	if slices.ContainsFunc(c.config.AllowedMethods, isAllowed) {
		return nil
	}
	return fmt.Errorf("%w: HTTP method not allowed: %s", ErrBlockedRequest, method)
}

// validateHeaderNames checks that none of the given header names are in the blocked list (case-insensitive).
func (c *httpClient) validateHeaderNames(names []string) error {
	blockedSet := make(map[string]struct{}, len(c.config.BlockedHeaders))
	for _, b := range c.config.BlockedHeaders {
		blockedSet[strings.ToLower(b)] = struct{}{}
	}
	for _, name := range names {
		if _, blocked := blockedSet[strings.ToLower(name)]; blocked {
			return fmt.Errorf("%w: HTTP header not allowed: %s", ErrBlockedRequest, name)
		}
	}
	return nil
}

func (c *httpClient) validateHeaders(headers map[string]string) error {
	return c.validateHeaderNames(slices.Collect(maps.Keys(headers)))
}

func (c *httpClient) validateMultiHeaders(multiHeaders map[string][]string) error {
	return c.validateHeaderNames(slices.Collect(maps.Keys(multiHeaders)))
}

// Send executes an http request that is always time limited by at least the
// default timeout.  Override the default timeout with a non-zero duration by
// passing a Timeout value on the request.
func (c *httpClient) Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
	if err := c.validateMethod(req.Method); err != nil {
		return nil, err
	}
	if len(req.MultiHeaders) > 0 {
		if err := c.validateMultiHeaders(req.MultiHeaders); err != nil {
			return nil, err
		}
	} else if err := c.validateHeaders(req.Headers); err != nil {
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

	requestStart := time.Now()
	trace, traceState := newClientTrace(ctx, req.Method, requestStart, c.metrics)
	timeoutCtx = httptrace.WithClientTrace(timeoutCtx, trace)

	r, err := http.NewRequestWithContext(timeoutCtx, req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		c.metrics.recordTotal(ctx, req.Method, 0, false, false, time.Since(requestStart))
		return nil, err
	}
	for k, values := range requestToNetHeader(req) {
		for _, v := range values {
			r.Header.Add(k, v)
		}
	}

	resp, err := c.client.Do(r)
	if err != nil {
		c.metrics.recordTotal(ctx, req.Method, 0, false, traceState.connReused.Load(), time.Since(requestStart))
		if isBlockedRequest(err) {
			c.lggr.Warnw("HTTP request blocked", "err", err)
			return nil, fmt.Errorf("%w: %w", ErrBlockedRequest, err)
		}
		c.lggr.Errorw("failed to send HTTP request", "err", err)
		return nil, errors.Join(err, ErrHTTPSend)
	}
	defer resp.Body.Close()

	n := maxReadBytes(readSize{defaultSize: c.config.MaxResponseBytes, requestSize: req.MaxResponseBytes})
	c.lggr.Debugw("max bytes to read from HTTP response", "bytes", n)

	reader := http.MaxBytesReader(nil, resp.Body, int64(n))
	body, err := io.ReadAll(reader)
	if err != nil {
		c.metrics.recordTotal(ctx, req.Method, resp.StatusCode, false, traceState.connReused.Load(), time.Since(requestStart))
		c.lggr.Errorw("failed to read HTTP response body", "err", err)
		return nil, errors.Join(err, ErrHTTPRead)
	}

	c.metrics.recordTotal(ctx, req.Method, resp.StatusCode, true, traceState.connReused.Load(), time.Since(requestStart))

	headers, multiHeaders := responseHeadersFromNetHeader(resp.Header)
	c.lggr.Debugw("received HTTP response", "statusCode", resp.StatusCode)
	return &HTTPResponse{
		Headers:      headers,
		MultiHeaders: multiHeaders,
		StatusCode:   resp.StatusCode,
		Body:         body,
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
	return min(sizes.defaultSize, sizes.requestSize)
}
