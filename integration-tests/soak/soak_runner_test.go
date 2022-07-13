package soak_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/remotetestrunner"
	networks "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/stretchr/testify/require"
)

func init() {
	logging.Init()
}

var baseEnvironmentConfig = &environment.Config{
	TTL: time.Hour * 720, // 30 days,
}

func TestOCRSoak(t *testing.T) {
	activeEVMNetwork := networks.SepoliaTestnet // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-ocr-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()
	staticValues["ETH_MAX_GAS_PRICE_WEI"] = "100000000000"
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

// Cannot boot Chainlink: fatal error instantiating application: failed to load EVM chainset: cannot create new chain with ID 11155111, config validation failed: EVM_GAS_FEE_CAP_DEFAULT (100000000000) must be less than or equal to ETH_MAX_GAS_PRICE_WEI (1000); ETH_MAX_GAS_PRICE_WEI must be greater than or equal to ETH_GAS_PRICE_DEFAULT
func TestKeeperSoak(t *testing.T) {
	activeEVMNetwork := networks.SepoliaTestnet // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-keeper-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig)

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()
	staticValues["ETH_MAX_GAS_PRICE_WEI"] = "100000000000"
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
