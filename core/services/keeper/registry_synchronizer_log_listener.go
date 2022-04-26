package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
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
		*registry1_2.KeeperRegistryConfigSet:
		// same mailbox because same action for config set and keepers updated
		wasOverCapacity = rs.mailRoom.mbSyncRegistry.Deliver(broadcast)
		mailboxName = "mbSyncRegistry"
	case *registry1_1.KeeperRegistryUpkeepCanceled,
		*registry1_2.KeeperRegistryUpkeepCanceled:
		wasOverCapacity = rs.mailRoom.mbUpkeepCanceled.Deliver(broadcast)
		mailboxName = "mbUpkeepCanceled"
	case *registry1_1.KeeperRegistryUpkeepRegistered,
		*registry1_2.KeeperRegistryUpkeepRegistered:
		wasOverCapacity = rs.mailRoom.mbUpkeepRegistered.Deliver(broadcast)
		mailboxName = "mbUpkeepRegistered"
	case *registry1_1.KeeperRegistryUpkeepPerformed,
		*registry1_2.KeeperRegistryUpkeepPerformed:
		wasOverCapacity = rs.mailRoom.mbUpkeepPerformed.Deliver(broadcast)
		mailboxName = "mbUpkeepPerformed"
	case *registry1_2.KeeperRegistryUpkeepGasLimitSet:
		wasOverCapacity = rs.mailRoom.mbUpkeepGasLimitSet.Deliver(broadcast)
		mailboxName = "mbUpkeepGasLimitSet"
	default:
		svcLogger.Warn("unexpected log type")
	}

	if wasOverCapacity {
		svcLogger.With("mailboxName", mailboxName).Errorf("mailbox is over capacity - dropped the oldest unprocessed item")
	}
}
