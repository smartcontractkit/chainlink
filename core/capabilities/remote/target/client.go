package target

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// client is a shim for remote target capabilities.
// It translates between capability API calls and network messages.
// Its responsibilities are:
//  1. Transmit capability requests to remote nodes according to a transmission schedule
//  2. Aggregate responses from remote nodes and return the aggregated response
//
// client communicates with corresponding server on remote nodes.
type client struct {
	services.StateMachine
	lggr                 logger.Logger
	remoteCapabilityInfo commoncap.CapabilityInfo
	localDONInfo         commoncap.DON
	dispatcher           types.Dispatcher
	requestTimeout       time.Duration

	messageIDToCallerRequest map[string]*request.ClientRequest
	mutex                    sync.Mutex
	stopCh                   services.StopChan
	wg                       sync.WaitGroup
}

var _ commoncap.TargetCapability = &client{}
var _ types.Receiver = &client{}
var _ services.Service = &client{}

func NewClient(remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, lggr logger.Logger) *client {
	return &client{
		lggr:                     lggr.Named("TargetClient"),
		remoteCapabilityInfo:     remoteCapabilityInfo,
		localDONInfo:             localDonInfo,
		dispatcher:               dispatcher,
		requestTimeout:           requestTimeout,
		messageIDToCallerRequest: make(map[string]*request.ClientRequest),
		stopCh:                   make(services.StopChan),
	}
}

func (c *client) Start(ctx context.Context) error {
	return c.StartOnce(c.Name(), func() error {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			c.checkForExpiredRequests()
		}()
		c.lggr.Info("TargetClient started")
		return nil
	})
}

func (c *client) Close() error {
	return c.StopOnce(c.Name(), func() error {
		close(c.stopCh)
		c.cancelAllRequests(errors.New("client closed"))
		c.wg.Wait()
		c.lggr.Info("TargetClient closed")
		return nil
	})
}

func (c *client) checkForExpiredRequests() {
	ticker := time.NewTicker(c.requestTimeout)
	defer ticker.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.expireRequests()
		}
	}
}

func (c *client) expireRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for messageID, req := range c.messageIDToCallerRequest {
		if req.Expired() {
			req.Cancel(errors.New("request expired"))
			delete(c.messageIDToCallerRequest, messageID)
		}
	}
}

func (c *client) cancelAllRequests(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, req := range c.messageIDToCallerRequest {
		req.Cancel(err)
	}
}

func (c *client) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.remoteCapabilityInfo, nil
}

func (c *client) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	// do nothing
	return nil
}

func (c *client) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	// do nothing
	return nil
}

func (c *client) Execute(ctx context.Context, capReq commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	req, err := c.executeRequest(ctx, capReq)
	if err != nil {
		return commoncap.CapabilityResponse{}, fmt.Errorf("failed to execute request: %w", err)
	}

	resp := <-req.ResponseChan()
	return resp.CapabilityResponse, resp.Err
}

func (c *client) executeRequest(ctx context.Context, capReq commoncap.CapabilityRequest) (*request.ClientRequest, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageIDForRequest(capReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID for request: %w", err)
	}

	c.lggr.Debugw("executing remote target", "messageID", messageID)

	if _, ok := c.messageIDToCallerRequest[messageID]; ok {
		return nil, fmt.Errorf("request for message ID %s already exists", messageID)
	}

	req, err := request.NewClientRequest(ctx, c.lggr, capReq, messageID, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher,
		c.requestTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create client request: %w", err)
	}

	c.messageIDToCallerRequest[messageID] = req
	return req, nil
}

func (c *client) Receive(ctx context.Context, msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageID(msg)
	if err != nil {
		c.lggr.Errorw("invalid message ID", "err", err, "id", remote.SanitizeLogString(string(msg.MessageId)))
		return
	}

	c.lggr.Debugw("Remote client target receiving message", "messageID", messageID)

	req := c.messageIDToCallerRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID ", "messageID", messageID)
		return
	}

	if err := req.OnMessage(ctx, msg); err != nil {
		c.lggr.Errorw("failed to add response to request", "messageID", messageID, "err", err)
	}
}

func GetMessageIDForRequest(req commoncap.CapabilityRequest) (string, error) {
	if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
		return "", fmt.Errorf("workflow ID is invalid: %w", err)
	}

	if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowExecutionID); err != nil {
		return "", fmt.Errorf("workflow execution ID is invalid: %w", err)
	}

	return req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID, nil
}

func (c *client) Ready() error {
	return nil
}

func (c *client) HealthReport() map[string]error {
	return nil
}

func (c *client) Name() string {
	return "TargetClient"
}
