package remote

import (
	"context"
	"sync"
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

	msgIDToExecuteRequest map[string]executeRequest
	requestTimeout        time.Duration

	receiveLock sync.Mutex
}

var _ types.Receiver = &remoteTargetReceiver{}

func NewRemoteTargetReceiver(ctx context.Context, lggr logger.Logger, underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo *capabilities.DON,
	workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration) *remoteTargetReceiver {

	receiver := &remoteTargetReceiver{
		underlying:   underlying,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		msgIDToExecuteRequest: map[string]executeRequest{},
		requestTimeout:        requestTimeout,

		lggr: lggr,
	}

	go func() {
		timer := time.NewTimer(requestTimeout)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				receiver.ExpireRequests(ctx)
			}
		}
	}()

	return receiver
}

type executeRequest struct {
	fromPeers        map[p2ptypes.PeerID]bool
	response         *types.MessageBody
	callingDonID     string
	firstRequestTime time.Time
}

func (r *remoteTargetReceiver) ExpireRequests(ctx context.Context) {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	for messageId, executeReq := range r.msgIDToExecuteRequest {
		if time.Since(executeReq.firstRequestTime) > r.requestTimeout {

			if executeReq.response == nil {
				responseMsg := &types.MessageBody{
					CapabilityId:    r.capInfo.ID,
					CapabilityDonId: r.localDonInfo.ID,
					CallerDonId:     executeReq.callingDonID,
					Method:          types.MethodExecute,
					MessageId:       []byte(messageId),
					// TODO sort out error codes - this should be a timeout error
					Error: types.Error_CAPABILITY_NOT_FOUND,
				}

				for peerID := range executeReq.fromPeers {
					if err := r.dispatcher.Send(peerID, responseMsg); err != nil {
						r.lggr.Errorw("failed to send time out response", "peer", peerID, "err", err)
					}
				}
			}

			delete(r.msgIDToExecuteRequest, messageId)
		}

	}

}

func (r *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	// TODO should the dispatcher be passing in a context?
	ctx := context.Background()

	// TODO Confirm threading semantics of dispatcher receive
	// TODO May want to have executor per message id to improve liveness
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	// TODO multithread this

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

	messageId := getMessageID(msg)

	executeReq, ok := r.msgIDToExecuteRequest[messageId]
	if !ok {
		executeReq = executeRequest{
			fromPeers:        map[p2ptypes.PeerID]bool{},
			callingDonID:     msg.CallerDonId,
			firstRequestTime: time.Now(),
		}
		r.msgIDToExecuteRequest[messageId] = executeReq
	}

	if executeReq.callingDonID != msg.CallerDonId {
		r.lggr.Warnw("received duplicate execute request from different don, ignoring", "peer", sender)
		return
	}

	if executeReq.fromPeers[sender] {
		r.lggr.Warnw("received duplicate execute request from peer, ignoring", "peer", sender)
		return
	}

	executeReq.fromPeers[sender] = true
	minRequiredRequests := int(callerDon.F + 1)
	if len(executeReq.fromPeers) >= minRequiredRequests {
		if executeReq.response == nil {

			responseMsg := &types.MessageBody{
				CapabilityId:    r.capInfo.ID,
				CapabilityDonId: r.localDonInfo.ID,
				CallerDonId:     msg.CallerDonId,
				Method:          types.MethodExecute,
				MessageId:       []byte(messageId),
			}

			capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
			if err == nil {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, r.requestTimeout)
				defer cancel()
				responseCh, err := r.underlying.Execute(ctxWithTimeout, capabilityRequest)
				if err == nil {
					// TODO handle the case where the capability returns a stream of responses
					response := <-responseCh
					responseMsg.Payload, err = pb.MarshalCapabilityResponse(response)
				} else {
					r.lggr.Errorw("failed to execute capability", "capabilityId", r.capInfo.ID, "err", err)
					// TODO set correct error code
					responseMsg.Error = types.Error_CAPABILITY_NOT_FOUND
				}
			} else {
				r.lggr.Errorw("failed to unmarshal capability request", "capabilityId", r.capInfo.ID, "err", err)
				// TODO set correct error code
				responseMsg.Error = types.Error_CAPABILITY_NOT_FOUND
			}

			executeReq.response = responseMsg

			for peerID := range executeReq.fromPeers {
				if err = r.dispatcher.Send(peerID, responseMsg); err != nil {
					r.lggr.Errorw("failed to send response", "peer", peerID, "err", err)
				}
			}
		} else {
			if err := r.dispatcher.Send(sender, executeReq.response); err != nil {
				r.lggr.Errorw("failed to send response", "peer", sender, "err", err)
			}
		}
	}

}

func getMessageID(msg *types.MessageBody) string {
	return string(msg.MessageId)
}
