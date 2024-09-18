package webapi

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const defaultSendChannelBufferSize = 1000

const triggerType = "web-trigger@1.0.0"

type Response struct {
	// TODO: what is the format for ACCEPTED, PENDING, COMPLETED status?
	// Status    string  `json:"status"`?
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Handles connections to the webapi trigger
type triggerConnectorHandler struct {
	services.StateMachine

	capabilities.CapabilityInfo
	connector connector.GatewayConnector
	lggr      logger.Logger
	mu        sync.Mutex
	// Will this have to get pulled into a store to have the topic and workflow ID?
	registeredWorkflows map[string]chan capabilities.TriggerResponse
	signerKey           *ecdsa.PrivateKey
}

var _ capabilities.TriggerCapability = (*triggerConnectorHandler)(nil)
var _ services.Service = &triggerConnectorHandler{}

func NewTrigger(config string, registry core.CapabilitiesRegistry, connector connector.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (job.ServiceCtx, error) {
	// TODO (CAPPL-22, CAPPL-24):
	//   - decode config
	//   - create an implementation of the capability API and add it to the Registry
	//   - create a handler and register it with Gateway Connector
	//   - manage trigger subscriptions
	//   - process incoming trigger events and related metadata

	handler := &triggerConnectorHandler{
		connector: connector,
		signerKey: signerKey,
		lggr:      lggr.Named("WorkflowConnectorHandler"),
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
//           sub-events: [
//             {
//               topics: ["daily_price_update"],
//               params: {
//                 bid: "101",
//                 ask: "102"
//               }
//             },
//             {
//               topics: ["daily_message", "summary"],
//               params: {
//                 message: "all good!",
//               }
//             },
//           ]
//         }
//       }
//     }
//   }

// from Web API Trigger Doc
// trigger_id          - ID of the trigger corresponding to the capability ID
// trigger_event_id    - uniquely identifies generated event (scoped to trigger_id and sender)xx
// timestamp           - timestamp of the event, needs to be within certain freshness to be processed
// sub_events {        - a list of per-topic-set components of this trigger event (can be of size 1)
//   topics            - [OPTIONAL] list of topics (strings) to be started by this event (affects all topics if empty)
//   workflow_owners   - [OPTIONAL] list of workflow owners allowed to receive this event (affects all workflows if empty)
//   params            - key-value pairs that will be used as trigger output in the workflow Engine (translated to values.Map)

type TriggerRequestPayload struct {
	TriggerId      string `json:"trigger_id"`
	TriggerEventId string `json:"trigger_event_id"`
	// how are timestamps defined?  ISO-8601 or UTC seconds or UTC ms?
	Timestamp uint        `json:"timestamp"`
	SubEvents []SubEvents `json:"sub_events"`
}

type SubEvents struct {
	Topics []string   `json:"topics"`
	Params values.Map `json:"params"`
}

func (h *triggerConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	body := &msg.Body
	fromAddr := ethCommon.HexToAddress(body.Sender)
	// TODO: apply allowlist and rate-limiting
	h.lggr.Debugw("handling gateway request", "id", gatewayID, "method", body.Method, "sender", fromAddr)
	var payload TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		return
	}

	// TODO: how to convert payload to *values.Map.  Parse directly to that instead of the structs?
	// Sri did wrappedPayload, err := values.WrapMap(log.Data), does that work in this case?

	// TODO: How/where to check timestamp for freshness

	switch body.Method {
	case workflow.MethodWebAPITrigger:
		h.lggr.Debugw("added MethodWebAPITrigger message", "payload", string(body.Payload))

		for triggerID, trigger := range h.registeredWorkflows {

			// TODO: CAPPL-24 extract the topic and then match the subscriber to this messages triggers.
			// TODO: CAPPL-24 check the topic to see if the method is a duplicate and the trigger has been sent, ie PENDING
			// TODO: Question asked in Web API trigger about checking for completed Triggers to return COMPLETED
			// TODO: update ID to conform to
			// "TriggerEventID used internally by the Engine is a pair (sender, trigger_event_id).
			// This is to protect against a case where two different authorized senders use the same event ID in their messages.
			tr := capabilities.TriggerResponse{
				Event: capabilities.TriggerEvent{
					TriggerType: triggerType,
					ID:          triggerID,
					Outputs:     payload, // must be *values.Map
				},
			}
			trigger <- tr
		}

		// TODO: ACCEPTED, PENDING, COMPLETED
		response := Response{Success: true}
		h.sendResponse(ctx, gatewayID, body, response)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayID, "method", body.Method)
	}
}

func (h *triggerConnectorHandler) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.registeredWorkflows[req.TriggerID]
	if ok {
		return nil, fmt.Errorf("triggerId %s already registered", req.TriggerID)
	}

	callbackCh := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	// TODO: CAPPL-24 how do we extract the topic and then define the trigger by that?
	// It's not TriggerID because TriggerID is concat of workflow ID and the trigger's index in the spec (what does that mean?)
	// I'm not sure if the workflow config comes in via the req.Config

	h.registeredWorkflows[req.TriggerID] = callbackCh

	return callbackCh, nil
}

func (h *triggerConnectorHandler) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	trigger, ok := h.registeredWorkflows[req.TriggerID]
	if ok {
		return fmt.Errorf("triggerId %s not registered", req.TriggerID)
	}

	// Close callback channel
	close(trigger)
	// Remove from triggers context
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
	if err = msg.Sign(h.signerKey); err != nil {
		return err
	}
	return h.connector.SendToGateway(ctx, gatewayID, msg)
}
