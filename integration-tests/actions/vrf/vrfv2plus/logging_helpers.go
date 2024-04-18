package vrfv2plus

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
)

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2_5.GetSubscription, subID *big.Int, coordinator contracts.VRFCoordinatorV2_5) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Str("Native Token Balance", assets.FormatWei(subscription.NativeBalance)).
		Str("Subscription ID", subID.String()).
		Str("Subscription Owner", subscription.SubOwner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsRequestedEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested,
	isNativeBilling bool,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Str("Keyhash", fmt.Sprintf("0x%x", randomWordsRequestedEvent.KeyHash)).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsFulfilledEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled,
	isNativeBilling bool,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Total Payment in Juels", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubID.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsRequestedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested,
	isNativeBilling bool,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Str("Keyhash", fmt.Sprintf("0x%x", randomWordsRequestedEvent.KeyHash)).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Str("TX Hash", randomWordsRequestedEvent.Raw.TxHash.String()).
		Uint64("BlockNumber", randomWordsRequestedEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsRequestedEvent.Raw.BlockHash.String()).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
	isNativeBilling bool,
) {
	l.Info().
		Bool("Native Billing", isNativeBilling).
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubId.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Uint64("BlockNumber", randomWordsFulfilledEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsFulfilledEvent.Raw.BlockHash.String()).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogMigrationCompletedEvent(l zerolog.Logger, migrationCompletedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, vrfv2PlusContracts *vrfcommon.VRFContracts) {
	l.Info().
		Str("Subscription ID", migrationCompletedEvent.SubId.String()).
		Str("Migrated From Coordinator", vrfv2PlusContracts.CoordinatorV2Plus.Address()).
		Str("Migrated To Coordinator", migrationCompletedEvent.NewCoordinator.String()).
		Msg("MigrationCompleted Event")
}

func LogSubDetailsAfterMigration(l zerolog.Logger, newCoordinator contracts.VRFCoordinatorV2PlusUpgradedVersion, subID *big.Int, migratedSubscription vrf_v2plus_upgraded_version.GetSubscription) {
	l.Info().
		Str("New Coordinator", newCoordinator.Address()).
		Str("Subscription ID", subID.String()).
		Str("Juels Balance", migratedSubscription.Balance.String()).
		Str("Native Token Balance", migratedSubscription.NativeBalance.String()).
		Str("Subscription Owner", migratedSubscription.SubOwner.String()).
		Interface("Subscription Consumers", migratedSubscription.Consumers).
		Msg("Subscription Data After Migration to New Coordinator")
}

func LogFulfillmentDetailsLinkBilling(
	l zerolog.Logger,
	wrapperConsumerJuelsBalanceBeforeRequest *big.Int,
	wrapperConsumerJuelsBalanceAfterRequest *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
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
}

func LogFulfillmentDetailsNativeBilling(
	l zerolog.Logger,
	wrapperConsumerBalanceBeforeRequestWei *big.Int,
	wrapperConsumerBalanceAfterRequestWei *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
	l.Info().
		Str("Consumer Balance Before Request", assets.FormatWei(wrapperConsumerBalanceBeforeRequestWei)).
		Str("Consumer Balance After Request", assets.FormatWei(wrapperConsumerBalanceAfterRequestWei)).
		Bool("Fulfilment Status", consumerStatus.Fulfilled).
		Str("Paid by Consumer Contract", assets.FormatWei(consumerStatus.Paid)).
		Str("Paid by Coordinator Sub", assets.FormatWei(randomWordsFulfilledEvent.Payment)).
		Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
		Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
		Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
		Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Msg("Random Words Request Fulfilment Details For Native Billing")
}

func LogRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID *big.Int,
	isNativeBilling bool,
	keyHash [32]byte,
	config *vrfv2plus_config.General) {
	l.Info().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Str("SubID", subID.String()).
		Bool("IsNativePayment", isNativeBilling).
		Uint16("MinimumConfirmations", *config.MinimumConfirmations).
		Uint32("CallbackGasLimit", *config.CallbackGasLimit).
		Uint32("NumberOfWords", *config.NumberOfWords).
		Str("KeyHash", fmt.Sprintf("0x%x", keyHash)).
		Uint16("RandomnessRequestCountPerRequest", *config.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *config.RandomnessRequestCountPerRequestDeviation).
		Msg("Requesting randomness")
}
