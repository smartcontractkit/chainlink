package synchronization_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/stretchr/testify/require"
)

func TestWebSocketClient_StartCloseStart(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	})
	require.NoError(t, wsclient.Close())

	// restart after client disconnect
	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client restarts", func() {
		<-wsserver.Connected
	}, 3*time.Second)
	require.NoError(t, wsclient.Close())
}

func TestWebSocketClient_ReconnectLoop(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	})

	// reconnect after server disconnect
	wsserver.WriteCloseMessage()
	cltest.CallbackOrTimeout(t, "ws client reconnects", func() {
		<-wsserver.Connected
	}, 3*time.Second)
	require.NoError(t, wsclient.Close())
}

func TestWebSocketClient_Send(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, wsclient.Start())
	defer wsclient.Close()

	expectation := `{"hello": "world"}`
	wsclient.Send([]byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.Received)
	})
}
