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

// Run the Keepers Benchmark test defined in ./tests/keeper_test.go
func TestKeeperBenchmark(t *testing.T) {
	activeEVMNetwork := networks.SimulatedEVM // Environment currently being used to soak test on

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

	benchmarkTestHelper(t, "@benchmark-keeper", testEnvironment, activeEVMNetwork)
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
		testEnvironment.AddHelm(chainlink.New(index, map[string]interface{}{
			"env": envVals,
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
				"capacity": "5Gi",
			},
		}))
	}
}

// builds tests, launches environment, and triggers the soak test to run
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
	}
	// Set evm network connection for remote runner
	for key, value := range activeEVMNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{
		"remote_test_runner": remoteRunnerValues,
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "512Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "512Mi",
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
