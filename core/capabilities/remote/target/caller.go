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

	messageIDToExecuteRequest map[string]*callerRequest
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
		messageIDToExecuteRequest: make(map[string]*callerRequest),
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

	for messageID, req := range c.messageIDToExecuteRequest {
		if time.Since(req.createdAt) > c.requestTimeout {
			req.cancelRequest("request timed out")
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

func (c *remoteTargetCaller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageIDForRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID for request: %w", err)
	}

	if _, ok := c.messageIDToExecuteRequest[messageID]; ok {
		return nil, fmt.Errorf("request for message ID %s already exists", messageID)
	}

	execRequest, err := NewCallerRequest(ctx, c.lggr, req, messageID, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher)

	c.messageIDToExecuteRequest[messageID] = execRequest

	return execRequest.ResponseChan(), nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID := GetMessageID(msg)
	sender := remote.ToPeerID(msg.Sender)

	req := c.messageIDToExecuteRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID ", "messageID", messageID, "sender", sender)
		return
	}

	if msg.Error != types.Error_OK {
		c.lggr.Warnw("received error response for pending request", "messageID", messageID, "sender", sender, "receiver", msg.Receiver, "error", msg.Error)
		return
	}

	if err := req.AddResponse(sender, msg.Payload); err != nil {
		c.lggr.Errorw("failed to add response to request", "messageID", messageID, "sender", sender, "err", err)
	}
}

func GetMessageIDForRequest(req commoncap.CapabilityRequest) (string, error) {
	if req.Metadata.WorkflowID == "" || req.Metadata.WorkflowExecutionID == "" {
		return "", errors.New("workflow ID and workflow execution ID must be set in request metadata")
	}

	return req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID, nil
}
