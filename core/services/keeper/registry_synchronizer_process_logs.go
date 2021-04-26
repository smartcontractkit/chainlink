package keeper

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func (rs *RegistrySynchronizer) processLogs(head models.Head) {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go rs.handleSyncRegistryLog(head, wg.Done)
	go rs.handleUpkeepCanceledLogs(head, wg.Done)
	go rs.handleUpkeepRegisteredLogs(head, wg.Done)
	go rs.handleUpkeepPerformedLogs(head, wg.Done)
	wg.Wait()
}

func (rs *RegistrySynchronizer) handleSyncRegistryLog(head models.Head, done func()) {
	defer done()
	oldEnough := isOldEnoughConstructor(head, rs.minConfirmations)
	i := rs.mailRoom.mbSyncRegistry.RetrieveIf(oldEnough)
	if i == nil {
		return
	}
	broadcast, ok := i.(log.Broadcast)
	if !ok {
		logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
		return
	}
	txHash := broadcast.RawLog().TxHash.Hex()
	logger.Debugw("RegistrySynchronizer: processing SyncRegistry log", "jobID", rs.job.ID, "txHash", txHash)
	was, err := broadcast.WasAlreadyConsumed()
	if err != nil {
		logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
		return
	}
	if was {
		return
	}
	_, err = rs.syncRegistry()
	if err != nil {
		logger.Error(errors.Wrapf(err, "RegistrySynchronizer: unable to sync registry, jobID: %d", rs.job.ID))
		return
	}
	err = broadcast.MarkConsumed()
	logger.ErrorIf((errors.Wrapf(err, "RegistrySynchronizer: unable to mark log as consumed, jobID: %d", rs.job.ID)))
}

func (rs *RegistrySynchronizer) handleUpkeepCanceledLogs(head models.Head, done func()) {
	defer done()
	oldEnough := isOldEnoughConstructor(head, rs.minConfirmations)
	for {
		i := rs.mailRoom.mbUpkeepCanceled.RetrieveIf(oldEnough)
		if i == nil {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		txHash := broadcast.RawLog().TxHash.Hex()
		logger.Debugw("RegistrySynchronizer: processing UpkeepCanceled log", "jobID", rs.job.ID, "txHash", txHash)
		was, err := broadcast.WasAlreadyConsumed()
		if err != nil {
			logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
			continue
		}
		if was {
			continue
		}
		log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepCanceled)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepCanceled log but got %T", log)
			continue
		}
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		affected, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, []int64{log.Id.Int64()})
		if err != nil {
			logger.Error(errors.Wrapf(err, "RegistrySynchronizer: unable to batch delete upkeeps, jobID: %d", rs.job.ID))
			continue
		}
		logger.Debugw(fmt.Sprintf("RegistrySynchronizer: deleted %v upkeep registrations", affected), "jobID", rs.job.ID, "txHash", txHash)
		err = broadcast.MarkConsumed()
		logger.ErrorIf((errors.Wrapf(err, "RegistrySynchronizer: unable to mark log as consumed, jobID: %d", rs.job.ID)))
	}
}

func (rs *RegistrySynchronizer) handleUpkeepRegisteredLogs(head models.Head, done func()) {
	defer done()
	oldEnough := isOldEnoughConstructor(head, rs.minConfirmations)
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	registry, err := rs.orm.RegistryForJob(ctx, rs.job.ID)
	if err != nil {
		logger.Error(errors.Wrapf(err, "RegistrySynchronizer: unable to find registry for job, jobID: %d", rs.job.ID))
		return
	}
	for {
		i := rs.mailRoom.mbUpkeepRegistered.RetrieveIf(oldEnough)
		if i == nil {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		txHash := broadcast.RawLog().TxHash.Hex()
		logger.Debugw("RegistrySynchronizer: processing UpkeepRegistered log", "jobID", rs.job.ID, "txHash", txHash)
		was, err := broadcast.WasAlreadyConsumed()
		if err != nil {
			logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
			continue
		}
		if was {
			continue
		}
		log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepRegistered)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepRegistered log but got %T", log)
			continue
		}
		err = rs.syncUpkeep(registry, log.Id.Int64())
		if err != nil {
			logger.Error(err)
			continue
		}
		err = broadcast.MarkConsumed()
		logger.ErrorIf((errors.Wrapf(err, "RegistrySynchronizer: unable to mark log as consumed, jobID: %d", rs.job.ID)))
	}
}

func (rs *RegistrySynchronizer) handleUpkeepPerformedLogs(head models.Head, done func()) {
	defer done()
	oldEnough := isOldEnoughConstructor(head, rs.minConfirmations)
	for {
		i := rs.mailRoom.mbUpkeepPerformed.RetrieveIf(oldEnough)
		if i == nil {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		txHash := broadcast.RawLog().TxHash.Hex()
		logger.Debugw("RegistrySynchronizer: processing UpkeepPerformed log", "jobID", rs.job.ID, "txHash", txHash)
		was, err := broadcast.WasAlreadyConsumed()
		if err != nil {
			logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
			continue
		}
		if was {
			continue
		}
		log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepPerformed)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepPerformed log but got %T", log)
			continue
		}
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		// set last run to 0 so that keeper can resume checkUpkeep()
		err = rs.orm.SetLastRunHeightForUpkeepOnJob(ctx, rs.job.ID, log.Id.Int64(), 0)
		if err != nil {
			logger.Error(err)
			continue
		}
		err = broadcast.MarkConsumed()
		logger.ErrorIf((errors.Wrapf(err, "RegistrySynchronizer: unable to mark log as consumed, jobID: %d", rs.job.ID)))
	}
}

func isOldEnoughConstructor(head models.Head, minConfirmations uint64) func(interface{}) bool {
	return func(i interface{}) bool {
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			return true // we want to get bad data out of the queue
		}
		logHeight := broadcast.RawLog().BlockNumber
		return (logHeight + uint64(minConfirmations) - 1) <= uint64(head.Number)
	}
}
