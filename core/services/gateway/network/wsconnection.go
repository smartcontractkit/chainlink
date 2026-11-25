package network

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// WSConnectionWrapper is a websocket connection abstraction that supports re-connects.
// I/O is separated from connection management:
//   - component doing writes can use the thread-safe Write() method
//   - component doing reads can listen on the ReadChannel()
//   - component managing connections can listen to connection-closed channels and call Reset()
//     to swap the underlying connection object
//
// The Wrapper can be used by a server expecting long-lived connections from a given client,
// as well as a client maintaining such long-lived connection with a given server.
// This fits the Gateway very well as servers accept connections only from a fixed set of nodes
// and conversely, nodes only connect to a fixed set of servers (Gateways).
//
// The concept of "pumps" is borrowed from https://github.com/smartcontractkit/wsrpc
// All methods are thread-safe.
type WSConnectionWrapper interface {
	job.ServiceCtx
	services.HealthReporter

	// Update underlying connection object. Return a channel that gets an error on connection close.
	// Cannot be called after Close().
	// To disconnect, pass an untyped nil. A nil pointer wrapped in the interface
	// (e.g. a nil *websocket.Conn) is non-nil, so it is treated as a live
	// connection and the read pump will panic when it reads from it.
	Reset(newConn WSConnection) <-chan error

	Write(ctx context.Context, msgType int, data []byte) error

	ReadChannel() <-chan ReadItem

	IsConnected() bool
}

type WSConnection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type wsConnectionHolder struct {
	WSConnection
}

type wsConnectionWrapper struct {
	services.StateMachine
	lggr logger.Logger

	conn atomic.Pointer[wsConnectionHolder]

	writeCh    chan writeItem
	readCh     chan ReadItem
	shutdownCh chan struct{}
	wg         services.WaitGroup
	metrics    WSConnectionObserver

	pingInterval time.Duration
	pongTimeout  time.Duration
	writeTimeout time.Duration
}

func (c *wsConnectionWrapper) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}

func (c *wsConnectionWrapper) Name() string { return c.lggr.Name() }

type ReadItem struct {
	MsgType int
	Data    []byte
}

type writeItem struct {
	MsgType int
	Data    []byte
	ErrCh   chan error
}

var _ WSConnectionWrapper = (*wsConnectionWrapper)(nil)

var (
	ErrNoActiveConnection = errors.New("no active connection")
	ErrWrapperShutdown    = errors.New("wrapper shutting down")
)

type WSConnectionOption func(*wsConnectionWrapper)

// WithKeepalive pings the peer every pingInterval and closes the connection
// when no pong arrives within pongTimeout. Without it, a peer whose host
// vanished is only noticed once the kernel gives up retransmitting (~15 min).
// It replaces the connection's pong handler, so don't combine it with one.
func WithKeepalive(pingInterval, pongTimeout time.Duration) WSConnectionOption {
	return func(c *wsConnectionWrapper) {
		c.pingInterval = pingInterval
		c.pongTimeout = pongTimeout
	}
}

// WithWriteTimeout fails a socket write that doesn't complete within timeout,
// which closes the connection and lets the owner reconnect.
func WithWriteTimeout(timeout time.Duration) WSConnectionOption {
	return func(c *wsConnectionWrapper) {
		c.writeTimeout = timeout
	}
}

// keepaliveConn is the part of *websocket.Conn that keepalive needs.
type keepaliveConn interface {
	SetReadDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	WriteControl(messageType int, data []byte, deadline time.Time) error
}

type writeDeadlineConn interface {
	SetWriteDeadline(t time.Time) error
}

func NewWSConnectionWrapper(lggr logger.Logger, opts ...WSConnectionOption) WSConnectionWrapper {
	return NewWSConnectionWrapperWithObserver(lggr, nil, opts...)
}

// NewWSConnectionWrapperWithObserver is NewWSConnectionWrapper with a metrics
// observer for the wrapper's blocking operations. A nil observer discards all
// observations.
func NewWSConnectionWrapperWithObserver(lggr logger.Logger, observer WSConnectionObserver, opts ...WSConnectionOption) WSConnectionWrapper {
	if observer == nil {
		observer = noopWSConnectionObserver{}
	}
	cw := &wsConnectionWrapper{
		lggr:       logger.Named(lggr, "WSConnectionWrapper"),
		writeCh:    make(chan writeItem),
		readCh:     make(chan ReadItem),
		shutdownCh: make(chan struct{}),
		metrics:    observer,
	}
	for _, opt := range opts {
		opt(cw)
	}
	return cw
}

func (c *wsConnectionWrapper) Start(_ context.Context) error {
	return c.StartOnce("WSConnectionWrapper", func() error {
		if err := c.wg.TryAdd(1); err != nil {
			return err // should never happen
		}
		// write pump runs until Shutdown() is called
		go c.writePump()
		return nil
	})
}

// Reset:
//  1. replaces the underlying connection and shuts the old one down
//  2. starts a new read goroutine that pushes received messages to readCh
//  3. returns channel that closes when connection closes, or nil if closed or newConn is nil.
func (c *wsConnectionWrapper) Reset(newConn WSConnection) <-chan error {
	var conn *wsConnectionHolder
	if newConn != nil {
		conn = &wsConnectionHolder{newConn}
	}

	oldConn := c.conn.Swap(conn)

	if oldConn != nil {
		oldConn.Close()
	}
	if newConn == nil {
		return nil
	}
	if err := c.wg.TryAdd(1); err != nil {
		return nil
	}
	closeCh := make(chan error, 1)
	readDone := make(chan struct{})
	c.startKeepalive(newConn, readDone)
	// readPump goroutine is tied to the lifecycle of the underlying conn object
	go c.readPump(conn, closeCh, readDone)
	return closeCh
}

// startKeepalive arms a read deadline that only pongs extend, and pings until stop is closed.
// An expired deadline fails the pending read, so readPump closes the connection.
func (c *wsConnectionWrapper) startKeepalive(conn WSConnection, stop <-chan struct{}) {
	kc, ok := conn.(keepaliveConn)
	if !ok || c.pingInterval <= 0 || c.pongTimeout <= 0 {
		return
	}
	extendDeadline := func(string) error {
		return kc.SetReadDeadline(time.Now().Add(c.pongTimeout))
	}
	if err := extendDeadline(""); err != nil {
		c.lggr.Errorw("failed to set read deadline", "error", err)
	}
	kc.SetPongHandler(extendDeadline)
	if err := c.wg.TryAdd(1); err != nil {
		return
	}
	go c.pingLoop(kc, stop)
}

func (c *wsConnectionWrapper) pingLoop(conn keepaliveConn, stop <-chan struct{}) {
	defer c.wg.Done()
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-c.shutdownCh:
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(c.pingInterval)); err != nil {
				c.lggr.Debugw("failed to send keepalive ping", "error", err)
			}
		}
	}
}

func (c *wsConnectionWrapper) armWriteDeadline(conn WSConnection) {
	if c.writeTimeout <= 0 {
		return
	}
	if wc, ok := conn.(writeDeadlineConn); ok {
		if err := wc.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			c.lggr.Errorw("failed to set write deadline", "error", err)
		}
	}
}

func (c *wsConnectionWrapper) Write(ctx context.Context, msgType int, data []byte) error {
	errCh := make(chan error, 1)
	// push to write channel
	c.metrics.AddPendingWriters(ctx, 1)
	queueStart := time.Now()
	select {
	case c.writeCh <- writeItem{msgType, data, errCh}:
		c.metrics.AddPendingWriters(ctx, -1)
		// The queue wait ends when the write pump accepts the item; the second
		// select below (waiting for the write result) is intentionally excluded.
		c.metrics.RecordWriteQueueWait(ctx, time.Since(queueStart))
	case <-c.shutdownCh:
		c.metrics.AddPendingWriters(ctx, -1)
		return ErrWrapperShutdown
	case <-ctx.Done():
		c.metrics.AddPendingWriters(ctx, -1)
		return ctx.Err()
	}
	// wait for write result
	select {
	case err := <-errCh:
		return err
	case <-c.shutdownCh:
		return ErrWrapperShutdown
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *wsConnectionWrapper) ReadChannel() <-chan ReadItem {
	return c.readCh
}

func (c *wsConnectionWrapper) IsConnected() bool {
	return c.conn.Load() != nil
}

func (c *wsConnectionWrapper) Close() error {
	return c.StopOnce("WSConnectionWrapper", func() error {
		close(c.shutdownCh)
		c.Reset(nil)
		c.wg.Wait()
		return nil
	})
}

func (c *wsConnectionWrapper) writePump() {
	defer c.wg.Done()
	for {
		select {
		case wsMsg := <-c.writeCh:
			// synchronization is a tradeoff for the ability to use a single write channel
			conn := c.conn.Load()
			if conn == nil {
				wsMsg.ErrCh <- ErrNoActiveConnection
				close(wsMsg.ErrCh)
				break
			}
			// Time exactly the socket write, recording on return (including
			// on error) so a stalled write still reports its full duration.
			writeStart := time.Now()
			c.armWriteDeadline(conn.WSConnection)
			err := conn.WriteMessage(wsMsg.MsgType, wsMsg.Data)
			c.metrics.RecordSocketWrite(context.Background(), time.Since(writeStart))
			if err != nil {
				c.lggr.Errorw("failed to write message", "msgType", wsMsg.MsgType, "dataLen", len(wsMsg.Data), "error", err)
				// A write failure (e.g. i/o timeout on a half-open TCP session) does not
				// produce a read error, so readPump would stay blocked and reconnectLoop
				// would never get the closeCh signal it needs to redial. Close the conn
				// here to unblock readPump and trigger the existing reconnect path.
				// If CAS fails, Reset already swapped the pointer and will close the old
				// conn itself — don't close it twice.
				if c.conn.CompareAndSwap(conn, nil) {
					if closeErr := conn.Close(); closeErr != nil {
						c.lggr.Errorw("error closing connection after write failure", "error", closeErr)
					}
				}
			}
			wsMsg.ErrCh <- err
			close(wsMsg.ErrCh)
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *wsConnectionWrapper) readPump(conn *wsConnectionHolder, closeCh chan<- error, done chan<- struct{}) {
	defer c.wg.Done()
	defer close(done)
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			c.lggr.Errorw("failed to read message, closing connection", "error", err)
			var closeErr error
			if c.conn.CompareAndSwap(conn, nil) {
				closeErr = conn.Close()
				if closeErr != nil {
					c.lggr.Errorw("error closing connection", "error", closeErr)
				}
			}
			closeCh <- closeErr
			close(closeCh)
			return
		}
		// Time only the handoff to the read channel consumer; ReadMessage above
		// is excluded because it blocks for normal idle time between messages.
		dispatchStart := time.Now()
		select {
		case c.readCh <- ReadItem{msgType, data}:
			c.metrics.RecordReadDispatchWait(context.Background(), time.Since(dispatchStart))
		case <-c.shutdownCh:
			var closeErr error
			if c.conn.CompareAndSwap(conn, nil) {
				closeErr = conn.Close()
				if closeErr != nil {
					c.lggr.Errorw("error closing connection during shutdown", "error", closeErr)
				}
			}
			closeCh <- closeErr
			close(closeCh)
			return
		}
	}
}
