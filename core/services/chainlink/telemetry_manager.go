package chainlink

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonsrv "github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

// telemetryManager generates the node's telemetry endpoints.
type telemetryManager interface {
	commonsrv.Service
	telemetry.MonitoringEndpointGenerator
}

// newTelemetryManager returns the DomainManager, which sends telemetry through Beholder's
// per-domain emitters, when TelemetryIngress.BeholderEnabled is set; otherwise the WSRPC Manager.
func newTelemetryManager(cfg config.TelemetryIngress, beholderClient *beholder.Client, csaKeystore keystore.CSA,
	csaPubKeyHex string, lggr logger.SugaredLogger) (telemetryManager, error) {
	if !cfg.BeholderEnabled() {
		return telemetry.NewManager(cfg, csaKeystore, lggr), nil
	}
	// Checked before beholderClient becomes an interface value, where a nil pointer is not nil.
	if beholderClient == nil {
		return nil, errors.New("TelemetryIngress.BeholderEnabled requires a Beholder client")
	}
	dm, err := telemetry.NewDomainManager(beholderClient, csaPubKeyHex, lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry domain manager: %w", err)
	}
	lggr.Info("Telemetry is sent through Beholder; TelemetryIngress.Endpoints and its sending settings are ignored")
	return dm, nil
}
