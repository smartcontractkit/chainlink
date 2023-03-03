package transport

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	mu sync.Mutex

	// config
	writeTimeout time.Duration

	// Underlying communication channel
	conn *websocket.Conn

	// The current state of the server transport
	state transportState

	// Callback function called when the transport is closed
	onClose func()

	// Communication channels
	write chan []byte
	read  chan []byte

	// A signal channel called when the reader encounters a websocket close error
	done chan struct{}
	// A signal channel called when the transport is closed
	interrupt chan struct{}
}

// newWebsocketServer server upgrades an HTTP connection to a websocket connection.
func newWebsocketServer(c *websocket.Conn, config *ServerConfig, onClose func()) *WebsocketServer {
	writeTimeout := defaultWriteTimeout
	if config.WriteTimeout != 0 {
		writeTimeout = config.WriteTimeout
	}

	s := &WebsocketServer{
		writeTimeout: writeTimeout,
		conn:         c,
		onClose:      onClose,
		write:        make(chan []byte),
		read:         make(chan []byte),
		done:         make(chan struct{}),
		interrupt:    make(chan struct{}),
	}

	go s.start()

	return s
}

// Read returns a channel which provides the messages as they are read.
func (s *WebsocketServer) Read() <-chan []byte {
	return s.read
}

// Write writes a message the websocket connection.
func (s *WebsocketServer) Write(msg []byte) error {
	// Send the message to the channel
	s.write <- msg

	return nil
}

// Close closes the websocket connection and cleans up pump goroutines. Notifies
// the caller with the onClose callback.
func (s *WebsocketServer) Close() error {
	s.mu.Lock()
	// Make sure we only Close once.
	if s.state == closing {
		s.mu.Unlock()

		return nil
	}

	s.state = closing

	// Close the write channel to stop the go routine
	close(s.interrupt)

	s.mu.Unlock()

	return nil
}

// start runs readPump in a goroutine and waits on writePump.
func (s *WebsocketServer) start() {
	defer func() {
		s.Close()
		s.onClose()
	}()

	// Set up reader
	go s.readPump()
	s.writePump()
}

// readPump pumps messages from the websocket connection.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s *WebsocketServer) readPump() {
	defer close(s.done)

	//nolint:errcheck
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(handlePong(s.conn))

	for {
		_, msg, err := s.conn.ReadMessage()
		// An error is provided when the websocket connection is closed,
		// allowing us to clean up the goroutine.
		if err != nil {
			log.Println("[wsrpc] Read error: ", err)

			break
		}

		s.read <- msg
	}
}

// writePump pumps messages from the server to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// server ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s *WebsocketServer) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			// When the read detects a websocket closure, it will close the done
			// channel so we can exit.
			return
		case msg := <-s.write:
			// Any error due to a closed connection will be immediately picked
			// up in the subsequent network message read or write.
			//nolint:errcheck
			s.conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
			if err := s.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Printf("[wsrpc] write error: %v\n", err)

				return
			}
		case <-ticker.C:
			// Any error due to a closed connection will be immediately picked
			// up in the subsequent network message read or write.
			if err := s.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(s.writeTimeout)); err != nil {
				s.conn.Close()

				return
			}
		case <-s.interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := s.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				return
			}
			s.conn.Close()
			select {
			case <-s.done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}
