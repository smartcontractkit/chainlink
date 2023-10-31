package types

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
)

type TelemetryService interface {
	Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error
}

type TelemetryClientEndpoint interface {
	SendLog(ctx context.Context, log []byte) error
}

type TelemetryClient interface {
	TelemetryService
	NewEndpoint(ctx context.Context, nework string, chainID string, contractID string, telemetryType string) (TelemetryClientEndpoint, error)
}

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(network, chainID, contractID, telemetryType string) commontypes.MonitoringEndpoint
}
