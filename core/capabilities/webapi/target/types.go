package target

const (
	AllAtOnce  string = "AllAtOnce"
	RoundRobin string = "RoundRobin"
)

type WorkflowInput struct {
	URL     string            `json:"url"`               // URL to query, only http and https protocols are supported.
	Method  string            `json:"method,omitempty"`  // HTTP verb, defaults to GET.
	Headers map[string]string `json:"headers,omitempty"` // HTTP headers, defaults to empty.
	Body    string            `json:"body,omitempty"`    // Base64-encoded binary body, defaults to empty.
}

type WorkflowConfig struct {
	TimeoutMs  uint32 `json:"timeoutMs,omitempty"`  // Timeout in milliseconds
	RetryCount uint8  `json:"retryCount,omitempty"` // Number of retries, defaults to 0.
	Schedule   string `json:"schedule,omitempty"`   // schedule, defaults to empty.
}
