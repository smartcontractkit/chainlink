package ccipdeployment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRMN(t *testing.T) {
	// TODO: needs to return RMN peerIDs.
	tenv := NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t))
	t.Log(tenv)
	// Use peerIDs to set RMN config.
	// Add a lane, send a message.
}
