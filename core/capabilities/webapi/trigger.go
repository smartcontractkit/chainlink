package trigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webapicapabilities"
)

const defaultSendChannelBufferSize = 1000

var webapiTriggerInfo = capabilities.MustNewCapabilityInfo(
	webapicapabilities.TriggerType,
	capabilities.CapabilityTypeTrigger,
	"A trigger to start workflow execution from a web api call",
)

type Input struct {
}
type TriggerConfig struct {
	AllowedSenders []string                 `toml:"allowedSenders"`
	AllowedTopics  []string                 `toml:"allowedTopics"`
	RateLimiter    common.RateLimiterConfig `toml:"rateLimiter"`
	// RequiredParams is advisory to the consumer, it is not enforced.
	RequiredParams []string `toml:"requiredParams"`
}

type webapiTrigger struct {
	allowedSenders map[string]bool
	allowedTopics  map[string]bool
	ch             chan<- capabilities.TriggerResponse
	config         TriggerConfig
	rateLimiter    *common.RateLimiter
}

type triggerConnectorHandler struct {
	services.StateMachine

	capabilities.CapabilityInfo
	capabilities.Validator[TriggerConfig, Input, capabilities.TriggerResponse]
	connector           connector.GatewayConnector
	lggr                logger.Logger
	mu                  sync.Mutex
	registeredWorkflows map[string]webapiTrigger
}

var _ capabilities.TriggerCapability = (*triggerConnectorHandler)(nil)
var _ services.Service = &triggerConnectorHandler{}

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, lggr logger.Logger) (*triggerConnectorHandler, error) {
	if connector == nil {
		return nil, errors.New("missing connector")
	}
	handler := &triggerConnectorHandler{
		Validator:           capabilities.NewValidator[TriggerConfig, Input, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: webapiTriggerInfo}),
		connector:           connector,
		registeredWorkflows: map[string]webapiTrigger{},
		lggr:                lggr.Named("WorkflowConnectorHandler"),
	}

	return handler, nil
}

// processTrigger iterates over each topic, checking against senders and rateLimits, then starting event processing and responding
func (h *triggerConnectorHandler) processTrigger(ctx context.Context, gatewayID string, body *api.MessageBody, sender ethCommon.Address, payload webapicapabilities.TriggerRequestPayload) {
	// Pass on the payload with the expectation that it's in an acceptable format for the executor
	wrappedPayload, err := values.WrapMap(payload)
	if err != nil {
		h.lggr.Errorw("Error wrapping payload")
		h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "Error wrapping payload"})
		return
	}
	matchedWorkflows := 0
	var response webapicapabilities.TriggerResponsePayload

	for _, trigger := range h.registeredWorkflows {
		topics := payload.Topics
		// empty topics means all topics
		if len(topics) == 0 {
			for k := range trigger.allowedTopics {
				topics = append(topics, k)
			}
		}

		for _, topic := range topics {
			if trigger.allowedTopics[topic] {
				matchedWorkflows++

				if !trigger.allowedSenders[sender.String()] {
					h.lggr.Errorw("Unauthorized Sender", "sender", sender.String(), "messageID", body.MessageId)
					h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "Unauthorized Sender"})
					return
				}
				if !trigger.rateLimiter.Allow(body.Sender) {
					h.lggr.Errorw("request rate-limited", sender.String(), "messageID", body.MessageId)
					h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "request rate-limited"})
					return
				}

				TriggerEventID := body.Sender + payload.TriggerEventID
				tr := capabilities.TriggerResponse{
					Event: capabilities.TriggerEvent{
						TriggerType: webapicapabilities.TriggerType,
						ID:          TriggerEventID,
						Outputs:     wrappedPayload,
					},
				}

				trigger.ch <- tr
				response = webapicapabilities.TriggerResponsePayload{Status: "ACCEPTED"}
				// Sending n topics that match a workflow with n allowedTopics, can only be triggered once.
				break
			}
		}
	}
	if matchedWorkflows == 0 {
		h.lggr.Errorw("No Matching Workflow Topics")
		response = webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "No Matching Workflow Topics"}
	}
	err = h.sendResponse(ctx, gatewayID, body, response)
	if err != nil {
		h.lggr.Errorw("Error sending response", "body", body, "response", response)
	}
}

func (h *triggerConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	// TODO: Validate Signature
	body := &msg.Body
	sender := ethCommon.HexToAddress(body.Sender)
	var payload webapicapabilities.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s decoding payload", err.Error()).Error()})
		return
	}

	switch body.Method {
	case webapicapabilities.MethodWebAPITrigger:
		h.processTrigger(ctx, gatewayID, body, sender, payload)
		return

	default:
		h.lggr.Errorw("unsupported method", "id", gatewayID, "method", body.Method)
		h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("unsupported method %s", body.Method).Error()})
	}
}

func (h *triggerConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	cfg := req.Config
	if cfg == nil {
		h.lggr.Errorw("config is required to register a web api trigger")
		return nil, errors.New("config is required to register a web api trigger")
	}

	reqConfig, err := h.ValidateConfig(cfg)
	if err != nil {
		h.lggr.Errorw("error unwrapping config", "err", err)
		return nil, err
	}

	if len(reqConfig.AllowedSenders) == 0 {
		h.lggr.Errorw("allowedSenders must have at least 1 entry")
		return nil, errors.New("allowedSenders must have at least 1 entry")
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, errBool := h.registeredWorkflows[req.TriggerID]
	if errBool {
		h.lggr.Errorf("triggerId %s already registered", req.TriggerID)
		return nil, fmt.Errorf("triggerId %s already registered", req.TriggerID)
	}

	rateLimiter, err := common.NewRateLimiter(reqConfig.RateLimiter)

	if err != nil {
		h.lggr.Errorw("error creating RateLimiter", "err", err, "RateLimiter config", reqConfig.RateLimiter)
		return nil, err
	}

	allowedSendersMap := map[string]bool{}
	for _, k := range reqConfig.AllowedSenders {
		allowedSendersMap[k] = true
	}

	allowedTopicsMap := map[string]bool{}
	for _, k := range reqConfig.AllowedTopics {
		allowedTopicsMap[k] = true
	}

	ch := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)

	h.registeredWorkflows[req.TriggerID] = webapiTrigger{
		allowedTopics:  allowedTopicsMap,
		allowedSenders: allowedSendersMap,
		ch:             ch,
		config:         *reqConfig,
		rateLimiter:    rateLimiter,
	}

	return ch, nil
}

func (h *triggerConnectorHandler) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	workflow, ok := h.registeredWorkflows[req.TriggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", req.TriggerID)
	}

	close(workflow.ch)
	delete(h.registeredWorkflows, req.TriggerID)
	return nil
}

func (h *triggerConnectorHandler) Start(ctx context.Context) error {
	return h.StartOnce("GatewayConnectorServiceWrapper", func() error {
		return h.connector.AddHandler([]string{"web_trigger"}, h)
	})
}
func (h *triggerConnectorHandler) Close() error {
	return h.StopOnce("GatewayConnectorServiceWrapper", func() error {
		return nil
	})
}

func (h *triggerConnectorHandler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *triggerConnectorHandler) Name() string {
	return "WebAPITrigger"
}

func (h *triggerConnectorHandler) sendResponse(ctx context.Context, gatewayID string, requestBody *api.MessageBody, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		h.lggr.Errorw("error marshalling payload", "err", err)
		payloadJSON, _ = json.Marshal(webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s marshalling payload", err.Error()).Error()})
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: requestBody.MessageId,
			DonId:     requestBody.DonId,
			Method:    requestBody.Method,
			Receiver:  requestBody.Sender,
			Payload:   payloadJSON,
		},
	}

	return h.connector.SendToGateway(ctx, gatewayID, msg)
}
