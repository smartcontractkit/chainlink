package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentBatchWrapper{}

// IngressAgentBatchWrapper provides monitoring endpoint generation for the telemetry batch client
type IngressAgentBatchWrapper struct {
	telemetryIngressBatchClient synchronization.TelemetryService
}

// NewIngressAgentBatchWrapper creates a new IngressAgentBatchWrapper with the provided telemetry batch client
func NewIngressAgentBatchWrapper(telemetryIngressBatchClient synchronization.TelemetryService) *IngressAgentBatchWrapper {
	return &IngressAgentBatchWrapper{telemetryIngressBatchClient}
}

// GenMonitoringEndpoint returns a new ingress batch agent instantiated with the batch client and a contractID
func (t *IngressAgentBatchWrapper) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return NewTypedIngressAgent(t.telemetryIngressBatchClient, network, chainID, contractID, telemType)
}

func (t *IngressAgentBatchWrapper) GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) MultitypeMonitoringEndpoint {
	return NewMultiIngressAgent(t.telemetryIngressBatchClient, network, chainID, contractID)
}

// TypedIngressAgentBatch allows for sending batch telemetry for a given contractID
type TypedIngressAgentBatch struct {
	*MultiIngressAgentBatch
	telemType synchronization.TelemetryType
}

// NewTypedIngressAgentBatch creates a new TypedIngressAgentBatch with the given batch client and contractID
func NewTypedIngressAgentBatch(telemetryIngressBatchClient synchronization.TelemetryService, network string, chainID string, contractID string, telemType synchronization.TelemetryType) *TypedIngressAgentBatch {
	m := NewMultiIngressAgentBatch(telemetryIngressBatchClient, network, chainID, contractID)
	return &TypedIngressAgentBatch{
		m,
		telemType,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *TypedIngressAgentBatch) SendLog(telemetry []byte) {
	t.telemetryIngressBatchClient.Send(context.Background(), telemetry, t.contractID, t.telemType)
}

// MultiIngressAgentBatch allows for sending batch telemetry for a given contractID
type MultiIngressAgentBatch struct {
	telemetryIngressBatchClient synchronization.TelemetryService
	network                     string
	chainID                     string
	contractID                  string
}

// NewMultiIngressAgentBatch creates a new MultiIngressAgentBatch with the given batch client and contractID
func NewMultiIngressAgentBatch(telemetryIngressBatchClient synchronization.TelemetryService, network string, chainID string, contractID string) *MultiIngressAgentBatch {
	return &MultiIngressAgentBatch{
		telemetryIngressBatchClient,
		network,
		chainID,
		contractID,
	}
}

// SendTypedLog sends a telemetry log to the ingress server
func (t *MultiIngressAgentBatch) SendTypedLog(telemType synchronization.TelemetryType, telemetry []byte) {
	t.telemetryIngressBatchClient.Send(context.Background(), telemetry, t.contractID, telemType)
}
