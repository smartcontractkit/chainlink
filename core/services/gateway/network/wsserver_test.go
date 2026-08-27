package network_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	WSTestHost = "localhost"
	WSTestPath = "/ws_test_path"
)

func startNewWSServer(t *testing.T, readTimeoutMillis uint32) (server network.WebSocketServer, acceptor *mocks.ConnectionAcceptor, url string) {
	return startNewWSServerWithPath(t, readTimeoutMillis, WSTestPath)
}

func startNewWSServerWithPath(t *testing.T, readTimeoutMillis uint32, path string) (server network.WebSocketServer, acceptor *mocks.ConnectionAcceptor, url string) {
	config := &network.WebSocketServerConfig{
		HTTPServerConfig: network.HTTPServerConfig{
			Host:                 WSTestHost,
			Port:                 0,
			Path:                 path,
			TLSEnabled:           false,
			ContentTypeHeader:    "application/jsonrpc",
			ReadTimeoutMillis:    readTimeoutMillis,
			WriteTimeoutMillis:   10_000,
			RequestTimeoutMillis: 10_000,
			MaxRequestBytes:      10_000,
		},
		HandshakeTimeoutMillis: 10_000,
	}

	acceptor = mocks.NewConnectionAcceptor(t)
	lggr := logger.Test(t)
	server, err := network.NewWebSocketServer(config, acceptor, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	servicetest.Run(t, server)

	port := server.GetPort()
	url = fmt.Sprintf("http://%s:%d%s", WSTestHost, port, path)
	return server, acceptor, url
}

func sendRequestWithHeader(t *testing.T, url, headerName, headerValue string) *http.Response {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	req.Header.Set(headerName, headerValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func TestWSServer_HandleRequest_AuthHeaderTooBig(t *testing.T) {
	t.Parallel()
	_, _, urlStr := startNewWSServer(t, 100_000)

	authHeader := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte("abcdefgh"), 64))
	resp := sendRequestWithHeader(t, urlStr, network.WsServerHandshakeAuthHeaderName, authHeader)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestWSServer_HandleRequest_AuthHeaderIncorrectlyBase64Encoded(t *testing.T) {
	t.Parallel()
	_, _, urlStr := startNewWSServer(t, 100_000)

	resp := sendRequestWithHeader(t, urlStr, network.WsServerHandshakeAuthHeaderName, "}}}")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestWSServer_HandleRequest_AuthHeaderInvalid(t *testing.T) {
	t.Parallel()
	_, acceptor, urlStr := startNewWSServer(t, 100_000)

	acceptor.On("StartHandshake", mock.Anything).Return("", []byte{}, errors.New("invalid auth header"))

	authHeader := base64.StdEncoding.EncodeToString([]byte("abcd"))
	resp := sendRequestWithHeader(t, urlStr, network.WsServerHandshakeAuthHeaderName, authHeader)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWSServer_RejectsHealthCheckRequestPath(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	config := &network.WebSocketServerConfig{HTTPServerConfig: network.HTTPServerConfig{Path: network.HealthCheckPath}}
	_, err := network.NewWebSocketServer(config, mocks.NewConnectionAcceptor(t), lggr, limits.Factory{Logger: lggr})
	require.EqualError(t, err, `WebSocket request path "/health" conflicts with health check path`)
}

func TestWSServer_HealthEndpointBypassesAuthentication(t *testing.T) {
	t.Parallel()

	_, _, baseURL := startNewWSServerWithPath(t, 100_000, "/")
	baseURL = strings.TrimSuffix(baseURL, "/")
	resp := sendRequestWithHeader(t, baseURL+network.HealthCheckPath, "unused", "unused")
	body, readErr := io.ReadAll(resp.Body)
	require.NoError(t, resp.Body.Close())
	require.NoError(t, readErr)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte(network.HealthCheckResponse), body)
}

func TestWSServer_UnauthenticatedConnectorPathStillRejected(t *testing.T) {
	t.Parallel()

	_, acceptor, urlStr := startNewWSServerWithPath(t, 100_000, "/")
	acceptor.On("StartHandshake", mock.Anything).Return("", []byte{}, errors.New("invalid auth header")).Once()
	resp := sendRequestWithHeader(t, urlStr, "unused", "unused")
	require.NoError(t, resp.Body.Close())
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWSServer_WSClient_DefaultConfig_Success(t *testing.T) {
	t.Parallel()
	_, acceptor, urlStr := startNewWSServerWithPath(t, 10_000, "/")

	waitCh := make(chan struct{})
	acceptor.On("StartHandshake", mock.Anything).Return("", []byte("challenge"), nil)
	acceptor.On("FinalizeHandshake", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		close(waitCh)
	})

	initiator := mocks.NewConnectionInitiator(t)
	initiator.On("NewAuthHeader", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return([]byte{}, nil)
	initiator.On("ChallengeResponse", mock.AnythingOfType("*context.cancelCtx"), mock.Anything, mock.Anything).Return([]byte{}, nil)

	client := network.NewWebSocketClient(network.WebSocketClientConfig{}, initiator, logger.Test(t))

	urlStr = strings.Replace(urlStr, "http", "ws", 1)
	parsedURL, err := url.Parse(urlStr)
	require.NoError(t, err)
	conn, err := client.Connect(t.Context(), parsedURL)
	require.NoError(t, err)
	require.NotNil(t, conn)

	<-waitCh
	require.NoError(t, conn.Close())
}

func TestWSServer_WSClient_DefaultConfig_Failure(t *testing.T) {
	t.Parallel()
	_, acceptor, urlStr := startNewWSServer(t, 10_000)

	waitCh := make(chan struct{})
	acceptor.On("StartHandshake", mock.Anything).Return("", []byte("challenge"), nil)
	acceptor.On("AbortHandshake", mock.Anything).Run(func(args mock.Arguments) {
		close(waitCh)
	})

	initiator := mocks.NewConnectionInitiator(t)
	initiator.On("NewAuthHeader", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return([]byte{}, nil)
	resp := make([]byte, 20000)
	initiator.On("ChallengeResponse", mock.AnythingOfType("*context.cancelCtx"), mock.Anything, mock.Anything).Return(resp, nil)

	client := network.NewWebSocketClient(network.WebSocketClientConfig{}, initiator, logger.Test(t))

	urlStr = strings.Replace(urlStr, "http", "ws", 1)
	parsedURL, err := url.Parse(urlStr)
	require.NoError(t, err)
	conn, err := client.Connect(t.Context(), parsedURL)
	require.NoError(t, err)
	require.NotNil(t, conn)

	<-waitCh
}
