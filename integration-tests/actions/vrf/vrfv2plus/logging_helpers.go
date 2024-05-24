package vrfv2plus

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
)

func LogRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID *big.Int,
	isNativeBilling bool,
	keyHash [32]byte,
	config *vrfv2plus_config.General,
	keyNum int,
) {
	l.Info().
		Int("KeyNum", keyNum).
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

func LogMigrationCompletedEvent(l zerolog.Logger, migrationCompletedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, coordinator contracts.Coordinator) {
	l.Info().
		Str("Subscription ID", migrationCompletedEvent.SubId.String()).
		Str("Migrated From Coordinator", coordinator.Address()).
		Str("Migrated To Coordinator", migrationCompletedEvent.NewCoordinator.String()).
		Msg("MigrationCompleted Event")
}

func LogSubDetailsAfterMigration(l zerolog.Logger, newCoordinator contracts.Coordinator, subID *big.Int, migratedSubscription vrf_v2plus_upgraded_version.GetSubscription) {
	l.Info().
		Str("New Coordinator", newCoordinator.Address()).
		Str("Subscription ID", subID.String()).
		Str("Juels Balance", migratedSubscription.Balance.String()).
		Str("Native Token Balance", migratedSubscription.NativeBalance.String()).
		Str("Subscription Owner", migratedSubscription.SubOwner.String()).
		Interface("Subscription Consumers", migratedSubscription.Consumers).
		Msg("Subscription Data After Migration to New Coordinator")
}
