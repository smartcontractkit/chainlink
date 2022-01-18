package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
)

func (rs *RegistrySynchronizer) JobID() int32 {
	return rs.job.ID
}

func (rs *RegistrySynchronizer) HandleLog(broadcast log.Broadcast) {
	eventLog := broadcast.DecodedLog()
	if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
		rs.logger.Errorf("HandleLog: ignoring nil value, type: %T", broadcast)
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
	case *keeper_registry_wrapper.KeeperRegistryKeepersUpdated:
		wasOverCapacity = rs.mailRoom.mbSyncRegistry.Deliver(broadcast) // same mailbox because same action
		mailboxName = "mbSyncRegistry"
	case *keeper_registry_wrapper.KeeperRegistryConfigSet:
		wasOverCapacity = rs.mailRoom.mbSyncRegistry.Deliver(broadcast) // same mailbox because same action
		mailboxName = "mbSyncRegistry"
	case *keeper_registry_wrapper.KeeperRegistryUpkeepCanceled:
		wasOverCapacity = rs.mailRoom.mbUpkeepCanceled.Deliver(broadcast)
		mailboxName = "mbUpkeepCanceled"
	case *keeper_registry_wrapper.KeeperRegistryUpkeepRegistered:
		wasOverCapacity = rs.mailRoom.mbUpkeepRegistered.Deliver(broadcast)
		mailboxName = "mbUpkeepRegistered"
	case *keeper_registry_wrapper.KeeperRegistryUpkeepPerformed:
		wasOverCapacity = rs.mailRoom.mbUpkeepPerformed.Deliver(broadcast)
		mailboxName = "mbUpkeepPerformed"
	default:
		svcLogger.Warn("unexpected log type")
	}

	if wasOverCapacity {
		svcLogger.With("mailboxName", mailboxName).Errorf("mailbox is over capacity - dropped the oldest unprocessed item")
	}
}
