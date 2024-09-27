package webapicapabilities

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type TargetRequestPayload struct {
	URL       string            `json:"url"`                 // URL to query, only http and https protocols are supported.
	Method    string            `json:"method,omitempty"`    // HTTP verb, defaults to GET.
	Headers   map[string]string `json:"headers,omitempty"`   // HTTP headers, defaults to empty.
	Body      []byte            `json:"body,omitempty"`      // HTTP request body
	TimeoutMs uint32            `json:"timeoutMs,omitempty"` // Timeout in milliseconds
}

type TargetResponsePayload struct {
	Success      bool              `json:"success"`                 // true if HTTP request was successful
	ErrorMessage string            `json:"error_message,omitempty"` // error message in case of failure
	StatusCode   uint8             `json:"statusCode"`              // HTTP status code
	Headers      map[string]string `json:"headers,omitempty"`       // HTTP headers
	Body         []byte            `json:"body,omitempty"`          // HTTP response body
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

// amendments for V1
// topics must be specified and a single one.
// workflow_owners is omitted
// params is advisory only, the Executor will parse the parameters.
type TriggerRequestPayload struct {
	TriggerID      string     `json:"trigger_id"`
	TriggerEventID string     `json:"trigger_event_id"`
	Timestamp      int64      `json:"timestamp"`
	Topics         []string   `json:"topics"`
	Params         values.Map `json:"params"`
}

type TriggerResponsePayload struct {
	ErrorMessage string `json:"error_message,omitempty"`
	// ERROR, ACCEPTED, PENDING, COMPLETED
	Status string `json:"status"`
}
