package synchronization

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
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
	services.ServiceCtx
	Url() url.URL
	Status() ConnectionStatus
	Send(context.Context, []byte, ...int)
	Receive(context.Context, ...time.Duration) ([]byte, error)
}

type NoopExplorerClient struct{}

// Url always returns underlying url.
func (NoopExplorerClient) Url() url.URL { return url.URL{} }

// Status always returns ConnectionStatusDisconnected.
func (NoopExplorerClient) Status() ConnectionStatus { return ConnectionStatusDisconnected }

// Start is a no-op
func (NoopExplorerClient) Start(context.Context) error { return nil }

// Close is a no-op
func (NoopExplorerClient) Close() error { return nil }

// Healthy is a no-op
func (NoopExplorerClient) Healthy() error { return nil }

// Ready is a no-op
func (NoopExplorerClient) Ready() error { return nil }

// Send is a no-op
func (NoopExplorerClient) Send(context.Context, []byte, ...int) {}

// Receive is a no-op
func (NoopExplorerClient) Receive(context.Context, ...time.Duration) ([]byte, error) { return nil, nil }

type explorerClient struct {
	utils.StartStopOnce
	conn             *websocket.Conn
	sendText         chan []byte
	sendBinary       chan []byte
	dropMessageCount atomic.Uint32
	receive          chan []byte
	sleeper          utils.Sleeper
	status           ConnectionStatus
	url              *url.URL
	accessKey        string
	secret           string
	lggr             logger.Logger

	chStop        chan struct{}
	wg            sync.WaitGroup
	writePumpDone chan struct{}

	statusMtx sync.RWMutex
}

// NewExplorerClient returns a stats pusher using a websocket for
// delivery.
func NewExplorerClient(url *url.URL, accessKey, secret string, lggr logger.Logger) ExplorerClient {
	return &explorerClient{
		url:       url,
		receive:   make(chan []byte),
		sleeper:   utils.NewBackoffSleeper(),
		status:    ConnectionStatusDisconnected,
		accessKey: accessKey,
		secret:    secret,
		lggr:      lggr.Named("ExplorerClient"),

		sendText:   make(chan []byte, SendBufferSize),
		sendBinary: make(chan []byte, SendBufferSize),
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
func (ec *explorerClient) Start(context.Context) error {
	return ec.StartOnce("Explorer client", func() error {
		ec.chStop = make(chan struct{})
		ec.wg.Add(1)
		go ec.connectAndWritePump()
		return nil
	})
}

// Send sends data asynchronously across the websocket if it's open, or
// holds it in a small buffer until connection, throwing away messages
// once buffer is full.
// func (ec *explorerClient) Receive(durationParams ...time.Duration) ([]byte, error) {
func (ec *explorerClient) Send(ctx context.Context, data []byte, messageTypes ...int) {
	messageType := ExplorerTextMessage
	if len(messageTypes) > 0 {
		messageType = messageTypes[0]
	}
	var send chan []byte
	switch messageType {
	case ExplorerTextMessage:
		send = ec.sendText
	case ExplorerBinaryMessage:
		send = ec.sendBinary
	default:
		log.Panicf("send on explorer client received unsupported message type %d", messageType)
	}
	select {
	case send <- data:
		ec.dropMessageCount.Store(0)
	case <-ctx.Done():
		return
	default:
		ec.logBufferFullWithExpBackoff(data)
	}
}

// logBufferFullWithExpBackoff logs messages at
// 1
// 2
// 4
// 8
// 16
// 32
// 64
// 100
// 200
// 300
// etc...
func (ec *explorerClient) logBufferFullWithExpBackoff(data []byte) {
	count := ec.dropMessageCount.Inc()
	if count > 0 && (count%100 == 0 || count&(count-1) == 0) {
		ec.lggr.Warnw("explorer client buffer full, dropping message", "data", data, "droppedCount", count)
	}
}

// Receive blocks the caller while waiting for a response from the server,
// returning the raw response bytes
func (ec *explorerClient) Receive(ctx context.Context, durationParams ...time.Duration) ([]byte, error) {
	duration := defaultReceiveTimeout
	if len(durationParams) > 0 {
		duration = durationParams[0]
	}

	select {
	case data := <-ec.receive:
		return data, nil
	case <-time.After(duration):
		return nil, ErrReceiveTimeout
	case <-ctx.Done():
		return nil, nil
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
func (ec *explorerClient) connectAndWritePump() {
	defer ec.wg.Done()
	ctx, cancel := utils.ContextFromChan(ec.chStop)
	defer cancel()
	for {
		select {
		case <-time.After(ec.sleeper.After()):
			ec.lggr.Infow("Connecting to explorer", "url", ec.url)
			err := ec.connect(ctx)
			if ctx.Err() != nil {
				return
			} else if err != nil {
				ec.setStatus(ConnectionStatusError)
				ec.lggr.Warn("Failed to connect to explorer (", ec.url.String(), "): ", err)
				break
			}

			ec.setStatus(ConnectionStatusConnected)

			ec.lggr.Infow("Connected to explorer", "url", ec.url)
			ec.sleeper.Reset()
			ec.writePumpDone = make(chan struct{})
			ec.wg.Add(1)
			go ec.readPump()
			ec.writePump()

		case <-ec.chStop:
			return
		}
	}
}

func (ec *explorerClient) setStatus(s ConnectionStatus) {
	ec.statusMtx.Lock()
	defer ec.statusMtx.Unlock()
	ec.status = s
}

// Inspired by https://github.com/gorilla/websocket/blob/master/examples/chat/client.go#L82
func (ec *explorerClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ec.wrapConnErrorIf(ec.conn.Close()) // exclusive responsibility to close ws conn
	}()

	for {
		select {
		case message, open := <-ec.sendText:
			if !open {
				ec.wrapConnErrorIf(ec.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			err := ec.writeMessage(message, websocket.TextMessage)
			if err != nil {
				ec.lggr.Warnw("websocketStatsPusher: error writing text message", "err", err)
				return
			}

		case message, open := <-ec.sendBinary:
			if !open {
				ec.wrapConnErrorIf(ec.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			err := ec.writeMessage(message, websocket.BinaryMessage)
			if err != nil {
				ec.lggr.Warnw("websocketStatsPusher: error writing binary message", "err", err)
				return
			}

		case <-ticker.C:
			ec.wrapConnErrorIf(ec.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := ec.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ec.wrapConnErrorIf(err)
				return
			}

		case <-ec.writePumpDone:
			return
		case <-ec.chStop:
			return
		}
	}
}

func (ec *explorerClient) writeMessage(message []byte, messageType int) error {
	ec.wrapConnErrorIf(ec.conn.SetWriteDeadline(time.Now().Add(writeWait)))
	writer, err := ec.conn.NextWriter(messageType)
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		return err
	}
	ec.lggr.Tracew("websocketStatsPusher successfully wrote message", "messageType", messageType, "message", message)

	return writer.Close()
}

func (ec *explorerClient) connect(ctx context.Context) error {
	authHeader := http.Header{}
	authHeader.Add("X-Explore-Chainlink-Accesskey", ec.accessKey)
	authHeader.Add("X-Explore-Chainlink-Secret", ec.secret)
	authHeader.Add("X-Explore-Chainlink-Core-Version", static.Version)
	authHeader.Add("X-Explore-Chainlink-Core-Sha", static.Sha)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, ec.url.String(), authHeader)
	if ctx.Err() != nil {
		return fmt.Errorf("websocketStatsPusher#connect context canceled: %w", ctx.Err())
	} else if err != nil {
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
func (ec *explorerClient) readPump() {
	defer ec.wg.Done()
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
				ec.lggr.Warnw("Unexpected close error on ExplorerClient", "err", err)
			}
			close(ec.writePumpDone)
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
		ec.lggr.Error(fmt.Sprintf("websocketStatsPusher: %v", err))
	}
}

func (ec *explorerClient) Close() error {
	return ec.StopOnce("Explorer client", func() error {
		close(ec.chStop)
		ec.wg.Wait()
		return nil
	})
}
