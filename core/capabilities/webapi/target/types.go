package target

import "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"

const (
	SingleNode string = "SingleNode"
	// TODO: AllAtOnce is not yet implemented
	AllAtOnce string = "AllAtOnce"
)

type Input struct {
	URL     string            `json:"url"`               // URL to query, only http and https protocols are supported.
	Method  string            `json:"method,omitempty"`  // HTTP verb, defaults to GET.
	Headers map[string]string `json:"headers,omitempty"` // HTTP headers, defaults to empty.
	Body    []byte            `json:"body,omitempty"`    // HTTP body, defaults to empty.
}

// WorkflowConfig is the configuration of the workflow that is passed in the workflow execute request
type WorkflowConfig struct {
	TimeoutMs    uint32 `toml:"timeoutMs" json:"timeoutMs" yaml:"timeoutMs" mapstructure:"timeoutMs"`             // Timeout in milliseconds
	RetryCount   uint8  `toml:"retryCount" json:"retryCount" yaml:"retryCount" mapstructure:"retryCount"`         // Number of retries, defaults to 0.
	DeliveryMode string `toml:"deliveryMode" json:"deliveryMode" yaml:"deliveryMode" mapstructure:"deliveryMode"` // DeliveryMode describes how request should be delivered to gateway nodes, defaults to SingleNode.
}

// Config is the configuration for the Target capability and handler
// TODO: handle retry configurations here CM-472
// Note that workflow executions have their own internal timeouts and retries set by the user
// that are separate from this configuration
type Config struct {
	RateLimiter common.RateLimiterConfig `toml:"rateLimiter" json:"rateLimiter" yaml:"rateLimiter" mapstructure:"rateLimiter"`
	WorkflowConfig
}
