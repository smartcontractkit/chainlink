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
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// remoteTargetCaller/Receiver are shims translating between capability API calls and network messages
type remoteTargetCaller struct {
	lggr                    logger.Logger
	remoteCapabilityInfo    commoncap.CapabilityInfo
	remoteCapabilityDonInfo capabilities.DON
	localDONInfo            capabilities.DON
	dispatcher              types.Dispatcher
	requestTimeout          time.Duration

	messageIDToExecuteRequest map[string]*callerExecuteRequest
	mutex                     sync.Mutex
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(ctx context.Context, lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, remoteCapabilityDonInfo capabilities.DON, localDonInfo capabilities.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration) *remoteTargetCaller {

	caller := &remoteTargetCaller{
		lggr:                      lggr,
		remoteCapabilityInfo:      remoteCapabilityInfo,
		remoteCapabilityDonInfo:   remoteCapabilityDonInfo,
		localDONInfo:              localDonInfo,
		dispatcher:                dispatcher,
		requestTimeout:            requestTimeout,
		messageIDToExecuteRequest: make(map[string]*callerExecuteRequest),
	}

	go func() {
		timer := time.NewTimer(requestTimeout)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				caller.ExpireRequests(ctx)
			}
		}
	}()

	return caller
}

func (c *remoteTargetCaller) ExpireRequests(ctx context.Context) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for messageID, req := range c.messageIDToExecuteRequest {
		if time.Since(req.creationTime) > c.requestTimeout {
			if !req.responseSent() {
				req.sendResponse(commoncap.CapabilityResponse{Err: errors.New("request timed out")})
			}
		}

		delete(c.messageIDToExecuteRequest, messageID)
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

func (c *remoteTargetCaller) Execute(parentCtx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	deterministicMessageID, err := getDeterministicMessageID(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get deterministic message ID from request: %w", err)
	}

	if _, ok := c.messageIDToExecuteRequest[deterministicMessageID]; ok {
		return nil, fmt.Errorf("request with message ID %s already exists", deterministicMessageID)
	}

	transmissionCtx, transmissionCancelFn := context.WithCancel(parentCtx)
	execRequest := newCallerExecuteRequest(transmissionCancelFn, int(c.remoteCapabilityDonInfo.F+1))

	c.messageIDToExecuteRequest[deterministicMessageID] = execRequest

	if err = c.transmitRequestWithMessageID(transmissionCtx, req, deterministicMessageID); err != nil {
		return nil, fmt.Errorf("failed to transmit request: %w", err)
	}

	return execRequest.responseCh, nil
}

func getDeterministicMessageID(req commoncap.CapabilityRequest) (string, error) {
	if req.Metadata.WorkflowID == "" || req.Metadata.WorkflowExecutionID == "" {
		return "", errors.New("workflow ID and workflow execution ID must be set in request metadata")
	}

	deterministicMessageID := req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID
	return deterministicMessageID, nil
}

// transmitRequestWithMessageID transmits a capability request to remote capabilities according to the transmission configuration
func (c *remoteTargetCaller) transmitRequestWithMessageID(ctx context.Context, req commoncap.CapabilityRequest, messageID string) error {
	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return fmt.Errorf("failed to marshal capability request: %w", err)
	}

	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	peerIDToDelay, err := transmission.GetPeerIDToTransmissionDelay(c.remoteCapabilityDonInfo.Members, c.localDONInfo.Config.SharedSecret, req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID, tc)
	if err != nil {
		return fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
	}

	for peerID, delay := range peerIDToDelay {
		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			message := &types.MessageBody{
				CapabilityId:    c.remoteCapabilityInfo.ID,
				CapabilityDonId: c.remoteCapabilityDonInfo.ID,
				CallerDonId:     c.localDONInfo.ID,
				Method:          types.MethodExecute,
				Payload:         rawRequest,
				MessageId:       []byte(messageID),
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				c.lggr.Debugw("executing delayed execution for peer", "peerID", peerID, "delay", delay)
				err = c.dispatcher.Send(peerID, message)
				if err != nil {
					c.lggr.Errorw("failed to send message", "peerID", peerID, "err", err)
				}
			}
		}(peerID, delay)
	}

	return nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID := GetMessageID(msg)

	req := c.messageIDToExecuteRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID", "messageID", messageID, "sender", msg.Sender)
		return
	}

	if msg.Error != types.Error_OK {
		c.lggr.Warnw("received error response for pending request", "messageID", messageID, "sender", msg.Sender, "receiver", msg.Receiver, "error", msg.Error)
		return
	}

	req.addResponse(msg.Payload)
}

type callerExecuteRequest struct {
	responseCh           chan commoncap.CapabilityResponse
	transmissionCancelFn context.CancelFunc
	creationTime         time.Time
	responseIDCount      map[[32]byte]int

	requiredIdenticalResponses int
	respSent                   bool
}

func newCallerExecuteRequest(transmissionCancelFn context.CancelFunc, requiredIdenticalResponses int) *callerExecuteRequest {
	return &callerExecuteRequest{
		responseCh:                 make(chan commoncap.CapabilityResponse, 1),
		transmissionCancelFn:       transmissionCancelFn,
		responseIDCount:            make(map[[32]byte]int),
		creationTime:               time.Now(),
		requiredIdenticalResponses: requiredIdenticalResponses,
	}
}

func (c *callerExecuteRequest) responseSent() bool {
	return c.respSent
}

// TODO addResponse assumes that only one response is received from each peer, if streaming responses need to be supported this will need to be updated
func (c *callerExecuteRequest) addResponse(response []byte) {
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
}

func (c *callerExecuteRequest) sendResponse(response commoncap.CapabilityResponse) {
	c.responseCh <- response
	close(c.responseCh)
	c.transmissionCancelFn()
	c.respSent = true
}
