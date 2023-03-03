package sm_yamux

import (
	"context"

	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-yamux/v2"
)

// conn implements mux.MuxedConn over yamux.Session.
type conn yamux.Session

// Close closes underlying yamux
func (c *conn) Close() error {
	return c.yamux().Close()
}

// IsClosed checks if yamux.Session is in closed state.
func (c *conn) IsClosed() bool {
	return c.yamux().IsClosed()
}

// OpenStream creates a new stream.
func (c *conn) OpenStream(ctx context.Context) (mux.MuxedStream, error) {
	s, err := c.yamux().OpenStream(ctx)
	if err != nil {
		return nil, err
	}

	return (*stream)(s), nil
}

// AcceptStream accepts a stream opened by the other side.
func (c *conn) AcceptStream() (mux.MuxedStream, error) {
	s, err := c.yamux().AcceptStream()
	return (*stream)(s), err
}

func (c *conn) yamux() *yamux.Session {
	return (*yamux.Session)(c)
}

var _ mux.MuxedConn = &conn{}
