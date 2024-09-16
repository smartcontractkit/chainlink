package webapi

import (
	"context"
	"errors"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type workflowConnectorHandler struct {
	services.StateMachine

	capabilities.CapabilityInfo
	connector connector.GatewayConnector
	registry  core.CapabilitiesRegistry
	lggr      logger.Logger
}

var _ capabilities.TriggerCapability = (*workflowConnectorHandler)(nil)
var _ services.Service = &workflowConnectorHandler{}

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, lggr logger.Logger) (job.ServiceCtx, error) {
	lggr.Debugw("-----NewTrigger", "connector", connector)

	// TODO (CAPPL-22, CAPPL-24):
	//   - decode config
	//   - create an implementation of the capability API and add it to the Registry
	//   - create a handler and register it with Gateway Connector
	return &workflowConnectorHandler{
		connector: connector,
		registry:  registry,
		lggr:      lggr.Named("WorkflowConnectorHandler"),
	}, nil
	//   - manage trigger subscriptions
	//   - process incoming trigger events and related metadata
}

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (h *workflowConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	h.lggr.Debugw("-----handleGatewayMessage")

	body := &msg.Body
	fromAddr := ethCommon.HexToAddress(body.Sender)
	// TODO: apply allowlist and rate-limiting
	h.lggr.Debugw("handling gateway request", "id", gatewayId, "method", body.Method, "address", fromAddr)

	switch body.Method {
	case workflow.MethodAddWorkflow:
		// TODO: add a new workflow spec and return success/failure
		// we need access to Job ORM or whatever CLO uses to fully launch a new spec
		h.lggr.Debugw("added workflow spec", "payload", string(body.Payload))
		// response := Response{Success: true}
		// h.sendResponse(ctx, gatewayId, body, response)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayId, "method", body.Method)
	}
}

// Generate a Trigger Event and start processing in the Engine.
// TriggerEventID used internally by the Engine is a pair (sender, trigger_event_id). This is to protect against a case where two different authorized senders use the same event ID in their messages.

// Register a new trigger
// Can register triggers before the service is actively scheduling
func (s *workflowConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	s.lggr.Debugw("-----RegisterTrigger")

	if req.Config == nil {
		return nil, errors.New("config is required to register a cron trigger")
	}
	return nil, nil
}

func (s *workflowConnectorHandler) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	return nil
}

func (s *workflowConnectorHandler) Start(ctx context.Context) error {
	s.lggr.Debugw("-----Start workflowConnectorHandler")

	return s.StartOnce("GatewayConnectorServiceWrapper", func() error {
		s.lggr.Debugw("-----StartOnce call to GatewayConnectorServiceWrapper")

		s.connector.AddHandler([]string{"add_workflow"}, s)
		s.lggr.Debugw("-----StartOnce call to GatewayConnectorServiceWrapper addHandler")
		return s.registry.Add(ctx, s)
	})
}
func (s *workflowConnectorHandler) Close() error {
	return nil
}

func (s *workflowConnectorHandler) Ready() error {
	return nil
}

func (s *workflowConnectorHandler) HealthReport() map[string]error {
	return map[string]error{s.Name(): nil}
}

func (s *workflowConnectorHandler) Name() string {
	return "WebAPITrigger"
}
