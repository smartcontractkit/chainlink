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
	"github.com/smartcontractkit/chainlink/core/store"
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

type NoopWebSocketClient struct{}

func (NoopWebSocketClient) Url() url.URL                             { return url.URL{} }
func (NoopWebSocketClient) Status() ConnectionStatus                 { return ConnectionStatusDisconnected }
func (NoopWebSocketClient) Start() error                             { return nil }
func (NoopWebSocketClient) Close() error                             { return nil }
func (NoopWebSocketClient) Send([]byte)                              {}
func (NoopWebSocketClient) Receive(...time.Duration) ([]byte, error) { return nil, nil }

type websocketClient struct {
	boot      *sync.Mutex
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

// NewWebSocketClient returns a stats pusher using a websocket for
// delivery.
func NewWebSocketClient(url *url.URL, accessKey, secret string) ExplorerClient {
	return &websocketClient{
		url:       url,
		send:      make(chan []byte),
		receive:   make(chan []byte),
		boot:      &sync.Mutex{},
		sleeper:   utils.NewBackoffSleeper(),
		status:    ConnectionStatusDisconnected,
		accessKey: accessKey,
		secret:    secret,

		closeRequested: make(chan struct{}),
		closed:         make(chan struct{}),
	}
}

// Url returns the URL the client was initialized with
func (w *websocketClient) Url() url.URL {
	return *w.url
}

// Status returns the current connection status
func (w *websocketClient) Status() ConnectionStatus {
	w.statusMtx.RLock()
	defer w.statusMtx.RUnlock()
	return w.status
}

// Start starts a write pump over a websocket.
func (w *websocketClient) Start() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.started {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go w.connectAndWritePump(ctx, wg)
	wg.Wait()
	w.started = true
	return nil
}

// Send sends data asynchronously across the websocket if it's open, or
// holds it in a small buffer until connection, throwing away messages
// once buffer is full.
func (w *websocketClient) Send(data []byte) {
	w.send <- data
}

// Receive blocks the caller while waiting for a response from the server,
// returning the raw response bytes
func (w *websocketClient) Receive(durationParams ...time.Duration) ([]byte, error) {
	duration := defaultReceiveTimeout
	if len(durationParams) > 0 {
		duration = durationParams[0]
	}

	select {
	case <-time.After(duration):
		return nil, ErrReceiveTimeout
	case data := <-w.receive:
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
// to clean up independent of itself by reducing shared state. i.e. a passed done, not w.done.
func (w *websocketClient) connectAndWritePump(parentCtx context.Context, wg *sync.WaitGroup) {
	doneWaiting := false
	logger.Infow("Connecting to explorer", "url", w.url)

	for {
		select {
		case <-parentCtx.Done():
			if !doneWaiting {
				wg.Done()
			}
			return
		case <-time.After(w.sleeper.After()):
			connectionCtx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			if err := w.connect(connectionCtx); err != nil {
				w.setStatus(ConnectionStatusError)
				if !doneWaiting {
					wg.Done()
				}
				logger.Warn("Failed to connect to explorer (", w.url.String(), "): ", err)
				break
			}

			w.setStatus(ConnectionStatusConnected)

			if !doneWaiting {
				wg.Done()
			}
			logger.Infow("Connected to explorer", "url", w.url)
			w.sleeper.Reset()
			go w.readPump(cancel)
			w.writePump(connectionCtx)
		}

		doneWaiting = true
	}
}

func (w *websocketClient) setStatus(s ConnectionStatus) {
	w.statusMtx.Lock()
	defer w.statusMtx.Unlock()
	w.status = s
}

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (w *websocketClient) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.wrapConnErrorIf(w.conn.Close()) // exclusive responsibility to close ws conn
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case message, open := <-w.send:
			if !open { // channel closed
				w.wrapConnErrorIf(w.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			err := w.writeMessage(message)
			if err != nil {
				logger.Error("websocketStatsPusher: ", err)
				return
			}
		case <-ticker.C:
			w.wrapConnErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				w.wrapConnErrorIf(err)
				return
			}
		}
	}
}

func (w *websocketClient) writeMessage(message []byte) error {
	w.wrapConnErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
	writer, err := w.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		return err
	}

	return writer.Close()
}

func (w *websocketClient) connect(ctx context.Context) error {
	authHeader := http.Header{}
	authHeader.Add("X-Explore-Chainlink-Accesskey", w.accessKey)
	authHeader.Add("X-Explore-Chainlink-Secret", w.secret)
	authHeader.Add("X-Explore-Chainlink-Core-Version", store.Version)
	authHeader.Add("X-Explore-Chainlink-Core-Sha", store.Sha)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, w.url.String(), authHeader)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#connect: %v", err)
	}

	w.conn = conn
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
func (w *websocketClient) readPump(cancel context.CancelFunc) {
	defer cancel()

	w.conn.SetReadLimit(maxMessageSize)
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, message, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
				logger.Warn(fmt.Sprintf("readPump: %v", err))
			}
			select {
			case <-w.closeRequested:
				w.closed <- struct{}{}
			case <-time.After(CloseTimeout):
				logger.Warn("websocket readPump failed to notify closer")
			}
			return
		}

		switch messageType {
		case websocket.TextMessage:
			w.receive <- message
		}
	}
}

func (w *websocketClient) wrapConnErrorIf(err error) {
	if err != nil && websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
		w.setStatus(ConnectionStatusError)
		logger.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (w *websocketClient) Close() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.started {
		w.cancel()
	}
	w.started = false
	select {
	case w.closeRequested <- struct{}{}:
		<-w.closed
	case <-time.After(CloseTimeout):
		logger.Warn("websocketClient.Close failed to be notified from readPump")
	}
	return nil
}
