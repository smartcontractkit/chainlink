package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRelay_OpenConnectBridgeAndClose(t *testing.T) {
	srv := NewServer(zerolog.Nop(), nil)
	httpServer := httptest.NewServer(srv.Handler())
	defer httpServer.Close()

	openResp := mustOpenRelay(t, httpServer.URL, openRelayRequest{
		Name:          "relay-critical-path",
		RequestedPort: 0,
	})
	require.NotEmpty(t, openResp.RelayID)
	require.Positive(t, openResp.BoundPort)

	wsConn := mustConnectRelayWS(t, httpServer.URL, openResp.RelayID)
	defer wsConn.Close()

	dialer := net.Dialer{}
	tcpConn, err := dialer.DialContext(context.Background(), "tcp", fmt.Sprintf("127.0.0.1:%d", openResp.BoundPort))
	require.NoError(t, err, "tcp client should connect to opened relay port")
	defer tcpConn.Close()

	_ = tcpConn.SetDeadline(time.Now().Add(3 * time.Second))
	_ = wsConn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_ = wsConn.SetWriteDeadline(time.Now().Add(3 * time.Second))

	_, err = tcpConn.Write([]byte("hello-from-tcp"))
	require.NoError(t, err, "writing to relay tcp side should succeed")

	msgType, payload, err := wsConn.ReadMessage()
	require.NoError(t, err, "relay should forward tcp payload to websocket")
	require.Equal(t, websocket.BinaryMessage, msgType)
	require.Equal(t, "hello-from-tcp", string(payload))

	err = wsConn.WriteMessage(websocket.BinaryMessage, []byte("hello-from-ws"))
	require.NoError(t, err, "writing to relay websocket side should succeed")

	buf := make([]byte, 64)
	n, err := tcpConn.Read(buf)
	require.NoError(t, err, "relay should forward websocket payload to tcp")
	require.Equal(t, "hello-from-ws", string(buf[:n]))

	closeResult := mustCloseRelay(t, httpServer.URL, openResp.RelayID)
	require.Equal(t, openResp.RelayID, closeResult["relayId"])
	require.Equal(t, true, closeResult["closed"])
	require.Equal(t, true, closeResult["found"])
}

func TestRelay_OpenIdempotentByRequestedPort(t *testing.T) {
	srv := NewServer(zerolog.Nop(), nil)
	httpServer := httptest.NewServer(srv.Handler())
	defer httpServer.Close()

	requestedPort := reserveFreePort(t)

	first := mustOpenRelay(t, httpServer.URL, openRelayRequest{
		Name:          "relay-first",
		RequestedPort: requestedPort,
	})
	second := mustOpenRelay(t, httpServer.URL, openRelayRequest{
		Name:          "relay-second",
		RequestedPort: requestedPort,
	})

	require.Equal(t, first.RelayID, second.RelayID, "same requested port should reuse existing relay")
	require.Equal(t, first.BoundPort, second.BoundPort)

	closeResult := mustCloseRelay(t, httpServer.URL, first.RelayID)
	require.Equal(t, true, closeResult["closed"])
	require.Equal(t, true, closeResult["found"])
}

func TestRelay_ConnectMissingRelayIDFails(t *testing.T) {
	srv := NewServer(zerolog.Nop(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/relay/connect", nil)
	rr := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), ErrCodeMissingComponentInput)
}

func mustOpenRelay(t *testing.T, baseURL string, req openRelayRequest) openRelayResponse {
	t.Helper()
	body, err := json.Marshal(req)
	require.NoError(t, err)
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, baseURL+"/v1/relay/open", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out openRelayResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return out
}

func mustCloseRelay(t *testing.T, baseURL, relayID string) map[string]any {
	t.Helper()
	body, err := json.Marshal(closeRelayRequest{RelayID: relayID})
	require.NoError(t, err)
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, baseURL+"/v1/relay/close", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

func mustConnectRelayWS(t *testing.T, baseURL, relayID string) *websocket.Conn {
	t.Helper()
	parsed, err := url.Parse(baseURL)
	require.NoError(t, err)
	wsURL := fmt.Sprintf("ws://%s/v1/relay/connect?relayId=%s", parsed.Host, relayID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "websocket bridge should connect")
	return conn
}

func reserveFreePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	addr, ok := ln.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return addr.Port
}
