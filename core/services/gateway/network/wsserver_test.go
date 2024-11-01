package network_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

const (
	WSTestHost = "localhost"
	WSTestPath = "/ws_test_path"
)

func startNewWSServer(t *testing.T, readTimeoutMillis uint32) (server network.WebSocketServer, acceptor *mocks.ConnectionAcceptor, url string) {
	config := &network.WebSocketServerConfig{
		HTTPServerConfig: network.HTTPServerConfig{
			Host:                 WSTestHost,
			Port:                 0,
			Path:                 "/ws_test_path",
			TLSEnabled:           false,
			ContentTypeHeader:    "application/jsonrpc",
			ReadTimeoutMillis:    readTimeoutMillis,
			WriteTimeoutMillis:   10_000,
			RequestTimeoutMillis: 10_000,
			MaxRequestBytes:      100_000,
		},
		HandshakeTimeoutMillis: 10_000,
	}

	acceptor = mocks.NewConnectionAcceptor(t)
	server = network.NewWebSocketServer(config, acceptor, logger.TestLogger(t))
	err := server.Start(testutils.Context(t))
	require.NoError(t, err)

	port := server.GetPort()
	url = fmt.Sprintf("http://%s:%d%s", WSTestHost, port, WSTestPath)
	return
}

func sendRequestWithHeader(t *testing.T, url string, headerName string, headerValue string) *http.Response {
	req, err := http.NewRequestWithContext(testutils.Context(t), "POST", url, bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	req.Header.Set(headerName, headerValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func TestWSServer_HandleRequest_AuthHeaderTooBig(t *testing.T) {
	t.Parallel()
	server, _, url := startNewWSServer(t, 100_000)
	defer server.Close()

	longString := "abcdefgh"
	for i := 0; i < 6; i++ {
		longString += longString
	}
	authHeader := base64.StdEncoding.EncodeToString([]byte(longString))
	resp := sendRequestWithHeader(t, url, network.WsServerHandshakeAuthHeaderName, authHeader)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestWSServer_HandleRequest_AuthHeaderIncorrectlyBase64Encoded(t *testing.T) {
	t.Parallel()
	server, _, url := startNewWSServer(t, 100_000)
	defer server.Close()

	resp := sendRequestWithHeader(t, url, network.WsServerHandshakeAuthHeaderName, "}}}")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestWSServer_HandleRequest_AuthHeaderInvalid(t *testing.T) {
	t.Parallel()
	server, acceptor, url := startNewWSServer(t, 100_000)
	defer server.Close()

	acceptor.On("StartHandshake", mock.Anything).Return("", []byte{}, errors.New("invalid auth header"))

	authHeader := base64.StdEncoding.EncodeToString([]byte("abcd"))
	resp := sendRequestWithHeader(t, url, network.WsServerHandshakeAuthHeaderName, authHeader)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
