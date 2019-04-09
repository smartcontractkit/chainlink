package synchronization_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/require"
)

func TestWebSocketStatsPusher_StartCloseStart(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	pusher := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, pusher.Start())
	cltest.CallbackOrTimeout(t, "stats pusher connects", func() {
		<-wsserver.Connected
	})
	require.NoError(t, pusher.Close())

	// restart after client disconnect
	require.NoError(t, pusher.Start())
	cltest.CallbackOrTimeout(t, "stats pusher restarts", func() {
		<-wsserver.Connected
	}, 3*time.Second)
	require.NoError(t, pusher.Close())
}

func TestWebSocketStatsPusher_ReconnectLoop(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	pusher := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, pusher.Start())
	cltest.CallbackOrTimeout(t, "stats pusher connects", func() {
		<-wsserver.Connected
	})

	// reconnect after server disconnect
	wsserver.WriteCloseMessage()
	cltest.CallbackOrTimeout(t, "stats pusher reconnects", func() {
		<-wsserver.Connected
	}, 3*time.Second)
	require.NoError(t, pusher.Close())
}

func TestWebSocketStatsPusher_Send(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	pusher := synchronization.NewWebSocketClient(wsserver.URL)
	require.NoError(t, pusher.Start())
	defer pusher.Close()

	expectation := `{"hello": "world"}`
	pusher.Send([]byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.Received)
	})
}
