package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint
}
