package synchronization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	// ErrReceiveTimeout is returned when no message is received after a
	// specified duration in Receive
	ErrReceiveTimeout = errors.New("timeout waiting for message")
)

// WebSocketClient encapsulates all the functionality needed to
// push run information to explorer.
type WebSocketClient interface {
	Start() error
	Close() error
	Send([]byte)
	Receive(...time.Duration) ([]byte, error)
}

type noopWebSocketClient struct{}

func (noopWebSocketClient) Start() error                             { return nil }
func (noopWebSocketClient) Close() error                             { return nil }
func (noopWebSocketClient) Send([]byte)                              {}
func (noopWebSocketClient) Receive(...time.Duration) ([]byte, error) { return nil, nil }

type websocketClient struct {
	boot      *sync.Mutex
	conn      *websocket.Conn
	cancel    context.CancelFunc
	send      chan []byte
	receive   chan []byte
	sleeper   utils.Sleeper
	started   bool
	url       *url.URL
	accessKey string
	secret    string
}

// NewWebSocketClient returns a stats pusher using a websocket for
// delivery.
func NewWebSocketClient(url *url.URL, accessKey, secret string) WebSocketClient {
	return &websocketClient{
		url:       url,
		send:      make(chan []byte),
		receive:   make(chan []byte),
		boot:      &sync.Mutex{},
		sleeper:   utils.NewBackoffSleeper(),
		accessKey: accessKey,
		secret:    secret,
	}
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
	wg.Done()
	logger.Infow("Connecting to explorer", "url", w.url)

	for {
		select {
		case <-parentCtx.Done():
			return
		case <-time.After(w.sleeper.After()):
			connectionCtx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			if err := w.connect(connectionCtx); err != nil {
				logger.Warn("Failed to connect to explorer (", w.url.String(), "): ", err)
				break
			}

			logger.Info("Connected to explorer at ", w.url.String())
			w.sleeper.Reset()
			go w.readPump(cancel)
			w.writePump(connectionCtx)
		}
	}
}

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (w *websocketClient) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wrapConnErrorIf(w.conn.Close()) // exclusive responsibility to close ws conn
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case message, open := <-w.send:
			if !open { // channel closed
				wrapConnErrorIf(w.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			err := w.writeMessage(message)
			if err != nil {
				logger.Error("websocketStatsPusher: ", err)
				return
			}
		case <-ticker.C:
			wrapConnErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				wrapConnErrorIf(err)
				return
			}
		}
	}
}

func (w *websocketClient) writeMessage(message []byte) error {
	wrapConnErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
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
	authHeader.Add("X-Explore-Chainlink-AccessKey", w.accessKey)
	authHeader.Add("X-Explore-Chainlink-Secret", w.secret)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, w.url.String(), authHeader)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#connect: %v", err)
	}

	w.conn = conn
	return nil
}

var expectedCloseMessages = []int{websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure}

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
			return
		}

		switch messageType {
		case websocket.TextMessage:
			w.receive <- message
		}
	}
}

func wrapConnErrorIf(err error) {
	if err != nil && websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
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
	return nil
}
