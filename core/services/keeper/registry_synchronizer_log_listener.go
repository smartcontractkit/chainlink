package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry1_3 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_3"
)

func (rs *RegistrySynchronizer) JobID() int32 {
	return rs.job.ID
}

func (rs *RegistrySynchronizer) HandleLog(broadcast log.Broadcast) {
	eventLog := broadcast.DecodedLog()
	if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
		rs.logger.Panicf("HandleLog: ignoring nil value, type: %T", broadcast)
		return
	}

	svcLogger := rs.logger.With(
		"logType", reflect.TypeOf(eventLog),
		"txHash", broadcast.RawLog().TxHash.Hex(),
	)

	svcLogger.Debug("received log, waiting for confirmations")

	var mailboxName string
	var wasOverCapacity bool
	switch eventLog.(type) {

	case *registry1_1.KeeperRegistryKeepersUpdated,
		*registry1_1.KeeperRegistryConfigSet,
		*registry1_2.KeeperRegistryKeepersUpdated,
		*registry1_2.KeeperRegistryConfigSet,
		*registry1_3.KeeperRegistryKeepersUpdated,
		*registry1_3.KeeperRegistryConfigSet:
		// same mailbox because same action for config set and keepers updated
		svcLogger.Debug("delivering to sync registry mailbox")
		wasOverCapacity = rs.mailRoom.mbSyncRegistry.Deliver(broadcast)
		mailboxName = "mbSyncRegistry"
	case *registry1_1.KeeperRegistryUpkeepCanceled,
		*registry1_2.KeeperRegistryUpkeepCanceled,
		*registry1_3.KeeperRegistryUpkeepCanceled:
		svcLogger.Debug("delivering to upkeep canceled mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepCanceled.Deliver(broadcast)
		mailboxName = "mbUpkeepCanceled"
	case *registry1_1.KeeperRegistryUpkeepRegistered,
		*registry1_2.KeeperRegistryUpkeepRegistered,
		*registry1_3.KeeperRegistryUpkeepRegistered:
		svcLogger.Debug("delivering to upkeep registered mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepRegistered.Deliver(broadcast)
		mailboxName = "mbUpkeepRegistered"
	case *registry1_1.KeeperRegistryUpkeepPerformed,
		*registry1_2.KeeperRegistryUpkeepPerformed,
		*registry1_3.KeeperRegistryUpkeepPerformed:
		svcLogger.Debug("delivering to upkeep performed mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepPerformed.Deliver(broadcast)
		mailboxName = "mbUpkeepPerformed"
	case *registry1_2.KeeperRegistryUpkeepGasLimitSet,
		*registry1_3.KeeperRegistryUpkeepGasLimitSet:
		svcLogger.Debug("delivering to upkeep gas limit set mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepGasLimitSet.Deliver(broadcast)
		mailboxName = "mbUpkeepGasLimitSet"
	case *registry1_2.KeeperRegistryUpkeepReceived,
		*registry1_3.KeeperRegistryUpkeepReceived:
		svcLogger.Debug("delivering to upkeep received mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepReceived.Deliver(broadcast)
		mailboxName = "mbUpkeepReceived"
	case *registry1_2.KeeperRegistryUpkeepMigrated,
		*registry1_3.KeeperRegistryUpkeepMigrated:
		svcLogger.Debug("delivering to upkeep migrated mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepMigrated.Deliver(broadcast)
		mailboxName = "mbUpkeepMigrated"
	case *registry1_3.KeeperRegistryUpkeepPaused:
		svcLogger.Debug("delivering to upkeep paused")
		wasOverCapacity = rs.mailRoom.mbUpkeepPaused.Deliver(broadcast)
		mailboxName = "mbUpkeepPaused"
	case *registry1_3.KeeperRegistryUpkeepUnpaused:
		svcLogger.Debug("delivering to upkeep unpaused mailbox")
		wasOverCapacity = rs.mailRoom.mbUpkeepUnpaused.Deliver(broadcast)
		mailboxName = "mbUpkeepUnpaused"
	case *registry1_3.KeeperRegistryUpkeepCheckDataUpdated:
		svcLogger.Debug("delivering to upkeep check data updated")
		wasOverCapacity = rs.mailRoom.mbUpkeepCheckDataUpdated.Deliver(broadcast)
		mailboxName = "mbUpkeepCheckDataUpdated"
	default:
		svcLogger.Warn("unexpected log type")
	}

	if wasOverCapacity {
		svcLogger.With("mailboxName", mailboxName).Errorf("mailbox is over capacity - dropped the oldest unprocessed item")
	}
}
