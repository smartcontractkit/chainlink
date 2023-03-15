package client

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

const pyroscopeTOML = `[Pyroscope]
ServerAddress = '%s'
Environment = '%s'`

// AddNetworksConfig adds EVM network configurations to a base config TOML. Useful for adding networks with default
// settings. See AddNetworkDetailedConfig for adding more detailed network configuration.
func AddNetworksConfig(baseTOML string, networks ...blockchain.EVMNetwork) string {
	networksToml := ""
	for _, network := range networks {
		networksToml = fmt.Sprintf("%s\n\n%s", networksToml, network.MustChainlinkTOML(""))
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", baseTOML, pyroscopeSettings(), networksToml)
}

// AddNetworkDetailedConfig adds EVM config to a base TOML. Also takes a detailed network config TOML where values like
// using transaction forwarders can be included.
// See https://github.com/smartcontractkit/chainlink/blob/develop/docs/CONFIG.md#EVM
func AddNetworkDetailedConfig(baseTOML, detailedNetworkConfig string, network blockchain.EVMNetwork) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", baseTOML, pyroscopeSettings(), network.MustChainlinkTOML(detailedNetworkConfig))
}

func pyroscopeSettings() string {
	pyroscopeServer := os.Getenv(config.EnvVarPyroscopeServer)
	pyroscopeEnv := os.Getenv(config.EnvVarPyroscopeEnvironment)
	if pyroscopeServer == "" {
		return ""
	}
	return fmt.Sprintf(pyroscopeTOML, pyroscopeServer, pyroscopeEnv)
}
