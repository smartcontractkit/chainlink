package store

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/utils"
)

// WebsocketClient encapsulates all the functionality needed to
// push run information to linkstats.
type WebsocketClient interface {
	Start() error
	Close() error
	Send([]byte)
}

// NewStatsPusher returns a functioning instance depending on the
// URL passed: nil is a noop instance, url assumes a websocket instance.
// No support for http.
func NewStatsPusher(url *url.URL) WebsocketClient {
	if url != nil {
		return NewWebsocketStatsPusher(url)
	}
	return noopStatsPusher{}
}

type noopStatsPusher struct{}

func (noopStatsPusher) Start() error { return nil }
func (noopStatsPusher) Close() error { return nil }
func (noopStatsPusher) Send([]byte)  {}

type websocketStatsPusher struct {
	boot    *sync.Mutex
	conn    *websocket.Conn
	cancel  context.CancelFunc
	send    chan []byte
	sleeper utils.Sleeper
	started bool
	url     url.URL
}

// NewWebsocketStatsPusher returns a stats pusher using a websocket for
// delivery.
func NewWebsocketStatsPusher(url *url.URL) WebsocketClient {
	return &websocketStatsPusher{
		url:     *url,
		send:    make(chan []byte, 100), // TODO: figure out a better buffer (circular FIFO?)
		boot:    &sync.Mutex{},
		sleeper: utils.NewBackoffSleeper(),
	}
}

// Start starts a write pump over a websocket.
func (w *websocketStatsPusher) Start() error {
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

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
// lexical confinement of done chan allows multiple connectAndWritePump routines
// to clean up independent of itself by reducing shared state. i.e. a passed done, not w.done.
func (w *websocketStatsPusher) connectAndWritePump(parentCtx context.Context, wg *sync.WaitGroup) {
	wg.Done()
	logger.Info("Connecting to linkstats at ", w.url.String())

	for {
		select {
		case <-parentCtx.Done():
			return
		case <-time.After(w.sleeper.After()):
			connectionCtx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			if err := w.connect(connectionCtx); err != nil {
				logger.Warn("Failed to connect to linkstats (", w.url.String(), "): ", err)
				break
			}

			logger.Info("Connected to linkstats at ", w.url.String())
			w.sleeper.Reset()
			go w.readPumpForControlMessages(cancel)
			w.writePump(connectionCtx)
		}
	}
}

var (
	newline = []byte{'\n'}
)

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (w *websocketStatsPusher) writePump(ctx context.Context) {
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
				return
			}

			wrapConnErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error("websocketStatsPusher: ", err)
				return
			}

			if _, err := writer.Write(message); err != nil {
				logger.Error("websocketStatsPusher: ", err)
				return
			}

			if err := writer.Close(); err != nil {
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

func (w *websocketStatsPusher) connect(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, w.url.String(), nil)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#connect: %v", err)
	}

	w.conn = conn
	return nil
}

var expectedCloseMessages = []int{websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure}

// readPumpForControlMessages listens on the websocket connection with the sole
// intention of handling control messages, like server disconnect.
// https://stackoverflow.com/a/48181794/639773
// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L56
func (w *websocketStatsPusher) readPumpForControlMessages(cancel context.CancelFunc) {
	defer cancel()

	w.conn.SetReadLimit(maxMessageSize)
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
				logger.Warn(fmt.Sprintf("readPumpForControlMessages: %v", err))
			}
			return
		}
	}
}

func wrapConnErrorIf(err error) {
	if err != nil && websocket.IsUnexpectedCloseError(err, expectedCloseMessages...) {
		logger.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (w *websocketStatsPusher) Close() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.started {
		w.cancel()
	}
	w.started = false
	return nil
}

// Send sends data asynchronously across the websocket if it's open, or
// holds it in a small buffer until connection, throwing away messages
// once buffer is full.
func (w *websocketStatsPusher) Send(data []byte) {
	select {
	case w.send <- data:
	default:
	}
}
