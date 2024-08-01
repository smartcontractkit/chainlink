package managed

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// runWithContractConfig runs fn with a contractConfig and manages its lifecycle
// as contractConfigs change according to contractConfigTracker. It also saves
// and restores contract configs using database.
func runWithContractConfig(
	ctx context.Context,

	contractConfigTracker types.ContractConfigTracker,
	database types.ConfigDatabase,
	fn func(context.Context, types.ContractConfig, loghelper.LoggerWithContext),
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	offchainConfigDigester types.OffchainConfigDigester,
) {
	rwcc := runWithContractConfigState{
		ctx,

		types.ConfigDigest{},
		contractConfigTracker,
		database,
		fn,
		localConfig,
		logger,

		prefixCheckConfigDigester{offchainConfigDigester},
		func() {},
		subprocesses.Subprocesses{},
		subprocesses.Subprocesses{},
	}
	rwcc.run()
}

type runWithContractConfigState struct {
	ctx context.Context

	configDigest          types.ConfigDigest
	contractConfigTracker types.ContractConfigTracker
	database              types.ConfigDatabase
	fn                    func(context.Context, types.ContractConfig, loghelper.LoggerWithContext)
	localConfig           types.LocalConfig
	logger                loghelper.LoggerWithContext

	configDigester prefixCheckConfigDigester
	fnCancel       context.CancelFunc
	fnSubs         subprocesses.Subprocesses
	otherSubs      subprocesses.Subprocesses
}

func (rwcc *runWithContractConfigState) run() {
	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	rwcc.restoreFromDatabase()

	// Only start tracking config after we attempted to load config from db
	chNewConfig := make(chan types.ContractConfig, 5)
	rwcc.otherSubs.Go(func() {
		TrackConfig(rwcc.ctx, rwcc.configDigester, rwcc.contractConfigTracker, rwcc.configDigest, rwcc.localConfig, rwcc.logger, chNewConfig)
	})

	for {
		select {
		case change := <-chNewConfig:
			rwcc.logger.Info("runWithContractConfig: switching between configs", commontypes.LogFields{
				"oldConfigDigest": rwcc.configDigest.Hex(),
				"newConfigDigest": change.ConfigDigest.Hex(),
			})
			rwcc.configChanged(change)
		case <-rwcc.ctx.Done():
			rwcc.logger.Info("runWithContractConfig: winding down", nil)
			rwcc.fnSubs.Wait()
			rwcc.otherSubs.Wait()
			rwcc.logger.Info("runWithContractConfig: exiting", nil)
			return // Exit managed event loop altogether
		}
	}
}

func (rwcc *runWithContractConfigState) restoreFromDatabase() {
	var contractConfig *types.ContractConfig
	ok := rwcc.otherSubs.BlockForAtMost(
		rwcc.ctx,
		rwcc.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			contractConfig = loadConfigFromDatabase(ctx, rwcc.database, rwcc.logger)
		},
	)
	if !ok {
		rwcc.logger.Error("runWithContractConfig: database timed out while attempting to restore configuration", commontypes.LogFields{
			"timeout": rwcc.localConfig.DatabaseTimeout,
		})
		return
	}

	if contractConfig == nil {
		rwcc.logger.Info("runWithContractConfig: found no configuration to restore", commontypes.LogFields{})
		return
	}

	rwcc.configChanged(*contractConfig)
}

func (rwcc *runWithContractConfigState) configChanged(contractConfig types.ContractConfig) {
	// Cease any operation from earlier configs
	rwcc.logger.Info("runWithContractConfig: winding down old configuration", commontypes.LogFields{
		"oldConfigDigest": rwcc.configDigest,
		"newConfigDigest": contractConfig.ConfigDigest,
	})
	rwcc.fnCancel()
	rwcc.fnSubs.Wait()
	rwcc.logger.Info("runWithContractConfig: closed old configuration", commontypes.LogFields{
		"oldConfigDigest": rwcc.configDigest,
		"newConfigDigest": contractConfig.ConfigDigest,
	})

	// note that there is an analogous check in TrackConfig, so this should never trigger.
	if err := rwcc.configDigester.CheckContractConfig(contractConfig); err != nil {
		rwcc.logger.Error("runWithContractConfig: detected corruption while attempting to change configuration", commontypes.LogFields{
			"err":            err,
			"contractConfig": contractConfig,
		})
		return
	}

	rwcc.configDigest = contractConfig.ConfigDigest

	fnCtx, fnCancel := context.WithCancel(rwcc.ctx)
	rwcc.fnCancel = fnCancel
	rwcc.fnSubs.Go(func() {
		defer fnCancel()
		rwcc.fn(
			fnCtx,
			contractConfig,
			rwcc.logger.MakeChild(commontypes.LogFields{"configDigest": contractConfig.ConfigDigest}),
		)
	})

	writeCtx, writeCancel := context.WithTimeout(rwcc.ctx, rwcc.localConfig.DatabaseTimeout)
	defer writeCancel()
	if err := rwcc.database.WriteConfig(writeCtx, contractConfig); err != nil {
		rwcc.logger.ErrorIfNotCanceled("runWithContractConfig: error writing new config to database", writeCtx, commontypes.LogFields{
			"configDigest": contractConfig.ConfigDigest,
			"config":       contractConfig,
			"error":        err,
		})
	}

}
