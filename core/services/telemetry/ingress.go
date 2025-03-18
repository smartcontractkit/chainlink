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

func (t *IngressAgentWrapper) GenMonitoringEndpoint(network, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return NewTypedIngressAgent(t.telemetryIngressClient, network, chainID, contractID, telemType)
}

func (t *IngressAgentWrapper) GenMultitypeMonitoringEndpoint(network, chainID, contractID string) MultitypeMonitoringEndpoint {
	return NewMultiIngressAgent(t.telemetryIngressClient, network, chainID, contractID)
}

type TypedIngressAgent struct {
	*MultiIngressAgent
	telemType synchronization.TelemetryType
}

func NewTypedIngressAgent(telemetryIngressClient synchronization.TelemetryService, network string, chainID string, contractID string, telemType synchronization.TelemetryType) *TypedIngressAgent {
	m := NewMultiIngressAgent(telemetryIngressClient, network, chainID, contractID)
	return &TypedIngressAgent{
		m,
		telemType,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *TypedIngressAgent) SendLog(telemetry []byte) {
	t.telemetryIngressClient.Send(context.Background(), telemetry, t.contractID, t.telemType)
}

type MultiIngressAgent struct {
	telemetryIngressClient synchronization.TelemetryService
	network                string
	chainID                string
	contractID             string
}

func NewMultiIngressAgent(telemetryIngressClient synchronization.TelemetryService, network string, chainID string, contractID string) *MultiIngressAgent {
	return &MultiIngressAgent{
		telemetryIngressClient,
		network,
		chainID,
		contractID,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *MultiIngressAgent) SendTypedLog(telemType synchronization.TelemetryType, telemetry []byte) {
	t.telemetryIngressClient.Send(context.Background(), telemetry, t.contractID, telemType)
}
