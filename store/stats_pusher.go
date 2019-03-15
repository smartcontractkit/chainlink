package store

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smartcontractkit/chainlink/logger"
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
	url     url.URL
	conn    *websocket.Conn
	send    chan []byte
	boot    *sync.Mutex
	started bool
}

// NewWebsocketStatsPusher returns a stats pusher using a websocket for
// delivery.
func NewWebsocketStatsPusher(url *url.URL) StatsPusher {
	return &websocketStatsPusher{
		url:  *url,
		send: make(chan []byte),
		boot: &sync.Mutex{},
	}
}

// Start starts a write pump over a websocket.
func (w *websocketStatsPusher) Start() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.started {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(w.url.String(), nil)
	if err != nil {
		return fmt.Errorf("websocketStatsPusher#Start(): %v", err)
	}

	w.conn = conn
	go w.writePump()
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
)

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
func (w *websocketStatsPusher) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wrapLoggerErrorIf(w.conn.Close())
	}()
	for {
		select {
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
			// batching sending for efficiency.
			n := len(w.send)
			for i := 0; i < n; i++ {
				additionalMsg, open := <-w.send
				if !open {
					break
				}
				_, err = writer.Write(additionalMsg)
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

func wrapLoggerErrorIf(err error) {
	if err != nil && !websocket.IsCloseError(err) {
		logger.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (w *websocketStatsPusher) Close() error {
	w.boot.Lock()
	defer w.boot.Unlock()

	if w.send != nil {
		close(w.send)
		w.send = nil
	}
	w.started = false
	return nil
}
