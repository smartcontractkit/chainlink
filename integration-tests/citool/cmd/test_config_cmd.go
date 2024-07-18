package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

var testConfigCmd = &cobra.Command{
	Use:   "test-config",
	Short: "Manage test config",
}

// OverrideConfig holds the configuration data for overrides
type OverrideConfig struct {
	ChainlinkImage                           string
	ChainlinkVersion                         string
	ChainlinkUpgradeImage                    string
	ChainlinkUpgradeVersion                  string
	ChainlinkPostgresVersion                 string
	SelectedNetworks                         []string
	PyroscopeEnabled                         bool
	PyroscopeServerURL                       string
	PyroscopeEnvironment                     string
	PyroscopeKey                             string
	LoggingTestLogCollect                    bool
	LoggingRunID                             string
	LoggingLogTargets                        []string
	LoggingLokiTenantID                      string
	LoggingLokiEndpoint                      string
	LoggingLokiBasicAuth                     string
	LoggingGrafanaBaseURL                    string
	LoggingGrafanaDashboardURL               string
	LoggingGrafanaBearerToken                string
	PrivateEthereumNetworkExecutionLayer     string
	PrivateEthereumNetworkEthereumVersion    string
	PrivateEthereumNetworkCustomDockerImages string
}

const (
	ChainlinkVersionFlag                        = "chainlink-version"
	ChainlinkUpgradeVersionFlag                 = "chainlink-upgrade-version"
	ChainlinkPostgresVersionFlag                = "chainlink-postgres-version"
	SelectedNetworksFlag                        = "selected-networks"
	FromBase64ConfigFlag                        = "from-base64-config"
	LoggingLokiBasicAuthFlag                    = "logging-loki-basic-auth"
	LoggingLokiEndpointFlag                     = "logging-loki-endpoint"
	LoggingRunIDFlag                            = "logging-run-id"
	LoggingLokiTenantIDFlag                     = "logging-loki-tenant-id"
	LoggingGrafanaBaseURLFlag                   = "logging-grafana-base-url"
	LoggingGrafanaDashboardURLFlag              = "logging-grafana-dashboard-url"
	LoggingGrafanaBearerTokenFlag               = "logging-grafana-bearer-token"
	LoggingLogTargetsFlag                       = "logging-log-targets"
	LoggingTestLogCollectFlag                   = "logging-test-log-collect"
	PyroscopeEnabledFlag                        = "pyroscope-enabled"
	PyroscopeServerURLFlag                      = "pyroscope-server-url"
	PyroscopeKeyFlag                            = "pyroscope-key"
	PyroscopeEnvironmentFlag                    = "pyroscope-environment"
	PrivateEthereumNetworkExecutionLayerFlag    = "private-ethereum-network-execution-layer"
	PrivateEthereumNetworkEthereumVersionFlag   = "private-ethereum-network-ethereum-version"
	PrivateEthereumNetworkCustomDockerImageFlag = "private-ethereum-network-custom-docker-image"
)

var oc OverrideConfig

func init() {
	cmds := []*cobra.Command{createTestConfigCmd}
	for _, c := range cmds {
		c.Flags().StringArrayVar(&oc.SelectedNetworks, SelectedNetworksFlag, nil, "Selected networks")
		c.Flags().StringVar(&oc.ChainlinkVersion, ChainlinkVersionFlag, "", "Chainlink version")
		c.Flags().StringVar(&oc.ChainlinkUpgradeVersion, ChainlinkUpgradeVersionFlag, "", "Chainlink upgrade version")
		c.Flags().StringVar(&oc.ChainlinkPostgresVersion, ChainlinkPostgresVersionFlag, "", "Chainlink Postgres version")
		c.Flags().BoolVar(&oc.PyroscopeEnabled, PyroscopeEnabledFlag, false, "Pyroscope enabled")
		c.Flags().StringVar(&oc.PyroscopeServerURL, PyroscopeServerURLFlag, "", "Pyroscope server URL")
		c.Flags().StringVar(&oc.PyroscopeKey, PyroscopeKeyFlag, "", "Pyroscope key")
		c.Flags().StringVar(&oc.PyroscopeEnvironment, PyroscopeEnvironmentFlag, "", "Pyroscope environment")
		c.Flags().BoolVar(&oc.LoggingTestLogCollect, LoggingTestLogCollectFlag, false, "Test log collect")
		c.Flags().StringVar(&oc.LoggingRunID, LoggingRunIDFlag, "", "Run ID")
		c.Flags().StringArrayVar(&oc.LoggingLogTargets, LoggingLogTargetsFlag, nil, "Logging.LogStream.LogTargets")
		c.Flags().StringVar(&oc.LoggingLokiEndpoint, LoggingLokiEndpointFlag, "", "")
		c.Flags().StringVar(&oc.LoggingLokiTenantID, LoggingLokiTenantIDFlag, "", "")
		c.Flags().StringVar(&oc.LoggingLokiBasicAuth, LoggingLokiBasicAuthFlag, "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaBaseURL, LoggingGrafanaBaseURLFlag, "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaDashboardURL, LoggingGrafanaDashboardURLFlag, "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaBearerToken, LoggingGrafanaBearerTokenFlag, "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkExecutionLayer, PrivateEthereumNetworkExecutionLayerFlag, "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkEthereumVersion, PrivateEthereumNetworkEthereumVersionFlag, "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkCustomDockerImages, PrivateEthereumNetworkCustomDockerImageFlag, "", "")

		c.PreRun = func(_ *cobra.Command, _ []string) {
			// Resolve selected networks environment variable if set
			if len(oc.SelectedNetworks) > 0 {
				_, hasEnvVar := utils.LookupEnvVarName(oc.SelectedNetworks[0])
				if hasEnvVar {
					selectedNetworks := utils.MustResolveEnvPlaceholder(oc.SelectedNetworks[0])
					oc.SelectedNetworks = strings.Split(selectedNetworks, ",")
				}
			}

			// Resolve all other environment variables
			oc.ChainlinkImage = utils.MustResolveEnvPlaceholder(oc.ChainlinkImage)
			oc.ChainlinkVersion = utils.MustResolveEnvPlaceholder(oc.ChainlinkVersion)
			oc.ChainlinkUpgradeImage = utils.MustResolveEnvPlaceholder(oc.ChainlinkUpgradeImage)
			oc.ChainlinkUpgradeVersion = utils.MustResolveEnvPlaceholder(oc.ChainlinkUpgradeVersion)
			oc.ChainlinkPostgresVersion = utils.MustResolveEnvPlaceholder(oc.ChainlinkPostgresVersion)
			oc.PyroscopeServerURL = utils.MustResolveEnvPlaceholder(oc.PyroscopeServerURL)
			oc.PyroscopeKey = utils.MustResolveEnvPlaceholder(oc.PyroscopeKey)
			oc.PyroscopeEnvironment = utils.MustResolveEnvPlaceholder(oc.PyroscopeEnvironment)
			oc.LoggingRunID = utils.MustResolveEnvPlaceholder(oc.LoggingRunID)
			oc.LoggingLokiTenantID = utils.MustResolveEnvPlaceholder(oc.LoggingLokiTenantID)
			oc.LoggingLokiEndpoint = utils.MustResolveEnvPlaceholder(oc.LoggingLokiEndpoint)
			oc.LoggingLokiBasicAuth = utils.MustResolveEnvPlaceholder(oc.LoggingLokiBasicAuth)
			oc.LoggingGrafanaBaseURL = utils.MustResolveEnvPlaceholder(oc.LoggingGrafanaBaseURL)
			oc.LoggingGrafanaDashboardURL = utils.MustResolveEnvPlaceholder(oc.LoggingGrafanaDashboardURL)
			oc.LoggingGrafanaBearerToken = utils.MustResolveEnvPlaceholder(oc.LoggingGrafanaBearerToken)
			oc.PrivateEthereumNetworkExecutionLayer = utils.MustResolveEnvPlaceholder(oc.PrivateEthereumNetworkExecutionLayer)
			oc.PrivateEthereumNetworkEthereumVersion = utils.MustResolveEnvPlaceholder(oc.PrivateEthereumNetworkEthereumVersion)
			oc.PrivateEthereumNetworkCustomDockerImages = utils.MustResolveEnvPlaceholder(oc.PrivateEthereumNetworkCustomDockerImages)
		}
	}
}
