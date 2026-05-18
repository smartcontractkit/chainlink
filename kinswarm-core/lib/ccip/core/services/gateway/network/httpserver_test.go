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

func startNewServer(t *testing.T, maxRequestBytes int64, readTimeoutMillis uint32) (server network.HttpServer, handler *mocks.HTTPRequestHandler, url string) {
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

func sendRequest(t *testing.T, url string, body []byte) *http.Response {
	req, err := http.NewRequestWithContext(testutils.Context(t), "POST", url, bytes.NewBuffer(body))
	require.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func TestHTTPServer_HandleRequest_Correct(t *testing.T) {
	server, handler, url := startNewServer(t, 100_000, 100_000)
	defer server.Close()

	handler.On("ProcessRequest", mock.Anything, mock.Anything).Return([]byte("response"), 200)

	resp := sendRequest(t, url, []byte("0123456789"))
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("response"), respBytes)
}

func TestHTTPServer_HandleRequest_RequestBodyTooBig(t *testing.T) {
	server, _, url := startNewServer(t, 5, 100_000)
	defer server.Close()

	resp := sendRequest(t, url, []byte("0123456789"))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	server, _, url := startNewServer(t, 100_000, 100_000)
	defer server.Close()

	url = strings.Replace(url, HTTPTestPath, network.HealthCheckPath, 1)
	resp := sendRequest(t, url, []byte{})
	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte(network.HealthCheckResponse), respBytes)
}
