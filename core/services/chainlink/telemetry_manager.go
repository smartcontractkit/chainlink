package chainlink

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

// newTelemetryManager returns the DomainManager, which sends telemetry through Beholder's
// per-domain emitters, when TelemetryIngress.BeholderEnabled is set; otherwise the WSRPCManager.
func newTelemetryManager(cfg config.TelemetryIngress, beholderClient *beholder.Client, csaKeystore keystore.CSA,
	csaPubKeyHex string, lggr logger.SugaredLogger) (telemetry.Manager, error) {
	if !cfg.BeholderEnabled() {
		return telemetry.NewWSRPCManager(cfg, csaKeystore, lggr), nil
	}
	// Checked before beholderClient becomes an interface value, where a nil pointer is not nil.
	if beholderClient == nil {
		return nil, errors.New("TelemetryIngress.BeholderEnabled requires a Beholder client")
	}
	// The manager's counters share the client's meter provider, as the emitters' chip
	// metrics already do. Nil on a hand-built client, where the option would panic.
	var opts []telemetry.DomainManagerOpt
	if beholderClient.MeterProvider != nil {
		opts = append(opts, telemetry.WithMeterProvider(beholderClient.MeterProvider))
	}
	dm, err := telemetry.NewDomainManager(beholderClient, csaPubKeyHex, lggr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry domain manager: %w", err)
	}
	lggr.Info("Telemetry is sent through Beholder; TelemetryIngress.Endpoints and its sending settings are ignored")
	return dm, nil
}
