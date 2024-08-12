package benchmark

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
	envclient "github.com/smartcontractkit/chainlink-testing-framework/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	sethutils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	ethcontracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

var (
	chainlinkResources = map[string]interface{}{
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
	dbResources = map[string]interface{}{
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
)

type NetworkConfig struct {
	upkeepSLA  int64
	blockTime  time.Duration
	deltaStage time.Duration
	funding    *big.Float
}

var defaultNetworkConfig = NetworkConfig{
	upkeepSLA:  int64(120),
	blockTime:  time.Second,
	deltaStage: time.Duration(0),
}

func TestAutomationBenchmark(t *testing.T) {
	l := logging.GetTestLogger(t)
	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err, "Error getting test type")

	config, err := tc.GetConfig([]string{testType}, tc.Keeper)
	require.NoError(t, err, "Error getting test config")

	testEnvironment, benchmarkNetwork := SetupAutomationBenchmarkEnv(t, &config)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}
	networkName := strings.ReplaceAll(benchmarkNetwork.Name, " ", "")
	testName := fmt.Sprintf("%s%s", networkName, *config.Keeper.Common.RegistryToTest)
	l.Info().Str("Test Name", testName).Msg("Running Benchmark Test")
	benchmarkTestNetwork := getNetworkConfig(&config)

	l.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")
	testNetwork := sethutils.MustReplaceSimulatedNetworkUrlWithK8(l, benchmarkNetwork, *testEnvironment)

	chainClient, err := sethutils.GetChainClientWithConfigFunction(&config, testNetwork, sethutils.OneEphemeralKeysLiveTestnetAutoFixFn)
	require.NoError(t, err, "Error getting Seth client")

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
				MaxRevertDataSize:    uint32(5_000),
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

func addRegistry(config *tc.TestConfig) []ethcontracts.KeeperRegistryVersion {
	switch *config.Keeper.Common.RegistryToTest {
	case "1_1":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_1_1}
	case "1_2":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_1_2}
	case "1_3":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_1_3}
	case "2_0":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_0}
	case "2_1":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_1}
	case "2_2":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_2}
	case "2_3":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_3}
	case "2_0-1_3":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_0, ethcontracts.RegistryVersion_1_3}
	case "2_1-2_0-1_3":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_1,
			ethcontracts.RegistryVersion_2_0, ethcontracts.RegistryVersion_1_3}
	case "2_2-2_1":
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_2, ethcontracts.RegistryVersion_2_1}
	case "2_0-Multiple":
		return repeatRegistries(ethcontracts.RegistryVersion_2_0, *config.Keeper.Common.NumberOfRegistries)
	case "2_1-Multiple":
		return repeatRegistries(ethcontracts.RegistryVersion_2_1, *config.Keeper.Common.NumberOfRegistries)
	case "2_2-Multiple":
		return repeatRegistries(ethcontracts.RegistryVersion_2_2, *config.Keeper.Common.NumberOfRegistries)
	default:
		return []ethcontracts.KeeperRegistryVersion{ethcontracts.RegistryVersion_2_0}
	}
}

func repeatRegistries(registryVersion ethcontracts.KeeperRegistryVersion, numberOfRegistries int) []ethcontracts.KeeperRegistryVersion {
	repeatedRegistries := make([]ethcontracts.KeeperRegistryVersion, 0)
	for i := 0; i < numberOfRegistries; i++ {
		repeatedRegistries = append(repeatedRegistries, registryVersion)
	}
	return repeatedRegistries
}

func getNetworkConfig(config *tc.TestConfig) NetworkConfig {
	evmNetwork := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	var nc NetworkConfig
	var ok bool
	if nc, ok = networkConfig[evmNetwork.Name]; !ok {
		nc = defaultNetworkConfig
	}

	if evmNetwork.Name == networks.SimulatedEVM.Name || evmNetwork.Name == networks.SimulatedEVMNonDev.Name {
		return nc
	}

	nc.funding = big.NewFloat(*config.Common.ChainlinkNodeFunding)

	return nc
}

var networkConfig = map[string]NetworkConfig{
	networks.SimulatedEVM.Name: {
		upkeepSLA:  int64(120), //2 minutes
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	networks.SimulatedEVMNonDev.Name: {
		upkeepSLA:  int64(120), //2 minutes
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	networks.GoerliTestnet.Name: {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	networks.SepoliaTestnet.Name: {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	networks.PolygonMumbai.Name: {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
	},
	networks.BaseSepolia.Name: {
		upkeepSLA:  int64(60),
		blockTime:  2 * time.Second,
		deltaStage: 20 * time.Second,
	},
	networks.ArbitrumSepolia.Name: {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
	networks.OptimismSepolia.Name: {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
	networks.LineaGoerli.Name: {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
	networks.GnosisChiado.Name: {
		upkeepSLA:  int64(120),
		blockTime:  6 * time.Second,
		deltaStage: 20 * time.Second,
	},
	networks.PolygonZkEvmCardona.Name: {
		upkeepSLA:  int64(120),
		blockTime:  time.Second,
		deltaStage: 20 * time.Second,
	},
}

func SetupAutomationBenchmarkEnv(t *testing.T, keeperTestConfig types.KeeperBenchmarkTestConfig) (*environment.Environment, blockchain.EVMNetwork) {
	l := logging.GetTestLogger(t)
	testNetwork := networks.MustGetSelectedNetworkConfig(keeperTestConfig.GetNetworkConfig())[0] // Environment currently being used to run benchmark test on
	blockTime := "1"
	numberOfNodes := *keeperTestConfig.GetKeeperConfig().Common.NumberOfNodes

	if strings.Contains(*keeperTestConfig.GetKeeperConfig().Common.RegistryToTest, "2_") {
		numberOfNodes++
	}

	networkName := strings.ReplaceAll(testNetwork.Name, " ", "-")
	networkName = strings.ReplaceAll(networkName, "_", "-")
	testNetwork.Name = networkName

	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s-%s",
			strings.ToLower(strings.Join(keeperTestConfig.GetConfigurationNames(), "")),
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
			strings.ReplaceAll(strings.ToLower(*keeperTestConfig.GetKeeperConfig().Common.RegistryToTest), "_", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})

	dbResources := dbResources
	chainlinkResources := chainlinkResources

	// Test can run on simulated, simulated-non-dev, testnets
	if testNetwork.Name == networks.SimulatedEVMNonDev.Name {
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
						"blocktime":      blockTime,
						"capacity":       "20Gi",
						"startGaslimit":  "20000000",
						"targetGasLimit": "30000000",
					},
				},
			}))
	}

	// TODO we need to update the image in CTF, the old one is not available anymore
	// deploy blockscout if running on simulated
	// if testNetwork.Simulated {
	// 	testEnvironment.
	// 		AddChart(blockscout.New(&blockscout.Props{
	// 			Name:    "geth-blockscout",
	// 			WsURL:   testNetwork.URLs[0],
	// 			HttpURL: testNetwork.HTTPURLs[0]}))
	// }
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
			txNodeInternalWs, err := testEnvironment.Fwd.FindPort(podName, "geth", "ws-rpc").As(envclient.RemoteConnection, envclient.WS)
			require.NoError(t, err, "Error finding WS ports")
			internalWsURLs = append(internalWsURLs, txNodeInternalWs)
			txNodeInternalHttp, err := testEnvironment.Fwd.FindPort(podName, "geth", "http-rpc").As(envclient.RemoteConnection, envclient.HTTP)
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
			ctfconfig.MustConfigOverrideChainlinkVersion(keeperTestConfig.GetChainlinkImageConfig(), target)
			ctfconfig.MightConfigOverridePyroscopeKey(keeperTestConfig.GetPyroscopeConfig(), target)
		}

		tomlConfig, err := actions.BuildTOMLNodeConfigForK8s(keeperTestConfig, testNetwork)
		require.NoError(t, err, "Error building TOML config")

		cd := chainlink.NewWithOverride(i, map[string]any{
			"toml":      tomlConfig,
			"chainlink": chainlinkResources,
			"db":        dbResources,
		}, keeperTestConfig.GetChainlinkImageConfig(), overrideFn)

		testEnvironment.AddHelm(cd)
	}
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, testNetwork
}
