package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &NoopAgent{}

type NoopAgent struct {
}

// SendLog sends a telemetry log to the ingress service
func (t *NoopAgent) SendLog(log []byte) {
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType, network string, chainID string) ocrtypes.MonitoringEndpoint {
	return t
}
