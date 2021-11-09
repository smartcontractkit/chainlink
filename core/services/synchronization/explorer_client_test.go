package synchronization_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocketClient_ReconnectLoop(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	}, 1*time.Second)

	// reconnect after server disconnect
	wsserver.WriteCloseMessage()
	cltest.CallbackOrTimeout(t, "ws client reconnects", func() {
		<-wsserver.Connected
	}, 3*time.Second)
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
	explorerClient := synchronization.NewExplorerClient(url, "accessKey", "secret")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

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

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

	expectation := `{"hello": "world"}`
	explorerClient.Send(context.Background(), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	}, 1*time.Second)
}

func TestWebSocketClient_Send_TextMessage(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

	expectation := `{"hello": "world"}`
	explorerClient.Send(context.Background(), []byte(expectation), synchronization.ExplorerTextMessage)
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	})
}

func TestWebSocketClient_Send_Binary(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

	address := common.HexToAddress("0xabc123")
	addressBytes := address.Bytes()
	explorerClient.Send(context.Background(), addressBytes, synchronization.ExplorerBinaryMessage)
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, addressBytes, <-wsserver.ReceivedBinary)
	})
}

func TestWebSocketClient_Send_Unsupported(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())

	assert.PanicsWithValue(t, "send on explorer client received unsupported message type -1", func() {
		explorerClient.Send(context.Background(), []byte(`{"hello": "world"}`), -1)
	})
	require.NoError(t, explorerClient.Close())
}

func TestWebSocketClient_Send_WithAck(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

	expectation := `{"hello": "world"}`
	explorerClient.Send(context.Background(), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
		err := wsserver.Broadcast(`{"result": 200}`)
		assert.NoError(t, err)
	})

	cltest.CallbackOrTimeout(t, "receive response", func() {
		response, err := explorerClient.Receive(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestWebSocketClient_Send_WithAckTimeout(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	require.NoError(t, explorerClient.Start())
	defer explorerClient.Close()

	expectation := `{"hello": "world"}`
	explorerClient.Send(context.Background(), []byte(expectation))
	cltest.CallbackOrTimeout(t, "receive stats", func() {
		require.Equal(t, expectation, <-wsserver.ReceivedText)
	})

	cltest.CallbackOrTimeout(t, "receive response", func() {
		_, err := explorerClient.Receive(context.Background(), 100*time.Millisecond)
		assert.Error(t, err)
		assert.Equal(t, err, synchronization.ErrReceiveTimeout)
	}, 300*time.Millisecond)
}

func TestWebSocketClient_Status_ConnectAndServerDisconnect(t *testing.T) {
	wsserver, cleanup := cltest.NewEventWebSocketServer(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	defer explorerClient.Close()
	assert.Equal(t, synchronization.ConnectionStatusDisconnected, explorerClient.Status())

	require.NoError(t, explorerClient.Start())
	cltest.CallbackOrTimeout(t, "ws client connects", func() {
		<-wsserver.Connected
	})
	cltest.NewGomegaWithT(t).Eventually(func() synchronization.ConnectionStatus {
		return explorerClient.Status()
	}).Should(gomega.Equal(synchronization.ConnectionStatusConnected))

	wsserver.WriteCloseMessage()
	wsserver.Close()

	time.Sleep(synchronization.CloseTimeout + (100 * time.Millisecond))

	cltest.NewGomegaWithT(t).Eventually(func() synchronization.ConnectionStatus {
		return explorerClient.Status()
	}).Should(gomega.Equal(synchronization.ConnectionStatusError))
}

func TestWebSocketClient_Status_ConnectError(t *testing.T) {
	badURL, err := url.Parse("http://badhost.com")
	require.NoError(t, err)

	errorexplorerClient := synchronization.NewExplorerClient(badURL, "", "")
	require.NoError(t, errorexplorerClient.Start())
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, synchronization.ConnectionStatusError, errorexplorerClient.Status())

}
