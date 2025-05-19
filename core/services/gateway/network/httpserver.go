package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
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
	CORSEnabled          bool
	CORSAllowedOrigins   []string
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

// split URL into: scheme, hostname, port
func (s *httpServer) splitURL(rawURL string) (string, string, string, error) {
	// lowercase the URL to avoid case sensitivity issues
	parsedURL, err := url.Parse(strings.ToLower((rawURL)))
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing URL: %w", err)
	}

	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		// if there's no port, the host itself is returned
		if parsedURL.Host != "" {
			return parsedURL.Scheme, parsedURL.Host, "", nil
		}
		return "", "", "", fmt.Errorf("error splitting host and port: %w", err)
	}

	return parsedURL.Scheme, host, port, nil
}

func (s *httpServer) isAllowedOrigin(origin string) bool {
	originScheme, originHost, originPort, err := s.splitURL(origin)
	if err != nil {
		s.lggr.Debug("error parsing origin URL", err)
		return false
	}
	for _, allowed := range s.config.CORSAllowedOrigins {
		// probably better to do this once when server starts and store it in a map
		// this is an easier solution so we don't have to apply more changes to the code
		// just need to be careful when specifying allowed origins in the config file
		allowedScheme, allowedHost, allowedPort, err := s.splitURL(allowed)
		if err != nil {
			s.lggr.Debug("error parsing allowed origin URL", err)
			continue
		}
		// skip if the scheme doesn't match at all
		if originScheme != allowedScheme {
			continue
		}
		// skip if the port doesn't match at all
		if originPort != allowedPort {
			continue
		}
		// check for exact host match (e.g., remix.com)
		if originHost == allowedHost {
			return true
		}
		// check for wildcard host match (e.g., *.remix.com)
		if strings.HasPrefix(allowedHost, "*.") {
			allowedHost = allowedHost[2:]
			if strings.HasSuffix(originHost, allowedHost) {
				return true
			}
		}
	}
	return false
}

func (s *httpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	if s.config.CORSEnabled {
		origin := r.Header.Get("Origin")
		if s.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		// handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

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
