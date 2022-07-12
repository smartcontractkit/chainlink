package benchmark_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/remotetestrunner"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/stretchr/testify/require"
)

var (
	// Keepers Benchmark EVM ensures that the test will use a custom simulated geth instance
	KeepersBenchmarkEVM *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:      "Simulated Geth",
		Simulated: true,
		ChainID:   1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

func init() {
	logging.Init()
}

func TestKeeperBenchmark(t *testing.T) {
	benchmarkTestHelper(t, "@benchmark-keeper", "benchmark-keeper", 6, KeepersBenchmarkEVM)
}

func benchmarkTestHelper(
	t *testing.T,
	testTag, namespacePrefix string,
	chainlinkReplicas int,
	evmNetwork *blockchain.EVMNetwork,
) {
	exeFile, exeFileSize, err := actions.BuildGoTests("./", "./tests", "../")
	require.NoError(t, err, "Error building go tests")
	env := environment.New(&environment.Config{
		TTL:             24 * time.Hour, // 1 day limit
		Labels:          []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5RemoteRunner)},
		NamespacePrefix: namespacePrefix,
	})

	remoteRunnerValues := map[string]interface{}{
		"test_name":      testTag,
		"env_namespace":  env.Cfg.Namespace,
		"test_file_size": fmt.Sprint(exeFileSize),
		"log_level":      "debug",
	}
	// Set evm network connection for remote runner
	for key, value := range evmNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{"remote_test_runner": remoteRunnerValues}

	// Set Chainlink vals
	chainlinkVals := map[string]interface{}{
		"replicas": chainlinkReplicas,
		"env": map[string]interface{}{
			"MIN_INCOMING_CONFIRMATIONS": "1",
			"KEEPER_TURN_FLAG_ENABLED":   "true",
		},
		"chainlink": map[string]interface{}{
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "4086Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "4086Mi",
				},
			},
		},
		"db": map[string]interface{}{
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "1024Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "1024Mi",
				},
			},
		},
	}

	err = env.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(remotetestrunner.New(remoteRunnerWrapper)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: evmNetwork.Name,
			Simulated:   evmNetwork.Simulated,
			Values: map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "1500m",
						"memory": "4086Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "1500m",
						"memory": "4086Mi",
					},
				},
			}})).
		AddHelm(chainlink.New(0, chainlinkVals)).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest(exeFile, env)
	require.NoError(t, err, "Error activating remote test")
}
