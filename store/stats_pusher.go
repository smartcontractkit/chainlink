package store

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// StatsPusher encapsulates all the functionality needed to
// push run information to linkstats.
type StatsPusher interface {
	Start() error
	Close() error
}

// NewStatsPusher returns a functioning instance depending on the
// URL passed: nil is a noop instance, url assumes a websocket instance.
// No support for http.
func NewStatsPusher(url *url.URL) StatsPusher {
	if url != nil {
		return NewWebsocketStatsPusher(url)
	}
	return noopStatsPusher{}
}

type noopStatsPusher struct{}

func (noopStatsPusher) Start() error { return nil }
func (noopStatsPusher) Close() error { return nil }

type websocketStatsPusher struct {
	boot    *sync.Mutex
	conn    *websocket.Conn
	done    chan struct{}
	send    chan []byte
	sleeper utils.Sleeper
	started bool
	url     url.URL
}

// NewWebsocketStatsPusher returns a stats pusher using a websocket for
// delivery.
func NewWebsocketStatsPusher(url *url.URL) StatsPusher {
	return &websocketStatsPusher{
		url:     *url,
		send:    make(chan []byte),
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

	w.done = make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go w.connectAndWritePump(wg, w.done)
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
func (w *websocketStatsPusher) connectAndWritePump(wg *sync.WaitGroup, done chan struct{}) {
	wg.Done()

	for {
		select {
		case <-done:
			return
		case <-time.After(w.sleeper.After()):
			if err := w.connect(); err != nil {
				logger.Warn("Inability to connect to linkstats: ", err)
				break
			}

			w.sleeper.Reset()

			serverDone := make(chan struct{})
			go w.readPumpForControlMessages(serverDone)
			w.writePump(done, serverDone)
		}
	}
}

var (
	newline = []byte{'\n'}
)

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (w *websocketStatsPusher) writePump(done chan struct{}, serverDone chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wrapLoggerErrorIf(w.conn.Close()) // exclusive responsibility to close ws conn
	}()
	for {
		select {
		case <-done:
			return
		case <-serverDone:
			return
		case message, open := <-w.send:
			if !open { // channel closed
				wrapLoggerErrorIf(w.conn.WriteMessage(websocket.CloseMessage, []byte{}))
				return
			}

			wrapLoggerErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error(err)
				return
			}
			_, err = writer.Write(message)
			wrapLoggerErrorIf(err)

			// Add queued messages to the current websocket message,
			// batch sending for efficiency.
			n := len(w.send)
			for i := 0; i < n; i++ {
				additionalMsg, open := <-w.send
				if !open {
					break
				}
				err = multierr.Append(
					utils.JustError(writer.Write(newline)),
					utils.JustError(writer.Write(additionalMsg)),
				)
				wrapLoggerErrorIf(err)
			}

			if err := writer.Close(); err != nil {
				logger.Error(err)
				return
			}
		case <-ticker.C:
			wrapLoggerErrorIf(w.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				wrapLoggerErrorIf(err)
				return
			}
		}
	}
}

func (w *websocketStatsPusher) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(w.url.String(), nil)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#connect: %v", err)
	}

	w.conn = conn
	return nil
}

// readPumpForControlMessages listens on the websocket connection with the sole
// intention of handling control messages, like server disconnect.
// https://stackoverflow.com/a/48181794/639773
// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L56
func (w *websocketStatsPusher) readPumpForControlMessages(serverDone chan struct{}) {
	w.conn.SetReadLimit(maxMessageSize)
	logger.WarnIf(w.conn.SetReadDeadline(time.Now().Add(pongWait)))
	w.conn.SetPongHandler(func(string) error {
		logger.WarnIf(w.conn.SetReadDeadline(time.Now().Add(pongWait)))
		return nil
	})

	for {
		_, _, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn(fmt.Sprintf("readPumpForControlMessages: %v", err))
			}
			close(serverDone)
			break
		}
	}
}

func wrapLoggerErrorIf(err error) {
	if err != nil && !websocket.IsCloseError(err) {
		logger.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (w *websocketStatsPusher) Close() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.done != nil {
		close(w.done)
		w.done = nil
	}
	w.started = false
	return nil
}
