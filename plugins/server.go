package plugins

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// StartServer returns a started Server.
// The caller is responsible for calling Server.Stop().
func StartServer(loggerName string) *Server {
	s := Server{
		// default prometheus.Registerer
		GRPCOpts: loop.SetupTelemetry(nil),
	}

	lggr, err := loop.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	lggr = logger.Named(lggr, loggerName)
	s.Logger = logger.Sugared(lggr)

	envCfg, err := GetEnvConfig()
	if err != nil {
		lggr.Fatalf("Failed to get environment configuration: %s\n", err)
	}
	s.PromServer = NewPromServer(envCfg.PrometheusPort(), lggr)
	err = s.PromServer.Start()
	if err != nil {
		lggr.Fatalf("Unrecoverable error starting prometheus server: %s", err)
	}

	s.Checker = services.NewChecker()
	err = s.Checker.Start()
	if err != nil {
		lggr.Fatalf("Failed to start health checker: %v", err)
	}

	return &s
}

// Server holds common plugin server fields.
type Server struct {
	loop.GRPCOpts
	Logger logger.SugaredLogger
	*PromServer
	services.Checker
}

// MustRegister registers the Checkable with services.Checker, or exits upon failure.
func (s *Server) MustRegister(name string, c services.Checkable) {
	err := s.Register(name, c)
	if err != nil {
		s.Logger.Fatalf("Failed to register %s with health checker: %v", name, err)
	}
}

// Stop closes resources and flushes logs.
func (s *Server) Stop() {
	s.Logger.ErrorIfFn(s.Checker.Close, "Failed to close health checker")
	s.Logger.ErrorIfFn(s.PromServer.Close, "error closing prometheus server")
	if err := s.Logger.Sync(); err != nil {
		fmt.Println("Failed to sync logger:", err)
	}
}
