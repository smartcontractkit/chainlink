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

const TriggerType = "web-trigger@1.0.0"

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
