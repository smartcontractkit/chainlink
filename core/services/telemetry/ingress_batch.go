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
func (t *IngressAgentBatchWrapper) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType, network string, chainID string) ocrtypes.MonitoringEndpoint {
	return NewIngressAgentBatch(t.telemetryIngressBatchClient, contractID, telemType, network, chainID)
}

// IngressAgentBatch allows for sending batch telemetry for a given contractID
type IngressAgentBatch struct {
	telemetryIngressBatchClient synchronization.TelemetryService
	contractID                  string
	telemType                   synchronization.TelemetryType
	network                     string
	chainID                     string
}

// NewIngressAgentBatch creates a new IngressAgentBatch with the given batch client and contractID
func NewIngressAgentBatch(telemetryIngressBatchClient synchronization.TelemetryService, contractID string, telemType synchronization.TelemetryType, network string, chainID string) *IngressAgentBatch {
	return &IngressAgentBatch{
		telemetryIngressBatchClient,
		contractID,
		telemType,
		network,
		chainID,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgentBatch) SendLog(telemetry []byte) {
	t.telemetryIngressBatchClient.Send(context.Background(), telemetry, t.contractID, t.telemType)
}
