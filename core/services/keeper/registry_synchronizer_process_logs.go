package keeper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
)

func (rs *RegistrySynchronizer) processLogs(ctx context.Context) {
	for _, broadcast := range rs.mbLogs.RetrieveAll() {
		eventLog := broadcast.DecodedLog()
		if eventLog == nil || reflect.ValueOf(eventLog).IsNil() {
			rs.logger.Panicf("processLogs: ignoring nil value, type: %T", eventLog)
			continue
		}

		was, err := rs.logBroadcaster.WasAlreadyConsumed(ctx, broadcast)
		if err != nil {
			rs.logger.Warn(errors.Wrap(err, "unable to check if log was consumed"))
			continue
		} else if was {
			continue
		}

		switch eventLog.(type) {
		case *registry1_1.KeeperRegistryKeepersUpdated,
			*registry1_1.KeeperRegistryConfigSet:
			err = rs.handleSyncRegistryLog(ctx, broadcast)

		case *registry1_1.KeeperRegistryUpkeepCanceled:
			err = rs.handleUpkeepCancelled(ctx, broadcast)

		case *registry1_1.KeeperRegistryUpkeepRegistered:
			err = rs.handleUpkeepRegistered(ctx, broadcast)

		case *registry1_1.KeeperRegistryUpkeepPerformed:
			err = rs.handleUpkeepPerformed(ctx, broadcast)

		default:
			rs.logger.Warn("unexpected log type")
			// Don't `continue` -- we still want to mark this log as consumed
		}

		if err != nil {
			if ctx.Err() != nil {
				return
			}
			rs.logger.Error(err)
		}

		err = rs.logBroadcaster.MarkConsumed(ctx, nil, broadcast)
		if err != nil {
			rs.logger.Error(errors.Wrapf(err, "unable to mark %T log as consumed, log: %v", broadcast.RawLog(), broadcast.String()))
		}
	}
}

func (rs *RegistrySynchronizer) handleSyncRegistryLog(ctx context.Context, broadcast log.Broadcast) error {
	rs.logger.Debugw("processing SyncRegistry log", "txHash", broadcast.RawLog().TxHash.Hex())

	_, err := rs.syncRegistry(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to sync registry")
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepCancelled(ctx context.Context, broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepCanceled log", "txHash", broadcast.RawLog().TxHash.Hex())

	cancelledID, err := rs.registryWrapper.GetCancelledUpkeepIDFromLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch cancelled upkeep ID from log")
	}

	affected, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, []big.Big{*big.New(cancelledID)})
	if err != nil {
		return errors.Wrap(err, "unable to batch delete upkeeps")
	}
	rs.logger.Debugw(fmt.Sprintf("deleted %v upkeep registrations", affected), "txHash", broadcast.RawLog().TxHash.Hex())
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepRegistered(ctx context.Context, broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepRegistered log", "txHash", broadcast.RawLog().TxHash.Hex())

	registry, err := rs.orm.RegistryForJob(ctx, rs.job.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find registry for job")
	}

	upkeepID, err := rs.registryWrapper.GetUpkeepIdFromRegistrationLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from registration log")
	}

	err = rs.syncUpkeep(ctx, &rs.registryWrapper, registry, big.New(upkeepID))
	if err != nil {
		return errors.Wrapf(err, "failed to sync upkeep, log: %v", broadcast.String())
	}
	return nil
}

func (rs *RegistrySynchronizer) handleUpkeepPerformed(ctx context.Context, broadcast log.Broadcast) error {
	rs.logger.Debugw("processing UpkeepPerformed log", "jobID", rs.job.ID, "txHash", broadcast.RawLog().TxHash.Hex())

	log, err := rs.registryWrapper.ParseUpkeepPerformedLog(broadcast)
	if err != nil {
		return errors.Wrap(err, "Unable to fetch upkeep ID from performed log")
	}
	rowsAffected, err := rs.orm.SetLastRunInfoForUpkeepOnJob(ctx, rs.job.ID, big.New(log.UpkeepID), int64(broadcast.RawLog().BlockNumber), types.EIP55AddressFromAddress(log.FromKeeper))
	if err != nil {
		return errors.Wrap(err, "failed to set last run to 0")
	}
	rs.logger.Debugw("updated db for UpkeepPerformed log",
		"jobID", rs.job.ID,
		"upkeepID", log.UpkeepID.String(),
		"blockNumber", int64(broadcast.RawLog().BlockNumber),
		"fromAddr", types.EIP55AddressFromAddress(log.FromKeeper),
		"rowsAffected", rowsAffected,
	)
	return nil
}
