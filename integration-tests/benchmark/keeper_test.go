package benchmark

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	env_client "github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
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

var (
	NumberOfNodes, _        = strconv.Atoi(getEnv("NUMBEROFNODES", "6"))
	RegistryToTest          = getEnv("REGISTRY", "2_1")
	NumberOfUpkeeps, _      = strconv.Atoi(getEnv("NUMBEROFUPKEEPS", "500"))
	CheckGasToBurn, _       = strconv.ParseInt(getEnv("CHECKGASTOBURN", "100000"), 0, 64)
	PerformGasToBurn, _     = strconv.ParseInt(getEnv("PERFORMGASTOBURN", "50000"), 0, 64)
	BlockRange, _           = strconv.ParseInt(getEnv("BLOCKRANGE", "3600"), 0, 64)
	BlockInterval, _        = strconv.ParseInt(getEnv("BLOCKINTERVAL", "20"), 0, 64)
	ChainlinkNodeFunding, _ = strconv.ParseFloat(getEnv("CHAINLINKNODEFUNDING", "0.5"), 64)
	MaxPerformGas, _        = strconv.ParseInt(getEnv("MAXPERFORMGAS", "5000000"), 0, 32)
	UpkeepGasLimit, _       = strconv.ParseInt(getEnv("UPKEEPGASLIMIT", fmt.Sprint(PerformGasToBurn+100000)), 0, 64)
	NumberOfRegistries, _   = strconv.Atoi(getEnv("NUMBEROFREGISTRIES", "1"))
	ForceSingleTxnKey, _    = strconv.ParseBool(getEnv("FORCESINGLETXNKEY", "false"))
	DeleteJobsOnEnd, _      = strconv.ParseBool(getEnv("DELETEJOBSONEND", "true"))
	RegistryAddress         = getEnv("REGISTRYADDRESS", "")
	RegistrarAddress        = getEnv("REGISTRARADDRESS", "")
	LinkTokenAddress        = getEnv("LINKTOKENADDRESS", "")
	EthFeedAddress          = getEnv("ETHFEEDADDRESS", "")
	GasFeedAddress          = getEnv("GASFEEDADDRESS", "")
)

type NetworkConfig struct {
	upkeepSLA  int64
	blockTime  time.Duration
	deltaStage time.Duration
	funding    *big.Float
}

func TestAutomationBenchmark(t *testing.T) {
	l := logging.GetTestLogger(t)
	testEnvironment, benchmarkNetwork := SetupAutomationBenchmarkEnv(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}
	networkName := strings.ReplaceAll(benchmarkNetwork.Name, " ", "")
	testName := fmt.Sprintf("%s%s", networkName, RegistryToTest)
	l.Info().Str("Test Name", testName).Str("Test Inputs", os.Getenv("TEST_INPUTS")).Msg("Running Benchmark Test")
	benchmarkTestNetwork := networkConfig[networkName]

	l.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")

	chainClient, err := blockchain.NewEVMClient(benchmarkNetwork, testEnvironment, l)
	require.NoError(t, err, "Error connecting to blockchain")
	registryVersions := addRegistry(RegistryToTest)
	keeperBenchmarkTest := testsetups.NewKeeperBenchmarkTest(t,
		testsetups.KeeperBenchmarkTestInputs{
			BlockchainClient: chainClient,
			RegistryVersions: registryVersions,
			KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
				PaymentPremiumPPB:    uint32(0),
				BlockCountPerTurn:    big.NewInt(100),
				CheckGasLimit:        uint32(45_000_000), //45M
				StalenessSeconds:     big.NewInt(90_000),
				GasCeilingMultiplier: uint16(2),
				MaxPerformGas:        uint32(MaxPerformGas),
				MinUpkeepSpend:       big.NewInt(0),
				FallbackGasPrice:     big.NewInt(2e11),
				FallbackLinkPrice:    big.NewInt(2e18),
				MaxCheckDataSize:     uint32(5_000),
				MaxPerformDataSize:   uint32(5_000),
			},
			Upkeeps: &testsetups.UpkeepConfig{
				NumberOfUpkeeps:     NumberOfUpkeeps,
				CheckGasToBurn:      CheckGasToBurn,
				PerformGasToBurn:    PerformGasToBurn,
				BlockRange:          BlockRange,
				BlockInterval:       BlockInterval,
				UpkeepGasLimit:      UpkeepGasLimit,
				FirstEligibleBuffer: 1,
			},
			Contracts: &testsetups.PreDeployedContracts{
				RegistrarAddress: RegistrarAddress,
				RegistryAddress:  RegistryAddress,
				LinkTokenAddress: LinkTokenAddress,
				EthFeedAddress:   EthFeedAddress,
				GasFeedAddress:   GasFeedAddress,
			},
			ChainlinkNodeFunding: benchmarkTestNetwork.funding,
			UpkeepSLA:            benchmarkTestNetwork.upkeepSLA,
			BlockTime:            benchmarkTestNetwork.blockTime,
			DeltaStage:           benchmarkTestNetwork.deltaStage,
			ForceSingleTxnKey:    ForceSingleTxnKey,
			DeleteJobsOnEnd:      DeleteJobsOnEnd,
		},
	)
	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})
	keeperBenchmarkTest.Setup(testEnvironment)
	keeperBenchmarkTest.Run()
}

func addRegistry(registryToTest string) []eth_contracts.KeeperRegistryVersion {
	switch registryToTest {
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
		return repeatRegistries(eth_contracts.RegistryVersion_2_0, NumberOfRegistries)
	case "2_1-Multiple":
		return repeatRegistries(eth_contracts.RegistryVersion_1_0, NumberOfRegistries)
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

var networkConfig = map[string]NetworkConfig{
	"SimulatedGeth": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	"geth": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: 30 * time.Second,
		funding:    big.NewFloat(100_000),
	},
	"GoerliTestnet": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
		funding:    big.NewFloat(ChainlinkNodeFunding),
	},
	"ArbitrumGoerli": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: time.Duration(0),
		funding:    big.NewFloat(ChainlinkNodeFunding),
	},
	"OptimismGoerli": {
		upkeepSLA:  int64(20),
		blockTime:  time.Second,
		deltaStage: time.Duration(0),
		funding:    big.NewFloat(ChainlinkNodeFunding),
	},
	"SepoliaTestnet": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
		funding:    big.NewFloat(ChainlinkNodeFunding),
	},
	"PolygonMumbai": {
		upkeepSLA:  int64(4),
		blockTime:  12 * time.Second,
		deltaStage: time.Duration(0),
		funding:    big.NewFloat(ChainlinkNodeFunding),
	},
}

func getEnv(key, fallback string) string {
	if inputs, ok := os.LookupEnv("TEST_INPUTS"); ok {
		values := strings.Split(inputs, ",")
		for _, value := range values {
			if strings.Contains(value, key) {
				return strings.Split(value, "=")[1]
			}
		}
	}
	return fallback
}

func SetupAutomationBenchmarkEnv(t *testing.T) (*environment.Environment, blockchain.EVMNetwork) {
	l := logging.GetTestLogger(t)
	testNetwork := networks.SelectedNetwork // Environment currently being used to run benchmark test on
	blockTime := "1"
	networkDetailTOML := `MinIncomingConfirmations = 1`

	if strings.Contains(RegistryToTest, "2_") {
		NumberOfNodes++
	}

	testType := strings.ToLower(getEnv("TEST_TYPE", "benchmark"))
	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s-%s",
			testType,
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
			strings.ReplaceAll(strings.ToLower(RegistryToTest), "_", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})
	// propagate TEST_INPUTS to remote runner
	if testEnvironment.WillUseRemoteRunner() {
		key := "TEST_INPUTS"
		err := os.Setenv(fmt.Sprintf("TEST_%s", key), os.Getenv(key))
		require.NoError(t, err, "failed to set the environment variable TEST_INPUTS for remote runner")
		key = "GRAFANA_DASHBOARD_URL"
		err = os.Setenv(fmt.Sprintf("TEST_%s", key), getEnv(key, ""))
		require.NoError(t, err, "failed to set the environment variable GRAFANA_DASHBOARD_URL for remote runner")
	}

	dbResources := performanceDbResources
	chainlinkResources := performanceChainlinkResources
	if testType == "soak" {
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
							"replicas": NumberOfNodes,
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
	for i := 0; i < NumberOfNodes; i++ {
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

	for i := 0; i < NumberOfNodes; i++ {
		testNetwork.HTTPURLs = []string{internalHttpURLs[i]}
		testNetwork.URLs = []string{internalWsURLs[i]}
		testEnvironment.AddHelm(chainlink.New(i, map[string]any{
			"toml":      client.AddNetworkDetailedConfig(keeperBenchmarkBaseTOML, networkDetailTOML, testNetwork),
			"chainlink": chainlinkResources,
			"db":        dbResources,
		}))
	}
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, testNetwork
}
