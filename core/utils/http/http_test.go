package http_test

import (
	"bytes"
	"io"
	netHttp "net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/http"
)

func TestUnrestrictedHTTPClient(t *testing.T) {
	t.Parallel()

	client := http.NewUnrestrictedHTTPClient()
	assert.True(t, client.Transport.(*netHttp.Transport).DisableCompression)
	client.Transport = newMockTransport()

	netReq, err := netHttp.NewRequestWithContext(testutils.Context(t), "GET", "http://localhost", bytes.NewReader([]byte{}))
	assert.NoError(t, err)

	req := &http.HTTPRequest{
		Client:  client,
		Request: netReq,
		Config:  http.HTTPRequestConfig{SizeLimit: 1000},
		Logger:  logger.Nop(),
	}

	response, statusCode, headers, err := req.SendRequest()
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, `{"foo":123}`, string(response))
}

type mockTransport struct{}

func newMockTransport() netHttp.RoundTripper {
	return &mockTransport{}
}

func (t *mockTransport) RoundTrip(req *netHttp.Request) (*netHttp.Response, error) {
	// Create mocked http.Response
	response := &netHttp.Response{
		Header:     make(netHttp.Header),
		Request:    req,
		StatusCode: netHttp.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")

	responseBody := `{"foo":123}`
	response.Body = io.NopCloser(strings.NewReader(responseBody))
	return response, nil
}
