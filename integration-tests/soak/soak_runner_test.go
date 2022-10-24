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
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-ocr-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))

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

	soakTestHelper(t, "@soak-ocr", testEnvironment, activeEVMNetwork)
}

// Run the OCR soak test defined in ./tests/ocr_test.go
func TestForwarderOCRSoak(t *testing.T) {
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-forwarder-ocr-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))

	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	staticValues := activeEVMNetwork.ChainlinkValuesMap()
	staticValues["ETH_USE_FORWARDERS"] = "true"
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

	soakTestHelper(t, "@soak-forwarder-ocr", testEnvironment, activeEVMNetwork)
}

// Run the keeper soak test defined in ./tests/keeper_test.go
func TestKeeperSoak(t *testing.T) {
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

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
	remoteRunnerValues := actions.BasicRunnerValuesSetup(
		testTag,
		testEnvironment.Cfg.Namespace,
		"./integration-tests/soak/tests",
	)
	// Set evm network connection for remote runner
	for key, value := range activeEVMNetwork.ToMap() {
		remoteRunnerValues[key] = value
	}
	remoteRunnerWrapper := map[string]interface{}{"remote_test_runner": remoteRunnerValues}

	err := testEnvironment.
		AddHelm(remotetestrunner.New(remoteRunnerWrapper)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: activeEVMNetwork.Name,
			Simulated:   activeEVMNetwork.Simulated,
			WsURLs:      activeEVMNetwork.URLs,
		})).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest("../../", testEnvironment)
	require.NoError(t, err, "Error activating remote test")
}
