package synchronization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"

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

// ExplorerClient encapsulates all the functionality needed to
// push run information to explorer.
type ExplorerClient interface {
	Url() url.URL
	Status() ConnectionStatus
	Start() error
	Close() error
	Send([]byte)
	Receive(...time.Duration) ([]byte, error)
}

type NoopExplorerClient struct{}

func (NoopExplorerClient) Url() url.URL                             { return url.URL{} }
func (NoopExplorerClient) Status() ConnectionStatus                 { return ConnectionStatusDisconnected }
func (NoopExplorerClient) Start() error                             { return nil }
func (NoopExplorerClient) Close() error                             { return nil }
func (NoopExplorerClient) Send([]byte)                              {}
func (NoopExplorerClient) Receive(...time.Duration) ([]byte, error) { return nil, nil }

type explorerClient struct {
	boot      *sync.RWMutex
	conn      *websocket.Conn
	cancel    context.CancelFunc
	send      chan []byte
	receive   chan []byte
	sleeper   utils.Sleeper
	started   bool
	status    ConnectionStatus
	url       *url.URL
	accessKey string
	secret    string

	closeRequested chan struct{}
	closed         chan struct{}

	statusMtx sync.RWMutex
}

// NewExplorerClient returns a stats pusher using a websocket for
// delivery.
func NewExplorerClient(url *url.URL, accessKey, secret string) ExplorerClient {
	return &explorerClient{
		url:       url,
		send:      make(chan []byte),
		receive:   make(chan []byte),
		boot:      new(sync.RWMutex),
		sleeper:   utils.NewBackoffSleeper(),
		status:    ConnectionStatusDisconnected,
		accessKey: accessKey,
		secret:    secret,

		closeRequested: make(chan struct{}),
		closed:         make(chan struct{}),
	}
}

// Url returns the URL the client was initialized with
func (ec *explorerClient) Url() url.URL {
	return *ec.url
}

// Status returns the current connection status
func (ec *explorerClient) Status() ConnectionStatus {
	ec.statusMtx.RLock()
	defer ec.statusMtx.RUnlock()
	return ec.status
}

// Start starts a write pump over a websocket.
func (ec *explorerClient) Start() error {
	ec.boot.Lock()
	defer ec.boot.Unlock()

	if ec.started {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	ec.cancel = cancel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go ec.connectAndWritePump(ctx, wg)
	wg.Wait()
	ec.started = true
	return nil
}

// Send sends data asynchronously across the websocket if it's open, or
// holds it in a small buffer until connection, throwing away messages
// once buffer is full.
func (ec *explorerClient) Send(data []byte) {
	ec.boot.RLock()
	defer ec.boot.RUnlock()
	if !ec.started {
		panic("send on unstarted explorer client")
	}
	ec.send <- data
}

// Receive blocks the caller while waiting for a response from the server,
// returning the raw response bytes
func (ec *explorerClient) Receive(durationParams ...time.Duration) ([]byte, error) {
	duration := defaultReceiveTimeout
	if len(durationParams) > 0 {
		duration = durationParams[0]
	}

	select {
	case <-time.After(duration):
		return nil, ErrReceiveTimeout
	case data := <-ec.receive:
		return data, nil
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// defaultReceiveTimeout is the default amount of time to wait for receipt of messages
	defaultReceiveTimeout = 30 * time.Second
)

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
// lexical confinement of done chan allows multiple connectAndWritePump routines
// to clean up independent of itself by reducing shared state. i.e. a passed done, not ec.done.
func (ec *explorerClient) connectAndWritePump(parentCtx context.Context, wg *sync.WaitGroup) {
	doneWaiting := false
	logger.Infow("Connecting to explorer", "url", ec.url)

	for {
		select {
		case <-parentCtx.Done():
			if !doneWaiting {
				wg.Done()
			}
			return
		case <-time.After(ec.sleeper.After()):
			connectionCtx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			if err := ec.connect(connectionCtx); err != nil {
				ec.setStatus(ConnectionStatusError)
				if !doneWaiting {
					wg.Done()
				}
				logger.Warn("Failed to connect to explorer (", ec.url.String(), "): ", err)
				break
			}

			ec.setStatus(ConnectionStatusConnected)

			if !doneWaiting {
				wg.Done()
			}
			logger.Infow("Connected to explorer", "url", ec.url)
			ec.sleeper.Reset()
			go ec.readPump(cancel)
			ec.writePump(connectionCtx)
		}

		doneWaiting = true
	}
}

func (ec *explorerClient) setStatus(s ConnectionStatus) {
	ec.statusMtx.Lock()
	defer ec.statusMtx.Unlock()
	ec.status = s
}

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (ec *explorerClient) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ec.wrapConnErrorIf(ec.conn.Close()) // exclusive responsibility to close ws conn
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case message, open := <-ec.send:
			if !open { // channel closed
				ec.wrapConnErrorIf(ec.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			err := ec.writeMessage(message)
			if err != nil {
				logger.Error("websocketStatsPusher: ", err)
				return
			}
		case <-ticker.C:
			ec.wrapConnErrorIf(ec.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := ec.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ec.wrapConnErrorIf(err)
				return
			}
		}
	}
}

func (ec *explorerClient) writeMessage(message []byte) error {
	ec.wrapConnErrorIf(ec.conn.SetWriteDeadline(time.Now().Add(writeWait)))
	writer, err := ec.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		return err
	}

	return writer.Close()
}

func (ec *explorerClient) connect(ctx context.Context) error {
	authHeader := http.Header{}
	authHeader.Add("X-Explore-Chainlink-Accesskey", ec.accessKey)
	authHeader.Add("X-Explore-Chainlink-Secret", ec.secret)
	authHeader.Add("X-Explore-Chainlink-Core-Version", static.Version)
	authHeader.Add("X-Explore-Chainlink-Core-Sha", static.Sha)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, ec.url.String(), authHeader)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#connect: %v", err)
	}

	ec.conn = conn
	return nil
}

var expectedCloseMessages = []int{websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure}

const CloseTimeout = 100 * time.Millisecond

// readPump listens on the websocket connection for control messages and
// response messages (text)
//
// For more details on how disconnection messages are handled, see:
//  * https://stackoverflow.com/a/48181794/639773
//  * https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L56
func (ec *explorerClient) readPump(cancel context.CancelFunc) {
	defer cancel()

	ec.conn.SetReadLimit(maxMessageSize)
	_ = ec.conn.SetReadDeadline(time.Now().Add(pongWait))
	ec.conn.SetPongHandler(func(string) error {
		_ = ec.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, message, err := ec.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
				logger.Warn(fmt.Sprintf("readPump: %v", err))
			}
			select {
			case <-ec.closeRequested:
				ec.closed <- struct{}{}
			case <-time.After(CloseTimeout):
				logger.Warn("websocket readPump failed to notify closer")
			}
			return
		}

		switch messageType {
		case websocket.TextMessage:
			ec.receive <- message
		}
	}
}

func (ec *explorerClient) wrapConnErrorIf(err error) {
	if err != nil && websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
		ec.setStatus(ConnectionStatusError)
		logger.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (ec *explorerClient) Close() error {
	ec.boot.Lock()
	defer ec.boot.Unlock()

	if ec.started {
		ec.cancel()
	}
	ec.started = false
	select {
	case ec.closeRequested <- struct{}{}:
		<-ec.closed
	case <-time.After(CloseTimeout):
		logger.Warn("websocketClient.Close failed to be notified from readPump")
	}
	return nil
}
