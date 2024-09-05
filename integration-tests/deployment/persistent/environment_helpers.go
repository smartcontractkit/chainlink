package persistent

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/internal/testutil"
	"testing"

	geth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/geth"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/hooks"
	seth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/seth"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	ccip_test_config "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	chainlink_test_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

// TODO in the future Seth config should be part of the test config
func EnvironmentConfigFromCCIPTestConfig(t *testing.T, testCfg ccip_test_config.Config, useSeth bool) (persistent_types.EnvironmentConfig, error) {
	envConfig := persistent_types.EnvironmentConfig{}
	evmChainConfig := persistent_types.ChainConfig{
		NewEVMChains:      make([]persistent_types.NewEVMChainProducer, 0),
		ExistingEVMChains: make([]persistent_types.ExistingEVMChainProducer, 0),
	}

	var getSimulatedNetworkFromTestConfig = func(testConfig ccip_test_config.Config, chainId uint64) (ctf_config.EthereumNetworkConfig, error) {
		for _, chainCfg := range testConfig.CCIP.Env.PrivateEthereumNetworks {
			if uint64(chainCfg.EthereumChainConfig.ChainID) == chainId {
				return *chainCfg, nil
			}
		}

		return ctf_config.EthereumNetworkConfig{}, fmt.Errorf("chain id %d not found in test config", chainId)
	}

	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	if err != nil {
		return envConfig, err
	}

	ls, err := logstream.NewLogStream(t, testCfg.CCIP.Env.Logging)
	if err != nil {
		return envConfig, err
	}

	evmHooks := hooks.DefaultPrivateEVMHooks{
		T:         t,
		LogStream: ls,
	}

	// here we are creating Seth config, but we should read it from the test config
	defaultSethConfig := seth.NewClientBuilder().BuildConfig()

	for _, network := range networks.MustGetSelectedNetworkConfig(testCfg.CCIP.Env.Network) {
		if network.Simulated {
			privateNetworkCfg, err := getSimulatedNetworkFromTestConfig(testCfg, uint64(network.ChainID))
			if err != nil {
				return envConfig, err
			}
			privateNetworkCfg.DockerNetworkNames = []string{dockerNetwork.Name}
			var chainConfig persistent_types.NewEVMChainProducer
			if !useSeth {
				chainConfig = geth_chain.CreateNewEVMChainWithGeth(privateNetworkCfg, &evmHooks)
			} else {
				chainConfig, err = seth_chain.CreateNewEVMChainWithSeth(privateNetworkCfg, *defaultSethConfig, &evmHooks)
				if err != nil {
					return envConfig, err
				}
			}
			evmChainConfig.NewEVMChains = append(evmChainConfig.NewEVMChains, chainConfig)
		} else {
			var chainConfig persistent_types.ExistingEVMChainProducer
			if !useSeth {
				chainConfig = geth_chain.CreateExistingEVMChainConfigWithGeth(network)
			} else {
				chainConfig, err = seth_chain.CreateExistingEVMChainWithSeth(network, *defaultSethConfig)
				if err != nil {
					return envConfig, err
				}
			}
			evmChainConfig.ExistingEVMChains = append(evmChainConfig.ExistingEVMChains, chainConfig)
		}
	}

	envConfig.ChainConfig = evmChainConfig

	donConfig := persistent_types.DONConfig{}
	if testCfg.CCIP.Env.NewCLCluster != nil {
		donConfig.NewDON = &persistent_types.NewDockerDONConfig{
			NewDONHooks:         hooks.NewDefaultDONHooksFromTestConfig(t, logging.GetTestLogger(t), ls, nil, testCfg.CCIP.Env.Logging),
			ChainlinkDeployment: testCfg.CCIP.Env.NewCLCluster,
			Options: persistent_types.DockerOptions{
				Networks: []string{dockerNetwork.Name},
			},
		}
	} else if testCfg.CCIP.Env.ExistingCLCluster != nil && testCfg.CCIP.Env.Mockserver != nil {
		donConfig.ExistingDON = &persistent_types.ExistingDONConfig{
			CLCluster:     testCfg.CCIP.Env.ExistingCLCluster,
			MockServerURL: testCfg.CCIP.Env.Mockserver,
		}
	} else {
		return envConfig, fmt.Errorf("either new or existing chainlink cluster config must be provided")
	}
	envConfig.DONConfig = donConfig
	envConfig.EnvironmentHooks = hooks.NewDefaultEnvironmentHooksFromTestConfig(t, logging.GetTestLogger(t), defaultSethConfig, testCfg.CCIP.Env.Logging)

	return envConfig, nil
}

// TODO in the future newClCluster and existingCluster (& mockServerUrl) should be part of the test config
func EnvironmentConfigFromChainlinkTestConfig(t *testing.T, testCfg chainlink_test_config.TestConfig, startNewCluster bool, existingCluster *ccip_test_config.CLCluster, mockServerUrl *string) (persistent_types.EnvironmentConfig, error) {
	envConfig := persistent_types.EnvironmentConfig{}
	evmChainConfig := persistent_types.ChainConfig{
		NewEVMChains:      make([]persistent_types.NewEVMChainProducer, 0),
		ExistingEVMChains: make([]persistent_types.ExistingEVMChainProducer, 0),
	}

	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	if err != nil {
		return envConfig, err
	}

	ls, err := logstream.NewLogStream(t, testCfg.Logging)
	if err != nil {
		return envConfig, err
	}

	evmHooks := hooks.DefaultPrivateEVMHooks{
		T:         t,
		LogStream: ls,
	}

	for _, network := range networks.MustGetSelectedNetworkConfig(testCfg.Network) {
		if network.Simulated {
			// we do not support more than 1 simulated network in chainlink tests
			privateNetworkCfg := testCfg.PrivateEthereumNetwork
			if privateNetworkCfg == nil {
				return envConfig, fmt.Errorf("private ethereum network config must be provided for simulated network")
			}
			privateNetworkCfg.DockerNetworkNames = []string{dockerNetwork.Name}
			chainConfig, err := seth_chain.CreateNewEVMChainWithSeth(*privateNetworkCfg, *testCfg.GetSethConfig(), &evmHooks)
			if err != nil {
				return envConfig, err
			}

			evmChainConfig.NewEVMChains = append(evmChainConfig.NewEVMChains, chainConfig)
		} else {
			chainConfig, err := seth_chain.CreateExistingEVMChainWithSeth(network, *testCfg.GetSethConfig())
			if err != nil {
				return envConfig, err
			}
			evmChainConfig.ExistingEVMChains = append(evmChainConfig.ExistingEVMChains, chainConfig)
		}
	}

	envConfig.ChainConfig = evmChainConfig

	donConfig := persistent_types.DONConfig{}
	if startNewCluster {
		donConfig.NewDON = &persistent_types.NewDockerDONConfig{
			NewDONHooks:         hooks.NewDefaultDONHooksFromTestConfig(t, logging.GetTestLogger(t), ls, nil, testCfg.Logging),
			ChainlinkDeployment: testutil.GetDefaultNewClusterConfigFromChainlinkTestConfig(testCfg),
			Options: persistent_types.DockerOptions{
				Networks: []string{dockerNetwork.Name},
			},
		}
	} else if existingCluster != nil && mockServerUrl != nil {
		donConfig.ExistingDON = &persistent_types.ExistingDONConfig{
			CLCluster:     existingCluster,
			MockServerURL: mockServerUrl,
		}
	} else {
		return envConfig, fmt.Errorf("either new or existing chainlink cluster (with mockserver) config must be provided")
	}

	envConfig.DONConfig = donConfig
	envConfig.EnvironmentHooks = hooks.NewDefaultEnvironmentHooksFromTestConfig(t, logging.GetTestLogger(t), testCfg.GetSethConfig(), testCfg.Logging)

	return envConfig, nil
}
