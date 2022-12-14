package client

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// AddNetworksConfig adds EVM network configurations to a base config TOML. Useful for adding networks with default
// settings. See AddNetworkDetailedConfig for adding more detailed network configuration.
func AddNetworksConfig(baseTOML string, networks ...blockchain.EVMNetwork) string {
	networksToml := ""
	for _, network := range networks {
		networksToml = fmt.Sprintf("%s\n\n%s", networksToml, network.MustChainlinkTOML(""))
	}
	return fmt.Sprintf("%s\n\n%s", baseTOML, networksToml)
}

// AddNetworkDetailedConfig adds EVM config to a base TOML. Also takes a detailed network config TOML where values like
// using transaction forwarders can be included.
// See https://github.com/smartcontractkit/chainlink/blob/develop/docs/CONFIG.md#EVM
func AddNetworkDetailedConfig(baseTOML, detailedNetworkConfig string, network blockchain.EVMNetwork) string {
	return fmt.Sprintf("%s\n\n%s", baseTOML, network.MustChainlinkTOML(detailedNetworkConfig))
}
