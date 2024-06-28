package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const DryRunFlag = "dry-run"

// fromBase64TOML is a base64 encoded TOML config to override
var fromBase64TOML string

var overrideTestConfigCmd = &cobra.Command{
	Use:   "override",
	Short: "Override base64 encoded TOML config with provided flags. Overrides only existing fields in the base config.",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool(DryRunFlag)

		decoded, err := base64.StdEncoding.DecodeString(fromBase64TOML)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding base64 TOML: %v\n", err)
			os.Exit(1)
		}
		var baseConfig ctf_config.TestConfig
		err = toml.Unmarshal(decoded, &baseConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling base64 TOML: %v\n", err)
			os.Exit(1)
		}

		cmd.Flags().Visit(func(flag *pflag.Flag) {
			switch flag.Name {
			case ChainlinkImageFlag:
				if baseConfig.ChainlinkImage != nil && baseConfig.ChainlinkImage.Image != nil {
					logIfDryRun(dryRun, "Found 'ChainlinkImage.Image' in config. Will override it with %s\n", ChainlinkImageFlag)
					baseConfig.ChainlinkImage.Image = &oc.ChainlinkImage
				} else {
					logIfDryRun(dryRun, "No 'ChainlinkImage.Image' in config. Will NOT OVERRIDE with %s\n", ChainlinkImageFlag)
				}
			case ChainlinkVersionFlag:
				if baseConfig.ChainlinkImage != nil && baseConfig.ChainlinkImage.Version != nil {
					logIfDryRun(dryRun, "Found 'ChainlinkImage.Version' in config. Will override it with %s\n", ChainlinkVersionFlag)
					baseConfig.ChainlinkImage.Version = &oc.ChainlinkVersion
				} else {
					logIfDryRun(dryRun, "No 'ChainlinkImage.Version' found in config. Will NOT OVERRIDE with %s\n", ChainlinkVersionFlag)
				}
			case ChainlinkPostgresVersionFlag:
				if baseConfig.ChainlinkImage != nil && baseConfig.ChainlinkImage.PostgresVersion != nil {
					logIfDryRun(dryRun, "Found 'ChainlinkImage.PostgresVersion' with %s\n", ChainlinkPostgresVersionFlag)
					baseConfig.ChainlinkImage.PostgresVersion = &oc.ChainlinkPostgresVersion
				} else {
					logIfDryRun(dryRun, "No 'ChainlinkImage.PostgresVersion' found in config. Will NOT OVERRIDE with %s\n", ChainlinkPostgresVersionFlag)
				}
			case SelectedNetworksFlag:
				if baseConfig.Network != nil && baseConfig.Network.SelectedNetworks != nil {
					logIfDryRun(dryRun, "Found 'Network.SelectedNetworks' in config. Will override it with %s", SelectedNetworksFlag)
					baseConfig.Network.SelectedNetworks = oc.SelectedNetworks
				} else {
					logIfDryRun(dryRun, "No 'Network.SelectedNetworks' found in config. Will NOT OVERRIDE with %s\n", SelectedNetworksFlag)
				}
			case LoggingLokiBasicAuthFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Loki != nil && baseConfig.Logging.Loki.BasicAuth != nil {
					logIfDryRun(dryRun, "Found 'Logging.Loki.BasicAuth' in config. Will override it with %s", LoggingLokiBasicAuthFlag)
					baseConfig.Logging.Loki.BasicAuth = &oc.LoggingLokiBasicAuth
				} else {
					logIfDryRun(dryRun, "No 'Logging.Loki' found in config. Will NOT OVERRIDE with %s\n", LoggingLokiBasicAuthFlag)
				}
			case LoggingLokiEndpointFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Loki != nil && baseConfig.Logging.Loki.Endpoint != nil {
					logIfDryRun(dryRun, "Found 'Logging.Loki.Endpoint' in config. Will override it with %s", LoggingLokiEndpointFlag)
					baseConfig.Logging.Loki.Endpoint = &oc.LoggingLokiEndpoint
				} else {
					logIfDryRun(dryRun, "No 'Logging.Loki' found in config. Will NOT OVERRIDE with %s\n", LoggingLokiEndpointFlag)
				}
			case LoggingLokiTenantIDFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Loki != nil && baseConfig.Logging.Loki.TenantId != nil {
					logIfDryRun(dryRun, "Found 'Logging.Loki.TenantId' in config. Will override it with %s", LoggingLokiTenantIDFlag)
					baseConfig.Logging.Loki.TenantId = &oc.LoggingLokiTenantID
				} else {
					logIfDryRun(dryRun, "No 'Logging.Loki' found in config. Will NOT OVERRIDE with %s\n", LoggingLokiTenantIDFlag)
				}
			case LoggingRunIDFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.RunId != nil {
					logIfDryRun(dryRun, "Found 'Logging.Loki.RunId' in config. Will override it with %s", LoggingRunIDFlag)
					baseConfig.Logging.RunId = &oc.LoggingRunID
				} else {
					logIfDryRun(dryRun, "No 'Logging' found in config. Will NOT OVERRIDE with %s\n", LoggingRunIDFlag)
				}
			case LoggingGrafanaBaseURLFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Grafana != nil && baseConfig.Logging.Grafana.BaseUrl != nil {
					logIfDryRun(dryRun, "Found 'Logging.Grafana.BaseUrl' in config. Will override it with %s", LoggingGrafanaBaseURLFlag)
					baseConfig.Logging.Grafana.BaseUrl = &oc.LoggingGrafanaBaseURL
				} else {
					logIfDryRun(dryRun, "No 'Logging.Grafana' found in config. Will NOT OVERRIDE with %s\n", LoggingGrafanaBaseURLFlag)
				}
			case LoggingGrafanaDashboardURLFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Grafana != nil && baseConfig.Logging.Grafana.DashboardUrl != nil {
					logIfDryRun(dryRun, "Found 'Logging.Grafana.BaseUrl' in config. Will override it with %s", LoggingGrafanaBaseURLFlag)
					baseConfig.Logging.Grafana.DashboardUrl = &oc.LoggingGrafanaDashboardURL
				} else {
					logIfDryRun(dryRun, "No 'Logging.Grafana' found in config. Will NOT OVERRIDE with %s\n", LoggingGrafanaDashboardURLFlag)
				}
			case LoggingGrafanaTokenFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.Grafana != nil && baseConfig.Logging.Grafana.BearerToken != nil {
					logIfDryRun(dryRun, "Found 'Logging.Grafana.BearerToken' in config. Will override it with %s", LoggingGrafanaTokenFlag)
					baseConfig.Logging.Grafana.BearerToken = &oc.LoggingGrafanaToken
				} else {
					logIfDryRun(dryRun, "No 'Logging.Grafana' found in config. Will NOT OVERRIDE with %s\n", LoggingGrafanaTokenFlag)
				}
			case LoggingLogTargetsFlag:
				if baseConfig.Logging != nil && baseConfig.Logging.LogStream != nil && baseConfig.Logging.LogStream.LogTargets != nil {
					logIfDryRun(dryRun, "Found 'Logging.LogStream.LogTargets' in config. Will override it with %s", LoggingLogTargetsFlag)
					baseConfig.Logging.LogStream.LogTargets = oc.LoggingLogTargets
				} else {
					logIfDryRun(dryRun, "No 'Logging.LogStream' found in config. Will NOT OVERRIDE with %s\n", LoggingLogTargetsFlag)
				}
			case PyroscopeEnabledFlag:
				if baseConfig.Pyroscope != nil && baseConfig.Pyroscope.Enabled != nil {
					logIfDryRun(dryRun, "Found 'Pyroscope.Enabled' in config. Will override it with %s", PyroscopeEnabledFlag)
					baseConfig.Pyroscope.Enabled = &oc.PyroscopeEnabled
				} else {
					logIfDryRun(dryRun, "No 'Pyroscope' found in config. Will NOT OVERRIDE with %s\n", PyroscopeEnabledFlag)
				}
			case PyroscopeServerURLFlag:
				if baseConfig.Pyroscope != nil {
					logIfDryRun(dryRun, "Found 'Pyroscope.ServerUrl' in config. Will override it with %s", PyroscopeServerURLFlag)
					baseConfig.Pyroscope.ServerUrl = &oc.PyroscopeServerURL
				} else {
					logIfDryRun(dryRun, "No 'Pyroscope' found in config. Will NOT OVERRIDE with %s\n", PyroscopeServerURLFlag)
				}
			case PyroscopeEnvironmentFlag:
				if baseConfig.Pyroscope != nil {
					logIfDryRun(dryRun, "Found 'Pyroscope.Environment' in config. Will override it with %s", PyroscopeEnvironmentFlag)
					baseConfig.Pyroscope.Environment = &oc.PyroscopeEnvironment
				} else {
					logIfDryRun(dryRun, "No 'Pyroscope' found in config. Will NOT OVERRIDE with %s\n", PyroscopeEnvironmentFlag)
				}
			case PyroscopeKeyFlag:
				if baseConfig.Pyroscope != nil {
					logIfDryRun(dryRun, "Found 'Pyroscope.Key' in config. Will override it with %s", PyroscopeKeyFlag)
					baseConfig.Pyroscope.Key = &oc.PyroscopeKey
				} else {
					logIfDryRun(dryRun, "No 'Pyroscope' found in config. Will NOT OVERRIDE with %s\n", PyroscopeKeyFlag)
				}
			case PrivateEthereumNetworkExecutionLayerFlag:
				if baseConfig.PrivateEthereumNetwork != nil {
					logIfDryRun(dryRun, "Found 'PrivateEthereumNetwork.ExecutionLayer' in config. Will override it with %s", PrivateEthereumNetworkExecutionLayerFlag)
					el := ctf_config.ExecutionLayer(oc.PrivateEthereumNetworkExecutionLayer)
					baseConfig.PrivateEthereumNetwork.ExecutionLayer = &el
				} else {
					logIfDryRun(dryRun, "No 'PrivateEthereumNetwork' found in config. Will NOT OVERRIDE with %s\n", PrivateEthereumNetworkExecutionLayerFlag)
				}
			case PrivateEthereumNetworkEthereumVersionFlag:
				if baseConfig.PrivateEthereumNetwork != nil {
					logIfDryRun(dryRun, "Found 'PrivateEthereumNetwork.EthereumVersion' in config. Will override it with %s", PrivateEthereumNetworkEthereumVersionFlag)
					ev := ctf_config.EthereumVersion(oc.PrivateEthereumNetworkEthereumVersion)
					baseConfig.PrivateEthereumNetwork.EthereumVersion = &ev
				} else {
					logIfDryRun(dryRun, "No 'PrivateEthereumNetwork' found in config. Will NOT OVERRIDE with %s\n", PrivateEthereumNetworkEthereumVersionFlag)
				}
			case PrivateEthereumNetworkCustomDockerImageFlag:
				if baseConfig.PrivateEthereumNetwork != nil {
					logIfDryRun(dryRun, "Found 'PrivateEthereumNetwork.CustomDockerImages' in config. Will override it with %s", PrivateEthereumNetworkCustomDockerImageFlag)
					customImages := map[ctf_config.ContainerType]string{"execution_layer": oc.PrivateEthereumNetworkCustomDockerImages}
					baseConfig.PrivateEthereumNetwork.CustomDockerImages = customImages
				} else {
					logIfDryRun(dryRun, "No 'PrivateEthereumNetwork' found in config. Will NOT OVERRIDE with %s\n", PrivateEthereumNetworkCustomDockerImageFlag)
				}
			default:
				fmt.Printf("Override not supported for flag: %s\n", flag.Name)
				os.Exit(1)
			}
		})

		if !dryRun {
			configToml, err := toml.Marshal(baseConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshalling TestConfig to TOML: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(configToml))
		}
	},
}

func init() {
	overrideTestConfigCmd.Flags().StringVar(&fromBase64TOML, FromBase64ConfigFlag, "", "Base64 encoded TOML config to override")
	overrideTestConfigCmd.MarkFlagRequired(FromBase64ConfigFlag)

	overrideTestConfigCmd.Flags().Bool(DryRunFlag, false, "Dry run mode")
}

func logIfDryRun(dryRun bool, format string, a ...interface{}) {
	if dryRun {
		fmt.Printf(format, a...)
	}
}
