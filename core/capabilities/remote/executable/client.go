package executable

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// client is a shim for remote executable capabilities.
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

	requestIDToCallerRequest map[string]*request.ClientRequest
	mutex                    sync.Mutex
	stopCh                   services.StopChan
	wg                       sync.WaitGroup
}

var _ commoncap.ExecutableCapability = &client{}
var _ types.Receiver = &client{}
var _ services.Service = &client{}

const expiryCheckInterval = 30 * time.Second

var (
	ErrRequestExpired                  = errors.New("request expired by executable client")
	ErrContextDoneBeforeResponseQuorum = errors.New("context done before remote client received a quorum of responses")
)

func NewClient(remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, lggr logger.Logger) *client {
	return &client{
		lggr:                     lggr.Named("ExecutableCapabilityClient"),
		remoteCapabilityInfo:     remoteCapabilityInfo,
		localDONInfo:             localDonInfo,
		dispatcher:               dispatcher,
		requestTimeout:           requestTimeout,
		requestIDToCallerRequest: make(map[string]*request.ClientRequest),
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
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			c.checkDispatcherReady()
		}()

		c.lggr.Info("ExecutableCapability Client started")
		return nil
	})
}

func (c *client) Close() error {
	return c.StopOnce(c.Name(), func() error {
		close(c.stopCh)
		c.cancelAllRequests(errors.New("client closed"))
		c.wg.Wait()
		c.lggr.Info("ExecutableCapability closed")
		return nil
	})
}

func (c *client) checkDispatcherReady() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			if err := c.dispatcher.Ready(); err != nil {
				c.cancelAllRequests(fmt.Errorf("dispatcher not ready: %w", err))
			}
		}
	}
}

func (c *client) checkForExpiredRequests() {
	tickerInterval := expiryCheckInterval
	if c.requestTimeout < tickerInterval {
		tickerInterval = c.requestTimeout
	}
	ticker := time.NewTicker(tickerInterval)
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

	for messageID, req := range c.requestIDToCallerRequest {
		if req.Expired() {
			req.Cancel(ErrRequestExpired)
			delete(c.requestIDToCallerRequest, messageID)
		}

		if c.dispatcher.Ready() != nil {
			c.cancelAllRequests(errors.New("dispatcher not ready"))
			return
		}
	}
}

func (c *client) cancelAllRequests(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, req := range c.requestIDToCallerRequest {
		req.Cancel(err)
	}
}

func (c *client) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.remoteCapabilityInfo, nil
}

func (c *client) RegisterToWorkflow(ctx context.Context, registerRequest commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (c *client) UnregisterFromWorkflow(ctx context.Context, unregisterRequest commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

func (c *client) Execute(ctx context.Context, capReq commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	req, err := request.NewClientExecuteRequest(ctx, c.lggr, capReq, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher,
		c.requestTimeout)
	if err != nil {
		return commoncap.CapabilityResponse{}, fmt.Errorf("failed to create client request: %w", err)
	}

	if err = c.sendRequest(req); err != nil {
		return commoncap.CapabilityResponse{}, fmt.Errorf("failed to send request: %w", err)
	}

	var respResult []byte
	var respErr error
	select {
	case resp := <-req.ResponseChan():
		respResult = resp.Result
		respErr = resp.Err
	case <-ctx.Done():
		// NOTE: ClientRequest will not block on sending to ResponseChan() because that channel is buffered (with size 1)
		return commoncap.CapabilityResponse{}, errors.Join(ErrContextDoneBeforeResponseQuorum, ctx.Err())
	}

	if respErr != nil {
		return commoncap.CapabilityResponse{}, fmt.Errorf("error executing request: %w", respErr)
	}

	capabilityResponse, err := pb.UnmarshalCapabilityResponse(respResult)
	if err != nil {
		return commoncap.CapabilityResponse{}, fmt.Errorf("failed to unmarshal capability response: %w", err)
	}

	return capabilityResponse, nil
}

func (c *client) sendRequest(req *request.ClientRequest) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lggr.Debugw("executing remote execute capability", "requestID", req.ID())

	if _, ok := c.requestIDToCallerRequest[req.ID()]; ok {
		return fmt.Errorf("request for ID %s already exists", req.ID())
	}

	c.requestIDToCallerRequest[req.ID()] = req
	return nil
}

func (c *client) Receive(ctx context.Context, msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageID(msg)
	if err != nil {
		c.lggr.Errorw("invalid message ID", "err", err, "id", remote.SanitizeLogString(string(msg.MessageId)))
		return
	}

	c.lggr.Debugw("Remote client executable receiving message", "messageID", messageID)

	req := c.requestIDToCallerRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID ", "messageID", messageID)
		return
	}

	if err := req.OnMessage(ctx, msg); err != nil {
		c.lggr.Errorw("failed to add response to request", "messageID", messageID, "err", err)
	}
}

func (c *client) Ready() error {
	return nil
}

func (c *client) HealthReport() map[string]error {
	return nil
}

func (c *client) Name() string {
	return c.lggr.Name()
}
