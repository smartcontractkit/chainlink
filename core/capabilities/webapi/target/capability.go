package target

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webcapabilities"
)

const ID = "web-api-target@1.0.0"

var _ capabilities.TargetCapability = &Capability{}

var capabilityInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTarget,
	"A target that sends HTTP requests to external clients via the Chainlink Gateway.",
)

type Capability struct {
	capabilityInfo   capabilities.CapabilityInfo
	connectorHandler *ConnectorHandler
	lggr             logger.Logger
	registry         core.CapabilitiesRegistry
	config           Config
	activeWorkflows  map[string]struct{} // tracks registered workflows
}

func NewCapability(config Config, registry core.CapabilitiesRegistry, connectorHandler *ConnectorHandler, lggr logger.Logger) (*Capability, error) {
	return &Capability{
		capabilityInfo:   capabilityInfo,
		config:           config,
		registry:         registry,
		connectorHandler: connectorHandler,
		activeWorkflows:  make(map[string]struct{}),
		lggr:             lggr,
	}, nil
}

func (c *Capability) Start(ctx context.Context) error {
	return c.registry.Add(ctx, c)
}

func (c *Capability) Close() error {
	return nil
}

func (c *Capability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilityInfo, nil
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
		webcapabilities.MethodWebAPITarget,
	}
	return strings.Join(messageID, "/"), nil
}

func (c *Capability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	c.lggr.Debugw("executing http target", "capabilityRequest", req)

	var input Input
	err := req.Inputs.UnwrapTo(&input)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	var workflowCfg WorkflowConfig
	err = req.Config.UnwrapTo(&workflowCfg)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	messageID, err := getMessageID(req)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	if _, ok := c.activeWorkflows[req.Metadata.WorkflowID]; !ok {
		return capabilities.CapabilityResponse{}, fmt.Errorf("workflow is not registered: %v", req.Metadata.WorkflowID)
	}

	payload := webcapabilities.TargetRequestPayload{
		URL:        input.URL,
		Method:     input.Method,
		Headers:    input.Headers,
		Body:       []byte(input.Body),
		TimeoutMs:  workflowCfg.TimeoutMs,
		RetryCount: workflowCfg.RetryCount,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	gatewayReq := &api.MessageBody{
		MessageId: messageID,
		Method:    webcapabilities.MethodWebAPITarget,
		Payload:   payloadBytes,
	}

	switch workflowCfg.Schedule {
	case SingleNode:
		// blocking call to handle single node request. waits for response from gateway
		resp, err := c.connectorHandler.HandleSingleNodeRequest(ctx, gatewayReq)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		c.lggr.Debugw("received gateway response", "resp", resp)
		var payload webcapabilities.TargetResponsePayload
		err = json.Unmarshal(resp.Body.Payload, &payload)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}

		// TODO: check target response format and fields
		values, err := values.NewMap(map[string]any{
			"statusCode": payload.StatusCode,
			"headers":    payload.Headers,
			"body":       string(payload.Body),
		})
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		return capabilities.CapabilityResponse{
			Value: values,
		}, nil
	default:
		return capabilities.CapabilityResponse{}, fmt.Errorf("unsupported schedule: %v", workflowCfg.Schedule)
	}
}

func (c *Capability) RegisterToWorkflow(ctx context.Context, req capabilities.RegisterToWorkflowRequest) error {
	if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
		return fmt.Errorf("workflow ID is invalid: %w", err)
	}
	c.activeWorkflows[req.Metadata.WorkflowID] = struct{}{}
	return nil
}

func (c *Capability) UnregisterFromWorkflow(ctx context.Context, req capabilities.UnregisterFromWorkflowRequest) error {
	// TODO: is best-effort removal of workflow sufficient here? or should it error if workflow is not found for some reason?
	if _, ok := c.activeWorkflows[req.Metadata.WorkflowID]; !ok {
		c.lggr.Warnw("workflow not found", "workflowID", req.Metadata.WorkflowID)
	} else {
		delete(c.activeWorkflows, req.Metadata.WorkflowID)
	}
	return nil
}
