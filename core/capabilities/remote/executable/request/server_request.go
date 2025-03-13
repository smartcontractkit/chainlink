package request

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"slices"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type response struct {
	response []byte
	error    types.Error
	errorMsg string
}

type ServerRequest struct {
	capability capabilities.ExecutableCapability

	capabilityPeerID p2ptypes.PeerID
	capabilityID     string
	capabilityDonID  uint32

	dispatcher types.Dispatcher

	requesters              map[p2ptypes.PeerID]bool
	responseSentToRequester map[p2ptypes.PeerID]bool

	createdTime time.Time

	response *response

	callingDon capabilities.DON

	requestMessageID string
	method           string
	requestTimeout   time.Duration

	mux  sync.Mutex
	lggr logger.Logger
}

var errExternalErrorMsg = errors.New("failed to execute capability")

func NewServerRequest(capability capabilities.ExecutableCapability, method string, capabilityID string, capabilityDonID uint32,
	capabilityPeerID p2ptypes.PeerID,
	callingDon capabilities.DON, requestID string,
	dispatcher types.Dispatcher, requestTimeout time.Duration, lggr logger.Logger) *ServerRequest {
	return &ServerRequest{
		capability:              capability,
		createdTime:             time.Now(),
		capabilityID:            capabilityID,
		capabilityDonID:         capabilityDonID,
		capabilityPeerID:        capabilityPeerID,
		dispatcher:              dispatcher,
		requesters:              map[p2ptypes.PeerID]bool{},
		responseSentToRequester: map[p2ptypes.PeerID]bool{},
		callingDon:              callingDon,
		requestMessageID:        requestID,
		method:                  method,
		requestTimeout:          requestTimeout,
		lggr: lggr.Named("ServerRequest").With(
			"requestID", requestID,
			"capabilityID", capabilityID,
		),
	}
}

func (e *ServerRequest) OnMessage(ctx context.Context, msg *types.MessageBody) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if msg.Sender == nil {
		return errors.New("sender missing from message")
	}

	requester, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return fmt.Errorf("failed to convert message sender to PeerID: %w", err)
	}

	if err := e.addRequester(requester); err != nil {
		return fmt.Errorf("failed to add requester to request: %w", err)
	}

	e.lggr.Debugw("OnMessage called for request", "calls", len(e.requesters),
		"hasResponse", e.response != nil, "requester", requester.String(), "minRequsters", e.callingDon.F+1)

	if e.minimumRequiredRequestsReceived() && !e.hasResponse() {
		switch e.method {
		case types.MethodExecute:
			e.executeRequest(ctx, msg, executeCapabilityRequest)
		default:
			e.setError(types.Error_INTERNAL_ERROR, "unknown method %s"+e.method)
		}
	}

	if err := e.sendResponses(); err != nil {
		return fmt.Errorf("failed to send responses: %w", err)
	}

	return nil
}

func (e *ServerRequest) Expired() bool {
	return time.Since(e.createdTime) > e.requestTimeout
}

func (e *ServerRequest) Cancel(err types.Error, msg string) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if !e.hasResponse() {
		e.setError(err, msg)
		if err := e.sendResponses(); err != nil {
			return fmt.Errorf("failed to send responses: %w", err)
		}
	}

	return nil
}

type executeFn func(ctx context.Context, lggr logger.Logger, capability capabilities.ExecutableCapability, payload []byte) ([]byte, error)

func (e *ServerRequest) executeRequest(ctx context.Context, msg *types.MessageBody, method executeFn) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, e.requestTimeout)
	defer cancel()

	responsePayload, err := method(ctxWithTimeout, e.lggr, e.capability, msg.Payload)
	if err != nil {
		e.setError(types.Error_INTERNAL_ERROR, err.Error())
	} else {
		e.setResult(responsePayload)
	}
}

func (e *ServerRequest) addRequester(from p2ptypes.PeerID) error {
	fromPeerInCallingDon := slices.Contains(e.callingDon.Members, from)

	if !fromPeerInCallingDon {
		return fmt.Errorf("request received from peer %s not in calling don", from)
	}

	if e.requesters[from] {
		return fmt.Errorf("request already received from peer %s", from)
	}

	e.requesters[from] = true

	return nil
}

func (e *ServerRequest) minimumRequiredRequestsReceived() bool {
	return len(e.requesters) >= int(e.callingDon.F+1)
}

func (e *ServerRequest) setResult(result []byte) {
	e.lggr.Debug("setting result on request")
	e.response = &response{
		response: result,
	}
}

func (e *ServerRequest) setError(err types.Error, errMsg string) {
	e.lggr.Debugw("setting error on request", "type", err, "error", errMsg)
	e.response = &response{
		error:    err,
		errorMsg: errMsg,
	}
}

func (e *ServerRequest) hasResponse() bool {
	return e.response != nil
}

func (e *ServerRequest) sendResponses() error {
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

func (e *ServerRequest) sendResponse(requester p2ptypes.PeerID) error {
	responseMsg := types.MessageBody{
		CapabilityId:    e.capabilityID,
		CapabilityDonId: e.capabilityDonID,
		CallerDonId:     e.callingDon.ID,
		Method:          types.MethodExecute,
		MessageId:       []byte(e.requestMessageID),
		Sender:          e.capabilityPeerID[:],
		Receiver:        requester[:],
	}

	if e.response.error != types.Error_OK {
		responseMsg.Error = e.response.error
		responseMsg.ErrorMsg = e.response.errorMsg
	} else {
		responseMsg.Payload = e.response.response
	}

	e.lggr.Debugw("Sending response", "receiver", requester)
	if err := e.dispatcher.Send(requester, &responseMsg); err != nil {
		return fmt.Errorf("failed to send response to dispatcher: %w", err)
	}

	e.responseSentToRequester[requester] = true

	return nil
}

func executeCapabilityRequest(ctx context.Context, lggr logger.Logger, capability capabilities.ExecutableCapability, payload []byte) ([]byte, error) {
	capabilityRequest, err := pb.UnmarshalCapabilityRequest(payload)
	if err != nil {
		lggr.Errorw("failed to unmarshal capability request", "err", err)
		return nil, errExternalErrorMsg
	}

	lggr = lggr.With("metadata", capabilityRequest.Metadata)

	lggr.Debugw("executing capability")
	capResponse, err := capability.Execute(ctx, capabilityRequest)

	if err != nil {
		lggr.Errorw("received execution error", "error", err)
		return nil, errExternalErrorMsg
	}

	responsePayload, err := pb.MarshalCapabilityResponse(capResponse)
	if err != nil {
		lggr.Errorw("failed to marshal capability request", "error", err)
		return nil, errExternalErrorMsg
	}

	lggr.Debug("received execution results")
	return responsePayload, nil
}
