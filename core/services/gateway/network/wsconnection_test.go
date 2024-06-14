package network_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var upgrader = websocket.Upgrader{}

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

func TestWSConnectionWrapper_ClientReconnect(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	// server
	ssl := &serverSideLogic{connWrapper: network.NewWSConnectionWrapper(lggr)}
	require.NoError(t, ssl.connWrapper.Start(ctx))
	s := httptest.NewServer(http.HandlerFunc(ssl.wsHandler))
	serverURL := "ws" + strings.TrimPrefix(s.URL, "http")
	defer s.Close()

	// client
	clientConnWrapper := network.NewWSConnectionWrapper(lggr)
	require.NoError(t, clientConnWrapper.Start(ctx))

	// connect, write a message, disconnect
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	clientConnWrapper.Reset(conn)
	writeErr := clientConnWrapper.Write(testutils.Context(t), websocket.TextMessage, []byte("hello"))
	require.NoError(t, writeErr)
	<-ssl.connWrapper.ReadChannel() // consumed by server
	conn.Close()

	// try to write without a connection
	writeErr = clientConnWrapper.Write(testutils.Context(t), websocket.TextMessage, []byte("failed send"))
	require.Error(t, writeErr)

	// re-connect, write another message, disconnect
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	clientConnWrapper.Reset(conn)
	writeErr = clientConnWrapper.Write(testutils.Context(t), websocket.TextMessage, []byte("hello again"))
	require.NoError(t, writeErr)
	<-ssl.connWrapper.ReadChannel() // consumed by server
	conn.Close()

	ssl.connWrapper.Close()
	clientConnWrapper.Close()
	clientConnWrapper.Close() // safe to call Close() twice
}
