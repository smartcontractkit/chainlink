package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	seth_chain "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctf_config_types "github.com/smartcontractkit/chainlink-testing-framework/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	ccipconfig "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCapReg_InMemory_Concurrent(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     4,
		Nodes:      4,
	})
	testDeployCapRegWithEnv_Concurrent(t, lggr, e)
}

func TestDeployCapReg_NewDevnet_Concurrent(t *testing.T) {
	lggr := logger.TestLogger(t)

	geth := ctf_config_types.ExecutionLayer_Geth
	eth1 := ctf_config_types.EthereumVersion_Eth1

	defaultSethConfig := seth.NewClientBuilder().WithGasPriceEstimations(false, 0, seth.Priority_Standard).BuildConfig()
	dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
	require.NoError(t, err)

	firstNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	firstNetworkConfig.ChainID = 1337

	firstChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &firstNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	secondNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	secondNetworkConfig.ChainID = 2337

	secondChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &secondNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	envConfig := persistent.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			NewEVMChains: []persistent_types.NewEVMChainConfig{firstChain, secondChain},
		},
	}

	e, err := persistent.NewEVMEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")
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

type TestAwareDONHooks struct {
	*testing.T
}

func (s *TestAwareDONHooks) PreStartupHook(nodes []*test_env.ClNode) error {
	for _, node := range nodes {
		node.SetTestLogger(s.T)
	}

	return nil
}

func (s *TestAwareDONHooks) PostStartupHook([]*test_env.ClNode) error {
	return nil
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

	firstChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &firstNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	secondNetworkConfig := ctf_config.MustGetDefaultChainConfig()
	secondNetworkConfig.ChainID = 2337

	secondChain, err := seth_chain.CreateNewEVMChainWithSeth(ctf_config.EthereumNetworkConfig{
		ExecutionLayer:      &geth,
		EthereumVersion:     &eth1,
		EthereumChainConfig: &secondNetworkConfig,
		DockerNetworkNames:  []string{dockerNetwork.Name},
	}, *defaultSethConfig)
	require.NoError(t, err, "Error creating new EVM chain with Seth")

	envConfig := persistent.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			NewEVMChains: []persistent_types.NewEVMChainConfig{firstChain, secondChain},
		},
		DONConfig: persistent.DONConfig{
			NewDON: &persistent.NewDockerDONConfig{
				NewDONHooks: &TestAwareDONHooks{t},
				ChainlinkDeployment: testconfig.ChainlinkDeployment{
					Common: &testconfig.Node{
						ChainlinkImage: &ctfconfig.ChainlinkImageConfig{
							Image:   ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
							Version: ptr.Ptr("2.13.0"),
						},
						DBImage: "795953128386.dkr.ecr.us-west-2.amazonaws.com/postgres",
						DBTag:   "15.6",
						BaseConfigTOML: `
[Feature]
LogPoller = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
HTTPWriteTimeout = '1m'

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[Database]
MaxIdleConns = 10
MaxOpenConns = 20
MigrateOnStartup = true

[OCR2]
Enabled = true
DefaultTransactionQueueDepth = 0

[OCR]
Enabled = false
DefaultTransactionQueueDepth = 0

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']
DeltaDial = '500ms'
DeltaReconcile = '5s'
`},
					NoOfNodes: ptr.Ptr(5),
				},
				Options: persistent.Options{
					Networks: []string{dockerNetwork.Name},
				},
			},
		},
	}

	e, err := persistent.NewEVMEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")
	testDeployCCIPContractsWithEnv(t, lggr, *e)
}

func TestDeployCCIPContractsNewDevnet_FromTestConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	testCfg := ccipconfig.GlobalTestConfig()
	require.NoError(t, testCfg.Validate(), "Error validating test config")

	// here we are creating Seth config, but we should read it from the test config
	defaultSethConfig := seth.NewClientBuilder().BuildConfig()

	chainCfg, err := persistent.EVMChainConfigFromTestConfig(*testCfg, defaultSethConfig)
	require.NoError(t, err, "Error creating chain config from test config")

	envConfig := persistent.EnvironmentConfig{
		ChainConfig: chainCfg,
	}

	e, err := persistent.NewEVMEnvironment(lggr, envConfig)
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

	envConfig := persistent.EnvironmentConfig{
		ChainConfig: persistent_types.ChainConfig{
			ExistingEVMChains: []persistent_types.ExistingEVMChainConfig{firstChain, secondChain},
		},
	}
	e, err := persistent.NewEVMEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")
	testDeployCCIPContractsWithEnv(t, lggr, *e)
}

func testDeployCCIPContractsWithEnv(t *testing.T, lggr logger.Logger, e deployment.Environment[deployment.Chain]) {
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
	}

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

func testDeployCapRegWithEnv_Concurrent(t *testing.T, lggr logger.Logger, e deployment.Environment[deployment.Chain]) {
	var ab deployment.AddressBook
	// Deploy all the CCIP contracts.
	ab, _, err := DeployCapReg_Concurrent(lggr, e.Chains, e.AllChainSelectors())
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
