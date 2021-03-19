package synchronization

import (
	"errors"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/gorilla/websocket"
)

var (
	// ErrReceiveTimeout is returned when no message is received after a
	// specified duration in Receive
	ErrReceiveTimeout = errors.New("timeout waiting for message")
)

type ConnectionStatus string

const (
	// ConnectionStatusDisconnected is the default state
	ConnectionStatusDisconnected = ConnectionStatus("disconnected")
	// ConnectionStatusConnected is used when the client is successfully connected
	ConnectionStatusConnected = ConnectionStatus("connected")
	// ConnectionStatusError is used when there is an error
	ConnectionStatusError = ConnectionStatus("error")
)

// SendBufferSize is the number of messages to keep in the buffer before dropping additional ones
const SendBufferSize = 100

const (
	ExplorerTextMessage   = websocket.TextMessage
	ExplorerBinaryMessage = websocket.BinaryMessage
)

// ExplorerClient encapsulates all the functionality needed to
// push run information to explorer.
type ExplorerClient interface {
	Url() url.URL
	Status() ConnectionStatus
	Start() error
	Close() error
	Send([]byte, ...int)
	Receive(...time.Duration) ([]byte, error)
}

type NoopExplorerClient struct{}

func (NoopExplorerClient) Url() url.URL                             { return url.URL{} }
func (NoopExplorerClient) Status() ConnectionStatus                 { return ConnectionStatusDisconnected }
func (NoopExplorerClient) Start() error                             { return nil }
func (NoopExplorerClient) Close() error                             { return nil }
func (NoopExplorerClient) Send([]byte, ...int)                      {}
func (NoopExplorerClient) Receive(...time.Duration) ([]byte, error) { return nil, nil }

type explorerClient struct {
}

// NewExplorerClient returns a stats pusher using a websocket for
// delivery.
func NewExplorerClient(url *url.URL, accessKey, secret string) ExplorerClient {
	return &explorerClient{}
}

// Url returns the URL the client was initialized with
func (ec *explorerClient) Url() url.URL {
	return url.URL{}
}

// Status returns the current connection status
func (ec *explorerClient) Status() ConnectionStatus {
	return ConnectionStatusConnected
}

// Start starts a write pump over a websocket.
func (ec *explorerClient) Start() error {
	logger.Debugw("explorerClient#Start")
	return nil
}

// Send sends data asynchronously across the websocket if it's open, or
// holds it in a small buffer until connection, throwing away messages
// once buffer is full.
// func (ec *explorerClient) Receive(durationParams ...time.Duration) ([]byte, error) {
func (ec *explorerClient) Send(data []byte, messageTypes ...int) {
	logger.Debugw("explorerClient#Send", "data", data, "messageTypes", messageTypes)
}

// Receive blocks the caller while waiting for a response from the server,
// returning the raw response bytes
func (ec *explorerClient) Receive(durationParams ...time.Duration) ([]byte, error) {
	return nil, nil
}

func (ec *explorerClient) Close() error {
	return nil
}
