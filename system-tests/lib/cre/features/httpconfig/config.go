package httpconfig

import (
	"fmt"
	"maps"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

// Build puts HTTP service settings in the registry spec_config consumed by the
// capability manager. It accepts the topology's existing PascalCase settings
// and the capability's native JSON field names.
func Build(input map[string]any, localOnly bool) (*capabilitiespb.CapabilityConfig, error) {
	config := make(map[string]any)
	for source, target := range map[string]string{
		"ProxyMode":                    "proxyMode",
		"GatewayConnection":            "gatewayConnection",
		"HTTPClient":                   "httpClient",
		"MetadataBatchSize":            "metadataBatchSize",
		"SendChannelBufferSize":        "sendChannelBufferSize",
		"MaxAuthorizedKeysPerWorkflow": "maxAuthorizedKeysPerWorkflow",
		"RequestCacheTTL":              "requestCacheTTL",
	} {
		if v, ok := input[target]; ok {
			config[target] = v
		} else if v, ok := input[source]; ok {
			config[target] = v
		}
	}
	// Metadata timeouts belong inside gatewayConnection in the capability API.
	gateway, _ := config["gatewayConnection"].(map[string]any)
	gateway = maps.Clone(gateway)
	for source, target := range map[string]string{
		"MaxPushMetadataDurationMs": "maxPushMetadataDurationMs",
		"MaxPullMetadataDurationMs": "maxPullMetadataDurationMs",
	} {
		if v, ok := input[source]; ok {
			if gateway == nil {
				gateway = make(map[string]any)
			}
			if _, exists := gateway[target]; !exists {
				gateway[target] = v
			}
		}
	}
	if len(gateway) > 0 {
		config["gatewayConnection"] = gateway
	}
	for _, direction := range []struct{ source, target string }{
		{"Incoming", "incomingRateLimiter"},
		{"Outgoing", "outgoingRateLimiter"},
	} {
		if v, ok := input[direction.target]; ok {
			config[direction.target] = v
			continue
		}
		limiter := make(map[string]any)
		for source, target := range map[string]string{
			"GlobalBurst": "globalBurst", "GlobalRPS": "globalRPS",
			"PerSenderBurst": "perSenderBurst", "PerSenderRPS": "perSenderRPS",
		} {
			if v, ok := input[direction.source+source]; ok {
				limiter[target] = v
			}
		}
		if len(limiter) > 0 {
			config[direction.target] = limiter
		}
	}
	specConfig, err := values.NewMap(config)
	if err != nil {
		return nil, fmt.Errorf("build HTTP spec_config: %w", err)
	}
	return &capabilitiespb.CapabilityConfig{
		LocalOnly:  localOnly,
		SpecConfig: values.ProtoMap(specConfig),
	}, nil
}
