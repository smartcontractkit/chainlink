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
	ExecutionError bool              `json:"executionError"`         // true if there were non-HTTP errors. false if HTTP request was sent regardless of status (2xx, 4xx, 5xx)
	ErrorMessage   string            `json:"errorMessage,omitempty"` // error message in case of failure
	StatusCode     int               `json:"statusCode,omitempty"`   // HTTP status code
	Headers        map[string]string `json:"headers,omitempty"`      // HTTP headers
	Body           []byte            `json:"body,omitempty"`         // HTTP response body
}

// https://gateway-us-1.chain.link/web-api-trigger
//
//	{
//	  jsonrpc: "2.0",
//	  id: "...",
//	  method: "web-api-trigger",
//	  params: {
//	    signature: "...",
//	    body: {
//	      don_id: "workflow_123",
//	      payload: {
//	        trigger_id: "web-api-trigger@1.0.0",
//	        trigger_event_id: "action_1234567890",
//	        timestamp: 1234567890,
//	        topics: ["daily_price_update"],
//	        params: {
//	          bid: "101",
//	          ask: "102"
//	        }
//	      }
//	    }
//	  }
//	}
//
// from Web API Trigger Doc, with modifications.
// trigger_id          - ID of the trigger corresponding to the capability ID
// trigger_event_id    - uniquely identifies generated event (scoped to trigger_id and sender)
// timestamp           - timestamp of the event (unix time), needs to be within certain freshness to be processed
// topics            - an array of a single topic (string) to be started by this event
// params            - key-value pairs for the workflow engine, untranslated.
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
