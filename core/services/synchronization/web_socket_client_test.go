package synchronization_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocketClient_StartCloseStart(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	}, 1*time.Second)
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

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	}, 1*time.Second)

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

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	require.NoError(t, wsclient.Start())
	defer wsclient.Close()

	expectation := `{"hello": "world"}`
	wsclient.Send([]byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.Received)
	})
}

func TestWebSocketClient_Authentication(t *testing.T) {
	headerChannel := make(chan http.Header)
	handler := func(w http.ResponseWriter, r *http.Request) {
		headerChannel <- r.Header
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := cltest.MustParseURL(server.URL)
	url.Scheme = "ws"
	wsclient := synchronization.NewWebSocketClient(url, "accessKey", "secret")
	require.NoError(t, wsclient.Start())
	defer wsclient.Close()

	cltest.CallbackOrTimeout(t, "receive authentication headers", func() {
		headers := <-headerChannel
		assert.Equal(t, []string{"accessKey"}, headers["X-Explore-Chainlink-Accesskey"])
		assert.Equal(t, []string{"secret"}, headers["X-Explore-Chainlink-Secret"])
		assert.Equal(t, []string{store.Version}, headers["X-Explore-Chainlink-Core-Version"])
		assert.Equal(t, []string{store.Sha}, headers["X-Explore-Chainlink-Core-Sha"])
	})
}

func TestWebSocketClient_SendWithAck(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	require.NoError(t, wsclient.Start())
	defer wsclient.Close()

	expectation := `{"hello": "world"}`
	wsclient.Send([]byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.Received)
		err := wsserver.Broadcast(`{"result": 200}`)
		assert.NoError(t, err)
	})

	cltest.CallbackOrTimeout(t, "receive response", func() {
		response, err := wsclient.Receive()
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestWebSocketClient_SendWithAckTimeout(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	require.NoError(t, wsclient.Start())
	defer wsclient.Close()

	expectation := `{"hello": "world"}`
	wsclient.Send([]byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.Received)
	})

	cltest.CallbackOrTimeout(t, "receive response", func() {
		_, err := wsclient.Receive(100 * time.Millisecond)
		assert.Error(t, err)
		assert.Equal(t, err, synchronization.ErrReceiveTimeout)
	}, 300*time.Millisecond)
}

func TestWebSocketClient_Status_ConnectAndServerDisconnect(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	wsclient := synchronization.NewWebSocketClient(wsserver.URL, "", "")
	defer wsclient.Close()
	assert.Equal(t, synchronization.ConnectionStatusDisconnected, wsclient.Status())

	require.NoError(t, wsclient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	})
	assert.Equal(t, synchronization.ConnectionStatusConnected, wsclient.Status())

	wsserver.WriteCloseMessage()
	wsserver.Close()

	// Must be longer than synchronization.closeTimeout
	time.Sleep(1100 * time.Millisecond)

	assert.Equal(t, synchronization.ConnectionStatusError, wsclient.Status())

}

func TestWebSocketClient_Status_ConnectError(t *testing.T) {
	badURL, err := url.Parse("http://badhost.com")
	require.NoError(t, err)

	errorwsclient := synchronization.NewWebSocketClient(badURL, "", "")
	require.NoError(t, errorwsclient.Start())
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, synchronization.ConnectionStatusError, errorwsclient.Status())

}
