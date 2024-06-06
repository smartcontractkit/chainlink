package target

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
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
	lggr                 logger.Logger
	remoteCapabilityInfo commoncap.CapabilityInfo
	localDONInfo         capabilities.DON
	dispatcher           types.Dispatcher
	requestTimeout       time.Duration

	messageIDToCallerRequest map[string]*request.ClientRequest
	mutex                    sync.Mutex
}

var _ commoncap.TargetCapability = &client{}
var _ types.Receiver = &client{}

func NewClient(ctx context.Context, lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration) *client {
	c := &client{
		lggr:                     lggr,
		remoteCapabilityInfo:     remoteCapabilityInfo,
		localDONInfo:             localDonInfo,
		dispatcher:               dispatcher,
		requestTimeout:           requestTimeout,
		messageIDToCallerRequest: make(map[string]*request.ClientRequest),
	}

	go func() {
		ticker := time.NewTicker(requestTimeout)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.expireRequests()
			}
		}
	}()

	return c
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

func (c *client) Execute(ctx context.Context, capReq commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageIDForRequest(capReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID for request: %w", err)
	}

	if _, ok := c.messageIDToCallerRequest[messageID]; ok {
		return nil, fmt.Errorf("request for message ID %s already exists", messageID)
	}

	req, err := request.NewClientRequest(ctx, c.lggr, capReq, messageID, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher,
		c.requestTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create client request: %w", err)
	}

	c.messageIDToCallerRequest[messageID] = req

	return req.ResponseChan(), nil
}

func (c *client) Receive(msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO should the dispatcher be passing in a context?
	ctx := context.Background()

	messageID := GetMessageID(msg)

	req := c.messageIDToCallerRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID ", "messageID", messageID)
		return
	}

	go func() {
		if err := req.OnMessage(ctx, msg); err != nil {
			c.lggr.Errorw("failed to add response to request", "messageID", messageID, "err", err)
		}
	}()
}

func GetMessageIDForRequest(req commoncap.CapabilityRequest) (string, error) {
	if req.Metadata.WorkflowID == "" || req.Metadata.WorkflowExecutionID == "" {
		return "", errors.New("workflow ID and workflow execution ID must be set in request metadata")
	}

	return req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID, nil
}

func (c *client) Start(ctx context.Context) error {
	return nil
}

func (c *client) Close() error {
	return nil
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
