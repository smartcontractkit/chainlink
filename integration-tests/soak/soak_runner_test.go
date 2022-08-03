package soak_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
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

// Run the OCR soak test defined in ./tests/ocr_test.go
func TestOCRSoak(t *testing.T) {
	activeEVMNetwork := networks.GoerliTestnet // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-ocr-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	) + "-erigon"
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()
	// Geth + Prysm
	// staticValues["ETH_URL"] = "ws://34.229.200.122/ws"
	// staticValues["HTTP_ETH_URL"] = "http://34.229.200.122/rpc"
	// Besu + Teku
	// staticValues["ETH_URL"] = "ws://3.85.208.9/ws"
	// staticValues["HTTP_ETH_URL"] = "http://3.85.208.9/rpc"
	// // Nethermind + Teku
	// staticValues["ETH_URL"] = "ws://18.215.162.179/ws"
	// staticValues["HTTP_ETH_URL"] = "http://18.215.162.179/rpc"
	// // Erigon
	// staticValues["ETH_URL"] = "wss://link-eth.getblock.io/erigon/goerli/archive/axej8woh-seej-6ash-4Yu7-eyib1495dhno/"
	// staticValues["HTTP_ETH_URL"] = "https://link-eth.getblock.io/erigon/goerli/archive/axej8woh-seej-6ash-4Yu7-eyib1495dhno/"
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

	soakTestHelper(t, "@soak-ocr", testEnvironment, activeEVMNetwork)
}

// Run the keeper soak test defined in ./tests/keeper_test.go
func TestKeeperSoak(t *testing.T) {
	activeEVMNetwork := networks.GeneralEVM() // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()
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

	soakTestHelper(t, "@soak-keeper", testEnvironment, activeEVMNetwork)
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
		testEnvironment.AddHelm(chainlink.New(index, map[string]interface{}{"env": envVals}))
	}
}

// builds tests, launches environment, and triggers the soak test to run
func soakTestHelper(
	t *testing.T,
	testTag string,
	testEnvironment *environment.Environment,
	activeEVMNetwork *blockchain.EVMNetwork,
) {
	exeFile, exeFileSize, err := actions.BuildGoTests("./", "./tests", "../")
	require.NoError(t, err, "Error building go tests")

	remoteRunnerValues := map[string]interface{}{
		"test_name":      testTag,
		"env_namespace":  testEnvironment.Cfg.Namespace,
		"test_file_size": fmt.Sprint(exeFileSize),
		"test_log_level": "debug",
	}
	// Set evm network connection for remote runner
	for key, value := range activeEVMNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{"remote_test_runner": remoteRunnerValues}

	err = testEnvironment.
		AddHelm(remotetestrunner.New(remoteRunnerWrapper)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: activeEVMNetwork.Name,
			Simulated:   activeEVMNetwork.Simulated,
			WsURLs:      activeEVMNetwork.URLs,
		})).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest(exeFile, testEnvironment)
	require.NoError(t, err, "Error activating remote test")
}
