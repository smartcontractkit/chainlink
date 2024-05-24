package remote

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type remoteTargetReceiver struct {
	peerID       p2ptypes.PeerID
	underlying   commoncap.TargetCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo capabilities.DON
	workflowDONs map[string]commoncap.DON
	dispatcher   types.Dispatcher
	lggr         logger.Logger

	requestMsgIDToRequest map[string]*remoteTargetCapabilityRequest
	requestTimeout        time.Duration

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

		requestMsgIDToRequest: map[string]*remoteTargetCapabilityRequest{},
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

func (r *remoteTargetReceiver) ExpireRequests(ctx context.Context) {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	for messageId, executeReq := range r.requestMsgIDToRequest {
		if time.Since(executeReq.createdTime) > r.requestTimeout {

			if !executeReq.hasResponse() {
				executeReq.setError(types.Error_TIMEOUT)
				if err := executeReq.sendResponseToAllRequesters(); err != nil {
					r.lggr.Errorw("failed to send timeout response to all requesters", "capabilityId", r.capInfo.ID, "err", err)
				}
			}

			delete(r.requestMsgIDToRequest, messageId)
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

	requester := ToPeerID(msg.Sender)
	messageId := GetMessageID(msg)

	if _, ok := r.requestMsgIDToRequest[messageId]; !ok {
		r.requestMsgIDToRequest[messageId] = newTargetCapabilityRequest(r.capInfo.ID, r.localDonInfo.ID, r.peerID,
			msg.CallerDonId, messageId, r.dispatcher)
	}

	request, ok := r.requestMsgIDToRequest[messageId]

	if err := request.addRequester(requester, msg.CallerDonId, messageId); err != nil {
		r.lggr.Errorw("failed to add request to response", "capabilityId", r.capInfo.ID, "sender",
			requester, "err", err)
		return
	}

	minRequiredRequests := int(callerDon.F + 1)
	if request.getRequestersCount() == minRequiredRequests {

		capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
		if err == nil {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, r.requestTimeout)
			defer cancel()
			capResponseCh, err := r.underlying.Execute(ctxWithTimeout, capabilityRequest)
			if err == nil {
				// TODO working on the assumption that the capability will only ever return one response from its channel (for now at least)
				capResponse := <-capResponseCh
				responsePayload, err := pb.MarshalCapabilityResponse(capResponse)
				if err != nil {
					r.lggr.Errorw("failed to marshal capability response", "capabilityId", r.capInfo.ID, "err", err)
					request.setError(types.Error_INTERNAL_ERROR)
				} else {
					request.setResult(responsePayload)
				}
			} else {
				r.lggr.Errorw("failed to execute capability", "capabilityId", r.capInfo.ID, "err", err)
				request.setError(types.Error_INTERNAL_ERROR)
			}
		} else {
			r.lggr.Errorw("failed to unmarshal capability request", "capabilityId", r.capInfo.ID, "err", err)
			request.setError(types.Error_INVALID_REQUEST)
		}

		if err := request.sendResponseToAllRequesters(); err != nil {
			r.lggr.Errorw("failed to send response to all requesters", "capabilityId", r.capInfo.ID, "err", err)
		}

	} else if request.getRequestersCount() > minRequiredRequests {
		if err := request.sendResponse(requester); err != nil {
			r.lggr.Errorw("failed to send response to requester", "capabilityId", r.capInfo.ID, "err", err)
		}
	}

}

type remoteTargetCapabilityRequest struct {
	id string

	capabilityPeerId p2ptypes.PeerID
	capabilityID     string
	capabilityDonID  string

	dispatcher types.Dispatcher

	requesters        map[p2ptypes.PeerID]bool
	responseReceivers map[p2ptypes.PeerID]bool

	createdTime time.Time

	response      []byte
	responseError types.Error

	initialRequestingDon string
	requestMessageID     string
}

func newTargetCapabilityRequest(capabilityID string, capabilityDonID string, capabilityPeerId p2ptypes.PeerID,
	callingDonID string, requestMessageID string,
	dispatcher types.Dispatcher) *remoteTargetCapabilityRequest {
	return &remoteTargetCapabilityRequest{
		id:                   uuid.New().String(),
		capabilityID:         capabilityID,
		capabilityDonID:      capabilityDonID,
		capabilityPeerId:     capabilityPeerId,
		dispatcher:           dispatcher,
		requesters:           map[p2ptypes.PeerID]bool{},
		responseReceivers:    map[p2ptypes.PeerID]bool{},
		createdTime:          time.Now(),
		initialRequestingDon: callingDonID,
		requestMessageID:     requestMessageID,
	}
}

func (e *remoteTargetCapabilityRequest) addRequester(from p2ptypes.PeerID, fromDonID string, requestMessageID string) error {
	if e.requesters[from] {
		return fmt.Errorf("request already received from peer %s", from)
	}

	if e.initialRequestingDon != fromDonID {
		return fmt.Errorf("received request from different initial requesting don %s, expected %s", fromDonID, e.initialRequestingDon)
	}

	if e.requestMessageID != requestMessageID {
		return fmt.Errorf("received request with different message id %s, expected %s", requestMessageID, e.requestMessageID)
	}

	e.requesters[from] = true

	return nil
}

func (e *remoteTargetCapabilityRequest) getRequestersCount() int {
	return len(e.requesters)
}

func (e *remoteTargetCapabilityRequest) setResult(result []byte) {
	e.response = result
}

func (e *remoteTargetCapabilityRequest) setError(err types.Error) {
	e.responseError = err
}

func (e *remoteTargetCapabilityRequest) hasResponse() bool {
	return e.response != nil || e.responseError != types.Error_OK
}

func (e *remoteTargetCapabilityRequest) sendResponseToAllRequesters() error {
	for requester := range e.requesters {
		if err := e.sendResponse(requester); err != nil {
			return fmt.Errorf("failed to send response to requester %s: %w", requester, err)
		}
	}

	return nil
}

func (e *remoteTargetCapabilityRequest) sendResponse(peer p2ptypes.PeerID) error {
	if err := e.validateResponseSendRequest(peer); err != nil {
		return fmt.Errorf("failed to validate response send request: %w", err)
	}

	responseMsg := types.MessageBody{
		CapabilityId:    e.capabilityID,
		CapabilityDonId: e.capabilityDonID,
		CallerDonId:     e.initialRequestingDon,
		Method:          types.MethodExecute,
		MessageId:       []byte(e.requestMessageID),
		Sender:          e.capabilityPeerId[:],
		Receiver:        peer[:],
	}

	if e.responseError != types.Error_OK {
		responseMsg.Error = e.responseError
	} else {
		responseMsg.Payload = e.response
	}

	if err := e.dispatcher.Send(peer, &responseMsg); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	e.responseReceivers[peer] = true

	return nil
}

func (e *remoteTargetCapabilityRequest) validateResponseSendRequest(peer p2ptypes.PeerID) error {
	if !e.hasResponse() {
		return fmt.Errorf("no response to send")
	}

	if e.responseReceivers[peer] {
		return fmt.Errorf("response already sent to peer")
	}

	return nil
}

func GetMessageID(msg *types.MessageBody) string {
	return string(msg.MessageId)
}
