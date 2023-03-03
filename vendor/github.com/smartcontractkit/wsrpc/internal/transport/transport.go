package transport

import (
	"context"
	"time"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/wsrpc/credentials"
)

const (
	// Time allowed to write a message to the connection.
	defaultWriteTimeout = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// ConnectOptions covers all relevant options for communicating with the server.
type ConnectOptions struct {
	WriteTimeout time.Duration
	// TransportCredentials stores the Authenticator required to setup a client
	// connection.
	TransportCredentials credentials.TransportCredentials
}

// ClientTransport is the common interface for wsrpc client-side transport
// implementations.
type ClientTransport interface {
	// Read reads a message from the stream
	Read() <-chan []byte

	// Write sends a message to the stream.
	Write(msg []byte) error

	// Close tears down this transport. Once it returns, the transport
	// should not be accessed any more.
	Close() error
}

// NewClientTransport establishes the transport with the required ConnectOptions
// and returns it to the caller.
func NewClientTransport(ctx context.Context, addr string, opts ConnectOptions, onClose func()) (ClientTransport, error) {
	return newWebsocketClient(ctx, addr, opts, onClose)
}

// state of transport.
type transportState int

const (
	// The default transport state.
	//
	// nolint is required because we don't actually use the var anywhere,
	// but it does represent a reachable transport.
	reachable transportState = iota //nolint:deadcode,varcheck
	closing
)

// ServerConfig consists of all the configurations to establish a server transport.
type ServerConfig struct {
	WriteTimeout time.Duration
}

// ServerTransport is the common interface for wsrpc server-side transport
// implementations.
type ServerTransport interface {
	// Read reads a message from the stream.
	Read() <-chan []byte

	// Write sends a message to the stream.
	Write(msg []byte) error

	// Close tears down the transport. Once it is called, the transport
	// should not be accessed any more.
	Close() error
}

// NewServerTransport creates a ServerTransport with conn or non-nil error
// if it fails.
func NewServerTransport(c *websocket.Conn, config *ServerConfig, onClose func()) (ServerTransport, error) {
	return newWebsocketServer(c, config, onClose), nil
}

func handlePong(conn *websocket.Conn) func(string) error {
	return func(msg string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}
