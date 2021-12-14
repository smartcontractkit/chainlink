package telemetry

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ MonitoringEndpointGenerator = &IngressAgentWrapper{}

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryIngressClient) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(contractID string) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, contractID)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
	contractID             string
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryIngressClient, contractID string) *IngressAgent {
	return &IngressAgent{
		telemetryIngressClient,
		contractID,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgent) SendLog(telemetry []byte) {
	payload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  telemetry,
		ContractID: t.contractID,
	}
	t.telemetryIngressClient.Send(payload)
}
