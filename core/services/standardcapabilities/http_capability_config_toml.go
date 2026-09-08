package standardcapabilities

import (
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// httpTriggerConfigKeys maps the http-trigger capability's JSON config keys.
// These mirror the ServiceConfig struct in the http_trigger capability binary.
const (
	keyMetadataBatchSize            = "metadataBatchSize"
	keySendChannelBufferSize        = "sendChannelBufferSize"
	keyMaxAuthorizedKeysPerWorkflow = "maxAuthorizedKeysPerWorkflow"
	keyRequestCacheTTL              = "requestCacheTTL"
	keyGatewayConnection            = "gatewayConnection"
	keyRetryConfig                  = "retryConfig"
	keyInitialIntervalMs            = "initialIntervalMs"
	keyMaxIntervalTimeMs            = "maxIntervalTimeMs"
	keyMaxElapsedTimeMs             = "maxElapsedTimeMs"
	keyMultiplier                   = "multiplier"
	keyMaxPushMetadataDurationMs    = "maxPushMetadataDurationMs"
	keyMaxPullMetadataDurationMs    = "maxPullMetadataDurationMs"
)

// httpActionConfigKeys maps the http-action capability's JSON config keys.
// These mirror the ServiceConfig struct in the http_action capability binary.
const (
	keyProxyMode      = "proxyMode"
	keyHttpClient     = "httpClient"
	keyBlockedIPs     = "blockedIPs"
	keyBlockedIPsCIDR = "blockedIPsCIDR"
	keyAllowedPorts   = "allowedPorts"
	keyAllowedSchemes = "allowedSchemes"
	keyAllowedIPs     = "allowedIPs"
	keyAllowedIPsCIDR = "allowedIPsCIDR"
)

// injectHTTPTriggerConfig merges [Capabilities.HTTPTrigger] TOML values into
// the http-trigger capability's config JSON. TOML values win; unset values
// fall back to job-spec values or the capability binary's defaults.
func injectHTTPTriggerConfig(lggr logger.Logger, cfg coreconfig.HTTPTriggerCapability, configJSON string) string {
	return injectHTTPCapabilityConfig(lggr, "http_trigger", configJSON, func(m map[string]any) {
		setIfNotZero(m, keyMetadataBatchSize, cfg.MetadataBatchSize())
		setIfNotZero(m, keySendChannelBufferSize, cfg.SendChannelBufferSize())
		setIfNotZero(m, keyMaxAuthorizedKeysPerWorkflow, cfg.MaxAuthorizedKeysPerWorkflow())
		setIfNotZero(m, keyRequestCacheTTL, cfg.RequestCacheTTL())

		retry := cfg.GatewayConnection().RetryConfig()
		retryMap := make(map[string]any)
		setIfNotZero(retryMap, keyInitialIntervalMs, retry.InitialIntervalMs())
		setIfNotZero(retryMap, keyMaxIntervalTimeMs, retry.MaxIntervalTimeMs())
		setIfNotZero(retryMap, keyMultiplier, retry.Multiplier())

		gw := cfg.GatewayConnection()
		gwMap := make(map[string]any)
		attachNested(gwMap, keyRetryConfig, retryMap)
		setIfNotZero(gwMap, keyMaxPushMetadataDurationMs, gw.MaxPushMetadataDurationMs())
		setIfNotZero(gwMap, keyMaxPullMetadataDurationMs, gw.MaxPullMetadataDurationMs())

		// Merge into any existing job-spec gatewayConnection rather than
		// replacing it, so job-spec values survive when TOML is unset.
		existingGw := nested(m, keyGatewayConnection)
		for k, v := range gwMap {
			existingGw[k] = v
		}
		if len(existingGw) == 0 {
			delete(m, keyGatewayConnection)
		}
	})
}

// injectHTTPActionConfig merges [Capabilities.HTTPAction] TOML values into
// the http-action capability's config JSON. TOML values win; unset values
// fall back to job-spec values or the capability binary's defaults.
//
// The httpClient network restrictions are only sourced from node TOML — they
// are never emitted into job specs, since they can be sensitive.
func injectHTTPActionConfig(lggr logger.Logger, cfg coreconfig.HTTPActionCapability, configJSON string) string {
	return injectHTTPCapabilityConfig(lggr, "http_action", configJSON, func(m map[string]any) {
		setIfNotZero(m, keyProxyMode, cfg.ProxyMode())

		gw := cfg.GatewayConnection()
		gwMap := make(map[string]any)
		setIfNotZero(gwMap, keyInitialIntervalMs, gw.InitialIntervalMs())
		setIfNotZero(gwMap, keyMaxElapsedTimeMs, gw.MaxElapsedTimeMs())
		setIfNotZero(gwMap, keyMultiplier, gw.Multiplier())

		// Merge into any existing job-spec gatewayConnection rather than
		// replacing it, so job-spec values survive when TOML is unset.
		existingGw := nested(m, keyGatewayConnection)
		for k, v := range gwMap {
			existingGw[k] = v
		}
		if len(existingGw) == 0 {
			delete(m, keyGatewayConnection)
		}

		hc := cfg.HTTPClient()
		hcMap := make(map[string]any)
		setIfNotZero(hcMap, keyBlockedIPs, hc.BlockedIPs())
		setIfNotZero(hcMap, keyBlockedIPsCIDR, hc.BlockedIPsCIDR())
		setIfNotZero(hcMap, keyAllowedPorts, hc.AllowedPorts())
		setIfNotZero(hcMap, keyAllowedSchemes, hc.AllowedSchemes())
		setIfNotZero(hcMap, keyAllowedIPs, hc.AllowedIPs())
		setIfNotZero(hcMap, keyAllowedIPsCIDR, hc.AllowedIPsCIDR())

		// Merge into any existing job-spec httpClient rather than replacing
		// it, so direct-mode settings survive when TOML is unset.
		existingHc := nested(m, keyHttpClient)
		for k, v := range hcMap {
			existingHc[k] = v
		}
		if len(existingHc) == 0 {
			delete(m, keyHttpClient)
		}
	})
}
