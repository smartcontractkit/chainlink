package plugins

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// NewStartedServer returns a started Server.
// The caller is responsible for calling Server.Stop().
func NewStartedServer(loggerName string) (*Server, error) {
	s, err := newServer(loggerName)
	if err != nil {
		return nil, err
	}
	err = s.start()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// MustNewStartedServer returns a new started Server like NewStartedServer, but logs and exits in the event of error.
// The caller is responsible for calling Server.Stop().
func MustNewStartedServer(loggerName string) *Server {
	s, err := newServer(loggerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %s\n", err)
		os.Exit(1)
	}
	err = s.start()
	if err != nil {
		s.Logger.Fatalf("Failed to start server: %s", err)
	}

	return s
}

// Server holds common plugin server fields.
type Server struct {
	loop.GRPCOpts
	Logger logger.SugaredLogger
	*PromServer
	services.Checker
}

func newServer(loggerName string) (*Server, error) {
	s := &Server{
		// default prometheus.Registerer
		GRPCOpts: loop.SetupTelemetry(nil),
	}

	lggr, err := loop.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %s", err)
	}
	lggr = logger.Named(lggr, loggerName)
	s.Logger = logger.Sugared(lggr)
	return s, nil
}

func (s *Server) start() error {
	envCfg, err := GetEnvConfig()
	if err != nil {
		return fmt.Errorf("error getting environment configuration: %w", err)
	}
	s.PromServer = NewPromServer(envCfg.PrometheusPort(), s.Logger)
	err = s.PromServer.Start()
	if err != nil {
		return fmt.Errorf("error starting prometheus server: %w", err)
	}

	s.Checker = services.NewChecker()
	err = s.Checker.Start()
	if err != nil {
		return fmt.Errorf("error starting health checker: %w", err)
	}

	return nil
}

// MustRegister registers the Checkable with services.Checker, or exits upon failure.
func (s *Server) MustRegister(c services.Checkable) {
	if err := s.Register(c); err != nil {
		s.Logger.Fatalf("Failed to register %s with health checker: %v", c.Name(), err)
	}
}

// Stop closes resources and flushes logs.
func (s *Server) Stop() {
	s.Logger.ErrorIfFn(s.Checker.Close, "Failed to close health checker")
	s.Logger.ErrorIfFn(s.PromServer.Close, "Failed to close prometheus server")
	if err := s.Logger.Sync(); err != nil {
		fmt.Println("Failed to sync logger:", err)
	}
}
