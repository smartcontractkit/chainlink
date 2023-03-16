package synchronization_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/static"
)

func TestWebSocketClient_ReconnectLoop(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	}, testutils.WaitTimeout(t))

	// reconnect after server disconnect
	wsserver.WriteCloseMessage()
	cltest.CallbackOrTimeout(t, "ws client reconnects", func() {
		<-wsserver.Disconnected
		<-wsserver.Connected
	}, testutils.WaitTimeout(t))
	require.NoError(t, explorerClient.Close())
}

func TestWebSocketClient_Authentication(t *testing.T) {
	headerChannel := make(chan http.Header, 1)
	handler := func(w http.ResponseWriter, r *http.Request) {
		headerChannel <- r.Header
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := cltest.MustParseURL(t, server.URL)
	url.Scheme = "ws"
	explorerClient := synchronization.NewExplorerClient(url, "accessKey", "secret", logger.TestLogger(t))
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	cltest.CallbackOrTimeout(t, "receive authentication headers", func() {
		headers := <-headerChannel
		assert.Equal(t, []string{"accessKey"}, headers["X-Explore-Chainlink-Accesskey"])
		assert.Equal(t, []string{"secret"}, headers["X-Explore-Chainlink-Secret"])
		assert.Equal(t, []string{static.Version}, headers["X-Explore-Chainlink-Core-Version"])
		assert.Equal(t, []string{static.Sha}, headers["X-Explore-Chainlink-Core-Sha"])
	})
}

func TestWebSocketClient_Send_DefaultsToTextMessage(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	expectation := `{"hello": "world"}`
	explorerClient.Send(testutils.Context(t), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	}, testutils.WaitTimeout(t))
}

func TestWebSocketClient_Send_TextMessage(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	expectation := `{"hello": "world"}`
	explorerClient.Send(testutils.Context(t), []byte(expectation), synchronization.ExplorerTextMessage)
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	}, testutils.WaitTimeout(t))
}

func TestWebSocketClient_Send_Binary(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	address := common.HexToAddress("0xabc123")
	addressBytes := address.Bytes()
	explorerClient.Send(testutils.Context(t), addressBytes, synchronization.ExplorerBinaryMessage)
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, addressBytes, <-wsserver.ReceivedBinary)
	}, testutils.WaitTimeout(t))
}

func TestWebSocketClient_Send_Unsupported(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))

	assert.PanicsWithValue(t, "send on explorer client received unsupported message type -1", func() {
		explorerClient.Send(testutils.Context(t), []byte(`{"hello": "world"}`), -1)
	})
	require.NoError(t, explorerClient.Close())
}

func TestWebSocketClient_Send_WithAck(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	expectation := `{"hello": "world"}`
	explorerClient.Send(testutils.Context(t), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
		err := wsserver.Broadcast(`{"result": 200}`)
		assert.NoError(t, err)
	}, testutils.WaitTimeout(t))

	cltest.CallbackOrTimeout(t, "receive response", func() {
		response, err := explorerClient.Receive(testutils.Context(t))
		assert.NoError(t, err)
		assert.NotNil(t, response)
	}, testutils.WaitTimeout(t))
}

func TestWebSocketClient_Send_WithAckTimeout(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()

	expectation := `{"hello": "world"}`
	explorerClient.Send(testutils.Context(t), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	}, testutils.WaitTimeout(t))

	cltest.CallbackOrTimeout(t, "receive response", func() {
		_, err := explorerClient.Receive(testutils.Context(t), 100*time.Millisecond)
		assert.ErrorIs(t, err, synchronization.ErrReceiveTimeout)
	}, testutils.WaitTimeout(t))
}

func TestWebSocketClient_Status_ConnectAndServerDisconnect(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := newTestExplorerClient(t, wsserver.URL)
	assert.Equal(t, synchronization.ConnectionStatusDisconnected, explorerClient.Status())

	require.NoError(t, explorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, explorerClient.Close()) }()
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	}, testutils.WaitTimeout(t))

	gomega.NewWithT(t).Eventually(func() synchronization.ConnectionStatus {
		return explorerClient.Status()
	}).Should(gomega.Equal(synchronization.ConnectionStatusConnected))

	// this triggers ConnectionStatusError and then the client gets reconnected
	wsserver.WriteCloseMessage()

	cltest.CallbackOrTimeout(t, "ws client disconnects and reconnects", func() {
		<-wsserver.Disconnected
		<-wsserver.Connected
	}, testutils.WaitTimeout(t))

	// expecting the client to reconnect
	gomega.NewWithT(t).Eventually(func() synchronization.ConnectionStatus {
		return explorerClient.Status()
	}).Should(gomega.Equal(synchronization.ConnectionStatusConnected))

	require.Equal(t, 1, wsserver.ConnectionsCount())
}

func TestWebSocketClient_Status_ConnectError(t *testing.T) {
	badURL, err := url.Parse("http://badhost.com")
	require.NoError(t, err)

	errorExplorerClient := newTestExplorerClient(t, badURL)
	require.NoError(t, errorExplorerClient.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, errorExplorerClient.Close()) }()

	gomega.NewWithT(t).Eventually(func() synchronization.ConnectionStatus {
		return errorExplorerClient.Status()
	}).Should(gomega.Equal(synchronization.ConnectionStatusError))
}

func newTestExplorerClient(t *testing.T, wsURL *url.URL) synchronization.ExplorerClient {
	return synchronization.NewExplorerClient(wsURL, "", "", logger.TestLogger(t))
}
