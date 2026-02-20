package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const relayIncomingQueueSize = 64

var relayIDSeq uint64

type relayRegistration struct {
	ID            string
	Name          string
	RequestedPort int
	Listener      net.Listener
	Incoming      chan net.Conn
	Closed        chan struct{}
}

type openRelayRequest struct {
	Name          string `json:"name"`
	RequestedPort int    `json:"requestedPort"`
}

type openRelayResponse struct {
	RelayID      string `json:"relayId"`
	RequestedPort int   `json:"requestedPort"`
	BoundPort    int    `json:"boundPort"`
}

type closeRelayRequest struct {
	RelayID string `json:"relayId"`
}

var relayWSUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func (s *Server) openRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	var req openRelayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidRequestBody, fmt.Sprintf("invalid relay open request body: %v", err), nil)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "relay name is required", nil)
		return
	}
	if req.RequestedPort < 0 || req.RequestedPort > 65535 {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("invalid requestedPort %d", req.RequestedPort), nil)
		return
	}

	// Idempotent open: same name+port returns the existing relay.
	s.relayMu.Lock()
	for _, relay := range s.relays {
		if relay.Name == req.Name && relay.RequestedPort == req.RequestedPort {
			s.relayMu.Unlock()
			s.respondJSONAny(w, http.StatusOK, openRelayResponse{
				RelayID:      relay.ID,
				RequestedPort: relay.RequestedPort,
				BoundPort:    listenerPort(relay.Listener),
			})
			return
		}
	}
	s.relayMu.Unlock()

	listenAddr := fmt.Sprintf("0.0.0.0:%d", req.RequestedPort)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to open relay listener: %v", err), nil)
		return
	}

	relayID := fmt.Sprintf("relay-%x", atomic.AddUint64(&relayIDSeq, 1))
	reg := &relayRegistration{
		ID:            relayID,
		Name:          req.Name,
		RequestedPort: req.RequestedPort,
		Listener:      ln,
		Incoming:      make(chan net.Conn, relayIncomingQueueSize),
		Closed:        make(chan struct{}),
	}

	s.relayMu.Lock()
	s.relays[relayID] = reg
	s.relayMu.Unlock()

	go s.acceptRelayConnections(reg)

	s.respondJSONAny(w, http.StatusOK, openRelayResponse{
		RelayID:      relayID,
		RequestedPort: req.RequestedPort,
		BoundPort:    listenerPort(ln),
	})
}

func (s *Server) closeRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	var req closeRelayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidRequestBody, fmt.Sprintf("invalid relay close request body: %v", err), nil)
		return
	}
	relayID := strings.TrimSpace(req.RelayID)
	if relayID == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "relayId is required", nil)
		return
	}

	s.relayMu.Lock()
	relay, ok := s.relays[relayID]
	if ok {
		delete(s.relays, relayID)
	}
	s.relayMu.Unlock()

	if !ok {
		s.respondJSONAny(w, http.StatusOK, map[string]any{"relayId": relayID, "closed": false, "found": false})
		return
	}
	close(relay.Closed)
	_ = relay.Listener.Close()
	drainAndCloseIncoming(relay.Incoming)

	s.respondJSONAny(w, http.StatusOK, map[string]any{"relayId": relayID, "closed": true, "found": true})
}

func (s *Server) connectRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}
	relayID := strings.TrimSpace(r.URL.Query().Get("relayId"))
	if relayID == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "relayId query parameter is required", nil)
		return
	}

	s.relayMu.Lock()
	relay, ok := s.relays[relayID]
	s.relayMu.Unlock()
	if !ok {
		s.respondError(w, http.StatusNotFound, ErrCodeDeployFailed, fmt.Sprintf("relay not found: %s", relayID), nil)
		return
	}

	wsConn, err := relayWSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	var incoming net.Conn
	select {
	case incoming = <-relay.Incoming:
	case <-relay.Closed:
		_ = wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "relay closed"), time.Now().Add(2*time.Second))
		return
	case <-r.Context().Done():
		return
	}
	if incoming == nil {
		return
	}
	defer incoming.Close()

	_ = bridgeWebSocketAndTCP(wsConn, incoming)
}

func (s *Server) acceptRelayConnections(relay *relayRegistration) {
	for {
		conn, err := relay.Listener.Accept()
		if err != nil {
			select {
			case <-relay.Closed:
				return
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return
		}

		select {
		case relay.Incoming <- conn:
		default:
			_ = conn.Close()
		}
	}
}

func bridgeWebSocketAndTCP(ws *websocket.Conn, tcp net.Conn) error {
	errCh := make(chan error, 2)

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := tcp.Read(buf)
			if n > 0 {
				if wErr := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); wErr != nil {
					errCh <- wErr
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	go func() {
		for {
			msgType, payload, err := ws.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
				continue
			}
			if len(payload) == 0 {
				continue
			}
			if _, wErr := tcp.Write(payload); wErr != nil {
				errCh <- wErr
				return
			}
		}
	}()

	err := <-errCh
	if err == nil || errors.Is(err, io.EOF) || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return nil
	}
	return err
}

func drainAndCloseIncoming(ch chan net.Conn) {
	for {
		select {
		case conn := <-ch:
			if conn != nil {
				_ = conn.Close()
			}
		default:
			return
		}
	}
}

func listenerPort(ln net.Listener) int {
	if ln == nil {
		return 0
	}
	_, portRaw, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		return 0
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return 0
	}
	return port
}
