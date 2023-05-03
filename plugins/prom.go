package plugins

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type PromServer struct {
	port     int
	srvr     *http.Server
	listener net.Listener
	lggr     logger.Logger
}

func NewPromServer(port int, lggr logger.Logger) *PromServer {

	return &PromServer{
		port: port,
		lggr: lggr,
		srvr: &http.Server{}, // Do not configure handler explicitly; want DefaultServerMux here and in Handle

	}
}

// Start start HTTP server on specified port to handle metrics requests
func (p *PromServer) Start() error {
	err := p.setupListener()
	if err != nil {
		return err
	}
	// this uses the default server mux
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		err := p.srvr.Serve(p.listener)
		if errors.Is(err, net.ErrClosed) {
			// ErrClose is expected on gracefully shutdown
			p.lggr.Warnf("%s closed", p.Name())
		} else {
			p.lggr.Errorf("%s: %w", p.Name(), err)
		}
	}()
	return nil
}

// Shutdown shutdowns down the underlying HTTP server. See [http.Server.Shutdown] for details
func (p *PromServer) Shutdown(ctx context.Context) error {
	return p.srvr.Shutdown(ctx)
}

// Name of the server
func (p *PromServer) Name() string {
	return fmt.Sprintf("%s-prom-server", p.lggr.Name())
}

func (p *PromServer) Addr() net.Addr {
	return p.listener.Addr()
}

// setupListener creates explicit listener so that we can resolve `:0` ports, which is needed for testing
// if we didn't need the resolved addr, or could pick a static port we could use p.srvr.ListenAndServer
func (p *PromServer) setupListener() error {
	addr := fmt.Sprintf(":%d", p.port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	p.listener = l
	return nil
}
