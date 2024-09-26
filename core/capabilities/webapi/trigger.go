package webapi

import (
	"context"
	"crypto/ecdsa"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
)

const defaultSendChannelBufferSize = 1000

const triggerType = "web-trigger@1.0.0"

var webapiTriggerInfo = capabilities.MustNewCapabilityInfo(
	triggerType,
	capabilities.CapabilityTypeTrigger,
	"A trigger to start workflow execution from a web api call",
)

type Input struct {
}
type TriggerConfig struct {
	AllowedSenders []string `toml:"allowedSenders"`
	AllowedTopics  []string `toml:"allowedTopics"`
	// RateLimiter    common.RateLimiterConfig `toml:"rateLimiter"`
	RateLimiter    *values.Map `toml:"rateLimiter"`
	RequiredParams []string    `toml:"requiredParams"`
}

type Response struct {
	ErrorMessage string `json:"error_message,omitempty"`
	// ERROR, ACCEPTED, PENDING, COMPLETED
	Status string `json:"status"`
}

type webapiTrigger struct {
	allowedSendersMap map[string]bool
	allowedTopicsMap  map[string]bool
	ch                chan<- capabilities.TriggerResponse
	config            TriggerConfig
	rateLimiter       *common.RateLimiter
}

type triggerConnectorHandler struct {
	services.StateMachine

	capabilities.CapabilityInfo
	capabilities.Validator[TriggerConfig, Input, capabilities.TriggerResponse]
	connector           connector.GatewayConnector
	lggr                logger.Logger
	mu                  sync.Mutex
	registeredWorkflows map[string]webapiTrigger
	signerKey           *ecdsa.PrivateKey
}

var _ capabilities.TriggerCapability = (*triggerConnectorHandler)(nil)
var _ services.Service = &triggerConnectorHandler{}

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*triggerConnectorHandler, error) {
	if connector == nil {
		return nil, errors.New("missing connector")
	}
	handler := &triggerConnectorHandler{
		Validator:           capabilities.NewValidator[TriggerConfig, Input, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: webapiTriggerInfo}),
		connector:           connector,
		signerKey:           signerKey,
		registeredWorkflows: map[string]webapiTrigger{},
		lggr:                lggr.Named("WorkflowConnectorHandler"),
	}

	return handler, nil
}

// Iterate over each topic, checking against senders and rateLimits, then starting event processing and responding
func (h *triggerConnectorHandler) processTrigger(ctx context.Context, gatewayID string, body *api.MessageBody, sender ethCommon.Address, payload workflow.TriggerRequestPayload) {
	// Pass on the payload with the expectation that that is acceptable format for the executor
	wrappedPayload, _ := values.WrapMap(payload)
	hasMatchedAWorkflow := false
	var response Response

	for _, trigger := range h.registeredWorkflows {
		topics := payload.Topics
		// empty topics means all topics
		if len(topics) == 0 {
			topics = []string{}
			for k := range trigger.allowedTopicsMap {
				topics = append(topics, k)
			}
		}

		for _, topic := range topics {

			if trigger.allowedTopicsMap[topic] {
				hasMatchedAWorkflow = true

				if !trigger.allowedSendersMap[sender.String()] {
					h.lggr.Errorw("Unauthorized Sender")
					h.sendResponse(ctx, gatewayID, body, Response{Status: "ERROR", ErrorMessage: "Unauthorized Sender"})
					return
				}
				if !trigger.rateLimiter.Allow(body.Sender) {
					h.lggr.Errorw("request rate-limited")
					h.sendResponse(ctx, gatewayID, body, Response{Status: "ERROR", ErrorMessage: "request rate-limited"})
					return
				}
				h.lggr.Debugw("processtrigger", "topic", topic, "allowedTopicsMap", trigger.allowedTopicsMap)

				TriggerEventID := body.Sender + payload.TriggerEventID
				tr := capabilities.TriggerResponse{
					Event: capabilities.TriggerEvent{
						TriggerType: triggerType,
						ID:          TriggerEventID,
						Outputs:     wrappedPayload,
					},
				}

				trigger.ch <- tr
				response = Response{Status: "ACCEPTED"}
				// Sending n topics that match a workflow with n allowedTopics, can only be triggered once.
				break
			}
		}
	}
	if !hasMatchedAWorkflow {
		h.lggr.Errorw("No Matching Workflow Topics")
		response = Response{Status: "ERROR", ErrorMessage: "No Matching Workflow Topics"}
	}
	_ = h.sendResponse(ctx, gatewayID, body, response)

}

// https://gateway-us-1.chain.link/web-trigger
//   {
//     jsonrpc: "2.0",
//     id: "...",
//     method: "web-trigger",
//     params: {
//       signature: "...",
//       body: {
//         don_id: "workflow_123",
//         payload: {
//           trigger_id: "web-trigger@1.0.0",
//           trigger_event_id: "action_1234567890",
//           timestamp: 1234567890,
//           topics: ["daily_price_update"],
//           params: {
//             bid: "101",
//             ask: "102"
//           }
//         }
//       }
//     }
//   }

// from Web API Trigger Doc
// trigger_id          - ID of the trigger corresponding to the capability ID
// trigger_event_id    - uniquely identifies generated event (scoped to trigger_id and sender)
// timestamp           - timestamp of the event (unix time), needs to be within certain freshness to be processed
// topics            - [OPTIONAL] list of topics (strings) to be started by this event (affects all topics if empty)
// workflow_owners   - [OPTIONAL] list of workflow owners allowed to receive this event (affects all workflows if empty)
// params            - key-value pairs that will be used as trigger output in the workflow Engine (translated to values.Map)

func (h *triggerConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	// TODO: Validate Signature
	body := &msg.Body
	sender := ethCommon.HexToAddress(body.Sender)
	h.lggr.Debugw("handling gateway request", "id", gatewayID, "method", body.Method, "sender", sender, "payload", body.Payload)
	var payload workflow.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		h.sendResponse(ctx, gatewayID, body, Response{Status: "ERROR", ErrorMessage: fmt.Errorf("error %s decoding payload", err.Error()).Error()})
		return
	}

	switch body.Method {
	case workflow.MethodWebAPITrigger:
		h.processTrigger(ctx, gatewayID, body, sender, payload)

	default:
		h.lggr.Errorw("unsupported method", "id", gatewayID, "method", body.Method)
		h.sendResponse(ctx, gatewayID, body, Response{Status: "ERROR", ErrorMessage: fmt.Errorf("unsupported method %s", body.Method).Error()})
	}
}

func (h *triggerConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	h.lggr.Debugw("RegisterTrigger", "req", req)
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

	var rateLimiterCfg common.RateLimiterConfig
	err = reqConfig.RateLimiter.UnwrapTo(&rateLimiterCfg)
	if err != nil {
		h.lggr.Errorw("error creating unwrapping RateLimiter", "err", err, "RateLimiter config", reqConfig.RateLimiter)
		return nil, err
	}

	rateLimiter, err := common.NewRateLimiter(rateLimiterCfg)
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
		allowedTopicsMap:  allowedTopicsMap,
		allowedSendersMap: allowedSendersMap,
		ch:                ch,
		config:            *reqConfig,
		rateLimiter:       rateLimiter,
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
		return err
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

	// TODO remove this and signerKey once Jin's PR is in.
	// if err = msg.Sign(h.signerKey); err != nil {
	// 	return err
	// }
	h.lggr.Debugw("Sending to Gateway", "msg", msg)

	return h.connector.SendToGateway(ctx, gatewayID, msg)
}
