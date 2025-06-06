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

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

const (
	HTTPTestHost = "localhost"
	HTTPTestPath = "/test_path"
)

func startNewServer(t *testing.T, maxRequestBytes int64, readTimeoutMillis uint32, enabledCORS bool, allowedOrigins []string) (server network.HttpServer, handler *mocks.HTTPRequestHandler, url string) {
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
	server = network.NewHttpServer(config, logger.TestLogger(t))
	server.SetHTTPRequestHandler(handler)
	err := server.Start(testutils.Context(t))
	require.NoError(t, err)

	port := server.GetPort()
	url = fmt.Sprintf("http://%s:%d%s", HTTPTestHost, port, HTTPTestPath)
	return
}

func sendRequest(t *testing.T, url string, body []byte, httpMethod string, origin *string) *http.Response {
	req, err := http.NewRequestWithContext(testutils.Context(t), httpMethod, url, bytes.NewBuffer(body))
	if origin != nil {
		req.Header.Set("Origin", *origin)
	}
	require.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func TestHTTPServer_HandleRequest_Correct(t *testing.T) {
	t.Parallel()
	server, handler, url := startNewServer(t, 100_000, 100_000, false, nil)
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, nil)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
}

func TestHTTPServer_HandleRequest_RequestBodyTooBig(t *testing.T) {
	t.Parallel()
	server, _, url := startNewServer(t, 5, 100_000, false, nil)
	defer server.Close()

	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, nil)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	t.Parallel()
	server, _, url := startNewServer(t, 100_000, 100_000, false, nil)
	defer server.Close()

	url = strings.Replace(url, HTTPTestPath, network.HealthCheckPath, 1)
	resp := sendRequest(t, url, []byte{}, http.MethodPost, nil)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte(network.HealthCheckResponse), respBytes)
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOrigin(t *testing.T) {
	t.Parallel()
	server, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://remix.ethereum.org"
	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOriginWildcards(t *testing.T) {
	t.Parallel()
	server, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://*.ethereum.org", "https://*.valid.domain.com", "http://*.gov"})
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://remix.ethereum.org"
	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "https://another.valid.domain.com"
	resp = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://example.gov"
	resp = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromAllowedOrigin_PreflightRequest(t *testing.T) {
	t.Parallel()
	server, _, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})
	defer server.Close()

	origin := "https://remix.ethereum.org"
	resp := sendRequest(t, url, []byte("0123456789"), http.MethodOptions, &origin)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, respBytes)
	require.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromNotAllowedOrigin(t *testing.T) {
	t.Parallel()
	server, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://remix.ethereum.org", "https://another.valid.origin.com"})
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://not.allowed.origin.com"
	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestHTTPServer_HandleRequest_CORSEnabled_FromNotAllowedOriginWildcards(t *testing.T) {
	t.Parallel()
	server, handler, url := startNewServer(t, 100_000, 100_000, true,
		[]string{"https://*.ethereum.org", "https://*.valid.domain.com", "http://example.gov:8080"})
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin := "https://ethereum.remix.org" // doesn't end with ethereum.org
	resp := sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://another.valid.domain.org" // http instead of https
	resp = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Headers"))

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	origin = "http://example.gov" // port missing
	resp = sendRequest(t, url, []byte("0123456789"), http.MethodPost, &origin)
	respBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Methods"))
	require.Equal(t, "", resp.Header.Get("Access-Control-Allow-Headers"))
}
