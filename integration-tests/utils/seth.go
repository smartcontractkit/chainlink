package utils

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
)

// MergeSethAndEvmNetworkConfigs merges EVMNetwork to Seth config. If Seth config already has Network settings,
// it will return unchanged Seth config that was passed to it. If the network is simulated, it will
// use Geth-specific settings. Otherwise it will use the chain ID to find the correct network settings.
// If no match is found it will return error.
func MergeSethAndEvmNetworkConfigs(evmNetwork blockchain.EVMNetwork, sethConfig seth.Config) (seth.Config, error) {
	if sethConfig.Network != nil {
		return sethConfig, nil
	}

	var sethNetwork *seth.Network

	for _, conf := range sethConfig.Networks {
		if evmNetwork.Simulated {
			if conf.Name == seth.GETH {
				conf.PrivateKeys = evmNetwork.PrivateKeys
				if len(conf.URLs) == 0 {
					conf.URLs = evmNetwork.URLs
				}
				// important since Besu doesn't support EIP-1559, but other EVM clients do
				conf.EIP1559DynamicFees = evmNetwork.SupportsEIP1559

				// might be needed for cases, when node is incapable of estimating gas limit (e.g. Geth < v1.10.0)
				if evmNetwork.DefaultGasLimit != 0 {
					conf.GasLimit = evmNetwork.DefaultGasLimit
				}

				sethNetwork = conf
				break
			}
		} else if conf.ChainID == fmt.Sprint(evmNetwork.ChainID) {
			conf.PrivateKeys = evmNetwork.PrivateKeys
			if len(conf.URLs) == 0 {
				conf.URLs = evmNetwork.URLs
			}

			sethNetwork = conf
			break
		}
	}

	if sethNetwork == nil {
		return seth.Config{}, fmt.Errorf("No matching EVM network found for chain ID %d. If it's a new network please define it as [Network.EVMNetworks.NETWORK_NAME] in TOML", evmNetwork.ChainID)
	}

	sethConfig.Network = sethNetwork

	return sethConfig, nil
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

// ValidateSethNetworkConfig validates the Seth network config
func ValidateSethNetworkConfig(cfg *seth.Network) error {
	if cfg == nil {
		return fmt.Errorf("Network cannot be nil")
	}
	if cfg.ChainID == "" {
		return fmt.Errorf("ChainID is required")
	}
	_, err := strconv.Atoi(cfg.ChainID)
	if err != nil {
		return fmt.Errorf("ChainID needs to be a number")
	}
	if cfg.URLs == nil || len(cfg.URLs) == 0 {
		return fmt.Errorf("URLs are required")
	}
	if cfg.PrivateKeys == nil || len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("PrivateKeys are required")
	}
	if cfg.TransferGasFee == 0 {
		return fmt.Errorf("TransferGasFee needs to be above 0. It's the gas fee for a simple transfer transaction")
	}
	if cfg.TxnTimeout.Duration() == 0 {
		return fmt.Errorf("TxnTimeout needs to be above 0. It's the timeout for a transaction")
	}
	if cfg.EIP1559DynamicFees {
		if cfg.GasFeeCap == 0 {
			return fmt.Errorf("GasFeeCap needs to be above 0. It's the maximum fee per gas for a transaction (including tip)")
		}
		if cfg.GasTipCap == 0 {
			return fmt.Errorf("GasTipCap needs to be above 0. It's the maximum tip per gas for a transaction")
		}
		if cfg.GasFeeCap <= cfg.GasTipCap {
			return fmt.Errorf("GasFeeCap needs to be above GasTipCap (as it is base fee + tip cap)")
		}
	} else {
		if cfg.GasPrice == 0 {
			return fmt.Errorf("GasPrice needs to be above 0. It's the price of gas for a transaction")
		}
	}

	return nil
}

// ValidateAddressesTypeAndNumber makes sure that minKeyNumber of epehemeral addresses are used on a simulated network and that no ephemeral addresses are used on a live network. Also, it makes sure that at least minKeyNumber of private keys are used on a live network.
func ValidateAddressesTypeAndNumber(sethCfg *seth.Config, minKeyNumber int) error {
	if sethCfg.IsSimulatedNetwork() {
		if int(*sethCfg.EphemeralAddrs) < minKeyNumber {
			return fmt.Errorf("You need to use at least %d ephemeral addresses on a simulated network, but %d was used. Please update '[Seth] ephemeral_addresses_number' field in the TOML config file", minKeyNumber, int(*sethCfg.EphemeralAddrs))
		}
		// take only the first key, all others are not funded in genesis and will crash the test ¯\_(ツ)_/¯
		sethCfg.Network.PrivateKeys = sethCfg.Network.PrivateKeys[0:1]
	} else {
		if int(*sethCfg.EphemeralAddrs) > 0 {
			msg := `You must not use any ephemeral addresses on a live network, but %d were set. All funds send to them will be lost. You should prefund some addresses instead manually (or using Seth CLI) and pass their private keys in network configuration or
			TOML keyfile. But if you really know what you are doing remove this check from the test.`
			return fmt.Errorf(msg, int(*sethCfg.EphemeralAddrs))
		}
		if len(sethCfg.Network.PrivateKeys) < minKeyNumber {
			msg := `You need to use at least %d addresses on a live network, but %d were passed. With current setting either concurrent deployment or load generation will surely fail. Please update '[Network.NETWORK_NAME] private_keys_secret=[]' field in the TOML config file
			or set SETH_KEYFILE_PATH env var with path to keyfile.toml generated by Seth key CLI.`
			return fmt.Errorf(msg, minKeyNumber, len(sethCfg.Network.PrivateKeys))
		}
	}

	return nil
}
