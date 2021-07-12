package telemetry

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

type Agent struct {
	explorerClient synchronization.ExplorerClient
}

// NewAgent returns a Agent which is just a thin wrapper over
// the explorerClient for now
func NewAgent(explorerClient synchronization.ExplorerClient) *Agent {
	return &Agent{explorerClient}
}

// SendLog sends a telemetry log to the explorer
func (t *Agent) SendLog(log []byte) {
	t.explorerClient.Send(context.Background(), log, synchronization.ExplorerBinaryMessage)
}
