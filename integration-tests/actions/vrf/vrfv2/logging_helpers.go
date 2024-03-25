package vrfv2

import (
	"fmt"

	"github.com/rs/zerolog"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"
)

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2.GetSubscription, subID uint64, coordinator contracts.VRFCoordinatorV2) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Uint64("Subscription ID", subID).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2,
	randomWordsRequestedEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Uint64("Subscription ID", randomWordsRequestedEvent.SubId).
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
	coordinator contracts.VRFCoordinatorV2,
	randomWordsFulfilledEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Uint64("BlockNumber", randomWordsFulfilledEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsFulfilledEvent.Raw.BlockHash.String()).
		Msg("RandomWordsFulfilled Event (TX metadata)")
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

func logRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID uint64,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	keyhash [32]byte,
) {
	l.Info().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Uint64("SubID", subID).
		Uint16("MinimumConfirmations", minimumConfirmations).
		Uint32("CallbackGasLimit", callbackGasLimit).
		Uint32("NumberOfWords", numberOfWords).
		Uint16("RandomnessRequestCountPerRequest", randomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", randomnessRequestCountPerRequestDeviation).
		Str("Keyhash", fmt.Sprintf("0x%x", keyhash)).
		Msg("Requesting randomness")
}
