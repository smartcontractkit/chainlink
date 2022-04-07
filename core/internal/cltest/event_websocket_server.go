package cltest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

// EventWebSocketServer is a web socket server designed specifically for testing
type EventWebSocketServer struct {
	*httptest.Server
	mutex          *sync.RWMutex // shared mutex for safe access to arrays/maps.
	t              *testing.T
	connections    []*websocket.Conn
	Connected      chan struct{}
	Disconnected   chan struct{}
	ReceivedText   chan string
	ReceivedBinary chan []byte
	URL            *url.URL
}

// NewEventWebSocketServer returns a new EventWebSocketServer
func NewEventWebSocketServer(t *testing.T) (*EventWebSocketServer, func()) {
	server := &EventWebSocketServer{
		mutex:          &sync.RWMutex{},
		t:              t,
		Connected:      make(chan struct{}, 1), // have buffer of one for easier assertions after the event
		Disconnected:   make(chan struct{}, 1), // have buffer of one for easier assertions after the event
		ReceivedText:   make(chan string, 100),
		ReceivedBinary: make(chan []byte, 100),
	}

	server.Server = httptest.NewServer(http.HandlerFunc(server.handler))
	u, err := url.Parse(server.Server.URL)
	if err != nil {
		t.Fatal("EventWebSocketServer: ", err)
	}
	u.Scheme = "ws"
	server.URL = u
	return server, func() {
		server.Close()
	}
}

func (wss EventWebSocketServer) ConnectionsCount() int {
	wss.mutex.RLock()
	defer wss.mutex.RUnlock()
	return len(wss.connections)
}

// Broadcast sends a message to every web socket client connected to the EventWebSocketServer
func (wss *EventWebSocketServer) Broadcast(message string) error {
	wss.mutex.RLock()
	defer wss.mutex.RUnlock()
	for _, connection := range wss.connections {
		err := connection.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return errors.Wrap(err, "error writing message to connection")
		}
	}

	return nil
}

// WriteCloseMessage tells connected clients to disconnect.
// Useful to emulate that the websocket server is shutting down without
// actually shutting down.
// This overcomes httptest.Server's inability to restart on the same URL:port.
func (wss *EventWebSocketServer) WriteCloseMessage() {
	wss.mutex.RLock()
	for _, connection := range wss.connections {
		err := connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			wss.t.Error(err)
		}
	}
	wss.mutex.RUnlock()
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	closeCodes = []int{websocket.CloseNormalClosure, websocket.CloseAbnormalClosure}
)

func (wss *EventWebSocketServer) handler(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wss.t.Fatal("EventWebSocketServer Upgrade: ", err)
	}

	wss.addConnection(conn)
	for {
		messageType, payload, err := conn.ReadMessage() // we only read
		if websocket.IsCloseError(err, closeCodes...) {
			wss.removeConnection(conn)
			return
		}
		if err != nil {
			wss.t.Fatal("EventWebSocketServer ReadMessage: ", err)
		}

		if messageType == websocket.TextMessage {
			select {
			case wss.ReceivedText <- string(payload):
			default:
			}
		} else if messageType == websocket.BinaryMessage {
			select {
			case wss.ReceivedBinary <- payload:
			default:
			}
		} else {
			wss.t.Fatal("EventWebSocketServer UnsupportedMessageType: ", messageType)
		}
	}
}

func (wss *EventWebSocketServer) addConnection(conn *websocket.Conn) {
	wss.mutex.Lock()
	wss.connections = append(wss.connections, conn)
	wss.mutex.Unlock()
	select { // broadcast connected event
	case wss.Connected <- struct{}{}:
	default:
	}
}

func (wss *EventWebSocketServer) removeConnection(conn *websocket.Conn) {
	newc := []*websocket.Conn{}
	wss.mutex.Lock()
	for _, connection := range wss.connections {
		if connection != conn {
			newc = append(newc, connection)
		}
	}
	wss.connections = newc
	wss.mutex.Unlock()
	select { // broadcast disconnected event
	case wss.Disconnected <- struct{}{}:
	default:
	}
}
