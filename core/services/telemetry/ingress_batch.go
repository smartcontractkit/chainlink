package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentBatchWrapper{}

// IngressAgentBatchWrapper provides monitoring endpoint generation for the telemetry batch client
type IngressAgentBatchWrapper struct {
	telemetryIngressBatchClient synchronization.TelemetryIngressBatchClient
}

// NewIngressAgentBatchWrapper creates a new IngressAgentBatchWrapper with the provided telemetry batch client
func NewIngressAgentBatchWrapper(telemetryIngressBatchClient synchronization.TelemetryIngressBatchClient) *IngressAgentBatchWrapper {
	return &IngressAgentBatchWrapper{telemetryIngressBatchClient}
}

// GenMonitoringEndpoint returns a new ingress batch agent instantiated with the batch client and a contractID
func (t *IngressAgentBatchWrapper) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return NewIngressAgentBatch(t.telemetryIngressBatchClient, contractID, telemType)
}

// IngressAgentBatch allows for sending batch telemetry for a given contractID
type IngressAgentBatch struct {
	telemetryIngressBatchClient synchronization.TelemetryIngressBatchClient
	contractID                  string
	telemType                   synchronization.TelemetryType
}

// NewIngressAgentBatch creates a new IngressAgentBatch with the given batch client and contractID
func NewIngressAgentBatch(telemetryIngressBatchClient synchronization.TelemetryIngressBatchClient, contractID string, telemType synchronization.TelemetryType) *IngressAgentBatch {
	return &IngressAgentBatch{
		telemetryIngressBatchClient,
		contractID,
		telemType,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgentBatch) SendLog(telemetry []byte) {
	payload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  telemetry,
		ContractID: t.contractID,
		TelemType:  t.telemType,
	}
	t.telemetryIngressBatchClient.Send(payload)
}
