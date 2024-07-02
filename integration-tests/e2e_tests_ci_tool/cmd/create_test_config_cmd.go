package cmd

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
)

var createTestConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a test config from the provided flags",
	Run: func(cmd *cobra.Command, _ []string) {
		var tc ctf_config.TestConfig

		var image, version, postgresVersion *string
		if cmd.Flags().Changed(ChainlinkImageFlag) {
			image = &oc.ChainlinkImage
		}
		if cmd.Flags().Changed(ChainlinkVersionFlag) {
			version = &oc.ChainlinkVersion
		}
		if cmd.Flags().Changed(ChainlinkPostgresVersionFlag) {
			version = &oc.ChainlinkPostgresVersion
		}
		if image != nil && version == nil || image == nil && version != nil {
			fmt.Fprintf(os.Stderr, "Error: both chainlink-image and chainlink-version must be set\n")
			os.Exit(1)
		}
		if image != nil && version != nil {
			tc.ChainlinkImage = &ctf_config.ChainlinkImageConfig{
				Image:           image,
				Version:         version,
				PostgresVersion: postgresVersion,
			}
		}

		var upgradeImage, upgradeVersion *string
		if cmd.Flags().Changed(ChainlinkUpgradeImageFlag) {
			upgradeImage = &oc.ChainlinkUpgradeImage
		}
		if cmd.Flags().Changed(ChainlinkUpgradeVersionFlag) {
			upgradeVersion = &oc.ChainlinkUpgradeVersion
		}
		if upgradeImage != nil || upgradeVersion != nil {
			tc.ChainlinkUpgradeImage = &ctf_config.ChainlinkImageConfig{
				Image:   upgradeImage,
				Version: upgradeVersion,
			}
		}

		var selectedNetworks *[]string
		if cmd.Flags().Changed(SelectedNetworksFlag) {
			selectedNetworks = &oc.SelectedNetworks
		}
		if selectedNetworks != nil {
			tc.Network = &ctf_config.NetworkConfig{
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
			tc.Pyroscope = &ctf_config.PyroscopeConfig{
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
			tc.Logging = &ctf_config.LoggingConfig{}
			tc.Logging.TestLogCollect = testLogCollect
			tc.Logging.RunId = loggingRunID
			if loggingLogTargets != nil {
				tc.Logging.LogStream = &ctf_config.LogStreamConfig{
					LogTargets: loggingLogTargets,
				}
			}
			if loggingLokiTenantID != nil || loggingLokiBasicAuth != nil || loggingLokiEndpoint != nil {
				tc.Logging.Loki = &ctf_config.LokiConfig{
					TenantId:  loggingLokiTenantID,
					BasicAuth: loggingLokiBasicAuth,
					Endpoint:  loggingLokiEndpoint,
				}
			}
			if loggingGrafanaBaseURL != nil || loggingGrafanaDashboardURL != nil || loggingGrafanaBearerToken != nil {
				tc.Logging.Grafana = &ctf_config.GrafanaConfig{
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
			tc.PrivateEthereumNetwork = &ctf_config.EthereumNetworkConfig{
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
