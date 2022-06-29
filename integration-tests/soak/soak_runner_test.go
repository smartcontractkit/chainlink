package soak_test

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

	networks "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/stretchr/testify/require"
)

func init() {
	logging.Init()
	networks.LoadNetworks("../.env")
}

func TestOCRSoak(t *testing.T) {
	soakTestHelper(t, "@soak-ocr", "soak-ocr", 6, networks.MetisTestNetwork)
}

func TestKeeperSoak(t *testing.T) {
	soakTestHelper(t, "@soak-keeper", "soak-keeper", 6, networks.SimulatedEVMNetwork)
}

func soakTestHelper(
	t *testing.T,
	testTag, namespacePrefix string,
	chainlinkReplicas int,
	evmNetwork *blockchain.EVMNetwork,
) {
	exeFile, exeFileSize, err := actions.BuildGoTests("./", "./tests", "../")
	require.NoError(t, err, "Error building go tests")
	env := environment.New(&environment.Config{
		TTL:             999 * time.Hour,
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
	}
	if !evmNetwork.Simulated {
		chainlinkVals["env"] = map[string]interface{}{
			"eth_url":      evmNetwork.URLs[0],
			"eth_chain_id": fmt.Sprint(networks.MetisTestNetwork.ChainID),
		}
	}

	err = env.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(remotetestrunner.New(remoteRunnerWrapper)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: evmNetwork.Name,
			Simulated:   evmNetwork.Simulated,
			WsURLs:      evmNetwork.URLs,
		})).
		AddHelm(chainlink.New(0, chainlinkVals)).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = actions.TriggerRemoteTest(exeFile, env)
	require.NoError(t, err, "Error activating remote test")
}
