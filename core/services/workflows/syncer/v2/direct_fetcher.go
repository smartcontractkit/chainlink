package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var ErrHostNotAllowed = errors.New("host is not in the allowed fetch hosts")

// NewDirectHTTPClient returns an HTTP client for fetching workflow artifacts without the Gateway.
// Only HTTPS on port 443 is allowed. Internal IPs are blocked and redirects are disabled.
// timeout applies to requests without their own timeout. maxResponseBytes caps every response; requests can only lower it.
func NewDirectHTTPClient(lggr logger.Logger, timeout time.Duration, maxResponseBytes uint32) (network.HTTPClient, error) {
	return network.NewHTTPClient(network.HTTPClientConfig{
		AllowedSchemes:   []string{"https"},
		AllowedPorts:     []int{443},
		AllowedMethods:   []string{http.MethodGet},
		DefaultTimeout:   timeout,
		MaxResponseBytes: maxResponseBytes,
	}, lggr)
}

// NewDirectFetcherFunc returns a FetcherFunc that downloads artifacts straight from the node.
// If allowedHosts is not empty, only those hosts (case-insensitive, exact match) can be fetched.
func NewDirectFetcherFunc(client network.HTTPClient, allowedHosts []string, lggr logger.Logger) types.FetcherFunc {
	allowed := make(map[string]struct{}, len(allowedHosts))
	for _, h := range allowedHosts {
		allowed[strings.ToLower(h)] = struct{}{}
	}
	lggr = logger.Named(lggr, "DirectFetcher")

	return func(ctx context.Context, _ string, req ghcapabilities.Request) ([]byte, error) {
		if req.WorkflowID == "" {
			return nil, errors.New("invalid call to fetch, must provide workflow ID")
		}

		u, err := url.Parse(req.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		if len(allowed) > 0 {
			if _, ok := allowed[strings.ToLower(u.Hostname())]; !ok {
				return nil, fmt.Errorf("%w: %s", ErrHostNotAllowed, u.Hostname())
			}
		}

		resp, err := client.Send(ctx, network.HTTPRequest{
			Method:           http.MethodGet,
			URL:              req.URL,
			Timeout:          time.Duration(req.TimeoutMs) * time.Millisecond,
			MaxResponseBytes: req.MaxResponseBytes,
		})
		if err != nil {
			return nil, fmt.Errorf("direct fetch failed: %w", err)
		}
		lggr.Debugw("fetched artifact", "host", u.Hostname(), "workflowID", req.WorkflowID, "status", resp.StatusCode)

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return resp.Body, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		}
		return resp.Body, nil
	}
}
