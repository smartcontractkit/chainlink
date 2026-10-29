package v2

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type stubHTTPClient struct {
	resp *network.HTTPResponse
	err  error
	got  []network.HTTPRequest
}

func (s *stubHTTPClient) Send(_ context.Context, req network.HTTPRequest) (*network.HTTPResponse, error) {
	s.got = append(s.got, req)
	return s.resp, s.err
}

func TestDirectFetcherFunc(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	req := ghcapabilities.Request{URL: "https://Artifacts.example.com/b?sig=x", WorkflowID: "wf", MaxResponseBytes: 42}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{resp: &network.HTTPResponse{StatusCode: 200, Body: []byte("data")}}
		body, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", req)
		require.NoError(t, err)
		assert.Equal(t, []byte("data"), body)
		require.Len(t, c.got, 1)
		assert.Equal(t, http.MethodGet, c.got[0].Method)
		assert.Equal(t, req.URL, c.got[0].URL)
		assert.Equal(t, uint32(42), c.got[0].MaxResponseBytes)
		assert.Zero(t, c.got[0].Timeout, "client default timeout applies")
	})

	t.Run("zero request cap defaults to 1GB", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{resp: &network.HTTPResponse{StatusCode: 200}}
		r := req
		r.MaxResponseBytes = 0
		_, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", r)
		require.NoError(t, err)
		assert.Equal(t, uint32(utils.GB), c.got[0].MaxResponseBytes)
	})

	t.Run("request timeout is passed through", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{resp: &network.HTTPResponse{StatusCode: 200}}
		r := req
		r.TimeoutMs = 1500
		_, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", r)
		require.NoError(t, err)
		assert.Equal(t, 1500*time.Millisecond, c.got[0].Timeout)
	})

	t.Run("non-2xx returns body and error", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{resp: &network.HTTPResponse{StatusCode: 403, Body: []byte("denied")}}
		body, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", req)
		require.ErrorContains(t, err, "request failed with status code: 403")
		assert.Equal(t, []byte("denied"), body)
	})

	t.Run("client error", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{err: errors.New("boom")}
		_, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", req)
		require.ErrorContains(t, err, "boom")
	})

	t.Run("workflow ID required", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{}
		_, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", ghcapabilities.Request{URL: req.URL})
		require.ErrorContains(t, err, "workflow ID")
		assert.Empty(t, c.got)
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{}
		_, err := NewDirectFetcherFunc(c, nil, lggr)(t.Context(), "id", ghcapabilities.Request{URL: "https://%zz", WorkflowID: "wf"})
		require.ErrorContains(t, err, "invalid URL")
		assert.Empty(t, c.got)
	})

	t.Run("allowed hosts", func(t *testing.T) {
		t.Parallel()
		c := &stubHTTPClient{resp: &network.HTTPResponse{StatusCode: 200}}
		f := NewDirectFetcherFunc(c, []string{"artifacts.EXAMPLE.com"}, lggr)

		_, err := f(t.Context(), "id", req)
		require.NoError(t, err)

		other := req
		other.URL = "https://evil.example.org/b"
		_, err = f(t.Context(), "id", other)
		require.ErrorIs(t, err, ErrHostNotAllowed)
		assert.Len(t, c.got, 1)
	})
}

func TestNewDirectHTTPClient_BlocksInternalAndNonHTTPS(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	client, err := NewDirectHTTPClient(lggr, 100_000_000, noop.NewMeterProvider().Meter("test"))
	require.NoError(t, err)
	f := NewDirectFetcherFunc(client, nil, lggr)

	for name, url := range map[string]string{
		"plain http":          "http://93.184.216.34/x",
		"loopback over https": "https://127.0.0.1/x",
		"metadata endpoint":   "https://169.254.169.254/latest/meta-data",
		"non-443 port":        "https://93.184.216.34:8443/x",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := f(t.Context(), "id", ghcapabilities.Request{URL: url, WorkflowID: "wf"})
			require.Error(t, err)
			assert.ErrorIs(t, err, network.ErrBlockedRequest)
		})
	}
}
