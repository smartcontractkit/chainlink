package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

type IHttpClient interface {
	// Get issue a GET request to the given url and return the response body and status code.
	Get(ctx context.Context, url string, timeout time.Duration) ([]byte, int, http.Header, error)
}

type HttpClient struct {
}

func (s *HttpClient) Get(ctx context.Context, url string, timeout time.Duration) ([]byte, int, http.Header, error) {
	// Use a timeout to guard against attestation API hanging, causing observation timeout and failing to make any progress.
	timeoutCtx, cancel := context.WithTimeoutCause(ctx, timeout, tokendata.ErrTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusBadRequest, nil, err
	}
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, http.StatusRequestTimeout, nil, tokendata.ErrTimeout
		}
		// On error, res is nil in most cases, do not read res.StatusCode, return BadRequest
		return nil, http.StatusBadRequest, nil, err
	}
	defer res.Body.Close()

	// Explicitly signal if the API is being rate limited
	if res.StatusCode == http.StatusTooManyRequests {
		return nil, res.StatusCode, res.Header, tokendata.ErrRateLimit
	}

	body, err := io.ReadAll(res.Body)
	return body, res.StatusCode, res.Header, err
}
