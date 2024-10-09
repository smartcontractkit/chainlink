package ccipdeployment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRMN(t *testing.T) {
	t.Skip("Local only")
	// TODO: needs to return RMN peerIDs.
	_, rmnCluster := NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t))
	for rmnNode, rmn := range rmnCluster.Nodes {
		t.Log(rmnNode, rmn.Proxy.PeerID, rmn.RMN.OffchainPublicKey, rmn.RMN.EVMOnchainPublicKey)
	}
	// Use peerIDs to set RMN config.
	// Add a lane, send a message.
}
