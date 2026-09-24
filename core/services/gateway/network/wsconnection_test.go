package network_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var errConnClosed = errors.New("connection closed")

type mockConnection struct {
	readChan  chan msg
	writeChan chan msg
	closed    chan struct{}
	closeOnce *sync.Once
	readErr   error
	writeErr  error
	// when non-nil, WriteMessage blocks until it is closed
	writeGate chan struct{}
}

type msg struct {
	messageType int
	p           []byte
}

func (m *mockConnection) ReadMessage() (messageType int, p []byte, err error) {
	select {
	case <-m.closed:
		return 0, nil, errConnClosed
	case msg := <-m.readChan:
		return msg.messageType, msg.p, m.readErr
	}
}

func (m *mockConnection) WriteMessage(messageType int, data []byte) error {
	if m.writeGate != nil {
		select {
		case <-m.closed:
			return errConnClosed
		case <-m.writeGate:
		}
	}

	if m.writeErr != nil {
		return m.writeErr
	}

	// Checked first because a select would pick randomly between a closed
	// conn and a writeChan with buffer space.
	select {
	case <-m.closed:
		return errConnClosed
	default:
	}

	msg := msg{messageType: messageType, p: data}
	select {
	case <-m.closed:
		return errConnClosed
	case m.writeChan <- msg:
		return nil
	}
}

func (m *mockConnection) Close() error {
	m.closeOnce.Do(func() { close(m.closed) })
	return nil
}

// newWebSocketPair returns two connected ends that share a close signal, so
// closing either end unblocks reads on both, like a real socket.
func newWebSocketPair(t *testing.T) (*mockConnection, *mockConnection) {
	t.Helper()
	chAtoB, chBtoA := make(chan msg, 10), make(chan msg, 10)
	closed, closeOnce := make(chan struct{}), &sync.Once{}

	return &mockConnection{closed: closed, closeOnce: closeOnce, readChan: chBtoA, writeChan: chAtoB},
		&mockConnection{closed: closed, closeOnce: closeOnce, readChan: chAtoB, writeChan: chBtoA}
}

// TestWSConnectionWrapper_WriteError_TriggersClose verifies that a write error
// (e.g. i/o timeout on a stale connection) causes the connection to be closed,
// which signals closeCh and allows reconnectLoop to re-establish a fresh connection.
func TestWSConnectionWrapper_WriteError_TriggersClose(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, clientConn := newWebSocketPair(t)

	// client
	clientConnWrapper := network.NewWSConnectionWrapper(lggr)
	servicetest.Run(t, clientConnWrapper)
	closeCh := clientConnWrapper.Reset(clientConn)

	clientConn.writeErr = errors.New("write error")

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

	clientConnWrapper := network.NewWSConnectionWrapper(logger.Test(t))
	servicetest.Run(t, clientConnWrapper)

	// connect, write a message, disconnect
	serverConn, clientConn := newWebSocketPair(t)
	clientConnWrapper.Reset(clientConn)
	require.NoError(t, clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("hello")))
	_, got, err := serverConn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), got)
	require.NoError(t, clientConn.Close())

	// try to write without a connection
	writeErr := clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("failed send"))
	require.Error(t, writeErr)

	// re-connect, write another message
	serverConn, clientConn = newWebSocketPair(t)
	clientConnWrapper.Reset(clientConn)
	require.NoError(t, clientConnWrapper.Write(t.Context(), websocket.TextMessage, []byte("hello again")))
	_, got, err = serverConn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, []byte("hello again"), got)
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

// newGatedWebSocketPair is newWebSocketPair with the server end's writes
// blocked until the returned unblock func is called.
func newGatedWebSocketPair(t *testing.T) (serverConn *mockConnection, unblock func()) {
	t.Helper()
	serverConn, _ = newWebSocketPair(t)
	gate := make(chan struct{})
	serverConn.writeGate = gate
	unblock = sync.OnceFunc(func() { close(gate) })
	t.Cleanup(unblock)
	return serverConn, unblock
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

	serverConn, unblock := newGatedWebSocketPair(t)
	connWrapper.Reset(serverConn)

	// The pump accepts the first write, then blocks in WriteMessage.
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
	unblock()

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

	serverConn, unblock := newGatedWebSocketPair(t)
	connWrapper.Reset(serverConn)

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

	unblock()
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
