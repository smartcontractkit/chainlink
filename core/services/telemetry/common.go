package telemetry

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// Manager is the node's telemetry service: it runs for the node's lifetime and generates
// the monitoring endpoints that jobs send telemetry through. WSRPCManager and DomainManager
// implement it.
type Manager interface {
	services.Service
	MonitoringEndpointGenerator
}

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint
	GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) MultitypeMonitoringEndpoint
}

type MultitypeMonitoringEndpoint interface {
	SendTypedLog(telemType synchronization.TelemetryType, log []byte)
}
