package remote

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// remoteTargetCaller/Receiver are shims translating between capability API calls and network messages
type remoteTargetCaller struct {
	capInfo    commoncap.CapabilityInfo
	donInfo    *capabilities.DON
	dispatcher types.Dispatcher
	lggr       logger.Logger
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(capInfo commoncap.CapabilityInfo, donInfo *capabilities.DON, dispatcher types.Dispatcher, lggr logger.Logger) *remoteTargetCaller {
	return &remoteTargetCaller{
		capInfo:    capInfo,
		donInfo:    donInfo,
		dispatcher: dispatcher,
		lggr:       lggr,
	}
}

func (c *remoteTargetCaller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.capInfo, nil
}

func (c *remoteTargetCaller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {

	/*
		if c.capInfo.DON == nil {
			return nil, errors.New("missing DON in capability info")
		}

		tc, err := workflows.ExtractTransmissionConfig(req.Config)
		if err != nil {
			return nil, err
		}

		n := len(c.capInfo.DON.Members)
		key := workflows.ScheduleSeed(c.donInfo.Config.SharedSecret, req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID)
		sched, err := workflows.Schedule(tc.Schedule, n)
		if err != nil {
			return nil, err
		}

		picked := permutation.Permutation(n, key)
		delay := workflows.DelayFor(d.Position, sched, picked, tc.DeltaStage)
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

		c.lggr.Debugw("not implemented - executing fake remote target capability", "capabilityId", c.capInfo.ID, "nMembers", len(c.donInfo.Members))
		for _, peerID := range c.donInfo.Members {
			m := &types.MessageBody{
				CapabilityId:    c.capInfo.ID,
				CapabilityDonId: c.donInfo.ID,
				Payload:         []byte{0x01, 0x02, 0x03},
			}
			err := c.dispatcher.Send(peerID, m)
			if err != nil {
				return nil, err
			}
		}


	*/
	// TODO: return a channel that will be closed when all responses are received
	return nil, nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {
	c.lggr.Debugw("not implemented - received message", "capabilityId", c.capInfo.ID, "payload", msg.Payload)
}
