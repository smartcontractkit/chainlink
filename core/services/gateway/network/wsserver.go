package network

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type WebSocketServer interface {
	job.ServiceCtx

	// Not thread-safe. Can be called after Start() returns.
	GetPort() int
}

type WebSocketServerConfig struct {
	HTTPServerConfig
	HandshakeTimeoutMillis uint32
}

type webSocketServer struct {
	services.StateMachine
	config            *WebSocketServerConfig
	listener          net.Listener
	server            *http.Server
	acceptor          ConnectionAcceptor
	upgrader          *websocket.Upgrader
	doneCh            chan struct{}
	cancelBaseContext context.CancelFunc
	lggr              logger.Logger
}

func NewWebSocketServer(config *WebSocketServerConfig, acceptor ConnectionAcceptor, lggr logger.Logger) WebSocketServer {
	baseCtx, cancelBaseCtx := context.WithCancel(context.Background())
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: time.Duration(config.HandshakeTimeoutMillis) * time.Millisecond,
	}
	server := &webSocketServer{
		config:            config,
		acceptor:          acceptor,
		upgrader:          upgrader,
		doneCh:            make(chan struct{}),
		cancelBaseContext: cancelBaseCtx,
		lggr:              lggr.Named("WebSocketServer"),
	}
	mux := http.NewServeMux()
	mux.Handle(config.Path, http.HandlerFunc(server.handleRequest))
	server.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:           mux,
		BaseContext:       func(net.Listener) context.Context { return baseCtx },
		ReadTimeout:       time.Duration(config.ReadTimeoutMillis) * time.Millisecond,
		ReadHeaderTimeout: time.Duration(config.ReadTimeoutMillis) * time.Millisecond,
		WriteTimeout:      time.Duration(config.WriteTimeoutMillis) * time.Millisecond,
	}
	return server
}

func (s *webSocketServer) GetPort() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *webSocketServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(WsServerHandshakeAuthHeaderName)
	if len(authHeader) > HandshakeEncodedAuthHeaderMaxLen {
		s.lggr.Errorw("received auth header is too large", "len", len(authHeader))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authBytes, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		s.lggr.Error("received auth header can't be base64-decoded", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	attemptId, challenge, err := s.acceptor.StartHandshake(authBytes)
	if err != nil {
		s.lggr.Error("received invalid auth header", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	challengeStr := base64.StdEncoding.EncodeToString(challenge)
	hdr := make(http.Header)
	hdr.Add(WsServerHandshakeChallengeHeaderName, challengeStr)
	conn, err := s.upgrader.Upgrade(w, r, hdr)
	if err != nil {
		s.lggr.Error("failed websocket upgrade", err)
		conn.Close()
		s.acceptor.AbortHandshake(attemptId)
		return
	}

	msgType, response, err := conn.ReadMessage()
	if err != nil || msgType != websocket.BinaryMessage {
		s.lggr.Error("invalid handshake message", msgType, err)
		conn.Close()
		s.acceptor.AbortHandshake(attemptId)
		return
	}

	if err = s.acceptor.FinalizeHandshake(attemptId, response, conn); err != nil {
		s.lggr.Error("unable to finalize handshake", err)
		conn.Close()
		s.acceptor.AbortHandshake(attemptId)
		return
	}
}

func (s *webSocketServer) Start(ctx context.Context) error {
	return s.StartOnce("GatewayWebSocketServer", func() error {
		s.lggr.Info("starting gateway WebSocket server")
		return s.runServer()
	})
}

func (s *webSocketServer) Close() error {
	return s.StopOnce("GatewayWebSocketServer", func() (err error) {
		s.lggr.Info("closing gateway WebSocket server")
		s.cancelBaseContext()
		err = s.server.Shutdown(context.Background())
		<-s.doneCh
		return
	})
}

func (s *webSocketServer) runServer() (err error) {
	s.listener, err = net.Listen("tcp", s.server.Addr)
	if err != nil {
		return
	}
	tlsEnabled := s.config.TLSEnabled
	go func() {
		if tlsEnabled {
			err := s.server.ServeTLS(s.listener, s.config.TLSCertPath, s.config.TLSKeyPath)
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway WS server closed with error:", err)
			}
		} else {
			err := s.server.Serve(s.listener)
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway WS server closed with error:", err)
			}
		}
		s.doneCh <- struct{}{}
	}()
	return
}
