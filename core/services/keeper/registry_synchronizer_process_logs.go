package keeper

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

func (rs *RegistrySynchronizer) processLogs() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go rs.handleSyncRegistryLog(wg.Done)
	go rs.handleUpkeepCanceledLogs(wg.Done)
	go rs.handleUpkeepRegisteredLogs(wg.Done)
	go rs.handleUpkeepPerformedLogs(wg.Done)
	wg.Wait()
}

func (rs *RegistrySynchronizer) handleSyncRegistryLog(done func()) {
	defer done()
	i, exists := rs.mailRoom.mbSyncRegistry.Retrieve()
	if !exists {
		return
	}
	broadcast, ok := i.(log.Broadcast)
	if !ok {
		logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
		return
	}
	txHash := broadcast.RawLog().TxHash.Hex()
	logger.Debugw("RegistrySynchronizer: processing SyncRegistry log", "jobID", rs.job.ID, "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
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
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast)
	logger.ErrorIf(errors.Wrapf(err, "RegistrySynchronizer: unable to mark SyncRegistryLog log as consumed, jobID: %d, log: %v", rs.job.ID, broadcast.String()))
}

func (rs *RegistrySynchronizer) handleUpkeepCanceledLogs(done func()) {
	defer done()
	for {
		i, exists := rs.mailRoom.mbUpkeepCanceled.Retrieve()
		if !exists {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.handleUpkeepCancelled(broadcast)
	}
}

func (rs *RegistrySynchronizer) handleUpkeepCancelled(broadcast log.Broadcast) {
	txHash := broadcast.RawLog().TxHash.Hex()
	logger.Debugw("RegistrySynchronizer: processing UpkeepCanceled log", "jobID", rs.job.ID, "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
		return
	}
	if was {
		return
	}
	log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepCanceled)
	if !ok {
		logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepCanceled log but got %T", log)
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	affected, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, []int64{log.Id.Int64()})
	if err != nil {
		logger.Error(errors.Wrapf(err, "RegistrySynchronizer: unable to batch delete upkeeps, jobID: %d", rs.job.ID))
		return
	}
	logger.Debugw(fmt.Sprintf("RegistrySynchronizer: deleted %v upkeep registrations", affected), "jobID", rs.job.ID, "txHash", txHash)

	ctx, cancel = postgres.DefaultQueryCtx()
	defer cancel()
	err = rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast)
	logger.ErrorIf(errors.Wrapf(err, "RegistrySynchronizer: unable to mark KeeperRegistryUpkeepCanceled log as consumed, jobID: %d, log: %v", rs.job.ID, broadcast.String()))
}

func (rs *RegistrySynchronizer) handleUpkeepRegisteredLogs(done func()) {
	defer done()
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	registry, err := rs.orm.RegistryForJob(ctx, rs.job.ID)
	if err != nil {
		logger.Error(errors.Wrapf(err, "RegistrySynchronizer: unable to find registry for job, jobID: %d", rs.job.ID))
		return
	}
	for {
		i, exists := rs.mailRoom.mbUpkeepRegistered.Retrieve()
		if !exists {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.HandleUpkeepRegistered(broadcast, registry)
	}
}

func (rs *RegistrySynchronizer) HandleUpkeepRegistered(broadcast log.Broadcast, registry Registry) {
	txHash := broadcast.RawLog().TxHash.Hex()
	logger.Debugw("RegistrySynchronizer: processing UpkeepRegistered log", "jobID", rs.job.ID, "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
		return
	}
	if was {
		return
	}
	log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepRegistered)
	if !ok {
		logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepRegistered log but got %T", log)
		return
	}
	err = rs.syncUpkeep(registry, log.Id.Int64())
	if err != nil {
		logger.Error(err)
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast)
	logger.ErrorIf(errors.Wrapf(err, "RegistrySynchronizer: unable to mark KeeperRegistryUpkeepRegistered log as consumed, jobID: %d, log: %v", rs.job.ID, broadcast.String()))
}

func (rs *RegistrySynchronizer) handleUpkeepPerformedLogs(done func()) {
	defer done()
	for {
		i, exists := rs.mailRoom.mbUpkeepPerformed.Retrieve()
		if !exists {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			logger.Errorf("RegistrySynchronizer: invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.handleUpkeepPerformed(broadcast)
	}
}

func (rs *RegistrySynchronizer) handleUpkeepPerformed(broadcast log.Broadcast) {
	txHash := broadcast.RawLog().TxHash.Hex()
	logger.Debugw("RegistrySynchronizer: processing UpkeepPerformed log", "jobID", rs.job.ID, "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		logger.Warn(errors.Wrapf(err, "RegistrySynchronizer: unable to check if log was consumed, jobID: %d", rs.job.ID))
		return
	}
	if was {
		return
	}
	log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepPerformed)
	if !ok {
		logger.Errorf("RegistrySynchronizer: invariant violation, expected UpkeepPerformed log but got %T", log)
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	db := rs.orm.DB.WithContext(ctx)
	// set last run to 0 so that keeper can resume checkUpkeep()
	err = rs.orm.SetLastRunHeightForUpkeepOnJob(db, rs.job.ID, log.Id.Int64(), 0)
	if err != nil {
		logger.Error(err)
		return
	}
	ctx, cancel = postgres.DefaultQueryCtx()
	defer cancel()
	err = rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast)
	logger.ErrorIf(errors.Wrapf(err, "RegistrySynchronizer: unable to mark KeeperRegistryUpkeepPerformed log as consumed, jobID: %d, log: %v", rs.job.ID, broadcast.String()))
}
