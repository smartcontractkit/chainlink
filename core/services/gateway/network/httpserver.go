package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name HttpServer --output ./mocks/ --case=underscore
type HttpServer interface {
	job.ServiceCtx

	// Not thread-safe. Should be done once before calling Start().
	SetHTTPRequestHandler(handler HTTPRequestHandler)
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
}

type httpServer struct {
	utils.StartStopOnce
	config            *HTTPServerConfig
	server            *http.Server
	handler           HTTPRequestHandler
	doneCh            chan struct{}
	cancelBaseContext context.CancelFunc
	lggr              logger.Logger
}

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

func (s *httpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	rawMessage, err := io.ReadAll(r.Body)
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

func (s *httpServer) Start(ctx context.Context) error {
	return s.StartOnce("GatewayHTTPServer", func() error {
		s.lggr.Info("starting gateway HTTP server")
		s.runServer()
		return nil
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

func (s *httpServer) runServer() {
	tlsEnabled := s.config.TLSEnabled
	go func() {
		if tlsEnabled {
			err := s.server.ListenAndServeTLS(s.config.TLSCertPath, s.config.TLSKeyPath)
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway server closed with error:", err)
			}
		} else {
			err := s.server.ListenAndServe()
			if err != http.ErrServerClosed {
				s.lggr.Error("gateway server closed with error:", err)
			}
		}
		s.doneCh <- struct{}{}
	}()
}
