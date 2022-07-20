package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ MonitoringEndpointGenerator = &NoopAgent{}

type NoopAgent struct {
}

// SendLog sends a telemetry log to the explorer
func (t *NoopAgent) SendLog(log []byte) {
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(contractID string) ocrtypes.MonitoringEndpoint {
	return t
}
