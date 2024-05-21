package remote

// here the only executes when it recieves a report from f + 1 nodes, can use the message cache to collect up these reports

// the chain write is waiting for f + 1 reports to be collected before it will execute the transmission

import (
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type remoteTargetReceiver struct {
	underlying   commoncap.TargetCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo *capabilities.DON
	workflowDONs map[string]commoncap.DON
	dispatcher   types.Dispatcher
	lggr         logger.Logger

	executeRequests map[[32]byte]requestCache
}

var _ types.Receiver = &remoteTargetReceiver{}

func NewRemoteTargetReceiver(underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo *capabilities.DON,
	workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, lggr logger.Logger) *remoteTargetReceiver {
	return &remoteTargetReceiver{
		underlying:   underlying,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		executeRequests: map[[32]byte]requestCache{},

		lggr: lggr,
	}
}

type requestCache struct {
	fromPeers map[p2ptypes.PeerID]bool
	response  *types.MessageBody
	callingDonID string
	firstRequestTime time.Time
}

func (r *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	if msg.Method != types.MethodExecute {
		r.lggr.Errorw("received request for unsupported method type", "method", msg.Method)
		return
	}

	callerDon, ok := r.workflowDONs[msg.CallerDonId]
	if !ok {
		r.lggr.Errorw("received a message from unsupported workflow DON", "capabilityId", r.capInfo.ID, "callerDonId", msg.CallerDonId)
		return
	}

	sender := ToPeerID(msg.Sender)

	var messageId [32]byte
	copy(messageId[:], msg.MessageId)

	rc, ok := r.executeRequests[messageId]
	if !ok {
		rc = requestCache{
			fromPeers: map[p2ptypes.PeerID]bool{},
			callingDonID: msg.CallerDonId,
			firstRequestTime: time.Now(),
		}
		r.executeRequests[messageId] = rc
	}

	if rc.callingDonID != msg.CallerDonId {
		r.lggr.Warnw("received duplicate execute request from different don, ignoring", "peer", sender)
		return
	}

	if rc.fromPeers[sender] {
		r.lggr.Warnw("received duplicate execute request from peer, ignoring", "peer", sender)
		return
	}

	rc.fromPeers[sender] = true
	minRequiredRequests := int(callerDon.F + 1)
	if len(rc.fromPeers) >= minRequiredRequests {
		if rc.response == nil {



			responseMsg := &types.MessageBody{
				CapabilityId:    r.capInfo.ID,
				CapabilityDonId: r.localDonInfo.ID,
				CallerDonId:     msg.CallerDonId,
				Method:          types.MethodExecute,
			}

			capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
			if err == nil {
				
			
				r.lggr.Errorw("failed to unmarshal capability request", "err", err)
				return
			} else {
				responseMsg.Error = types.Error_CAPABILITY_NOT_FOUND

			}
		}
			

			r.underlying.Execute(msg.Payload, func(response []byte) {
			
				
				
			r.lggr.Warnw("received enough execute requests, but no response was provided")
			return
		} else {
			if err := r.dispatcher.Send(sender, rc.response); err != nil {
				r.lggr.Errorw("failed to send response", "peer", sender, "err", err)
			}
		}
	}

}
