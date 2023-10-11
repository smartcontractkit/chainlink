package soak

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/l2/zksync"
	"github.com/stretchr/testify/require"
)

// WIP
func TestOCRZKSync(t *testing.T) {
	l := logging.GetTestLogger(t)

	l1RpcUrl, isSet := os.LookupEnv("L1_RPC_URL")
	require.Equal(t, isSet, true, "L1_RPC_URL should be defined")

	testDuration, isSet := os.LookupEnv("TEST_DURATION")
	require.Equal(t, isSet, true, "TEST_DURATION should be defined")

	timeBetweenRounds, isSet := os.LookupEnv("OCR_TIME_BETWEEN_ROUNDS")
	require.Equal(t, isSet, true, "OCR_TIME_BETWEEN_ROUNDS should be defined")

	testEnvironment, testNetwork, err := zksync.SetupOCRTest(t)
	require.NoError(t, err, "Deploying env should not fail")
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
		}
		if answer.Int64() != int64(10) {
			l.Error().Int64("Expected answer to be 10 but got", answer.Int64())
		}
		round++
	}
}
