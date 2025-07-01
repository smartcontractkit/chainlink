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

func NewLocalExecutableCapability(lggr logger.Logger, capabilityID string, localDON capabilities.Node, underlying capabilities.TargetCapability) *LocalExecutableCapability {
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
	l.lggr.Errorw("METERING_LOGS: starting to set peer2peerID in response metadata",
		"capability_id", l.capabilityID,
		"response_metadata_metering_length", len(response.Metadata.Metering),
		"local_node_peer_id", l.localNode.PeerID,
		"local_node_workflow_don_id", l.localNode.WorkflowDON.ID)

	if len(response.Metadata.Metering) == 1 {
		l.lggr.Errorw("METERING_LOGS: setting peer2peerID for single metering entry",
			"capability_id", l.capabilityID,
			"original_peer2peer_id", response.Metadata.Metering[0].Peer2PeerID,
			"new_peer2peer_id", l.localNode.PeerID.String())

		response.Metadata.Metering[0].Peer2PeerID = l.localNode.PeerID.String()

		l.lggr.Errorw("METERING_LOGS: successfully set peer2peerID",
			"capability_id", l.capabilityID,
			"final_peer2peer_id", response.Metadata.Metering[0].Peer2PeerID)
	} else {
		l.lggr.Errorw("METERING_LOGS: skipping peer2peerID setting - unexpected metering length",
			"capability_id", l.capabilityID,
			"metering_length", len(response.Metadata.Metering),
			"expected_length", 1)
	}

	return response, nil
}
