package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	DefaultGlobalRPS      = 100.0
	DefaultGlobalBurst    = 100
	DefaultPerSenderRPS   = 100.0
	DefaultPerSenderBurst = 100
	DefaultWorkflowRPS    = 5.0
	DefaultWorkflowBurst  = 50
	defaultFetchTimeoutMs = 20_000

	errorOutgoingRatelimitGlobal   = "global limit of gateways requests has been exceeded"
	errorOutgoingRatelimitWorkflow = "workflow exceeded limit of gateways requests"
	errorIncomingRatelimitGlobal   = "message from gateway exceeded global rate limit"
	errorIncomingRatelimitSender   = "message from gateway exceeded per sender rate limit"
)

var _ connector.GatewayConnectorHandler = &OutgoingConnectorHandler{}

type OutgoingConnectorHandler struct {
	services.StateMachine
	gc                  connector.GatewayConnector
	method              string
	lggr                logger.Logger
	incomingRateLimiter *common.RateLimiter
	outgoingRateLimiter *common.RateLimiter
	responses           *responses
	selectorOpts        []func(*RoundRobinSelector)
	metrics             *metrics
}

func NewOutgoingConnectorHandler(gc connector.GatewayConnector, config ServiceConfig, method string, lgger logger.Logger, opts ...func(*RoundRobinSelector)) (*OutgoingConnectorHandler, error) {
	outgoingRLCfg := outgoingRateLimiterConfigDefaults(config.OutgoingRateLimiter)
	outgoingRateLimiter, err := common.NewRateLimiter(outgoingRLCfg)
	if err != nil {
		return nil, err
	}
	incomingRLCfg := incomingRateLimiterConfigDefaults(config.RateLimiter)
	incomingRateLimiter, err := common.NewRateLimiter(incomingRLCfg)
	if err != nil {
		return nil, err
	}

	if !validMethod(method) {
		return nil, fmt.Errorf("invalid outgoing connector handler method: %s", method)
	}

	m, err := newMetrics(method)
	if err != nil {
		return nil, err
	}

	return &OutgoingConnectorHandler{
		gc:                  gc,
		method:              method,
		responses:           newResponses(),
		outgoingRateLimiter: outgoingRateLimiter,
		incomingRateLimiter: incomingRateLimiter,
		lggr:                lgger,
		selectorOpts:        opts,
		metrics:             m,
	}, nil
}

// HandleSingleNodeRequest sends a request to first available gateway node and blocks until response is received
// TODO: handle retries
func (c *OutgoingConnectorHandler) HandleSingleNodeRequest(ctx context.Context, messageID string, req capabilities.Request) (*api.Message, error) {
	start := time.Now()

	m, err := c.handleSingleNodeRequest(ctx, messageID, req)

	totalDuration := time.Since(start)
	status := "fail"
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		status = "timeout"
	case err == nil:
		status = "success"
	}
	c.metrics.recordSingleNodeRequestDuration(ctx, totalDuration, status, req.WorkflowID)

	return m, err
}

func (c *OutgoingConnectorHandler) handleSingleNodeRequest(ctx context.Context, messageID string, req capabilities.Request) (*api.Message, error) {
	lggr := logger.With(c.lggr, "messageID", messageID, "workflowID", req.WorkflowID)
	workflowAllow, globalAllow := c.outgoingRateLimiter.AllowVerbose(req.WorkflowID)
	if !workflowAllow {
		return nil, errors.New(errorOutgoingRatelimitWorkflow)
	}
	if !globalAllow {
		return nil, errors.New(errorOutgoingRatelimitGlobal)
	}

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

	lggr.Debugw("sending request to gateway")

	body := &api.MessageBody{
		MessageId: messageID,
		DonId:     c.gc.DonID(),
		Method:    c.method,
		Payload:   payload,
	}

	start := time.Now()
	selectedGateway, err := c.awaitConnection(ctx, awaitContext{
		messageID:  messageID,
		workflowID: req.WorkflowID,
	})
	c.metrics.recordAwaitConnectionDuration(ctx, time.Since(start), req.WorkflowID, selectedGateway, err == nil)
	if err != nil {
		return nil, err
	}

	if err := c.gc.SignAndSendToGateway(ctx, selectedGateway, body); err != nil {
		return nil, errors.Wrap(err, "failed to send request to gateway")
	}

	select {
	case resp := <-ch:
		switch resp.Body.Method {
		case api.MethodInternalError:
			var errPayload api.JsonRPCError
			err := json.Unmarshal(resp.Body.Payload, &errPayload)
			if err != nil {
				lggr.Errorw("failed to unmarshal err payload", "err", err)
				return nil, errors.New("unknown internal error")
			}
			return nil, errors.New(errPayload.Message)
		default:
			lggr.Debugw("received response from gateway")
			return resp, nil
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// awaitContext are context values useful for tracing the logs of awaiting connections.
type awaitContext struct {
	gateway    string
	workflowID string
	messageID  string
}

// awaitConnection attempts to establish a connection to an available gateway.  It iterates through available gateways
// using a round robin selector, connecting to the first available.  The method respects the provided context, allowing for
// cancellation or timeout.
func (c *OutgoingConnectorHandler) awaitConnection(ctx context.Context, md awaitContext) (string, error) {
	lggr := logger.With(c.lggr, "messageID", md.messageID, "workflowID", md.workflowID)
	selector := NewRoundRobinSelector(c.gc.GatewayIDs(), c.selectorOpts...)
	attempts := make(map[string]int)
	backoff := 10 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			gateway, err := selector.NextGateway()
			if err != nil {
				return "", fmt.Errorf("failed to select gateway: %w", err)
			}

			md.gateway = gateway

			if attempts[gateway] > 0 {
				if allGatewaysAttempted(attempts) {
					lggr.Warnw("all available gateway nodes attempted without connection, backing off", "waitTime", backoff)

					select {
					case <-ctx.Done():
						return "", ctx.Err()
					case <-time.After(backoff):
						// backoff completed, update state and continue with next iteration
						attempts = make(map[string]int)
						backoff *= 2
					}
				}
			}

			attempts[gateway]++

			lggr.Infow("selected gateway, awaiting connection", "selectedGateway", gateway)

			if err := c.attemptGatewayConnection(ctx, md); err != nil {
				lggr.Warnw("failed to await connection to gateway node, retrying", "selectedGateway", gateway, "error", err)
				continue
			}

			lggr.Debugw("connected successfully", "selectedGateway", gateway)
			return gateway, nil
		}
	}
}

// allGatewaysAttempted checks if all available gateways have been attempted.
func allGatewaysAttempted(attempts map[string]int) bool {
	for _, count := range attempts {
		if count == 0 {
			return false
		}
	}
	return true
}

// attemptGatewayConnection waits to connect to a gateway with a new child context
func (c *OutgoingConnectorHandler) attemptGatewayConnection(ctx context.Context, md awaitContext) error {
	lggr := logger.With(c.lggr, "messageID", md.messageID, "workflowID", md.workflowID, "selectedGateway", md.gateway)
	timeout := 1_000 * time.Millisecond

	lggr.Debugw("awaiting connection", "timeout", timeout)

	// create a new child context to wait on gateway connection
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := c.gc.AwaitConnection(ctxWithTimeout, md.gateway); err != nil {
		return fmt.Errorf("gateway connection failed: %w", err)
	}
	return nil
}

// HandleGatewayMessage processes incoming messages from the Gateway,
// which are in response to a HandleSingleNodeRequest call.
func (c *OutgoingConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	body := &msg.Body
	l := logger.With(c.lggr, "gatewayID", gatewayID, "method", body.Method, "messageID", msg.Body.MessageId)

	ch, ok := c.responses.get(body.MessageId)
	if !ok {
		l.Warnw("no response channel found; this may indicate that the node timed out the request")
		return
	}

	senderAllow, globalAllow := c.incomingRateLimiter.AllowVerbose(body.Sender)
	errJSON := api.JsonRPCError{
		Code:    500,
		Message: "",
	}
	if !senderAllow {
		errJSON.Message = errorIncomingRatelimitSender
	}
	if !globalAllow {
		if errJSON.Message == "" {
			errJSON.Message = errorIncomingRatelimitGlobal
		} else {
			errJSON.Message += "\n" + errorIncomingRatelimitGlobal
		}
	}

	if errJSON.Message != "" {
		l.Errorw("request rate-limited")
		errPayload, err := json.Marshal(errJSON)
		if err != nil {
			l.Errorw("failed to marshal err payload", "err", err)
		}
		errMsg := api.Message{
			Body: api.MessageBody{
				MessageId: body.MessageId,
				Method:    api.MethodInternalError,
				Payload:   errPayload,
			},
		}
		ch <- &errMsg
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

func incomingRateLimiterConfigDefaults(config common.RateLimiterConfig) common.RateLimiterConfig {
	if config.GlobalBurst == 0 {
		config.GlobalBurst = DefaultGlobalBurst
	}
	if config.GlobalRPS == 0 {
		config.GlobalRPS = DefaultGlobalRPS
	}
	if config.PerSenderBurst == 0 {
		config.PerSenderBurst = DefaultPerSenderBurst
	}
	if config.PerSenderRPS == 0 {
		config.PerSenderRPS = DefaultPerSenderRPS
	}
	return config
}
func outgoingRateLimiterConfigDefaults(config common.RateLimiterConfig) common.RateLimiterConfig {
	if config.GlobalBurst == 0 {
		config.GlobalBurst = DefaultGlobalBurst
	}
	if config.GlobalRPS == 0 {
		config.GlobalRPS = DefaultGlobalRPS
	}
	if config.PerSenderBurst == 0 {
		config.PerSenderBurst = DefaultWorkflowBurst
	}
	if config.PerSenderRPS == 0 {
		config.PerSenderRPS = DefaultWorkflowRPS
	}
	return config
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

type metrics struct {
	handleDuration    metric.Int64Histogram
	awaitConnDuration metric.Int64Histogram
	method            string
}

func (m *metrics) recordSingleNodeRequestDuration(ctx context.Context, d time.Duration, status string, wid string) {
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("status", status),
		attribute.String("workflowID", wid),
		attribute.String("method", m.method),
	))
}

func (m *metrics) recordAwaitConnectionDuration(ctx context.Context, d time.Duration, wid string, gateway string, success bool) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.awaitConnDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("gateway", gateway),
		attribute.String("workflowID", wid),
		attribute.String("success", successStr),
		attribute.String("method", m.method),
	))
}

func newMetrics(method string) (*metrics, error) {
	h, err := beholder.GetMeter().Int64Histogram("platform_outgoing_connector_handler_single_node_request_duration_ms")
	if err != nil {
		return nil, err
	}

	a, err := beholder.GetMeter().Int64Histogram("platform_outgoing_connector_handler_await_conn_duration_ms")
	if err != nil {
		return nil, err
	}

	return &metrics{handleDuration: h, awaitConnDuration: a}, nil
}
