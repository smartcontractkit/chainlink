package common

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_owner"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func LogSubDetails(
	l zerolog.Logger,
	subscription contracts.Subscription,
	subID string,
	coordinator contracts.Coordinator,
) {
	log := l.Info().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Str("Subscription ID", subID).
		Str("Subscription Owner", subscription.SubOwner.String()).
		Interface("Subscription Consumers", subscription.Consumers)
	if subscription.NativeBalance != nil {
		log = log.Str("Native Token Balance", assets.FormatWei(subscription.NativeBalance))
	}
	log.Msg("Subscription Data")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.Coordinator,
	randomWordsRequestedEvent *contracts.CoordinatorRandomWordsRequested,
	isNativeBilling bool,
	keyNum int,
) {
	l.Info().
		Int("KeyNum", keyNum).
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId).
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
	coordinator contracts.Coordinator,
	randomWordsFulfilledEvent *contracts.CoordinatorRandomWordsFulfilled,
	isNativeBilling bool,
	keyNum int,
) {
	l.Info().
		Int("KeyNum", keyNum).
		Bool("Native Billing", isNativeBilling).
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubId).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Uint64("BlockNumber", randomWordsFulfilledEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsFulfilledEvent.Raw.BlockHash.String()).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogFulfillmentDetailsLinkBilling(
	l zerolog.Logger,
	wrapperConsumerJuelsBalanceBeforeRequest *big.Int,
	wrapperConsumerJuelsBalanceAfterRequest *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *contracts.CoordinatorRandomWordsFulfilled,
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
	randomWordsFulfilledEvent *contracts.CoordinatorRandomWordsFulfilled,
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

func LogRandomWordsForcedEvent(
	l zerolog.Logger,
	vrfOwner contracts.VRFOwner,
	randomWordsForcedEvent *vrf_owner.VRFOwnerRandomWordsForced,
) {
	l.Debug().
		Str("VRFOwner", vrfOwner.Address()).
		Uint64("Sub ID", randomWordsForcedEvent.SubId).
		Str("TX Hash", randomWordsForcedEvent.Raw.TxHash.String()).
		Str("Request ID", randomWordsForcedEvent.RequestId.String()).
		Str("Sender", randomWordsForcedEvent.Sender.String()).
		Msg("RandomWordsForced Event (TX metadata)")
}
