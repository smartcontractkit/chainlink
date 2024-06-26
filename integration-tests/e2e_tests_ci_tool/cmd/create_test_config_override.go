package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/spf13/cobra"
)

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
	LoggingGrafanaToken                      string
	PrivateEthereumNetworkExecutionLayer     string
	PrivateEthereumNetworkEthereumVersion    string
	PrivateEthereumNetworkCustomDockerImages string
}

var oc OverrideConfig

var createTestConfigOverrideCmd = &cobra.Command{
	Use:   "create-test-config-override",
	Short: "Processes configurations and outputs templated results for test config overrides",
	Run: func(cmd *cobra.Command, args []string) {
		var config ctf_config.TestConfig

		var image, version, postgresVersion *string
		if cmd.Flags().Changed("chainlink-image") {
			image = &oc.ChainlinkImage
		}
		if cmd.Flags().Changed("chainlink-version") {
			version = &oc.ChainlinkVersion
		}
		if cmd.Flags().Changed("chainlink-postgres-version") {
			version = &oc.ChainlinkPostgresVersion
		}
		if image != nil && version == nil || image == nil && version != nil {
			fmt.Fprintf(os.Stderr, "Error: both chainlink-image and chainlink-version must be set\n")
			os.Exit(1)
		}
		if image != nil && version != nil {
			config.ChainlinkImage = &ctf_config.ChainlinkImageConfig{
				Image:           image,
				Version:         version,
				PostgresVersion: postgresVersion,
			}
		}

		var selectedNetworks *[]string
		if cmd.Flags().Changed("selected-networks") {
			selectedNetworks = &oc.SelectedNetworks
		}
		if selectedNetworks != nil {
			config.Network = &ctf_config.NetworkConfig{
				SelectedNetworks: oc.SelectedNetworks,
			}
		}

		var peryscopeEnabled *bool
		var pyroscopeServerURL, pyroscopeEnvironment, pyroscopeKey *string
		if cmd.Flags().Changed("pyroscope-enabled") {
			peryscopeEnabled = &oc.PyroscopeEnabled
		}
		if cmd.Flags().Changed("peryscope-server-url") {
			pyroscopeServerURL = &oc.PyroscopeServerURL
		}
		if cmd.Flags().Changed("peryscope-server-key") {
			pyroscopeKey = &oc.PyroscopeKey
		}
		if cmd.Flags().Changed("peryscope-environment") {
			pyroscopeEnvironment = &oc.PyroscopeEnvironment
		}
		if peryscopeEnabled != nil {
			config.Pyroscope = &ctf_config.PyroscopeConfig{
				Enabled:     peryscopeEnabled,
				ServerUrl:   pyroscopeServerURL,
				Environment: pyroscopeEnvironment,
				Key:         pyroscopeKey,
			}
		}

		var testLogCollect *bool
		if cmd.Flags().Changed("logging-test-log-collect") {
			testLogCollect = &oc.LoggingTestLogCollect
		}
		var loggingRunID *string
		if cmd.Flags().Changed("logging-run-id") {
			loggingRunID = &oc.LoggingRunID
		}
		var loggingLogTargets []string
		if cmd.Flags().Changed("logging-log-targets") {
			loggingLogTargets = oc.LoggingLogTargets
		}
		var loggingLokiTenantID *string
		if cmd.Flags().Changed("logging-loki-tenant-id") {
			loggingLokiTenantID = &oc.LoggingLokiTenantID
		}
		var loggingLokiBasicAuth *string
		if cmd.Flags().Changed("logging-loki-basic-auth") {
			loggingLokiBasicAuth = &oc.LoggingLokiBasicAuth
		}
		var loggingLokiEndpoint *string
		if cmd.Flags().Changed("logging-loki-endpoint") {
			loggingLokiEndpoint = &oc.LoggingLokiEndpoint
		}
		var loggingGrafanaBaseURL *string
		if cmd.Flags().Changed("logging-grafana-base-url") {
			loggingGrafanaBaseURL = &oc.LoggingGrafanaBaseURL
		}
		var loggingGrafanaDashboardURL *string
		if cmd.Flags().Changed("logging-grafana-dashboard-url") {
			loggingGrafanaDashboardURL = &oc.LoggingGrafanaDashboardURL
		}
		var loggingGrafanaToken *string
		if cmd.Flags().Changed("logging-grafana-token") {
			loggingGrafanaToken = &oc.LoggingGrafanaToken
		}

		if testLogCollect != nil || loggingRunID != nil || loggingLogTargets != nil || loggingLokiEndpoint != nil || loggingLokiTenantID != nil || loggingLokiBasicAuth != nil || loggingGrafanaBaseURL != nil || loggingGrafanaDashboardURL != nil || loggingGrafanaToken != nil {
			config.Logging = &ctf_config.LoggingConfig{}
			config.Logging.TestLogCollect = testLogCollect
			config.Logging.RunId = loggingRunID
			if loggingLogTargets != nil {
				config.Logging.LogStream = &ctf_config.LogStreamConfig{
					LogTargets: loggingLogTargets,
				}
			}
			if loggingLokiTenantID != nil || loggingLokiBasicAuth != nil || loggingLokiEndpoint != nil {
				config.Logging.Loki = &ctf_config.LokiConfig{
					TenantId:  loggingLokiTenantID,
					BasicAuth: loggingLokiBasicAuth,
					Endpoint:  loggingLokiEndpoint,
				}
			}
			if loggingGrafanaBaseURL != nil || loggingGrafanaDashboardURL != nil || loggingGrafanaToken != nil {
				config.Logging.Grafana = &ctf_config.GrafanaConfig{
					BaseUrl:      loggingGrafanaBaseURL,
					DashboardUrl: loggingGrafanaDashboardURL,
					BearerToken:  loggingGrafanaToken,
				}
			}
		}

		var privateEthereumNetworkExecutionLayer *string
		if cmd.Flags().Changed("private-ethereum-network-execution-layer") {
			privateEthereumNetworkExecutionLayer = &oc.PrivateEthereumNetworkExecutionLayer
		}
		var privateEthereumNetworkEthereumVersion *string
		if cmd.Flags().Changed("private-ethereum-network-ethereum-version") {
			privateEthereumNetworkEthereumVersion = &oc.PrivateEthereumNetworkEthereumVersion
		}
		var privateEthereumNetworkCustomDockerImage *string
		if cmd.Flags().Changed("private-ethereum-network-custom-docker-image") {
			privateEthereumNetworkCustomDockerImage = &oc.PrivateEthereumNetworkCustomDockerImages
		}
		if privateEthereumNetworkExecutionLayer != nil || privateEthereumNetworkEthereumVersion != nil || privateEthereumNetworkCustomDockerImage != nil {
			var el ctf_config.ExecutionLayer
			if privateEthereumNetworkExecutionLayer != nil {
				el = ctf_config.ExecutionLayer(*privateEthereumNetworkExecutionLayer)
			}
			var ev ctf_config.EthereumVersion
			if privateEthereumNetworkEthereumVersion != nil {
				ev = ctf_config.EthereumVersion(*privateEthereumNetworkEthereumVersion)
			}
			var customImages map[ctf_config.ContainerType]string
			if privateEthereumNetworkCustomDockerImage != nil {
				customImages = map[ctf_config.ContainerType]string{"execution_layer": *privateEthereumNetworkCustomDockerImage}
			}
			config.PrivateEthereumNetwork = &ctf_config.EthereumNetworkConfig{
				ExecutionLayer:     &el,
				EthereumVersion:    &ev,
				CustomDockerImages: customImages,
			}
		}

		configToml, err := toml.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling TestConfig to TOML: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(configToml))
	},
}

var maskTestConfigOverrideCmd = &cobra.Command{
	Use:   "mask-test-config-override",
	Short: "Processes configurations and outputs templated results for test config overrides",
	Run: func(cmd *cobra.Command, args []string) {
		maskSecret("chainlink image", oc.ChainlinkImage)
		maskSecret("pyroscope server url", oc.PyroscopeServerURL)
		maskSecret("pyroscope key", oc.PyroscopeKey)
		maskSecret("loki tenant id", oc.LoggingLokiTenantID)
		maskSecret("loki endpoint", oc.LoggingLokiEndpoint)
		maskSecret("loki basic auth", oc.LoggingLokiBasicAuth)
		maskSecret("loki grafana token", oc.LoggingGrafanaToken)
	},
}

func maskSecret(description string, secret string) {
	if secret != "" {
		fmt.Printf("# Mask '%s' value\n", description)
		fmt.Printf("echo ::add-mask::%s\n", secret)
	}
}

func init() {
	cmds := []*cobra.Command{createTestConfigOverrideCmd, maskTestConfigOverrideCmd}
	for _, c := range cmds {
		c.Flags().StringArrayVar(&oc.SelectedNetworks, "selected-networks", nil, "Selected networks")
		c.Flags().StringVar(&oc.ChainlinkImage, "chainlink-image", "", "Chainlink image")
		c.Flags().StringVar(&oc.ChainlinkVersion, "chainlink-version", "", "Chainlink version")
		c.Flags().StringVar(&oc.ChainlinkPostgresVersion, "chainlink-postgres-version", "", "Chainlink Postgres version")
		c.Flags().BoolVar(&oc.PyroscopeEnabled, "pyroscope-enabled", false, "Pyroscope enabled")
		c.Flags().StringVar(&oc.PyroscopeServerURL, "pyroscope-server-url", "", "Pyroscope server URL")
		c.Flags().StringVar(&oc.PyroscopeKey, "pyroscope-key", "", "Pyroscope key")
		c.Flags().StringVar(&oc.PyroscopeEnvironment, "pyroscope-environment", "", "Pyroscope environment")
		c.Flags().BoolVar(&oc.LoggingTestLogCollect, "logging-test-log-collect", false, "Test log collect")
		c.Flags().StringVar(&oc.LoggingRunID, "logging-run-id", "", "Run ID")
		c.Flags().StringArrayVar(&oc.LoggingLogTargets, "logging-log-targets", nil, "Logging.LogStream.LogTargets")
		c.Flags().StringVar(&oc.LoggingLokiEndpoint, "logging-loki-endpoint", "", "")
		c.Flags().StringVar(&oc.LoggingLokiTenantID, "logging-loki-tenant-id", "", "")
		c.Flags().StringVar(&oc.LoggingLokiBasicAuth, "logging-loki-basic-auth", "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaBaseURL, "logging-grafana-base-url", "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaDashboardURL, "logging-grafana-dashboard-url", "", "")
		c.Flags().StringVar(&oc.LoggingGrafanaToken, "logging-grafana-token", "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkExecutionLayer, "private-ethereum-network-execution-layer", "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkEthereumVersion, "private-ethereum-network-ethereum-version", "", "")
		c.Flags().StringVar(&oc.PrivateEthereumNetworkCustomDockerImages, "private-ethereum-network-custom-docker-image", "", "")

		c.PreRun = func(cmd *cobra.Command, args []string) {
			// Resolve selected networks environment variable if set
			_, hasEnvVar := lookupEnvVarName(oc.SelectedNetworks[0])
			if hasEnvVar {
				selectedNetworks := mustResolveEnvPlaceholder(oc.SelectedNetworks[0])
				oc.SelectedNetworks = strings.Split(selectedNetworks, ",")
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
			oc.LoggingGrafanaToken = mustResolveEnvPlaceholder(oc.LoggingGrafanaToken)
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
