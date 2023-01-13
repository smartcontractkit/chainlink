package benchmark_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/remotetestrunner"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

func init() {
	logging.Init()
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

var (
	baseEnvironmentConfig = &environment.Config{
		TTL: time.Hour * 720, // 30 days,
	}
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
		"capacity": "20Gi",
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
		"capacity": "20Gi",
	}
)

func TestAutomationBenchmark(t *testing.T) {
	registryToTest := getEnv("AUTOMATION_REGISTRY_TO_TEST", "Registry_2_0")
	var numberOfNodes, _ = strconv.Atoi(getEnv("AUTOMATION_NUMBER_OF_NODES", "6"))
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to run benchmark test on
	blockTime := "1"

	baseTOML := `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	networkDetailTOML := `MinIncomingConfirmations = 1`

	if registryToTest == "Registry_2_0" {
		numberOfNodes += 1
		blockTime = "12"
	}

	testType := strings.ToLower(getEnv("TEST_TYPE", "benchmark"))
	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"automation-%s-%s-%s",
		testType,
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
		strings.ReplaceAll(strings.ToLower(registryToTest), "_", "-"),
	)
	dbResources := performanceDbResources
	chainlinkResources := performanceChainlinkResources
	if testType == "soak" {
		chainlinkResources = soakChainlinkResources
		dbResources = soakDbResources
	}

	testEnvironment := environment.New(baseEnvironmentConfig)
	for i := 0; i < numberOfNodes; i++ {
		testEnvironment.
			AddHelm(chainlink.New(i, map[string]interface{}{
				"toml":      client.AddNetworkDetailedConfig(baseTOML, networkDetailTOML, activeEVMNetwork),
				"chainlink": chainlinkResources,
				"db":        dbResources,
			}))
	}

	networkTestName := strings.ReplaceAll(activeEVMNetwork.Name, " ", "")
	testName := fmt.Sprintf("TestKeeperBenchmark%s%s", networkTestName, registryToTest)
	log.Info().Str("Test Name", testName).Msg("Running Benchmark Test")
	benchmarkTestHelper(t, testName, testEnvironment, activeEVMNetwork, blockTime, numberOfNodes)
}

// builds tests, launches environment, and triggers the benchmark test to run
func benchmarkTestHelper(
	t *testing.T,
	testName string,
	testEnvironment *environment.Environment,
	activeEVMNetwork blockchain.EVMNetwork,
	blockTime string,
	nodeReplicas int,
) {
	testDirectory := "./benchmark/tests"
	log.Info().
		Str("Test Name", testName).
		Str("Directory", testDirectory).
		Str("Namespace", testEnvironment.Cfg.Namespace).
		Msg("Benchmark Test")
	remoteRunnerValues := map[string]interface{}{
		"test_name":             testName,
		"env_namespace":         testEnvironment.Cfg.Namespace,
		"test_dir":              testDirectory,
		"test_log_level":        "debug",
		"grafana_dashboard_url": getEnv("GRAFANA_DASHBOARD_URL", ""),
		"TEST_INPUTS":           os.Getenv("TEST_INPUTS"),
		"SELECTED_NETWORKS":     os.Getenv("SELECTED_NETWORKS"),
	}
	// Set evm network connection for remote runner
	for key, value := range activeEVMNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{
		"remote_test_runner": remoteRunnerValues,
	}

	if activeEVMNetwork.Simulated {
		testEnvironment.
			AddChart(blockscout.New(&blockscout.Props{
				Name:    "geth-blockscout",
				WsURL:   activeEVMNetwork.URL,
				HttpURL: activeEVMNetwork.HTTPURLs[0]}))
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
				"geth": map[string]interface{}{
					"blocktime": blockTime,
				},
			},
		})).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest("../../", testEnvironment)
	require.NoError(t, err, "Error activating remote test")
}
