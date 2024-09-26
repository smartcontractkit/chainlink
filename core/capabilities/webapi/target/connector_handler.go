package target

import (
	"context"
	"encoding/json"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webapicapabilities"
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
func (c *ConnectorHandler) HandleSingleNodeRequest(ctx context.Context, messageID string, payload []byte) (*api.Message, error) {
	ch := make(chan *api.Message, 1)
	c.responseChsMu.Lock()
	c.responseChs[messageID] = ch
	c.responseChsMu.Unlock()
	l := logger.With(c.lggr, "messageID", messageID)
	l.Debugw("sending request to gateway")

	body := &api.MessageBody{
		MessageId: messageID,
		DonId:     c.gc.DonID(),
		Method:    webapicapabilities.MethodWebAPITarget,
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

func (c *ConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
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
	case webapicapabilities.MethodWebAPITarget:
		var payload webapicapabilities.TargetResponsePayload
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

func (c *ConnectorHandler) Start(ctx context.Context) error {
	return c.gc.AddHandler([]string{webapicapabilities.MethodWebAPITarget}, c)
}

func (c *ConnectorHandler) Close() error {
	return nil
}
