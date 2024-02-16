package testsetups

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/seth"
)

// MustDecorateSethConfigWithNetwork decorates Seth config with legacy EVMNetwork settings. If Seth config
// already has Network settings, it will return without doing anything. If the network is simulated, it will
// use Geth-specific settings. Otherwise it will use the chain ID to find the correct network settings.
// If no match is found it will return an error.
func MustDecorateSethConfigWithNetwork(evmNetwork *blockchain.EVMNetwork, sethConfig *seth.Config) error {
	if evmNetwork == nil {
		panic("evmNetwork must not be nil")
	}

	if sethConfig == nil {
		panic("sethConfig must not be nil")
	}

	if sethConfig.Network != nil {
		return nil
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

			sethConfig.Network = conf
			break
		}
	}

	if sethNetwork == nil {
		return fmt.Errorf("Could not find any Seth network settings for chain ID %d", evmNetwork.ChainID)
	}

	sethConfig.Network = sethNetwork

	return nil
}
