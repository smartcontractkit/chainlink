package smoke

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
)

const (
	SethRootKeyIndex = 0
)

// vrfv2CleanUpFn is a cleanup function that captures pointers from context, in which it's called and uses them to clean up the test environment
var vrfv2CleanUpFn = func(
	t **testing.T,
	sethClient **seth.Client,
	config **tc.TestConfig,
	testEnv **test_env.CLClusterTestEnv,
	vrfContracts **vrfcommon.VRFContracts,
	subIDsForCancellingAfterTest *[]uint64,
	walletAddress **string,
) func() {
	return func() {
		logger := logging.GetTestLogger(*t)
		testConfig := **config
		network := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0]
		if network.Simulated {
			logger.Info().
				Str("Network Name", network.Name).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfContracts != nil && *sethClient != nil {
				if *testConfig.VRFv2.General.CancelSubsAfterTestRun {
					client := *sethClient
					var returnToAddress string
					if walletAddress == nil || *walletAddress == nil {
						returnToAddress = client.MustGetRootKeyAddress().Hex()
					} else {
						returnToAddress = **walletAddress
					}
					//cancel subs and return funds to sub owner
					vrfv2.CancelSubsAndReturnFunds(testcontext.Get(*t), *vrfContracts, returnToAddress, *subIDsForCancellingAfterTest, logger)
				}
			} else {
				logger.Error().Msg("VRF Contracts and/or Seth client are nil. Cannot execute cleanup")
			}
		}
		if !*testConfig.VRFv2.General.UseExistingEnv {
			if *testEnv == nil {
				logger.Error().Msg("Test environment is nil. Cannot execute cleanup")
				return
			}
			if err := (*testEnv).Cleanup(test_env.CleanupOpts{TestName: (*t).Name()}); err != nil {
				logger.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
}

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	configPtr := &config
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &testEnv, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	t.Run("Request Randomness", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForRequestRandomness, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subIDForRequestRandomness, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		subBalanceBeforeRequest := subscription.Balance

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForRequestRandomness,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
	t.Run("VRF Node waits block confirmation number specified by the consumer before sending fulfilment on-chain", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2.General

		consumers, subIDs, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		expectedBlockNumberWait := uint16(10)
		testConfig.MinimumConfirmations = ptr.Ptr[uint16](expectedBlockNumberWait)
		randomWordsRequestedEvent, randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subID,
			vrfKey,
			*testConfig.MinimumConfirmations,
			*testConfig.CallbackGasLimit,
			*testConfig.NumberOfWords,
			*testConfig.RandomnessRequestCountPerRequest,
			*testConfig.RandomnessRequestCountPerRequestDeviation,
			testConfig.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		// check that VRF node waited at least the number of blocks specified by the consumer in the rand request min confs field
		blockNumberWait := randomWordsRequestedEvent.Raw.BlockNumber - randomWordsFulfilledEvent.Raw.BlockNumber
		require.GreaterOrEqual(t, blockNumberWait, uint64(expectedBlockNumberWait))

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})
	t.Run("CL Node VRF Job Runs", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForJobRuns, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")

		subIDForJobRuns := subIDsForJobRuns[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForJobRuns)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subIDForJobRuns, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForJobRuns...)

		jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		_, _, err = vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForJobRuns,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))
	})
	t.Run("Direct Funding", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		wrapperContracts, wrapperSubID, err := vrfv2.SetupVRFV2WrapperEnvironment(
			testcontext.Get(t),
			sethClient,
			&configCopy,
			vrfContracts.LinkToken,
			vrfContracts.MockETHLINKFeed,
			vrfContracts.CoordinatorV2,
			vrfKey.KeyHash,
			1,
		)
		require.NoError(t, err)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, *wrapperSubID)

		wrapperConsumer := wrapperContracts.LoadTestConsumers[0]

		wrapperConsumerJuelsBalanceBeforeRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")

		wrapperSubscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceBeforeRequest := wrapperSubscription.Balance

		// Request Randomness and wait for fulfillment event
		randomWordsFulfilledEvent, err := vrfv2.DirectFundingRequestRandomnessAndWaitForFulfillment(
			l,
			wrapperConsumer,
			vrfContracts.CoordinatorV2,
			*wrapperSubID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "Error requesting randomness and waiting for fulfilment")

		// Check wrapper subscription balance
		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		wrapperSubscription, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceAfterRequest := wrapperSubscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		// Check status of randomness request within the wrapper consumer contract
		consumerStatus, err := wrapperConsumer.GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "Error getting randomness request status")
		require.True(t, consumerStatus.Fulfilled)

		// Check wrapper consumer LINK balance
		expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)
		wrapperConsumerJuelsBalanceAfterRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")
		require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)

		// Check random word count
		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
		for _, w := range consumerStatus.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		l.Info().
			Str("Consumer Balance Before Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceBeforeRequest).Link()).
			Str("Consumer Balance After Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceAfterRequest).Link()).
			Bool("Fulfilment Status", consumerStatus.Fulfilled).
			Str("Paid by Consumer Contract (Link)", (*commonassets.Link)(consumerStatus.Paid).Link()).
			Str("Paid by Coordinator Sub (Link)", (*commonassets.Link)(randomWordsFulfilledEvent.Payment).Link()).
			Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
			Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
			Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
			Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
			Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
			Msg("Random Words Fulfilment Details For Link Billing")
	})
	t.Run("Oracle Withdraw", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForOracleWithDraw, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")

		subIDForOracleWithdraw := subIDsForOracleWithDraw[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOracleWithdraw)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subIDForOracleWithdraw, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForOracleWithDraw...)

		_, fulfilledEventLink, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForOracleWithdraw,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err)

		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceLinkBeforeOracleWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		l.Info().
			Str("Returning to", sethClient.MustGetRootKeyAddress().Hex()).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfContracts.CoordinatorV2.OracleWithdraw(sethClient.MustGetRootKeyAddress(), amountToWithdrawLink)
		require.NoError(t, err, "Error withdrawing LINK from coordinator to default wallet")

		defaultWalletBalanceLinkAfterOracleWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		require.Equal(
			t,
			1,
			defaultWalletBalanceLinkAfterOracleWithdraw.Cmp(defaultWalletBalanceLinkBeforeOracleWithdraw),
			"LINK funds were not returned after oracle withdraw",
		)
	})
	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		_, subIDsForCancelling, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForCancelling := subIDsForCancelling[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subIDForCancelling, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForCancelling...)

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")

		cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		txGasUsed := new(big.Int).SetUint64(cancellationTx.Receipt.GasUsed)
		// we don't have that information for older Geth versions
		if cancellationTx.Receipt.EffectiveGasPrice == nil {
			cancellationTx.Receipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
		}
		cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTx.Receipt.EffectiveGasPrice)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Effective Gas Price", cancellationTx.Receipt.EffectiveGasPrice.String()).
			Uint64("Gas Used", cancellationTx.Receipt.GasUsed).
			Msg("Cancellation TX Receipt")

		l.Info().
			Str("Returned Subscription Amount Link", cancellationEvent.Amount.String()).
			Uint64("SubID", cancellationEvent.SubId).
			Str("Returned to", cancellationEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, cancellationEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		testWalletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedLinkActual := new(big.Int).Sub(testWalletBalanceLinkAfterSubCancelling, testWalletBalanceLinkBeforeSubCancelling)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Msg("Sub funds returned")

		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")
	})
	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		// Underfund subscription to force fulfillments to fail
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDsForOwnerCancelling, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForOwnerCancelling := subIDsForOwnerCancelling[0]
		subscriptionForCancelling, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscriptionForCancelling, strconv.FormatUint(subIDForOwnerCancelling, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForOwnerCancelling...)

		// No GetActiveSubscriptionIds function available - skipping check

		pendingRequestsExist, err := vrfContracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		// Request randomness - should fail due to underfunded subscription
		randomWordsFulfilledEventTimeout := 5 * time.Second
		_, _, err = vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForOwnerCancelling,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			randomWordsFulfilledEventTimeout,
			0,
		)
		require.Error(t, err, "Error should occur while waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfContracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfilfulled requests due to low sub balance")

		walletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForOwnerCancelling).
			Str("Returning funds to", sethClient.MustGetRootKeyAddress().Hex()).
			Msg("Canceling subscription and returning funds to subscription owner")

		// Call OwnerCancelSubscription
		cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2.OwnerCancelSubscription(subIDForOwnerCancelling)
		require.NoError(t, err, "Error canceling subscription")

		txGasUsed := new(big.Int).SetUint64(cancellationTx.Receipt.GasUsed)
		// we don't have that information for older Geth versions
		if cancellationTx.Receipt.EffectiveGasPrice == nil {
			cancellationTx.Receipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
		}
		cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTx.Receipt.EffectiveGasPrice)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Effective Gas Price", cancellationTx.Receipt.EffectiveGasPrice.String()).
			Uint64("Gas Used", cancellationTx.Receipt.GasUsed).
			Msg("Cancellation TX Receipt")

		l.Info().
			Str("Returned Subscription Amount Link", cancellationEvent.Amount.String()).
			Uint64("SubID", cancellationEvent.SubId).
			Str("Returned to", cancellationEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, cancellationEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		walletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		// Verify that subscription was deleted from Coordinator contract
		_, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		l.Info().
			Str("Expected error message", err.Error())
		require.Error(t, err, "Error did not occur when fetching deleted subscription from the Coordinator after owner cancelation")

		subFundsReturnedLinkActual := new(big.Int).Sub(walletBalanceLinkAfterSubCancelling, walletBalanceLinkBeforeSubCancelling)
		l.Info().
			Str("Wallet Balance Before Owner Cancelation", walletBalanceLinkBeforeSubCancelling.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Str("Wallet Balance After Owner Cancelation", walletBalanceLinkAfterSubCancelling.String()).
			Msg("Sub funds returned")

		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

		// Again, there is no GetActiveSubscriptionIds method on the v2 Coordinator contract, so we can't double check that the cancelled
		// subID is no longer in the list of active subs
	})
}

func TestVRFv2MultipleSendingKeys(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	if err != nil {
		t.Fatal(err)
	}
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	configPtr := &config
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &testEnv, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          2,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDsForMultipleSendingKeys, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForMultipleSendingKeys := subIDsForMultipleSendingKeys[0]
		subscriptionForMultipleSendingKeys, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForMultipleSendingKeys)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscriptionForMultipleSendingKeys, strconv.FormatUint(subIDForMultipleSendingKeys, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForMultipleSendingKeys...)

		txKeys, _, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < newEnvConfig.NumberOfTxKeysToCreate+1; i++ {
			_, randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
				l,
				consumers[0],
				vrfContracts.CoordinatorV2,
				subIDForMultipleSendingKeys,
				vrfKey,
				*configCopy.VRFv2.General.MinimumConfirmations,
				*configCopy.VRFv2.General.CallbackGasLimit,
				*configCopy.VRFv2.General.NumberOfWords,
				*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
				*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
				configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
				0,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
			fulfillmentTx, _, err := sethClient.Client.TransactionByHash(testcontext.Get(t), randomWordsFulfilledEvent.Raw.TxHash)
			require.NoError(t, err, "error getting tx from hash")
			fulfillmentTxFromAddress, err := actions.GetTxFromAddress(fulfillmentTx)
			require.NoError(t, err, "error getting tx from address")
			fulfillmentTxFromAddresses = append(fulfillmentTxFromAddresses, fulfillmentTxFromAddress)
		}
		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(fulfillmentTxFromAddresses))
		var txKeyAddresses []string
		for _, txKey := range txKeys.Data {
			txKeyAddresses = append(txKeyAddresses, txKey.Attributes.Address)
		}
		less := func(a, b string) bool { return a < b }
		equalIgnoreOrder := cmp.Diff(txKeyAddresses, fulfillmentTxFromAddresses, cmpopts.SortSlices(less)) == ""
		require.True(t, equalIgnoreOrder)
	})
}

func TestVRFOwner(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	configPtr := &config
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &testEnv, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     true,
		UseTestCoordinator:              true,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, _, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	t.Run("Request Randomness With Force-Fulfill", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDsForForceFulfill, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForForceFulfill := subIDsForForceFulfill[0]
		subscriptionForMultipleSendingKeys, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForForceFulfill)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscriptionForMultipleSendingKeys, strconv.FormatUint(subIDForForceFulfill, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForForceFulfill...)

		vrfCoordinatorOwner, err := vrfContracts.CoordinatorV2.GetOwner(testcontext.Get(t))
		require.NoError(t, err)
		require.Equal(t, vrfContracts.VRFOwner.Address(), vrfCoordinatorOwner.String())

		err = vrfContracts.LinkToken.Transfer(
			consumers[0].Address(),
			conversions.EtherToWei(big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink)),
		)
		require.NoError(t, err, "error transferring link to consumer contract")

		consumerLinkBalance, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), consumers[0].Address())
		require.NoError(t, err, "error getting consumer link balance")
		l.Info().
			Str("Balance", conversions.WeiToEther(consumerLinkBalance).String()).
			Str("Consumer", consumers[0].Address()).
			Msg("Consumer Link Balance")

		err = vrfContracts.MockETHLINKFeed.SetBlockTimestampDeduction(big.NewInt(3))
		require.NoError(t, err)

		// test and assert
		_, randFulfilledEvent, _, err := vrfv2.RequestRandomnessWithForceFulfillAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			vrfContracts.VRFOwner,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			conversions.EtherToWei(big.NewFloat(5)),
			common.HexToAddress(vrfContracts.LinkToken.Address()),
			time.Minute*2,
		)
		require.NoError(t, err, "error requesting randomness with force-fulfillment and waiting for fulfilment")
		require.Equal(t, 0, randFulfilledEvent.Payment.Cmp(big.NewInt(0)), "Forced Fulfilled Randomness's Payment should be 0")

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		coordinatorConfig, err := vrfContracts.CoordinatorV2.GetConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator config")

		coordinatorFeeConfig, err := vrfContracts.CoordinatorV2.GetFeeConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator fee config")

		coordinatorFallbackWeiPerUnitLinkConfig, err := vrfContracts.CoordinatorV2.GetFallbackWeiPerUnitLink(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator FallbackWeiPerUnitLink")

		require.Equal(t, *configCopy.VRFv2.General.StalenessSeconds, coordinatorConfig.StalenessSeconds)
		require.Equal(t, *configCopy.VRFv2.General.GasAfterPaymentCalculation, coordinatorConfig.GasAfterPaymentCalculation)
		require.Equal(t, *configCopy.VRFv2.General.MinimumConfirmations, coordinatorConfig.MinimumRequestConfirmations)
		require.Equal(t, *configCopy.VRFv2.General.FulfillmentFlatFeeLinkPPMTier1, coordinatorFeeConfig.FulfillmentFlatFeeLinkPPMTier1)
		require.Equal(t, *configCopy.VRFv2.General.ReqsForTier2, coordinatorFeeConfig.ReqsForTier2.Int64())
		require.Equal(t, *configCopy.VRFv2.General.FallbackWeiPerUnitLink, coordinatorFallbackWeiPerUnitLinkConfig.String())
	})
}

func TestVRFV2WithBHS(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	vrfv2Config := config.VRFv2
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	//decrease default span for checking blockhashes for unfulfilled requests
	vrfv2Config.General.BHSJobWaitBlocks = ptr.Ptr(2)
	vrfv2Config.General.BHSJobLookBackBlocks = ptr.Ptr(20)
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &testEnv, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		if os.Getenv("TEST_UNSKIP") != "true" {
			t.Skip("Skipped due to long execution time. Should be run on-demand on live testnet with TEST_UNSKIP=\"true\".")
		}
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		consumers, subIDsForBHS, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForBHS := subIDsForBHS[0]
		subscriptionForBHS, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForBHS)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscriptionForBHS, strconv.FormatUint(subIDForBHS, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForBHS...)

		randomWordsRequestedEvent, err := vrfv2.RequestRandomness(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForBHS,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			SethRootKeyIndex,
		)
		require.NoError(t, err, "error requesting randomness")

		vrfcommon.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2, randomWordsRequestedEvent, false, 0)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 256 blocks
		_, err = actions.WaitForBlockNumberToBe(
			testcontext.Get(t),
			randRequestBlockNumber+uint64(257),
			sethClient,
			&wg,
			nil,
			configCopy.VRFv2.General.WaitFor256BlocksTimeout.Duration,
			l,
		)
		wg.Wait()
		require.NoError(t, err)
		err = vrfv2.FundSubscriptions(big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink), vrfContracts.LinkToken, vrfContracts.CoordinatorV2, subIDsForBHS)
		require.NoError(t, err, "error funding subscriptions")

		randomWordsFulfilledEvent, err := vrfv2.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2,
			randomWordsRequestedEvent.RequestId,
			randomWordsRequestedEvent.Raw.BlockNumber,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfcommon.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2, randomWordsFulfilledEvent, false, 0)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDsForBHS, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForBHS := subIDsForBHS[0]
		subscriptionForBHS, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForBHS)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscriptionForBHS, strconv.FormatUint(subIDForBHS, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForBHS...)

		randomWordsRequestedEvent, err := vrfv2.RequestRandomness(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForBHS,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			SethRootKeyIndex,
		)
		require.NoError(t, err, "error requesting randomness")
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		_, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(
			testcontext.Get(t),
			randRequestBlockNumber+uint64(*configCopy.VRFv2.General.BHSJobWaitBlocks),
			sethClient,
			&wg,
			nil,
			time.Minute*1,
			l,
		)
		wg.Wait()
		require.NoError(t, err, "error waiting for blocknumber to be")

		metrics, err := consumers[0].GetLoadTestMetrics(testcontext.Get(t))
		require.Equal(t, 0, metrics.RequestCount.Cmp(big.NewInt(1)))
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)))
		gom := gomega.NewGomegaWithT(t)

		if !*configCopy.VRFv2.General.UseExistingEnv {
			l.Info().Msg("Checking BHS Node's transactions")
			var clNodeTxs *client.TransactionsData
			var txHash string
			gom.Eventually(func(g gomega.Gomega) {
				clNodeTxs, _, err = nodeTypeToNodeMap[vrfcommon.BHS].CLNode.API.ReadTransactions()
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting CL Node transactions")
				g.Expect(len(clNodeTxs.Data)).Should(gomega.BeNumerically("==", 1), "Expected 1 tx posted by BHS Node, but found %d", len(clNodeTxs.Data))
				txHash = clNodeTxs.Data[0].Attributes.Hash
				l.Info().
					Str("TX Hash", txHash).
					Int("Number of TXs", len(clNodeTxs.Data)).
					Msg("BHS Node txs")
			}, "2m", "1s").Should(gomega.Succeed())

			require.Equal(t, strings.ToLower(vrfContracts.BHS.Address()), strings.ToLower(clNodeTxs.Data[0].Attributes.To))

			bhsStoreTx, _, err := sethClient.Client.TransactionByHash(testcontext.Get(t), common.HexToHash(txHash))
			require.NoError(t, err, "error getting tx from hash")

			bhsStoreTxInputData, err := actions.DecodeTxInputData(blockhash_store.BlockhashStoreABI, bhsStoreTx.Data())
			require.NoError(t, err, "error decoding tx input data")
			l.Info().
				Str("Block Number", bhsStoreTxInputData["n"].(*big.Int).String()).
				Msg("BHS Node's Store Blockhash for Blocknumber Method TX")
			require.Equal(t, randRequestBlockNumber, bhsStoreTxInputData["n"].(*big.Int).Uint64())
		} else {
			l.Warn().Msg("Skipping BHS Node's transactions check as existing env is used")
		}
		var randRequestBlockHash [32]byte
		gom.Eventually(func(g gomega.Gomega) {
			randRequestBlockHash, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting blockhash for a blocknumber which was stored in BHS contract")
		}, "2m", "1s").Should(gomega.Succeed())
		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}

func TestVRFV2NodeReorg(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	if !network.Simulated {
		t.Skip("Skipped since Reorg test could only be run on Simulated chain.")
	}
	chainID := network.ChainID

	configPtr := &config
	chainlinkNodeLogScannerSettings := test_env.GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(
		testreporters.NewAllowedLogMessage(
			"Got very old block.",
			"Test is expecting a reorg to occur",
			zapcore.DPanicLevel,
			testreporters.WarnAboutAllowedMsgs_No),
		testreporters.NewAllowedLogMessage(
			"Reorg greater than finality depth detected",
			"Test is expecting a reorg to occur",
			zapcore.DPanicLevel,
			testreporters.WarnAboutAllowedMsgs_No),
	)
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: chainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, _, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2 universe")

	consumers, subIDs, err := vrfv2.SetupNewConsumersAndSubs(
		sethClient,
		vrfContracts.CoordinatorV2,
		config,
		vrfContracts.LinkToken,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	subID := subIDs[0]
	subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")
	vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

	t.Run("Reorg on fulfillment", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		configCopy.VRFv2.General.MinimumConfirmations = ptr.Ptr[uint16](10)

		//1. request randomness and wait for fulfillment for blockhash from Reorged Fork
		randomWordsRequestedEvent, randomWordsFulfilledEventOnReorgedFork, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err)

		// rewind chain to block number after the request was made, but before the request was fulfilled
		rewindChainToBlock := randomWordsRequestedEvent.Raw.BlockNumber + 1

		rpcUrl, err := vrfcommon.GetRPCUrl(env, chainID)
		require.NoError(t, err, "error getting rpc url")

		//2. rewind chain by n number of blocks - basically, mimicking reorg scenario
		latestBlockNumberAfterReorg, err := vrfcommon.RewindSimulatedChainToBlockNumber(testcontext.Get(t), sethClient, rpcUrl, rewindChainToBlock, l)
		require.NoError(t, err, fmt.Sprintf("error rewinding chain to block number %d", rewindChainToBlock))

		//3.1 ensure that chain is reorged and latest block number is greater than the block number when request was made
		require.Greater(t, latestBlockNumberAfterReorg, randomWordsRequestedEvent.Raw.BlockNumber)

		//3.2 ensure that chain is reorged and latest block number is less than the block number when fulfilment was performed
		require.Less(t, latestBlockNumberAfterReorg, randomWordsFulfilledEventOnReorgedFork.Raw.BlockNumber)

		//4. wait for the fulfillment which VRF Node will generate for Canonical chain
		_, err = vrfv2.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2,
			randomWordsRequestedEvent.RequestId,
			randomWordsRequestedEvent.Raw.BlockNumber,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)

		require.NoError(t, err, "error waiting for randomness fulfilled event")
	})

	t.Run("Reorg on rand request", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		//1. set minimum confirmations to higher value so that we can be sure that request won't be fulfilled before reorg
		configCopy.VRFv2.General.MinimumConfirmations = ptr.Ptr[uint16](6)

		//2. request randomness
		randomWordsRequestedEvent, err := vrfv2.RequestRandomness(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			SethRootKeyIndex,
		)
		require.NoError(t, err)

		// rewind chain to block number before the randomness request was made
		rewindChainToBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber - 3

		rpcUrl, err := vrfcommon.GetRPCUrl(env, chainID)
		require.NoError(t, err, "error getting rpc url")

		//3. rewind chain by n number of blocks - basically, mimicking reorg scenario
		latestBlockNumberAfterReorg, err := vrfcommon.RewindSimulatedChainToBlockNumber(testcontext.Get(t), sethClient, rpcUrl, rewindChainToBlockNumber, l)
		require.NoError(t, err, fmt.Sprintf("error rewinding chain to block number %d", rewindChainToBlockNumber))

		//4. ensure that chain is reorged and latest block number is less than the block number when request was made
		require.Less(t, latestBlockNumberAfterReorg, randomWordsRequestedEvent.Raw.BlockNumber)

		//5. ensure that rand request is not fulfilled for the request which was made on reorged fork
		// For context - when performing debug_setHead on geth simulated chain and therefore rewinding chain to a previous block,
		//then tx that was mined after reorg will not appear in canonical chain contrary to real world scenario
		//Hence, we only verify that VRF node will not generate fulfillment for the reorged fork request
		_, err = vrfv2.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2,
			randomWordsRequestedEvent.RequestId,
			randomWordsRequestedEvent.Raw.BlockNumber,
			time.Second*10,
			l,
		)
		require.Error(t, err, "fulfillment should not be generated for the request which was made on reorged fork on Simulated Chain")
	})
}

func TestVRFv2BatchFulfillmentEnabledDisabled(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	chainID := network.ChainID

	configPtr := &config
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2CleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest, nil),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2 universe")

	//batchMaxGas := config.MaxGasLimit() (2.5 mill) + 400_000 = 2.9 mill
	//callback gas limit set by consumer = 500k
	// so 4 requests should be fulfilled inside 1 tx since 500k*4 < 2.9 mill

	batchFulfilmentMaxGas := *config.VRFv2.General.MaxGasLimitCoordinatorConfig + 400_000
	config.VRFv2.General.CallbackGasLimit = ptr.Ptr(uint32(500_000))

	expectedNumberOfFulfillmentsInsideOneBatchFulfillment := (batchFulfilmentMaxGas / *config.VRFv2.General.CallbackGasLimit) - 1
	randRequestCount := expectedNumberOfFulfillmentsInsideOneBatchFulfillment

	t.Run("Batch Fulfillment Enabled", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")

		//ensure that no job present on the node
		err = actions.DeleteJobs([]*client.ChainlinkClient{vrfNode.CLNode.API})
		require.NoError(t, err)

		batchFullfillmentEnabled := true
		// create job with batch fulfillment enabled
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2.Address(),
			BatchCoordinatorAddress:       vrfContracts.BatchCoordinatorV2.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       batchFullfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2.General.VRFJobSimulationBlock,
			VRFOwnerConfig: &vrfcommon.VRFOwnerConfig{
				UseVRFOwner: false,
			},
		}

		l.Info().
			Msg("Creating VRFV2 Job with `batchFulfillmentEnabled = true`")
		job, err := vrfv2.CreateVRFV2Job(
			vrfNode.CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, "error creating job with higher timeout")
		vrfNode.Job = job

		consumers, subIDs, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		configCopy.VRFv2.General.RandomnessRequestCountPerRequest = ptr.Ptr(uint16(randRequestCount))

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		var wgAllRequestsFulfilled sync.WaitGroup
		wgAllRequestsFulfilled.Add(1)
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), consumers[0], 2*time.Minute, &wgAllRequestsFulfilled)
		require.NoError(t, err)
		wgAllRequestsFulfilled.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Request/Fulfilment Stats")

		clNodeTxs, resp, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTransactions()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		var batchFulfillmentTxs []client.TransactionData
		for _, tx := range clNodeTxs.Data {
			if common.HexToAddress(tx.Attributes.To).Cmp(common.HexToAddress(vrfContracts.BatchCoordinatorV2.Address())) == 0 {
				batchFulfillmentTxs = append(batchFulfillmentTxs, tx)
			}
		}
		// verify that all fulfillments should be inside one tx
		require.Equal(t, 1, len(batchFulfillmentTxs))

		fulfillmentTx, _, err := sethClient.Client.TransactionByHash(testcontext.Get(t), randomWordsFulfilledEvent.Raw.TxHash)
		require.NoError(t, err, "error getting tx from hash")

		fulfillmentTXToAddress := fulfillmentTx.To().String()
		l.Info().
			Str("Actual Fulfillment Tx To Address", fulfillmentTXToAddress).
			Str("BatchCoordinatorV2 Address", vrfContracts.BatchCoordinatorV2.Address()).
			Msg("Fulfillment Tx To Address should be the BatchCoordinatorV2 Address when batch fulfillment is enabled")

		// verify that VRF node sends fulfillments via BatchCoordinator contract
		require.Equal(t, vrfContracts.BatchCoordinatorV2.Address(), fulfillmentTXToAddress, "Fulfillment Tx To Address should be the BatchCoordinatorV2 Address when batch fulfillment is enabled")

		// verify that all fulfillments should be inside one tx
		// This check is disabled for live testnets since each testnet has different gas usage for similar tx
		if network.Simulated {
			fulfillmentTxReceipt, err := sethClient.Client.TransactionReceipt(testcontext.Get(t), fulfillmentTx.Hash())
			require.NoError(t, err)
			randomWordsFulfilledLogs, err := contracts.ParseRandomWordsFulfilledLogs(vrfContracts.CoordinatorV2, fulfillmentTxReceipt.Logs)
			require.NoError(t, err)
			require.Equal(t, 1, len(batchFulfillmentTxs))
			require.Equal(t, int(randRequestCount), len(randomWordsFulfilledLogs))
		}
	})
	t.Run("Batch Fulfillment Disabled", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")
		//ensure that no job present on the node
		err = actions.DeleteJobs([]*client.ChainlinkClient{vrfNode.CLNode.API})
		require.NoError(t, err)

		batchFullfillmentEnabled := false

		//create job with batchFulfillmentEnabled = false
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2.Address(),
			BatchCoordinatorAddress:       vrfContracts.BatchCoordinatorV2.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       batchFullfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2.General.VRFJobSimulationBlock,
			VRFOwnerConfig: &vrfcommon.VRFOwnerConfig{
				UseVRFOwner: false,
			},
		}

		l.Info().
			Msg("Creating VRFV2 Job with `batchFulfillmentEnabled = false`")
		job, err := vrfv2.CreateVRFV2Job(
			vrfNode.CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, "error creating job with higher timeout")
		vrfNode.Job = job

		consumers, subIDs, err := vrfv2.SetupNewConsumersAndSubs(
			sethClient,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		configCopy.VRFv2.General.RandomnessRequestCountPerRequest = ptr.Ptr(uint16(randRequestCount))

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		var wgAllRequestsFulfilled sync.WaitGroup
		wgAllRequestsFulfilled.Add(1)
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), consumers[0], 2*time.Minute, &wgAllRequestsFulfilled)
		require.NoError(t, err)
		wgAllRequestsFulfilled.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Request/Fulfilment Stats")

		fulfillmentTx, _, err := sethClient.Client.TransactionByHash(testcontext.Get(t), randomWordsFulfilledEvent.Raw.TxHash)
		require.NoError(t, err, "error getting tx from hash")

		fulfillmentTXToAddress := fulfillmentTx.To().String()
		l.Info().
			Str("Actual Fulfillment Tx To Address", fulfillmentTXToAddress).
			Str("CoordinatorV2 Address", vrfContracts.CoordinatorV2.Address()).
			Msg("Fulfillment Tx To Address should be the CoordinatorV2 Address when batch fulfillment is disabled")

		// verify that VRF node sends fulfillments via Coordinator contract
		require.Equal(t, vrfContracts.CoordinatorV2.Address(), fulfillmentTXToAddress, "Fulfillment Tx To Address should be the CoordinatorV2 Address when batch fulfillment is disabled")

		clNodeTxs, resp, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTransactions()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)

		var singleFulfillmentTxs []client.TransactionData
		for _, tx := range clNodeTxs.Data {
			if common.HexToAddress(tx.Attributes.To).Cmp(common.HexToAddress(vrfContracts.CoordinatorV2.Address())) == 0 {
				singleFulfillmentTxs = append(singleFulfillmentTxs, tx)
			}
		}
		// verify that all fulfillments should be in separate txs
		require.Equal(t, int(randRequestCount), len(singleFulfillmentTxs))
	})
}
