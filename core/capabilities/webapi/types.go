package webapi

import (
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/common"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
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
	IncomingRateLimiter gwcommon.RateLimiterConfig `toml:"incomingRateLimiter" json:"incomingRateLimiter" yaml:"incomingRateLimiter" mapstructure:"incomingRateLimiter"`
	OutgoingRateLimiter common.RateLimiterConfig   `toml:"outgoingRateLimiter" json:"outgoingRateLimiter" yaml:"outgoingRateLimiter" mapstructure:"outgoingRateLimiter"`
}
