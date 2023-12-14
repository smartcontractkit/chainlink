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
//
//go:generate mockery --quiet --name Checker --output ./mocks/ --case=underscore
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

type InBackupHealthReport struct {
	server http.Server
	lggr   logger.Logger
}

func (i *InBackupHealthReport) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), time.Second)
	defer shutdownRelease()
	if err := i.server.Shutdown(shutdownCtx); err != nil {
		i.lggr.Errorf("InBackupHealthReport shutdown error: %v", err)
	}
	i.lggr.Info("InBackupHealthReport shutdown complete")
}

func (i *InBackupHealthReport) Start() {
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		i.lggr.Info("Starting InBackupHealthReport")
		if err := i.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			i.lggr.Errorf("InBackupHealthReport server error: %v", err)
		}
	}()
}

// NewInBackupHealthReport creates a new InBackupHealthReport that will serve the /health endpoint, useful for
// preventing shutdowns due to health-checks when running long backup tasks
func NewInBackupHealthReport(port uint16, lggr logger.Logger) *InBackupHealthReport {
	return &InBackupHealthReport{
		server: http.Server{Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: time.Second * 5},
		lggr:   lggr,
	}
}
