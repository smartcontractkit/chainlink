package workflow

// TODO: reconcile with Jin's PR.
import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

const (
	MethodWebAPITrigger = "web_trigger"
)

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
	lggr.Debugw("-------HNewWorkflowHandler")

	return &workflowHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           lggr.Named("WorkflowHandler." + donConfig.DonId),
	}, nil
}

func (d *workflowHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	d.lggr.Debugw("-------HandleUserMessage", "method", msg.Body.Method)

	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()

	// TODO: apply allowlist and rate-limiting here.
	if msg.Body.Method != MethodWebAPITrigger {
		d.lggr.Errorw("unsupported method", "method", msg.Body.Method)
		return fmt.Errorf("unsupported method")
	}

	var err error
	// Send to all nodes.
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
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
