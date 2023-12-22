package testconfig

import (
	"github.com/smartcontractkit/wasp"
)

// TODO until we add TOML config support to WASP
func LokiConfigFromToml(config *TestConfig) *wasp.LokiConfig {
	lokiConfig := wasp.NewEnvLokiConfig()
	lokiConfig.BasicAuth = *config.Logging.Loki.LokiBasicAuth
	lokiConfig.TenantID = *config.Logging.Loki.LokiTenantId
	lokiConfig.URL = *config.Logging.Loki.LokiUrl

	return lokiConfig
}
