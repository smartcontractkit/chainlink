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

func TestKeeperBenchmark(t *testing.T) {
	registryToTest := os.Getenv("AUTOMATION_REGISTRY_TO_TEST")
	KeeperBenchmark(t, registryToTest)
}

// Run the Keepers Benchmark test defined in ./tests/keeper_test.go
func KeeperBenchmark(t *testing.T, registryToTest string) {
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to run benchmark test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"benchmark-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := map[string]interface{}{
		"ETH_URL":      activeEVMNetwork.URLs[0],
		"ETH_HTTP_URL": activeEVMNetwork.HTTPURLs[0],
	}

	keeperBenchmarkValues := map[string]interface{}{
		"MIN_INCOMING_CONFIRMATIONS": "1",
		"KEEPER_TURN_FLAG_ENABLED":   "true",
		"CHAINLINK_DEV":              "false",
	}

	testTag := "simulated"

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

	if !activeEVMNetwork.Simulated {
		staticValues = map[string]interface{}{
			"KEEPER_REGISTRY_SYNC_INTERVAL": "",
			"ETH_URL":                       "",
			"ETH_CHAIN_ID":                  "",
		}
		dynamicValues = dynamicValues_EvmNodes
		if activeEVMNetwork.Name == "Goerli Testnet" {
			keeperBenchmarkValues = map[string]interface{}{
				"MIN_INCOMING_CONFIRMATIONS":     "1",
				"KEEPER_TURN_FLAG_ENABLED":       "true",
				"CHAINLINK_DEV":                  "false",
				"ETH_MAX_IN_FLIGHT_TRANSACTIONS": "3",
				"ETH_MAX_QUEUED_TRANSACTIONS":    "15",
				"ETH_GAS_BUMP_TX_DEPTH":          "3",
			}
			testTag = "goerli"
		}
		if activeEVMNetwork.Name == "Arbitrum Goerli" || activeEVMNetwork.Name == "Optimism Goerli" {
			keeperBenchmarkValues = map[string]interface{}{
				"KEEPER_TURN_FLAG_ENABLED":       "true",
				"CHAINLINK_DEV":                  "false",
				"ETH_MAX_IN_FLIGHT_TRANSACTIONS": "",
				"ETH_MAX_QUEUED_TRANSACTIONS":    "",
				"ETH_GAS_BUMP_TX_DEPTH":          "",
			}
			testTag = "arbitrum-goerli"
		}
		if activeEVMNetwork.Name == "Optimism Goerli" {
			testTag = "optimistic-goerli"
		}
	}

	testTag = "@" + testTag + "-" + registryToTest

	mergo.Merge(&staticValues, &keeperBenchmarkValues)

	addSeparateChainlinkDeployments(testEnvironment, staticValues, dynamicValues)

	benchmarkTestHelper(t, testTag+" @benchmark-keeper", testEnvironment, activeEVMNetwork)
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
		testEnvironment.AddHelm(chainlink.NewVersioned(index, "0.0.11", chartValues))
	}
}

// builds tests, launches environment, and triggers the benchmark test to run
func benchmarkTestHelper(
	t *testing.T,
	testTag string,
	testEnvironment *environment.Environment,
	activeEVMNetwork *blockchain.EVMNetwork,
) {
	remoteRunnerValues := map[string]interface{}{
		"focus":                 testTag,
		"env_namespace":         testEnvironment.Cfg.Namespace,
		"test_dir":              "./suite/benchmark/tests",
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
	}

	err := testEnvironment.
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
	err = actions.TriggerRemoteTest("../../", testEnvironment)
	require.NoError(t, err, "Error activating remote test")
}
