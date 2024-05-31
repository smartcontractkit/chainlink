package transmission

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// LocalTargetCapability handles the transmission protocol required for a target capability that exists in the same don as
// the caller.
type LocalTargetCapability struct {
	peerID     p2ptypes.PeerID
	don        capabilities.DON
	underlying capabilities.TargetCapability
}

func NewLocalTargetCapability(peerID p2ptypes.PeerID, don capabilities.DON, underlying capabilities.TargetCapability) *LocalTargetCapability {
	return &LocalTargetCapability{
		peerID:     peerID,
		don:        don,
		underlying: underlying,
	}
}

func (l *LocalTargetCapability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return l.underlying.Info(ctx)
}

func (l *LocalTargetCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return l.underlying.RegisterToWorkflow(ctx, request)
}

func (l *LocalTargetCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return l.underlying.UnregisterFromWorkflow(ctx, request)
}

func (l *LocalTargetCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if req.Config == nil || req.Config.Underlying["schedule"] == nil {
		return l.underlying.Execute(ctx, req)
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
		return l.underlying.Execute(ctx, req)
	}
}
