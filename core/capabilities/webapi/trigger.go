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

type Input struct {
}
type TriggerConfig struct {
	AllowedSenders []ethCommon.Address      `toml:"allowedSenders"`
	AllowedTopics  []string                 `toml:"allowedTopics"`
	RateLimiter    common.RateLimiterConfig `toml:"rateLimiter"`
	RequiredParams []string                 `toml:"requiredParams"`
}

// TODO: Question asked in Web API Trigger about this structure
// https://docs.google.com/document/d/1mCTAo-ix-P923eUlh4SloZfBN9PCvgf90oHWbmykjsc/edit?disco=AAABWF65hAM
type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
	Status       string `json:"status"`
}

type webapiTrigger struct {
	allowedSendersMap map[string]bool
	allowedTopicsMap  map[string]bool
	ch                chan<- capabilities.TriggerResponse
	config            TriggerConfig
	rateLimiter       *common.RateLimiter
}

// Handles connections to the webapi trigger
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

// TODO: From Design doc,
// Once connected to a Gateway, each connector handler periodically sends metadata messages containing aggregated
// config for all registered workflow specs using web-trigger.

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*triggerConnectorHandler, error) {
	// TODO (CAPPL-22, CAPPL-24):
	//   - decode config
	//   - create an implementation of the capability API and add it to the Registry
	//   - create a handler and register it with Gateway Connector
	//   - manage trigger subscriptions
	//   - process incoming trigger events and related metadata

	handler := &triggerConnectorHandler{
		connector:           connector,
		signerKey:           signerKey,
		registeredWorkflows: map[string]webapiTrigger{},
		lggr:                lggr.Named("WorkflowConnectorHandler"),
	}

	return handler, nil
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
	body := &msg.Body
	sender := ethCommon.HexToAddress(body.Sender)
	h.lggr.Debugw("handling gateway request", "id", gatewayID, "method", body.Method, "sender", sender, "payload", body.Payload)
	var payload workflow.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		return
	}

	switch body.Method {
	case workflow.MethodWebAPITrigger:
		h.lggr.Debugw("added MethodWebAPITrigger message", "payload", string(body.Payload))
		// Sample payload
		// payload: {
		// 	trigger_id: "web-trigger@1.0.0",
		// 	trigger_event_id: "action_1234567890",
		// 	timestamp: 1234567890,
		// 	topics: ["daily_price_update"],
		// 	params: {
		// 		bid: "101",
		// 		ask: "102"
		// 	}
		// }

		// I don't see anywhere there is a mapping to a different structure such as
		// payload: {
		//   topic: "daily_price_update"
		// 	 params: {
		// 		 bid: "101",
		// 		 ask: "102"
		// 	 }

		// So this just passes it on with the expectation that that is acceptable format for the executor

		wrappedPayload, _ := values.WrapMap(payload)

		for _, trigger := range h.registeredWorkflows {

			// TODO: Question asked in Web API trigger about checking for completed Triggers to return COMPLETED
			// "TriggerEventID used internally by the Engine is a pair (sender, trigger_event_id).
			// This is to protect against a case where two different authorized senders use the same event ID in their messages.

			// TODO: how do we know PENDING state, that is node received the event but processing hasn't finished.

			topics := payload.Topics

			for _, topic := range topics {
				if !trigger.allowedTopicsMap[topic] {
					return
				}
			}
			sender := ethCommon.HexToAddress(body.Sender).String()

			if !trigger.allowedSendersMap[sender] {
				h.lggr.Errorw("Unauthorized Sender")
				return
			}
			if !trigger.rateLimiter.Allow(body.Sender) {
				h.lggr.Errorw("request rate-limited")
				return
			}
			// TODO: CAPPL-24 check the topic to see if the method is a duplicate and the trigger has been sent, ie PENDING

			TriggerEventID := body.Sender + payload.TriggerEventId
			tr := capabilities.TriggerResponse{
				Event: capabilities.TriggerEvent{
					TriggerType: triggerType,
					ID:          TriggerEventID,
					Outputs:     wrappedPayload,
				},
			}

			trigger.ch <- tr

			// TODO: PENDING
			response := Response{Success: true, Status: "ACCEPTED"}
			h.sendResponse(ctx, gatewayID, body, response)
		}
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayID, "method", body.Method)
	}
}

func (h *triggerConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	h.lggr.Debugw("RegisterTrigger", "req", req)
	cfg := req.Config
	if cfg == nil {
		return nil, errors.New("config is required to register a web api trigger")
	}
	// This is failing in the call to ConfigSchema due to nil pointer
	// github.com/smartcontractkit/chainlink-common/pkg/capabilities.(*Validator[...]).ConfigSchema(0x104b4b36c?)
	// on this line
	// 	return schemaWith(*v.ConfigReflector, v.configType, v.schemas, "config", v.Info)

	// Sri does this style in PR
	// https://github.com/smartcontractkit/chainlink/pull/14308/files#diff-05dcc9ad24b80011a87a133958bbe7c5ed83e814a928e9d3257a95fbde981b2cR100
	//
	// I've tried simplifying to topics, a simple string array without success.
	// Not sure how the generics are done and how I could pass in a logger to see what value is nil.
	// Also not sure how the validator gets the configType, schemas, etc.
	// I've initialized similarly to PR 14308 with
	// 	capabilities.Validator[TriggerConfig, Input, capabilities.TriggerResponse]

	reqConfig, err := h.ValidateConfig(cfg)
	if err != nil {
		h.lggr.Errorw("error validating config", "err", err)

		return nil, err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.registeredWorkflows[req.TriggerID]
	if ok {
		return nil, fmt.Errorf("triggerId %s already registered", req.TriggerID)
	}

	ch := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	rateLimiter, err := common.NewRateLimiter(reqConfig.RateLimiter)
	if err != nil {
		return nil, err
	}
	allowedSendersMap := map[string]bool{}
	for _, k := range reqConfig.AllowedSenders {
		allowedSendersMap[k.String()] = true
	}

	allowedTopicsMap := map[string]bool{}
	for _, k := range reqConfig.AllowedTopics {
		allowedTopicsMap[k] = true
	}

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
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		h.lggr.Errorw("error marshallig payload", "err", err)
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: requestBody.MessageId,
			DonId:     requestBody.DonId,
			Method:    requestBody.Method,
			Receiver:  requestBody.Sender,
			Payload:   payloadJson,
		},
	}

	// TODO remove this and signerKey once Jin's PR is in.
	// if err = msg.Sign(h.signerKey); err != nil {
	// 	return err
	// }
	h.lggr.Debugw("Sending to Gateway", "msg", msg)

	return h.connector.SendToGateway(ctx, gatewayID, msg)
}
