package target

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
)

var _ connector.GatewayConnectorHandler = &Handler{}

type Handler struct {
	gc            connector.GatewayConnector
	lggr          logger.Logger
	responseChs   map[string]chan *api.Message
	responseChsMu *sync.Mutex
}

func NewHandler(gc connector.GatewayConnector, responseChs map[string]chan *api.Message, responseChsMu *sync.Mutex, lgger logger.Logger) (*Handler, error) {
	return &Handler{
		gc:            gc,
		responseChs:   responseChs,
		responseChsMu: responseChsMu,
		lggr:          lgger,
	}, nil
}

func (c *Handler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	body := &msg.Body
	l := logger.With(c.lggr, "id", gatewayId, "method", body.Method, "messageID", msg.Body.MessageId)
	l.Debugw("handling gateway request")

	// TODO: check allowlist and enforce rate-limiting
	switch body.Method {
	case workflow.MethodWebAPITarget:
		var payload workflow.TargetResponsePayload
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

func (l *Handler) Start(ctx context.Context) error {
	return l.gc.AddHandler([]string{workflow.MethodWebAPITarget}, l)
}

func (l *Handler) Close() error {
	return nil
}
