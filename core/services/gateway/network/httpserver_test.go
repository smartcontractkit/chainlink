package network_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

const (
	HTTPTestHost = "localhost"
	HTTPTestPath = "/test_path"
)

func startNewServer(t *testing.T, maxRequestBytes int64, readTimeoutMillis uint32, enabledCORS bool, allowedOrigins []string) (server network.HTTPServer, handler *mocks.HTTPRequestHandler, url string) {
	t.Helper()
	config := &network.HTTPServerConfig{
		Host:                 HTTPTestHost,
		Port:                 0,
		Path:                 HTTPTestPath,
		TLSEnabled:           false,
		ContentTypeHeader:    "application/jsonrpc",
		ReadTimeoutMillis:    readTimeoutMillis,
		WriteTimeoutMillis:   10_000,
		RequestTimeoutMillis: 10_000,
		MaxRequestBytes:      maxRequestBytes,
		CORSEnabled:          enabledCORS,
		CORSAllowedOrigins:   allowedOrigins,
	}

	handler = mocks.NewHTTPRequestHandler(t)
	lggr := logger.Test(t)
	server, err := network.NewHTTPServer(config, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	server.SetHTTPRequestHandler(handler)
	servicetest.Run(t, server)

	port := server.GetPort()
	url = fmt.Sprintf("http://%s:%d%s", HTTPTestHost, port, HTTPTestPath)
	return
}

func sendRequest(t *testing.T, url string, body []byte, httpMethod string, origin *string) (*http.Response, []byte) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), httpMethod, url, bytes.NewBuffer(body))
	if origin != nil {
		req.Header.Set("Origin", *origin)
	}
	require.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, respBytes
}

func TestHTTPServer_HandleRequest_Correct(t *testing.T) {
	t.Parallel()
	_, handler, url := startNewServer(t, 100_000, 100_000, false, nil)

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodPost, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
}

func TestHTTPServer_HandleRequest_RequestBodyTooBig(t *testing.T) {
	t.Parallel()
	_, _, url := startNewServer(t, 5, 100_000, false, nil)

	resp, _ := sendRequest(t, url, []byte("0123456789"), http.MethodPost, nil)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	t.Parallel()
	_, _, url := startNewServer(t, 100_000, 100_000, false, nil)

	url = strings.Replace(url, HTTPTestPath, network.HealthCheckPath, 1)
	resp, respBytes := sendRequest(t, url, []byte{}, http.MethodPost, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte(network.HealthCheckResponse), respBytes)
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOrigin(t *testing.T) {
	t.Parallel()
	_, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://remix.ethereum.org"
	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOriginWildcards(t *testing.T) {
	t.Parallel()
	_, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://*.ethereum.org", "https://*.valid.domain.com", "http://*.gov"})

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://remix.ethereum.org"
	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "https://another.valid.domain.com"
	resp, respBytes = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://example.gov"
	resp, respBytes = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOrigin_PreflightRequest(t *testing.T) {
	t.Parallel()
	_, _, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})

	origin := "https://remix.ethereum.org"
	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodOptions, &origin)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromNotAllowedOrigin(t *testing.T) {
	t.Parallel()
	_, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://not.allowed.origin.com"
	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromNotAllowedOriginWildcards(t *testing.T) {
	t.Parallel()
	_, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://*.ethereum.org", "https://*.valid.domain.com", "http://example.gov:8080"})

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://ethereum.remix.org" // doesn't end with ethereum.org
	resp, respBytes := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://another.valid.domain.org" // http instead of https
	resp, respBytes = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://example.gov" // port missing
	resp, respBytes = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
}
