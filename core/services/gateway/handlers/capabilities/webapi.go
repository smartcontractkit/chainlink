package capabilities

import "errors"

type Request struct {
	URL       string            `json:"url"`                 // URL to query, only http and https protocols are supported.
	Method    string            `json:"method,omitempty"`    // HTTP verb, defaults to GET.
	Headers   map[string]string `json:"headers,omitempty"`   // HTTP headers, defaults to empty.
	Body      []byte            `json:"body,omitempty"`      // HTTP request body
	TimeoutMs uint32            `json:"timeoutMs,omitempty"` // Timeout in milliseconds

	// Maximum number of bytes to read from the response body.  If the gateway max response size is smaller than this value, the gateway max response size will be used.
	MaxResponseBytes uint32 `json:"maxBytes,omitempty"`
	WorkflowID       string
}

type Response struct {
	ExecutionError bool              `json:"executionError"`         // true if there were non-HTTP errors. false if HTTP request was sent regardless of status (2xx, 4xx, 5xx)
	ErrorMessage   string            `json:"errorMessage,omitempty"` // error message in case of failure
	StatusCode     int               `json:"statusCode,omitempty"`   // HTTP status code
	Headers        map[string]string `json:"headers,omitempty"`      // HTTP headers
	Body           []byte            `json:"body,omitempty"`         // HTTP response body
}

// Validate ensures the Response struct is consistent.
func (r Response) Validate() error {
	if r.ExecutionError {
		if r.ErrorMessage == "" {
			return errors.New("executionError is true but errorMessage is empty")
		}
		if r.StatusCode != 0 || len(r.Headers) > 0 || len(r.Body) > 0 {
			return errors.New("executionError is true but response details (statusCode, headers, body) are populated")
		}
		return nil
	}

	if r.StatusCode < 100 || r.StatusCode > 599 {
		return errors.New("statusCode must be a valid HTTP status code (100-599)")
	}

	return nil
}

type TriggerResponsePayload struct {
	ErrorMessage string `json:"error_message,omitempty"`
	// ERROR, ACCEPTED, PENDING, COMPLETED
	Status string `json:"status"`
}
