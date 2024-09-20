package workflow

// TODO: reconcile with Jin's PR.
import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

const (
	MethodWebAPITrigger = "web_trigger"
)

type TriggerRequestPayload struct {
	TriggerId      string     `json:"trigger_id"`
	TriggerEventId string     `json:"trigger_event_id"`
	Timestamp      int64      `json:"timestamp"`
	Topics         []string   `json:"topics"`
	Params         values.Map `json:"params"`
}

type workflowHandler struct {
	donConfig      *config.DONConfig
	don            handlers.DON
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
}

var _ handlers.Handler = (*workflowHandler)(nil)

func NewWorkflowHandler(donConfig *config.DONConfig, don handlers.DON, lggr logger.Logger) (*workflowHandler, error) {
	return &workflowHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           lggr.Named("WorkflowHandler." + donConfig.DonId),
	}, nil
}

// TODO: how is the HTTP response with success, status etc. generated?
func (d *workflowHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()
	body := msg.Body
	var payload TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		d.lggr.Errorw("error decoding payload", "err", err)
		return err
	}

	currentTime := time.Now()
	// TODO: check against h.config.MaxAllowedMessageAgeSec
	if currentTime.Unix()-3000 > payload.Timestamp {
		// TODO: fix up with error handling update in design doc
		return fmt.Errorf("message too stale")
	}
	// TODO: apply allowlist and rate-limiting here.
	if msg.Body.Method != MethodWebAPITrigger {
		d.lggr.Errorw("unsupported method", "method", body.Method)
		return fmt.Errorf("unsupported method")
	}

	// Send to all nodes.
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
	// TODO: CAPPL-21 send correct response.
	return err
}

func (d *workflowHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	d.mu.Lock()
	savedCb, found := d.savedCallbacks[msg.Body.MessageId]
	delete(d.savedCallbacks, msg.Body.MessageId)
	d.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		// TODO: in practice, we should wait for at least 2F+1 nodes to respond and then return an aggregated response
		// back to the user.
		savedCb.callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (d *workflowHandler) Start(context.Context) error {
	return nil
}

func (d *workflowHandler) Close() error {
	return nil
}
