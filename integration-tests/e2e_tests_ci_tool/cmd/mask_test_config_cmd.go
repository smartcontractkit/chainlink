package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/seth"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
)

var maskTestConfigCmd = &cobra.Command{
	Use:   "mask-secrets",
	Short: "Run 'echo ::add-mask::${secret}' for all secrets in the test config",
	Run: func(_ *cobra.Command, _ []string) {
		decoded, err := base64.StdEncoding.DecodeString(fromBase64TOML)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding base64 TOML: %v\n", err)
			os.Exit(1)
		}
		var config ctf_config.TestConfig
		err = toml.Unmarshal(decoded, &config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling base64 TOML: %v\n", err)
			os.Exit(1)
		}

		showCmdInfo()

		if config.ChainlinkImage != nil {
			mustMaskSecret("chainlink image", safeDeref(config.ChainlinkImage.Image))
			mustMaskSecret("chainlink version", safeDeref(config.ChainlinkImage.Version))
		}
		if config.Pyroscope != nil {
			mustMaskSecret("pyroscope server url", safeDeref(config.Pyroscope.ServerUrl))
			mustMaskSecret("pyroscope key", safeDeref(config.Pyroscope.Key))
		}
		if config.Logging != nil && config.Logging.Loki != nil {
			// mustMaskSecret("loki tenant id", config.Logging.RunId)
			mustMaskSecret("loki endpoint", safeDeref(config.Logging.Loki.Endpoint))
			mustMaskSecret("loki tenant id", safeDeref(config.Logging.Loki.TenantId))
			mustMaskSecret("loki basic auth", safeDeref(config.Logging.Loki.BasicAuth))
			mustMaskSecret("loki bearer token", safeDeref(config.Logging.Loki.BearerToken))
		}
		if config.Logging != nil && config.Logging.Grafana != nil {
			mustMaskSecret("loki grafana token", safeDeref(config.Logging.Grafana.BearerToken))
		}
		if config.Network != nil && config.Network.RpcHttpUrls != nil {
			for _, urls := range config.Network.RpcHttpUrls {
				for _, url := range urls {
					mustMaskSecret("rpc http url", url)
				}
			}
		}
		if config.Network != nil && config.Network.RpcWsUrls != nil {
			for _, urls := range config.Network.RpcWsUrls {
				for _, url := range urls {
					mustMaskSecret("rpc ws url", url)
				}
			}
		}
		if config.Network != nil && config.Network.WalletKeys != nil {
			for _, keys := range config.Network.WalletKeys {
				for _, key := range keys {
					mustMaskSecret("wallet key", key)
				}
			}
		}
		// Mask EVMNetworks config
		if config.Network != nil && config.Network.EVMNetworks != nil {
			for _, network := range config.Network.EVMNetworks {
				mustMaskEVMNetworkConfig(network)
			}
		}
		if config.WaspConfig != nil {
			mustMaskSecret("wasp repo image version uri", safeDeref(config.WaspConfig.RepoImageVersionURI))
		}
		// Mask Seth config
		if config.Seth != nil {
			mustMaskSethNetworkConfig(config.Seth.Network)
			for _, network := range config.Seth.Networks {
				mustMaskSethNetworkConfig(network)
			}
		}
	},
}

func mustMaskEVMNetworkConfig(network *blockchain.EVMNetwork) {
	if network != nil {
		for _, url := range network.HTTPURLs {
			mustMaskSecret("network rpc url", url)
		}
		for _, url := range network.URLs {
			mustMaskSecret("network ws url", url)
		}
		for _, key := range network.PrivateKeys {
			mustMaskSecret("network private key", key)
		}
		mustMaskSecret("network url", network.URL)
	}
}

func mustMaskSethNetworkConfig(network *seth.Network) {
	if network != nil {
		for _, url := range network.URLs {
			mustMaskSecret("network url", url)
		}
		for _, key := range network.PrivateKeys {
			mustMaskSecret("network private key", key)
		}
	}
}

func mustMaskSecret(description string, secret string) {
	if secret != "" {
		fmt.Printf("Mask '%s'\n", description)
		fmt.Printf("::add-mask::%s\n", secret)
		cmd := exec.Command("bash", "-c", "echo ::add-mask::'$0'", "_", secret)
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to mask secret '%s'", description)
			os.Exit(1)
		}
	}
}

func safeDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func showCmdInfo() {
	data := [][]string{
		{"This command masks ONLY SELECTED secrets in the test config!"},
		{"Secrets inside Chainlink Node config are NOT masked"},
		{"Please ensure that are not present in the test config!"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"WARNING!"})
	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetColWidth(100)
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetTablePadding("\t")

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
	fmt.Println()
}

func init() {
	maskTestConfigCmd.Flags().StringVar(&fromBase64TOML, FromBase64ConfigFlag, "", "Base64 encoded TOML config to override")
	err := maskTestConfigCmd.MarkFlagRequired(FromBase64ConfigFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
