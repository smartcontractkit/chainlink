package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/require"
)

func TestWebsocketStatsPusher_New(t *testing.T) {
	wsserver, cleanup := cltest.NewCountingWebsocketServer(t)
	defer cleanup()

	pusher := store.NewWebsocketStatsPusher(wsserver.URL)
	require.NoError(t, pusher.Start())
	cltest.CallbackOrTimeout(t, "stats pusher connects", func() {
		<-wsserver.Connected
	})
	require.NoError(t, pusher.Close())
}
