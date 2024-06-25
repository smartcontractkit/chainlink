package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

//go:generate mockery --quiet --name MonitoringEndpointGenerator --output ./mocks --case=underscore
type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint
}
