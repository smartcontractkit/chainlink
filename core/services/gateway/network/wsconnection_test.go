package network_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

	conn, resp, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

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
	conn, resp, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
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
	conn, resp, err = websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	t.Cleanup(func() { assert.NoError(t, conn.Close()) })
	clientConnWrapper.Reset(conn)
	writeErr = clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("hello again"))
	require.NoError(t, writeErr)
	<-ssl.connWrapper.ReadChannel() // consumed by server
}

// recordingObserver is a WSConnectionObserver that captures every observation
// for assertions.
type recordingObserver struct {
	mu            sync.Mutex
	queueWaits    []time.Duration
	socketWrites  []time.Duration
	dispatchWaits []time.Duration
	pending       int64
}

func (o *recordingObserver) RecordWriteQueueWait(_ context.Context, d time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.queueWaits = append(o.queueWaits, d)
}

func (o *recordingObserver) RecordSocketWrite(_ context.Context, d time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.socketWrites = append(o.socketWrites, d)
}

func (o *recordingObserver) RecordReadDispatchWait(_ context.Context, d time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.dispatchWaits = append(o.dispatchWaits, d)
}

func (o *recordingObserver) AddPendingWriters(_ context.Context, delta int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.pending += delta
}

func (o *recordingObserver) queueWaitSamples() []time.Duration {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]time.Duration(nil), o.queueWaits...)
}

func (o *recordingObserver) socketWriteSamples() []time.Duration {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]time.Duration(nil), o.socketWrites...)
}

func (o *recordingObserver) dispatchWaitSamples() []time.Duration {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]time.Duration(nil), o.dispatchWaits...)
}

func (o *recordingObserver) pendingWriters() int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.pending
}

// connGate blocks writes on demand, to hold the write pump inside WriteMessage.
type connGate struct {
	mu      sync.RWMutex
	blocker chan struct{} // when non-nil, Write blocks until it is closed
}

func (g *connGate) block() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.blocker == nil {
		g.blocker = make(chan struct{})
	}
}

func (g *connGate) unblock() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.blocker != nil {
		close(g.blocker)
		g.blocker = nil
	}
}

type gatedConn struct {
	net.Conn
	gate *connGate
}

func (c *gatedConn) Write(b []byte) (int, error) {
	c.gate.mu.RLock()
	blocker := c.gate.blocker
	c.gate.mu.RUnlock()
	if blocker != nil {
		<-blocker
	}
	return c.Conn.Write(b)
}

type gatedListener struct {
	net.Listener
	gate *connGate
}

func (l *gatedListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &gatedConn{Conn: c, gate: l.gate}, nil
}

// newGatedWebSocketPair is newWebSocketPair with a gate on the server-side
// connection's writes.
func newGatedWebSocketPair(t *testing.T) (serverConn *websocket.Conn, gate *connGate, clientConn *websocket.Conn) {
	t.Helper()
	gate = &connGate{}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	connCh := make(chan *websocket.Conn, 1)
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		connCh <- c
	})}
	go func() { _ = srv.Serve(&gatedListener{Listener: ln, gate: gate}) }()
	t.Cleanup(func() { srv.Close() })

	client, resp, err := websocket.DefaultDialer.Dial("ws://"+ln.Addr().String(), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	t.Cleanup(func() { client.Close() })

	select {
	case serverConn = <-connCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for WebSocket upgrade")
	}
	t.Cleanup(func() { serverConn.Close() })
	return serverConn, gate, client
}

// TestWSConnectionWrapper_Metrics_WriteQueueWaitAndSocketWrite blocks the write
// pump inside the socket write and verifies that a queued writer's queue wait
// grows by the blocked time, that the blocked socket write records its full
// duration once it returns, and that each operation records exactly once.
func TestWSConnectionWrapper_Metrics_WriteQueueWaitAndSocketWrite(t *testing.T) {
	t.Parallel()

	obs := &recordingObserver{}
	connWrapper := network.NewWSConnectionWrapperWithObserver(logger.Test(t), obs)
	servicetest.Run(t, connWrapper)

	serverConn, gate, _ := newGatedWebSocketPair(t)
	connWrapper.Reset(serverConn)

	// The pump accepts the first write, then blocks in WriteMessage.
	gate.block()
	t.Cleanup(gate.unblock)
	firstErr := make(chan error, 1)
	go func() { firstErr <- connWrapper.Write(t.Context(), websocket.BinaryMessage, []byte("first")) }()
	require.Eventually(t, func() bool { return len(obs.queueWaitSamples()) == 1 },
		5*time.Second, time.Millisecond, "write pump should accept the first write and block in the socket write")

	// The second write cannot be accepted while the pump is blocked, so it
	// waits in Write's first select and is counted as a pending writer.
	secondErr := make(chan error, 1)
	go func() { secondErr <- connWrapper.Write(t.Context(), websocket.BinaryMessage, []byte("second")) }()
	require.Eventually(t, func() bool { return obs.pendingWriters() == 1 },
		5*time.Second, time.Millisecond, "queued writer should be counted as pending")

	const blockFor = 250 * time.Millisecond
	time.Sleep(blockFor)
	gate.unblock()

	require.NoError(t, <-firstErr)
	require.NoError(t, <-secondErr)

	// One queue-wait sample per write; the blocked writer's wait includes the delay.
	queueWaits := obs.queueWaitSamples()
	require.Len(t, queueWaits, 2, "each completed write must record its queue wait exactly once")
	require.Less(t, queueWaits[0], queueWaits[1])
	require.GreaterOrEqual(t, queueWaits[1], 200*time.Millisecond, "blocked writer's queue wait must include the pump blockage")

	// One socket-write sample per message; the blocked write includes the delay.
	socketWrites := obs.socketWriteSamples()
	require.Len(t, socketWrites, 2, "each WriteMessage must record its duration exactly once")
	require.GreaterOrEqual(t, socketWrites[0], 200*time.Millisecond, "blocked socket write must record its full duration")

	require.Zero(t, obs.pendingWriters(), "pending writers must return to zero")
}

// TestWSConnectionWrapper_Metrics_PendingWritersCanceled verifies that a writer
// leaving the first select via context cancellation decrements the pending
// writer count, records no queue-wait sample, and still returns the context error.
func TestWSConnectionWrapper_Metrics_PendingWritersCanceled(t *testing.T) {
	t.Parallel()

	obs := &recordingObserver{}
	connWrapper := network.NewWSConnectionWrapperWithObserver(logger.Test(t), obs)
	servicetest.Run(t, connWrapper)

	serverConn, gate, _ := newGatedWebSocketPair(t)
	connWrapper.Reset(serverConn)

	gate.block()
	t.Cleanup(gate.unblock)
	firstErr := make(chan error, 1)
	go func() { firstErr <- connWrapper.Write(t.Context(), websocket.BinaryMessage, []byte("first")) }()
	require.Eventually(t, func() bool { return len(obs.queueWaitSamples()) == 1 },
		5*time.Second, time.Millisecond, "write pump should accept the first write and block in the socket write")

	ctx, cancel := context.WithCancel(t.Context())
	writeErr := make(chan error, 1)
	go func() { writeErr <- connWrapper.Write(ctx, websocket.BinaryMessage, []byte("canceled")) }()
	require.Eventually(t, func() bool { return obs.pendingWriters() == 1 },
		5*time.Second, time.Millisecond, "queued writer should be counted as pending")

	cancel()
	require.ErrorIs(t, <-writeErr, context.Canceled)
	require.Eventually(t, func() bool { return obs.pendingWriters() == 0 },
		5*time.Second, time.Millisecond, "cancellation must decrement the pending writer count")
	require.Len(t, obs.queueWaitSamples(), 1, "a canceled writer must not record a queue-wait sample")

	gate.unblock()
	require.NoError(t, <-firstErr)
}

// TestWSConnectionWrapper_Metrics_ReadDispatchWait verifies that a message that
// was read successfully but not consumed promptly records a dispatch wait equal
// to the consumer blockage, exactly once.
func TestWSConnectionWrapper_Metrics_ReadDispatchWait(t *testing.T) {
	t.Parallel()

	obs := &recordingObserver{}
	connWrapper := network.NewWSConnectionWrapperWithObserver(logger.Test(t), obs)
	servicetest.Run(t, connWrapper)

	serverConn, clientConn := newWebSocketPair(t)
	connWrapper.Reset(serverConn)

	require.NoError(t, clientConn.WriteMessage(websocket.TextMessage, []byte("hello")))

	// Do not consume from ReadChannel for a controlled delay: the read pump is
	// blocked handing the item to readCh.
	const holdFor = 250 * time.Millisecond
	time.Sleep(holdFor)
	item := <-connWrapper.ReadChannel()
	require.Equal(t, websocket.TextMessage, item.MsgType)
	require.Equal(t, []byte("hello"), item.Data)

	require.Eventually(t, func() bool {
		waits := obs.dispatchWaitSamples()
		return len(waits) == 1 && waits[0] >= 200*time.Millisecond
	}, 5*time.Second, time.Millisecond, "blocked read consumer must produce one dispatch-wait sample covering the delay")
}
