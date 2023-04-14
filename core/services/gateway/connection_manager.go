package gateway

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectionManager interface {
	AddConnection(address string, conn *websocket.Conn)
	SendToAll(message *Message)
}

type connManager struct {
	conns   map[string]*connection
	connMu  sync.Mutex
	handler Handler
}

type connection struct {
	wsConn *websocket.Conn
	state  string
	done   chan bool
}

func NewConnectionManager(donConfig *GatewayDONConfig, handler Handler) ConnectionManager {
	mgr := connManager{handler: handler, conns: make(map[string]*connection)}
	/*for _, member := range donConfig.Members {
		mgr.conns[member.SignerAddress] = &connection{done: make(chan bool)}
	}*/
	return &mgr
}

func (m *connManager) AddConnection(address string, conn *websocket.Conn) {
	// TODO check if exists, only allow one
	m.connMu.Lock()
	defer m.connMu.Unlock()
	done := make(chan bool)
	m.conns[address] = &connection{wsConn: conn, state: "connected", done: make(chan bool)}
	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			msg, err := Decode(message)
			if err != nil {
				break
			}
			m.handler.HandleNodeMessage(msg, address)
		}
		done <- true
	}()
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		hbMsg, _ := Encode(&Message{Method: "heartbeat"})
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				conn.WriteMessage(websocket.TextMessage, hbMsg)
			}
		}
	}()
}

func (m *connManager) SendToAll(message *Message) {
	msgBytes, _ := Encode(message)
	m.connMu.Lock()
	defer m.connMu.Unlock()
	for _, conn := range m.conns {
		conn.wsConn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
