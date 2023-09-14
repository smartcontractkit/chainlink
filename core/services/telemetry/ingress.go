package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentWrapper{}

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryIngressClient) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType, network string, chainID string) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, contractID, telemType, network, chainID)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
	contractID             string
	telemType              synchronization.TelemetryType
	network                string
	chainID                string
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryIngressClient, contractID string, telemType synchronization.TelemetryType, network string, chainID string) *IngressAgent {
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
	payload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  telemetry,
		ContractID: t.contractID,
		TelemType:  t.telemType,
	}
	t.telemetryIngressClient.Send(payload)
}
