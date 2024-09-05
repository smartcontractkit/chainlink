package ccipdeployment

import (
	"encoding/json"
	"fmt"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/internal/testutil"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/hooks"
	seth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/seth"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TODO this will not work now, because we only have 1 key for that environment, is there a way to use multiple?
// TODO if not, it will be hard to have tests that require concurrent deployments and can run on persistent and in-memory environments
func TestDeployCapReg_InMemory_Concurrent(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	testDeployCapRegWithEnv_Concurrent(t, lggr, e)
}

func TestDeployCapReg_NewDevnet_Concurrent(t *testing.T) {
	lggr := logger.TestLogger(t)

	geth := ctf_config_types.ExecutionLayer_Geth
	eth1 := ctf_config_types.EthereumVersion_Eth1

	defaultSethConfig := seth.NewClientBuilder().WithGasPriceEstimations(false, 0, seth.Priority_Standard).BuildConfig()
	two := int64(2)
	defaultSethConfig.EphemeralAddrs = &two
	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	require.NoError(t, err)

	firstNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	firstNetworkConfig.ChainID = 1337

	ls, err := logstream.NewLogStream(t, &ctf_config.LoggingConfig{
		TestLogCollect: ptr.Ptr(true),
		LogStream: &ctf_config.LogStreamConfig{
			LogTargets:            []string{fmt.Sprint(logstream.File)},
			LogProducerTimeout:    ptr.Ptr(blockchain.StrDuration{Duration: 1 * time.Minute}),
			LogProducerRetryLimit: ptr.Ptr(uint(5)),
		},
	})
	require.NoError(t, err, "Error creating new log stream")

	defaultNewEVMHooks := hooks.DefaultPrivateEVMHooks{
		T:         t,
		LogStream: ls,
	}

	firstChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &firstNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig, &defaultNewEVMHooks)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	secondNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	secondNetworkConfig.ChainID = 2337

	secondChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &secondNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig, &defaultNewEVMHooks)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	envConfig := persistent_types.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			NewEVMChains: []persistent_types.NewEVMChainProducer{firstChain, secondChain},
		},
		DONConfig: persistent_types.DONConfig{
			NewDON: &persistent_types.NewDockerDONConfig{
				ChainlinkDeployment: testutil.GetDefaultNewChainlinkClusterConfig(),
				Options: persistent_types.DockerOptions{
					Networks: []string{dockerNetwork.Name},
				},
			},
		},
	}

	e, err := persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new environment")
	testDeployCapRegWithEnv_Concurrent(t, lggr, *e)
}

func TestDeployCCIPContractsInMemory(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      4,
	})
	testDeployCCIPContractsWithEnv(t, lggr, e)
}

func TestDeployCCIPContractsNewDevnet(t *testing.T) {
	lggr := logger.TestLogger(t)

	geth := ctf_config_types.ExecutionLayer_Geth
	eth1 := ctf_config_types.EthereumVersion_Eth1

	defaultSethConfig := seth.NewClientBuilder().WithGasPriceEstimations(false, 0, "").BuildConfig()
	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	require.NoError(t, err)

	firstNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	firstNetworkConfig.ChainID = 1337

	ls, err := logstream.NewLogStream(t, &ctf_config.LoggingConfig{
		TestLogCollect: ptr.Ptr(true),
		LogStream: &ctf_config.LogStreamConfig{
			LogTargets:            []string{fmt.Sprint(logstream.File)},
			LogProducerTimeout:    ptr.Ptr(blockchain.StrDuration{Duration: 1 * time.Minute}),
			LogProducerRetryLimit: ptr.Ptr(uint(5)),
		},
	})
	require.NoError(t, err, "Error creating new log stream")
	defaultNewEVMHooks := hooks.DefaultPrivateEVMHooks{
		T:         t,
		LogStream: ls,
	}

	firstChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &firstNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig, &defaultNewEVMHooks)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	secondNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	secondNetworkConfig.ChainID = 2337

	secondChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &secondNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig, &defaultNewEVMHooks)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	envConfig := persistent_types.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			NewEVMChains: []persistent_types.NewEVMChainProducer{firstChain, secondChain},
		},
		DONConfig: persistent_types.DONConfig{
			NewDON: &persistent_types.NewDockerDONConfig{
				ChainlinkDeployment: testutil.GetDefaultNewChainlinkClusterConfig(),
				Options: persistent_types.DockerOptions{
					Networks: []string{dockerNetwork.Name},
				},
			},
		},
	}

	e, err := persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")
	testDeployCCIPContractsWithEnv(t, lggr, *e)
}

// TODO: update urls before running
func TestDeployCCIPContractsExistingDevnet(t *testing.T) {
	lggr := logger.TestLogger(t)
	defaultSethConfig := seth.NewClientBuilder().BuildConfig()

	firstChain, err := seth_chain.CreateExistingEVMChainWithSeth(
		blockchain.EVMNetwork{
			Name:        "SomeChain_1337",
			ChainID:     1337,
			URLs:        []string{"ws://127.0.0.1:57163"},
			HTTPURLs:    []string{"ws://127.0.0.1:57162"},
			PrivateKeys: []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}, // default Geth PK
		},
		*defaultSethConfig,
	)
	require.NoError(t, err, "Error creating existing EVM chain with Seth")

	secondChain, err := seth_chain.CreateExistingEVMChainWithSeth(
		blockchain.EVMNetwork{
			Name:        "SomeChain_2337",
			ChainID:     2337,
			URLs:        []string{"ws://127.0.0.1:57251"},
			HTTPURLs:    []string{"ws://127.0.0.1:57161"},
			PrivateKeys: []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}, // default Geth PK
		},
		*defaultSethConfig,
	)
	require.NoError(t, err, "Error creating existing EVM chain with Seth")

	envConfig := persistent_types.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			ExistingEVMChains: []persistent_types.ExistingEVMChainProducer{firstChain, secondChain},
		},
	}
	e, err := persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")
	testDeployCCIPContractsWithEnv(t, lggr, *e)
}

func testDeployCCIPContractsWithEnv(t *testing.T, lggr logger.Logger, e deployment.Environment) {
	var ab deployment.AddressBook
	// Deploy all the CCIP contracts.
	for _, chain := range e.AllChainSelectors() {
		capRegAddresses, _, err := DeployCapReg(lggr, e.Chains, chain)
		require.NoError(t, err)
		s, err := LoadOnchainState(e, capRegAddresses)
		require.NoError(t, err)
		newAb, err := DeployCCIPContracts(e, DeployCCIPContractConfig{
			HomeChainSel:     chain,
			CCIPOnChainState: s,
		})
		require.NoError(t, err)
		if ab == nil {
			ab = newAb
		} else {
			mergeErr := ab.Merge(newAb)
			require.NoError(t, mergeErr)
		}
		_ = s
	}

	_ = ab

	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.Snapshot(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func testDeployCapRegWithEnv_Concurrent(t *testing.T, lggr logger.Logger, e deployment.Environment) {
	var ab deployment.AddressBook
	// Deploy all the CCIP contracts.
	ab, _, err := DeployCapReg_Concurrent(lggr, e.Chains, 2)
	require.NoError(t, err)

	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.Snapshot(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TestJobSpecGeneration(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	js, err := NewCCIPJobSpecs(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for node, jb := range js {
		fmt.Println(node, jb)
	}
	// TODO: Add job assertions
}
