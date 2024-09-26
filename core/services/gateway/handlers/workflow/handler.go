package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	MethodAddWorkflow = "add_workflow"
	// NOTE: more methods will go here: CRUD for workflow specs; HTTP trigger/action/target; etc.
	MethodWebAPITarget = "web_api_target"
)

type workflowHandler struct {
	donConfig      *config.DONConfig
	don            handlers.DON
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
	httpClient     network.HttpClient // for outgoing requests to users
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
}

var _ handlers.Handler = (*workflowHandler)(nil)

func NewWorkflowHandler(donConfig *config.DONConfig, don handlers.DON, httpClient network.HttpClient, lggr logger.Logger) (*workflowHandler, error) {
	return &workflowHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		httpClient:     httpClient,
		lggr:           lggr.Named("WorkflowHandler." + donConfig.DonId),
	}, nil
}

func (d *workflowHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()

	// TODO: apply allowlist and rate-limiting here.
	if msg.Body.Method != MethodAddWorkflow {
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

func (d *workflowHandler) handleWebAPITargetMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	var targetPayload TargetRequestPayload
	err := json.Unmarshal(msg.Body.Payload, &targetPayload)
	if err != nil {
		return err
	}
	// send message to target
	req := network.HttpRequest{
		Method:     targetPayload.Method,
		URL:        targetPayload.URL,
		Headers:    targetPayload.Headers,
		Body:       targetPayload.Body,
		Timeout:    time.Duration(targetPayload.TimeoutMs) * time.Millisecond,
		RetryCount: targetPayload.RetryCount,
	}
	resp, err := d.httpClient.Send(ctx, req)
	if err != nil {
		return err
	}
	respPayload := TargetResponsePayload{
		StatusCode: uint8(resp.StatusCode),
		Headers:    resp.Headers,
		Body:       resp.Body,
	}
	payload, err := json.Marshal(respPayload)
	if err != nil {
		return err
	}

	respMsg := &api.Message{
		Body: api.MessageBody{
			MessageId: msg.Body.MessageId,
			Method:    msg.Body.Method,
			DonId:     msg.Body.DonId,
			Payload:   payload,
		},
	}

	// TODO: check if respMsg is needs to be signed? may not be needed since WS connection between gateway and node are already verified
	err = d.don.SendToNode(ctx, nodeAddr, respMsg)
	if err != nil {
		return err
	}
	return nil
}

func (d *workflowHandler) handleNodeResponse(ctx context.Context, msg *api.Message) error {
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

func (d *workflowHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	// signature verification not needed because websocket connection is verified
	// TODO: rate limiting
	switch msg.Body.Method {
	case MethodWebAPITarget:
		return d.handleWebAPITargetMessage(ctx, msg, nodeAddr)
	case MethodAddWorkflow:
		return d.handleNodeResponse(ctx, msg)
	default:
		return fmt.Errorf("unsupported method: %s", msg.Body.Method)
	}
}

func (d *workflowHandler) Start(context.Context) error {
	return nil
}

func (d *workflowHandler) Close() error {
	return nil
}
