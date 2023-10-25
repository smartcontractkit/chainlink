package types

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
)

type Telemetry interface {
	Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error
}

// MonitoringEndpointGenerator almost identical to synchronization.MonitoringEndpointGenerator except for the telemetry type
type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(network string, chainID string, contractID string, telemetryType string) commontypes.MonitoringEndpoint
}
