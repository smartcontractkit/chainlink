package network

import (
	"context"
	"errors"
	"sync/atomic"

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
	Reset(newConn *websocket.Conn) <-chan error

	Write(ctx context.Context, msgType int, data []byte) error

	ReadChannel() <-chan ReadItem
}

type wsConnectionWrapper struct {
	services.StateMachine
	lggr logger.Logger

	conn atomic.Pointer[websocket.Conn]

	writeCh    chan writeItem
	readCh     chan ReadItem
	shutdownCh chan struct{}
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

func NewWSConnectionWrapper(lggr logger.Logger) WSConnectionWrapper {
	cw := &wsConnectionWrapper{
		lggr:       logger.Named(lggr, "WSConnectionWrapper"),
		writeCh:    make(chan writeItem),
		readCh:     make(chan ReadItem),
		shutdownCh: make(chan struct{}),
	}
	return cw
}

func (c *wsConnectionWrapper) Start(_ context.Context) error {
	return c.StartOnce("WSConnectionWrapper", func() error {
		// write pump runs until Shutdown() is called
		go c.writePump()
		return nil
	})
}

// Reset:
//  1. replaces the underlying connection and shuts the old one down
//  2. starts a new read goroutine that pushes received messages to readCh
//  3. returns channel that closes when connection closes
func (c *wsConnectionWrapper) Reset(newConn *websocket.Conn) <-chan error {
	oldConn := c.conn.Swap(newConn)

	if oldConn != nil {
		oldConn.Close()
	}
	if newConn == nil {
		return nil
	}
	closeCh := make(chan error, 1)
	// readPump goroutine is tied to the lifecycle of the underlying conn object
	go c.readPump(newConn, closeCh)
	return closeCh
}

func (c *wsConnectionWrapper) Write(ctx context.Context, msgType int, data []byte) error {
	errCh := make(chan error, 1)
	// push to write channel
	select {
	case c.writeCh <- writeItem{msgType, data, errCh}:
		break
	case <-c.shutdownCh:
		return ErrWrapperShutdown
	case <-ctx.Done():
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

func (c *wsConnectionWrapper) Close() error {
	return c.StopOnce("WSConnectionWrapper", func() error {
		close(c.shutdownCh)
		c.Reset(nil)
		return nil
	})
}

func (c *wsConnectionWrapper) writePump() {
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
			wsMsg.ErrCh <- conn.WriteMessage(wsMsg.MsgType, wsMsg.Data)
			close(wsMsg.ErrCh)
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *wsConnectionWrapper) readPump(conn *websocket.Conn, closeCh chan<- error) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			closeCh <- conn.Close()
			close(closeCh)
			return
		}
		select {
		case c.readCh <- ReadItem{msgType, data}:
		case <-c.shutdownCh:
			closeCh <- conn.Close()
			close(closeCh)
			return
		}
	}
}
