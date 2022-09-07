package keeper

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry1_3 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (rs *RegistrySynchronizer) processLogs() {
	for _, broadcast := range rs.mbLogs.RetrieveAll() {
		eventLog := broadcast.DecodedLog()
		if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
			rs.logger.Panicf("processLogs: ignoring nil value, type: %T", eventLog)
			continue
		}

		was, err := rs.logBroadcaster.WasAlreadyConsumed(broadcast)
		if err != nil {
			rs.logger.Warn(errors.Wrap(err, "unable to check if log was consumed"))
			continue
		} else if was {
			continue
		}

		switch eventLog.(type) {
		case *registry1_1.KeeperRegistryKeepersUpdated,
			*registry1_1.KeeperRegistryConfigSet,
			*registry1_2.KeeperRegistryKeepersUpdated,
			*registry1_2.KeeperRegistryConfigSet,
			*registry1_3.KeeperRegistryKeepersUpdated,
			*registry1_3.KeeperRegistryConfigSet:
			err = rs.handleSyncRegistryLog(broadcast)

		case *registry1_1.KeeperRegistryUpkeepCanceled,
			*registry1_2.KeeperRegistryUpkeepCanceled,
			*registry1_3.KeeperRegistryUpkeepCanceled:
			err = rs.handleUpkeepCancelled(broadcast)

		case *registry1_1.KeeperRegistryUpkeepRegistered,
			*registry1_2.KeeperRegistryUpkeepRegistered,
			*registry1_3.KeeperRegistryUpkeepRegistered:
			err = rs.handleUpkeepRegistered(broadcast)

		case *registry1_1.KeeperRegistryUpkeepPerformed,
			*registry1_2.KeeperRegistryUpkeepPerformed,
			*registry1_3.KeeperRegistryUpkeepPerformed:
			err = rs.handleUpkeepPerformed(broadcast)

		case *registry1_2.KeeperRegistryUpkeepGasLimitSet,
			*registry1_3.KeeperRegistryUpkeepGasLimitSet:
			err = rs.handleUpkeepGasLimitSet(broadcast)

		case *registry1_2.KeeperRegistryUpkeepReceived,
			*registry1_3.KeeperRegistryUpkeepReceived:
			err = rs.handleUpkeepReceived(broadcast)

		case *registry1_2.KeeperRegistryUpkeepMigrated,
			*registry1_3.KeeperRegistryUpkeepMigrated:
			err = rs.handleUpkeepMigrated(broadcast)

		case *registry1_3.KeeperRegistryUpkeepPaused:
			err = rs.handleUpkeepPaused(broadcast)

		case *registry1_3.KeeperRegistryUpkeepUnpaused:
			err = rs.handleUpkeepUnpaused(broadcast)

		case *registry1_3.KeeperRegistryUpkeepCheckDataUpdated:
			err = rs.handleUpkeepCheckDataUpdated(broadcast)

		default:
			rs.logger.Warn("unexpected log type")
			// Don't `continue` -- we still want to mark this log as consumed
		}

		if err != nil {
			rs.logger.Error(err)
		}

		err = rs.logBroadcaster.MarkConsumed(broadcast)
		if err != nil {
			rs.logger.Error(errors.Wrapf(err, "unable to mark %T log as consumed, log: %v", broadcast.RawLog(), broadcast.String()))
		}
	}
}

func (rs *RegistrySynchronizer) handleSyncRegistryLog(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing SyncRegistry log", "txHash", broadcast.RawLog().TxHash.Hex())

	_, err := rs.syncRegistry()
	if err != nil {
		return errors.Wrap(err, "unable to sync registry")
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepCancelled(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepCanceled log", "txHash", broadcast.RawLog().TxHash.Hex())

	cancelledID, err := rs.registryWrapper.GetCancelledUpkeepIDFromLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch cancelled upkeep ID from log")
	}

	affected, err := rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, []utils.Big{*utils.NewBig(cancelledID)})
	if err != nil {
		return errors.Wrap(err, "unable to batch delete upkeeps")
	}
	rs.logger.Debugw(fmt.Sprintf("deleted %v upkeep registrations", affected), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepRegistered(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepRegistered log", "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	upkeepID, err := rs.registryWrapper.GetUpkeepIdFromRegistrationLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from registration log")
	}

	err = rs.syncUpkeep(&rs.registryWrapper, registry, utils.NewBig(upkeepID))
	if err != nil {
		return errors.Wrapf(err, "failed to sync upkeep, log: %v", broadcast.String())
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepPerformed(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepPerformed log", "jobID", rs.job.ID, "txHash", broadcast.RawLog().TxHash.Hex())

	log, err := rs.registryWrapper.ParseUpkeepPerformedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from performed log")
	}
	rowsAffected, err := rs.orm.SetLastRunInfoForUpkeepOnJob(rs.job.ID, utils.NewBig(log.UpkeepID), int64(broadcast.RawLog().BlockNumber), ethkey.EIP55AddressFromAddress(log.FromKeeper))
	if err != nil {
		return errors.Wrap(err, "failed to set last run to 0")
	}
	rs.logger.Debugw("updated db for UpkeepPerformed log",
		"jobID", rs.job.ID,
		"upkeepID", log.UpkeepID.String(),
		"blockNumber", int64(broadcast.RawLog().BlockNumber),
		"fromAddr", ethkey.EIP55AddressFromAddress(log.FromKeeper),
		"rowsAffected", rowsAffected,
	)
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepGasLimitSet(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepGasLimitSet log", "jobID", rs.job.ID, "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	upkeepID, err := rs.registryWrapper.GetIDFromGasLimitSetLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from gas limit set log")
	}

	err = rs.syncUpkeep(&rs.registryWrapper, registry, utils.NewBig(upkeepID))
	if err != nil {
		return errors.Wrapf(err, "failed to sync upkeep, log: %v", broadcast.String())
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepReceived(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepReceived log", "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	upkeepID, err := rs.registryWrapper.GetUpkeepIdFromReceivedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from received log")
	}

	err = rs.syncUpkeep(&rs.registryWrapper, registry, utils.NewBig(upkeepID))
	if err != nil {
		return errors.Wrapf(err, "failed to sync upkeep, log: %v", broadcast.String())
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepMigrated(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepMigrated log", "txHash", broadcast.RawLog().TxHash.Hex())

	migratedID, err := rs.registryWrapper.GetUpkeepIdFromMigratedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch migrated upkeep ID from log")
	}

	affected, err := rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, []utils.Big{*utils.NewBig(migratedID)})
	if err != nil {
		return errors.Wrap(err, "unable to batch delete upkeeps")
	}
	rs.logger.Debugw(fmt.Sprintf("deleted %v upkeep registrations", affected), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepPaused(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepPaused log", "txHash", broadcast.RawLog().TxHash.Hex())

	pausedUpkeepId, err := rs.registryWrapper.GetUpkeepIdFromUpkeepPausedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from upkeep paused log")
	}

	_, err = rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, []utils.Big{*utils.NewBig(pausedUpkeepId)})
	if err != nil {
		return errors.Wrap(err, "unable to batch delete upkeeps")
	}
	rs.logger.Debugw(fmt.Sprintf("paused upkeep %s", pausedUpkeepId.String()), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepUnpaused(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepUnpaused log", "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	unpausedUpkeepId, err := rs.registryWrapper.GetUpkeepIdFromUpkeepUnpausedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from upkeep unpaused log")
	}

	err = rs.syncUpkeep(&rs.registryWrapper, registry, utils.NewBig(unpausedUpkeepId))
	if err != nil {
		return errors.Wrapf(err, "failed to sync upkeep, log: %s", broadcast.String())
	}
	rs.logger.Debugw(fmt.Sprintf("unpaused upkeep %s", unpausedUpkeepId.String()), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepCheckDataUpdated(broadcast log.Broadcast) error {
	rs.logger.Debugw("processing Upkeep check data updated log", "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	updateLog, err := rs.registryWrapper.ParseUpkeepCheckDataUpdatedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to parse update log from upkeep check data updated log")
	}

	err = rs.syncUpkeep(&rs.registryWrapper, registry, utils.NewBig(updateLog.UpkeepID))
	if err != nil {
		return errors.Wrapf(err, "unable to update check data for upkeep %s", updateLog.UpkeepID.String())
	}

	rs.logger.Debugw(fmt.Sprintf("updated check data for upkeep %s", updateLog.UpkeepID.String()), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}
