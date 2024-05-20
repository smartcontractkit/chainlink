package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type executionStrategy interface {
	Apply(ctx context.Context, l logger.Logger, cap capabilities.CallbackCapability, req capabilities.CapabilityRequest) (values.Value, error)
}

var _ executionStrategy = immediateExecution{}

type immediateExecution struct{}

func (i immediateExecution) Apply(ctx context.Context, lggr logger.Logger, cap capabilities.CallbackCapability, req capabilities.CapabilityRequest) (values.Value, error) {
	l, err := capabilities.ExecuteSync(ctx, cap, req)
	if err != nil {
		return nil, err
	}

	// `ExecuteSync` returns a `values.List` even if there was
	// just one return value. If that is the case, let's unwrap the
	// single value to make it easier to use in -- for example -- variable interpolation.
	if len(l.Underlying) > 1 {
		return l, nil
	}

	return l.Underlying[0], nil
}

var _ executionStrategy = scheduledExecution{}

type scheduledExecution struct {
	DON      *capabilities.DON
	PeerID   *p2ptypes.PeerID
	Position int
}

// scheduledExecution generates a pseudo-random transmission schedule,
// and delays execution until a node is required to transmit.
func (d scheduledExecution) Apply(ctx context.Context, lggr logger.Logger, cap capabilities.CallbackCapability, req capabilities.CapabilityRequest) (values.Value, error) {
	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	info, err := cap.Info(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	// Case 1: Local DON
	case info.DON == nil:

		peerIDToTransmissionDelay, err := transmission.GetPeerIDToTransmissionDelay(d.DON.Members, d.DON.Config.SharedSecret,
			req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID, tc)
		if err != nil {
			return nil, fmt.Errorf("failed to get peer ID to transmission delay map: %w", err)
		}

		delay := peerIDToTransmissionDelay[*d.PeerID]
		if delay == nil {
			lggr.Debugw("skipping transmission: node is not included in schedule")
			return nil, nil
		}

		lggr.Debugf("execution delayed by %+v", *delay)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(*delay):
			lggr.Debugw("executing delayed execution")
			return immediateExecution{}.Apply(ctx, lggr, cap, req)
		}
	// Case 2: Remote DON
	default:

		// In this case just execute immediately on the capability and the shims will handle the scheduling and f+1 aggregation

		// TODO: fill in the remote DON case once consensus has been reach on what to do.
		lggr.Debugw("remote DON transmission not implemented: using immediate execution")
		return immediateExecution{}.Apply(ctx, lggr, cap, req)
	}
}
