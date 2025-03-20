package webapi

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	SingleNode string = "SingleNode"
	// TODO: AllAtOnce is not yet implemented
	AllAtOnce string = "AllAtOnce"
)

// ServiceConfig is the configuration for the Target capability and handler
// TODO: handle retry configurations here CM-472
// Note that workflow executions have their own internal timeouts and retries set by the user
// that are separate from this configuration
type ServiceConfig struct {
	// RateLimiter configuration for messages incoming to this node from the gateway.
	// The sender is a Gateway node, which is identified by the Gateway ID.
	RateLimiter common.RateLimiterConfig `toml:"incomingRateLimiter" json:"incomingRateLimiter" yaml:"incomingRateLimiter" mapstructure:"incomingRateLimiter"`
	// RateLimiter configuration for outgoing messages from this node to the gateway.
	// The sender is a workflow, which is identified by the Workflow ID.
	OutgoingRateLimiter common.RateLimiterConfig `toml:"outgoingRateLimiter" json:"outgoingRateLimiter" yaml:"outgoingRateLimiter" mapstructure:"outgoingRateLimiter"`
}
