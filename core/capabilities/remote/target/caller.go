package target

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	sender := remote.ToPeerID(msg.Sender)

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
