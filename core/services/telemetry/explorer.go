package telemetry

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &ExplorerAgent{}

type ExplorerAgent struct {
	explorerClient synchronization.ExplorerClient
}

// NewExplorerAgent returns a Agent which is just a thin wrapper over
// the explorerClient for now
func NewExplorerAgent(explorerClient synchronization.ExplorerClient) *ExplorerAgent {
	return &ExplorerAgent{explorerClient}
}

// SendLog sends a telemetry log to the explorer
func (t *ExplorerAgent) SendLog(log []byte) {
	t.explorerClient.Send(context.Background(), log, synchronization.ExplorerBinaryMessage)
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *ExplorerAgent) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return t
}
