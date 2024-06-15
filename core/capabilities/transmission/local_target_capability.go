package transmission

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// LocalTargetCapability handles the transmission protocol required for a target capability that exists in the same don as
// the caller.
type LocalTargetCapability struct {
	lggr logger.Logger
	capabilities.TargetCapability
	localNode capabilities.Node
}

func NewLocalTargetCapability(lggr logger.Logger, localDON capabilities.Node, underlying capabilities.TargetCapability) *LocalTargetCapability {
	return &LocalTargetCapability{
		TargetCapability: underlying,
		lggr:             lggr,
		localNode:        localDON,
	}
}

func (l *LocalTargetCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if l.localNode.PeerID == nil || l.localNode.WorkflowDON.ID == "" {
		l.lggr.Debugf("empty DON info, executing immediately")
		return l.TargetCapability.Execute(ctx, req)
	}

	if req.Config == nil || req.Config.Underlying["schedule"] == nil {
		l.lggr.Debug("no schedule found, executing immediately")
		return l.TargetCapability.Execute(ctx, req)
	}

	peerIDToTransmissionDelay, err := GetPeerIDToTransmissionDelay(l.localNode.WorkflowDON.Members, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay map: %w", err)
	}

	delay, existsForPeerID := peerIDToTransmissionDelay[*l.localNode.PeerID]
	if !existsForPeerID {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
		return l.TargetCapability.Execute(ctx, req)
	}
}
