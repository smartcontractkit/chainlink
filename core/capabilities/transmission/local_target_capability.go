package transmission

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// LocalTargetCapability handles the transmission protocol required for a target capability that exists in the same don as
// the caller.
type LocalTargetCapability struct {
	lggr logger.Logger
	capabilities.TargetCapability
	localNode    capabilities.Node
	capabilityID string
}

func NewLocalTargetCapability(lggr logger.Logger, capabilityID string, localDON capabilities.Node, underlying capabilities.TargetCapability) *LocalTargetCapability {
	return &LocalTargetCapability{
		TargetCapability: underlying,
		capabilityID:     capabilityID,
		lggr:             lggr,
		localNode:        localDON,
	}
}

func (l *LocalTargetCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if l.localNode.PeerID == nil || l.localNode.WorkflowDON.ID == 0 {
		l.lggr.Debugf("empty DON info, executing immediately")
		return l.TargetCapability.Execute(ctx, req)
	}

	if req.Config == nil || req.Config.Underlying["schedule"] == nil {
		l.lggr.Debug("no schedule found, executing immediately")
		return l.TargetCapability.Execute(ctx, req)
	}

	peerIDToTransmissionDelay, err := GetPeerIDToTransmissionDelay(l.localNode.WorkflowDON.Members, req)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("capability id: %s failed to get peer ID to transmission delay map: %w", l.capabilityID, err)
	}

	delay, existsForPeerID := peerIDToTransmissionDelay[*l.localNode.PeerID]
	if !existsForPeerID {
		return capabilities.CapabilityResponse{}, nil
	}

	select {
	case <-ctx.Done():
		return capabilities.CapabilityResponse{}, ctx.Err()
	case <-time.After(delay):
		return l.TargetCapability.Execute(ctx, req)
	}
}
