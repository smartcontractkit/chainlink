package soak_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
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
	"github.com/smartcontractkit/chainlink/integration-tests/client"
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

	replicas := 6
	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	baseTOML := `[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))
	for i := 0; i < replicas; i++ {
		testEnvironment.AddHelm(chainlink.New(i, map[string]interface{}{
			"toml": client.AddNetworksConfig(baseTOML, activeEVMNetwork),
		}))
	}

	soakTestHelper(t, testEnvironment, activeEVMNetwork)
}

// Run the OCR soak test defined in ./tests/ocr_test.go
func TestForwarderOCRSoak(t *testing.T) {
	activeEVMNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

	baseEnvironmentConfig.NamespacePrefix = fmt.Sprintf(
		"soak-forwarder-ocr-%s",
		strings.ReplaceAll(strings.ToLower(activeEVMNetwork.Name), " ", "-"),
	)

	replicas := 6
	// Values you want each node to have the exact same of (e.g. eth_chain_id)
	baseTOML := `[OCR]
Enabled = true

[Feature]
LogPoller = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`
	networkDetailTOML := `[EVM.Transactions]
ForwardersEnabled = true`
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))
	for i := 0; i < replicas; i++ {
		testEnvironment.AddHelm(chainlink.New(i, map[string]interface{}{
			"toml": client.AddNetworkDetailedConfig(baseTOML, networkDetailTOML, activeEVMNetwork),
		}))
	}
	// List of distinct Chainlink nodes to launch, and their distinct values (blank interface for none)

	soakTestHelper(t, testEnvironment, activeEVMNetwork)
}

// builds tests, launches environment, and triggers the soak test to run
func soakTestHelper(
	t *testing.T,
	testEnvironment *environment.Environment,
	activeEVMNetwork blockchain.EVMNetwork,
) {
	testDirectory := "./soak/tests"
	log.Info().
		Str("Name", t.Name()).
		Str("Directory", testDirectory).
		Str("Namespace", testEnvironment.Cfg.Namespace).
		Msg("Soak Test")
	remoteRunnerValues := actions.BasicRunnerValuesSetup(
		t.Name(),
		testEnvironment.Cfg.Namespace,
		testDirectory,
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
