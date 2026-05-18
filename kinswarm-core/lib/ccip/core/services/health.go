package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ Checker = (*services.HealthChecker)(nil)

// Checker provides a service which can be probed for system health.
type Checker interface {
	// Register a service for health checks.
	Register(service services.HealthReporter) error
	// Unregister a service.
	Unregister(name string) error
	// IsReady returns the current readiness of the system.
	// A system is considered ready if all checks are passing (no errors)
	IsReady() (ready bool, errors map[string]error)
	// IsHealthy returns the current health of the system.
	// A system is considered healthy if all checks are passing (no errors)
	IsHealthy() (healthy bool, errors map[string]error)

	Start() error
	Close() error
}

type StartUpHealthReport struct {
	server http.Server
	lggr   logger.Logger
	mux    *http.ServeMux
}

func (i *StartUpHealthReport) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), time.Second)
	defer shutdownRelease()
	if err := i.server.Shutdown(shutdownCtx); err != nil {
		i.lggr.Errorf("StartUpHealthReport shutdown error: %v", err)
	}
	i.lggr.Info("StartUpHealthReport shutdown complete")
}

func (i *StartUpHealthReport) Start() {
	go func() {
		i.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		i.lggr.Info("Starting StartUpHealthReport")
		if err := i.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			i.lggr.Errorf("StartUpHealthReport server error: %v", err)
		}
	}()
}

// NewStartUpHealthReport creates a new StartUpHealthReport that will serve the /health endpoint, useful for
// preventing shutdowns due to health-checks when running long backup tasks or migrations
func NewStartUpHealthReport(port uint16, lggr logger.Logger) *StartUpHealthReport {
	mux := http.NewServeMux()
	return &StartUpHealthReport{
		lggr:   lggr,
		mux:    mux,
		server: http.Server{Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: time.Second * 5, Handler: mux},
	}
}
