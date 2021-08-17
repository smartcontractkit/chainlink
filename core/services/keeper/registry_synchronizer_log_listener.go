package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

func (rs *RegistrySynchronizer) JobID() int32 {
	return rs.job.ID
}

func (rs *RegistrySynchronizer) HandleLog(broadcast log.Broadcast) {
	eventLog := broadcast.DecodedLog()
	if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
		logger.Errorf("RegistrySynchronizer: HandleLog: ignoring nil value, type: %T", broadcast)
		return
	}

	logger.Debugw(
		"RegistrySynchronizer: received log, waiting for confirmations",
		"jobID", rs.job.ID,
		"logType", reflect.TypeOf(eventLog),
		"txHash", broadcast.RawLog().TxHash.Hex(),
	)

	var mailboxName string
	var wasOverCapacity bool
	switch eventLog := eventLog.(type) {
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
		logger.Warnf("unexpected log type %T", eventLog)
	}
	if wasOverCapacity {
		logger.Errorf("RegistrySynchronizer: %v mailbox is over capacity - dropped the oldest unprocessed item", mailboxName)
	}
}
