package utils

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
)

// MergeSethAndEvmNetworkConfigs merges EVMNetwork to Seth config. If Seth config already has Network settings,
// it will return unchanged Seth config that was passed to it. If the network is simulated, it will
// use Geth-specific settings. Otherwise it will use the chain ID to find the correct network settings.
// If no match is found it will use default settings (currently based on Sepolia network settings).
func MergeSethAndEvmNetworkConfigs(l zerolog.Logger, evmNetwork blockchain.EVMNetwork, sethConfig seth.Config) seth.Config {
	if sethConfig.Network != nil {
		return sethConfig
	}

	var sethNetwork *seth.Network

	for _, conf := range sethConfig.Networks {
		if evmNetwork.Simulated {
			if conf.Name == seth.GETH {
				conf.PrivateKeys = evmNetwork.PrivateKeys
				conf.URLs = evmNetwork.URLs
				// important since Besu doesn't support EIP-1559, but other EVM clients do
				conf.EIP1559DynamicFees = evmNetwork.SupportsEIP1559

				sethNetwork = conf
				break
			}
		} else if conf.ChainID == fmt.Sprint(evmNetwork.ChainID) {
			conf.PrivateKeys = evmNetwork.PrivateKeys
			conf.URLs = evmNetwork.URLs

			sethNetwork = conf
			break
		}
	}

	if sethNetwork == nil {
		//TODO in the future we could run gas estimator here
		l.Warn().
			Int64("chainID", evmNetwork.ChainID).
			Msg("Could not find any Seth network settings for chain ID. Using default network settings")
		sethNetwork = &seth.Network{}
		sethNetwork.PrivateKeys = evmNetwork.PrivateKeys
		sethNetwork.URLs = evmNetwork.URLs
		sethNetwork.EIP1559DynamicFees = evmNetwork.SupportsEIP1559
		sethNetwork.ChainID = fmt.Sprint(evmNetwork.ChainID)
		// Sepolia settings
		sethNetwork.GasLimit = 14_000_000
		sethNetwork.GasPrice = 1_000_000_000
		sethNetwork.GasFeeCap = 25_000_000_000
		sethNetwork.GasTipCap = 5_000_000_000
		sethNetwork.TransferGasFee = 21_000
		sethNetwork.TxnTimeout = seth.MustMakeDuration(evmNetwork.Timeout.Duration)
	}

	sethConfig.Network = sethNetwork

	return sethConfig
}

// MustReplaceSimulatedNetworkUrlWithK8 replaces the simulated network URL with the K8 URL and returns the network.
// If the network is not simulated, it will return the network unchanged.
func MustReplaceSimulatedNetworkUrlWithK8(l zerolog.Logger, network blockchain.EVMNetwork, testEnvironment environment.Environment) blockchain.EVMNetwork {
	if !network.Simulated {
		return network
	}

	if _, ok := testEnvironment.URLs["Simulated Geth"]; !ok {
		for k := range testEnvironment.URLs {
			l.Info().Str("Network", k).Msg("Available networks")
		}
		panic("no network settings for Simulated Geth")
	}
	network.URLs = testEnvironment.URLs["Simulated Geth"]

	return network
}
