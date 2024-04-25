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

type remoteTargetReceiver struct {
	capInfo    commoncap.CapabilityInfo
	donInfo    *capabilities.DON
	dispatcher types.Dispatcher
	lggr       logger.Logger
}

var _ types.Receiver = &remoteTargetReceiver{}

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

func (c *remoteTargetCaller) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
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

	// TODO: return a channel that will be closed when all responses are received
	return nil, nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {
	c.lggr.Debugw("not implemented - received message", "capabilityId", c.capInfo.ID, "payload", msg.Payload)
}

func NewRemoteTargetReceiver(capInfo commoncap.CapabilityInfo, donInfo *capabilities.DON, dispatcher types.Dispatcher, lggr logger.Logger) *remoteTargetReceiver {
	return &remoteTargetReceiver{
		capInfo:    capInfo,
		donInfo:    donInfo,
		dispatcher: dispatcher,
		lggr:       lggr,
	}
}

func (c *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	c.lggr.Debugw("not implemented - received message", "capabilityId", c.capInfo.ID, "payload", msg.Payload)
}
