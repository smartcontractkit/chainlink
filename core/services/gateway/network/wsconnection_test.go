package network_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var upgrader = websocket.Upgrader{}

func newWebSocketPair(t *testing.T) (*websocket.Conn, *websocket.Conn) {
	t.Helper()

	type upgradeResult struct {
		conn *websocket.Conn
		err  error
	}
	serverConnCh := make(chan upgradeResult, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		serverConnCh <- upgradeResult{conn: conn, err: err}
	}))
	t.Cleanup(server.Close)

	clientConn, resp, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.NoError(t, err)
	result := <-serverConnCh
	require.NoError(t, result.err)
	serverConn := result.conn
	t.Cleanup(func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
	})
	return serverConn, clientConn
}

type serverSideLogic struct {
	connWrapper network.WSConnectionWrapper
}

func (ssl *serverSideLogic) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// one wsConnWrapper per client
	ssl.connWrapper.Reset(c)
}

// TestWSConnectionWrapper_WriteError_TriggersClose verifies that a write error
// (e.g. i/o timeout on a stale connection) causes the connection to be closed,
// which signals closeCh and allows reconnectLoop to re-establish a fresh connection.
func TestWSConnectionWrapper_WriteError_TriggersClose(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	// server — accepts one connection
	serverConn := network.NewWSConnectionWrapper(lggr)
	servicetest.Run(t, serverConn)
	ssl := &serverSideLogic{connWrapper: serverConn}
	s := httptest.NewServer(http.HandlerFunc(ssl.wsHandler))
	serverURL := "ws" + strings.TrimPrefix(s.URL, "http")
	defer s.Close()

	// client
	clientConnWrapper := network.NewWSConnectionWrapper(lggr)
	servicetest.Run(t, clientConnWrapper)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)

	closeCh := clientConnWrapper.Reset(conn)

	// Set the write deadline to the past to simulate a write i/o timeout without
	// affecting the read side — this mimics a half-open / zombie TCP connection
	// where writes time out but the kernel hasn't detected a read error yet.
	require.NoError(t, conn.SetWriteDeadline(time.Now().Add(-time.Second)))

	// The write must fail (write deadline already expired).
	writeErr := clientConnWrapper.Write(t.Context(), websocket.BinaryMessage, []byte("data"))
	require.Error(t, writeErr, "write should fail due to expired write deadline")

	// The write failure must cause the connection to be closed so that
	// the reconnect loop can establish a fresh connection.
	select {
	case <-closeCh:
		// correct: write error triggered connection close, reconnect can proceed
	case <-time.After(5 * time.Second):
		t.Fatal("closeCh was not signaled after write error; stale connection will block reconnect indefinitely")
	}
	require.False(t, clientConnWrapper.IsConnected())
}

func TestWSConnectionWrapper_ConnectedState(t *testing.T) {
	t.Parallel()

	connWrapper := network.NewWSConnectionWrapper(logger.Test(t))
	servicetest.Run(t, connWrapper)
	require.False(t, connWrapper.IsConnected())

	serverConn, clientConn := newWebSocketPair(t)
	closeCh := connWrapper.Reset(serverConn)
	require.True(t, connWrapper.IsConnected())

	require.NoError(t, clientConn.Close())
	select {
	case <-closeCh:
	case <-time.After(5 * time.Second):
		t.Fatal("connection close was not observed")
	}
	require.False(t, connWrapper.IsConnected())
}

func TestWSConnectionWrapper_OldReadPumpCannotClearReplacement(t *testing.T) {
	t.Parallel()

	connWrapper := network.NewWSConnectionWrapper(logger.Test(t))
	servicetest.Run(t, connWrapper)

	serverConnA, _ := newWebSocketPair(t)
	closeA := connWrapper.Reset(serverConnA)
	serverConnB, clientConnB := newWebSocketPair(t)
	connWrapper.Reset(serverConnB)

	select {
	case <-closeA:
	case <-time.After(5 * time.Second):
		t.Fatal("old connection close was not observed")
	}
	require.True(t, connWrapper.IsConnected())

	payload := []byte("replacement connection")
	require.NoError(t, connWrapper.Write(t.Context(), websocket.TextMessage, payload))
	msgType, got, err := clientConnB.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, msgType)
	require.Equal(t, payload, got)
}

func TestWSConnectionWrapper_ClientReconnect(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	// server
	wsConn := network.NewWSConnectionWrapper(lggr)
	servicetest.Run(t, wsConn)
	ssl := &serverSideLogic{connWrapper: wsConn}
	s := httptest.NewServer(http.HandlerFunc(ssl.wsHandler))
	serverURL := "ws" + strings.TrimPrefix(s.URL, "http")
	defer s.Close()

	// client
	clientConnWrapper := network.NewWSConnectionWrapper(lggr)
	servicetest.Run(t, clientConnWrapper)

	// connect, write a message, disconnect
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	func() {
		defer func() { assert.NoError(t, conn.Close()) }()
		clientConnWrapper.Reset(conn)
		writeErr := clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("hello"))
		require.NoError(t, writeErr)
		<-ssl.connWrapper.ReadChannel() // consumed by server
	}()

	// try to write without a connection
	writeErr := clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("failed send"))
	require.Error(t, writeErr)

	// re-connect, write another message, disconnect
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, conn.Close()) })
	clientConnWrapper.Reset(conn)
	writeErr = clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("hello again"))
	require.NoError(t, writeErr)
	<-ssl.connWrapper.ReadChannel() // consumed by server
}
