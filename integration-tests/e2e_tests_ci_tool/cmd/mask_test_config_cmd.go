package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/seth"
	"github.com/spf13/cobra"
)

var maskTestConfigCmd = &cobra.Command{
	Use:   "mask-secrets",
	Short: "Run 'echo ::add-mask::${secret}' for all secrets in the test config",
	Run: func(cmd *cobra.Command, args []string) {
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
			mustMaskSecret("chainlink image", config.ChainlinkImage.Image)
			mustMaskSecret("chainlink version", config.ChainlinkImage.Version)
		}
		if config.Pyroscope != nil {
			mustMaskSecret("pyroscope server url", config.Pyroscope.ServerUrl)
			mustMaskSecret("pyroscope key", config.Pyroscope.Key)
		}
		if config.Logging != nil && config.Logging.Loki != nil {
			// mustMaskSecret("loki tenant id", config.Logging.RunId)
			mustMaskSecret("loki endpoint", config.Logging.Loki.Endpoint)
			mustMaskSecret("loki tenant id", config.Logging.Loki.TenantId)
			mustMaskSecret("loki basic auth", config.Logging.Loki.BasicAuth)
			mustMaskSecret("loki bearer token", config.Logging.Loki.BearerToken)
		}
		if config.Logging != nil && config.Logging.Grafana != nil {
			mustMaskSecret("loki grafana token", config.Logging.Grafana.BearerToken)
		}
		if config.Network != nil && config.Network.RpcHttpUrls != nil {
			for _, urls := range config.Network.RpcHttpUrls {
				for _, url := range urls {
					mustMaskSecret("rpc http url", &url)
				}
			}
		}
		if config.Network != nil && config.Network.RpcWsUrls != nil {
			for _, urls := range config.Network.RpcWsUrls {
				for _, url := range urls {
					mustMaskSecret("rpc ws url", &url)
				}
			}
		}
		if config.Network != nil && config.Network.WalletKeys != nil {
			for _, keys := range config.Network.WalletKeys {
				for _, key := range keys {
					mustMaskSecret("wallet key", &key)
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
			mustMaskSecret("wasp repo image version uri", config.WaspConfig.RepoImageVersionURI)
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
			mustMaskSecret("network rpc url", &url)
		}
		for _, url := range network.URLs {
			mustMaskSecret("network ws url", &url)
		}
		for _, key := range network.PrivateKeys {
			mustMaskSecret("network private key", &key)
		}
		mustMaskSecret("network url", &network.URL)
	}
}

func mustMaskSethNetworkConfig(network *seth.Network) {
	if network != nil {
		for _, url := range network.URLs {
			mustMaskSecret("network url", &url)
		}
		for _, key := range network.PrivateKeys {
			mustMaskSecret("network private key", &key)
		}
	}
}

func mustMaskSecret(description string, secret *string) {
	if secret != nil && *secret != "" {
		fmt.Printf("Mask '%s'\n", description)
		fmt.Printf("::add-mask::%s\n", *secret)
		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo ::add-mask::%s", *secret))
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to mask secret '%s'", description)
			os.Exit(1)
		}
	}
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
	maskTestConfigCmd.MarkFlagRequired(FromBase64ConfigFlag)
}
