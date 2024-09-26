package target

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
)

var _ capabilities.TargetCapability = &Capability{}

type Capability struct {
	gc            connector.GatewayConnector
	lggr          logger.Logger
	registry      core.CapabilitiesRegistry
	config        Config
	responseChs   map[string]chan *api.Message
	responseChsMu *sync.Mutex
}

func NewCapability(config Config, registry core.CapabilitiesRegistry, gc connector.GatewayConnector, lggr logger.Logger,
	responseChs map[string]chan *api.Message, responseChsMu *sync.Mutex) (*Capability, error) {
	return &Capability{
		responseChs:   make(map[string]chan *api.Message),
		responseChsMu: responseChsMu,
	}, nil
}

func (c *Capability) Start(ctx context.Context) error {
	return c.registry.Add(ctx, c)
}

func (c *Capability) Close() error {
	return nil
}

func (c *Capability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.CapabilityInfo{}, nil
}

func getMessageID(req capabilities.CapabilityRequest) (string, error) {
	if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
		return "", fmt.Errorf("workflow ID is invalid: %w", err)
	}
	if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowExecutionID); err != nil {
		return "", fmt.Errorf("workflow execution ID is invalid: %w", err)
	}
	messageID := []string{
		req.Metadata.WorkflowID,
		req.Metadata.WorkflowExecutionID,
		workflow.MethodWebAPITarget,
	}
	return strings.Join(messageID, "/"), nil
}

func (c *Capability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	c.lggr.Debugw("executing http target", "capabilityRequest", req)

	var input WorkflowInput
	err := req.Inputs.UnwrapTo(&input)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	var workflowCfg WorkflowConfig
	err = req.Config.UnwrapTo(&workflowCfg)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	payload := workflow.TargetRequestPayload{
		URL:        input.URL,
		Method:     input.Method,
		Headers:    input.Headers,
		Body:       []byte(input.Body),
		TimeoutMs:  workflowCfg.TimeoutMs,
		RetryCount: workflowCfg.RetryCount,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	messageID, err := getMessageID(req)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	gatewayReq := api.SendRequest{
		MessageId: messageID,
		Method:    workflow.MethodWebAPITarget,
		Payload:   payloadJson,
	}

	switch workflowCfg.Schedule {
	case RoundRobin:
		resp, err := c.handleRoundRobinRequest(ctx, gatewayReq, workflowCfg.RetryCount)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		return resp, nil
	default:
		return capabilities.CapabilityResponse{}, fmt.Errorf("unsupported schedule: %v", workflowCfg.Schedule)
	}
}

func (c *Capability) handleRoundRobinRequest(ctx context.Context, req api.SendRequest, retryCount uint8) (capabilities.CapabilityResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutMs)*time.Millisecond)
	defer cancel()
	ch := make(chan *api.Message, 1)
	c.responseChsMu.Lock()
	c.responseChs[req.MessageId] = ch
	c.responseChsMu.Unlock()

	success := false
	for i := uint8(0); i < c.config.RetryCount; i++ {
		l := logger.With(c.lggr, "messageId", req.MessageId, "attempt", i+1)
		l.Debugw("sending request to gateway")
		err := c.gc.SendToGatewayRoundRobin(ctx, req)
		if err != nil {
			l.Warnw("failed to send request to gateway", "err", err)
			continue
		}
		success = true
	}

	if !success {
		return capabilities.CapabilityResponse{}, fmt.Errorf("failed to send request to gateway for messageID: %s", req.MessageId)
	}

	select {
	case resp := <-ch:
		var payload workflow.TargetResponsePayload
		err := json.Unmarshal(resp.Body.Payload, &payload)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		// TODO: check target response format and fields
		values, err := values.NewMap(map[string]any{
			"statusCode": payload.StatusCode,
			"headers":    payload.Headers,
			"body":       string(payload.Body),
		})
		return capabilities.CapabilityResponse{
			Value: values,
		}, nil
	case <-ctx.Done():
		return capabilities.CapabilityResponse{}, ctx.Err()
	}
}

func (c *Capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	// TODO: validate registration request
	return nil
}

func (c *Capability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	// TODO: what should this do? stop processing any incoming requests to the target?
	return nil
}
