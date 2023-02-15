package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentWrapper{}

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryIngressClient) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, contractID, telemType)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
	contractID             string
	telemType              synchronization.TelemetryType
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryIngressClient, contractID string, telemType synchronization.TelemetryType) *IngressAgent {
	return &IngressAgent{
		telemetryIngressClient,
		contractID,
		telemType,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgent) SendLog(telemetry []byte) {
	payload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  telemetry,
		ContractID: t.contractID,
		TelemType:  t.telemType,
	}
	t.telemetryIngressClient.Send(payload)
}
