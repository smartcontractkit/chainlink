package benchmark

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	env_client "github.com/smartcontractkit/chainlink-testing-framework/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/integration-tests/types"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	keeperBenchmarkBaseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]
[Keeper]
TurnLookBack = 0
[WebServer]
HTTPWriteTimeout = '1h'`

	simulatedEVMNonDevTOML = `
Enabled = true
FinalityDepth = 50
LogPollInterval = '1s'

[EVM.HeadTracker]
HistoryDepth = 100

[EVM.GasEstimator]
Mode = 'FixedPrice'
LimitDefault = 5_000_000`

	performanceChainlinkResources = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "4Gi",
			},
		},
	}
	performanceDbResources = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "1Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "1Gi",
			},
		},
		"stateful": true,
		"capacity": "10Gi",
	}

	soakChainlinkResources = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "350m",
				"memory": "1Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "350m",
				"memory": "1Gi",
			},
		},
	}
	soakDbResources = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "250m",
				"memory": "256Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "250m",
				"memory": "256Mi",
			},
		},
		"stateful": true,
		"capacity": "10Gi",
	}
)

type NetworkConfig struct {
	upkeepSLA  int64
	blockTime  time.Duration
	deltaStage time.Duration
	funding    *big.Float
}

func TestAutomationBenchmark(t *testing.T) {
	l := logging.GetTestLogger(t)
	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err, "Error getting test type")

	config, err := tc.GetConfig(testType, tc.Keeper)
	require.NoError(t, err, "Error getting test config")

	testEnvironment, benchmarkNetwork := SetupAutomationBenchmarkEnv(t, &config)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}
	networkName := strings.ReplaceAll(benchmarkNetwork.Name, " ", "")
	testName := fmt.Sprintf("%s%s", networkName, *config.Keeper.Common.RegistryToTest)
	l.Info().Str("Test Name", testName).Msg("Running Benchmark Test")
	benchmarkTestNetwork := getNetworkConfig(networkName, &config)

	l.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")

	chainClient, err := blockchain.NewEVMClient(benchmarkNetwork, testEnvironment, l)
	require.NoError(t, err, "Error connecting to blockchain")
	registryVersions := addRegistry(&config)
	keeperBenchmarkTest := testsetups.NewKeeperBenchmarkTest(t,
		testsetups.KeeperBenchmarkTestInputs{
			BlockchainClient: chainClient,
			RegistryVersions: registryVersions,
			KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
				PaymentPremiumPPB:    uint32(0),
				FlatFeeMicroLINK:     uint32(40000),
				BlockCountPerTurn:    big.NewInt(100),
				CheckGasLimit:        uint32(45_000_000), //45M
				StalenessSeconds:     big.NewInt(90_000),
				GasCeilingMultiplier: uint16(2),
				MaxPerformGas:        uint32(*config.Keeper.Common.MaxPerformGas),
				MinUpkeepSpend:       big.NewInt(0),
				FallbackGasPrice:     big.NewInt(2e11),
				FallbackLinkPrice:    big.NewInt(2e18),
				MaxCheckDataSize:     uint32(5_000),
				MaxPerformDataSize:   uint32(5_000),
			},
			Upkeeps: &testsetups.UpkeepConfig{
				NumberOfUpkeeps:     *config.Keeper.Common.NumberOfUpkeeps,
				CheckGasToBurn:      *config.Keeper.Common.CheckGasToBurn,
				PerformGasToBurn:    *config.Keeper.Common.PerformGasToBurn,
				BlockRange:          *config.Keeper.Common.BlockRange,
				BlockInterval:       *config.Keeper.Common.BlockInterval,
				UpkeepGasLimit:      *config.Keeper.Common.UpkeepGasLimit,
				FirstEligibleBuffer: 1,
			},
			Contracts: &testsetups.PreDeployedContracts{
				RegistrarAddress: *config.Keeper.Common.RegistrarAddress,
				RegistryAddress:  *config.Keeper.Common.RegistryAddress,
				LinkTokenAddress: *config.Keeper.Common.LinkTokenAddress,
				EthFeedAddress:   *config.Keeper.Common.EthFeedAddress,
				GasFeedAddress:   *config.Keeper.Common.GasFeedAddress,
			},
			ChainlinkNodeFunding: benchmarkTestNetwork.funding,
			UpkeepSLA:            benchmarkTestNetwork.upkeepSLA,
			BlockTime:            benchmarkTestNetwork.blockTime,
			DeltaStage:           benchmarkTestNetwork.deltaStage,
			ForceSingleTxnKey:    *config.Keeper.Common.ForceSingleTxKey,
			DeleteJobsOnEnd:      *config.Keeper.Common.DeleteJobsOnEnd,
		},
	)
	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})
	keeperBenchmarkTest.Setup(testEnvironment, &config)
	keeperBenchmarkTest.Run()
}

func addRegistry(config *tc.TestConfig) []eth_contracts.KeeperRegistryVersion {
	switch *config.Keeper.Common.RegistryToTest {
	case "1_1":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_1}
	case "1_2":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2}
	case "1_3":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3}
	case "2_0":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0}
	case "2_1":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_1}
	case "2_0-1_3":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0, eth_contracts.RegistryVersion_1_3}
	case "2_1-2_0-1_3":
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_1,
			eth_contracts.RegistryVersion_2_0, eth_contracts.RegistryVersion_1_3}
	case "2_0-Multiple":
		return repeatRegistries(eth_contracts.RegistryVersion_2_0, *config.Keeper.Common.NumberOfRegistries)
	case "2_1-Multiple":
		return repeatRegistries(eth_contracts.RegistryVersion_2_1, *config.Keeper.Common.NumberOfRegistries)
	default:
		return []eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0}
	}
}

func repeatRegistries(registryVersion eth_contracts.KeeperRegistryVersion, numberOfRegistries int) []eth_contracts.KeeperRegistryVersion {
	repeatedRegistries := make([]eth_contracts.KeeperRegistryVersion, 0)
	for i := 0; i < numberOfRegistries; i++ {
		repeatedRegistries = append(repeatedRegistries, registryVersion)
	}
	return repeatedRegistries
}

func getNetworkConfig(networkName string, config *tc.TestConfig) NetworkConfig {
	var nc NetworkConfig
	var ok bool
	if nc, ok = networkConfig[networkName]; !ok {
		return NetworkConfig{}
	}

	if networkName == "SimulatedGeth" || networkName == "geth" {
		return nc
	}

	nc.funding = big.NewFloat(*config.Common.ChainlinkNodeFunding)

	return nc
}

var networkConfig = map[string]NetworkConfig{
	"SimulatedGeth": {
		upkeepSLA:  int64(120), //2 minutes
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	"geth": {
		upkeepSLA:  int64(120), //2 minutes
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	"GoerliTestnet": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	"ArbitrumGoerli": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: time.Duration(0),
	},
	"OptimismGoerli": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: time.Duration(0),
	},
	"SepoliaTestnet": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	"PolygonMumbai": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	"BaseGoerli": {
		upkeepSLA:  int64(60),
		blockTime:  2 * time.Second,
		deltaStage: 20 * time.Second,
	},
	"ArbitrumSepolia": {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
	"LineaGoerli": {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
}

func SetupAutomationBenchmarkEnv(t *testing.T, keeperTestConfig types.KeeperBenchmarkTestConfig) (*environment.Environment, blockchain.EVMNetwork) {
	l := logging.GetTestLogger(t)
	testNetwork := networks.MustGetSelectedNetworkConfig(keeperTestConfig.GetNetworkConfig())[0] // Environment currently being used to run benchmark test on
	blockTime := "1"
	networkDetailTOML := `MinIncomingConfirmations = 1`

	numberOfNodes := *keeperTestConfig.GetKeeperConfig().Common.NumberOfNodes

	if strings.Contains(*keeperTestConfig.GetKeeperConfig().Common.RegistryToTest, "2_") {
		numberOfNodes++
	}

	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s-%s",
			strings.ToLower(keeperTestConfig.GetConfigurationName()),
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
			strings.ReplaceAll(strings.ToLower(*keeperTestConfig.GetKeeperConfig().Common.RegistryToTest), "_", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})

	dbResources := performanceDbResources
	chainlinkResources := performanceChainlinkResources
	if strings.ToLower(keeperTestConfig.GetConfigurationName()) == "soak" {
		chainlinkResources = soakChainlinkResources
		dbResources = soakDbResources
	}

	// Test can run on simulated, simulated-non-dev, testnets
	if testNetwork.Name == networks.SimulatedEVMNonDev.Name {
		networkDetailTOML = simulatedEVMNonDevTOML
		testEnvironment.
			AddHelm(reorg.New(&reorg.Props{
				NetworkName: testNetwork.Name,
				Values: map[string]interface{}{
					"geth": map[string]interface{}{
						"tx": map[string]interface{}{
							"replicas": numberOfNodes,
						},
						"miner": map[string]interface{}{
							"replicas": 2,
						},
					},
				},
			}))
	} else {
		testEnvironment.
			AddHelm(ethereum.New(&ethereum.Props{
				NetworkName: testNetwork.Name,
				Simulated:   testNetwork.Simulated,
				WsURLs:      testNetwork.URLs,
				Values: map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"cpu":    "4000m",
							"memory": "4Gi",
						},
						"limits": map[string]interface{}{
							"cpu":    "4000m",
							"memory": "4Gi",
						},
					},
					"geth": map[string]interface{}{
						"blocktime": blockTime,
					},
				},
			}))
	}

	// deploy blockscout if running on simulated
	if testNetwork.Simulated {
		testEnvironment.
			AddChart(blockscout.New(&blockscout.Props{
				Name:    "geth-blockscout",
				WsURL:   testNetwork.URLs[0],
				HttpURL: testNetwork.HTTPURLs[0]}))
	}
	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	if testEnvironment.WillUseRemoteRunner() {
		return testEnvironment, testNetwork
	}

	// separate RPC urls per CL node
	internalWsURLs := make([]string, 0)
	internalHttpURLs := make([]string, 0)
	for i := 0; i < numberOfNodes; i++ {
		// for simulated-nod-dev each CL node gets its own RPC node
		if testNetwork.Name == networks.SimulatedEVMNonDev.Name {
			podName := fmt.Sprintf("%s-ethereum-geth:%d", testNetwork.Name, i)
			txNodeInternalWs, err := testEnvironment.Fwd.FindPort(podName, "geth", "ws-rpc").As(env_client.RemoteConnection, env_client.WS)
			require.NoError(t, err, "Error finding WS ports")
			internalWsURLs = append(internalWsURLs, txNodeInternalWs)
			txNodeInternalHttp, err := testEnvironment.Fwd.FindPort(podName, "geth", "http-rpc").As(env_client.RemoteConnection, env_client.HTTP)
			require.NoError(t, err, "Error finding HTTP ports")
			internalHttpURLs = append(internalHttpURLs, txNodeInternalHttp)
			// for testnets with more than 1 RPC nodes
		} else if len(testNetwork.URLs) > 1 {
			internalWsURLs = append(internalWsURLs, testNetwork.URLs[i%len(testNetwork.URLs)])
			internalHttpURLs = append(internalHttpURLs, testNetwork.HTTPURLs[i%len(testNetwork.URLs)])
			// for simulated and testnets with 1 RPC node
		} else {
			internalWsURLs = append(internalWsURLs, testNetwork.URLs[0])
			internalHttpURLs = append(internalHttpURLs, testNetwork.HTTPURLs[0])
		}
	}
	l.Debug().Strs("internalWsURLs", internalWsURLs).Strs("internalHttpURLs", internalHttpURLs).Msg("internalURLs")

	for i := 0; i < numberOfNodes; i++ {
		testNetwork.HTTPURLs = []string{internalHttpURLs[i]}
		testNetwork.URLs = []string{internalWsURLs[i]}

		var overrideFn = func(_ interface{}, target interface{}) {
			ctf_config.MustConfigOverrideChainlinkVersion(keeperTestConfig.GetChainlinkImageConfig(), target)
			ctf_config.MightConfigOverridePyroscopeKey(keeperTestConfig.GetPyroscopeConfig(), target)
		}

		cd := chainlink.NewWithOverride(i, map[string]any{
			"toml":      networks.AddNetworkDetailedConfig(keeperBenchmarkBaseTOML, keeperTestConfig.GetPyroscopeConfig(), networkDetailTOML, testNetwork),
			"chainlink": chainlinkResources,
			"db":        dbResources,
		}, keeperTestConfig.GetChainlinkImageConfig(), overrideFn)

		testEnvironment.AddHelm(cd)
	}
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, testNetwork
}
