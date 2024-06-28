package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var testConfigCmd = &cobra.Command{
	Use:   "test-config",
	Short: "Manage test config",
}

// OverrideConfig holds the configuration data for overrides
type OverrideConfig struct {
	ChainlinkImage                           string
	ChainlinkVersion                         string
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
	ChainlinkImageFlag                          = "chainlink-image"
	ChainlinkVersionFlag                        = "chainlink-version"
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
	cmds := []*cobra.Command{createTestConfigCmd, overrideTestConfigCmd}
	for _, c := range cmds {
		c.Flags().StringArrayVar(&oc.SelectedNetworks, SelectedNetworksFlag, nil, "Selected networks")
		c.Flags().StringVar(&oc.ChainlinkImage, ChainlinkImageFlag, "", "Chainlink image")
		c.Flags().StringVar(&oc.ChainlinkVersion, ChainlinkVersionFlag, "", "Chainlink version")
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
				_, hasEnvVar := lookupEnvVarName(oc.SelectedNetworks[0])
				if hasEnvVar {
					selectedNetworks := mustResolveEnvPlaceholder(oc.SelectedNetworks[0])
					oc.SelectedNetworks = strings.Split(selectedNetworks, ",")
				}
			}

			// Resolve all other environment variables
			oc.ChainlinkImage = mustResolveEnvPlaceholder(oc.ChainlinkImage)
			oc.ChainlinkVersion = mustResolveEnvPlaceholder(oc.ChainlinkVersion)
			oc.ChainlinkPostgresVersion = mustResolveEnvPlaceholder(oc.ChainlinkPostgresVersion)
			oc.PyroscopeServerURL = mustResolveEnvPlaceholder(oc.PyroscopeServerURL)
			oc.PyroscopeKey = mustResolveEnvPlaceholder(oc.PyroscopeKey)
			oc.PyroscopeEnvironment = mustResolveEnvPlaceholder(oc.PyroscopeEnvironment)
			oc.LoggingRunID = mustResolveEnvPlaceholder(oc.LoggingRunID)
			oc.LoggingLokiTenantID = mustResolveEnvPlaceholder(oc.LoggingLokiTenantID)
			oc.LoggingLokiEndpoint = mustResolveEnvPlaceholder(oc.LoggingLokiEndpoint)
			oc.LoggingLokiBasicAuth = mustResolveEnvPlaceholder(oc.LoggingLokiBasicAuth)
			oc.LoggingGrafanaBaseURL = mustResolveEnvPlaceholder(oc.LoggingGrafanaBaseURL)
			oc.LoggingGrafanaDashboardURL = mustResolveEnvPlaceholder(oc.LoggingGrafanaDashboardURL)
			oc.LoggingGrafanaBearerToken = mustResolveEnvPlaceholder(oc.LoggingGrafanaBearerToken)
			oc.PrivateEthereumNetworkExecutionLayer = mustResolveEnvPlaceholder(oc.PrivateEthereumNetworkExecutionLayer)
			oc.PrivateEthereumNetworkEthereumVersion = mustResolveEnvPlaceholder(oc.PrivateEthereumNetworkEthereumVersion)
			oc.PrivateEthereumNetworkCustomDockerImages = mustResolveEnvPlaceholder(oc.PrivateEthereumNetworkCustomDockerImages)
		}
	}
}

// mustResolveEnvPlaceholder checks if the input string is an environment variable placeholder and resolves it.
func mustResolveEnvPlaceholder(input string) string {
	envVarName, hasEnvVar := lookupEnvVarName(input)
	if hasEnvVar {
		value, set := os.LookupEnv(envVarName)
		if !set {
			fmt.Fprintf(os.Stderr, "Error resolving '%s'. Environment variable '%s' not set or is empty\n", input, envVarName)
			os.Exit(1)
		}
		return value
	}
	return input
}

func lookupEnvVarName(input string) (string, bool) {
	re := regexp.MustCompile(`^\${{ env\.([a-zA-Z_]+) }}$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}
