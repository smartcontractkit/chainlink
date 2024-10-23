package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
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

func NewOutgoingConnectorHandler(gc connector.GatewayConnector, config ServiceConfig, method string, lgger logger.Logger) (*OutgoingConnectorHandler, error) {
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
	case capabilities.MethodWebAPITarget, capabilities.MethodComputeAction:
		body := &msg.Body
		var payload capabilities.Response
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

func validMethod(method string) bool {
	switch method {
	case capabilities.MethodWebAPITarget, capabilities.MethodComputeAction:
		return true
	default:
		return false
	}
}
