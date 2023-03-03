package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebsocketClient implements the ClientTransport interface with websockets.
type WebsocketClient struct {
	ctx context.Context

	// Config
	writeTimeout time.Duration

	// Underlying communication channel
	conn *websocket.Conn

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

// newWebsocketClient establishes the transport with the required ConnectOptions
// and returns it to the caller.
func newWebsocketClient(ctx context.Context, addr string, opts ConnectOptions, onClose func()) (_ *WebsocketClient, err error) {
	writeTimeout := defaultWriteTimeout
	if opts.WriteTimeout != 0 {
		writeTimeout = opts.WriteTimeout
	}

	d := websocket.Dialer{
		TLSClientConfig:  opts.TransportCredentials.Config,
		HandshakeTimeout: 45 * time.Second,
	}

	url := fmt.Sprintf("wss://%s", addr)
	conn, _, err := d.DialContext(ctx, url, http.Header{})
	if err != nil {
		return nil, fmt.Errorf("[wsrpc] error while dialing %w", err)
	}

	c := &WebsocketClient{
		ctx:          ctx,
		writeTimeout: writeTimeout,
		conn:         conn,
		onClose:      onClose,
		write:        make(chan []byte), // Should this be buffered?
		read:         make(chan []byte), // Should this be buffered?
		done:         make(chan struct{}),
		interrupt:    make(chan struct{}),
	}

	// Start go routines to establish the read/write channels
	go c.start()

	return c, nil
}

// Read returns a channel which provides the messages as they are read.
func (c *WebsocketClient) Read() <-chan []byte {
	return c.read
}

// Write writes a message the websocket connection.
func (c *WebsocketClient) Write(msg []byte) error {
	c.write <- msg

	return nil
}

// Close closes the websocket connection and cleans up pump goroutines.
func (c *WebsocketClient) Close() error {
	close(c.interrupt)

	return nil
}

// start run readPump in a goroutine and waits on writePump.
func (c WebsocketClient) start() {
	defer c.onClose()

	// Set up reader
	go c.readPump()

	c.writePump()
}

// readPump pumps messages from the websocket connection. When a websocket
// connection closure is detected through a read error, it closes the done
// channel to shutdown writePump.
//
// The application runs readPump in a per-connection goroutine. This ensures
// that there is at most one reader on a connection by executing all reads from
// this goroutine.
func (c *WebsocketClient) readPump() {
	defer close(c.done)

	//nolint:errcheck
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(handlePong(c.conn))

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[wsrpc] Read error: ", err)

			return
		}

		c.read <- msg
	}
}

// writePump pumps messages from the client to the websocket connection.
//
// A goroutine running writePump is started for each connection. This ensures
// that there is at most one writer to a connection by executing all writes
// from this goroutine.
func (c *WebsocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			// When the read detects a websocket closure, it will close the done
			// channel so we can exit
			return
		case msg := <-c.write: // Write the message
			// Any error due to a closed connection will be immediately picked
			// up in the subsequent network message read or write.
			//nolint:errcheck
			c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Printf("[wsrpc] write error: %v\n", err)

				c.conn.Close()

				return
			}
		case <-ticker.C:
			// Any error due to a closed connection will be immediately picked
			// up in the subsequent network message read or write.
			if err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(c.writeTimeout)); err != nil {
				c.conn.Close()

				return
			}
		case <-c.interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				return
			}
			c.conn.Close()
			select {
			case <-c.done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}
