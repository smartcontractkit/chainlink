package keeper

import (
	"reflect"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func (rs *RegistrySynchronizer) OnConnect() {}

func (rs *RegistrySynchronizer) OnDisconnect() {}

func (rs *RegistrySynchronizer) JobID() models.JobID {
	return models.JobID{}
}

func (rs *RegistrySynchronizer) JobIDV2() int32 {
	return rs.job.ID
}

func (rs *RegistrySynchronizer) IsV2Job() bool {
	return true
}

func (rs *RegistrySynchronizer) HandleLog(broadcast log.Broadcast) {
	log := broadcast.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Errorf("RegistrySynchronizer: HandleLog: ignoring nil value, type: %T", broadcast)
		return
	}

	logger.Debugw(
		"RegistrySynchronizer: received log, waiting for confirmations",
		"jobID", rs.job.ID,
		"logType", reflect.TypeOf(log),
		"txHash", broadcast.RawLog().TxHash.Hex(),
	)

	switch log := log.(type) {
	case *keeper_registry_wrapper.KeeperRegistryKeepersUpdated:
		rs.mailRoom.mbSyncRegistry.Deliver(broadcast) // same mailbox because same action
	case *keeper_registry_wrapper.KeeperRegistryConfigSet:
		rs.mailRoom.mbSyncRegistry.Deliver(broadcast) // same mailbox because same action
	case *keeper_registry_wrapper.KeeperRegistryUpkeepCanceled:
		rs.mailRoom.mbUpkeepCanceled.Deliver(broadcast)
	case *keeper_registry_wrapper.KeeperRegistryUpkeepRegistered:
		rs.mailRoom.mbUpkeepRegistered.Deliver(broadcast)
	case *keeper_registry_wrapper.KeeperRegistryUpkeepPerformed:
		rs.mailRoom.mbUpkeepPerformed.Deliver(broadcast)
	default:
		logger.Warnf("unexpected log type %T", log)
	}
}
