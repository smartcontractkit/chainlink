package plugins

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

type PromServer struct {
	port        int
	srvrDone    chan struct{} // closed when the http server is done
	srvr        *http.Server
	tcpListener *net.TCPListener
	lggr        logger.Logger

	handler http.Handler
}

type PromServerOpt func(*PromServer)

func WithHandler(h http.Handler) PromServerOpt {
	return func(s *PromServer) {
		s.handler = h
	}
}

func NewPromServer(port int, lggr logger.Logger, opts ...PromServerOpt) *PromServer {

	s := &PromServer{
		port:     port,
		lggr:     lggr,
		srvrDone: make(chan struct{}),
		srvr: &http.Server{
			// reasonable default based on typical prom poll interval of 15s.
			ReadTimeout: 5 * time.Second,
		},

		handler: promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start start HTTP server on specified port to handle metrics requests
func (p *PromServer) Start() error {
	p.lggr.Debugf("Starting prom server on port %d", p.port)
	err := p.setupListener()
	if err != nil {
		return err
	}

	http.Handle("/metrics", p.handler)

	go func() {
		defer close(p.srvrDone)
		err := p.srvr.Serve(p.tcpListener)
		if errors.Is(err, net.ErrClosed) {
			// ErrClose is expected on gracefully shutdown
			p.lggr.Warnf("%s closed", p.Name())
		} else {
			p.lggr.Errorf("%s: %s", p.Name(), err)
		}

	}()
	return nil
}

// Close shutdowns down the underlying HTTP server. See [http.Server.Close] for details
func (p *PromServer) Close() error {
	err := p.srvr.Shutdown(context.Background())
	<-p.srvrDone
	return err
}

// Name of the server
func (p *PromServer) Name() string {
	return fmt.Sprintf("%s-prom-server", p.lggr.Name())
}

// Port is the resolved port and is only known after Start().
// returns -1 before it is resolved or if there was an error during resolution.
func (p *PromServer) Port() int {
	if p.tcpListener == nil {
		return -1
	}
	// always safe to cast because we explicitly have a tcp listener
	// there is direct access to Port without the addr casting
	// Note: addr `:0` is not resolved to non-zero port until ListenTCP is called
	// net.ResolveTCPAddr sounds promising, but doesn't work in practice
	return p.tcpListener.Addr().(*net.TCPAddr).Port

}

// setupListener creates explicit listener so that we can resolve `:0` port, which is needed for testing
// if we didn't need the resolved addr, or could pick a static port we could use p.srvr.ListenAndServer
func (p *PromServer) setupListener() error {

	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: p.port,
	})
	if err != nil {
		return err
	}

	p.tcpListener = l
	return nil
}
