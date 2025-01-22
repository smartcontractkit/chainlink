package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	defaultFetchTimeoutMs = 20_000
)

var _ connector.GatewayConnectorHandler = &OutgoingConnectorHandler{}

type OutgoingConnectorHandler struct {
	services.StateMachine
	gc          connector.GatewayConnector
	method      string
	lggr        logger.Logger
	rateLimiter *common.RateLimiter
	responses   *responses
}

func NewOutgoingConnectorHandler(gc connector.GatewayConnector, config ServiceConfig, method string, lgger logger.Logger) (*OutgoingConnectorHandler, error) {
	rateLimiter, err := common.NewRateLimiter(config.RateLimiter)
	if err != nil {
		return nil, err
	}

	if !validMethod(method) {
		return nil, fmt.Errorf("invalid outgoing connector handler method: %s", method)
	}

	return &OutgoingConnectorHandler{
		gc:          gc,
		method:      method,
		responses:   newResponses(),
		rateLimiter: rateLimiter,
		lggr:        lgger,
	}, nil
}

// HandleSingleNodeRequest sends a request to first available gateway node and blocks until response is received
// TODO: handle retries
func (c *OutgoingConnectorHandler) HandleSingleNodeRequest(ctx context.Context, messageID string, req capabilities.Request) (*api.Message, error) {
	// set default timeout if not provided for all outgoing requests
	if req.TimeoutMs == 0 {
		req.TimeoutMs = defaultFetchTimeoutMs
	}

	// Create a subcontext with the timeout plus some margin for the gateway to process the request
	timeoutDuration := time.Duration(req.TimeoutMs) * time.Millisecond
	margin := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration+margin)
	defer cancel()

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fetch request: %w", err)
	}

	ch, err := c.responses.new(messageID)
	if err != nil {
		return nil, fmt.Errorf("duplicate message received for ID: %s", messageID)
	}
	defer c.responses.cleanup(messageID)

	l := logger.With(c.lggr, "messageID", messageID)
	l.Debugw("sending request to gateway")

	body := &api.MessageBody{
		MessageId: messageID,
		DonId:     c.gc.DonID(),
		Method:    c.method,
		Payload:   payload,
	}

	// simply, send request to first available gateway node from sorted list
	// this allows for deterministic selection of gateway node receiver for easier debugging
	gatewayIDs := c.gc.GatewayIDs()
	if len(gatewayIDs) == 0 {
		return nil, errors.New("no gateway nodes available")
	}
	sort.Strings(gatewayIDs)

	selectedGateway := gatewayIDs[0]

	l.Infow("selected gateway, awaiting connection", "gatewayID", selectedGateway)

	if err := c.gc.AwaitConnection(ctx, selectedGateway); err != nil {
		return nil, errors.Wrap(err, "await connection canceled")
	}

	if err := c.gc.SignAndSendToGateway(ctx, selectedGateway, body); err != nil {
		return nil, errors.Wrap(err, "failed to send request to gateway")
	}

	select {
	case resp := <-ch:
		l.Debugw("received response from gateway")
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *OutgoingConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	body := &msg.Body
	l := logger.With(c.lggr, "gatewayID", gatewayID, "method", body.Method, "messageID", msg.Body.MessageId)
	if !c.rateLimiter.Allow(body.Sender) {
		// error is logged here instead of warning because if a message from gateway is rate-limited,
		// the workflow will eventually fail with timeout as there are no retries in place yet
		c.lggr.Errorw("request rate-limited")
		return
	}
	l.Debugw("handling gateway request")
	switch body.Method {
	case capabilities.MethodWebAPITarget, capabilities.MethodComputeAction, capabilities.MethodWorkflowSyncer:
		body := &msg.Body
		var payload capabilities.Response
		err := json.Unmarshal(body.Payload, &payload)
		if err != nil {
			l.Errorw("failed to unmarshal payload", "err", err)
			return
		}
		ch, ok := c.responses.get(body.MessageId)
		if !ok {
			l.Warnw("no response channel found; this may indicate that the node timed out the request")
			return
		}
		select {
		case ch <- msg:
			return
		case <-ctx.Done():
			return
		}
	default:
		l.Errorw("unsupported method")
	}
}

func (c *OutgoingConnectorHandler) Start(ctx context.Context) error {
	return c.StartOnce("OutgoingConnectorHandler", func() error {
		return c.gc.AddHandler([]string{c.method}, c)
	})
}

func (c *OutgoingConnectorHandler) Close() error {
	return c.StopOnce("OutgoingConnectorHandler", func() error {
		return nil
	})
}

func (c *OutgoingConnectorHandler) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}

func (c *OutgoingConnectorHandler) Name() string {
	return c.lggr.Name()
}

func validMethod(method string) bool {
	switch method {
	case capabilities.MethodWebAPITarget, capabilities.MethodComputeAction, capabilities.MethodWorkflowSyncer:
		return true
	default:
		return false
	}
}

func newResponses() *responses {
	return &responses{
		chs: map[string]chan *api.Message{},
	}
}

type responses struct {
	chs map[string]chan *api.Message
	mu  sync.RWMutex
}

func (r *responses) new(id string) (chan *api.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.chs[id]
	if ok {
		return nil, fmt.Errorf("already have response for id: %s", id)
	}

	// Buffered so we don't wait if sending
	ch := make(chan *api.Message, 1)
	r.chs[id] = ch
	return ch, nil
}

func (r *responses) cleanup(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.chs, id)
}

func (r *responses) get(id string) (chan *api.Message, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.chs[id]
	return ch, ok
}
