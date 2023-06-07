package soak

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestVRFV2Soak(t *testing.T) {
	var testInputs testsetups.VRFV2SoakTestInputs
	err := envconfig.Process("VRFV2", &testInputs)
	require.NoError(t, err, "Error reading VRFV2 soak test inputs")
	vrfSubscriptionFundingAmountInLink := testInputs.SubscriptionFunding
	chainlinkNodeFundingAmountEth := testInputs.ChainlinkNodeFunding
	randomnessRequestCountPerRequest := testInputs.RandomnessRequestCountPerRequest

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
		//todo - what should be TTL for soak test?
		time.Minute*60,
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

	vrfV2Contracts, chainlinkNodesAfterRedeployment, vrfV2jobs, _ := vrfv2_actions.SetupVRFV2Universe(
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

	consumerContract := vrfV2Contracts.LoadTestConsumer

	vrfV2SoakTest := testsetups.NewVRFV2SoakTest(&testsetups.VRFV2SoakTestInputs{
		BlockchainClient:     chainClient,
		TestDuration:         testInputs.TestDuration,
		ChainlinkNodeFunding: chainlinkNodeFundingAmountEth,
		SubscriptionFunding:  vrfSubscriptionFundingAmountInLink,
		StopTestOnError:      testInputs.StopTestOnError,
		RequestsPerMinute:    testInputs.RequestsPerMinute,
		ConsumerContract:     consumerContract,
		TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) error {
			// request randomness
			err = consumerContract.RequestRandomness(
				vrfV2jobs[0].KeyHash,
				vrfv2_constants.SubID,
				vrfv2_constants.MinimumConfirmations,
				vrfv2_constants.CallbackGasLimit,
				vrfv2_constants.NumberOfWords,
				uint16(randomnessRequestCountPerRequest),
			)
			if err != nil {
				return fmt.Errorf("error occurred Requesting Randomness, error: %w", err)
			}

			l.Info().
				Int("Request Number", requestNumber).
				Int("Randomness Request Count Per Request", randomnessRequestCountPerRequest).
				Msg("Randomness requested")

			printDebugData(l, vrfV2Contracts, chainlinkNodesAfterRedeployment)

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

func printDebugData(l zerolog.Logger, vrfV2Contracts vrfv2_actions.VRFV2Contracts, chainlinkNodesAfterRedeployment []*client.Chainlink) {
	subscription, err := vrfV2Contracts.Coordinator.GetSubscription(nil, vrfv2_constants.SubID)
	if err != nil {
		l.Error().Err(err).
			Uint64("Subscription ID", vrfv2_constants.SubID).
			Interface("Coordinator Address", vrfV2Contracts.Coordinator.Address()).
			Msg("error occurred Getting Subscription Data from a Coordinator Contract")
	}
	l.Debug().Interface("Data", subscription).Uint64("Subscription ID", vrfv2_constants.SubID).Msg("Subscription Data")
	remainingSubBalanceInLink := new(big.Float).Quo(new(big.Float).SetInt(subscription.Balance), big.NewFloat(1e18))
	l.Debug().Interface("Balance", remainingSubBalanceInLink).Msg("Remaining Balance in Link for a subscription")
	nativeTokenPrimaryKey, err := chainlinkNodesAfterRedeployment[0].ReadPrimaryETHKey()
	if err != nil {
		l.Error().Err(err).Msg("error occurred reading Native Token Primary Key from Chainlink Node")
	}
	ethBalance, ok := new(big.Int).SetString(nativeTokenPrimaryKey.Attributes.ETHBalance, 10)
	if !ok {
		l.Error().Interface("Balance", nativeTokenPrimaryKey.Attributes.ETHBalance).Msg("error occurred converting Native Token Primary Key from Chainlink Node")
	}
	remainingNativeTokenPrimaryKeyBalanceInETH := new(big.Float).Quo(new(big.Float).SetInt(ethBalance), big.NewFloat(1e18))
	l.Debug().
		Interface("Balance", remainingNativeTokenPrimaryKeyBalanceInETH).
		Interface("Key Address", nativeTokenPrimaryKey.Attributes.Address).
		Msg("Remaining Balance for a Native Token Primary Key of Chainlink Node")
}
