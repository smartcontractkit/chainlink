package smoke

import (
	"fmt"
	"math"
	"math/big"
	"os"
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

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
)

// vrfv2PlusCleanUpFn is a cleanup function that captures pointers from context, in which it's called and uses them to clean up the test environment
var vrfv2PlusCleanUpFn = func(
	t **testing.T,
	sethClient **seth.Client,
	config **tc.TestConfig,
	testEnv **test_env.CLClusterTestEnv,
	vrfContracts **vrfcommon.VRFContracts,
	subIDsForCancellingAfterTest *[]*big.Int,
) func() {
	return func() {
		l := logging.GetTestLogger(*t)
		testConfig := **config
		network := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0]
		if network.Simulated {
			l.Info().
				Str("Network Name", network.Name).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfContracts != nil && *sethClient != nil {
				if *testConfig.VRFv2Plus.General.CancelSubsAfterTestRun {
					client := *sethClient
					returnToAddress := client.MustGetRootKeyAddress().Hex()
					//cancel subs and return funds to sub owner
					vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(*t), *vrfContracts, returnToAddress, *subIDsForCancellingAfterTest, l)
				}
			} else {
				l.Error().Msg("VRF Contracts and/or Seth client are nil. Cannot execute cleanup")
			}
		}
		if !*testConfig.VRFv2Plus.General.UseExistingEnv {
			if *testEnv == nil {
				l.Error().Msg("Test environment is nil. Cannot execute cleanup")
				return
			}
			if err := (*testEnv).Cleanup(test_env.CleanupOpts{TestName: (*t).Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
}

func TestVRFv2Plus(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")

	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2Plus universe")

	t.Run("Link Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false

		consumers, subIDsForRequestRandomness, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subIDForRequestRandomness.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		subBalanceBeforeRequest := subscription.Balance

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subIDForRequestRandomness,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		require.False(t, randomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, randomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, int(*configCopy.VRFv2Plus.General.NumberOfWords), len(status.RandomWords))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
		t.Run("Verify Billing", func(t *testing.T) {
			actualSubPaymentJuels := randomWordsFulfilledEvent.Payment

			expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, actualSubPaymentJuels)
			subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := subscription.Balance
			require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

			//todo - need to refactor the test so that when running on live testnet and deploying a new environment,
			// we would load real Link Token and Link/ETH feed addresses from the config
		})

	})
	t.Run("Native Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2Plus.General
		var isNativeBilling = true

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		subNativeTokenBalanceBeforeRequest := subscription.NativeBalance

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		require.False(t, randomWordsFulfilledEvent.OnlyPremium)
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment)
		require.True(t, randomWordsFulfilledEvent.Success)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, int(*testConfig.NumberOfWords), len(status.RandomWords))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
		t.Run("Verify Billing", func(t *testing.T) {
			actualSubPaymentWei := randomWordsFulfilledEvent.Payment
			expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, actualSubPaymentWei)
			subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
			require.NoError(t, err)
			subBalanceAfterRequest := subscription.NativeBalance
			require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

			// verify that the actual sub payment is within the expected range - this cannot be checked on SIMULATED env
			coordinatorConfig, err := vrfContracts.CoordinatorV2Plus.GetConfig(testcontext.Get(t))
			require.NoError(t, err, "error getting coordinator config")
			//check that Coordinator config is correct
			require.Equal(t, *configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM, coordinatorConfig.FulfillmentFlatFeeNativePPM)
			require.Equal(t, *configCopy.VRFv2Plus.General.NativePremiumPercentage, coordinatorConfig.NativePremiumPercentage)
			require.Equal(t, *configCopy.VRFv2Plus.General.GasAfterPaymentCalculation, coordinatorConfig.GasAfterPaymentCalculation)
			// check that the actual sub payment is within the expected range and is according to the coordinator config's billing settings
			fulfillmentTxHash := randomWordsFulfilledEvent.Raw.TxHash
			fulfillmentTxReceipt, err := sethClient.Client.TransactionReceipt(testcontext.Get(t), fulfillmentTxHash)
			require.NoError(t, err, "error getting fulfilment tx receipt")

			txGasUsed := new(big.Int).SetUint64(fulfillmentTxReceipt.GasUsed)
			// we don't have that information for older Geth versions
			if fulfillmentTxReceipt.EffectiveGasPrice == nil {
				fulfillmentTxReceipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
			}
			fulfillmentTxFeeWei := new(big.Int).Mul(txGasUsed, fulfillmentTxReceipt.EffectiveGasPrice)

			premiumFeeInDecimal := new(big.Float).Quo(big.NewFloat(float64(*configCopy.VRFv2Plus.General.CoordinatorNativePremiumPercentage)), big.NewFloat(100))
			premiumFeeRatio := new(big.Float).Add(premiumFeeInDecimal, big.NewFloat(1))
			// since - tx fee * (premium fee in decimal + 1) = expected sub payment
			fulfillmentTxFeeWeiFloat64, _ := fulfillmentTxFeeWei.Float64()
			expectedSubPaymentWeiWithoutFlatFee := new(big.Float).Mul(big.NewFloat(fulfillmentTxFeeWeiFloat64), premiumFeeRatio)
			flatFee := new(big.Float).Mul(big.NewFloat(float64(*configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM)), big.NewFloat(1e12))
			expectedSubPaymentWeiWitFlatFee := new(big.Float).Add(expectedSubPaymentWeiWithoutFlatFee, flatFee)
			vrfv2plus.LogPaymentDetails(l, fulfillmentTxFeeWei, fulfillmentTxReceipt, actualSubPaymentWei, expectedSubPaymentWeiWitFlatFee, configCopy)
			actualSubPaymentWeiFloat, _ := actualSubPaymentWei.Float64()
			expectedSubPaymentWeiFloat, _ := expectedSubPaymentWeiWitFlatFee.Float64()

			//verify that the actual sub payment is more than TX fee
			require.Equal(t, 1, actualSubPaymentWei.Cmp(fulfillmentTxFeeWei), "the actual sub payment has to be more than the TX fee")

			tolerance := *testConfig.SubBillingTolerance
			isWithinTolerance, diff := actions.WithinTolerance(actualSubPaymentWeiFloat, expectedSubPaymentWeiFloat, tolerance)
			require.True(t, isWithinTolerance, fmt.Sprintf("Expected the actual sub payment to be within %f tolerance of the expected sub payment. Diff: %f", tolerance, diff))
		})
	})
	t.Run("Direct Funding", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperUniverse(
			testcontext.Get(t),
			sethClient,
			vrfContracts,
			&configCopy,
			vrfKey.KeyHash,
			1,
			l,
		)
		require.NoError(t, err)

		t.Run("Link Billing", func(t *testing.T) {
			configCopy := config.MustCopy().(tc.TestConfig)
			testConfig := configCopy.VRFv2Plus.General
			var isNativeBilling = false

			wrapperConsumerJuelsBalanceBeforeRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.WrapperConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.Balance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.WrapperConsumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				wrapperSubID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.Balance
			require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.WrapperConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)

			wrapperConsumerJuelsBalanceAfterRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.WrapperConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfcommon.LogFulfillmentDetailsLinkBilling(l, wrapperConsumerJuelsBalanceBeforeRequest, wrapperConsumerJuelsBalanceAfterRequest, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, int(*testConfig.NumberOfWords), len(consumerStatus.RandomWords))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
		t.Run("Native Billing", func(t *testing.T) {
			configCopy := config.MustCopy().(tc.TestConfig)
			testConfig := configCopy.VRFv2Plus.General
			var isNativeBilling = true

			wrapperConsumerBalanceBeforeRequestWei, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.WrapperConsumers[0].Address()), nil)
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.NativeBalance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.WrapperConsumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				wrapperSubID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceWei := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.NativeBalance
			require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.WrapperConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerWeiBalance := new(big.Int).Sub(wrapperConsumerBalanceBeforeRequestWei, consumerStatus.Paid)

			wrapperConsumerBalanceAfterRequestWei, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.WrapperConsumers[0].Address()), nil)
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerWeiBalance, wrapperConsumerBalanceAfterRequestWei)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfcommon.LogFulfillmentDetailsNativeBilling(l, wrapperConsumerBalanceBeforeRequestWei, wrapperConsumerBalanceAfterRequestWei, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, int(*testConfig.NumberOfWords), len(consumerStatus.RandomWords))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
	})
	t.Run("VRF Node waits block confirmation number specified by the consumer before sending fulfilment on-chain", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2Plus.General
		var isNativeBilling = true

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		expectedBlockNumberWait := uint16(10)
		testConfig.MinimumConfirmations = ptr.Ptr[uint16](expectedBlockNumberWait)
		randomWordsRequestedEvent, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			testConfig,
			l,
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
		var isNativeBilling = false

		consumers, subIDsForRequestRandomness, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subIDForRequestRandomness.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subIDForRequestRandomness,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))
	})
	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		_, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceNativeBeforeSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), testWalletAddress, nil)
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")

		cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, testWalletAddress)
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
			Str("Returned Subscription Amount Native", cancellationEvent.AmountLink.String()).
			Str("Returned Subscription Amount Link", cancellationEvent.AmountLink.String()).
			Str("SubID", cancellationEvent.SubId.String()).
			Str("Returned to", cancellationEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceNative, cancellationEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
		require.Equal(t, subBalanceLink, cancellationEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		testWalletBalanceNativeAfterSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), testWalletAddress, nil)
		require.NoError(t, err)

		testWalletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedNativeActual := new(big.Int).Sub(testWalletBalanceNativeAfterSubCancelling, testWalletBalanceNativeBeforeSubCancelling)
		subFundsReturnedLinkActual := new(big.Int).Sub(testWalletBalanceLinkAfterSubCancelling, testWalletBalanceLinkBeforeSubCancelling)

		subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
		deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
		l.Info().
			Str("Sub Balance - Native", subBalanceNative.String()).
			Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
			Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Msg("Sub funds returned")

		//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
		//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
		require.Equal(t, 1, testWalletBalanceNativeAfterSubCancelling.Cmp(testWalletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")
		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

	})
	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2Plus.General

		//underfund subs in order rand fulfillments to fail
		testConfig.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))
		testConfig.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
		activeSubscriptionIdsBeforeSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err)

		require.True(t, it_utils.BigIntSliceContains(activeSubscriptionIdsBeforeSubCancellation, subID))

		pendingRequestsExist, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout = ptr.Ptr(blockchain.StrDuration{Duration: 5 * time.Second})
		_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfulfilled rand requests due to low sub balance")

		walletBalanceNativeBeforeSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
		require.NoError(t, err)

		walletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", sethClient.MustGetRootKeyAddress().Hex()).
			Msg("Canceling subscription and returning funds to subscription owner")

		cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2Plus.OwnerCancelSubscription(subID)
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
			Str("Returned Subscription Amount Native", cancellationEvent.AmountNative.String()).
			Str("Returned Subscription Amount Link", cancellationEvent.AmountLink.String()).
			Str("SubID", cancellationEvent.SubId.String()).
			Str("Returned to", cancellationEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceNative, cancellationEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
		require.Equal(t, subBalanceLink, cancellationEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		walletBalanceNativeAfterSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
		require.NoError(t, err)

		walletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedNativeActual := new(big.Int).Sub(walletBalanceNativeAfterSubCancelling, walletBalanceNativeBeforeSubCancelling)
		subFundsReturnedLinkActual := new(big.Int).Sub(walletBalanceLinkAfterSubCancelling, walletBalanceLinkBeforeSubCancelling)

		subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
		deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
		l.Info().
			Str("Sub Balance - Native", subBalanceNative.String()).
			Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
			Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Str("walletBalanceNativeBeforeSubCancelling", walletBalanceNativeBeforeSubCancelling.String()).
			Str("walletBalanceNativeAfterSubCancelling", walletBalanceNativeAfterSubCancelling.String()).
			Msg("Sub funds returned")

		//todo - need to use different wallet for each test to verify exact amount of Native/LINK returned
		//todo - as defaultWallet is used in other tests in parallel which might affect the balance - TT-684
		//require.Equal(t, 1, walletBalanceNativeAfterSubCancelling.Cmp(walletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")

		//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
		//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

		activeSubscriptionIdsAfterSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error getting active subscription ids")

		require.False(
			t,
			it_utils.BigIntSliceContains(activeSubscriptionIdsAfterSubCancellation, subID),
			"Active subscription ids should not contain sub id after sub cancellation",
		)
	})
	t.Run("Owner Withdraw", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		_, fulfilledEventLink, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err)

		_, fulfilledEventNative, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err)
		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceNativeBeforeWithdraw, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
		require.NoError(t, err)

		defaultWalletBalanceLinkBeforeWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		l.Info().
			Str("Returning to", sethClient.MustGetRootKeyAddress().Hex()).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfContracts.CoordinatorV2Plus.Withdraw(
			common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()),
		)
		require.NoError(t, err, "error withdrawing LINK from coordinator to default wallet")
		amountToWithdrawNative := fulfilledEventNative.Payment

		l.Info().
			Str("Returning to", sethClient.MustGetRootKeyAddress().Hex()).
			Str("Amount", amountToWithdrawNative.String()).
			Msg("Invoking Oracle Withdraw for Native")

		err = vrfContracts.CoordinatorV2Plus.WithdrawNative(
			common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()),
		)
		require.NoError(t, err, "error withdrawing Native tokens from coordinator to default wallet")

		defaultWalletBalanceNativeAfterWithdraw, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
		require.NoError(t, err)

		defaultWalletBalanceLinkAfterWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
		require.NoError(t, err)

		//not possible to verify exact amount of Native/LINK returned as defaultWallet is used in other tests in parallel which might affect the balance
		require.Equal(t, 1, defaultWalletBalanceNativeAfterWithdraw.Cmp(defaultWalletBalanceNativeBeforeWithdraw), "Native funds were not returned after oracle withdraw native")
		require.Equal(t, 1, defaultWalletBalanceLinkAfterWithdraw.Cmp(defaultWalletBalanceLinkBeforeWithdraw), "LINK funds were not returned after oracle withdraw")
	})
}

func TestVRFv2PlusMultipleSendingKeys(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          2,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = true

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		txKeys, _, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < newEnvConfig.NumberOfTxKeysToCreate+1; i++ {
			_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
				consumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				subID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
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

func TestVRFv2PlusMigration(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	// Migrate subscription from old coordinator to new coordinator, verify if balances
	// are moved correctly and requests can be made successfully in the subscription in
	// new coordinator
	t.Run("Test migration of Subscription Billing subID", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			2,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		require.Equal(t, subID, activeSubIdsOldCoordinatorBeforeMigration[0])

		oldSubscriptionBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := contracts.DeployVRFCoordinatorV2PlusUpgradedVersion(sethClient, vrfContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfKey.VRFKey, newCoordinator, assets.GWei(*configCopy.VRFv2Plus.General.CLNodeMaxGasPriceGWei).ToInt().Uint64())
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		err = newCoordinator.SetConfig(
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.MaxGasLimitCoordinatorConfig,
			*configCopy.VRFv2Plus.General.StalenessSeconds,
			*configCopy.VRFv2Plus.General.GasAfterPaymentCalculation,
			big.NewInt(*configCopy.VRFv2Plus.General.LinkNativeFeedResponse),
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM,
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeLinkDiscountPPM,
			*configCopy.VRFv2Plus.General.NativePremiumPercentage,
			*configCopy.VRFv2Plus.General.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(vrfContracts.LinkToken.Address(), vrfContracts.MockETHLINKFeed.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.VRFKey.Data.ID,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       false,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		_, migrationCompletedEvent, err := vrfContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())
		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfContracts.CoordinatorV2Plus)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		migratedSubscription, err := newCoordinator.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

		//Verify that Coordinators were updated in Consumers
		for _, consumer := range consumers {
			coordinatorAddressInConsumerAfterMigration, err := consumer.GetCoordinator(testcontext.Get(t))
			require.NoError(t, err, "error getting Coordinator from Consumer contract")
			require.Equal(t, newCoordinator.Address(), coordinatorAddressInConsumerAfterMigration.String())
			l.Info().
				Str("Consumer", consumer.Address()).
				Str("Coordinator", coordinatorAddressInConsumerAfterMigration.String()).
				Msg("Coordinator Address in Consumer After Migration")
		}

		//Verify old and migrated subs
		require.Equal(t, oldSubscriptionBeforeMigration.NativeBalance, migratedSubscription.NativeBalance)
		require.Equal(t, oldSubscriptionBeforeMigration.Balance, migratedSubscription.Balance)
		require.Equal(t, oldSubscriptionBeforeMigration.SubOwner, migratedSubscription.SubOwner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		// If (subscription billing), numActiveSub should be 0 after migration in oldCoordinator
		require.Error(t, err, "error not occurred getting active sub ids. Should occur since it should revert when sub id array is empty")

		activeSubIdsMigratedCoordinator, err := newCoordinator.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsMigratedCoordinator, 1, "Active Sub Ids length is not equal to 1 for Migrated Coordinator after migration")
		require.Equal(t, subID, activeSubIdsMigratedCoordinator[0])

		//Verify that total balances changed for Link and Eth for new and old coordinator
		expectedLinkTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.Balance, migratedCoordinatorLinkTotalBalanceBeforeMigration)
		expectedEthTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.NativeBalance, migratedCoordinatorEthTotalBalanceBeforeMigration)

		expectedLinkTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorLinkTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.Balance)
		expectedEthTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorEthTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.NativeBalance)
		require.Equal(t, 0, expectedLinkTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorEthTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedLinkTotalBalanceForOldCoordinator.Cmp(oldCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForOldCoordinator.Cmp(oldCoordinatorEthTotalBalanceAfterMigration))

		//Verify rand requests fulfills with Link Token billing
		_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			newCoordinator,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		//Verify rand requests fulfills with Native Token billing
		_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[1],
			newCoordinator,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	})

	// Migrate wrapper subscription from old coordinator to new coordinator, verify if balances
	// are moved correctly and requests can be made successfully in the subscription in
	// new coordinator
	t.Run("Test migration of direct billing using VRFV2PlusWrapper subID", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperUniverse(
			testcontext.Get(t),
			sethClient,
			vrfContracts,
			&configCopy,
			vrfKey.KeyHash,
			1,
			l,
		)
		require.NoError(t, err)
		subID := wrapperSubID

		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		activeSubID := activeSubIdsOldCoordinatorBeforeMigration[0]
		require.Equal(t, subID, activeSubID)

		oldSubscriptionBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := contracts.DeployVRFCoordinatorV2PlusUpgradedVersion(sethClient, vrfContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfKey.VRFKey, newCoordinator, assets.GWei(*configCopy.VRFv2Plus.General.CLNodeMaxGasPriceGWei).ToInt().Uint64())
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		err = newCoordinator.SetConfig(
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.MaxGasLimitCoordinatorConfig,
			*configCopy.VRFv2Plus.General.StalenessSeconds,
			*configCopy.VRFv2Plus.General.GasAfterPaymentCalculation,
			big.NewInt(*configCopy.VRFv2Plus.General.LinkNativeFeedResponse),
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM,
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeLinkDiscountPPM,
			*configCopy.VRFv2Plus.General.NativePremiumPercentage,
			*configCopy.VRFv2Plus.General.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(vrfContracts.LinkToken.Address(), vrfContracts.MockETHLINKFeed.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.VRFKey.Data.ID,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       false,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		// Migrate wrapper's sub using coordinator's migrate method
		_, migrationCompletedEvent, err := vrfContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())
		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfContracts.CoordinatorV2Plus)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		migratedSubscription, err := newCoordinator.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

		// Verify that Coordinators were updated in Consumers- Consumer in this case is the VRFV2PlusWrapper
		coordinatorAddressInConsumerAfterMigration, err := wrapperContracts.VRFV2PlusWrapper.Coordinator(testcontext.Get(t))
		require.NoError(t, err, "error getting Coordinator from Consumer contract- VRFV2PlusWrapper")
		require.Equal(t, newCoordinator.Address(), coordinatorAddressInConsumerAfterMigration.String())
		l.Info().
			Str("Consumer-VRFV2PlusWrapper", wrapperContracts.VRFV2PlusWrapper.Address()).
			Str("Coordinator", coordinatorAddressInConsumerAfterMigration.String()).
			Msg("Coordinator Address in VRFV2PlusWrapper After Migration")

		//Verify old and migrated subs
		require.Equal(t, oldSubscriptionBeforeMigration.NativeBalance, migratedSubscription.NativeBalance)
		require.Equal(t, oldSubscriptionBeforeMigration.Balance, migratedSubscription.Balance)
		require.Equal(t, oldSubscriptionBeforeMigration.SubOwner, migratedSubscription.SubOwner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		// If (subscription billing) or (direct billing and numActiveSubs is 0 before this test) -> numActiveSub should be 0 after migration in oldCoordinator
		require.Error(t, err, "error not occurred getting active sub ids. Should occur since it should revert when sub id array is empty")

		activeSubIdsMigratedCoordinator, err := newCoordinator.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsMigratedCoordinator, 1, "Active Sub Ids length is not equal to 1 for Migrated Coordinator after migration")
		require.Equal(t, subID, activeSubIdsMigratedCoordinator[0])

		//Verify that total balances changed for Link and Eth for new and old coordinator
		expectedLinkTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.Balance, migratedCoordinatorLinkTotalBalanceBeforeMigration)
		expectedEthTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.NativeBalance, migratedCoordinatorEthTotalBalanceBeforeMigration)

		expectedLinkTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorLinkTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.Balance)
		expectedEthTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorEthTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.NativeBalance)
		require.Equal(t, 0, expectedLinkTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorEthTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedLinkTotalBalanceForOldCoordinator.Cmp(oldCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForOldCoordinator.Cmp(oldCoordinatorEthTotalBalanceAfterMigration))

		// Verify rand requests fulfills with Link Token billing
		isNativeBilling := false
		randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
			wrapperContracts.WrapperConsumers[0],
			newCoordinator,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		consumerStatus, err := wrapperContracts.WrapperConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, consumerStatus.Fulfilled)

		// Verify rand requests fulfills with Native Token billing
		isNativeBilling = true
		randomWordsFulfilledEvent, err = vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
			wrapperContracts.WrapperConsumers[0],
			newCoordinator,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		consumerStatus, err = wrapperContracts.WrapperConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, consumerStatus.Fulfilled)
	})
}

func TestVRFV2PlusWithBHS(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	//decrease default span for checking blockhashes for unfulfilled requests
	vrfv2PlusConfig.General.BHSJobWaitBlocks = ptr.Ptr(2)
	vrfv2PlusConfig.General.BHSJobLookBackBlocks = ptr.Ptr(20)
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	var isNativeBilling = true
	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		if os.Getenv("TEST_UNSKIP") != "true" {
			t.Skip("Skipped due to long execution time. Should be run on-demand on live testnet with TEST_UNSKIP=\"true\".")
		}
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		randomWordsRequestedEvent, err := vrfv2plus.RequestRandomness(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness")

		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)

		const waitForNumberOfBlocks = 257
		desiredBlockNumberReached := make(chan bool)
		go func() {
			//Wait at least 256 blocks
			_, err = actions.WaitForBlockNumberToBe(
				testcontext.Get(t),
				randRequestBlockNumber+waitForNumberOfBlocks,
				sethClient,
				&wg,
				desiredBlockNumberReached,
				configCopy.VRFv2Plus.General.WaitFor256BlocksTimeout.Duration,
				l,
			)
			require.NoError(t, err)
		}()

		if *configCopy.VRFv2Plus.General.GenerateTXsOnChain {
			wg.Add(1)
			go func() {
				_, err := actions.ContinuouslyGenerateTXsOnChain(sethClient, desiredBlockNumberReached, &wg, l)
				require.NoError(t, err)
			}()
		}
		wg.Wait()

		err = vrfv2plus.FundSubscriptions(
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			subIDs,
			*configCopy.VRFv2Plus.General.SubscriptionBillingType,
		)
		require.NoError(t, err, "error funding subscriptions")
		randomWordsFulfilledEvent, err := vrfv2plus.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2Plus,
			randomWordsRequestedEvent.RequestId,
			subID,
			randomWordsRequestedEvent.Raw.BlockNumber,
			isNativeBilling,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
			0,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfcommon.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsFulfilledEvent, isNativeBilling, 0)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		randRequestBlockHash, err := vrfContracts.BHS.GetBlockHash(testcontext.Get(t), new(big.Int).SetUint64(randRequestBlockNumber))
		require.NoError(t, err, "error getting blockhash for a blocknumber which was stored in BHS contract")

		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Str("BHS Contract", vrfContracts.BHS.Address()).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		randomWordsRequestedEvent, err := vrfv2plus.RequestRandomness(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness")
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		_, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), new(big.Int).SetUint64(randRequestBlockNumber))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		if *configCopy.VRFv2Plus.General.BHSJobWaitBlocks < 0 {
			t.Fatalf("negative job wait blocks: %d", *configCopy.VRFv2Plus.General.BHSJobWaitBlocks)
		}
		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(
			testcontext.Get(t),
			randRequestBlockNumber+uint64(*configCopy.VRFv2Plus.General.BHSJobWaitBlocks)+10, //nolint:gosec // G115 false positive
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

		if !*configCopy.VRFv2Plus.General.UseExistingEnv {
			l.Info().Msg("Checking BHS Node's transactions")
			var clNodeTxs *nodeclient.TransactionsData
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
			randRequestBlockHash, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), new(big.Int).SetUint64(randRequestBlockNumber))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting blockhash for a blocknumber which was stored in BHS contract")
		}, "2m", "1s").Should(gomega.Succeed())
		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}

func TestVRFV2PlusWithBHF(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	// BHF job config - NOTE - not possible to create BHF job spec with waitBlocks being less than 256 blocks
	config.VRFv2Plus.General.BHFJobWaitBlocks = ptr.Ptr(260)
	config.VRFv2Plus.General.BHFJobLookBackBlocks = ptr.Ptr(500)
	config.VRFv2Plus.General.BHFJobPollPeriod = ptr.Ptr(blockchain.StrDuration{Duration: time.Second * 30})
	config.VRFv2Plus.General.BHFJobRunTimeout = ptr.Ptr(blockchain.StrDuration{Duration: time.Minute * 24})
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	chainlinkNodeLogScannerSettings := test_env.GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(testreporters.NewAllowedLogMessage(
		"Pipeline error",
		"Test is expecting this error to occur",
		zapcore.DPanicLevel,
		testreporters.WarnAboutAllowedMsgs_No))
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: chainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(
		testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err)

	var isNativeBilling = true
	t.Run("BHF Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		// Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		randomWordsRequestedEvent, err := vrfv2plus.RequestRandomness(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness")

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
			configCopy.VRFv2Plus.General.WaitFor256BlocksTimeout.Duration,
			l,
		)
		wg.Wait()
		require.NoError(t, err)
		l.Info().Float64("SubscriptionFundingAmountNative", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative).
			Float64("SubscriptionFundingAmountLink", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink).
			Msg("Funding subscription")
		err = vrfv2plus.FundSubscriptions(
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			subIDs,
			*configCopy.VRFv2Plus.General.SubscriptionBillingType,
		)
		require.NoError(t, err, "error funding subscriptions")

		randomWordsFulfilledEvent, err := vrfv2plus.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2Plus,
			randomWordsRequestedEvent.RequestId,
			subID,
			randomWordsRequestedEvent.Raw.BlockNumber,
			isNativeBilling,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
			0,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfcommon.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsFulfilledEvent, isNativeBilling, 0)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		clNodeTxs, _, err := nodeTypeToNodeMap[vrfcommon.BHF].CLNode.API.ReadTransactions()
		require.NoError(t, err, "error fetching txns from BHF node")
		batchBHSTxFound := false
		for _, tx := range clNodeTxs.Data {
			if strings.EqualFold(tx.Attributes.To, vrfContracts.BatchBHS.Address()) {
				batchBHSTxFound = true
			}
		}
		require.True(t, batchBHSTxFound)

		randRequestBlockHash, err := vrfContracts.BHS.GetBlockHash(testcontext.Get(t), new(big.Int).SetUint64(randRequestBlockNumber))
		require.NoError(t, err, "error getting blockhash for a blocknumber which was stored in BHS contract")

		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}

func TestVRFv2PlusReplayAfterTimeout(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	// 1. Add job spec with requestTimeout = 5 seconds
	timeout := time.Second * 5
	config.VRFv2Plus.General.VRFJobRequestTimeout = ptr.Ptr(blockchain.StrDuration{Duration: timeout})
	config.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
	config.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	t.Run("Timed out request fulfilled after node restart with replay", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			2,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		// 2. create request but without fulfilment - e.g. simulation failure (insufficient balance in the sub, )
		initialReqRandomWordsRequestedEvent, err := vrfv2plus.RequestRandomness(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for requested event")

		// 3. wait for the request timeout (1s more) duration
		time.Sleep(timeout + 1*time.Second)

		fundingLinkAmt := big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink)
		fundingNativeAmt := big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative)
		// 4. fund sub so that node can fulfill request
		err = vrfv2plus.FundSubscriptions(
			fundingLinkAmt,
			fundingNativeAmt,
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			[]*big.Int{subID},
			*configCopy.VRFv2Plus.General.SubscriptionBillingType,
		)
		require.NoError(t, err, "error funding subs after request timeout")

		// 5. no fulfilment should happen since timeout+1 seconds passed in the job
		pendingReqExists, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err, "error fetching PendingRequestsExist from coordinator")
		require.True(t, pendingReqExists, "pendingRequest must exist since subID was underfunded till request timeout")

		// 6. remove job and add new job with requestTimeout = 1 hour
		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")
		resp, err := vrfNode.CLNode.API.DeleteJob(vrfNode.Job.Data.ID)
		require.NoError(t, err, "error deleting job after timeout")
		require.Equal(t, resp.StatusCode, 204)

		configCopy.VRFv2Plus.General.VRFJobRequestTimeout = ptr.Ptr(blockchain.StrDuration{Duration: time.Duration(time.Hour * 1)})
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2Plus.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       false,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
			VRFOwnerConfig:                nil,
		}

		go func() {
			l.Info().
				Msg("Creating VRFV2 Plus Job with higher timeout (1hr)")
			job, err := vrfv2plus.CreateVRFV2PlusJob(
				vrfNode.CLNode.API,
				vrfJobSpecConfig,
			)
			require.NoError(t, err, "error creating job with higher timeout")
			vrfNode.Job = job
		}()

		// 7. Check if initial req in underfunded sub is fulfilled now, since it has been topped up and timeout increased
		l.Info().Str("reqID", initialReqRandomWordsRequestedEvent.RequestId.String()).
			Str("subID", subID.String()).
			Msg("Waiting for initalReqRandomWordsFulfilledEvent")
		initalReqRandomWordsFulfilledEvent, err := vrfv2plus.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2Plus,
			initialReqRandomWordsRequestedEvent.RequestId,
			subID,
			initialReqRandomWordsRequestedEvent.Raw.BlockNumber,
			isNativeBilling,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
			0,
		)
		require.NoError(t, err, "error waiting for initial request RandomWordsFulfilledEvent")

		require.NoError(t, err, "error waiting for fulfilment of old req")
		require.False(t, initalReqRandomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, initalReqRandomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, initalReqRandomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		// Get request status
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), initalReqRandomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})
}

func TestVRFv2PlusPendingBlockSimulationAndZeroConfirmationDelays(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}

	// override config with minConf = 0 and use pending block for simulation
	config.VRFv2Plus.General.MinimumConfirmations = ptr.Ptr[uint16](0)
	config.VRFv2Plus.General.VRFJobSimulationBlock = ptr.Ptr[string]("pending")

	env, vrfContracts, vrfKey, _, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
		testcontext.Get(t),
		sethClient,
		vrfContracts.CoordinatorV2Plus,
		config,
		vrfContracts.LinkToken,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	subID := subIDs[0]
	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")
	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

	var isNativeBilling = true

	l.Info().Uint16("minimumConfirmationDelay", *config.VRFv2Plus.General.MinimumConfirmations).Msg("Minimum Confirmation Delay")

	// test and assert
	_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		consumers[0],
		vrfContracts.CoordinatorV2Plus,
		vrfKey,
		subID,
		isNativeBilling,
		config.VRFv2Plus.General,
		l,
		0,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

	status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	require.NoError(t, err, "error getting rand request status")
	require.True(t, status.Fulfilled)
	l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
}

func TestVRFv2PlusNodeReorg(t *testing.T) {
	t.Skip("Flakey", "https://smartcontract-it.atlassian.net/browse/DEVSVCS-829")
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	if !network.Simulated {
		t.Skip("Skipped since Reorg test could only be run on Simulated chain.")
	}
	chainID := network.ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
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
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: chainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, _, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2Plus universe")

	var isNativeBilling = true
	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
		testcontext.Get(t),
		sethClient,
		vrfContracts.CoordinatorV2Plus,
		config,
		vrfContracts.LinkToken,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	subID := subIDs[0]
	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")
	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

	t.Run("Reorg on fulfillment", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		configCopy.VRFv2Plus.General.MinimumConfirmations = ptr.Ptr[uint16](10)

		//1. request randomness and wait for fulfillment for blockhash from Reorged Fork
		randomWordsRequestedEvent, randomWordsFulfilledEventOnReorgedFork, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
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
		_, err = vrfv2plus.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2Plus,
			randomWordsRequestedEvent.RequestId,
			subID,
			randomWordsRequestedEvent.Raw.BlockNumber,
			isNativeBilling,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
			0,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
	})

	t.Run("Reorg on rand request", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		//1. set minimum confirmations to higher value so that we can be sure that request won't be fulfilled before reorg
		configCopy.VRFv2Plus.General.MinimumConfirmations = ptr.Ptr[uint16](6)

		//2. request randomness
		randomWordsRequestedEvent, err := vrfv2plus.RequestRandomness(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
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
		_, err = vrfv2plus.WaitRandomWordsFulfilledEvent(
			vrfContracts.CoordinatorV2Plus,
			randomWordsRequestedEvent.RequestId,
			subID,
			randomWordsRequestedEvent.Raw.BlockNumber,
			isNativeBilling,
			time.Second*10,
			l,
			0,
		)
		require.Error(t, err, "fulfillment should not be generated for the request which was made on reorged fork on Simulated Chain")
	})

}

func TestVRFv2PlusBatchFulfillmentEnabledDisabled(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient                   *seth.Client
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	chainID := network.ChainID
	configPtr := &config

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  vrfv2PlusCleanUpFn(&t, &sethClient, &configPtr, &env, &vrfContracts, &subIDsForCancellingAfterTest),
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	env, vrfContracts, vrfKey, nodeTypeToNodeMap, sethClient, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2Plus universe")

	//batchMaxGas := config.MaxGasLimit() (2.5 mill) + 400_000 = 2.9 mill
	//callback gas limit set by consumer = 500k
	// so 4 requests should be fulfilled inside 1 tx since 500k*4 < 2.9 mill

	batchFulfilmentMaxGas := *config.VRFv2Plus.General.MaxGasLimitCoordinatorConfig + 400_000
	config.VRFv2Plus.General.CallbackGasLimit = ptr.Ptr(uint32(500_000))

	expectedNumberOfFulfillmentsInsideOneBatchFulfillment := (batchFulfilmentMaxGas / *config.VRFv2Plus.General.CallbackGasLimit) - 1
	randRequestCount := expectedNumberOfFulfillmentsInsideOneBatchFulfillment

	t.Run("Batch Fulfillment Enabled", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = true

		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")

		//ensure that no job present on the node
		err = actions.DeleteJobs([]*nodeclient.ChainlinkClient{vrfNode.CLNode.API})
		require.NoError(t, err)

		batchFullfillmentEnabled := true
		// create job with batch fulfillment enabled
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2Plus.Address(),
			BatchCoordinatorAddress:       vrfContracts.BatchCoordinatorV2Plus.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       batchFullfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
			VRFOwnerConfig:                nil,
		}

		l.Info().
			Msg("Creating VRFV2 Plus Job with `batchFulfillmentEnabled = true`")
		job, err := vrfv2plus.CreateVRFV2PlusJob(
			vrfNode.CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, "error creating job with higher timeout")
		vrfNode.Job = job

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		if randRequestCount > math.MaxUint16 {
			t.Fatalf("rand request count overflows uint16: %d", randRequestCount)
		}
		configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest = ptr.Ptr(uint16(randRequestCount))

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
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
		var batchFulfillmentTxs []nodeclient.TransactionData
		for _, tx := range clNodeTxs.Data {
			if common.HexToAddress(tx.Attributes.To).Cmp(common.HexToAddress(vrfContracts.BatchCoordinatorV2Plus.Address())) == 0 {
				batchFulfillmentTxs = append(batchFulfillmentTxs, tx)
			}
		}
		fulfillmentTx, _, err := sethClient.Client.TransactionByHash(testcontext.Get(t), randomWordsFulfilledEvent.Raw.TxHash)
		require.NoError(t, err, "error getting tx from hash")

		fulfillmentTXToAddress := fulfillmentTx.To().String()
		l.Info().
			Str("Actual Fulfillment Tx To Address", fulfillmentTXToAddress).
			Str("BatchCoordinatorV2Plus Address", vrfContracts.BatchCoordinatorV2Plus.Address()).
			Msg("Fulfillment Tx To Address should be the BatchCoordinatorV2Plus Address when batch fulfillment is enabled")

		// verify that VRF node sends fulfillments via BatchCoordinator contract
		require.Equal(t, vrfContracts.BatchCoordinatorV2Plus.Address(), fulfillmentTXToAddress, "Fulfillment Tx To Address should be the BatchCoordinatorV2Plus Address when batch fulfillment is enabled")

		// verify that all fulfillments should be inside one tx
		// This check is disabled for live testnets since each testnet has different gas usage for similar tx
		if network.Simulated {
			fulfillmentTxReceipt, err := sethClient.Client.TransactionReceipt(testcontext.Get(t), fulfillmentTx.Hash())
			require.NoError(t, err)
			randomWordsFulfilledLogs, err := contracts.ParseRandomWordsFulfilledLogs(vrfContracts.CoordinatorV2Plus, fulfillmentTxReceipt.Logs)
			require.NoError(t, err)
			require.Equal(t, 1, len(batchFulfillmentTxs))
			require.Equal(t, int(randRequestCount), len(randomWordsFulfilledLogs))
		}
	})
	t.Run("Batch Fulfillment Disabled", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = true

		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")
		//ensure that no job present on the node
		err = actions.DeleteJobs([]*nodeclient.ChainlinkClient{vrfNode.CLNode.API})
		require.NoError(t, err)

		batchFullfillmentEnabled := false

		//create job with batchFulfillmentEnabled = false
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2Plus.Address(),
			BatchCoordinatorAddress:       vrfContracts.BatchCoordinatorV2Plus.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       batchFullfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
			VRFOwnerConfig:                nil,
		}

		l.Info().
			Msg("Creating VRFV2 Plus Job with `batchFulfillmentEnabled = false`")
		job, err := vrfv2plus.CreateVRFV2PlusJob(
			vrfNode.CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, "error creating job with higher timeout")
		vrfNode.Job = job

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			sethClient,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		if randRequestCount > math.MaxUint16 {
			t.Fatalf("rand request count overflows uint16: %d", randRequestCount)
		}
		configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest = ptr.Ptr(uint16(randRequestCount))

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
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
			Str("CoordinatorV2Plus Address", vrfContracts.CoordinatorV2Plus.Address()).
			Msg("Fulfillment Tx To Address should be the CoordinatorV2Plus Address when batch fulfillment is disabled")

		// verify that VRF node sends fulfillments via Coordinator contract
		require.Equal(t, vrfContracts.CoordinatorV2Plus.Address(), fulfillmentTXToAddress, "Fulfillment Tx To Address should be the CoordinatorV2Plus Address when batch fulfillment is disabled")

		clNodeTxs, resp, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTransactions()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)

		var singleFulfillmentTxs []nodeclient.TransactionData
		for _, tx := range clNodeTxs.Data {
			if common.HexToAddress(tx.Attributes.To).Cmp(common.HexToAddress(vrfContracts.CoordinatorV2Plus.Address())) == 0 {
				singleFulfillmentTxs = append(singleFulfillmentTxs, tx)
			}
		}
		// verify that all fulfillments should be in separate txs
		require.Equal(t, int(randRequestCount), len(singleFulfillmentTxs))
	})
}
