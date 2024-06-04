package transmission

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// LocalTargetCapability handles the transmission protocol required for a target capability that exists in the same don as
// the caller.
type LocalTargetCapability struct {
	lggr logger.Logger
	capabilities.TargetCapability
	peerID p2ptypes.PeerID
	don    capabilities.DON
}

func NewLocalTargetCapability(lggr logger.Logger, peerID p2ptypes.PeerID, don capabilities.DON, underlying capabilities.TargetCapability) *LocalTargetCapability {
	return &LocalTargetCapability{
		TargetCapability: underlying,
		lggr:             lggr,
		peerID:           peerID,
		don:              don,
	}
}

func (l *LocalTargetCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if req.Config == nil || req.Config.Underlying["schedule"] == nil {
		l.lggr.Debug("no schedule found, executing immediately")
		return l.TargetCapability.Execute(ctx, req)
	}

	tc, err := ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	peerIDToTransmissionDelay, err := GetPeerIDToTransmissionDelay(l.don.Members, l.don.Config.SharedSecret,
		req.Metadata.WorkflowID+req.Metadata.WorkflowExecutionID, tc)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay map: %w", err)
	}

	delay, existsForPeerID := peerIDToTransmissionDelay[l.peerID]
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
