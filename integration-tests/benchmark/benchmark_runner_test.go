package benchmark_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/imdario/mergo"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/remotetestrunner"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
)

func init() {
	logging.Init()
}

var baseEnvironmentConfig = &environment.Config{
	TTL: time.Hour * 720, // 30 days,
}

var dynamicValues_EvmNodes = []map[string]interface{}{
	{
		"EVM_NODES": os.Getenv("EVM_NODES_0"),
	},
	{
		"EVM_NODES": os.Getenv("EVM_NODES_1"),
	},
	{
		"EVM_NODES": os.Getenv("EVM_NODES_2"),
	},
	{
		"EVM_NODES": os.Getenv("EVM_NODES_3"),
	},
	{
		"EVM_NODES": os.Getenv("EVM_NODES_4"),
	},
	{
		"EVM_NODES": os.Getenv("EVM_NODES_5"),
	},
}

var chainlinkPerformance = map[string]interface{}{
	"chainlink": map[string]interface{}{
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
	},
	"db": map[string]interface{}{
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
		"capacity": "20Gi",
	},
}

var chainlinkSoak = map[string]interface{}{
	"chainlink": map[string]interface{}{
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
	},
	"db": map[string]interface{}{
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
		"capacity": "20Gi",
	},
}

func TestKeeperBenchmarkSimulated(t *testing.T) {
	KeeperBenchmarkSimulated(t, "@simulated-registry-1-3 @benchmark-keeper")
}

func TestKeeperBenchmarkSimulated1_2(t *testing.T) {
	KeeperBenchmarkSimulated(t, "@simulated-registry1-2 @benchmark-keeper")
}

func TestKeeperBenchmarkSimulatedMultiple(t *testing.T) {
	KeeperBenchmarkSimulated(t, "@simulated-multiple-registries @benchmark-keeper")
}

func TestKeeperBenchmarkGoerli(t *testing.T) {
	KeeperBenchmarkGoerli(t, "@goerli-registry-1-3 @benchmark-keeper")
}

func TestKeeperBenchmarkGoerli1_2(t *testing.T) {
	KeeperBenchmarkGoerli(t, "@goerli-registry1-2 @benchmark-keeper")
}

func TestKeeperBenchmarkGoerliMultiple(t *testing.T) {
	KeeperBenchmarkGoerli(t, "@goerli-multiple-registries @benchmark-keeper")
}

func TestKeeperBenchmarkArbitrumGoerli(t *testing.T) {
	KeeperBenchmarkArbitrumGoerli(t, "@arbitrum-goerli-registry-1-3 @benchmark-keeper")
}

func TestKeeperBenchmarkOptimisticGoerli(t *testing.T) {
	KeeperBenchmarkOptimisticGoerli(t, "@optimistic-goerli-registry-1-3 @benchmark-keeper")
}

// Run the Keepers Benchmark test defined in ./tests/keeper_test.go
func KeeperBenchmarkSimulated(t *testing.T, testTag string) {
	activeEVMNetwork := networks.SimulatedEVM // Environment currently being used to run benchmark test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"benchmark-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()

	keeperBenchmarkValues := map[string]interface{}{
		"MIN_INCOMING_CONFIRMATIONS": "1",
		"KEEPER_TURN_FLAG_ENABLED":   "true",
		"CHAINLINK_DEV":              "false",
		//"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD": "2000000",
	}
	mergo.Merge(&staticValues, &keeperBenchmarkValues)
	// List of distinct Chainlink nodes to launch, and their distinct values (blank interface for none)
	dynamicValues := []map[string]interface{}{
		{
			"dynamic_value": "0",
		},
		{
			"dynamic_value": "1",
		},
		{
			"dynamic_value": "2",
		},
		{
			"dynamic_value": "3",
		},
		{
			"dynamic_value": "4",
		},
		{
			"dynamic_value": "5",
		},
	}
	addSeparateChainlinkDeployments(testEnvironment, staticValues, dynamicValues)

	benchmarkTestHelper(t, testTag, testEnvironment, activeEVMNetwork)
}

func KeeperBenchmarkGoerli(t *testing.T, testTag string) {
	activeEVMNetwork := networks.GoerliTestnet // Environment currently being used to run benchmark test on
	activeEVMNetwork.MinimumConfirmations = 1
	activeEVMNetwork.Timeout = 10 * time.Minute

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"benchmark-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := map[string]interface{}{
		"KEEPER_REGISTRY_SYNC_INTERVAL":  "",
		"ETH_URL":                        "",
		"ETH_CHAIN_ID":                   "",
		"ETH_MAX_IN_FLIGHT_TRANSACTIONS": "3",
		"ETH_MAX_QUEUED_TRANSACTIONS":    "15",
		"ETH_GAS_BUMP_TX_DEPTH":          "3",
	}

	keeperBenchmarkValues := map[string]interface{}{
		"MIN_INCOMING_CONFIRMATIONS": "1",
		"KEEPER_TURN_FLAG_ENABLED":   "true",
		"CHAINLINK_DEV":              "false",
	}
	mergo.Merge(&staticValues, &keeperBenchmarkValues)
	addSeparateChainlinkDeployments(testEnvironment, staticValues, dynamicValues_EvmNodes)

	benchmarkTestHelper(t, testTag, testEnvironment, activeEVMNetwork)
}

func KeeperBenchmarkArbitrumGoerli(t *testing.T, testTag string) {
	activeEVMNetwork := networks.ArbitrumGoerli // Environment currently being used to run benchmark test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"benchmark-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := map[string]interface{}{
		"KEEPER_REGISTRY_SYNC_INTERVAL":  "",
		"ETH_URL":                        "",
		"ETH_CHAIN_ID":                   "",
		"ETH_MAX_IN_FLIGHT_TRANSACTIONS": "",
		"ETH_MAX_QUEUED_TRANSACTIONS":    "",
		"ETH_GAS_BUMP_TX_DEPTH":          "",
	}

	keeperBenchmarkValues := map[string]interface{}{
		"KEEPER_TURN_FLAG_ENABLED": "true",
		"CHAINLINK_DEV":            "false",
	}
	mergo.Merge(&staticValues, &keeperBenchmarkValues)
	addSeparateChainlinkDeployments(testEnvironment, staticValues, dynamicValues_EvmNodes)

	benchmarkTestHelper(t, testTag, testEnvironment, activeEVMNetwork)
}

func KeeperBenchmarkOptimisticGoerli(t *testing.T, testTag string) {
	activeEVMNetwork := networks.OptimismGoerli // Environment currently being used to run benchmark test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"benchmark-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := map[string]interface{}{
		"KEEPER_REGISTRY_SYNC_INTERVAL":  "",
		"ETH_URL":                        "",
		"ETH_CHAIN_ID":                   "",
		"ETH_MAX_IN_FLIGHT_TRANSACTIONS": "",
		"ETH_MAX_QUEUED_TRANSACTIONS":    "",
		"ETH_GAS_BUMP_TX_DEPTH":          "",
	}

	keeperBenchmarkValues := map[string]interface{}{
		"KEEPER_TURN_FLAG_ENABLED": "true",
		"CHAINLINK_DEV":            "false",
	}
	mergo.Merge(&staticValues, &keeperBenchmarkValues)
	addSeparateChainlinkDeployments(testEnvironment, staticValues, dynamicValues_EvmNodes)

	benchmarkTestHelper(t, testTag, testEnvironment, activeEVMNetwork)
}

// adds distinct Chainlink deployments to the test environment, using staticVals on all of them, while distributing
// a single dynamicVal to each Chainlink deployment
func addSeparateChainlinkDeployments(
	testEnvironment *environment.Environment,
	staticValues map[string]interface{},
	dynamicValueList []map[string]interface{},
) {
	for index, dynamicValues := range dynamicValueList {
		envVals := map[string]interface{}{}
		for key, value := range staticValues {
			envVals[key] = value
		}
		for key, value := range dynamicValues {
			envVals[key] = value
		}
		chartValues := map[string]interface{}{
			"env": envVals,
		}
		chartResources := chainlinkPerformance
		testType, testTypeExists := os.LookupEnv("TEST_TYPE")
		if testTypeExists && strings.ToLower(testType) == "soak" {
			chartResources = chainlinkSoak
		}
		mergo.Merge(&chartValues, &chartResources)
		testEnvironment.AddHelm(chainlink.New(index, chartValues))
	}
}

// builds tests, launches environment, and triggers the benchmark test to run
func benchmarkTestHelper(
	t *testing.T,
	testTag string,
	testEnvironment *environment.Environment,
	activeEVMNetwork *blockchain.EVMNetwork,
) {
	exeFile, exeFileSize, err := actions.BuildGoTests("./", "./tests", "../")
	require.NoError(t, err, "Error building go tests")

	remoteRunnerValues := map[string]interface{}{
		"test_name":             testTag,
		"env_namespace":         testEnvironment.Cfg.Namespace,
		"test_file_size":        fmt.Sprint(exeFileSize),
		"test_log_level":        "debug",
		"grafana_dashboard_url": os.Getenv("GRAFANA_DASHBOARD_URL"),
		"NUMBEROFCONTRACTS":     os.Getenv("NUMBEROFCONTRACTS"),
		"CHECKGASTOBURN":        os.Getenv("CHECKGASTOBURN"),
		"PERFORMGASTOBURN":      os.Getenv("PERFORMGASTOBURN"),
		"BLOCKRANGE":            os.Getenv("BLOCKRANGE"),
		"BLOCKINTERVAL":         os.Getenv("BLOCKINTERVAL"),
		"CHAINLINKNODEFUNDING":  os.Getenv("CHAINLINKNODEFUNDING"),
	}
	// Set evm network connection for remote runner
	for key, value := range activeEVMNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{
		"remote_test_runner": remoteRunnerValues,
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "250m",
				"memory": "1Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "250m",
				"memory": "1Gi",
			},
		},
	}

	err = testEnvironment.
		AddHelm(remotetestrunner.New(remoteRunnerWrapper)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: activeEVMNetwork.Name,
			Simulated:   activeEVMNetwork.Simulated,
			WsURLs:      activeEVMNetwork.URLs,
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
			},
		})).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest(exeFile, testEnvironment)
	require.NoError(t, err, "Error activating remote test")
}
