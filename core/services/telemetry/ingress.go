package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentWrapper{}

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryService
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryService) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType, network string, chainID string) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, contractID, telemType, network, chainID)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryService
	contractID             string
	telemType              synchronization.TelemetryType
	network                string
	chainID                string
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryService, contractID string, telemType synchronization.TelemetryType, network string, chainID string) *IngressAgent {
	return &IngressAgent{
		telemetryIngressClient,
		contractID,
		telemType,
		network,
		chainID,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgent) SendLog(telemetry []byte) {
	t.telemetryIngressClient.Send(context.Background(), telemetry, t.contractID, t.telemType)
}
