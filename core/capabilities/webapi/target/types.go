package target

import "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"

const (
	AllAtOnce  string = "AllAtOnce"
	SingleNode string = "SingleNode"
)

type Input struct {
	URL     string            `json:"url"`               // URL to query, only http and https protocols are supported.
	Method  string            `json:"method,omitempty"`  // HTTP verb, defaults to GET.
	Headers map[string]string `json:"headers,omitempty"` // HTTP headers, defaults to empty.
	Body    string            `json:"body,omitempty"`    // Base64-encoded binary body, defaults to empty.
}

// WorkflowConfig is the configuration of the workflow that is passed in the workflow execute request
type WorkflowConfig struct {
	TimeoutMs  uint32 `json:"timeoutMs,omitempty"`  // Timeout in milliseconds
	RetryCount uint8  `json:"retryCount,omitempty"` // Number of retries, defaults to 0.
	Schedule   string `json:"schedule,omitempty"`   // schedule, defaults to empty.
}

// CapabilityConfigConfig is the configuration for the Target capability and handler
// TODO: handle retry configurations here
// Note that workflow executions have their own internal timeouts and retries set by the user
// that are separate from this configuration
type Config struct {
	RateLimiter common.RateLimiterConfig `toml:"rateLimiter"`
}
