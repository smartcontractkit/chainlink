package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(contractID string) ocrtypes.MonitoringEndpoint
}
