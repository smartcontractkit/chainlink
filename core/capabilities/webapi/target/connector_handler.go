package target

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webcapabilities"
)

var _ connector.GatewayConnectorHandler = &ConnectorHandler{}

type ConnectorHandler struct {
	gc            connector.GatewayConnector
	lggr          logger.Logger
	responseChs   map[string]chan *api.Message
	responseChsMu sync.Mutex
	rateLimiter   *common.RateLimiter
}

func NewConnectorHandler(gc connector.GatewayConnector, config Config, lgger logger.Logger) (*ConnectorHandler, error) {
	rateLimiter, err := common.NewRateLimiter(config.RateLimiter)
	if err != nil {
		return nil, err
	}
	responseChs := make(map[string]chan *api.Message)
	return &ConnectorHandler{
		gc:            gc,
		responseChs:   responseChs,
		responseChsMu: sync.Mutex{},
		rateLimiter:   rateLimiter,
		lggr:          lgger,
	}, nil
}

// HandleSingleNodeRequest sends a request to first available gateway node and blocks until response is received
// TODO: handle retries and timeouts
func (c *ConnectorHandler) HandleSingleNodeRequest(ctx context.Context, msg *api.MessageBody) (*api.Message, error) {
	ch := make(chan *api.Message, 1)
	c.responseChsMu.Lock()
	c.responseChs[msg.MessageId] = ch
	c.responseChsMu.Unlock()
	l := logger.With(c.lggr, "messageId", msg.MessageId)
	l.Debugw("sending request to gateway")

	err := c.gc.SendToAvailableGateway(ctx, msg)
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

func (c *ConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	body := &msg.Body
	l := logger.With(c.lggr, "gatewayID", gatewayID, "method", body.Method, "messageID", msg.Body.MessageId)
	if !c.rateLimiter.Allow(body.Sender) {
		c.lggr.Errorw("request rate-limited")
		return
	}
	l.Debugw("handling gateway request")
	switch body.Method {
	case webcapabilities.MethodWebAPITarget:
		var payload webcapabilities.TargetResponsePayload
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
		ch <- msg
	default:
		l.Errorw("unsupported method")
	}
}

func (c *ConnectorHandler) Start(ctx context.Context) error {
	return c.gc.AddHandler([]string{webcapabilities.MethodWebAPITarget}, c)
}

func (c *ConnectorHandler) Close() error {
	return nil
}
