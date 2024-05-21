package remote

// here the only executes when it recieves a report from f + 1 nodes, can use the message cache to collect up these reports

// the chain write is waiting for f + 1 reports to be collected before it will execute the transmission

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type remoteTargetReceiver struct {
	underlying commoncap.TargetCapability
	capInfo    commoncap.CapabilityInfo
	donInfo    *capabilities.DON
	dispatcher types.Dispatcher
	lggr       logger.Logger
}

var _ types.Receiver = &remoteTargetReceiver{}

func NewRemoteTargetReceiver(underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, donInfo *capabilities.DON, dispatcher types.Dispatcher, lggr logger.Logger) *remoteTargetReceiver {
	return &remoteTargetReceiver{
		underlying: underlying,
		capInfo:    capInfo,
		donInfo:    donInfo,
		dispatcher: dispatcher,
		lggr:       lggr,
	}
}

func (c *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	c.lggr.Debugw("not implemented - received message", "capabilityId", c.capInfo.ID, "payload", msg.Payload)
}
