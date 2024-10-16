package capabilities

type Request struct {
	URL       string            `json:"url"`                 // URL to query, only http and https protocols are supported.
	Method    string            `json:"method,omitempty"`    // HTTP verb, defaults to GET.
	Headers   map[string]string `json:"headers,omitempty"`   // HTTP headers, defaults to empty.
	Body      []byte            `json:"body,omitempty"`      // HTTP request body
	TimeoutMs uint32            `json:"timeoutMs,omitempty"` // Timeout in milliseconds
}

type Response struct {
	ExecutionError bool              `json:"executionError"`         // true if there were non-HTTP errors. false if HTTP request was sent regardless of status (2xx, 4xx, 5xx)
	ErrorMessage   string            `json:"errorMessage,omitempty"` // error message in case of failure
	StatusCode     int               `json:"statusCode,omitempty"`   // HTTP status code
	Headers        map[string]string `json:"headers,omitempty"`      // HTTP headers
	Body           []byte            `json:"body,omitempty"`         // HTTP response body
}

type TriggerResponsePayload struct {
	ErrorMessage string `json:"error_message,omitempty"`
	// ERROR, ACCEPTED, PENDING, COMPLETED
	Status string `json:"status"`
}
