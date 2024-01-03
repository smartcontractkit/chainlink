package testconfig

import (
	"github.com/smartcontractkit/wasp"
)

// TODO until we add TOML config support to WASP
func LokiConfigFromToml(config *TestConfig) *wasp.LokiConfig {
	lokiConfig := wasp.NewEnvLokiConfig()
	lokiConfig.BasicAuth = *config.Logging.Loki.BasicAuth
	lokiConfig.TenantID = *config.Logging.Loki.TenantId
	lokiConfig.URL = *config.Logging.Loki.Url

	return lokiConfig
}
