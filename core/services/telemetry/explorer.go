package telemetry

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

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
