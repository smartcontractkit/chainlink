package testconfig

import (
	"github.com/smartcontractkit/wasp"
)

// TODO until we add TOML config support to WASP
func LokiConfigFromToml(globalConfig GlobalTestConfig) *wasp.LokiConfig {
	lokiConfig := wasp.NewEnvLokiConfig()
	lokiConfig.BasicAuth = *globalConfig.GetLoggingConfig().Loki.BasicAuth
	lokiConfig.TenantID = *globalConfig.GetLoggingConfig().Loki.TenantId
	lokiConfig.URL = *globalConfig.GetLoggingConfig().Loki.Endpoint

	return lokiConfig
}
