package cmd

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctfconfigtypes "github.com/smartcontractkit/chainlink-testing-framework/config/types"
)

var createTestConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a test config from the provided flags",
	Run: func(cmd *cobra.Command, _ []string) {
		var tc ctfconfig.TestConfig

		var version, postgresVersion *string
		if cmd.Flags().Changed(ChainlinkVersionFlag) {
			version = &oc.ChainlinkVersion
		}
		if cmd.Flags().Changed(ChainlinkPostgresVersionFlag) {
			version = &oc.ChainlinkPostgresVersion
		}
		if version != nil || postgresVersion != nil {
			tc.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{
				Version:         version,
				PostgresVersion: postgresVersion,
			}
		}

		var upgradeVersion *string
		if cmd.Flags().Changed(ChainlinkUpgradeVersionFlag) {
			upgradeVersion = &oc.ChainlinkUpgradeVersion
		}
		if upgradeVersion != nil {
			tc.ChainlinkUpgradeImage = &ctfconfig.ChainlinkImageConfig{
				Version: upgradeVersion,
			}
		}

		var selectedNetworks *[]string
		if cmd.Flags().Changed(SelectedNetworksFlag) {
			selectedNetworks = &oc.SelectedNetworks
		}
		if selectedNetworks != nil {
			tc.Network = &ctfconfig.NetworkConfig{
				SelectedNetworks: oc.SelectedNetworks,
			}
		}

		var peryscopeEnabled *bool
		var pyroscopeServerURL, pyroscopeEnvironment, pyroscopeKey *string
		if cmd.Flags().Changed(PyroscopeEnabledFlag) {
			peryscopeEnabled = &oc.PyroscopeEnabled
		}
		if cmd.Flags().Changed(PyroscopeServerURLFlag) {
			pyroscopeServerURL = &oc.PyroscopeServerURL
		}
		if cmd.Flags().Changed(PyroscopeKeyFlag) {
			pyroscopeKey = &oc.PyroscopeKey
		}
		if cmd.Flags().Changed(PyroscopeEnvironmentFlag) {
			pyroscopeEnvironment = &oc.PyroscopeEnvironment
		}
		if peryscopeEnabled != nil {
			tc.Pyroscope = &ctfconfig.PyroscopeConfig{
				Enabled:     peryscopeEnabled,
				ServerUrl:   pyroscopeServerURL,
				Environment: pyroscopeEnvironment,
				Key:         pyroscopeKey,
			}
		}

		var testLogCollect *bool
		if cmd.Flags().Changed(LoggingTestLogCollectFlag) {
			testLogCollect = &oc.LoggingTestLogCollect
		}
		var loggingRunID *string
		if cmd.Flags().Changed(LoggingRunIDFlag) {
			loggingRunID = &oc.LoggingRunID
		}
		var loggingLogTargets []string
		if cmd.Flags().Changed(LoggingLogTargetsFlag) {
			loggingLogTargets = oc.LoggingLogTargets
		}
		var loggingLokiTenantID *string
		if cmd.Flags().Changed(LoggingLokiTenantIDFlag) {
			loggingLokiTenantID = &oc.LoggingLokiTenantID
		}
		var loggingLokiBasicAuth *string
		if cmd.Flags().Changed(LoggingLokiBasicAuthFlag) {
			loggingLokiBasicAuth = &oc.LoggingLokiBasicAuth
		}
		var loggingLokiEndpoint *string
		if cmd.Flags().Changed(LoggingLokiEndpointFlag) {
			loggingLokiEndpoint = &oc.LoggingLokiEndpoint
		}
		var loggingGrafanaBaseURL *string
		if cmd.Flags().Changed(LoggingGrafanaBaseURLFlag) {
			loggingGrafanaBaseURL = &oc.LoggingGrafanaBaseURL
		}
		var loggingGrafanaDashboardURL *string
		if cmd.Flags().Changed(LoggingGrafanaDashboardURLFlag) {
			loggingGrafanaDashboardURL = &oc.LoggingGrafanaDashboardURL
		}
		var loggingGrafanaBearerToken *string
		if cmd.Flags().Changed(LoggingGrafanaBearerTokenFlag) {
			loggingGrafanaBearerToken = &oc.LoggingGrafanaBearerToken
		}

		if testLogCollect != nil || loggingRunID != nil || loggingLogTargets != nil || loggingLokiEndpoint != nil || loggingLokiTenantID != nil || loggingLokiBasicAuth != nil || loggingGrafanaBaseURL != nil || loggingGrafanaDashboardURL != nil || loggingGrafanaBearerToken != nil {
			tc.Logging = &ctfconfig.LoggingConfig{}
			tc.Logging.TestLogCollect = testLogCollect
			tc.Logging.RunId = loggingRunID
			if loggingLogTargets != nil {
				tc.Logging.LogStream = &ctfconfig.LogStreamConfig{
					LogTargets: loggingLogTargets,
				}
			}
			if loggingLokiTenantID != nil || loggingLokiBasicAuth != nil || loggingLokiEndpoint != nil {
				tc.Logging.Loki = &ctfconfig.LokiConfig{
					TenantId:  loggingLokiTenantID,
					BasicAuth: loggingLokiBasicAuth,
					Endpoint:  loggingLokiEndpoint,
				}
			}
			if loggingGrafanaBaseURL != nil || loggingGrafanaDashboardURL != nil || loggingGrafanaBearerToken != nil {
				tc.Logging.Grafana = &ctfconfig.GrafanaConfig{
					BaseUrl:      loggingGrafanaBaseURL,
					DashboardUrl: loggingGrafanaDashboardURL,
					BearerToken:  loggingGrafanaBearerToken,
				}
			}
		}

		var privateEthereumNetworkExecutionLayer *string
		if cmd.Flags().Changed(PrivateEthereumNetworkExecutionLayerFlag) {
			privateEthereumNetworkExecutionLayer = &oc.PrivateEthereumNetworkExecutionLayer
		}
		var privateEthereumNetworkEthereumVersion *string
		if cmd.Flags().Changed(PrivateEthereumNetworkEthereumVersionFlag) {
			privateEthereumNetworkEthereumVersion = &oc.PrivateEthereumNetworkEthereumVersion
		}
		var privateEthereumNetworkCustomDockerImage *string
		if cmd.Flags().Changed(PrivateEthereumNetworkCustomDockerImageFlag) {
			privateEthereumNetworkCustomDockerImage = &oc.PrivateEthereumNetworkCustomDockerImages
		}
		if privateEthereumNetworkExecutionLayer != nil || privateEthereumNetworkEthereumVersion != nil || privateEthereumNetworkCustomDockerImage != nil {
			var el ctfconfigtypes.ExecutionLayer
			if privateEthereumNetworkExecutionLayer != nil {
				el = ctfconfigtypes.ExecutionLayer(*privateEthereumNetworkExecutionLayer)
			}
			var ev ctfconfigtypes.EthereumVersion
			if privateEthereumNetworkEthereumVersion != nil {
				ev = ctfconfigtypes.EthereumVersion(*privateEthereumNetworkEthereumVersion)
			}
			var customImages map[ctfconfig.ContainerType]string
			if privateEthereumNetworkCustomDockerImage != nil {
				customImages = map[ctfconfig.ContainerType]string{"execution_layer": *privateEthereumNetworkCustomDockerImage}
			}
			tc.PrivateEthereumNetwork = &ctfconfig.EthereumNetworkConfig{
				ExecutionLayer:     &el,
				EthereumVersion:    &ev,
				CustomDockerImages: customImages,
			}
		}

		configToml, err := toml.Marshal(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling TestConfig to TOML: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(configToml))
	},
}
