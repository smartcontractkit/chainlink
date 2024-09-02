package persistent

import (
	"fmt"
	geth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/geth"
	seth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/seth"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	ccipconfig "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
)

// TODO in the future Seth config should be part of the test config
func EVMChainConfigFromTestConfig(testCfg ccipconfig.Config, sethConfig *seth.Config) (persistent_types.ChainConfig, error) {
	evmChainConfig := persistent_types.ChainConfig{
		NewEVMChains:      make([]persistent_types.NewEVMChainConfig, 0),
		ExistingEVMChains: make([]persistent_types.ExistingEVMChainConfig, 0),
	}

	var getSimulatedNetworkFromTestConfig = func(testConfig ccipconfig.Config, chainId uint64) (ctf_config.EthereumNetworkConfig, error) {
		for _, chainCfg := range testConfig.CCIP.Env.PrivateEthereumNetworks {
			if uint64(chainCfg.EthereumChainConfig.ChainID) == chainId {
				return *chainCfg, nil
			}
		}

		return ctf_config.EthereumNetworkConfig{}, fmt.Errorf("chain id %d not found in test config", chainId)
	}

	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	if err != nil {
		return evmChainConfig, err
	}

	for _, network := range networks.MustGetSelectedNetworkConfig(testCfg.CCIP.Env.Network) {
		if network.Simulated {
			privateNetworkCfg, err := getSimulatedNetworkFromTestConfig(testCfg, uint64(network.ChainID))
			if err != nil {
				return evmChainConfig, err
			}
			privateNetworkCfg.DockerNetworkNames = []string{dockerNetwork.Name}
			var chainConfig persistent_types.NewEVMChainConfig
			if sethConfig == nil {
				chainConfig = geth_chain.CreateNewEVMChainWithGeth(&privateNetworkCfg)
			} else {
				chainConfig, err = seth_chain.CreateNewEVMChainWithSeth(privateNetworkCfg, *sethConfig)
				if err != nil {
					return evmChainConfig, err
				}
			}
			evmChainConfig.NewEVMChains = append(evmChainConfig.NewEVMChains, chainConfig)
		} else {
			var chainConfig persistent_types.ExistingEVMChainConfig
			if sethConfig == nil {
				chainConfig = geth_chain.CreateExistingEVMChainConfigWithGeth(network)
			} else {
				chainConfig, err = seth_chain.CreateExistingEVMChainWithSeth(network, *sethConfig)
				if err != nil {
					return evmChainConfig, err
				}
			}
			evmChainConfig.ExistingEVMChains = append(evmChainConfig.ExistingEVMChains, chainConfig)
		}
	}

	return evmChainConfig, nil
}
