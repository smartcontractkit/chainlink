package transmission

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// LocalExecutableCapability handles the transmission protocol required for a target capability that exists in the same don as
// the caller.
type LocalExecutableCapability struct {
	lggr logger.Logger
	capabilities.ExecutableCapability
	localNode    capabilities.Node
	capabilityID string
}

func NewLocalExecutableCapability(lggr logger.Logger, capabilityID string, localDON capabilities.Node, underlying capabilities.ExecutableCapability) *LocalExecutableCapability {
	return &LocalExecutableCapability{
		ExecutableCapability: underlying,
		capabilityID:         capabilityID,
		lggr:                 lggr,
		localNode:            localDON,
	}
}

func (l *LocalExecutableCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if l.localNode.PeerID == nil || l.localNode.WorkflowDON.ID == 0 {
		l.lggr.Debugf("empty DON info, executing immediately")
		return l.ExecutableCapability.Execute(ctx, req)
	}

	response, err := l.ExecutableCapability.Execute(ctx, req)
	if err != nil {
		return response, err
	}

	// Set peer2peerID in the response metadata for local capabilities
	if len(response.Metadata.Metering) == 1 {
		response.Metadata.Metering[0].Peer2PeerID = l.localNode.PeerID.String()
	}

	return response, nil
}
