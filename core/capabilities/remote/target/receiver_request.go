package target

import (
	"context"
	"fmt"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type receiverRequest struct {
	lggr logger.Logger

	capability capabilities.TargetCapability

	capabilityPeerId p2ptypes.PeerID
	capabilityID     string
	capabilityDonID  string

	dispatcher types.Dispatcher

	requesters              map[p2ptypes.PeerID]bool
	responseSentToRequester map[p2ptypes.PeerID]bool

	createdTime time.Time

	response      []byte
	responseError types.Error

	callingDon       commoncap.DON
	requestMessageID string

	requestTimeout time.Duration
}

func NewReceiverRequest(lggr logger.Logger, capability capabilities.TargetCapability, capabilityID string, capabilityDonID string, capabilityPeerId p2ptypes.PeerID,
	callingDon commoncap.DON, requestMessageID string,
	dispatcher types.Dispatcher, requestTimeout time.Duration) *receiverRequest {
	return &receiverRequest{
		lggr:                    lggr,
		capability:              capability,
		createdTime:             time.Now(),
		capabilityID:            capabilityID,
		capabilityDonID:         capabilityDonID,
		capabilityPeerId:        capabilityPeerId,
		dispatcher:              dispatcher,
		requesters:              map[p2ptypes.PeerID]bool{},
		responseSentToRequester: map[p2ptypes.PeerID]bool{},
		callingDon:              callingDon,
		requestMessageID:        requestMessageID,
		requestTimeout:          requestTimeout,
	}
}

func (e *receiverRequest) receive(ctx context.Context, msg *types.MessageBody) error {
	requester := remote.ToPeerID(msg.Sender)
	if err := e.addRequester(requester); err != nil {
		return fmt.Errorf("failed to add requester to request: %w", err)
	}

	if e.minimumRequiredRequestsReceived() && !e.hasResponse() {
		e.executeRequest(ctx, msg.Payload)
	}

	if err := e.sendResponses(); err != nil {
		return fmt.Errorf("failed to send response to requesters: %w", err)
	}

	return nil
}

func (e *receiverRequest) executeRequest(ctx context.Context, payload []byte) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, e.requestTimeout)
	defer cancel()

	capabilityRequest, err := pb.UnmarshalCapabilityRequest(payload)
	if err != nil {
		e.setError(types.Error_INVALID_REQUEST)
		e.lggr.Errorw("failed to unmarshal capability request", "err", err)
	}

	capResponseCh, err := e.capability.Execute(ctxWithTimeout, capabilityRequest)

	if err != nil {
		e.setError(types.Error_INTERNAL_ERROR)
		e.lggr.Errorw("failed to execute capability", "err", err)
	}

	// TODO working on the assumption that the capability will only ever return one response from its channel (for now at least)
	capResponse := <-capResponseCh
	responsePayload, err := pb.MarshalCapabilityResponse(capResponse)
	if err != nil {
		e.setError(types.Error_INTERNAL_ERROR)
		e.lggr.Errorw("failed to marshal capability response", "err", err)
	}

	e.setResult(responsePayload)
}

func (e *receiverRequest) addRequester(from p2ptypes.PeerID) error {

	fromPeerInCallingDon := false
	for _, member := range e.callingDon.Members {
		if member == from {
			fromPeerInCallingDon = true
			break
		}
	}

	if !fromPeerInCallingDon {
		return fmt.Errorf("request received from peer %s not in calling don", from)
	}

	if e.requesters[from] {
		return fmt.Errorf("request already received from peer %s", from)
	}

	e.requesters[from] = true

	return nil
}

func (e *receiverRequest) minimumRequiredRequestsReceived() bool {
	return len(e.requesters) >= int(e.callingDon.F+1)
}

func (e *receiverRequest) setResult(result []byte) {
	e.response = result
}

func (e *receiverRequest) setError(err types.Error) {
	e.responseError = err
}

func (e *receiverRequest) hasResponse() bool {
	return e.response != nil || e.responseError != types.Error_OK
}

func (e *receiverRequest) sendResponses() error {
	if e.hasResponse() {
		for requester := range e.requesters {
			if !e.responseSentToRequester[requester] {
				e.responseSentToRequester[requester] = true
				if err := e.sendResponse(requester); err != nil {
					return fmt.Errorf("failed to send response to requester %s: %w", requester, err)
				}
			}
		}
	}

	return nil
}

func (e *receiverRequest) sendResponse(requester p2ptypes.PeerID) error {

	responseMsg := types.MessageBody{
		CapabilityId:    e.capabilityID,
		CapabilityDonId: e.capabilityDonID,
		CallerDonId:     e.callingDon.ID,
		Method:          types.MethodExecute,
		MessageId:       []byte(e.requestMessageID),
		Sender:          e.capabilityPeerId[:],
		Receiver:        requester[:],
	}

	if e.responseError != types.Error_OK {
		responseMsg.Error = e.responseError
	} else {
		responseMsg.Payload = e.response
	}

	if err := e.dispatcher.Send(requester, &responseMsg); err != nil {
		return fmt.Errorf("failed to send response to dispatcher: %w", err)
	}

	e.responseSentToRequester[requester] = true

	return nil
}
