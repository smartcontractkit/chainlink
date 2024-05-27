package remote

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// remoteTargetCaller/Receiver are shims translating between capability API calls and network messages
type remoteTargetCaller struct {
	lggr                 logger.Logger
	remoteCapabilityInfo commoncap.CapabilityInfo
	localDONInfo         capabilities.DON
	dispatcher           types.Dispatcher
	requestTimeout       time.Duration

	requestIDToExecuteRequest map[string]*callerExecuteRequest
	mutex                     sync.Mutex
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(ctx context.Context, lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration) *remoteTargetCaller {

	caller := &remoteTargetCaller{
		lggr:                      lggr,
		remoteCapabilityInfo:      remoteCapabilityInfo,
		localDONInfo:              localDonInfo,
		dispatcher:                dispatcher,
		requestTimeout:            requestTimeout,
		requestIDToExecuteRequest: make(map[string]*callerExecuteRequest),
	}

	go func() {
		timer := time.NewTimer(requestTimeout)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				caller.ExpireRequests()
			}
		}
	}()

	return caller
}

func (c *remoteTargetCaller) ExpireRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for messageID, req := range c.requestIDToExecuteRequest {
		if time.Since(req.createdAt) > c.requestTimeout {
			req.cancelRequest("request timed out")
		}

		delete(c.requestIDToExecuteRequest, messageID)
	}
}

func (c *remoteTargetCaller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.remoteCapabilityInfo, nil
}

func (c *remoteTargetCaller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	requestID, err := GetRequestID(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get request ID: %w", err)
	}

	if _, ok := c.requestIDToExecuteRequest[requestID]; ok {
		return nil, fmt.Errorf("request with ID %s already exists", requestID)
	}

	execRequest, err := newCallerExecuteRequest(ctx, c.lggr, req, requestID, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher)

	c.requestIDToExecuteRequest[requestID] = execRequest

	return execRequest.responseCh, nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	requestID := GetMessageID(msg)
	sender := ToPeerID(msg.Sender)

	req := c.requestIDToExecuteRequest[requestID]
	if req == nil {
		c.lggr.Warnw("received response for unknown request ID", "requestID", requestID, "sender", sender)
		return
	}

	if msg.Error != types.Error_OK {
		c.lggr.Warnw("received error response for pending request", "requestID", requestID, "sender", sender, "receiver", msg.Receiver, "error", msg.Error)
		return
	}

	if err := req.addResponse(sender, msg.Payload); err != nil {
		c.lggr.Errorw("failed to add response to request", "requestID", requestID, "sender", sender, "err", err)
	}
}

// Move this into common?
func GetRequestID(req commoncap.CapabilityRequest) (string, error) {
	if req.Metadata.WorkflowID == "" || req.Metadata.WorkflowExecutionID == "" {
		return "", errors.New("workflow ID and workflow execution ID must be set in request metadata")
	}

	return req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID, nil
}

type callerExecuteRequest struct {
	transmissionCtx      context.Context
	responseCh           chan commoncap.CapabilityResponse
	transmissionCancelFn context.CancelFunc
	createdAt            time.Time
	responseIDCount      map[[32]byte]int
	responseReceived     map[p2ptypes.PeerID]bool

	requiredIdenticalResponses int

	respSent bool
}

func newCallerExecuteRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest, messageID string,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher) (*callerExecuteRequest, error) {

	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	peerIDToTransmissionDelay, err := transmission.GetPeerIDToTransmissionDelay(remoteCapabilityDonInfo.Members, localDonInfo.Config.SharedSecret,
		messageID, tc)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
	}

	transmissionCtx, transmissionCancelFn := context.WithCancel(ctx)
	responseReceived := make(map[p2ptypes.PeerID]bool)
	for peerID, delay := range peerIDToTransmissionDelay {
		responseReceived[peerID] = false
		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			message := &types.MessageBody{
				CapabilityId:    remoteCapabilityInfo.ID,
				CapabilityDonId: remoteCapabilityDonInfo.ID,
				CallerDonId:     localDonInfo.ID,
				Method:          types.MethodExecute,
				Payload:         rawRequest,
				MessageId:       []byte(messageID),
			}

			select {
			case <-transmissionCtx.Done():
				return
			case <-time.After(delay):
				err = dispatcher.Send(peerID, message)
				if err != nil {
					lggr.Errorw("failed to send message", "peerID", peerID, "err", err)
				}
			}
		}(peerID, delay)
	}

	return &callerExecuteRequest{
		createdAt:                  time.Now(),
		transmissionCancelFn:       transmissionCancelFn,
		requiredIdenticalResponses: int(remoteCapabilityDonInfo.F + 1),
		responseIDCount:            make(map[[32]byte]int),
		responseReceived:           responseReceived,
		responseCh:                 make(chan commoncap.CapabilityResponse, 1),
	}, nil
}

func (c *callerExecuteRequest) responseSent() bool {
	return c.respSent
}

// TODO addResponse assumes that only one response is received from each peer, if streaming responses need to be supported this will need to be updated
func (c *callerExecuteRequest) addResponse(sender p2ptypes.PeerID, response []byte) error {
	if _, ok := c.responseReceived[sender]; !ok {
		return fmt.Errorf("response from peer %s not expected", sender)
	}

	if c.responseReceived[sender] {
		return fmt.Errorf("response from peer %s already received", sender)
	}

	c.responseReceived[sender] = true

	payloadId := sha256.Sum256(response)
	c.responseIDCount[payloadId]++

	if c.responseIDCount[payloadId] == c.requiredIdenticalResponses {
		capabilityResponse, err := pb.UnmarshalCapabilityResponse(response)
		if err != nil {
			c.sendResponse(commoncap.CapabilityResponse{Err: fmt.Errorf("failed to unmarshal capability response: %w", err)})
		} else {
			c.sendResponse(commoncap.CapabilityResponse{Value: capabilityResponse.Value})
		}
	}

	return nil
}

func (c *callerExecuteRequest) sendResponse(response commoncap.CapabilityResponse) {
	c.responseCh <- response
	close(c.responseCh)
	c.transmissionCancelFn()
	c.respSent = true
}

func (c *callerExecuteRequest) cancelRequest(reason string) {
	c.transmissionCancelFn()
	if !c.responseSent() {
		c.sendResponse(commoncap.CapabilityResponse{Err: errors.New(reason)})
	}
}
