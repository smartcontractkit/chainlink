package target

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type remoteTargetReceiver struct {
	lggr         logger.Logger
	peerID       p2ptypes.PeerID
	underlying   commoncap.TargetCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo capabilities.DON
	workflowDONs map[string]commoncap.DON
	dispatcher   types.Dispatcher

	requestIDToRequest map[string]*receiverRequest
	requestTimeout     time.Duration

	receiveLock sync.Mutex
}

var _ types.Receiver = &remoteTargetReceiver{}

func NewRemoteTargetReceiver(ctx context.Context, lggr logger.Logger, peerID p2ptypes.PeerID, underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON,
	workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration) *remoteTargetReceiver {

	receiver := &remoteTargetReceiver{
		underlying:   underlying,
		peerID:       peerID,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		requestIDToRequest: map[string]*receiverRequest{},
		requestTimeout:     requestTimeout,

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
				receiver.ExpireRequests()
			}
		}
	}()

	return receiver
}

func (r *remoteTargetReceiver) ExpireRequests() {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	for messageId, executeReq := range r.requestIDToRequest {
		if time.Since(executeReq.createdTime) > r.requestTimeout {

			if !executeReq.hasResponse() {
				executeReq.setError(types.Error_TIMEOUT)
				if err := executeReq.sendResponses(); err != nil {
					r.lggr.Errorw("failed to send timeout response to all requesters", "capabilityId", r.capInfo.ID, "err", err)
				}
			}

			delete(r.requestIDToRequest, messageId)
		}

	}

}

func (r *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	// TODO should the dispatcher be passing in a context?
	ctx := context.Background()

	// TODO Confirm threading semantics of dispatcher Receive
	// TODO May want to have executor per message id to improve liveness
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	// TODO multithread this

	if msg.Method != types.MethodExecute {
		r.lggr.Errorw("received request for unsupported method type", "method", msg.Method)
		return
	}

	// A request is uniquely identified by the message id and the hash of the payload
	messageId := GetMessageID(msg)
	hash := sha256.Sum256(msg.Payload)
	requestID := messageId + hex.EncodeToString(hash[:])

	if _, ok := r.requestIDToRequest[requestID]; !ok {
		if callingDon, ok := r.workflowDONs[msg.CallerDonId]; ok {
			r.requestIDToRequest[requestID] = NewReceiverRequest(r.lggr, r.underlying, r.capInfo.ID, r.localDonInfo.ID, r.peerID,
				callingDon, messageId, r.dispatcher, r.requestTimeout)
		} else {
			r.lggr.Errorw("received request from unregistered workflow don", "donId", msg.CallerDonId)
			return
		}
	}

	request := r.requestIDToRequest[requestID]

	err := request.Receive(ctx, msg)
	if err != nil {
		r.lggr.Errorw("request failed to Receive new message", "request", request, "err", err)
	}

}

func GetMessageID(msg *types.MessageBody) string {
	return string(msg.MessageId)
}
