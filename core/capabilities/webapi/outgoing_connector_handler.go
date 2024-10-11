package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

var _ connector.GatewayConnectorHandler = &OutgoingConnectorHandler{}

type OutgoingConnectorHandler struct {
	gc            connector.GatewayConnector
	method        string
	lggr          logger.Logger
	responseChs   map[string]chan *api.Message
	responseChsMu sync.Mutex
	rateLimiter   *common.RateLimiter
}

// Config is the configuration for OutgoingConnectorHandler.
// Currently used by the WebApi Target and Compute Action capability & handler
// TODO: handle retry configurations here CM-472
// Note that workflow executions have their own internal timeouts and retries set by the user
// that are separate from this configuration
type Config struct {
	RateLimiter common.RateLimiterConfig `toml:"rateLimiter"`
}

func NewOutgoingConnectorHandler(gc connector.GatewayConnector, config Config, method string, lgger logger.Logger) (*OutgoingConnectorHandler, error) {
	rateLimiter, err := common.NewRateLimiter(config.RateLimiter)
	if err != nil {
		return nil, err
	}

	if !validMethod(method) {
		return nil, fmt.Errorf("invalid outgoing connector handler method: %s", method)
	}

	responseChs := make(map[string]chan *api.Message)
	return &OutgoingConnectorHandler{
		gc:            gc,
		method:        method,
		responseChs:   responseChs,
		responseChsMu: sync.Mutex{},
		rateLimiter:   rateLimiter,
		lggr:          lgger,
	}, nil
}

// HandleSingleNodeRequest sends a request to first available gateway node and blocks until response is received
// TODO: handle retries and timeouts
func (c *OutgoingConnectorHandler) HandleSingleNodeRequest(ctx context.Context, messageID string, payload []byte) (*api.Message, error) {
	ch := make(chan *api.Message, 1)
	c.responseChsMu.Lock()
	c.responseChs[messageID] = ch
	c.responseChsMu.Unlock()
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

	err := c.gc.SignAndSendToGateway(ctx, gatewayIDs[0], body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to gateway")
	}

	select {
	case resp := <-ch:
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
	case capabilities.MethodWebAPITarget:
		var payload capabilities.TargetResponsePayload
		err := json.Unmarshal(body.Payload, &payload)
		if err != nil {
			l.Errorw("failed to unmarshal payload", "err", err)
			return
		}
		c.responseChsMu.Lock()
		defer c.responseChsMu.Unlock()
		ch, ok := c.responseChs[body.MessageId]
		if !ok {
			l.Errorw("no response channel found")
			return
		}
		select {
		case ch <- msg:
			delete(c.responseChs, body.MessageId)
		case <-ctx.Done():
			return
		}
	case capabilities.MethodComputeAction:
		var payload sdk.FetchResponse
		err := json.Unmarshal(body.Payload, &payload)
		if err != nil {
			l.Errorw("failed to unmarshal payload", "err", err)
			return
		}
		c.responseChsMu.Lock()
		defer c.responseChsMu.Unlock()
		ch, ok := c.responseChs[body.MessageId]
		if !ok {
			l.Errorw("no response channel found")
			return
		}
		select {
		case ch <- msg:
			delete(c.responseChs, body.MessageId)
		case <-ctx.Done():
			return
		}
	default:
		l.Errorw("unsupported method")
	}
}

func (c *OutgoingConnectorHandler) Start(ctx context.Context) error {
	return c.gc.AddHandler([]string{c.method}, c)
}

func (c *OutgoingConnectorHandler) Close() error {
	return nil
}

func (c *OutgoingConnectorHandler) CreateFetcher(workflowID, workflowExecutionID string) func(req *wasmpb.FetchRequest) (*wasmpb.FetchResponse, error) {
	return func(req *wasmpb.FetchRequest) (*wasmpb.FetchResponse, error) {
		if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
			return nil, fmt.Errorf("workflow ID %q is invalid: %w", workflowID, err)
		}
		if err := validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
			return nil, fmt.Errorf("workflow execution ID %q is invalid: %w", workflowExecutionID, err)
		}

		messageID := strings.Join([]string{
			workflowID,
			workflowExecutionID,
			ghcapabilities.MethodComputeAction,
		}, "/")

		fields := req.Headers.GetFields()
		headersReq := make(map[string]any, len(fields))
		for k, v := range fields {
			headersReq[k] = v
		}

		payloadBytes, err := json.Marshal(sdk.FetchRequest{
			URL:       req.Url,
			Method:    req.Method,
			Headers:   headersReq,
			Body:      req.Body,
			TimeoutMs: req.TimeoutMs,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal fetch request: %w", err)
		}

		resp, err := c.HandleSingleNodeRequest(context.Background(), messageID, payloadBytes)
		if err != nil {
			return nil, err
		}

		c.lggr.Debugw("received gateway response", "resp", resp)
		var response wasmpb.FetchResponse
		err = json.Unmarshal(resp.Body.Payload, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal fetch response: %w", err)
		}
		return &response, nil
	}
}

func validMethod(method string) bool {
	validMethods := []string{capabilities.MethodWebAPITarget, capabilities.MethodWebAPITrigger, capabilities.MethodComputeAction}

	for _, vm := range validMethods {
		if vm == method {
			return true
		}
	}

	return false
}
