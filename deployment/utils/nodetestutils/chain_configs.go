package nodetestutils

import (
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// createSuiChainConfig creates a Chainlink node defined Sui chain configuration.
func createSuiChainConfig(chainID string, chain cldf_sui.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "sui-localnet"
	chainConfig["NetworkNameFull"] = "sui-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}

// createTronChainConfig creates a Chainlink node defined Tron chain configuration.
func createTronChainConfig(chainID string, chain cldf_tron.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "tron-local"
	chainConfig["NetworkNameFull"] = "tron-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}

// createAptosChainConfig creates a Chainlink node defined Aptos chain configuration.
func createAptosChainConfig(chainID string, chain cldf_aptos.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "localnet"
	chainConfig["NetworkNameFull"] = "aptos-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
