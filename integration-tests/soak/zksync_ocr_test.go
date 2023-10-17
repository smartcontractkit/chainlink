package soak

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/L2/zksync"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/stretchr/testify/require"
)

// WIP
func TestOCRZKSync(t *testing.T) {
	l := logging.GetTestLogger(t)

	_, insideRunner := os.LookupEnv("INSIDE_REMOTE_RUNNER")

	l1RpcUrl, isSet := os.LookupEnv("L1_RPC_URL")
	require.Equal(t, isSet, true, "L1_RPC_URL should be defined")

	testDuration, isSet := os.LookupEnv("TEST_DURATION")
	require.Equal(t, isSet, true, "TEST_DURATION should be defined")

	timeBetweenRounds, isSet := os.LookupEnv("OCR_TIME_BETWEEN_ROUNDS")
	require.Equal(t, isSet, true, "OCR_TIME_BETWEEN_ROUNDS should be defined")

	gauntletBinary, isSet := os.LookupEnv("GAUNTLET_LOCAL_BINARY")
	require.Equal(t, isSet, true, "GAUNTLET_LOCAL_BINARY should be defined")

	testEnvironment, testNetwork := DeployEnvironment(t)
	pl, err := testEnvironment.Client.ListPods(testEnvironment.Cfg.Namespace, "job-name=remote-test-runner")
	require.NoError(t, err)

	fmt.Println(pl)
	if !insideRunner {
		_, _, _, err = testEnvironment.Client.CopyToPod(testEnvironment.Cfg.Namespace, gauntletBinary, fmt.Sprintf("%s/%s:/gauntlet-evm-zksync-linux-x64", testEnvironment.Cfg.Namespace, pl.Items[0].Name), "remote-test-runner-node")
		require.NoError(t, err, "Error uploading to pod")
	}

	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	// Adding L1 URL to HTTPURLs
	testNetwork.HTTPURLs = append(testNetwork.HTTPURLs, l1RpcUrl)

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")

	// Gauntlet Setup
	zkClient, err := zksync.Setup(os.Getenv("ZK_SYNC_GOERLI_HTTP_URLS"), chainClient.GetDefaultWallet().PrivateKey(), chainClient)
	require.NoError(t, err, "Creating ZKSync client should not fail")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")

	err = zkClient.DeployOCRFeed(testEnvironment, chainClient, chainlinkNodes, testNetwork, l)
	require.NoError(t, err, "Error deploying OCR FEED")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.DebugLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	duration, err := time.ParseDuration(testDuration)
	require.NoError(t, err, "Error parsing test duration")

	waitBetweenRounds, err := time.ParseDuration(timeBetweenRounds)
	require.NoError(t, err, "Error parsing time between rounds duration")

	endTime := time.Now().Add(duration)
	round := 1
	for ; time.Now().Before(endTime); time.Sleep(waitBetweenRounds) {
		l.Info().Msg(fmt.Sprintf("Starting round %d", round))
		answer, err := zkClient.RequestOCRRound(int64(round), 10, l)
		if err != nil {
			l.Error().Err(err)
			continue
		}
		if answer.Int64() != int64(10) {
			l.Error().Int64("Expected answer to be 10 but got", answer.Int64())
		}
		round++
	}
}

func DeployEnvironment(t *testing.T) (*environment.Environment, blockchain.EVMNetwork) {
	network := networks.SelectedNetwork // Environment currently being used to soak test on
	nsPre := "soak-ocr-"

	nsPre = fmt.Sprintf("%s%s", nsPre, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"))
	baseEnvironmentConfig := &environment.Config{
		TTL:             time.Hour * 720, // 30 days,
		NamespacePrefix: nsPre,
		Test:            t,
	}

	cd := chainlink.New(0, map[string]any{
		"replicas": 6,
		"toml":     client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, config.DefaultOCRNetworkDetailTomlConfig, network),
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
	})

	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})).
		AddHelm(cd)
	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, network
}
