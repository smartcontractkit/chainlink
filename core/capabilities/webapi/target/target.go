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
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webapicapabilities"
)

const ID = "web-api-target@1.0.0"

var _ capabilities.TargetCapability = &Capability{}

var capabilityInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTarget,
	"A target that sends HTTP requests to external clients via the Chainlink Gateway.",
)

const (
	DefaultDeliveryMode = SingleNode
	DefaultHTTPMethod   = "GET"
	DefaultTimeoutMs    = 30000
	MaxTimeoutMs        = 600000
)

// Capability is a target capability that sends HTTP requests to external clients via the Chainlink Gateway.
type Capability struct {
	capabilityInfo   capabilities.CapabilityInfo
	connectorHandler *ConnectorHandler
	lggr             logger.Logger
	registry         core.CapabilitiesRegistry
	config           ServiceConfig
}

func NewCapability(config ServiceConfig, registry core.CapabilitiesRegistry, connectorHandler *ConnectorHandler, lggr logger.Logger) (*Capability, error) {
	return &Capability{
		capabilityInfo:   capabilityInfo,
		config:           config,
		registry:         registry,
		connectorHandler: connectorHandler,
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
		req.Metadata.WorkflowExecutionID,
		webapicapabilities.MethodWebAPITarget,
	}
	return strings.Join(messageID, "/"), nil
}

// defaultIfNil is a helper function to handle nil pointers and provide default values
func defaultIfNil[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}

func getPayload(input webapicap.TargetPayload, cfg webapicap.TargetConfig) (webapicapabilities.TargetRequestPayload, error) {
	method := defaultIfNil(input.Method, DefaultHTTPMethod)
	body := defaultIfNil(input.Body, "")
	timeoutMs := defaultIfNil(cfg.TimeoutMs, DefaultTimeoutMs)
	if timeoutMs < 0 || timeoutMs > MaxTimeoutMs {
		return webapicapabilities.TargetRequestPayload{}, fmt.Errorf("timeoutMs must be between 0 and %d", MaxTimeoutMs)
	}

	return webapicapabilities.TargetRequestPayload{
		URL:       input.Url,
		Method:    method,
		Headers:   input.Headers,
		Body:      []byte(body),
		TimeoutMs: timeoutMs,
	}, nil
}

func (c *Capability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	c.lggr.Debugw("executing http target", "capabilityRequest", req)

	var input webapicap.TargetPayload
	err := req.Inputs.UnwrapTo(&input)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	var workflowCfg webapicap.TargetConfig
	err = req.Config.UnwrapTo(&workflowCfg)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	messageID, err := getMessageID(req)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	payload, err := getPayload(input, workflowCfg)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	// Default to SingleNode delivery mode
	deliveryMode := defaultIfNil(workflowCfg.DeliveryMode, SingleNode)

	switch deliveryMode {
	case SingleNode:
		// blocking call to handle single node request. waits for response from gateway
		resp, err := c.connectorHandler.HandleSingleNodeRequest(ctx, messageID, payloadBytes)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		c.lggr.Debugw("received gateway response", "resp", resp)
		var payload webapicapabilities.TargetResponsePayload
		err = json.Unmarshal(resp.Body.Payload, &payload)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}

		// TODO: check target response format and fields CM-473
		values, err := values.NewMap(map[string]any{
			"statusCode": payload.StatusCode,
			"headers":    payload.Headers,
			"body":       payload.Body,
		})
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		return capabilities.CapabilityResponse{
			Value: values,
		}, nil
	default:
		return capabilities.CapabilityResponse{}, fmt.Errorf("unsupported delivery mode: %v", workflowCfg.DeliveryMode)
	}
}

func (c *Capability) RegisterToWorkflow(ctx context.Context, req capabilities.RegisterToWorkflowRequest) error {
	// Workflow engine guarantees registration requests are valid
	// TODO: handle retry configuration CM-472
	return nil
}

func (c *Capability) UnregisterFromWorkflow(ctx context.Context, req capabilities.UnregisterFromWorkflowRequest) error {
	// Workflow engine guarantees deregistration requests are valid
	return nil
}
