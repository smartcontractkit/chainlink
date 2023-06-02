package soak

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestVRFV2Soak(t *testing.T) {
	var testInputs testsetups.VRFV2SoakTestInputs
	err := envconfig.Process("VRFV2", &testInputs)
	require.NoError(t, err, "Error reading VRFV2 soak test inputs")
	vrfSubscriptionFundingAmountInLink := testInputs.SubscriptionFunding
	chainlinkNodeFundingAmountEth := testInputs.ChainlinkNodeFunding

	waitForRandRequestStatusToBeFulfilledTimeout := time.Second * 40
	l := utils.GetTestLogger(t)

	testInputs.SetForRemoteRunner()
	testNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

	testEnvironment := vrfv2_actions.SetupVRFV2Environment(
		t,
		testNetwork,
		config.BaseVRFV2NetworkDetailTomlConfig,
		"",
		"soak-vrfv2",
		"",
		time.Minute*20,
	)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)
	chainClient.ParallelTransactions(true)

	mockETHLINKFeed, err := contractDeployer.DeployMockETHLINKFeed(vrfv2_constants.LinkEthFeedResponse)
	require.NoError(t, err)
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err)

	vrfV2Contracts, chainlinkNodesAfterRedeployment, vrfV2jobs, testEnvironmentAfterRedeployment := vrfv2_actions.SetupVRFV2Universe(
		t,
		linkToken,
		mockETHLINKFeed,
		contractDeployer,
		chainClient,
		chainlinkNodes,
		testNetwork,
		testEnvironment,
		chainlinkNodeFundingAmountEth,
		vrfSubscriptionFundingAmountInLink,
		"soak-vrfv2",
		time.Hour*1,
	)

	vrfV2SoakTest := testsetups.NewVRFV2SoakTest(&testsetups.VRFV2SoakTestInputs{
		BlockchainClient:     chainClient,
		TestDuration:         testInputs.TestDuration,
		ChainlinkNodeFunding: chainlinkNodeFundingAmountEth,
		SubscriptionFunding:  vrfSubscriptionFundingAmountInLink,
		StopTestOnError:      testInputs.StopTestOnError,
		RequestsPerMinute:    testInputs.RequestsPerMinute,
		TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int, wg *sync.WaitGroup) error {

			concurrentEVMClient, err := blockchain.ConcurrentEVMClient(testNetwork, testEnvironmentAfterRedeployment, chainClient)
			vrfV2Contracts.LoadTestConsumer.ChangeEVMClient(concurrentEVMClient)
			if err != nil {
				return fmt.Errorf("error occurred creating ConcurrentEVMClient, error: %w", err)
			}

			wg.Add(1)
			// request randomness
			err = vrfV2Contracts.LoadTestConsumer.RequestRandomness(
				vrfV2jobs[0].KeyHash,
				vrfv2_constants.SubID,
				uint16(vrfv2_constants.MinimumConfirmations),
				vrfv2_constants.CallbackGasLimit,
				vrfv2_constants.NumberOfWords,
				1,
			)
			if err != nil {
				return fmt.Errorf("error occurred Requesting Randomness, error: %w", err)
			}

			err = concurrentEVMClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error occurred waiting on chain events, error: %w", err)
			}

			lastRequestID, err := vrfV2Contracts.LoadTestConsumer.GetLastRequestId(context.Background())
			if err != nil {
				return fmt.Errorf("error occurred getting Last Request ID, error: %w", err)
			}

			l.Info().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

			_, err = WaitForRandRequestToBeFulfilled(
				vrfV2Contracts.LoadTestConsumer,
				lastRequestID,
				waitForRandRequestStatusToBeFulfilledTimeout,
				wg,
				t,
			)

			if err != nil {
				return fmt.Errorf("error occurred waiting for Randomness Request Status with request ID: %v, error: %w", lastRequestID.String(), err)
			}

			l.Info().
				Int("Request Number", requestNumber).
				Str("RequestID", lastRequestID.String()).
				Msg("Randomness fulfilled")
			return nil
		},
	},
		chainlinkNodesAfterRedeployment)

	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(vrfV2SoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	vrfV2SoakTest.Setup(t, testEnvironment)
	l.Info().Msg("Set up soak test")
	vrfV2SoakTest.Run(t)
}

func WaitForRandRequestToBeFulfilled(
	consumer contracts.VRFv2LoadTestConsumer,
	lastRequestID *big.Int,
	timeout time.Duration,
	wg *sync.WaitGroup,
	t *testsetups.VRFV2SoakTest,
) (vrf_load_test_with_metrics.GetRequestStatus, error) {
	requestStatusChannel := make(chan vrf_load_test_with_metrics.GetRequestStatus)
	requestStatusErrorsChannel := make(chan error)

	// set the requests to only run for a certain amount of time
	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return vrf_load_test_with_metrics.GetRequestStatus{}, fmt.Errorf("timeout waiting for new transmission event")

		case <-ticker.C:
			go getRandomnessRequestStatus(
				consumer,
				lastRequestID,
				requestStatusChannel,
				requestStatusErrorsChannel,
			)

		case requestStatus := <-requestStatusChannel:
			if requestStatus.Fulfilled == true {
				t.NumberOfFulfillments++
				wg.Done()
				return requestStatus, nil
			}

		case err := <-requestStatusErrorsChannel:
			wg.Done()
			return vrf_load_test_with_metrics.GetRequestStatus{}, err
		}
	}
}

func getRandomnessRequestStatus(
	consumer contracts.VRFv2LoadTestConsumer,
	lastRequestID *big.Int,
	requestStatusChannel chan vrf_load_test_with_metrics.GetRequestStatus,
	requestStatusErrorsChannel chan error,
) {
	requestStatus, err := consumer.GetRequestStatus(context.Background(), lastRequestID)

	if err != nil {
		requestStatusErrorsChannel <- fmt.Errorf("error occurred getting Request Status for requestID: %g", lastRequestID)
	}
	requestStatusChannel <- requestStatus
}
