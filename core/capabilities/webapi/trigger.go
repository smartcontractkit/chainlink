package trigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webapicapabilities"
)

const defaultSendChannelBufferSize = 1000

const TriggerType = "web-trigger@1.0.0"

var webapiTriggerInfo = capabilities.MustNewCapabilityInfo(
	TriggerType,
	capabilities.CapabilityTypeTrigger,
	"A trigger to start workflow execution from a web api call",
)

type webapiTrigger struct {
	allowedSenders map[string]bool
	allowedTopics  map[string]bool
	ch             chan<- capabilities.TriggerResponse
	config         webapicap.TriggerConfig
	rateLimiter    *common.RateLimiter
	rawConfig      *values.Map
}

type TriggerConnectorHandler struct {
	services.StateMachine

	capabilities.CapabilityInfo
	capabilities.Validator[webapicap.TriggerConfig, struct{}, capabilities.TriggerResponse]
	connector           connector.GatewayConnector
	lggr                logger.Logger
	mu                  sync.Mutex
	RegisteredWorkflows map[string]webapiTrigger
}

var _ capabilities.TriggerCapability = (*TriggerConnectorHandler)(nil)
var _ services.Service = &TriggerConnectorHandler{}

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, lggr logger.Logger) (*TriggerConnectorHandler, error) {
	if connector == nil {
		return nil, errors.New("missing connector")
	}
	handler := &TriggerConnectorHandler{
		Validator:           capabilities.NewValidator[webapicap.TriggerConfig, struct{}, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: webapiTriggerInfo}),
		connector:           connector,
		RegisteredWorkflows: map[string]webapiTrigger{},
		lggr:                lggr.Named("WorkflowConnectorHandler"),
	}

	return handler, nil
}

// processTrigger iterates over each topic, checking against senders and rateLimits, then starting event processing and responding
func (h *TriggerConnectorHandler) processTrigger(ctx context.Context, gatewayID string, body *api.MessageBody, sender ethCommon.Address, payload webapicapabilities.TriggerRequestPayload) error {
	// Pass on the payload with the expectation that it's in an acceptable format for the executor
	wrappedPayload, err := values.WrapMap(payload)
	if err != nil {
		return fmt.Errorf("error wrapping payload %s", err)
	}
	topics := payload.Topics

	// empty topics is error for V1
	if len(topics) == 0 {
		return fmt.Errorf("empty Workflow Topics")
	}

	// workflows that have matched topics
	matchedWorkflows := 0
	// workflows that have matched topic and passed all checks
	fullyMatchedWorkflows := 0
	for _, trigger := range h.RegisteredWorkflows {
		for _, topic := range topics {
			if trigger.allowedTopics[topic] {
				matchedWorkflows++
				if !trigger.allowedSenders[sender.String()] {
					err = fmt.Errorf("unauthorized Sender %s, messageID %s", sender.String(), body.MessageId)
					h.lggr.Debugw(err.Error())
					continue
				}
				if !trigger.rateLimiter.Allow(body.Sender) {
					err = fmt.Errorf("request rate-limited for sender %s, messageID %s", sender.String(), body.MessageId)
					continue
				}
				fullyMatchedWorkflows++
				TriggerEventID := body.Sender + payload.TriggerEventID
				tr := capabilities.TriggerResponse{
					Event: capabilities.TriggerEvent{
						TriggerType: TriggerType,
						ID:          TriggerEventID,
						Outputs:     wrappedPayload,
					},
				}
				select {
				case <-ctx.Done():
					return nil
				case trigger.ch <- tr:
					// Sending n topics that match a workflow with n allowedTopics, can only be triggered once.
					break
				}
			}
		}
	}
	if matchedWorkflows == 0 {
		return fmt.Errorf("no Matching Workflow Topics")
	}

	if fullyMatchedWorkflows > 0 {
		return nil
	}
	return err
}

func (h *TriggerConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	// TODO: Validate Signature
	body := &msg.Body
	sender := ethCommon.HexToAddress(body.Sender)
	var payload webapicapabilities.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		err = h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s decoding payload", err.Error()).Error()})
		if err != nil {
			h.lggr.Errorw("error sending response", "err", err)
		}
		return
	}

	switch body.Method {
	case webapicapabilities.MethodWebAPITrigger:
		h.lggr.Debugw("Processing web api trigger", "gatewayid", gatewayID, "method", body.Method)
		resp := h.processTrigger(ctx, gatewayID, body, sender, payload)
		var response webapicapabilities.TriggerResponsePayload
		if resp == nil {
			response = webapicapabilities.TriggerResponsePayload{Status: "ACCEPTED"}
		} else {
			response = webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: resp.Error()}
			h.lggr.Errorw("Error processing trigger", "gatewayID", gatewayID, "body", body, "response", resp)
		}
		err = h.sendResponse(ctx, gatewayID, body, response)
		if err != nil {
			h.lggr.Errorw("Error sending response", "body", body, "response", response, "err", err)
		}
		return

	default:
		h.lggr.Errorw("unsupported method", "id", gatewayID, "method", body.Method)
		err = h.sendResponse(ctx, gatewayID, body, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("unsupported method %s", body.Method).Error()})
		if err != nil {
			h.lggr.Errorw("error sending response", "err", err)
		}
	}
}

// Periodically update the gateways with the state of the workflow triggers.
// Send the allowList for each workflow.

func (h *TriggerConnectorHandler) UpdateGateways(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var workflowConfigs = make(map[string]*values.Map)
	for triggerID, trigger := range h.RegisteredWorkflows {
		workflowConfigs[triggerID] = trigger.rawConfig
	}

	payloadJSON, err := json.Marshal(workflowConfigs)
	if err != nil {
		h.lggr.Errorw("error marshalling payload", "err", err)
		payloadJSON, _ = json.Marshal(webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s marshalling payload", err.Error()).Error()})
	}
	for gatewayID := range h.connector.GatewayIDs() {
		// convert gatewayID to string

		gatewayIDStr := strconv.Itoa(gatewayID)
		body := api.MessageBody{
			MessageId: types.RandomID().String(),
			DonId:     h.connector.DonID(),
			Method:    webapicapabilities.MethodWebAPITriggerUpdateMetadata,
			Receiver:  gatewayIDStr,
			Payload:   payloadJSON,
		}
		err = h.connector.SignAndSendToGateway(ctx, gatewayIDStr, &body)
		if err != nil {
			h.lggr.Errorw("error sending message", "err", err)
		}
	}
	return nil
}

func (h *TriggerConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	cfg := req.Config
	if cfg == nil {
		return nil, errors.New("config is required to register a web api trigger")
	}

	reqConfig, err := h.ValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	if len(reqConfig.AllowedSenders) == 0 {
		return nil, errors.New("allowedSenders must have at least 1 entry")
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, errBool := h.RegisteredWorkflows[req.TriggerID]
	if errBool {
		return nil, fmt.Errorf("triggerId %s already registered", req.TriggerID)
	}

	rateLimiterConfig := reqConfig.RateLimiter
	commonRateLimiter := common.RateLimiterConfig{
		GlobalRPS:      rateLimiterConfig.GlobalRPS,
		GlobalBurst:    int(rateLimiterConfig.GlobalBurst),
		PerSenderRPS:   rateLimiterConfig.PerSenderRPS,
		PerSenderBurst: int(rateLimiterConfig.PerSenderBurst),
	}

	rateLimiter, err := common.NewRateLimiter(commonRateLimiter)
	if err != nil {
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

	h.RegisteredWorkflows[req.TriggerID] = webapiTrigger{
		allowedTopics:  allowedTopicsMap,
		allowedSenders: allowedSendersMap,
		ch:             ch,
		rawConfig:      cfg,
		config:         *reqConfig,
		rateLimiter:    rateLimiter,
	}

	return ch, nil
}

func (h *TriggerConnectorHandler) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	workflow, ok := h.RegisteredWorkflows[req.TriggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", req.TriggerID)
	}

	close(workflow.ch)
	delete(h.RegisteredWorkflows, req.TriggerID)
	return nil
}

func (h *TriggerConnectorHandler) Start(ctx context.Context) error {
	return h.StartOnce("GatewayConnectorServiceWrapper", func() error {
		return h.connector.AddHandler([]string{"web_api_trigger"}, h)
	})
}
func (h *TriggerConnectorHandler) Close() error {
	return h.StopOnce("GatewayConnectorServiceWrapper", func() error {
		return nil
	})
}

func (h *TriggerConnectorHandler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *TriggerConnectorHandler) Name() string {
	return "WebAPITrigger"
}

func (h *TriggerConnectorHandler) sendResponse(ctx context.Context, gatewayID string, requestBody *api.MessageBody, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		h.lggr.Errorw("error marshalling payload", "err", err)
		payloadJSON, _ = json.Marshal(webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s marshalling payload", err.Error()).Error()})
	}

	body := api.MessageBody{
		MessageId: requestBody.MessageId,
		DonId:     requestBody.DonId,
		Method:    requestBody.Method,
		Receiver:  requestBody.Sender,
		Payload:   payloadJSON,
	}

	return h.connector.SignAndSendToGateway(ctx, gatewayID, &body)
}
