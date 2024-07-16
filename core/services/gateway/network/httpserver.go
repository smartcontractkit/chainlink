package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type HttpServer interface {
	job.ServiceCtx

	// Not thread-safe. Should be called once, before Start() is called.
	SetHTTPRequestHandler(handler HTTPRequestHandler)

	// Not thread-safe. Can be called after Start() returns.
	GetPort() int
}

type HTTPRequestHandler interface {
	ProcessRequest(ctx context.Context, rawRequest []byte) (rawResponse []byte, httpStatusCode int)
}

type HTTPServerConfig struct {
	Host                 string
	Port                 uint16
	TLSEnabled           bool
	TLSCertPath          string
	TLSKeyPath           string
	Path                 string
	ContentTypeHeader    string
	ReadTimeoutMillis    uint32
	WriteTimeoutMillis   uint32
	RequestTimeoutMillis uint32
	MaxRequestBytes      int64
}

type httpServer struct {
	services.StateMachine
	config            *HTTPServerConfig
	listener          net.Listener
	server            *http.Server
	handler           HTTPRequestHandler
	doneCh            chan struct{}
	cancelBaseContext context.CancelFunc
	lggr              logger.Logger
}

const (
	HealthCheckPath     = "/health"
	HealthCheckResponse = "OK"
)

func NewHttpServer(config *HTTPServerConfig, lggr logger.Logger) HttpServer {
	baseCtx, cancelBaseCtx := context.WithCancel(context.Background())
	server := &httpServer{
		config:            config,
		doneCh:            make(chan struct{}),
		cancelBaseContext: cancelBaseCtx,
		lggr:              lggr.Named("WebSocketServer"),
	}
	mux := http.NewServeMux()
	mux.Handle(config.Path, http.HandlerFunc(server.handleRequest))
	mux.Handle(HealthCheckPath, http.HandlerFunc(server.handleHealthCheck))
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

func (s *httpServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(HealthCheckResponse))
	if err != nil {
		s.lggr.Debug("error when writing response for healthcheck", err)
	}
}

func (s *httpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	source := http.MaxBytesReader(nil, r.Body, s.config.MaxRequestBytes)
	rawMessage, err := io.ReadAll(source)
	if err != nil {
		s.lggr.Error("error reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestCtx := r.Context()
	if s.config.RequestTimeoutMillis > 0 {
		var cancel context.CancelFunc
		requestCtx, cancel = context.WithTimeout(requestCtx, time.Duration(s.config.RequestTimeoutMillis)*time.Millisecond)
		defer cancel()
	}
	rawResponse, httpStatusCode := s.handler.ProcessRequest(requestCtx, rawMessage)

	w.Header().Set("Content-Type", s.config.ContentTypeHeader)
	w.WriteHeader(httpStatusCode)
	_, err = w.Write(rawResponse)
	if err != nil {
		s.lggr.Error("error when writing response", err)
	}
}

func (s *httpServer) SetHTTPRequestHandler(handler HTTPRequestHandler) {
	s.handler = handler
}

func (s *httpServer) GetPort() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *httpServer) Start(ctx context.Context) error {
	return s.StartOnce("GatewayHTTPServer", func() error {
		s.lggr.Info("starting gateway HTTP server")
		return s.runServer()
	})
}

func (s *httpServer) Close() error {
	return s.StopOnce("GatewayHTTPServer", func() (err error) {
		s.lggr.Info("closing gateway HTTP server")
		s.cancelBaseContext()
		err = s.server.Shutdown(context.Background())
		<-s.doneCh
		return
	})
}

func (s *httpServer) runServer() (err error) {
	s.listener, err = net.Listen("tcp", s.server.Addr)
	if err != nil {
		return
	}
	tlsEnabled := s.config.TLSEnabled

	go func() {
		if tlsEnabled {
			err := s.server.ServeTLS(s.listener, s.config.TLSCertPath, s.config.TLSKeyPath)
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway server closed with error:", err)
			}
		} else {
			err := s.server.Serve(s.listener)
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway server closed with error:", err)
			}
		}
		s.doneCh <- struct{}{}
	}()
	return
}
