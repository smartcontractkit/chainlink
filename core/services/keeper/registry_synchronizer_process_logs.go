package keeper

import (
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
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
		rs.logger.Errorf("invariant violation, expected log.Broadcast but got %T", broadcast)
		return
	}
	txHash := broadcast.RawLog().TxHash.Hex()
	rs.logger.Debugw("processing SyncRegistry log", "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		rs.logger.With("error", err).Warn("unable to check if log was consumed")
		return
	}
	if was {
		return
	}
	_, err = rs.syncRegistry()
	if err != nil {
		rs.logger.With("error", err).Error("unable to sync registry")
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err := rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast); err != nil {
		rs.logger.With("error", err).Errorf("unable to mark SyncRegistryLog log as consumed, log: %v", broadcast.String())
	}
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
			rs.logger.Errorf("invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.handleUpkeepCancelled(broadcast)
	}
}

func (rs *RegistrySynchronizer) handleUpkeepCancelled(broadcast log.Broadcast) {
	txHash := broadcast.RawLog().TxHash.Hex()
	rs.logger.Debugw("processing UpkeepCanceled log", "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		rs.logger.With("error", err).Error("unable to check if log was consumed")
		return
	}
	if was {
		return
	}
	broadcastedLog, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepCanceled)
	if !ok {
		rs.logger.Errorf("invariant violation, expected UpkeepCanceled log but got %T", broadcastedLog)
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	affected, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, []int64{broadcastedLog.Id.Int64()})
	if err != nil {
		rs.logger.With("error", err).Error("unable to batch delete upkeeps")
		return
	}
	rs.logger.Debugw(fmt.Sprintf("deleted %v upkeep registrations", affected), "txHash", txHash)

	ctx, cancel = postgres.DefaultQueryCtx()
	defer cancel()
	if err := rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast); err != nil {
		rs.logger.With("error", err).Errorf("unable to mark KeeperRegistryUpkeepCanceled log as consumed,  log: %v", broadcast.String())
	}
}

func (rs *RegistrySynchronizer) handleUpkeepRegisteredLogs(done func()) {
	defer done()
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	registry, err := rs.orm.RegistryForJob(ctx, rs.job.ID)
	if err != nil {
		rs.logger.With("error", err).Error("unable to find registry for job")
		return
	}
	for {
		i, exists := rs.mailRoom.mbUpkeepRegistered.Retrieve()
		if !exists {
			return
		}
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			rs.logger.Errorf("invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.HandleUpkeepRegistered(broadcast, registry)
	}
}

func (rs *RegistrySynchronizer) HandleUpkeepRegistered(broadcast log.Broadcast, registry Registry) {
	txHash := broadcast.RawLog().TxHash.Hex()
	rs.logger.Debugw("processing UpkeepRegistered log", "txHash", txHash)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		rs.logger.With("error", err).Error("unable to check if log was consumed")
		return
	}
	if was {
		return
	}
	broadcastedLog, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepRegistered)
	if !ok {
		rs.logger.Errorf("invariant violation, expected UpkeepRegistered log but got %T", broadcastedLog)
		return
	}
	err = rs.syncUpkeep(registry, broadcastedLog.Id.Int64())
	if err != nil {
		rs.logger.With("error", err).Error("failed to sync upkeep, log: %v", broadcast.String())
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err := rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast); err != nil {
		rs.logger.With("error", err).Errorf("unable to mark KeeperRegistryUpkeepRegistered log as consumed, log: %v", broadcast.String())
	}
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
			rs.logger.Errorf("invariant violation, expected log.Broadcast but got %T", broadcast)
			continue
		}
		rs.handleUpkeepPerformed(broadcast)
	}
}

func (rs *RegistrySynchronizer) handleUpkeepPerformed(broadcast log.Broadcast) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	rawLog := broadcast.RawLog().TxHash.Hex()

	rs.logger.Debugw("processing UpkeepPerformed log", "txHash", rawLog)
	was, err := rs.logBroadcaster.WasAlreadyConsumed(rs.orm.DB, broadcast)
	if err != nil {
		rs.logger.With("error", err).Warn("unable to check if log was consumed")
		return
	}

	if was {
		return
	}

	log, ok := broadcast.DecodedLog().(*keeper_registry_wrapper.KeeperRegistryUpkeepPerformed)
	if !ok {
		rs.logger.Errorf("invariant violation, expected UpkeepPerformed log but got %T", log)
		return
	}

	if log.From.Hex() != rs.job.KeeperSpec.FromAddress.Hex() {
		return
	}

	// set last run to 0 so that keeper can resume checkUpkeep()
	err = rs.orm.SetLastRunHeightForUpkeepOnJob(ctx, rs.job.ID, log.Id.Int64(), 0)
	if err != nil {
		rs.logger.With("error", err).Error("failed to set last run to 0")
		return
	}

	ctx, cancel = postgres.DefaultQueryCtx()
	defer cancel()

	if err := rs.logBroadcaster.MarkConsumed(rs.orm.DB.WithContext(ctx), broadcast); err != nil {
		rs.logger.With("error", err).With("log", broadcast.String()).Error("unable to mark KeeperRegistryUpkeepPerformed log as consumed")
	}
}
