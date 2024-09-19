package webcapabilities

type TargetRequestPayload struct {
	URL        string            `json:"url"`                  // URL to query, only http and https protocols are supported.
	Method     string            `json:"method,omitempty"`     // HTTP verb, defaults to GET.
	Headers    map[string]string `json:"headers,omitempty"`    // HTTP headers, defaults to empty.
	Body       []byte            `json:"body,omitempty"`       // HTTP request body
	TimeoutMs  uint32            `json:"timeoutMs,omitempty"`  // Timeout in milliseconds
	RetryCount uint8             `json:"retryCount,omitempty"` // Number of retries, defaults to 0.
}

type TargetResponsePayload struct {
	Success      bool              `json:"success"`                 // true if HTTP request was successful
	ErrorMessage string            `json:"error_message,omitempty"` // error message in case of failure
	StatusCode   uint8             `json:"statusCode"`              // HTTP status code
	Headers      map[string]string `json:"headers,omitempty"`       // HTTP headers
	Body         []byte            `json:"body,omitempty"`          // HTTP response body
}
