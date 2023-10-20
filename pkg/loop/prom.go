package loop

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

type PromServerOpts struct {
	Handler http.Handler
}

func NewPromServer(port int, lggr logger.Logger) *PromServer {
	return PromServerOpts{}.New(port, lggr)
}

func (o PromServerOpts) New(port int, lggr logger.Logger) *PromServer {
	s := &PromServer{
		port:     port,
		lggr:     lggr,
		srvrDone: make(chan struct{}),
		srvr: &http.Server{
			// reasonable default based on typical prom poll interval of 15s.
			ReadTimeout: 5 * time.Second,
		},

		handler: o.Handler,
	}
	if s.handler == nil {
		s.handler = promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		)
	}
	return s
}

// Start starts HTTP server on specified port to handle metrics requests
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

// Close shuts down the underlying HTTP server. See [http.Server.Close] for details
func (p *PromServer) Close() error {
	err := p.srvr.Shutdown(context.Background())
	<-p.srvrDone
	return err
}

// Name of the server
func (p *PromServer) Name() string {
	return fmt.Sprintf("%s-prom-server", p.lggr.Name())
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
